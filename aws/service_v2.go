package aws

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

// https://github.com/aws/aws-sdk-go-v2/issues/543
type NoOpRateLimit struct{}

func (NoOpRateLimit) AddTokens(uint) error { return nil }
func (NoOpRateLimit) GetToken(context.Context, uint) (func() error, error) {
	return noOpToken, nil
}
func noOpToken() error { return nil }

// IAMClient returns the service client for AWS IAM service
func IAMClient(ctx context.Context, d *plugin.QueryData) (*iam.Client, error) {
	// have we already created and cached the service?
	serviceCacheKey := "iam-v2"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*iam.Client), nil
	}
	// so it was not in cache - create service
	cfg, err := getSessionV2(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}

	svc := iam.NewFromConfig(*cfg)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// S3Client returns the service client for AWS S3 service
func S3Client(ctx context.Context, d *plugin.QueryData, region string) (*s3.Client, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed S3Service")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("s3-v2-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*s3.Client), nil
	}

	awsConfig := GetConfig(d.Connection)

	// so it was not in cache - create service
	cfg, err := getSessionV2(ctx, d, region)
	if err != nil {
		return nil, err
	}

	var svc *s3.Client

	if awsConfig.S3ForcePathStyle != nil {
		svc = s3.NewFromConfig(*cfg, func(o *s3.Options) {
			o.UsePathStyle = *awsConfig.S3ForcePathStyle
		})
	} else {
		svc = s3.NewFromConfig(*cfg)
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// SNSClient returns the service client for AWS SNS service
func SNSClient(ctx context.Context, d *plugin.QueryData) (*sns.Client, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	if region == "" {
		return nil, fmt.Errorf("region must be passed SNSV2Service")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("sns-v2-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*sns.Client), nil
	}

	// so it was not in cache - create service
	cfg, err := getSessionV2(ctx, d, region)
	if err != nil {
		return nil, err
	}

	svc := sns.NewFromConfig(*cfg)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

func getSessionV2(ctx context.Context, d *plugin.QueryData, region string) (*aws.Config, error) {
	awsConfig := GetConfig(d.Connection)

	// As per the logic used in retryRules of NewConnectionErrRetryer, default to minimum delay of 25ms and maximum
	// number of retries as 9 (our default). The default maximum delay will not be more than approximately 3 minutes to avoid Steampipe
	// waiting too long to return results
	maxRetries := 9
	var minRetryDelay time.Duration = 25 * time.Millisecond // Default minimum delay

	// Set max retry count from config file or env variable (config file has precedence)
	if awsConfig.MaxErrorRetryAttempts != nil {
		maxRetries = *awsConfig.MaxErrorRetryAttempts
	} else if os.Getenv("AWS_MAX_ATTEMPTS") != "" {
		maxRetriesEnvVar, err := strconv.Atoi(os.Getenv("AWS_MAX_ATTEMPTS"))
		if err != nil || maxRetriesEnvVar < 1 {
			panic("invalid value for environment variable \"AWS_MAX_ATTEMPTS\". It should be an integer value greater than or equal to 1")
		}
		maxRetries = maxRetriesEnvVar
	}

	// Set min delay time from config file
	if awsConfig.MinErrorRetryDelay != nil {
		minRetryDelay = time.Duration(*awsConfig.MinErrorRetryDelay) * time.Millisecond
	}

	if maxRetries < 1 {
		panic("\nconnection config has invalid value for \"max_error_retry_attempts\", it must be greater than or equal to 1. Edit your connection configuration file and then restart Steampipe.")
	}
	if minRetryDelay < 1 {
		panic("\nconnection config has invalid value for \"min_error_retry_delay\", it must be greater than or equal to 1. Edit your connection configuration file and then restart Steampipe.")
	}

	return getSessionV2WithMaxRetries(ctx, d, region, maxRetries, minRetryDelay)
}

func getSessionV2WithMaxRetries(ctx context.Context, d *plugin.QueryData, region string, maxRetries int, minRetryDelay time.Duration) (*aws.Config, error) {
	sessionCacheKey := fmt.Sprintf("session-v2-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*aws.Config), nil
	}

	retryer := retry.NewStandard(func(o *retry.StandardOptions) {
		// reseting state of rand to generate different random values
		rand.Seed(time.Now().UnixNano())
		o.MaxAttempts = maxRetries
		o.MaxBackoff = 5 * time.Minute
		o.RateLimiter = NoOpRateLimit{} // With no rate limiter
		o.Backoff = NewExponentialJitterBackoff(minRetryDelay, maxRetries)
	})

	awsConfig := GetConfig(d.Connection)
	configOptions := []func(*config.LoadOptions) error{
		config.WithRegion(region),
		config.WithRetryer(func() aws.Retryer {
			return retryer
		}),
	}

	// handle custom endpoint URL, if any
	var awsEndpointUrl string

	awsEndpointUrl = os.Getenv("AWS_ENDPOINT_URL")
	if awsConfig.EndpointUrl != nil {
		awsEndpointUrl = *awsConfig.EndpointUrl
	}

	if awsEndpointUrl != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpointUrl,
				SigningRegion: region,
			}, nil
		})

		configOptions = append(configOptions, config.WithEndpointResolverWithOptions(customResolver))
	}

	// awsConfig.S3ForcePathStyle - Moved to service specific client (i.e. in S3V2Client)

	if awsConfig.Profile != nil {
		configOptions = append(configOptions, config.WithSharedConfigProfile(aws.ToString(awsConfig.Profile)))
	}

	if awsConfig.AccessKey != nil && awsConfig.SecretKey == nil {
		return nil, fmt.Errorf("Partial credentials found in connection config, missing: secret_key")
	} else if awsConfig.SecretKey != nil && awsConfig.AccessKey == nil {
		return nil, fmt.Errorf("Partial credentials found in connection config, missing: access_key")
	} else if awsConfig.AccessKey != nil && awsConfig.SecretKey != nil {
		var provider credentials.StaticCredentialsProvider

		if awsConfig.SessionToken != nil {
			provider = credentials.NewStaticCredentialsProvider(*awsConfig.AccessKey, *awsConfig.SecretKey, *awsConfig.SessionToken)
		} else {
			provider = credentials.NewStaticCredentialsProvider(*awsConfig.AccessKey, *awsConfig.SecretKey, "")
		}
		configOptions = append(configOptions, config.WithCredentialsProvider(provider))
	}

	cfg, err := config.LoadDefaultConfig(ctx, configOptions...)
	if err != nil {
		plugin.Logger(ctx).Error("getAwsConfigWithMaxRetries", "load_default_config", err)
		return nil, err
	}
	d.ConnectionManager.Cache.Set(sessionCacheKey, &cfg)

	return &cfg, err
}

