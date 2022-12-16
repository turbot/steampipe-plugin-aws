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
	"github.com/aws/aws-sdk-go-v2/service/amplify"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/appconfig"
	"github.com/aws/aws-sdk-go-v2/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go-v2/service/auditmanager"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudsearch"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
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
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/fsx"
	"github.com/aws/aws-sdk-go-v2/service/glacier"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/aws/aws-sdk-go-v2/service/health"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/inspector"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2"
	"github.com/aws/aws-sdk-go-v2/service/kinesisvideo"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/macie2"
	"github.com/aws/aws-sdk-go-v2/service/mediastore"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/networkfirewall"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/ram"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	"github.com/aws/aws-sdk-go-v2/service/route53resolver"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/serverlessapplicationrepository"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/waf"
	"github.com/aws/aws-sdk-go-v2/service/wafregional"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	"github.com/aws/aws-sdk-go-v2/service/workspaces"
	"github.com/aws/aws-sdk-go/aws/endpoints"

	amplifyEndpoint "github.com/aws/aws-sdk-go/service/amplify"
	auditmanagerEndpoint "github.com/aws/aws-sdk-go/service/auditmanager"
	backupEndpoint "github.com/aws/aws-sdk-go/service/backup"
	cloudsearchEndpoint "github.com/aws/aws-sdk-go/service/cloudsearch"
	codeartifactEndpoint "github.com/aws/aws-sdk-go/service/codeartifact"
	codebuildEndpoint "github.com/aws/aws-sdk-go/service/codebuild"
	codecommitEndpoint "github.com/aws/aws-sdk-go/service/codecommit"
	codepipelineEndpoint "github.com/aws/aws-sdk-go/service/codepipeline"
	daxEndpoint "github.com/aws/aws-sdk-go/service/dax"
	directoryserviceEndpoint "github.com/aws/aws-sdk-go/service/directoryservice"
	dlmEndpoint "github.com/aws/aws-sdk-go/service/dlm"
	dynamodbEndpoint "github.com/aws/aws-sdk-go/service/dynamodb"
	eksEndpoint "github.com/aws/aws-sdk-go/service/eks"
	elasticbeanstalkEndpoint "github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	emrEndpoint "github.com/aws/aws-sdk-go/service/emr"
	eventbridgeEndpoint "github.com/aws/aws-sdk-go/service/eventbridge"
	fsxEndpoint "github.com/aws/aws-sdk-go/service/fsx"
	glacierEndpoint "github.com/aws/aws-sdk-go/service/glacier"
	inspectorEndpoint "github.com/aws/aws-sdk-go/service/inspector"
	kafkaEndpoint "github.com/aws/aws-sdk-go/service/kafka"
	kinesisanalyticsv2Endpoint "github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	kinesisvideoEndpoint "github.com/aws/aws-sdk-go/service/kinesisvideo"
	kmsEndpoint "github.com/aws/aws-sdk-go/service/kms"
	lambdaEndpoint "github.com/aws/aws-sdk-go/service/lambda"
	lightsailEndpoint "github.com/aws/aws-sdk-go/service/lightsail"
	macie2Endpoint "github.com/aws/aws-sdk-go/service/macie2"
	mediastoreEndpoint "github.com/aws/aws-sdk-go/service/mediastore"
	networkfirewallEndpoint "github.com/aws/aws-sdk-go/service/networkfirewall"
	pinpointEndpoint "github.com/aws/aws-sdk-go/service/pinpoint"
	pricingEndpoint "github.com/aws/aws-sdk-go/service/pricing"
	redshiftserverlessEndpoint "github.com/aws/aws-sdk-go/service/redshiftserverless"
	route53resolverEndpoint "github.com/aws/aws-sdk-go/service/route53resolver"
	sagemakerEndpoint "github.com/aws/aws-sdk-go/service/sagemaker"
	securityhubEndpoint "github.com/aws/aws-sdk-go/service/securityhub"
	serverlessrepoEndpoint "github.com/aws/aws-sdk-go/service/serverlessapplicationrepository"
	servicequotasEndpoint "github.com/aws/aws-sdk-go/service/servicequotas"
	sesEndpoint "github.com/aws/aws-sdk-go/service/ses"
	ssmEndpoint "github.com/aws/aws-sdk-go/service/ssm"
	wafregionalEnpoint "github.com/aws/aws-sdk-go/service/wafregional"
	wafv2Enpoint "github.com/aws/aws-sdk-go/service/wafv2"
	wellarchitectedEndpoint "github.com/aws/aws-sdk-go/service/wellarchitected"
	workspacesEndpoint "github.com/aws/aws-sdk-go/service/workspaces"
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

