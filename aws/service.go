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
	"github.com/aws/aws-sdk-go/service/accessanalyzer"
	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/auditmanager"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/aws/aws-sdk-go/service/cloudcontrolapi"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/databasemigrationservice"
	"github.com/aws/aws-sdk-go/service/dax"
	"github.com/aws/aws-sdk-go/service/directoryservice"
	"github.com/aws/aws-sdk-go/service/dlm"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecrpublic"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/fsx"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/globalaccelerator"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/identitystore"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	"github.com/aws/aws-sdk-go/service/kinesisvideo"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/macie2"
	"github.com/aws/aws-sdk-go/service/mediastore"
	"github.com/aws/aws-sdk-go/service/neptune"
	"github.com/aws/aws-sdk-go/service/networkfirewall"
	"github.com/aws/aws-sdk-go/service/opensearchservice"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/pinpoint"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/aws/aws-sdk-go/service/ram"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53domains"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/aws/aws-sdk-go/service/serverlessapplicationrepository"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/aws/aws-sdk-go/service/wellarchitected"
	"github.com/aws/aws-sdk-go/service/workspaces"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func AccessAnalyzerService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*accessanalyzer.AccessAnalyzer, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, accessanalyzer.EndpointsID)
	if err != nil {
		return nil, err
	}
	return accessanalyzer.New(sess), nil
}

func AmplifyService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*amplify.Amplify, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, amplify.EndpointsID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	return amplify.New(sess), nil
}

func ApplicationAutoScalingService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*applicationautoscaling.ApplicationAutoScaling, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, applicationautoscaling.EndpointsID)
	if err != nil {
		return nil, err
	}
	return applicationautoscaling.New(sess), nil
}

func AuditManagerService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*auditmanager.AuditManager, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, auditmanager.EndpointsID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	return auditmanager.New(sess), nil
}

func AutoScalingService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*autoscaling.AutoScaling, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, autoscaling.EndpointsID)
	if err != nil {
		return nil, err
	}
	return autoscaling.New(sess), nil
}

func BackupService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*backup.Backup, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, backup.EndpointsID)
	if err != nil {
		return nil, err
	}
	return backup.New(sess), nil
}

func CloudControlService(ctx context.Context, d *plugin.QueryData) (*cloudcontrolapi.CloudControlApi, error) {
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
	serviceCacheKey := fmt.Sprintf("cloudcontrolapi-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudcontrolapi.CloudControlApi), nil
	}

	sess, err := getSessionWithMaxRetries(ctx, d, region, 4, 25*time.Millisecond)
	if err != nil {
		return nil, err
	}
	svc := cloudcontrolapi.New(sess)

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

func CodeBuildService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*codebuild.CodeBuild, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, codebuild.EndpointsID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	return codebuild.New(sess), nil
}

func CodeCommitService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*codecommit.CodeCommit, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, codecommit.EndpointsID)
	if err != nil {
		return nil, err
	}
	return codecommit.New(sess), nil
}

func CodePipelineService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*codepipeline.CodePipeline, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, codepipeline.EndpointsID)
	if err != nil {
		return nil, err
	}
	return codepipeline.New(sess), nil
}

func CloudFrontService(ctx context.Context, d *plugin.QueryData) (*cloudfront.CloudFront, error) {
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return cloudfront.New(sess), nil
}

func CloudFormationService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*cloudformation.CloudFormation, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, cloudformation.EndpointsID)
	if err != nil {
		return nil, err
	}
	return cloudformation.New(sess), nil
}

func CloudWatchService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*cloudwatch.CloudWatch, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, cloudwatch.EndpointsID)
	if err != nil {
		return nil, err
	}
	return cloudwatch.New(sess), nil
}

func CloudWatchLogsService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*cloudwatchlogs.CloudWatchLogs, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, cloudwatchlogs.EndpointsID)
	if err != nil {
		return nil, err
	}
	return cloudwatchlogs.New(sess), nil
}

func CloudTrailService(ctx context.Context, d *plugin.QueryData, region string) (*cloudtrail.CloudTrail, error) {
	sess, err := getSessionForRegion(ctx, d, region)
	if err != nil {
		return nil, err
	}
	return cloudtrail.New(sess), nil
}

func ConfigService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*configservice.ConfigService, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, configservice.EndpointsID)
	if err != nil {
		return nil, err
	}
	return configservice.New(sess), nil
}

func CostExplorerService(ctx context.Context, d *plugin.QueryData) (*costexplorer.CostExplorer, error) {
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return costexplorer.New(sess), nil
}

func DAXService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*dax.DAX, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, dax.EndpointsID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	return dax.New(sess), nil
}

