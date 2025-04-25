package aws

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/accessanalyzer"
	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acmpca"
	"github.com/aws/aws-sdk-go-v2/service/amplify"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/appconfig"
	"github.com/aws/aws-sdk-go-v2/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	"github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/athena"
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
	"github.com/aws/aws-sdk-go-v2/service/codestarnotifications"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costoptimizationhub"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"
	"github.com/aws/aws-sdk-go-v2/service/dax"
	"github.com/aws/aws-sdk-go-v2/service/directoryservice"
	"github.com/aws/aws-sdk-go-v2/service/dlm"
	"github.com/aws/aws-sdk-go-v2/service/docdb"
	"github.com/aws/aws-sdk-go-v2/service/drs"
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
	"github.com/aws/aws-sdk-go-v2/service/fms"
	"github.com/aws/aws-sdk-go-v2/service/fsx"
	"github.com/aws/aws-sdk-go-v2/service/glacier"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/aws/aws-sdk-go-v2/service/health"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/inspector"
	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/keyspaces"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2"
	"github.com/aws/aws-sdk-go-v2/service/kinesisvideo"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/lakeformation"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/macie2"
	"github.com/aws/aws-sdk-go-v2/service/mediastore"
	"github.com/aws/aws-sdk-go-v2/service/memorydb"
	"github.com/aws/aws-sdk-go-v2/service/mgn"
	"github.com/aws/aws-sdk-go-v2/service/mq"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/networkfirewall"
	"github.com/aws/aws-sdk-go-v2/service/oam"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
	"github.com/aws/aws-sdk-go-v2/service/pipes"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"
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
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securitylake"
	"github.com/aws/aws-sdk-go-v2/service/serverlessapplicationrepository"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/shield"
	"github.com/aws/aws-sdk-go-v2/service/simspaceweaver"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssmincidents"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/support"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"github.com/aws/aws-sdk-go-v2/service/transfer"
	"github.com/aws/aws-sdk-go-v2/service/waf"
	"github.com/aws/aws-sdk-go-v2/service/wafregional"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	"github.com/aws/aws-sdk-go-v2/service/workspaces"
	"github.com/aws/smithy-go/logging"
	"github.com/hashicorp/go-hclog"
	"github.com/rs/dnscache"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"golang.org/x/sync/semaphore"

	amplifyEndpoint "github.com/aws/aws-sdk-go/service/amplify"
	apigatewayv2Endpoint "github.com/aws/aws-sdk-go/service/apigatewayv2"
	appRunnerEndpoint "github.com/aws/aws-sdk-go/service/apprunner"
	appsyncv2Endpoint "github.com/aws/aws-sdk-go/service/appsync"
	auditmanagerEndpoint "github.com/aws/aws-sdk-go/service/auditmanager"
	backupEndpoint "github.com/aws/aws-sdk-go/service/backup"
	cloudsearchEndpoint "github.com/aws/aws-sdk-go/service/cloudsearch"
	codeartifactEndpoint "github.com/aws/aws-sdk-go/service/codeartifact"
	codebuildEndpoint "github.com/aws/aws-sdk-go/service/codebuild"
	codecommitEndpoint "github.com/aws/aws-sdk-go/service/codecommit"
	codepipelineEndpoint "github.com/aws/aws-sdk-go/service/codepipeline"
	cognitoidentityEndpoint "github.com/aws/aws-sdk-go/service/cognitoidentity"
	cognitoidentityproviderEndpoint "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	daxEndpoint "github.com/aws/aws-sdk-go/service/dax"
	directoryserviceEndpoint "github.com/aws/aws-sdk-go/service/directoryservice"
	dlmEndpoint "github.com/aws/aws-sdk-go/service/dlm"
	drsEndpoint "github.com/aws/aws-sdk-go/service/drs"
	dynamodbEndpoint "github.com/aws/aws-sdk-go/service/dynamodb"
	eksEndpoint "github.com/aws/aws-sdk-go/service/eks"
	elasticbeanstalkEndpoint "github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	emrEndpoint "github.com/aws/aws-sdk-go/service/emr"
	eventbridgeEndpoint "github.com/aws/aws-sdk-go/service/eventbridge"
	fsxEndpoint "github.com/aws/aws-sdk-go/service/fsx"
	glacierEndpoint "github.com/aws/aws-sdk-go/service/glacier"
	healthEndpoint "github.com/aws/aws-sdk-go/service/health"
	inspectorEndpoint "github.com/aws/aws-sdk-go/service/inspector"
	inspector2Endpoint "github.com/aws/aws-sdk-go/service/inspector2"
	iotEndpoint "github.com/aws/aws-sdk-go/service/iot"
	kafkaEndpoint "github.com/aws/aws-sdk-go/service/kafka"
	keyspacesEndpoint "github.com/aws/aws-sdk-go/service/keyspaces"
	kinesisanalyticsv2Endpoint "github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	kinesisvideoEndpoint "github.com/aws/aws-sdk-go/service/kinesisvideo"
	kmsEndpoint "github.com/aws/aws-sdk-go/service/kms"
	lambdaEndpoint "github.com/aws/aws-sdk-go/service/lambda"
	lightsailEndpoint "github.com/aws/aws-sdk-go/service/lightsail"
	macie2Endpoint "github.com/aws/aws-sdk-go/service/macie2"
	mediastoreEndpoint "github.com/aws/aws-sdk-go/service/mediastore"
	memoryDBEndpoint "github.com/aws/aws-sdk-go/service/memorydb"
	mgnEndpoint "github.com/aws/aws-sdk-go/service/mgn"
	mqEndpoint "github.com/aws/aws-sdk-go/service/mq"
	networkfirewallEndpoint "github.com/aws/aws-sdk-go/service/networkfirewall"
	oamEndpoint "github.com/aws/aws-sdk-go/service/oam"
	pinpointEndpoint "github.com/aws/aws-sdk-go/service/pinpoint"
	pipesEndpoint "github.com/aws/aws-sdk-go/service/pipes"
	pricingEndpoint "github.com/aws/aws-sdk-go/service/pricing"
	rdsEndpoint "github.com/aws/aws-sdk-go/service/rds"
	redshiftserverlessEndpoint "github.com/aws/aws-sdk-go/service/redshiftserverless"
	resourceexplorer2Endpoint "github.com/aws/aws-sdk-go/service/resourceexplorer2"
	route53resolverEndpoint "github.com/aws/aws-sdk-go/service/route53resolver"
	sagemakerEndpoint "github.com/aws/aws-sdk-go/service/sagemaker"
	schedulerEndpoint "github.com/aws/aws-sdk-go/service/scheduler"
	securityhubEndpoint "github.com/aws/aws-sdk-go/service/securityhub"
	securitylakeEndpoint "github.com/aws/aws-sdk-go/service/securitylake"
	serverlessrepoEndpoint "github.com/aws/aws-sdk-go/service/serverlessapplicationrepository"
	servicequotasEndpoint "github.com/aws/aws-sdk-go/service/servicequotas"
	sesEndpoint "github.com/aws/aws-sdk-go/service/ses"
	simspaceWeaverEndpoint "github.com/aws/aws-sdk-go/service/simspaceweaver"
	ssmEndpoint "github.com/aws/aws-sdk-go/service/ssm"
	ssmIncidentsEndpoint "github.com/aws/aws-sdk-go/service/ssmincidents"
	ssoEndpoint "github.com/aws/aws-sdk-go/service/sso"
	timestreamwriteEndpoint "github.com/aws/aws-sdk-go/service/timestreamwrite"
	transferEndpoint "github.com/aws/aws-sdk-go/service/transfer"
	wafregionalEndpoint "github.com/aws/aws-sdk-go/service/wafregional"
	wafv2Endpoint "github.com/aws/aws-sdk-go/service/wafv2"
	wellarchitectedEndpoint "github.com/aws/aws-sdk-go/service/wellarchitected"
	workspacesEndpoint "github.com/aws/aws-sdk-go/service/workspaces"
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