// AccessAnalyzerClient returns the service connection for AWS IAM Access Analyzer service
func AccessAnalyzerClient(ctx context.Context, d *plugin.QueryData) (*accessanalyzer.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return accessanalyzer.NewFromConfig(*cfg), nil
}

func AccountClient(ctx context.Context, d *plugin.QueryData) (*account.Client, error) {
	// AWS Account APIs can be called from any single region. It's best to use
	// the client region of the user (e.g. closest to them)
	clientRegion, err := getClientRegion(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	// Use the client region as default
	// If empty, use the default region, i.e. us-east-1
	queryRegion := clientRegion
	if queryRegion == "" {
		queryRegion, err = getDefaultRegion(ctx, d, nil)
		if err != nil {
			return nil, err
		}
	}

	cfg, err := getClient(ctx, d, queryRegion)
	if err != nil {
		return nil, err
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

func AmplifyClient(ctx context.Context, d *plugin.QueryData) (*amplify.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, amplifyEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return amplify.NewFromConfig(*cfg), nil
}

func APIGatewayClient(ctx context.Context, d *plugin.QueryData) (*apigateway.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return apigateway.NewFromConfig(*cfg), nil
}

func APIGatewayV2Client(ctx context.Context, d *plugin.QueryData) (*apigatewayv2.Client, error) {
	// APIGatewayV2's endpoint ID is the same as APIGateway's, but me-central-1 does not support API Gateway v2 yet
	//cfg, err := getClientForQuerySupportedRegion(ctx, d, apigatewayv2Endpoint.EndpointsID)

	region := d.KeyColumnQualString(matrixKeyRegion)
	validRegions := []string{"af-south-1", "ap-east-1", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ap-southeast-3", "ca-central-1", "eu-central-1", "eu-north-1", "eu-south-1", "eu-west-1", "eu-west-2", "eu-west-3", "me-south-1", "sa-east-1", "us-east-1", "us-east-2", "us-west-1", "us-west-2", "cn-north-1", "cn-northwest-1", "us-gov-east-1", "us-gov-west-1", "us-iso-east-1"}

	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
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

func ApplicationAutoScalingClient(ctx context.Context, d *plugin.QueryData) (*applicationautoscaling.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return applicationautoscaling.NewFromConfig(*cfg), nil
}

func AuditManagerClient(ctx context.Context, d *plugin.QueryData) (*auditmanager.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, auditmanagerEndpoint.EndpointsID)
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
	cfg, err := getClientForQuerySupportedRegion(ctx, d, backupEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
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
		return nil, fmt.Errorf("CloudControlClient called without a region in QueryData")
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

func CodeCommitClient(ctx context.Context, d *plugin.QueryData) (*codecommit.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, codecommitEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return codecommit.NewFromConfig(*cfg), nil
}

func CloudFormationClient(ctx context.Context, d *plugin.QueryData) (*cloudformation.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return cloudformation.NewFromConfig(*cfg), nil
}

func CloudFrontClient(ctx context.Context, d *plugin.QueryData) (*cloudfront.Client, error) {
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return cloudfront.NewFromConfig(*cfg), nil
}

func CloudSearchClient(ctx context.Context, d *plugin.QueryData) (*cloudsearch.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, cloudsearchEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return cloudsearch.NewFromConfig(*cfg), nil
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

func CodeArtifactClient(ctx context.Context, d *plugin.QueryData) (*codeartifact.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, codeartifactEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return codeartifact.NewFromConfig(*cfg), nil
}

func CodeBuildClient(ctx context.Context, d *plugin.QueryData) (*codebuild.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, codebuildEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
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
	cfg, err := getClientForQuerySupportedRegion(ctx, d, codepipelineEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
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
	cfg, err := getClientForDefaultRegion(ctx, d)
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
	cfg, err := getClientForQuerySupportedRegion(ctx, d, daxEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return dax.NewFromConfig(*cfg), nil
}

func DirectoryServiceClient(ctx context.Context, d *plugin.QueryData) (*directoryservice.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, directoryserviceEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return directoryservice.NewFromConfig(*cfg), nil
}

func DLMClient(ctx context.Context, d *plugin.QueryData) (*dlm.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, dlmEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
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
	cfg, err := getClientForQuerySupportedRegion(ctx, d, dynamodbEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
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

// EC2RegionsClient returns the service connection for Amazon EC2 with a specific region endpoint
func EC2RegionsClient(ctx context.Context, d *plugin.QueryData, region string) (*ec2.Client, error) {
	cfg, err := getClient(ctx, d, region)
	if err != nil {
		return nil, err
	}
	return ec2.NewFromConfig(*cfg), nil
}

// EC2RegionsClientWithMaxRetires returns the service connection for Amazon EC2 with a specific region endpoint and the capability of modifying the default retry mechanism.
func EC2RegionsClientWithMaxRetires(ctx context.Context, d *plugin.QueryData, region string) (*ec2.Client, error) {
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

func ECRClient(ctx context.Context, d *plugin.QueryData) (*ecr.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return ecr.NewFromConfig(*cfg), nil
}

// ECRPublicClient returns the service connection for Amazon ECR Public endpoints
func ECRPublicClient(ctx context.Context, d *plugin.QueryData) (*ecrpublic.Client, error) {
	// Amazon ecr-public actions are only supported when providing us-east-1.
	// Amazon ECR Public repositories can be created in many other regions, but
	// the client requires authentication in the us-east-1
	// https://docs.aws.amazon.com/AmazonECR/latest/public/getting-started-cli.html
	//
	// As of December 16, 2022, Amazon ecr-public actions are not supported in AWS US Gov Cloud
	// Using the gov cloud creds will results in an error
	// api error UnrecognizedClientException: The security token included in the request is invalid.
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}

	return ecrpublic.NewFromConfig(*cfg), nil
}

func ECSClient(ctx context.Context, d *plugin.QueryData) (*ecs.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return ecs.NewFromConfig(*cfg), nil
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

func EKSClient(ctx context.Context, d *plugin.QueryData) (*eks.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, eksEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return eks.NewFromConfig(*cfg), nil
}

func ElastiCacheClient(ctx context.Context, d *plugin.QueryData) (*elasticache.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return elasticache.NewFromConfig(*cfg), nil
}

func ElasticBeanstalkClient(ctx context.Context, d *plugin.QueryData) (*elasticbeanstalk.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, elasticbeanstalkEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return elasticbeanstalk.NewFromConfig(*cfg), nil
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

func ElasticsearchClient(ctx context.Context, d *plugin.QueryData) (*elasticsearchservice.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return elasticsearchservice.NewFromConfig(*cfg), nil
}

func EMRClient(ctx context.Context, d *plugin.QueryData) (*emr.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, emrEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return emr.NewFromConfig(*cfg), nil
}

func EventBridgeClient(ctx context.Context, d *plugin.QueryData) (*eventbridge.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, eventbridgeEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return eventbridge.NewFromConfig(*cfg), nil
}

func FirehoseClient(ctx context.Context, d *plugin.QueryData) (*firehose.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return firehose.NewFromConfig(*cfg), nil
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

func GlacierClient(ctx context.Context, d *plugin.QueryData) (*glacier.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, glacierEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return glacier.NewFromConfig(*cfg), nil
}

func GlobalAcceleratorClient(ctx context.Context, d *plugin.QueryData) (*globalaccelerator.Client, error) {
	// Global Accelerator is a global service that supports endpoints in multiple AWS Regions but you must specify
	// the us-west-2 (Oregon) Region to create or update accelerators.
	cfg, err := getClient(ctx, d, "us-west-2")
	if err != nil {
		return nil, err
	}
	return globalaccelerator.NewFromConfig(*cfg), nil
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

func GuardDutyClient(ctx context.Context, d *plugin.QueryData) (*guardduty.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return guardduty.NewFromConfig(*cfg), nil
}

func HealthClient(ctx context.Context, d *plugin.QueryData) (*health.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return health.NewFromConfig(*cfg), nil
}

func IAMClient(ctx context.Context, d *plugin.QueryData) (*iam.Client, error) {
	cfg, err := getClientForDefaultRegion(ctx, d)
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
	cfg, err := getClientForQuerySupportedRegion(ctx, d, inspectorEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return inspector.NewFromConfig(*cfg), nil
}

func KafkaClient(ctx context.Context, d *plugin.QueryData) (*kafka.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, kafkaEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return kafka.NewFromConfig(*cfg), nil
}

func KinesisClient(ctx context.Context, d *plugin.QueryData) (*kinesis.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return kinesis.NewFromConfig(*cfg), nil
}

func KinesisAnalyticsV2Client(ctx context.Context, d *plugin.QueryData) (*kinesisanalyticsv2.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, kinesisanalyticsv2Endpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return kinesisanalyticsv2.NewFromConfig(*cfg), nil
}

func KinesisVideoClient(ctx context.Context, d *plugin.QueryData) (*kinesisvideo.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, kinesisvideoEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return kinesisvideo.NewFromConfig(*cfg), nil
}

func KMSClient(ctx context.Context, d *plugin.QueryData) (*kms.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, kmsEndpoint.EndpointsID)
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

func LightsailClient(ctx context.Context, d *plugin.QueryData) (*lightsail.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, lightsailEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return lightsail.NewFromConfig(*cfg), nil
}

func Macie2Client(ctx context.Context, d *plugin.QueryData) (*macie2.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, macie2Endpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return macie2.NewFromConfig(*cfg), nil
}

func MediaStoreClient(ctx context.Context, d *plugin.QueryData) (*mediastore.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, mediastoreEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return mediastore.NewFromConfig(*cfg), nil
}

func NeptuneClient(ctx context.Context, d *plugin.QueryData) (*neptune.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return neptune.NewFromConfig(*cfg), nil
}

func NetworkFirewallClient(ctx context.Context, d *plugin.QueryData) (*networkfirewall.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, networkfirewallEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return networkfirewall.NewFromConfig(*cfg), nil
}

func OpenSearchClient(ctx context.Context, d *plugin.QueryData) (*opensearch.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return opensearch.NewFromConfig(*cfg), nil
}

func OrganizationClient(ctx context.Context, d *plugin.QueryData) (*organizations.Client, error) {
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return organizations.NewFromConfig(*cfg), nil
}

func PinpointClient(ctx context.Context, d *plugin.QueryData) (*pinpoint.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, pinpointEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return pinpoint.NewFromConfig(*cfg), nil
}

func PricingClient(ctx context.Context, d *plugin.QueryData) (*pricing.Client, error) {
	// Get Pricing API supported regions
	pricingAPISupportedRegions, err := GetSupportedRegionsForClient(ctx, d, pricingEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}

	// Get the client region for AWS API calls
	// Typically this should be the region closest to the user
	clientRegion, err := getClientRegion(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	// Pricing API is a global API that supports only us-east-1 and ap-south-1 regions
	// If a preferred region is set using client_region, or in the AWS config files,
	// and the API supports that region, use that as the endpoint
	// As of Dec 13, 2022, AWS Pricing API only supports in AWS Commercial Cloud
	// Default set to us-east-1 for now
	queryRegion := clientRegion
	if !helpers.StringSliceContains(pricingAPISupportedRegions, queryRegion) {
		queryRegion, err = getDefaultRegion(ctx, d, nil)
		if err != nil {
			return nil, err
		}
	}

	cfg, err := getClient(ctx, d, queryRegion)
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

func RDSDBProxyClient(ctx context.Context, d *plugin.QueryData) (*rds.Client, error) {

	// AWS RDS DD Proxy's endpoint ID is the same as AWS RDS's, but me-central-1 and ap-northeast-3 does not support for it
	//cfg, err := getClientForQuerySupportedRegion(ctx, d, rdxEndpoint.EndpointsID)

	region := d.KeyColumnQualString(matrixKeyRegion)
	validRegions := []string{"af-south-1", "ap-east-1", "ap-northeast-1", "ap-northeast-2", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ap-southeast-3", "ca-central-1", "eu-central-1", "eu-north-1", "eu-south-1", "eu-west-1", "eu-west-2", "eu-west-3", "me-south-1", "sa-east-1", "us-east-1", "us-east-2", "us-west-1", "us-west-2", "cn-north-1", "cn-northwest-1", "us-gov-east-1", "us-gov-west-1", "us-iso-east-1"}

	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return rds.NewFromConfig(*cfg), nil
}

func RedshiftClient(ctx context.Context, d *plugin.QueryData) (*redshift.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return redshift.NewFromConfig(*cfg), nil
}

func RedshiftServerlessClient(ctx context.Context, d *plugin.QueryData) (*redshiftserverless.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, redshiftserverlessEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return redshiftserverless.NewFromConfig(*cfg), nil
}

func ResourceExplorerClient(ctx context.Context, d *plugin.QueryData, region string) (*resourceexplorer2.Client, error) {
	// https://aws.amazon.com/about-aws/whats-new/2022/11/announcing-aws-resource-explorer/
	// AWS Resource Explorer is generally available in the following AWS Regions, with more Regions coming soon: US East (Ohio), US East (N. Virginia), US West (N. California), US West (Oregon), Asia Pacific (Mumbai), Asia Pacific (Osaka), Asia Pacific (Seoul), Asia Pacific (Singapore), Asia Pacific (Sydney), Asia Pacific (Tokyo), Canada (Central), Europe (Frankfurt), Europe (Ireland), Europe (London), Europe (Paris), Europe (Stockholm), and South America (São Paulo).
	var resourceExplorerRegions = []string{"ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ca-central-1", "eu-central-1", "eu-north-1", "eu-west-1", "eu-west-2", "eu-west-3", "sa-east-1", "us-east-1", "us-east-2", "us-west-1", "us-west-2"}

	if region == "" {
		return nil, fmt.Errorf("region must be passed ResourceExplorerClient")
	}

	// If not a supported region return nil client
	if !helpers.StringSliceContains(resourceExplorerRegions, region) {
		return nil, nil
	}

	cfg, err := getClient(ctx, d, region)
	if err != nil {
		return nil, err
	}
	return resourceexplorer2.NewFromConfig(*cfg), nil
}

func ResourceGroupsTaggingClient(ctx context.Context, d *plugin.QueryData) (*resourcegroupstaggingapi.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return resourcegroupstaggingapi.NewFromConfig(*cfg), nil
}

func Route53Client(ctx context.Context, d *plugin.QueryData) (*route53.Client, error) {
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return route53.NewFromConfig(*cfg), nil
}

func Route53DomainsClient(ctx context.Context, d *plugin.QueryData) (*route53domains.Client, error) {
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return route53domains.NewFromConfig(*cfg), nil
}

func Route53ResolverClient(ctx context.Context, d *plugin.QueryData) (*route53resolver.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, route53resolverEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return route53resolver.NewFromConfig(*cfg), nil
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
	cfg, err := getClientForRegion(ctx, d, region)
	if err != nil {
		return nil, err
	}
	return s3control.NewFromConfig(*cfg), nil
}

func SageMakerClient(ctx context.Context, d *plugin.QueryData) (*sagemaker.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, sagemakerEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return sagemaker.NewFromConfig(*cfg), nil
}

func SecretsManagerClient(ctx context.Context, d *plugin.QueryData) (*secretsmanager.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return secretsmanager.NewFromConfig(*cfg), nil
}

func SecurityHubClient(ctx context.Context, d *plugin.QueryData) (*securityhub.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, securityhubEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return securityhub.NewFromConfig(*cfg), nil
}

// Added for using middleware for migrating table "aws_securityhub_member"
// See https://github.com/aws/aws-sdk-go-v2/issues/1884#issuecomment-1278567756 for more info
func SecurityHubClientConfig(ctx context.Context, d *plugin.QueryData) (*aws.Config, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, securityhubEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return cfg, nil
}

func SESClient(ctx context.Context, d *plugin.QueryData) (*ses.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, sesEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return ses.NewFromConfig(*cfg), nil
}

func ServerlessApplicationRepositoryClient(ctx context.Context, d *plugin.QueryData) (*serverlessapplicationrepository.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, serverlessrepoEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return serverlessapplicationrepository.NewFromConfig(*cfg), nil
}

func ServiceQuotasClient(ctx context.Context, d *plugin.QueryData) (*servicequotas.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, servicequotasEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return servicequotas.NewFromConfig(*cfg), nil
}

func StepFunctionsClient(ctx context.Context, d *plugin.QueryData) (*sfn.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return sfn.NewFromConfig(*cfg), nil
}

func SNSClient(ctx context.Context, d *plugin.QueryData) (*sns.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return sns.NewFromConfig(*cfg), nil
}

func SSMClient(ctx context.Context, d *plugin.QueryData) (*ssm.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, ssmEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return ssm.NewFromConfig(*cfg), nil
}

func SQSClient(ctx context.Context, d *plugin.QueryData) (*sqs.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return sqs.NewFromConfig(*cfg), nil
}

func STSClient(ctx context.Context, d *plugin.QueryData) (*sts.Client, error) {
	// STS is available in each region, so we can use the client_region
	// closest to the user.
	cfg, err := getClientForClientRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return sts.NewFromConfig(*cfg), nil
}

func SSOAdminClient(ctx context.Context, d *plugin.QueryData) (*ssoadmin.Client, error) {
	// https://github.com/aws/aws-sdk-go/blob/main/aws/endpoints/defaults.go#L17417
	cfg, err := getClientForQuerySupportedRegion(ctx, d, "portal.sso")
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return ssoadmin.NewFromConfig(*cfg), nil
}

func WAFClient(ctx context.Context, d *plugin.QueryData) (*waf.Client, error) {
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return waf.NewFromConfig(*cfg), nil
}

func WAFRegionalClient(ctx context.Context, d *plugin.QueryData) (*wafregional.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, wafregionalEnpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return wafregional.NewFromConfig(*cfg), nil
}

func WAFV2Client(ctx context.Context, d *plugin.QueryData, region string) (*wafv2.Client, error) {
	validRegions, err := GetSupportedRegionsForClient(ctx, d, wafv2Enpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if !helpers.StringSliceContains(validRegions, region) {
		// We choose to ignore unsupported regions rather than returning an error
		// for them - it's a better user experience. So, return a nil session rather
		// than an error. The caller must handle this case.
		return nil, nil
	}
	cfg, err := getClientForRegion(ctx, d, region)
	if err != nil {
		return nil, err
	}
	return wafv2.NewFromConfig(*cfg), nil
}

func WellArchitectedClient(ctx context.Context, d *plugin.QueryData) (*wellarchitected.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, wellarchitectedEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return wellarchitected.NewFromConfig(*cfg), nil
}

func WorkspacesClient(ctx context.Context, d *plugin.QueryData) (*workspaces.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, workspacesEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return workspaces.NewFromConfig(*cfg), nil
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
		return nil, fmt.Errorf("getClientForQuerySupportedRegion called without a region in QueryData")
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
	return getClient(ctx, d, region)
}

// Helper function to get the session for the preferred region in this partition
func getClientForClientRegion(ctx context.Context, d *plugin.QueryData) (*aws.Config, error) {
	r, err := getClientRegion(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return getClient(ctx, d, r)
}

// Helper function to get the session for the default region in this partition
func getClientForDefaultRegion(ctx context.Context, d *plugin.QueryData) (*aws.Config, error) {
	r, err := getDefaultRegion(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return getClient(ctx, d, r)
}

// Helper function to get the session for a region set in query data
func getClientForQueryRegion(ctx context.Context, d *plugin.QueryData) (*aws.Config, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	if region == "" {
		return nil, fmt.Errorf("getClientForQuerySupportedRegion called without a region in QueryData")
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

// GetSupportedRegionsForClient lists valid regions for a service based on service ID
func GetSupportedRegionsForClient(ctx context.Context, d *plugin.QueryData, serviceID string) ([]string, error) {
	var partitionName string
	var partition endpoints.Partition

	// If valid regions list is already available in cache, return it
	cacheKey := fmt.Sprintf("supported-regions-%s", serviceID)
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]string), nil
	}

	// get the partition of the AWS account plugin is connected to
	// using the cached version of getCommonColumns
	commonColumnData, err := getCommonColumnsCached(ctx, d, nil)
	if err != nil {
		plugin.Logger(ctx).Error("GetSupportedRegionsForClient", "unable to get partition name", err)
		return nil, err
	}
	partitionName = commonColumnData.(*awsCommonColumnData).Partition

	// Get AWS partition based on the partition name
	switch partitionName {
	case endpoints.AwsPartitionID:
		partition = endpoints.AwsPartition()
	case endpoints.AwsCnPartitionID:
		partition = endpoints.AwsCnPartition()
	case endpoints.AwsIsoPartitionID:
		partition = endpoints.AwsIsoPartition()
	case endpoints.AwsUsGovPartitionID:
		partition = endpoints.AwsUsGovPartition()
	case endpoints.AwsIsoBPartitionID:
		partition = endpoints.AwsIsoBPartition()
	default:
		plugin.Logger(ctx).Error("service_v2.GetSupportedRegionsForClient", "invalid_partition_error", fmt.Errorf("%s is an invalid partition", partitionName))
		return nil, fmt.Errorf("service_v2.GetSupportedRegionsForClient:: '%s' is an invalid partition", partitionName)
	}

	var validRegions []string

	// Get the list of the service regions based on the service ID
	services := partition.Services()
	serviceInfo, ok := services[serviceID]
	if !ok {
		return nil, fmt.Errorf("service_v2.SupportedRegionsForClient called with invalid service ID: %s", serviceID)
	}

	regions := serviceInfo.Regions()
	for rs := range regions {
		validRegions = append(validRegions, rs)
	}

	// Save valid regions in the cache
	d.ConnectionManager.Cache.Set(cacheKey, validRegions)
	return validRegions, nil
}
