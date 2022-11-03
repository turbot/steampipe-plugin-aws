package aws

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/backup"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func BackupService(ctx context.Context, d *plugin.QueryData) (*backup.Backup, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, backup.EndpointsID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	return backup.New(sess), nil
}

func getSession(ctx context.Context, d *plugin.QueryData, region string) (*session.Session, error) {

	sessionCacheKey := fmt.Sprintf("session-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*session.Session), nil
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

	sess, err := getSessionWithMaxRetries(ctx, d, region, maxRetries, minRetryDelay)
	if err != nil {
		plugin.Logger(ctx).Error("getClient.getSessionWithMaxRetries", "region", region, "err", err)
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

func getSessionWithMaxRetries(ctx context.Context, d *plugin.QueryData, region string, maxRetries int, minRetryDelay time.Duration) (*session.Session, error) {

	// get aws config info
	awsConfig := GetConfig(d.Connection)

	// session default configuration
	sessionOptions := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region:     &region,
			MaxRetries: aws.Int(maxRetries),
			Retryer:    NewConnectionErrRetryer(maxRetries, minRetryDelay, ctx),
		},
	}

	// handle custom endpoint URL, if any
	var awsEndpointUrl string

	awsEndpointUrl = os.Getenv("AWS_ENDPOINT_URL")

	if awsConfig.EndpointUrl != nil {
		awsEndpointUrl = *awsConfig.EndpointUrl
	}

	if awsEndpointUrl != "" {
		sessionOptions.Config.Endpoint = aws.String(awsEndpointUrl)
	}

	if awsConfig.S3ForcePathStyle != nil {
		sessionOptions.Config.S3ForcePathStyle = awsConfig.S3ForcePathStyle
	}

	if awsConfig.Profile != nil {
		sessionOptions.Profile = *awsConfig.Profile
	}

	if awsConfig.AccessKey != nil && awsConfig.SecretKey == nil {
		return nil, fmt.Errorf("Partial credentials found in connection config, missing: secret_key")
	} else if awsConfig.SecretKey != nil && awsConfig.AccessKey == nil {
		return nil, fmt.Errorf("Partial credentials found in connection config, missing: access_key")
	} else if awsConfig.AccessKey != nil && awsConfig.SecretKey != nil {

		sessionOptions.Config.Credentials = credentials.NewStaticCredentials(
			*awsConfig.AccessKey, *awsConfig.SecretKey, "",
		)

		if awsConfig.SessionToken != nil {
			sessionOptions.Config.Credentials = credentials.NewStaticCredentials(
				*awsConfig.AccessKey, *awsConfig.SecretKey, *awsConfig.SessionToken,
			)
		}
	}

	sess, err := session.NewSessionWithOptions(sessionOptions)
	if err != nil {
		plugin.Logger(ctx).Error("getSessionWithMaxRetries.NewSessionWithOptions", "sessionOptions", sessionOptions, "err", err)
		return nil, err
	}

	return sess, nil
}

// Get a session for the region defined in query data, but only after checking it's
// a supported region for the given serviceID.
func getSessionForQuerySupportedRegion(ctx context.Context, d *plugin.QueryData, serviceID string) (*session.Session, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	if region == "" {
		return nil, fmt.Errorf("getSessionForQuerySupportedRegion called without a region in QueryData")
	}
	validRegions, err := GetSupportedRegionsForClient(ctx, d, serviceID)
	if err != nil {
		return nil, err
	}
	if !helpers.StringSliceContains(validRegions, region) {
		// We choose to ignore unsupported regions rather than returning an error
		// for them - it's a better user experience. So, return a nil session rather
		// than an error. The caller must handle this case.
		return nil, nil
	}
	// Supported region, so get and return the session
	return getSession(ctx, d, region)
}

// GetDefaultAwsRegion returns the default region for AWS partiton
// if not set by Env variable or in aws profile
func GetDefaultAwsRegion(d *plugin.QueryData) string {
	allAwsRegions := getAllAwsRegions()

	// have we already created and cached the service?
	serviceCacheKey := "GetDefaultAwsRegion"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(string)
	}

	// get aws config info
	awsConfig := GetConfig(d.Connection)

	var regions []string
	var region string

	if awsConfig.Regions != nil {
		regions = awsConfig.Regions
		region = regions[0]
	} else {
		session, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})
		if err != nil {
			panic(err)
		}
		if session != nil && session.Config != nil {
			region = *session.Config.Region
		}

		if region != "" {
			regions = []string{region}
		}
	}

	invalidPatterns := []string{}
	for _, namePattern := range regions {
		validRegions := []string{}
		for _, validRegion := range allAwsRegions {
			if ok, _ := path.Match(namePattern, validRegion); ok {
				validRegions = append(validRegions, validRegion)
			}
		}

		// Region items with wildcards that match on 0 regions should not be
		// considered invalid
		if len(validRegions) == 0 && !strings.ContainsAny(namePattern, "?*") {
			invalidPatterns = append(invalidPatterns, namePattern)
		}
	}

	if len(invalidPatterns) > 0 {
		panic("\nconnection config has invalid \"regions\": " + strings.Join(invalidPatterns, ", ") + ". Edit your connection configuration file and then restart Steampipe.")
	}

	// most of the global services like IAM, S3, Route 53, etc. in all cloud types target these regions
	if strings.HasPrefix(region, "us-gov") && !helpers.StringSliceContains(allAwsRegions, region) {
		region = "us-gov-west-1"
	} else if strings.HasPrefix(region, "cn") && !helpers.StringSliceContains(allAwsRegions, region) {
		region = "cn-northwest-1"
	} else if strings.HasPrefix(region, "us-isob") && !helpers.StringSliceContains(allAwsRegions, region) {
		region = "us-isob-east-1"
	} else if strings.HasPrefix(region, "us-iso") && !helpers.StringSliceContains(allAwsRegions, region) {
		region = "us-iso-east-1"
	} else if !helpers.StringSliceContains(allAwsRegions, region) {
		region = "us-east-1"
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, region)
	return region
}