func DatabaseMigrationService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*databasemigrationservice.DatabaseMigrationService, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, databasemigrationservice.EndpointsID)
	if err != nil {
		return nil, err
	}
	return databasemigrationservice.New(sess), nil
}

func DirectoryService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*directoryservice.DirectoryService, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, directoryservice.EndpointsID)
	if err != nil {
		return nil, err
	}
	return directoryservice.New(sess), nil
}

func DLMService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*dlm.DLM, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, dlm.EndpointsID)
	if err != nil {
		return nil, err
	}
	return dlm.New(sess), nil
}

func DynamoDBService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*dynamodb.DynamoDB, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, dynamodb.EndpointsID)
	if err != nil {
		return nil, err
	}
	return dynamodb.New(sess), nil
}

func EC2Service(ctx context.Context, d *plugin.QueryData, region string) (*ec2.EC2, error) {
	sess, err := getSessionForRegion(ctx, d, region)
	if err != nil {
		return nil, err
	}
	return ec2.New(sess), nil
}

func Ec2RegionsService(ctx context.Context, d *plugin.QueryData, region string) (*ec2.EC2, error) {
	// We can query EC2 for the list of supported regions. But, if credentials
	// are insufficient this query will retry many times, so we create a special
	// client with a small number of retries to prevent hangs.
	// Note - This is not cached, but usually the result of using this service will be.
	sess, err := getSessionWithMaxRetries(ctx, d, region, 4, 25*time.Millisecond)
	if err != nil {
		return nil, err
	}
	svc := ec2.New(sess)
	return svc, nil
}

func ECRService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*ecr.ECR, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, ecr.EndpointsID)
	if err != nil {
		return nil, err
	}
	return ecr.New(sess), nil
}

func EcrPublicService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*ecrpublic.ECRPublic, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, ecrpublic.EndpointsID)
	if err != nil {
		return nil, err
	}
	return ecrpublic.New(sess), nil
}

func ECSService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*ecs.ECS, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, ecs.EndpointsID)
	if err != nil {
		return nil, err
	}
	return ecs.New(sess), nil
}

func EFSService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*efs.EFS, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, efs.EndpointsID)
	if err != nil {
		return nil, err
	}
	return efs.New(sess), nil
}

func FSxService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*fsx.FSx, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, fsx.EndpointsID)
	if err != nil {
		return nil, err
	}
	return fsx.New(sess), nil
}

func EKSService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*eks.EKS, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, eks.EndpointsID)
	if err != nil {
		return nil, err
	}
	return eks.New(sess), nil
}

func ElasticBeanstalkService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*elasticbeanstalk.ElasticBeanstalk, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, elasticbeanstalk.EndpointsID)
	if err != nil {
		return nil, err
	}
	return elasticbeanstalk.New(sess), nil
}

func ElastiCacheService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*elasticache.ElastiCache, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, elasticache.EndpointsID)
	if err != nil {
		return nil, err
	}
	return elasticache.New(sess), nil
}

func ElasticsearchService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*elasticsearchservice.ElasticsearchService, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, elasticsearchservice.EndpointsID)
	if err != nil {
		return nil, err
	}
	return elasticsearchservice.New(sess), nil
}

func ELBv2Service(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*elbv2.ELBV2, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, elbv2.EndpointsID)
	if err != nil {
		return nil, err
	}
	return elbv2.New(sess), nil
}

func EventBridgeService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*eventbridge.EventBridge, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, eventbridge.EndpointsID)
	if err != nil {
		return nil, err
	}
	return eventbridge.New(sess), nil
}

func EMRService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*emr.EMR, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, emr.EndpointsID)
	if err != nil {
		return nil, err
	}
	return emr.New(sess), nil
}

func FirehoseService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*firehose.Firehose, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, firehose.EndpointsID)
	if err != nil {
		return nil, err
	}
	return firehose.New(sess), nil
}

func GlacierService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*glacier.Glacier, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, glacier.EndpointsID)
	if err != nil {
		return nil, err
	}
	return glacier.New(sess), nil
}

func GlobalAcceleratorService(ctx context.Context, d *plugin.QueryData) (*globalaccelerator.GlobalAccelerator, error) {
	// Global Accelerator is a global service that supports endpoints in multiple AWS Regions but you must specify
	// the us-west-2 (Oregon) Region to create or update accelerators.
	sess, err := getSession(ctx, d, "us-west-2")
	if err != nil {
		return nil, err
	}
	return globalaccelerator.New(sess), nil
}

func GlueService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*glue.Glue, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, glue.EndpointsID)
	if err != nil {
		return nil, err
	}
	return glue.New(sess), nil
}

func GuardDutyService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*guardduty.GuardDuty, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, guardduty.EndpointsID)
	if err != nil {
		return nil, err
	}
	return guardduty.New(sess), nil
}

