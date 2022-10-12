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
	"github.com/aws/aws-sdk-go-v2/service/accessanalyzer"
	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/appconfig"
	"github.com/aws/aws-sdk-go-v2/service/auditmanager"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"
	"github.com/aws/aws-sdk-go-v2/service/dax"
	"github.com/aws/aws-sdk-go-v2/service/directoryservice"
	"github.com/aws/aws-sdk-go-v2/service/dlm"
	"github.com/aws/aws-sdk-go-v2/service/docdb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
  "github.com/aws/aws-sdk-go-v2/service/fsx"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
  "github.com/aws/aws-sdk-go-v2/service/inspector"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/ram"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/waf"
	"github.com/aws/aws-sdk-go/aws/endpoints"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"

	fsxEndpoint "github.com/aws/aws-sdk-go/service/fsx"
	lambdaEndpoint "github.com/aws/aws-sdk-go/service/lambda"
)

// https://github.com/aws/aws-sdk-go-v2/issues/543
type NoOpRateLimit struct{}

func (NoOpRateLimit) AddTokens(uint) error { return nil }
func (NoOpRateLimit) GetToken(context.Context, uint) (func() error, error) {
	return noOpToken, nil
}
func noOpToken() error { return nil }

// AccessAnalyzerClient returns the service connection for AWS IAM Access Analyzer service
func AccessAnalyzerClient(ctx context.Context, d *plugin.QueryData) (*accessanalyzer.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return accessanalyzer.NewFromConfig(*cfg), nil
}

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

func AuditManagerClient(ctx context.Context, d *plugin.QueryData) (*auditmanager.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, "auditmanager")
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return auditmanager.NewFromConfig(*cfg), nil
}

func AutoScalingClient(ctx context.Context, d *plugin.QueryData) (*autoscaling.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return autoscaling.NewFromConfig(*cfg), nil
}

func BackupClient(ctx context.Context, d *plugin.QueryData) (*backup.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return backup.NewFromConfig(*cfg), nil
}