// Function from https://github.com/panther-labs/panther/blob/v1.16.0/pkg/awsretry/connection_retryer.go
func NewConnectionErrRetryer(maxRetries int, minRetryDelay time.Duration, ctx context.Context) *ConnectionErrRetryer {
	rand.Seed(time.Now().UnixNano()) // reseting state of rand to generate different random values
	return &ConnectionErrRetryer{
		ctx: ctx,
		DefaultRetryer: client.DefaultRetryer{
			NumMaxRetries: maxRetries,    // MUST be set or all retrying is skipped!
			MinRetryDelay: minRetryDelay, // Set minimum retry delay
		},
	}
}

// ConnectionErrRetryer wraps the SDK's built in DefaultRetryer adding customization
// to retry `connection reset by peer` errors.
// Note: This retryer should be used for either idempotent operations or for operations
// where performing duplicate requests to AWS is acceptable.
// See also: https://github.com/aws/aws-sdk-go/issues/3027#issuecomment-567269161
type ConnectionErrRetryer struct {
	client.DefaultRetryer
	ctx context.Context
}

func (r ConnectionErrRetryer) ShouldRetry(req *request.Request) bool {
	if req.Error != nil {
		if strings.Contains(req.Error.Error(), "connection reset by peer") {
			return true
		}

		var awsErr awserr.Error
		if errors.As(req.Error, &awsErr) {
			/*
				If no credentials are set or an invalid profile is provided, the AWS SDK
				will attempt to authenticate using all known methods. This takes a while
				since it will attempt to reach the EC2 metadata service and will continue
				to retry on connection errors, e.g.,
				awsErr.OrigErr()="Put "http://169.254.169.254/latest/api/token": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
				awsErr.OrigErr()="Get "http://169.254.169.254/latest/meta-data/iam/security-credentials/": dial tcp 169.254.169.254:80: connect: no route to host"
				To reduce the time to fail, limit the number of retries for these errors specifically.
			*/
			if awsErr.OrigErr() != nil {
				if strings.Contains(awsErr.OrigErr().Error(), "http://169.254.169.254/latest") && req.RetryCount > 3 {
					return false
				}
			}
		}
	}

	// Fallback to SDK's built in retry rules
	return r.DefaultRetryer.ShouldRetry(req)
}

// Customize the RetryRules to implement exponential backoff retry
func (d ConnectionErrRetryer) RetryRules(r *request.Request) time.Duration {
	retryCount := r.RetryCount
	minDelay := d.MinRetryDelay

	// If errors are caused by load, retries can be ineffective if all API request retry at the same time.
	// To avoid this problem added a jitter of "+/-20%" with delay time.
	// For example, if the delay is 25ms, the final delay could be between 20 and 30ms.
	var jitter = float64(rand.Intn(120-80)+80) / 100

	// Creates a new exponential backoff using the starting value of
	// minDelay and (minDelay * 3^retrycount) * jitter on each failure
	// For example, with a min delay time of 25ms: 23.25ms, 63ms, 238.5ms, 607.4ms, 2s, 5.22s, 20.31s..., up to max.
	retryTime := time.Duration(int(float64(int(minDelay.Nanoseconds())*int(math.Pow(3, float64(retryCount)))) * jitter))

	// Cap retry time at 5 minutes to avoid too long a wait
	if retryTime > time.Duration(5*time.Minute) {
		retryTime = time.Duration(5 * time.Minute)
	}

	return retryTime
}