// AccountClient is used to query general information about an AWS account.
func AccountClient(ctx context.Context, d *plugin.QueryData) (*account.Client, error) {
	// Use the client region - service is global but available in all regions.
	cfg, err := getClientForDefaultRegion(ctx, d)
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

func ACMPCAClient(ctx context.Context, d *plugin.QueryData) (*acmpca.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return acmpca.NewFromConfig(*cfg), nil
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
	// API Gateway V2 has the same endpoint information in the SDK as API Gateway, but
	// is actually available in less regions. We have to manually remove them
	// here.
	// Source - https://www.aws-services.info/apigatewayv2.html
	excludeRegions := []string{
		"ap-south-2",     // Hyderabad
		"ap-southeast-3", // Jakarta
		"ap-southeast-4", // Melbourne
		"eu-central-2",   // Zurich
		"eu-south-2",     // Spain
		"il-central-1",   // Israel (Tel Aviv)
	}
	cfg, err := getClientForQuerySupportedRegionWithExclusions(ctx, d, apigatewayv2Endpoint.EndpointsID, excludeRegions)
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

func AppRunnerClient(ctx context.Context, d *plugin.QueryData) (*apprunner.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, appRunnerEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return apprunner.NewFromConfig(*cfg), nil
}

func AppStreamClient(ctx context.Context, d *plugin.QueryData) (*appstream.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return appstream.NewFromConfig(*cfg), nil
}

func AppSyncClient(ctx context.Context, d *plugin.QueryData) (*appsync.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, appsyncv2Endpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return appsync.NewFromConfig(*cfg), nil
}

func AthenaClient(ctx context.Context, d *plugin.QueryData) (*athena.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return athena.NewFromConfig(*cfg), nil
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
	region := d.EqualsQualString(matrixKeyRegion)
	cfg, err := getClientWithMaxRetries(ctx, d, region, 4, 25*time.Millisecond)
	if err != nil {
		return nil, err
	}
	return cloudcontrol.NewFromConfig(*cfg), nil
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
	// CloudFront a global service with a single DNS endpoint
	// (cloudfront.amazonaws.com).
	// https://docs.aws.amazon.com/general/latest/gr/cf_region.html
	// So, while requests will go to the global endpoint, we can still prefer /
	// reuse the client region.
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

func CodeStarNotificationsClient(ctx context.Context, d *plugin.QueryData) (*codestarnotifications.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return codestarnotifications.NewFromConfig(*cfg), nil
}

func CognitoIdentityClient(ctx context.Context, d *plugin.QueryData) (*cognitoidentity.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, cognitoidentityEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return cognitoidentity.NewFromConfig(*cfg), nil
}

func CognitoIdentityProviderClient(ctx context.Context, d *plugin.QueryData) (*cognitoidentityprovider.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, cognitoidentityproviderEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return cognitoidentityprovider.NewFromConfig(*cfg), nil
}

func ConfigClient(ctx context.Context, d *plugin.QueryData) (*configservice.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return configservice.NewFromConfig(*cfg), nil
}

func CostExplorerClient(ctx context.Context, d *plugin.QueryData) (*costexplorer.Client, error) {
	// Cost Explorer is a global service that operates from a single
	// region (ce.us-east-1.amazonaws.com).
	// https://docs.aws.amazon.com/general/latest/gr/billing.html
	// Testing shows it works with either the default or client region object,
	// so use client region for higher reuse.
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return costexplorer.NewFromConfig(*cfg), nil
}

func CostOptimizationHubClient(ctx context.Context, d *plugin.QueryData) (*costoptimizationhub.Client, error) {
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return costoptimizationhub.NewFromConfig(*cfg), nil
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

func DRSClient(ctx context.Context, d *plugin.QueryData) (*drs.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, drsEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return drs.NewFromConfig(*cfg), nil
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

// Get an EC2 client for a specific region. Used by various hydrate functions
// pulling data from regions other than the query region.
func EC2ClientForRegion(ctx context.Context, d *plugin.QueryData, region string) (*ec2.Client, error) {
	cfg, err := getClient(ctx, d, region)
	if err != nil {
		return nil, err
	}
	return ec2.NewFromConfig(*cfg), nil
}

// Get an EC2 client with a small number of retries. Used in very specific
// situations like listing regions where fast failure is preferred over a long
// retry/backoff loop. Do not use for general tables.
func EC2LowRetryClientForRegion(ctx context.Context, d *plugin.QueryData, region string) (*ec2.Client, error) {
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

func ECRPublicClient(ctx context.Context, d *plugin.QueryData) (*ecrpublic.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
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

func FMSClient(ctx context.Context, d *plugin.QueryData) (*fms.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return fms.NewFromConfig(*cfg), nil
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
	// Global Accelerator is a global service with a single DNS endpoint
	// (globalaccelerator.amazonaws.com). As of 2023-01-18, it's only
	// available in the us-west-2 region. It doesn't resolve if we use
	// client region, and it's not using the default region, so we have no
	// choice but to hard code it here.
	// https://docs.aws.amazon.com/general/latest/gr/global_accelerator.html
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
	// Get Health API supported regions
	healthAPISupportedRegions, err := listRegionsForService(ctx, d, healthEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}

	// Get the client region for AWS API calls
	// Typically this should be the region closest to the user
	clientRegion, err := getDefaultRegion(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	// Health API is a global API that supports only us-east-1 and us-east-2 regions
	// in `aws` partition and us-gov-west-1 in `aws-gov` partition.
	// If a preferred region is set using default_region, or in the AWS config files,
	// and the API supports that region, use that as the endpoint.
	// As of April 12, 2025, AWS Health API only works in AWS Commercial Cloud and GovCloud.
	queryRegion := clientRegion
	if !slices.Contains(healthAPISupportedRegions, queryRegion) {
		queryRegion, err = getLastResortRegion(ctx, d, nil)
		if err != nil {
			return nil, err
		}
	}

	cfg, err := getClient(ctx, d, queryRegion)
	if err != nil {
		return nil, err
	}

	return health.NewFromConfig(*cfg), nil
}

func IAMClient(ctx context.Context, d *plugin.QueryData) (*iam.Client, error) {
	// IAM a global service with a single DNS endpoint (iam.amazonaws.com).
	// https://docs.aws.amazon.com/general/latest/gr/iam-service.html
	// So, while requests will go to the global endpoint, we can still prefer /
	// reuse the client region.
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

func Inspector2Client(ctx context.Context, d *plugin.QueryData) (*inspector2.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, inspector2Endpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return inspector2.NewFromConfig(*cfg), nil
}

func IoTClient(ctx context.Context, d *plugin.QueryData) (*iot.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, iotEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return iot.NewFromConfig(*cfg), nil
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

func KeyspacesClient(ctx context.Context, d *plugin.QueryData) (*keyspaces.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, keyspacesEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return keyspaces.NewFromConfig(*cfg), nil
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

func LakeFormationClient(ctx context.Context, d *plugin.QueryData) (*lakeformation.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return lakeformation.NewFromConfig(*cfg), nil
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

func MemoryDBClient(ctx context.Context, d *plugin.QueryData) (*memorydb.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, memoryDBEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return memorydb.NewFromConfig(*cfg), nil
}

func MGNClient(ctx context.Context, d *plugin.QueryData) (*mgn.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, mgnEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return mgn.NewFromConfig(*cfg), nil
}

func MQClient(ctx context.Context, d *plugin.QueryData) (*mq.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, mqEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return mq.NewFromConfig(*cfg), nil
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

func OAMClient(ctx context.Context, d *plugin.QueryData) (*oam.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, oamEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return oam.NewFromConfig(*cfg), nil
}

func OpenSearchClient(ctx context.Context, d *plugin.QueryData) (*opensearch.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return opensearch.NewFromConfig(*cfg), nil
}

func OrganizationClient(ctx context.Context, d *plugin.QueryData) (*organizations.Client, error) {
	// Organizations is a global service that operates from a single
	// region (organizations.us-east-1.amazonaws.com).
	// https://docs.aws.amazon.com/general/latest/gr/ao.html
	// So, we must specify the default region rather than the client region.
	// Testing shows it works with either the default or client region object,
	// so use client region for higher reuse.
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

func PipesClient(ctx context.Context, d *plugin.QueryData) (*pipes.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, pipesEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return pipes.NewFromConfig(*cfg), nil
}

func PricingClient(ctx context.Context, d *plugin.QueryData) (*pricing.Client, error) {

	// The Pricing API is different from other services. It's a global service,
	// but only available from two regions:
	// - us-east-1
	// - ap-south-1
	// There is a big latency difference between these regions, so we do our
	// best here to use the region you've chosen.  This could be smarter (e.g.
	// choose closest), but for now it just tries to use your client region if
	// it can and otherwise falls back to the default region us-east-1.

	// Get Pricing API supported regions
	pricingAPISupportedRegions, err := listRegionsForService(ctx, d, pricingEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}

	// Get the client region for AWS API calls
	// Typically this should be the region closest to the user
	clientRegion, err := getDefaultRegion(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	// Pricing API is a global API that supports only us-east-1 and ap-south-1 regions.
	// If a preferred region is set using default_region, or in the AWS config files,
	// and the API supports that region, use that as the endpoint.
	// As of Dec 13, 2022, AWS Pricing API only works in AWS Commercial Cloud.
	queryRegion := clientRegion
	if !slices.Contains(pricingAPISupportedRegions, queryRegion) {
		queryRegion, err = getLastResortRegion(ctx, d, nil)
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

// QuickSightClient returns the client for the Amazon QuickSight service
func QuickSightClient(ctx context.Context, d *plugin.QueryData) (*quicksight.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return quicksight.NewFromConfig(*cfg), nil
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

func RDSDBRecommendationClient(ctx context.Context, d *plugin.QueryData) (*rds.Client, error) {
	// RDS DB Recommendation has the same endpoint information in the SDK as RDS, but
	// is actually available in less regions. We have to manually remove them
	// here.
	// Source - https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UserRecommendationsView.html
	excludeRegions := []string{
		"af-south-1",     // Africa (Cape Town)
		"ap-east-1",      // Asia Pacific (Hong Kong)
		"ap-northeast-3", // Asia Pacific (Osaka)
		"ap-southeast-3", // Asia Pacific (Jakarta)
		"eu-north-1",     // Europe (Stockholm)
		"eu-south-1",     // Europe (Milan)
		"me-central-1",   // Middle East (UAE)
		"me-south-1",     // Middle East (Bahrain)
		"us-gov-west-1",  // AWS GovCloud (US-West)
		"us-gov-east-1",  // AWS GovCloud (US-East)
		"cn-north-1",     // China (Beijing)
		"cn-northwest-1", // China (Ningxia)
	}
	excludeRegions = append(excludeRegions, getRegionByPartition("aws-cn")...)
	excludeRegions = append(excludeRegions, getRegionByPartition("aws-us-gov")...)
	cfg, err := getClientForQuerySupportedRegionWithExclusions(ctx, d, rdsEndpoint.EndpointsID, excludeRegions)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return rds.NewFromConfig(*cfg), nil
}

func RDSDBProxyClient(ctx context.Context, d *plugin.QueryData) (*rds.Client, error) {
	// RDS DB Proxy has the same endpoint information in the SDK as RDS, but
	// is actually available in less regions. We have to manually remove them
	// here.
	// Source - https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RDS_Fea_Regions_DB-eng.Feature.RDSProxy.html
	excludeRegions := []string{
		"ap-south-2",     // Hyderabad
		"ap-southeast-3", // Jakarta
		"ap-southeast-4", // Melbourne
		"eu-central-2",   // Zurich
		"eu-south-2",     // Spain
		"me-central-1",   // UAE
	}
	excludeRegions = append(excludeRegions, getRegionByPartition("aws-cn")...)
	excludeRegions = append(excludeRegions, getRegionByPartition("aws-us-gov")...)
	cfg, err := getClientForQuerySupportedRegionWithExclusions(ctx, d, rdsEndpoint.EndpointsID, excludeRegions)
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

	// Resource Explorer can be called in any specific region. It's normally
	// used from the client region, but the region can be optionally provided
	// as a qual in the query. So, in this function we take the desired region
	// and check that a valid endpoint for the service before using it.

	if region == "" {
		return nil, fmt.Errorf("ResourceExplorerClient requires a region")
	}

	// Get the list of supported regions for the service
	resourceExplorerRegions, err := listRegionsForService(ctx, d, resourceexplorer2Endpoint.EndpointsID)
	if err != nil {
		return nil, fmt.Errorf("ResourceExplorerClient: failed to get supported regions")
	}

	// Verify the requested region is supported, otherwise return nil which
	// will mean zero results are returned for the query.
	if !slices.Contains(resourceExplorerRegions, region) {
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
	// Route53 is a global service with a single DNS endpoint
	// (route53.amazonaws.com).
	// https://docs.aws.amazon.com/general/latest/gr/r53.html
	// So, while requests will go to the global endpoint, but we can still
	// prefer / reuse the client region.
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return route53.NewFromConfig(*cfg), nil
}

func Route53DomainsClient(ctx context.Context, d *plugin.QueryData) (*route53domains.Client, error) {
	// Route53 Domains is a global service that operates from a single
	// region (route53domains.us-east-1.amazonaws.com).
	// https://docs.aws.amazon.com/general/latest/gr/r53.html
	// So, we must specify the default region rather than the client region.
	cfg, err := getClientForLastResortRegion(ctx, d)
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

	// Depending on their configuration, the S3 client may need to be configured
	// to use path-style addressing.
	awsSpcConfig := GetConfig(d.Connection)
	if awsSpcConfig.S3ForcePathStyle != nil {
		svc = s3.NewFromConfig(*cfg, func(o *s3.Options) {
			o.UsePathStyle = *awsSpcConfig.S3ForcePathStyle
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

// All requests to create or maintain Multi-Region Access Points are routed to the US West (Oregon) Region. so we have no choice but to hard code it here.
// This is true regardless of which Region you are in when making the request, or what Regions the Multi-Region Access Point supports.
// https://docs.aws.amazon.com/AmazonS3/latest/userguide/ManagingMultiRegionAccessPoints.html
// S3 multi-region access point supports in China but not in US Gov or US ISO
func S3ControlMultiRegionAccessClient(ctx context.Context, d *plugin.QueryData) (*s3control.Client, error) {
	cfg, err := getClient(ctx, d, "us-west-2")
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

func SchedulerClient(ctx context.Context, d *plugin.QueryData) (*scheduler.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, schedulerEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return scheduler.NewFromConfig(*cfg), nil
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

func SecurityLakeClient(ctx context.Context, d *plugin.QueryData) (*securitylake.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, securitylakeEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return securitylake.NewFromConfig(*cfg), nil
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

func ServiceCatalogClient(ctx context.Context, d *plugin.QueryData) (*servicecatalog.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return servicecatalog.NewFromConfig(*cfg), nil
}

func ServiceDiscoveryClient(ctx context.Context, d *plugin.QueryData) (*servicediscovery.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return servicediscovery.NewFromConfig(*cfg), nil
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

func SimSpaceWeaverClient(ctx context.Context, d *plugin.QueryData) (*simspaceweaver.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, simspaceWeaverEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return simspaceweaver.NewFromConfig(*cfg), nil
}

func StepFunctionsClient(ctx context.Context, d *plugin.QueryData) (*sfn.Client, error) {
	cfg, err := getClientForQueryRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return sfn.NewFromConfig(*cfg), nil
}

func ShieldClient(ctx context.Context, d *plugin.QueryData) (*shield.Client, error) {
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		return nil, nil
	}
	return shield.NewFromConfig(*cfg), nil
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

func SSMIncidentsClient(ctx context.Context, d *plugin.QueryData) (*ssmincidents.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, ssmIncidentsEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return ssmincidents.NewFromConfig(*cfg), nil
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
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return sts.NewFromConfig(*cfg), nil
}

func SSOAdminClient(ctx context.Context, d *plugin.QueryData) (*ssoadmin.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, ssoEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return ssoadmin.NewFromConfig(*cfg), nil
}

func SupportClient(ctx context.Context, d *plugin.QueryData) (*support.Client, error) {
	// AWS Support is a global service. This means that any endpoint that you use will update your support cases in the Support Center Console.
	// For example, if you use the US East (N. Virginia) endpoint to create a case, you can use the US West (Oregon) or Europe (Ireland) endpoint to add a correspondence to the same case.
	// https://docs.aws.amazon.com/awssupport/latest/user/about-support-api.html#endpoint
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return support.NewFromConfig(*cfg), nil
}

func TransferClient(ctx context.Context, d *plugin.QueryData) (*transfer.Client, error) {
	// AWS Transfer Family
	cfg, err := getClientForQuerySupportedRegion(ctx, d, transferEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return transfer.NewFromConfig(*cfg), nil
}

func TimestreamwriteClient(ctx context.Context, d *plugin.QueryData) (*timestreamwrite.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, timestreamwriteEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return timestreamwrite.NewFromConfig(*cfg), nil
}

func WAFClient(ctx context.Context, d *plugin.QueryData) (*waf.Client, error) {
	// WAF Classic a global service with a single DNS endpoint
	// (waf.amazonaws.com).
	// https://docs.aws.amazon.com/general/latest/gr/waf-classic.html
	// So, while requests will go to the global endpoint, we can still prefer /
	// reuse the client region.
	cfg, err := getClientForDefaultRegion(ctx, d)
	if err != nil {
		return nil, err
	}
	return waf.NewFromConfig(*cfg), nil
}

func WAFRegionalClient(ctx context.Context, d *plugin.QueryData) (*wafregional.Client, error) {
	cfg, err := getClientForQuerySupportedRegion(ctx, d, wafregionalEndpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return wafregional.NewFromConfig(*cfg), nil
}

func WAFV2Client(ctx context.Context, d *plugin.QueryData, region string) (*wafv2.Client, error) {
	var cfg *aws.Config
	var err error
	// For WAFv2 resources of type CloudFront, we are building the region metrix including a region value 'global'.
	// We need to pass the the region value 'us-east-1' to get the cloudfront resource types.
	// getClientForQuerySupportedRegion function removes the invalid region(global) while building the client for which we need the below check
	if region == "global" {
		cfg, err = getClient(ctx, d, "us-east-1")
	} else {
		cfg, err = getClientForQuerySupportedRegion(ctx, d, wafv2Endpoint.EndpointsID)
	}
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
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

// Get a session for the region defined in query data, but only after checking
// it's a supported region for the given serviceID.
func getClientForQuerySupportedRegion(ctx context.Context, d *plugin.QueryData, serviceID string) (*aws.Config, error) {
	return getClientForQuerySupportedRegionWithExclusions(ctx, d, serviceID, []string{})
}

// Get a session for the region defined in query data, but only after checking
// it's a supported region for the given serviceID and that it's not in the
// list of excluded regions. This is useful for cases where the service regions
// in the AWS SDK are actually wrong (e.g. APIGatewayV2 has less regions than
// APIGateway but uses the same service definition.)
func getClientForQuerySupportedRegionWithExclusions(ctx context.Context, d *plugin.QueryData, serviceID string, excludeRegions []string) (*aws.Config, error) {

	// Verify we have good region data
	region := d.EqualsQualString(matrixKeyRegion)
	if region == "" {
		return nil, fmt.Errorf("getClientForQuerySupportedRegion called without a region in QueryData")
	}

	// Work out which regions are valid for this service
	validRegions, err := listRegionsForService(ctx, d, serviceID)
	if err != nil {
		return nil, err
	}

	// Remove the excluded regions from the valid list
	validRegions = helpers.RemoveFromStringSlice(validRegions, excludeRegions...)

	if !slices.Contains(validRegions, region) {
		// We choose to ignore unsupported regions rather than returning an error
		// for them - it's a better user experience. So, return a nil session rather
		// than an error. The caller must handle this case.
		return nil, nil
	}

	// Supported region, so get and return the session
	return getClient(ctx, d, region)
}

// Helper function to get the session for the preferred region in this partition
func getClientForDefaultRegion(ctx context.Context, d *plugin.QueryData) (*aws.Config, error) {
	r, err := getDefaultRegion(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return getClient(ctx, d, r)
}

// Helper function to get the session for the last resort region in this partition
func getClientForLastResortRegion(ctx context.Context, d *plugin.QueryData) (*aws.Config, error) {
	r, err := getLastResortRegion(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return getClient(ctx, d, r)
}

// Helper function to get the session for a region set in query data
func getClientForQueryRegion(ctx context.Context, d *plugin.QueryData) (*aws.Config, error) {
	region := d.EqualsQualString(matrixKeyRegion)
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

// Get the AWS client for a given region. This is cached on a per-connection-region
// basis internally.
func getClient(ctx context.Context, d *plugin.QueryData, region string) (*aws.Config, error) {
	// Create custom hydrate data to pass through the region. Hydrate data
	// is normally per-column, but we can hijack it for this case to pass
	// through the context we need.
	h := &plugin.HydrateData{Item: region}
	i, err := getClientCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	return i.(*aws.Config), nil
}

// Cached form of getClient, using the per-connection and parallel safe
// Memoize() method.
var getClientCached = plugin.HydrateFunc(getClientUncached).Memoize(memoize.WithCacheKeyFunction(getClientCacheKey))

// getClient is per-region, but Memoize() is per-connection, so a setup
// a custom cache key with region information in it.
func getClientCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Extract the region from the hydrate data. This is not per-row data,
	// but a clever pass through of context for our case.
	region := h.Item.(string)
	key := fmt.Sprintf("getClient-%s", region)
	return key, nil
}

// getClientUncached is the actual implementation of getClient, which should
// be run only once per region per connection. Do not call this directly, use
// getClient instead.
func getClientUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Extract the region from the hydrate data. This is not per-row data,
	// but a clever pass through of context for our case.
	region := h.Item.(string)

	plugin.Logger(ctx).Debug("getClientUncached", "connection_name", d.Connection.Name, "region", region, "status", "starting")

	awsSpcConfig := GetConfig(d.Connection)

	// As per the logic used in retryRules of NewConnectionErrRetryer, default to minimum delay of 25ms and maximum
	// number of retries as 9 (our default). The default maximum delay will not be more than approximately 3 minutes to avoid Steampipe
	// waiting too long to return results
	maxRetries := 9
	var minRetryDelay time.Duration = 25 * time.Millisecond // Default minimum delay

	// Set max retry count from config file or env variable (config file has precedence)
	if awsSpcConfig.MaxErrorRetryAttempts != nil {
		maxRetries = *awsSpcConfig.MaxErrorRetryAttempts
	} else if os.Getenv("AWS_MAX_ATTEMPTS") != "" {
		maxRetriesEnvVar, err := strconv.Atoi(os.Getenv("AWS_MAX_ATTEMPTS"))
		if err != nil || maxRetriesEnvVar < 1 {
			panic("invalid value for environment variable \"AWS_MAX_ATTEMPTS\". It should be an integer value greater than or equal to 1")
		}
		maxRetries = maxRetriesEnvVar
	}

	// Set min delay time from config file
	if awsSpcConfig.MinErrorRetryDelay != nil {
		minRetryDelay = time.Duration(*awsSpcConfig.MinErrorRetryDelay) * time.Millisecond
	}

	if maxRetries < 1 {
		panic("connection config has invalid value for \"max_error_retry_attempts\", it must be greater than or equal to 1")
	}
	if minRetryDelay < 1 {
		panic("connection config has invalid value for \"min_error_retry_delay\", it must be greater than or equal to 1")
	}

	sess, err := getClientWithMaxRetries(ctx, d, region, maxRetries, minRetryDelay)
	if err != nil {
		plugin.Logger(ctx).Error("getClientUncached", "region", region, "err", err)
		return nil, err
	}

	plugin.Logger(ctx).Debug("getClientUncached", "connection_name", d.Connection.Name, "region", region, "status", "done")
	return sess, err
}

func getClientWithMaxRetries(ctx context.Context, d *plugin.QueryData, region string, maxRetries int, minRetryDelay time.Duration) (*aws.Config, error) {

	plugin.Logger(ctx).Debug("getClientWithMaxRetries", "connection_name", d.Connection.Name, "region", region, "status", "starting")

	if region == "" {
		return nil, fmt.Errorf("getClientWithMaxRetries called with an empty region")
	}

	// Start with the shared config for the account, and then customize
	// for this specific region etc.
	baseCfg, err := getBaseClientForAccount(ctx, d)
	if err != nil {
		return nil, err
	}
	cfg := baseCfg.Copy()
	plugin.Logger(ctx).Debug("getClientWithMaxRetries", "connection_name", d.Connection.Name, "config_region", cfg.Region, "status", "copy_base_config")

	// Set the region for this client
	// Note: The region set directly in cfg.Region will not be used by the AWS
	// SDK when making background sts:AssumeRole API calls for IAM role
	// authentication. So even if we set a region here but the AWS SDK could not
	// resolve a region and no region was passed into the base config's options,
	// a signing error will be thrown for API calls with this client, e.g.,
	// Error: operation error CloudFront: ListDistributions, failed to sign request: failed to retrieve credentials: failed to refresh cached credentials, operation error STS: AssumeRole, failed to resolve service endpoint, an AWS region is required, but was not found
	cfg.Region = region
	plugin.Logger(ctx).Debug("getClientWithMaxRetries", "connection_name", d.Connection.Name, "config_region", cfg.Region, "status", "set_client_region")

	// Add the retryer definition
	retryer := retry.NewStandard(func(o *retry.StandardOptions) {
		// reseting state of rand to generate different random values
		rand.New(rand.NewSource(time.Now().UnixNano()))
		o.MaxAttempts = maxRetries
		o.MaxBackoff = 5 * time.Minute
		o.RateLimiter = NoOpRateLimit{} // With no rate limiter
		o.Backoff = NewExponentialJitterBackoff(minRetryDelay, maxRetries)
	})
	cfg.Retryer = func() aws.Retryer {
		// UnknownError is the code returned for a 408 from the aws go sdk, these can be frequent on large accounts especially around SNS Topics, etc.
		additionalErrors := []string{"UnknownError"}
		return retry.AddWithErrorCodes(retryer, additionalErrors...)
	}

	// Plugin level config
	awsSpcConfig := GetConfig(d.Connection)

	// If there is a custom endpoint, use it
	var awsEndpointUrl string
	awsEndpointUrl = os.Getenv("AWS_ENDPOINT_URL")
	if awsSpcConfig.EndpointUrl != nil {
		awsEndpointUrl = *awsSpcConfig.EndpointUrl
		if awsEndpointUrl != "" {
			customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           awsEndpointUrl,
					SigningRegion: region,
				}, nil
			})
			newCfg, err := config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(customResolver))
			if err != nil {
				plugin.Logger(ctx).Error("service.getClientWithMaxRetries", "connection_error", err)
				return nil, err
			}
			newCfg.Retryer = cfg.Retryer
			newCfg.Region = cfg.Region
			cfg = newCfg
		}
	}

	plugin.Logger(ctx).Debug("getClientWithMaxRetries", "connection_name", d.Connection.Name, "region", region, "status", "done")

	return &cfg, err
}

// Helper function to get an AWS config object for each connection. This object
// is then copied and shared across regions. This approach avoids unnecssary
// creation work for sessions, particularly when using a shared service like IDMS.
// The client is actually cached indefinitely which is desirable since the AWS
// SDK will automatically refresh as needed. We used to time this cache out
// (every 5 mins), but that caused a lof of session resets and termination of
// things like SSO sessions before they expired.
// Previously we'd create a new session for each region, but this leads to
// throttling on the IDMS service - consider 10 connections with 10 regions, that's
// 100 sessions which leads to 300 IDMS calls in very quick succession. It was
// even worse when the cache was not safe against parallel runs, causing many to
// be recreated. Using a base client creation and combining with the safety of
// Memoize() is a much better approach.
func getBaseClientForAccount(ctx context.Context, d *plugin.QueryData) (*aws.Config, error) {
	tmp, err := getBaseClientForAccountCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return tmp.(*aws.Config), nil
}

// Initialize a single HTTP client that is optimized for Steampipe and shared
// across all AWS SDK clients. We have hundreds of AWS SDK clients (one per
// account region) that are all sharing this same HTTP client - creating shared
// caching and controls over parallelism.
//
// The AWS SDK defaults are good, but not great for our highly parallel use in
// Steampipe. Specific problems this client aims to solve:
// 1. DNS floods - performing thousands of simultaneous API calls creates a DNS
// lookup for each one (even if the same domain). This can overwhelm the DNS
// server and cause "no such host" errors.
// 2. HTTP connection floods - the AWS SDK defaults to no limit on the number of
// HTTP connections per host. Thousands of connections created simultaneously to
// the same host is hard on both the client and the target server.
// 3. DNS caching - Golang does not cache DNS lookups by default. We end up
// looking up the same host thousands of times both within a query and across
// queries.
func initializeHTTPClient() aws.HTTPClient {

	// DNS lookup floods are a real problem with highly parallel AWS SDK calls. Every
	// API request leads to a DNS lookup by default (since Go doesn't cache them). We
	// employ a DNS lookup cache, but we also need to limit the number of parallel DNS
	// requests to avoid overwhelming the underlying DNS server. For example, listing
	// S3 buckets will create 2 DNS lookup requests per bucket which is a lot of
	// pressure on the DNS layer of your network.
	// This setting will limit the number of parallel DNS lookups. An appropriate setting
	// depends on the capabilities of your DNS server. The default is 25, which is low
	// enough for a Macbook M1 to work without "no such host" errors when using the cgo
	// network stack. It's high enough to work great in most cases, except maybe massive
	// S3 bucket listing (which is rare). Notably on the same Macbook M1, when the plugin
	// is compiled using netgo (our default on Mac) DNS lookups will succeed with virtually
	// no upper limit on this setting. So, bottom line, 25 is a guess to try and ensure
	// it works reliably and optimally enough.
	dnsLookupMaxParallel := readEnvVarToInt("STEAMPIPE_AWS_DNS_LOOKUP_MAX_PARALLEL", 25)

	// The DNS cache will be refreshed at this interval. A refresh means that
	// any unused entries are removed and any entries that were used since the
	// last refresh will be re-looked up to ensure they are current.
	// This setting should be large enough to get the benefit of caching and short
	// enough to prevent stale entries from being used for too long.
	// Set to 0 to disable the refresh completely (not a good idea).
	// Set to -1 to disable the DNS cache completely (the AWS default).
	dnsCacheRefreshIntervalSecs := readEnvVarToInt("STEAMPIPE_AWS_DNS_CACHE_REFRESH_INTERVAL_SECS", 300)

	// This is the maximum number of HTTPS API connections used for each host
	// (e.g.  iam.amazonaws.com). We want a number that is high enough to do a
	// lot of parallel work, but not so high that we have an excess number of
	// sockets open.
	// There is a trade off here. Tables like S3 have a lot of hosts - i.e. two
	// per bucket (one for the central region to get the creation time and one
	// for the actual bucket region), while services like IAM use a single host
	// for all queries.
	// Set to 0 to remove the limit (which is the AWS SDK default).
	httpTransportMaxConnsPerHost := readEnvVarToInt("STEAMPIPE_AWS_HTTP_TRANSPORT_MAX_CONNS_PER_HOST", 5000)

	// Our DNS resolver should automatically refresh itself on this schedule.
	var resolver = &dnscache.Resolver{}
	if dnsCacheRefreshIntervalSecs > 0 {
		go func() {
			t := time.NewTicker(time.Duration(dnsCacheRefreshIntervalSecs) * time.Second)
			defer t.Stop()
			for range t.C {
				resolver.Refresh(true)
			}
		}()
	}

	// The AWS SDK has a special "buildable" HTTP client so it can be combined
	// with specific options such as custom certificate bundles. It matches the
	// interface of a HTTPClient, but has specific approaches for setting
	// transport options etc. Our goal is to use the default AWS settings (e.g.
	// timeouts, etc) as much as possible and just override the specific
	// behavior of parallelism for DNS lookups and HTTP requests.
	client := awshttp.NewBuildableClient()

	// Limit the max connections per host, but only if set. The AWS SDK default
	// is no limit.
	if httpTransportMaxConnsPerHost > 0 {
		client = client.WithTransportOptions(func(tr *http.Transport) {
			tr.MaxConnsPerHost = httpTransportMaxConnsPerHost
		})
	}

	// Use a DNS cache if it's set, otherwise we just avoid changing the dialer behavior
	// of the AWS HTTP client.
	if dnsCacheRefreshIntervalSecs >= 0 {

		// A semaphore is used to control the number of parallel DNS lookups.
		sem := semaphore.NewWeighted(int64(dnsLookupMaxParallel))

		// A dialer for testing connections
		dialer := client.GetDialer()

		client = client.WithTransportOptions(func(tr *http.Transport) {
			tr.DialContext = func(ctx context.Context, network string, addr string) (conn net.Conn, err error) {

				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				// Acquire a semaphore slot, blocking until one is available.
				if err := sem.Acquire(ctx, 1); err != nil {
					return nil, err
				}

				// Actually resolve the host, using a cached result if possible.
				// Returns an array of IPs for the host.
				ips, err := resolver.LookupHost(ctx, host)

				// Release the semaphore, even if there was an error.
				sem.Release(1)

				// If there was an error during lookup, we give up immediately.
				if err != nil {
					return nil, err
				}

				// Now, look through the IP addresses until we manage to create a good connection.
				// This is less optimal than the parallelized native golang approach, but good
				// enough and much simpler. Comparison - https://cs.opensource.google/go/go/+/refs/tags/go1.21.5:src/net/dial.go;l=454-507
				for _, ip := range ips {
					conn, err = dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
					if err == nil {
						break
					}
				}

				return
			}
		})
	}

	return client
}

var sharedHTTPClient = initializeHTTPClient()

// Cached form of the base client.
// This cache HAS A 30 DAY EXPIRATION! This is because the AWS SDK will
// automatically refresh credentials as needed from this cached object.
// If we expire the cache regularly we are causing SSO sessions to end
// prematurely, and causing the AWS SDK to refresh credentials more often
// using the IDMS service etc.
var getBaseClientForAccountCached = plugin.HydrateFunc(getBaseClientForAccountUncached).Memoize(memoize.WithTtl(time.Hour * 24 * 30))

// Do the actual work of creating an AWS config object for reuse across many
// regions. This client has the minimal reusable configuration on it, so it
// can be modified in the higher level client functions.
func getBaseClientForAccountUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	plugin.Logger(ctx).Debug("getBaseClientForAccountUncached", "connection_name", d.Connection.Name, "status", "starting")

	awsSpcConfig := GetConfig(d.Connection)

	var configOptions []func(*config.LoadOptions) error

	// Note about region config: We deliberately do not set a region when
	// creating the config. It will load the best region for this account
	// (profile) from the ~/.aws shared config / env vars / etc. This is
	// important, as this base client is even used when trying to guess the
	// default region for the user based on these settings.

	if awsSpcConfig.Profile != nil {
		profile := aws.ToString(awsSpcConfig.Profile)
		plugin.Logger(ctx).Debug("getBaseClientForAccountUncached", "connection_name", d.Connection.Name, "status", "profile_found", "profile", profile)
		configOptions = append(configOptions, config.WithSharedConfigProfile(profile))
	}

	if awsSpcConfig.AccessKey != nil && awsSpcConfig.SecretKey == nil {
		return nil, fmt.Errorf("partial credentials found in connection config, missing: secret_key")
	} else if awsSpcConfig.SecretKey != nil && awsSpcConfig.AccessKey == nil {
		return nil, fmt.Errorf("partial credentials found in connection config, missing: access_key")
	} else if awsSpcConfig.AccessKey != nil && awsSpcConfig.SecretKey != nil {
		plugin.Logger(ctx).Debug("getBaseClientForAccountUncached", "connection_name", d.Connection.Name, "status", "key_pair_found")
		sessionToken := ""
		if awsSpcConfig.SessionToken != nil {
			plugin.Logger(ctx).Debug("getBaseClientForAccountUncached", "connection_name", d.Connection.Name, "status", "session_token_found")
			sessionToken = *awsSpcConfig.SessionToken
		}
		provider := credentials.NewStaticCredentialsProvider(*awsSpcConfig.AccessKey, *awsSpcConfig.SecretKey, sessionToken)
		configOptions = append(configOptions, config.WithCredentialsProvider(provider))
	}

	plugin.Logger(ctx).Debug("getBaseClientForAccountUncached", "connection_name", d.Connection.Name, "status", "loading_config")
	if plugin.Logger(ctx).GetLevel() <= hclog.Debug {
		logger := plugin.Logger(ctx)
		configOptions = append(configOptions, config.WithLogger(NewHCLoggerToSmithyLoggerWrapper(&logger)))
		configOptions = append(configOptions, config.WithClientLogMode(aws.LogRetries))
	}

	// NOTE: EC2 metadata service IMDS throttling and retries
	//
	// With a large number of connections all being setup on a single machine using
	// IMDS credentials, we can hit the IMDS throttling limits. We only query IMDS
	// once per connection (3 API calls under the hood), but it still throttles once
	// over 200 connections or so (estimate, rate limits vary).
	//
	// I was unable to find a way to setup automatic retries and the information
	// available online as of 2023-01-26 is limited. Best links I could find:
	// * IDMS service - https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/ec2/imds#pkg-overview
	// * (Broken) example with ec2rolecreds - https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds#Provider
	// * (Fixed) issue that it didn't work - https://github.com/aws/aws-sdk-go-v2/issues/1296
	//
	// // This code ran, but didn't seem to respect the idms.Options{...}
	// // The debugHTTPClient did work to show requests / throttling.
	// retryer := retry.NewStandard(func(o *retry.StandardOptions) {
	//   // reseting state of rand to generate different random values
	//   rand.Seed(time.Now().UnixNano())
	//   o.MaxAttempts = 50
	//   o.MaxBackoff = 5 * time.Minute
	//   o.RateLimiter = NoOpRateLimit{} // With no rate limiter
	//   o.Backoff = NewExponentialJitterBackoff(100*time.Millisecond, 5)
	//   log.Printf("[WARN] retryer!")
	// })
	// configOptions = append(configOptions, config.WithEC2RoleCredentialOptions(func(opts *ec2rolecreds.Options) {
	//   // debugHTTPClient per https://github.com/aws/aws-sdk-go-v2/issues/1296
	//   opts.Client = imds.New(imds.Options{Retryer: retryer, ClientLogMode: aws.LogRetries | aws.LogRequest}, withDebugHTTPClient())
	// }))

	configOptions = append(configOptions, config.WithHTTPClient(sharedHTTPClient))

	cfg, err := config.LoadDefaultConfig(ctx, configOptions...)
	if err != nil {
		plugin.Logger(ctx).Error("getBaseClientForAccountUncached", "connection_name", d.Connection.Name, "load_default_config_error", err)
		return nil, err
	}

	// Even though we create a client per region and set the region during that
	// step, we need to pass a region in the config options if the AWS SDK could
	// not resolve a region from environment variables or the AWS config.
	// This region is used by the AWS SDK when making background sts:AssumeRole
	// API calls for IAM role authentication; if it's not set here, a signing
	// error is thrown for API calls with this client, e.g.,
	// Error: operation error CloudFront: ListDistributions, failed to sign request: failed to retrieve credentials: failed to refresh cached credentials, operation error STS: AssumeRole, failed to resolve service endpoint, an AWS region is required, but was not found
	if cfg.Region == "" {
		defaultRegion, err := getDefaultRegionFromConfig(ctx, d, nil)
		if err != nil {
			plugin.Logger(ctx).Error("getBaseClientForAccountUncached", "connection_name", d.Connection.Name, "get_default_region_error", err)
			return nil, err
		}

		plugin.Logger(ctx).Debug("getBaseClientForAccountUncached", "connection_name", d.Connection.Name, "region", defaultRegion, "status", "set_default_region")
		configOptions = append(configOptions, config.WithRegion(defaultRegion))
		cfg, err = config.LoadDefaultConfig(ctx, configOptions...)
		if err != nil {
			plugin.Logger(ctx).Error("getBaseClientForAccountUncached", "connection_name", d.Connection.Name, "load_default_config_error", err)
			return nil, err
		}
	}

	plugin.Logger(ctx).Debug("getBaseClientForAccountUncached", "connection_name", d.Connection.Name, "status", "done")

	return &cfg, err

}

// HCLoggerToSmithyLoggerWrapper wraps an hclog Logger in order to pass it as an AWS SDK smithy Logger
type HCLoggerToSmithyLoggerWrapper struct {
	hclogger *hclog.Logger
}

func (logger *HCLoggerToSmithyLoggerWrapper) Logf(classification logging.Classification, format string, v ...interface{}) {
	(*logger.hclogger).Debug(fmt.Sprintf(format, v...))
}

func NewHCLoggerToSmithyLoggerWrapper(l *hclog.Logger) *HCLoggerToSmithyLoggerWrapper {
	return &HCLoggerToSmithyLoggerWrapper{l}
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

	// Low level method to log retries since we don't have context etc here.
	// Logging is helpful for visibility into retries and choke points in using
	// the API.
	log.Printf("[INFO] BackoffDelay: attempt=%d, retryTime=%s, err=%v", attempt, retryTime.String(), err)

	return retryTime, nil
}
