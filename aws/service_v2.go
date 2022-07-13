package aws

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

// SNSV2Service returns the service connection for AWS SNS service
func SNSV2Client(ctx context.Context, d *plugin.QueryData) (*sns.Client, error) {
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

// type ConnectionErrRetryerV2 struct {
// 	retry.Standard
// 	ctx context.Context
// }

// func NewConnectionErrRetryerV2(maxRetries int, minRetryDelay time.Duration, ctx context.Context) *ConnectionErrRetryerV2 {
// 	rand.Seed(time.Now().UnixNano()) // reseting state of rand to generate different random values
// 	return &ConnectionErrRetryer{
// 		ctx: ctx,
// 		DefaultRetryer: client.DefaultRetryer{
// 			NumMaxRetries: maxRetries,    // MUST be set or all retrying is skipped!
// 			MinRetryDelay: minRetryDelay, // Set minimum retry delay
// 		},
// 	}
// }

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
	sessionCacheKey := fmt.Sprintf("session-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*aws.Config), nil
	}

	// If session was not in cache - create a session and save to cache

	// get aws config info
	awsConfig := GetConfig(d.Connection)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithRetryMaxAttempts(maxRetries),
		// config.WithRetryer(NewConnectionErrRetryerV2(maxRetries, minRetryDelay, ctx)), TODO
	)
	if err != nil {
		return nil, err
	}

	// session default configuration
	// sessionOptions := session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// 	Config: aws.Config{
	// 		Region:     &region,
	// 		MaxRetries: aws.Int(maxRetries),
	// 		Retryer:    NewConnectionErrRetryerV2(maxRetries, minRetryDelay, ctx),
	// 	},
	// }

	// handle custom endpoint URL, if any
	// var awsEndpointUrl string

	// awsEndpointUrl = os.Getenv("AWS_ENDPOINT_URL")

	// if awsConfig.EndpointUrl != nil {
	// 	config.LoadOptions.EndpointResolverWithOptions
	// 	awsEndpointUrl = *awsConfig.EndpointUrl
	// }

	// if awsEndpointUrl != "" {
	// 	sessionOptions.Config.Endpoint = aws.String(awsEndpointUrl)
	// }

	// if awsConfig.S3ForcePathStyle != nil {
	// 	sessionOptions.Config.S3ForcePathStyle = awsConfig.S3ForcePathStyle
	// }

	if awsConfig.Profile != nil {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithSharedConfigProfile(*awsConfig.Profile),
		)
	}

	if awsConfig.AccessKey != nil && awsConfig.SecretKey == nil {
		return nil, fmt.Errorf("Partial credentials found in connection config, missing: secret_key")
	} else if awsConfig.SecretKey != nil && awsConfig.AccessKey == nil {
		return nil, fmt.Errorf("Partial credentials found in connection config, missing: access_key")
	} else if awsConfig.AccessKey != nil && awsConfig.SecretKey != nil {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			*awsConfig.AccessKey, *awsConfig.SecretKey, "",
		)))
		// sessionOptions.Config.Credentials = credentials.NewStaticCredentials(
		// 	*awsConfig.AccessKey, *awsConfig.SecretKey, "",
		// )

		if awsConfig.SessionToken != nil {
			cfg, err = config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				*awsConfig.AccessKey, *awsConfig.SecretKey, *awsConfig.SessionToken,
			)))
		}
	}

	return &cfg, err
}
