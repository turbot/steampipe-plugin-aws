package aws

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
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
	// d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

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

	// If session was not in cache - create a session and save to cache

	// get aws config info
	awsConfig := GetConfig(d.Connection)
	// ratelimiter := ratelimit.NewTokenRateLimit(500)
	retryer := retry.NewStandard(func(o *retry.StandardOptions) {
		o.MaxAttempts = maxRetries
		o.RateLimiter = NoOpRateLimit{}
		// o.RateLimiter = ratelimiter
		// backoff := retry.NewExponentialJitterBackoff(5 * time.Minute)
		backoff := NewExponentialJitterBackoff(5 * time.Minute)
		o.Backoff = backoff
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithRetryer(func() aws.Retryer {
			return retry.AddWithMaxBackoffDelay(retryer, minRetryDelay)
		}),
	)
	if err != nil {
		return nil, err
	}

	// handle custom endpoint URL, if any
	var awsEndpointUrl string

	awsEndpointUrl = os.Getenv("AWS_ENDPOINT_URL")
	if awsConfig.EndpointUrl != nil {
		awsEndpointUrl = *awsConfig.EndpointUrl
	}

	var customResolver aws.EndpointResolverWithOptionsFunc
	if awsEndpointUrl != "" {
		customResolver = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpointUrl,
				SigningRegion: region,
			}, nil
		})
	}

	if awsEndpointUrl != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(customResolver),
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxBackoffDelay(retryer, minRetryDelay)
			}),
		)
	}

	// awsConfig.S3ForcePathStyle - Moved to service specific client (i.e. in S3V2Client)

	if awsConfig.Profile != nil {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithSharedConfigProfile(*awsConfig.Profile),
			config.WithRegion(region),
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxBackoffDelay(retryer, minRetryDelay)
			},
			),
		)
	}

	if awsConfig.AccessKey != nil && awsConfig.SecretKey == nil {
		return nil, fmt.Errorf("Partial credentials found in connection config, missing: secret_key")
	} else if awsConfig.SecretKey != nil && awsConfig.AccessKey == nil {
		return nil, fmt.Errorf("Partial credentials found in connection config, missing: access_key")
	} else if awsConfig.AccessKey != nil && awsConfig.SecretKey != nil {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			*awsConfig.AccessKey, *awsConfig.SecretKey, "",
		)),
			config.WithRegion(region),
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxBackoffDelay(retryer, minRetryDelay)
			}))

		if awsConfig.SessionToken != nil {
			cfg, err = config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				*awsConfig.AccessKey, *awsConfig.SecretKey, *awsConfig.SessionToken,
			)),
				config.WithRegion(region),
				config.WithRetryer(func() aws.Retryer {
					return retry.AddWithMaxBackoffDelay(retryer, minRetryDelay)
				}))
		}
	}

	return &cfg, err
}

// ExponentialJitterBackoff provides backoff delays with jitter based on the
// number of attempts.
type ExponentialJitterBackoff struct {
	maxBackoff time.Duration
	// precomputed number of attempts needed to reach max backoff.
	maxBackoffAttempts float64

	randFloat64 func() (float64, error)
}

// NewExponentialJitterBackoff returns an ExponentialJitterBackoff configured
// for the max backoff.
func NewExponentialJitterBackoff(maxBackoff time.Duration) *ExponentialJitterBackoff {
	return &ExponentialJitterBackoff{
		maxBackoff: maxBackoff,
		maxBackoffAttempts: math.Log2(
			float64(maxBackoff) / float64(time.Second)),
		randFloat64: CryptoRandFloat64,
	}
}

// BackoffDelay returns the duration to wait before the next attempt should be
// made. Returns an error if unable get a duration.
func (j *ExponentialJitterBackoff) BackoffDelay(attempt int, err error) (time.Duration, error) {
	log.Printf("[WARN] ***************** attempt: %d\n", attempt)
	if attempt > int(j.maxBackoffAttempts) {
		return j.maxBackoff, nil
	}

	b, err := j.randFloat64()
	if err != nil {
		return 0, err
	}

	// [0.0, 1.0) * 2 ^ attempts
	ri := int64(1 << uint64(attempt))
	delaySeconds := b * float64(ri)

	delay := FloatSecondsDur(delaySeconds)
	log.Printf("[WARN] ***************** delay: %v\n", delay)
	return delay, nil
}

func CryptoRandFloat64() (float64, error) {
	return Float64(Reader)
}

// Float64 returns a float64 read from an io.Reader source. The returned float will be between [0.0, 1.0).
func Float64(reader io.Reader) (float64, error) {
	bi, err := rand.Int(reader, floatMaxBigInt)
	if err != nil {
		return 0, fmt.Errorf("failed to read random value, %v", err)
	}

	return float64(bi.Int64()) / (1 << 53), nil
}

var Reader io.Reader

var floatMaxBigInt = big.NewInt(1 << 53)

// FloatSecondsDur converts a fractional seconds to duration.
func FloatSecondsDur(v float64) time.Duration {
	return time.Duration(v * float64(time.Second))
}
