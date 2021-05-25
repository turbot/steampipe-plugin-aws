package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/accessanalyzer"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/auditmanager"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/databasemigrationservice"
	"github.com/aws/aws-sdk-go/service/dax"
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
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	"github.com/aws/aws-sdk-go/service/kinesisvideo"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/aws/aws-sdk-go/service/wellarchitected"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// AccessAnalyzerService returns the service connection for AWS IAM Access Analyzer service
func AccessAnalyzerService(ctx context.Context, d *plugin.QueryData, region string) (*accessanalyzer.AccessAnalyzer, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed AccessAnalyzerService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("accessanalyzer-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*accessanalyzer.AccessAnalyzer), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := accessanalyzer.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// ACMService returns the service connection for AWS ACM service
func ACMService(ctx context.Context, d *plugin.QueryData, region string) (*acm.ACM, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed ACMService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("acm-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*acm.ACM), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := acm.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// APIGatewayService returns the service connection for AWS API Gateway service
func APIGatewayService(ctx context.Context, d *plugin.QueryData, region string) (*apigateway.APIGateway, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed APIGateway")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("apigateway-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*apigateway.APIGateway), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := apigateway.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// APIGatewayV2Service returns the service connection for AWS API Gateway V2 service
func APIGatewayV2Service(ctx context.Context, d *plugin.QueryData, region string) (*apigatewayv2.ApiGatewayV2, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed APIGatewayV2Service")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("apigatewayv2-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*apigatewayv2.ApiGatewayV2), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := apigatewayv2.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// AutoScalingService returns the service connection for AWS AutoScaling service
func AutoScalingService(ctx context.Context, d *plugin.QueryData, region string) (*autoscaling.AutoScaling, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed AutoScalingService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("autoscaling-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*autoscaling.AutoScaling), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := autoscaling.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// ApplicationAutoScalingService returns the service connection for AWS Application Auto Scaling service
func ApplicationAutoScalingService(ctx context.Context, d *plugin.QueryData, region string) (*applicationautoscaling.ApplicationAutoScaling, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed ApplicationAutoScalingService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("applicationautoscaling-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*applicationautoscaling.ApplicationAutoScaling), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := applicationautoscaling.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// AuditManagerService returns the service connection for AWS Audit Manager service
func AuditManagerService(ctx context.Context, d *plugin.QueryData, region string) (*auditmanager.AuditManager, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed AuditManagerService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("auditmanager-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*auditmanager.AuditManager), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := auditmanager.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// BackupService returns the service connection for AWS Backup service
func BackupService(ctx context.Context, d *plugin.QueryData, region string) (*backup.Backup, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed BackupService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("backup-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*backup.Backup), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := backup.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CodeBuildService returns the service connection for AWS CodeBuild service
func CodeBuildService(ctx context.Context, d *plugin.QueryData, region string) (*codebuild.CodeBuild, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed CodeBuildService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("codebuild-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*codebuild.CodeBuild), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := codebuild.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CloudFrontService returns the service connection for AWS CloudFront service
func CloudFrontService(ctx context.Context, d *plugin.QueryData) (*cloudfront.CloudFront, error) {
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("cloudfront")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudfront.CloudFront), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	svc := cloudfront.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CloudFormationService returns the service connection for AWS CloudFormation service
func CloudFormationService(ctx context.Context, d *plugin.QueryData, region string) (*cloudformation.CloudFormation, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed CloudFormationService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("cloudformation-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudformation.CloudFormation), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := cloudformation.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// CloudWatchService returns the service connection for AWS Cloud Watch service
func CloudWatchService(ctx context.Context, d *plugin.QueryData, region string) (*cloudwatch.CloudWatch, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed CloudWatchService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("cloudwatch-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudwatch.CloudWatch), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := cloudwatch.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// CloudWatchLogsService returns the service connection for AWS Cloud Watch Logs service
func CloudWatchLogsService(ctx context.Context, d *plugin.QueryData, region string) (*cloudwatchlogs.CloudWatchLogs, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed CloudWatchLogsService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("cloudwatchlogs-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudwatchlogs.CloudWatchLogs), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := cloudwatchlogs.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// CloudTrailService returns the service connection for AWS CloudTrail service
func CloudTrailService(ctx context.Context, d *plugin.QueryData, region string) (*cloudtrail.CloudTrail, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed CloudTrailService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("cloudtrail-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudtrail.CloudTrail), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := cloudtrail.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// DaxService returns the service connection for AWS DAX service
func DaxService(ctx context.Context, d *plugin.QueryData, region string) (*dax.DAX, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed DaxService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("dax-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*dax.DAX), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := dax.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// DatabaseMigrationService returns the service connection for AWS Database Migration service
func DatabaseMigrationService(ctx context.Context, d *plugin.QueryData, region string) (*databasemigrationservice.DatabaseMigrationService, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed DatabaseMigrationService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("databasemigrationservice-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*databasemigrationservice.DatabaseMigrationService), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := databasemigrationservice.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// DynamoDbService returns the service connection for AWS DynamoDb service
func DynamoDbService(ctx context.Context, d *plugin.QueryData, region string) (*dynamodb.DynamoDB, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed DynamoDbService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("dynamodb-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*dynamodb.DynamoDB), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := dynamodb.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// Ec2Service returns the service connection for AWS EC2 service
func Ec2Service(ctx context.Context, d *plugin.QueryData, region string) (*ec2.EC2, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed Ec2Service")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("ec2-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*ec2.EC2), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := ec2.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// EcrService returns the service connection for AWS ECR service
func EcrService(ctx context.Context, d *plugin.QueryData, region string) (*ecr.ECR, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed EcrService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("ecr-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*ecr.ECR), nil
	}

	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := ecr.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// EcrPublicService returns the service connection for AWS ECRPublic service
func EcrPublicService(ctx context.Context, d *plugin.QueryData, region string) (*ecrpublic.ECRPublic, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed EcrPublicService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("ecrpublic-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*ecrpublic.ECRPublic), nil
	}

	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := ecrpublic.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// EcsService returns the service connection for AWS ECS service
func EcsService(ctx context.Context, d *plugin.QueryData, region string) (*ecs.ECS, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed EcsService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("ecs-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*ecs.ECS), nil
	}

	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// EfsService returns the service connection for AWS Elastic File System service
func EfsService(ctx context.Context, d *plugin.QueryData, region string) (*efs.EFS, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed EfsService")
	}

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("efs-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*efs.EFS), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := efs.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// EksService returns the service connection for AWS EKS service
func EksService(ctx context.Context, d *plugin.QueryData, region string) (*eks.EKS, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed EksService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("eks-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*eks.EKS), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := eks.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// ElasticBeanstalkService returns the service connection for AWS ElasticBeanstalk service
func ElasticBeanstalkService(ctx context.Context, d *plugin.QueryData, region string) (*elasticbeanstalk.ElasticBeanstalk, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed ElasticBeanstalkService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("elasticbeanstalk-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*elasticbeanstalk.ElasticBeanstalk), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := elasticbeanstalk.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// ElastiCacheService returns the service connection for AWS ElastiCache service
func ElastiCacheService(ctx context.Context, d *plugin.QueryData, region string) (*elasticache.ElastiCache, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed ElastiCache")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("elasticache-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*elasticache.ElastiCache), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := elasticache.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ElasticsearchService returns the service connection for AWS Elasticsearch service
func ElasticsearchService(ctx context.Context, d *plugin.QueryData, region string) (*elasticsearchservice.ElasticsearchService, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed ElasticsearchService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("elasticache-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*elasticsearchservice.ElasticsearchService), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := elasticsearchservice.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ELBv2Service returns the service connection for AWS EC2 service
func ELBv2Service(ctx context.Context, d *plugin.QueryData, region string) (*elbv2.ELBV2, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed ELBv2Service")
	}

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("elbv2-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*elbv2.ELBV2), nil
	}

	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := elbv2.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// ELBService returns the service connection for AWS ELB Classic service
func ELBService(ctx context.Context, d *plugin.QueryData, region string) (*elb.ELB, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed ELBv2Service")
	}

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("elb-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*elb.ELB), nil
	}

	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := elb.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// EventBridgeService returns the service connection for AWS EventBridge service
func EventBridgeService(ctx context.Context, d *plugin.QueryData, region string) (*eventbridge.EventBridge, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed EventBridgeService")
	}

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("eventbridge-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*eventbridge.EventBridge), nil
	}

	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := eventbridge.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// EmrService returns the service connection for AWS EMR service
func EmrService(ctx context.Context, d *plugin.QueryData, region string) (*emr.EMR, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed EmrService")
	}

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("emr-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*emr.EMR), nil
	}

	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := emr.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// FirehoseService returns the service connection for AWS Kinesis Firehose service
func FirehoseService(ctx context.Context, d *plugin.QueryData, region string) (*firehose.Firehose, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed FirehoseService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("firehose-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*firehose.Firehose), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := firehose.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// GlacierService returns the service connection for AWS Glacier service
func GlacierService(ctx context.Context, d *plugin.QueryData, region string) (*glacier.Glacier, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed GlacierService")
	}

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("glacier-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*glacier.Glacier), nil
	}

	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := glacier.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// GlueService returns the service connection for AWS Glue service
func GlueService(ctx context.Context, d *plugin.QueryData, region string) (*glue.Glue, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed GlueService")
	}

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("glue-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*glue.Glue), nil
	}

	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := glue.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// GuardDutyService returns the service connection for AWS GuardDuty service
func GuardDutyService(ctx context.Context, d *plugin.QueryData, region string) (*guardduty.GuardDuty, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed GuardDutyService")
	}

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("guardduty-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*guardduty.GuardDuty), nil
	}

	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := guardduty.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// IAMService returns the service connection for AWS IAM service
func IAMService(ctx context.Context, d *plugin.QueryData) (*iam.IAM, error) {
	// have we already created and cached the service?
	serviceCacheKey := "iam"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*iam.IAM), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	// svc := iam.New(session.New(&aws.Config{MaxRetries: aws.Int(10)}))
	svc := iam.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// InspectorService returns the service connection for AWS Inspector service
func InspectorService(ctx context.Context, d *plugin.QueryData, region string) (*inspector.Inspector, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed InspectorService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("inspector-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*inspector.Inspector), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := inspector.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// KinesisService returns the service connection for AWS Kinesis service
func KinesisService(ctx context.Context, d *plugin.QueryData, region string) (*kinesis.Kinesis, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed KinesisService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("kinesis-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*kinesis.Kinesis), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := kinesis.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// KinesisAnalyticsV2Service returns the service connection for AWS Kinesis AnalyticsV2 service
func KinesisAnalyticsV2Service(ctx context.Context, d *plugin.QueryData, region string) (*kinesisanalyticsv2.KinesisAnalyticsV2, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed KinesisAnalyticsV2Service")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("kinesisanalyticsv2-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*kinesisanalyticsv2.KinesisAnalyticsV2), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := kinesisanalyticsv2.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// KinesisVideoService returns the service connection for AWS Kinesis Video service
func KinesisVideoService(ctx context.Context, d *plugin.QueryData, region string) (*kinesisvideo.KinesisVideo, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed Kinesis Video")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("kinesisvideo-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*kinesisvideo.KinesisVideo), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := kinesisvideo.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// KMSService returns the service connection for AWS KMS service
func KMSService(ctx context.Context, d *plugin.QueryData, region string) (*kms.KMS, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed KMSService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("kms-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*kms.KMS), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := kms.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// LambdaService returns the service connection for AWS Lambda service
func LambdaService(ctx context.Context, d *plugin.QueryData, region string) (*lambda.Lambda, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed LambdaService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("lambda-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*lambda.Lambda), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := lambda.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// OrganizationService returns the service connection for AWS Organization service
func OrganizationService(ctx context.Context, d *plugin.QueryData) (*organizations.Organizations, error) {
	// have we already created and cached the service?
	serviceCacheKey := "Organization"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*organizations.Organizations), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	svc := organizations.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// ConfigService returns the service connection for AWS Config  service
func ConfigService(ctx context.Context, d *plugin.QueryData, region string) (*configservice.ConfigService, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed ConfigService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("configservice-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*configservice.ConfigService), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := configservice.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// RDSService returns the service connection for AWS RDS service
func RDSService(ctx context.Context, d *plugin.QueryData, region string) (*rds.RDS, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed RDSService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("rds-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*rds.RDS), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := rds.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// RedshiftService returns the service connection for AWS Redshift service
func RedshiftService(ctx context.Context, d *plugin.QueryData, region string) (*redshift.Redshift, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed Redshift")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("redshift-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*redshift.Redshift), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := redshift.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// Route53ResolverService returns the service connection for AWS route53resolver service
func Route53ResolverService(ctx context.Context, d *plugin.QueryData, region string) (*route53resolver.Route53Resolver, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed Route53Resolver")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("route53resolver-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*route53resolver.Route53Resolver), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := route53resolver.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// Route53Service returns the service connection for AWS route53 service
func Route53Service(ctx context.Context, d *plugin.QueryData) (*route53.Route53, error) {
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("route53")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*route53.Route53), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	svc := route53.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// SecretsManagerService returns the service connection for AWS secretsManager service
func SecretsManagerService(ctx context.Context, d *plugin.QueryData, region string) (*secretsmanager.SecretsManager, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed SecretsManagerService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("secretsManager-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*secretsmanager.SecretsManager), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := secretsmanager.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// SecurityHubService returns the service connection for AWS securityHub service
func SecurityHubService(ctx context.Context, d *plugin.QueryData, region string) (*securityhub.SecurityHub, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed SecurityHubService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("securityHub-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*securityhub.SecurityHub), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := securityhub.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// S3ControlService returns the service connection for AWS s3control service
func S3ControlService(ctx context.Context, d *plugin.QueryData, region string) (*s3control.S3Control, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed S3ControlService")
	}

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("s3control-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*s3control.S3Control), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := s3control.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// S3Service returns the service connection for AWS S3 service
func S3Service(ctx context.Context, d *plugin.QueryData, region string) (*s3.S3, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed S3Service")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("s3-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*s3.S3), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// SageMakerService returns the service connection for AWS SageMaker service
func SageMakerService(ctx context.Context, d *plugin.QueryData, region string) (*sagemaker.SageMaker, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed SageMakerService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("sagemaker-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*sagemaker.SageMaker), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := sagemaker.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// SNSService returns the service connection for AWS SNS service
func SNSService(ctx context.Context, d *plugin.QueryData, region string) (*sns.SNS, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed SNSService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("sns-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*sns.SNS), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := sns.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// SQSService returns the service connection for AWS SQS service
func SQSService(ctx context.Context, d *plugin.QueryData, region string) (*sqs.SQS, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed SQSService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("sqs-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*sqs.SQS), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := sqs.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// SsmService returns the service connection for AWS SSM service
func SsmService(ctx context.Context, d *plugin.QueryData, region string) (*ssm.SSM, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed SsmService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("ssm-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*ssm.SSM), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := ssm.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// StsService returns the service connection for AWS STS service
func StsService(ctx context.Context, d *plugin.QueryData) (*sts.STS, error) {
	// have we already created and cached the service?
	serviceCacheKey := "sts"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*sts.STS), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	svc := sts.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// WAFService returns the service connection for AWS WAF service
func WAFService(ctx context.Context, d *plugin.QueryData) (*waf.WAF, error) {

	// have we already created and cached the service?
	serviceCacheKey := "waf"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*waf.WAF), nil
	}

	// so it was not in cache - create service
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
	if err != nil {
		return nil, err
	}
	svc := waf.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// WAFv2Service returns the service connection for AWS WAFv2 service
func WAFv2Service(ctx context.Context, d *plugin.QueryData, region string) (*wafv2.WAFV2, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed WAFv2")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("wafv2-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*wafv2.WAFV2), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := wafv2.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

// WellArchitectedService returns the service connection for AWS Well-Architected service
func WellArchitectedService(ctx context.Context, d *plugin.QueryData, region string) (*wellarchitected.WellArchitected, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed WellArchitectedService")
	}
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("wellarchitected-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*wellarchitected.WellArchitected), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, region)
	if err != nil {
		return nil, err
	}
	svc := wellarchitected.New(sess)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)

	return svc, nil
}

func getSession(ctx context.Context, d *plugin.QueryData, region string) (*session.Session, error) {
	// get aws config info
	awsConfig := GetConfig(d.Connection)
	sessionOptions := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	if &awsConfig != nil {
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
	}

	// TODO is it correct to always pass region to session?
	// have we cached a session?
	sessionCacheKey := fmt.Sprintf("session-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*session.Session), nil
	}

	// so it was not in cache - create a session
	// sess, err := session.NewSession(&aws.Config{Region: &region, MaxRetries: aws.Int(10)})

	sessionOptions.Config.Region = &region
	sessionOptions.Config.MaxRetries = aws.Int(10)

	// so it was not in cache - create a session
	sess, err := session.NewSessionWithOptions(sessionOptions)
	if err != nil {
		return nil, err
	}
	// save session in cache
	d.ConnectionManager.Cache.Set(sessionCacheKey, sess)

	return sess, nil
}

// GetDefaultRegion returns the default region used
func GetDefaultRegion() string {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	session, err := session.NewSession(aws.NewConfig())
	if err != nil {
		panic(err)
	}

	region := *session.Config.Region
	if region == "" {
		// get aws config info
		panic("\n\n'regions' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	return region
}

// GetDefaultAwsRegion returns the default region for AWS partiton
// if not set by Env variable or in aws profile or i
func GetDefaultAwsRegion(d *plugin.QueryData) string {
	// have we already created and cached the service?
	serviceCacheKey := "GetDefaultAwsRegion"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(string)
	}

	// get aws config info
	awsConfig := GetConfig(d.Connection)

	var regions []string
	var region string

	if &awsConfig != nil && awsConfig.Regions != nil {
		regions = GetConfig(d.Connection).Regions
	}

	if len(getInvalidRegions(regions)) < 1 {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
		session, err := session.NewSession(aws.NewConfig())
		if err != nil {
			panic(err)
		}
		region = *session.Config.Region
	} else {
		// Set the first region in regions list to be default region
		region = regions[0]

		// check if it is a valid region
		if len(getInvalidRegions([]string{region})) > 0 {
			panic("\n\nConnection config have invalid region: " + region + ". Edit your connection configuration file and then restart Steampipe")
		}
	}

	if region == "" {
		region = "us-east-1"
	}
	return region
}