func IAMService(ctx context.Context, d *plugin.QueryData) (*iam.IAM, error) {
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return iam.New(sess), nil
}

func IdentityStoreService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*identitystore.IdentityStore, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, identitystore.EndpointsID)
	if err != nil {
		return nil, err
	}
	return identitystore.New(sess), nil
}

func InspectorService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*inspector.Inspector, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, inspector.EndpointsID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	return inspector.New(sess), nil
}

func KinesisService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*kinesis.Kinesis, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, kinesis.EndpointsID)
	if err != nil {
		return nil, err
	}
	return kinesis.New(sess), nil
}

func KinesisAnalyticsV2Service(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*kinesisanalyticsv2.KinesisAnalyticsV2, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, kinesisanalyticsv2.EndpointsID)
	if err != nil {
		return nil, err
	}
	return kinesisanalyticsv2.New(sess), nil
}

func KinesisVideoService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*kinesisvideo.KinesisVideo, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, kinesisvideo.EndpointsID)
	if err != nil {
		return nil, err
	}
	return kinesisvideo.New(sess), nil
}

func KMSService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*kms.KMS, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, kms.EndpointsID)
	if err != nil {
		return nil, err
	}
	return kms.New(sess), nil
}

func LambdaService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*lambda.Lambda, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, lambda.EndpointsID)
	if err != nil {
		return nil, err
	}
	return lambda.New(sess), nil
}

func Macie2Service(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*macie2.Macie2, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, macie2.EndpointsID)
	if err != nil {
		return nil, err
	}
	return macie2.New(sess), nil
}

func MediaStoreService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*mediastore.MediaStore, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, mediastore.EndpointsID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	return mediastore.New(sess), nil
}

func NeptuneService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*neptune.Neptune, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, neptune.EndpointsID)
	if err != nil {
		return nil, err
	}
	return neptune.New(sess), nil
}

func NetworkFirewallService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*networkfirewall.NetworkFirewall, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, networkfirewall.EndpointsID)
	if err != nil {
		return nil, err
	}
	return networkfirewall.New(sess), nil
}

func PinpointService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*pinpoint.Pinpoint, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, pinpoint.EndpointsID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	return pinpoint.New(sess), nil
}

func OpenSearchService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*opensearchservice.OpenSearchService, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, opensearchservice.EndpointsID)
	if err != nil {
		return nil, err
	}
	return opensearchservice.New(sess), nil
}

func OrganizationService(ctx context.Context, d *plugin.QueryData) (*organizations.Organizations, error) {
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return organizations.New(sess), nil
}

func PricingService(ctx context.Context, d *plugin.QueryData) (*pricing.Pricing, error) {
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return pricing.New(sess), nil
}

func RAMService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*ram.RAM, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, ram.EndpointsID)
	if err != nil {
		return nil, err
	}
	return ram.New(sess), nil
}

func RDSService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*rds.RDS, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, rds.EndpointsID)
	if err != nil {
		return nil, err
	}
	return rds.New(sess), nil
}

func RedshiftService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*redshift.Redshift, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, redshift.EndpointsID)
	if err != nil {
		return nil, err
	}
	return redshift.New(sess), nil
}

func Route53DomainsService(ctx context.Context, d *plugin.QueryData) (*route53domains.Route53Domains, error) {
	// Route53Domains only operate in us-east-1
	sess, err := getSession(ctx, d, "us-east-1")
	if err != nil {
		return nil, err
	}
	return route53domains.New(sess), nil
}

func Route53ResolverService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*route53resolver.Route53Resolver, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, route53resolver.EndpointsID)
	if err != nil {
		return nil, err
	}
	return route53resolver.New(sess), nil
}

func Route53Service(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*route53.Route53, error) {
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return route53.New(sess), nil
}

func SecretsManagerService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*secretsmanager.SecretsManager, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, secretsmanager.EndpointsID)
	if err != nil {
		return nil, err
	}
	return secretsmanager.New(sess), nil
}

func SecurityHubService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*securityhub.SecurityHub, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, securityhub.EndpointsID)
	if err != nil {
		return nil, err
	}
	return securityhub.New(sess), nil
}

func S3ControlService(ctx context.Context, d *plugin.QueryData, region string) (*s3control.S3Control, error) {
	sess, err := getSessionForRegion(ctx, d, region)
	if err != nil {
		return nil, err
	}
	return s3control.New(sess), nil
}

func SageMakerService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*sagemaker.SageMaker, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, sagemaker.EndpointsID)
	if err != nil {
		return nil, err
	}
	return sagemaker.New(sess), nil
}