// ExponentialJitterBackoff provides backoff delays with jitter based on the
// number of attempts.
type ExponentialJitterBackoff struct {
	minDelay           time.Duration
	maxBackoffAttempts int
}

// NewExponentialJitterBackoff returns an ExponentialJitterBackoff configured
// for the max backoff.
func NewExponentialJitterBackoff(minDelay time.Duration, maxAttempts int) *ExponentialJitterBackoff {
	return &ExponentialJitterBackoff{minDelay, maxAttempts}
}

// BackoffDelay returns the duration to wait before the next attempt should be
// made. Returns an error if unable get a duration.
func (j *ExponentialJitterBackoff) BackoffDelay(attempt int, err error) (time.Duration, error) {
	minDelay := j.minDelay

	// The calculatted jitter will be between [0.8, 1.2)
	var jitter = float64(rand.Intn(120-80)+80) / 100

	retryTime := time.Duration(int(float64(int(minDelay.Nanoseconds())*int(math.Pow(3, float64(attempt)))) * jitter))
	// log.Printf("[INFO] *******INSIDE CODE********** Attempt: %d, DelayInSeconds: %f, Delay: %v", attempt, retryTime.Seconds(), retryTime)

	// Cap retry time at 5 minutes to avoid too long a wait
	if retryTime > time.Duration(5*time.Minute) {
		retryTime = time.Duration(5 * time.Minute)
	}

	return retryTime, nil
}
