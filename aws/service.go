package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

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

// ElasticacheService returns the service connection for AWS Elasticache service
func ElasticacheService(ctx context.Context, d *plugin.QueryData, region string) (*elasticache.ElastiCache, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed Elasticache")
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

// S3ControlService returns the service connection for AWS s3control service
func S3ControlService(ctx context.Context, d *plugin.QueryData) (*s3control.S3Control, error) {
	// have we already created and cached the service?
	serviceCacheKey := "s3control"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*s3control.S3Control), nil
	}
	// so it was not in cache - create service
	sess, err := getSession(ctx, d, GetDefaultAwsRegion(d))
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
