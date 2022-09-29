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
	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/appconfig"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/dax"
	"github.com/aws/aws-sdk-go-v2/service/docdb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go/aws/endpoints"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

// https://github.com/aws/aws-sdk-go-v2/issues/543
type NoOpRateLimit struct{}

func (NoOpRateLimit) AddTokens(uint) error { return nil }
func (NoOpRateLimit) GetToken(context.Context, uint) (func() error, error) {
	return noOpToken, nil
}
func noOpToken() error { return nil }

func AccountClient(ctx context.Context, d *plugin.QueryData) (*account.Client, error) {
	cfg, err := getClient(ctx, d, "Account")
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return account.NewFromConfig(*cfg), nil
}

func ACMClient(ctx context.Context, d *plugin.QueryData) (*acm.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return acm.NewFromConfig(*cfg), nil
}

func APIGatewayClient(ctx context.Context, d *plugin.QueryData) (*apigateway.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return apigateway.NewFromConfig(*cfg), nil
}

func APIGatewayV2Client(ctx context.Context, d *plugin.QueryData) (*apigatewayv2.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return apigatewayv2.NewFromConfig(*cfg), nil
}

func AppConfigClient(ctx context.Context, d *plugin.QueryData) (*appconfig.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return appconfig.NewFromConfig(*cfg), nil
}

func AutoScalingClient(ctx context.Context, d *plugin.QueryData) (*autoscaling.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return autoscaling.NewFromConfig(*cfg), nil
}

func CodeDeployClient(ctx context.Context, d *plugin.QueryData) (*codedeploy.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return codedeploy.NewFromConfig(*cfg), nil
}

func CodeArtifactClient(ctx context.Context, d *plugin.QueryData) (*codeartifact.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, "codeartifact")
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return codeartifact.NewFromConfig(*cfg), nil
}

// CostExplorerClient returns the connection client for AWS Cost Explorer service
func CostExplorerClient(ctx context.Context, d *plugin.QueryData) (*costexplorer.Client, error) {
	cfg, err := getClient(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}

	return costexplorer.NewFromConfig(*cfg), nil
}

func DaxClient(ctx context.Context, d *plugin.QueryData) (*dax.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, endpoints.DaxServiceID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return dax.NewFromConfig(*cfg), nil
}

func DocDBClient(ctx context.Context, d *plugin.QueryData) (*docdb.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return docdb.NewFromConfig(*cfg), nil
}

func DynamoDbClient(ctx context.Context, d *plugin.QueryData) (*dynamodb.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(*cfg), nil
}

func EC2Client(ctx context.Context, d *plugin.QueryData) (*ec2.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return ec2.NewFromConfig(*cfg), nil
}

func ElasticBeanstalkClient(ctx context.Context, d *plugin.QueryData) (*elasticbeanstalk.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return elasticbeanstalk.NewFromConfig(*cfg), nil
}

func ELBClient(ctx context.Context, d *plugin.QueryData) (*elb.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return elb.NewFromConfig(*cfg), nil
}

func ELBV2Client(ctx context.Context, d *plugin.QueryData) (*elbv2.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return elbv2.NewFromConfig(*cfg), nil
}

func IAMClient(ctx context.Context, d *plugin.QueryData) (*iam.Client, error) {
	cfg, err := getClient(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return iam.NewFromConfig(*cfg), nil
}

func KafkaClient(ctx context.Context, d *plugin.QueryData) (*kafka.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, "kafka")
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return kafka.NewFromConfig(*cfg), nil
}

func RedshiftClient(ctx context.Context, d *plugin.QueryData) (*redshift.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return redshift.NewFromConfig(*cfg), nil
}

func RedshiftServerlessClient(ctx context.Context, d *plugin.QueryData) (*redshiftserverless.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, "redshift-serverless")
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return redshiftserverless.NewFromConfig(*cfg), nil
}

func S3Client(ctx context.Context, d *plugin.QueryData, region string) (*s3.Client, error) {
	cfg, err := getClientForRegion(ctx, d, region)
	if err != nil {
		return nil, err
	}

	var svc *s3.Client

	awsConfig := GetConfig(d.Connection)
	if awsConfig.S3ForcePathStyle != nil {
		svc = s3.NewFromConfig(*cfg, func(o *s3.Options) {
			o.UsePathStyle = *awsConfig.S3ForcePathStyle
		})
	} else {
		svc = s3.NewFromConfig(*cfg)
	}

	return svc, nil
}

func S3ControlClient(ctx context.Context, d *plugin.QueryData, region string) (*s3control.Client, error) {
	cfg, err := getClient(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return s3control.NewFromConfig(*cfg), nil
}

func SNSClient(ctx context.Context, d *plugin.QueryData) (*sns.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return sns.NewFromConfig(*cfg), nil
}

func getClient(ctx context.Context, d *plugin.QueryData, region string) (*aws.Config, error) {

	sessionCacheKey := fmt.Sprintf("session-v2-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*aws.Config), nil
	}

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

	sess, err := getClientWithMaxRetries(ctx, d, region, maxRetries, minRetryDelay)
	if err != nil {
		plugin.Logger(ctx).Error("getService.getClientWithMaxRetries", "region", region, "err", err)
	} else {
		// Caching sessions saves about 10ms, which is significant when there are
		// multiple instantiations (per account region) and when doing queries that
		// often take <100ms total. But, it's not that important compared to having
		// fresh credentials all the time. So, set a short cache length to ensure
		// we don't get tripped up by credential rotation on short lived roles etc.
		// The minimum assume role time is 15 minutes, so 5 minutes feels like a
		// reasonable balance - I certainly wouldn't do longer.
		d.ConnectionManager.Cache.SetWithTTL(sessionCacheKey, sess, 5*time.Minute)
	}

	return sess, err
}

// Get a session for the region defined in query data, but only after checking it's
// a supported region for the given serviceID.
func getClientForQuerySupportedRegion(ctx context.Context, d *plugin.QueryData, serviceID string) (*aws.Config, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	if region == "" {
		return nil, fmt.Errorf("getSessionForQueryRegion called without a region in QueryData")
	}
	validRegions := SupportedRegionsForService(ctx, d, serviceID)
	if !helpers.StringSliceContains(validRegions, region) {
		// We choose to ignore unsupported regions rather than returning an error
		// for them - it's a better user experience. So, return a nil session rather
		// than an error. The caller must handle this case.
		return nil, nil
	}
	// Supported region, so get and return the session
	return getClient(ctx, d, region)
}

// Helper function to get the session for a region set in query data
func getClientForQueryRegion(ctx context.Context, d *plugin.QueryData) (*aws.Config, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	if region == "" {
		return nil, fmt.Errorf("getSessionForQueryRegion called without a region in QueryData")
	}
	return getClient(ctx, d, region)
}

// Helper function to get the session for a specific region
func getClientForRegion(ctx context.Context, d *plugin.QueryData, region string) (*aws.Config, error) {
	if region == "" {
		return nil, fmt.Errorf("getSessionForRegion called with an empty region")
	}
	return getClient(ctx, d, region)
}

func getClientWithMaxRetries(ctx context.Context, d *plugin.QueryData, region string, maxRetries int, minRetryDelay time.Duration) (*aws.Config, error) {

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

	// Cap retry time at 5 minutes to avoid too long a wait
	if retryTime > time.Duration(5*time.Minute) {
		retryTime = time.Duration(5 * time.Minute)
	}

	return retryTime, nil
}