func CloudControlClient(ctx context.Context, d *plugin.QueryData) (*cloudcontrol.Client, error) {
	// CloudControl returns GeneralServiceException in a lot of situations, which
	// AWS SDK treats as retryable. This is frustrating because we end up retrying
	// many times for things that will never work.
	// So, we use a specific client configuration for CloudControl with a smaller
	// number of retries to avoid hangs. In effect, this service IGNORES the retry
	// configuration in aws.spc - but, good enough for something that is rarely used
	// anyway.
	region := d.KeyColumnQualString(matrixKeyRegion)
	if region == "" {
		return nil, fmt.Errorf("CloudControlService called without a region in QueryData")
	}

	// Use a service level cache since we are going around the standard
	// getSession with its caching.
	serviceCacheKey := fmt.Sprintf("cloudcontrol-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudcontrol.Client), nil
	}

	cfg, err := getClientWithMaxRetries(ctx, d, region, 4, 25*time.Millisecond)
	if err != nil {
		return nil, err
	}
	svc := cloudcontrol.NewFromConfig(*cfg)

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

func CloudFrontClient(ctx context.Context, d *plugin.QueryData) (*cloudfront.Client, error) {
	cfg, err := getClient(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return cloudfront.NewFromConfig(*cfg), nil
}

func CloudTrailClient(ctx context.Context, d *plugin.QueryData) (*cloudtrail.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return cloudtrail.NewFromConfig(*cfg), nil
}

func CloudTrailRegionsClient(ctx context.Context, d *plugin.QueryData, region string) (*cloudtrail.Client, error) {
	cfg, err := getClient(ctx, d, region)
	if err != nil {
		return nil, err
	}
	return cloudtrail.NewFromConfig(*cfg), nil
}

func CloudWatchClient(ctx context.Context, d *plugin.QueryData) (*cloudwatch.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return cloudwatch.NewFromConfig(*cfg), nil
}

func CloudWatchLogsClient(ctx context.Context, d *plugin.QueryData) (*cloudwatchlogs.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return cloudwatchlogs.NewFromConfig(*cfg), nil
}

func CloudFormationClient(ctx context.Context, d *plugin.QueryData) (*cloudformation.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return cloudformation.NewFromConfig(*cfg), nil
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

func CodeBuildClient(ctx context.Context, d *plugin.QueryData) (*codebuild.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, endpoints.CodebuildServiceID)
	if err != nil {
		return nil, err
	}
	return codebuild.NewFromConfig(*cfg), nil
}

func CodeDeployClient(ctx context.Context, d *plugin.QueryData) (*codedeploy.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return codedeploy.NewFromConfig(*cfg), nil
}

// CodePipelineClient returns the service connection for AWS CodePipeline service
func CodePipelineClient(ctx context.Context, d *plugin.QueryData) (*codepipeline.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return codepipeline.NewFromConfig(*cfg), nil
}

func ConfigClient(ctx context.Context, d *plugin.QueryData) (*configservice.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return configservice.NewFromConfig(*cfg), nil
}

// CostExplorerClient returns the connection client for AWS Cost Explorer service
func CostExplorerClient(ctx context.Context, d *plugin.QueryData) (*costexplorer.Client, error) {
	cfg, err := getClient(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return costexplorer.NewFromConfig(*cfg), nil
}

func DatabaseMigrationClient(ctx context.Context, d *plugin.QueryData) (*databasemigrationservice.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return databasemigrationservice.NewFromConfig(*cfg), nil
}

func DAXClient(ctx context.Context, d *plugin.QueryData) (*dax.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return dax.NewFromConfig(*cfg), nil
}

func DirectoryServiceClient(ctx context.Context, d *plugin.QueryData) (*directoryservice.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return directoryservice.NewFromConfig(*cfg), nil
}

func DLMClient(ctx context.Context, d *plugin.QueryData) (*dlm.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return dlm.NewFromConfig(*cfg), nil
}

func DocDBClient(ctx context.Context, d *plugin.QueryData) (*docdb.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return docdb.NewFromConfig(*cfg), nil
}

func DynamoDBClient(ctx context.Context, d *plugin.QueryData) (*dynamodb.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(*cfg), nil
}

func EC2Client(ctx context.Context, d *plugin.QueryData) (*ec2.Client,
	error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return ec2.NewFromConfig(*cfg), nil
}

func ECRClient(ctx context.Context, d *plugin.QueryData) (*ecr.Client,
	error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return ecr.NewFromConfig(*cfg), nil
}

func ECSClient(ctx context.Context, d *plugin.QueryData) (*ecs.Client,
	error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return ecs.NewFromConfig(*cfg), nil
}

func ECRPublicClient(ctx context.Context, d *plugin.QueryData) (*ecrpublic.Client,
	error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return ecrpublic.NewFromConfig(*cfg), nil
}

func Ec2RegionsClient(ctx context.Context, d *plugin.QueryData, region string) (*ec2.Client, error) {
	// We can query EC2 for the list of supported regions. But, if credentials
	// are insufficient this query will retry many times, so we create a special
	// client with a small number of retries to prevent hangs.
	// Note - This is not cached, but usually the result of using this service will be.
	cfg, err := getClientWithMaxRetries(ctx, d, region, 4, 25*time.Millisecond)
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

func ElastiCacheClient(ctx context.Context, d *plugin.QueryData) (*elasticache.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return elasticache.NewFromConfig(*cfg), nil
}

func ELBClient(ctx context.Context, d *plugin.QueryData) (*elasticloadbalancing.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return elasticloadbalancing.NewFromConfig(*cfg), nil
}

func ELBV2Client(ctx context.Context, d *plugin.QueryData) (*elasticloadbalancingv2.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return elasticloadbalancingv2.NewFromConfig(*cfg), nil
}

func EFSClient(ctx context.Context, d *plugin.QueryData) (*efs.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
  	
  if cfg == nil {
		return nil, nil
	}
	return efs.NewFromConfig(*cfg), nil
}

func EventBridgeClient(ctx context.Context, d *plugin.QueryData) (*eventbridge.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return eventbridge.NewFromConfig(*cfg), nil
}

func FSxClient(ctx context.Context, d *plugin.QueryData) (*fsx.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, fsxEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return fsx.NewFromConfig(*cfg), nil
}

func GlueClient(ctx context.Context, d *plugin.QueryData) (*glue.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return glue.NewFromConfig(*cfg), nil
}

func IAMClient(ctx context.Context, d *plugin.QueryData) (*iam.Client, error) {
	cfg, err := getClient(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return iam.NewFromConfig(*cfg), nil
}


func IdentityStoreClient(ctx context.Context, d *plugin.QueryData) (*identitystore.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return identitystore.NewFromConfig(*cfg), nil
}

func InspectorClient(ctx context.Context, d *plugin.QueryData) (*inspector.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return inspector.NewFromConfig(*cfg), nil
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

func KMSClient(ctx context.Context, d *plugin.QueryData) (*kms.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, "kms")
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return kms.NewFromConfig(*cfg), nil
}

func LambdaClient(ctx context.Context, d *plugin.QueryData) (*lambda.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, lambdaEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return lambda.NewFromConfig(*cfg), nil
}

func OrganizationClient(ctx context.Context, d *plugin.QueryData) (*organizations.Client, error) {
	cfg, err := getClient(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return organizations.NewFromConfig(*cfg), nil
}

func RedshiftClient(ctx context.Context, d *plugin.QueryData) (*redshift.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return redshift.NewFromConfig(*cfg), nil
}

func PricingServiceClient(ctx context.Context, d *plugin.QueryData) (*pricing.Client, error) {
	cfg, err := getClient(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return pricing.NewFromConfig(*cfg), nil
}

func RAMClient(ctx context.Context, d *plugin.QueryData) (*ram.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return ram.NewFromConfig(*cfg), nil
}

func RDSClient(ctx context.Context, d *plugin.QueryData) (*rds.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return rds.NewFromConfig(*cfg), nil
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

func SageMakerClient(ctx context.Context, d *plugin.QueryData) (*sagemaker.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return sagemaker.NewFromConfig(*cfg), nil
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

func WAFClient(ctx context.Context, d *plugin.QueryData) (*waf.Client, error) {
	cfg, err := getClient(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return waf.NewFromConfig(*cfg), nil
}

func Route53DomainsClient(ctx context.Context, d *plugin.QueryData) (*route53domains.Client, error) {
	cfg, err := getClient(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return route53domains.NewFromConfig(*cfg), nil
}

func Route53Client(ctx context.Context, d *plugin.QueryData) (*route53.Client, error) {
	cfg, err := getClient(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return route53.NewFromConfig(*cfg), nil
}

func SQSClient(ctx context.Context, d *plugin.QueryData) (*sqs.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return sqs.NewFromConfig(*cfg), nil
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