func ServerlessApplicationRepositoryService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*serverlessapplicationrepository.ServerlessApplicationRepository, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, serverlessapplicationrepository.EndpointsID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	return serverlessapplicationrepository.New(sess), nil
}

func SESService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*ses.SES, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, ses.EndpointsID)
	if err != nil {
		return nil, err
	}
	return ses.New(sess), nil
}

func SNSService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*sns.SNS, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, sns.EndpointsID)
	if err != nil {
		return nil, err
	}
	return sns.New(sess), nil
}

// TODO
// ServiceQuotasService returns the service connection for AWS ServiceQuotas service
func ServiceQuotasService(ctx context.Context, d *plugin.QueryData) (*servicequotas.ServiceQuotas, error) {
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("servicequotas-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*servicequotas.ServiceQuotas), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, "")
	if err != nil {
		return nil, err
	}
	svc := servicequotas.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

func ServiceQuotasRegionalService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*servicequotas.ServiceQuotas, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, servicequotas.EndpointsID)
	if err != nil {
		return nil, err
	}
	return servicequotas.New(sess), nil
}

func SQSService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*sqs.SQS, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, sqs.EndpointsID)
	if err != nil {
		return nil, err
	}
	return sqs.New(sess), nil
}

func SSMService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*ssm.SSM, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, ssm.EndpointsID)
	if err != nil {
		return nil, err
	}
	return ssm.New(sess), nil
}

func SSOAdminService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*ssoadmin.SSOAdmin, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, ssoadmin.EndpointsID)
	if err != nil {
		return nil, err
	}
	return ssoadmin.New(sess), nil
}

func StepFunctionsService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*sfn.SFN, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, sfn.EndpointsID)
	if err != nil {
		return nil, err
	}
	return sfn.New(sess), nil
}

func STSService(ctx context.Context, d *plugin.QueryData) (*sts.STS, error) {
	// TODO - Should STS be regional instead?
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return sts.New(sess), nil
}

func TaggingResourceService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*resourcegroupstaggingapi.ResourceGroupsTaggingAPI, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, resourcegroupstaggingapi.EndpointsID)
	if err != nil {
		return nil, err
	}
	return resourcegroupstaggingapi.New(sess), nil
}

func WAFService(ctx context.Context, d *plugin.QueryData) (*waf.WAF, error) {
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	return waf.New(sess), nil
}

func WAFRegionalService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*wafregional.WAFRegional, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, wafregional.EndpointsID)
	if err != nil {
		return nil, err
	}
	return wafregional.New(sess), nil
}

func WAFv2Service(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, region string) (*wafv2.WAFV2, error) {
	validRegions := SupportedRegionsForService(ctx, d, h, wafv2.EndpointsID)
	if !helpers.StringSliceContains(validRegions, region) {
		// We choose to ignore unsupported regions rather than returning an error
		// for them - it's a better user experience. So, return a nil session rather
		// than an error. The caller must handle this case.
		return nil, nil
	}
	sess, err := getSessionForRegion(ctx, d, region)
	if err != nil {
		return nil, err
	}
	return wafv2.New(sess), nil
}

func WellArchitectedService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*wellarchitected.WellArchitected, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, wellarchitected.EndpointsID)
	if err != nil {
		return nil, err
	}
	return wellarchitected.New(sess), nil
}

func WorkspacesService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*workspaces.WorkSpaces, error) {
	sess, err := getSessionForQuerySupportedRegion(ctx, d, h, workspaces.EndpointsID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	return workspaces.New(sess), nil
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
func getSessionForQuerySupportedRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, serviceID string) (*session.Session, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	if region == "" {
		return nil, fmt.Errorf("getSessionForQueryRegion called without a region in QueryData")
	}
	validRegions := SupportedRegionsForService(ctx, d, h, serviceID)
	if !helpers.StringSliceContains(validRegions, region) {
		// We choose to ignore unsupported regions rather than returning an error
		// for them - it's a better user experience. So, return a nil session rather
		// than an error. The caller must handle this case.
		return nil, nil
	}
	// Supported region, so get and return the session
	return getSession(ctx, d, region)
}

// Helper function to get the session for a region set in query data
// func getSessionForQueryRegion(ctx context.Context, d *plugin.QueryData) (*session.Session, error) {
// 	region := d.KeyColumnQualString(matrixKeyRegion)
// 	if region == "" {
// 		return nil, fmt.Errorf("getSessionForQueryRegion called without a region in QueryData")
// 	}
// 	return getSession(ctx, d, region)
// }

// Helper function to get the session for a specific region
func getSessionForRegion(ctx context.Context, d *plugin.QueryData, region string) (*session.Session, error) {
	if region == "" {
		return nil, fmt.Errorf("getSessionForRegion called with an empty region")
	}
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
