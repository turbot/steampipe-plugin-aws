package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

const (
	APIGatewayServiceID                      = "apigateway"              // API Gateway
	APIGatewayV2ServiceID                    = "apigateway"              // ApiGatewayV2
	AccessAnalyzerServiceID                  = "access-analyzer"         // AccessAnalyzer
	ACMServiceID                             = "acm"                     // ACM
	AmplifyServiceID                         = "amplify"                 // Amplify
	ApplicationAutoScalingServiceID          = "application-autoscaling" // Application Auto Scaling
	AuditManagerServiceID                    = "auditmanager"            // AuditManager
	AutoScalingServiceID                     = "autoscaling"             // AutoScaling
	BackupServiceID                          = "backup"                  // Backup
	CloudControlServiceID                    = "cloudcontrolapi"         // CloudControl
	CloudFormationServiceID                  = "cloudformation"          // CloudFormation
	CloudFrontServiceID                      = "cloudfront"              // CloudFront
	CloudTrailServiceID                      = "cloudtrail"              // CloudTrail
	CloudWatchLogsServiceID                  = "logs"                    // CloudWatch Logs
	CodeBuildServiceID                       = "codebuild"               // CodeBuild
	CodeCommitServiceID                      = "codecommit"              // CodeCommit
	CodePipelineServiceID                    = "codepipeline"            // CodePipeline
	ConfigServiceID                          = "config"                  // Config Service
	CostExplorerServiceID                    = "ce"                      // Cost Explorer
	DAXServiceID                             = "dax"                     // DAX
	DLMServiceID                             = "dlm"                     // DLM
	DatabaseMigrationServiceID               = "dms"                     // Database Migration Service
	DirectoryServiceID                       = "ds"                      // Directory Service
	DynamoDBServiceID                        = "dynamodb"                // DynamoDB
	EC2ServiceID                             = "ec2"                     // EC2
	ECRPublicServiceID                       = "api.ecr-public"          // ECR PUblic Service
	ECRServiceID                             = "api.ecr"                 // ECR Service
	ECSServiceID                             = "ecs"                     // ECS Service
	EFSServiceID                             = "elasticfilesystem"       // EFS Service
	EKSServiceID                             = "eks"                     // EKS Service
	EMRServiceID                             = "elasticmapreduce"        // EMR Service
	ElastiCacheServiceID                     = "elasticache"             // ElastiCache Service
	ElasticBeanstalkServiceID                = "elasticbeanstalk"        // Elastic Beanstalk Service
	ElasticLoadBalancingServiceID            = "elasticloadbalancing"    // Elastic Load Balancing Service
	ElasticLoadBalancingV2ServiceID          = "elasticloadbalancing"    // Elastic Load Balancing v2 Service
	ElasticsearchServiceID                   = "es"                      // Elasticsearch Service
	EventBridgeServiceID                     = "events"                  // EventBridge Service
	FSxServiceID                             = "fsx"                     // FSx Service
	FirehoseServiceID                        = "firehose"                // Firehose Service
	GlacierServiceID                         = "glacier"                 // Glacier Service
	GlobalAcceleratorServiceID               = "globalaccelerator"       // Global Accelerator Service
	GlueServiceID                            = "glue"                    // Glue Service
	GuardDutyServiceID                       = "guardduty"               // GuardDuty Service
	IAMServiceID                             = "iam"                     // IAM Service
	IdentityStoreServiceID                   = "identitystore"           // Identity Store Service
	InspectorServiceID                       = "inspector"               // Inspector Service
	KMSServiceID                             = "kms"                     // KMS Service
	KinesisAnalyticsV2ServiceID              = "kinesisanalytics"        // Kinesis Analytics V2 Service
	KinesisServiceID                         = "kinesis"                 // Kinesis Service
	KinesisVideoServiceID                    = "kinesisvideo"            // Kinesis Video Service
	LambdaServiceID                          = "lambda"                  // Lambda Service
	Macie2ServiceID                          = "macie2"                  // Macie2 Service
	Neptune                                  = "rds"                     // Neptune Service
	NetworkFirewallServiceID                 = "network-firewall"        // Network Firewall Service
	OpenSearchServiceID                      = "es"                      // OpenSearch Service
	OrganizationsServiceID                   = "organizations"           // Organizations Service
	PinpointServiceID                        = "pinpoint"                // Pinpoint Service
	PricingServiceID                         = "api.pricing"             // Pricing Service
	RAMServiceID                             = "ram"                     // RAM Service
	RDSServiceID                             = "rds"                     // RDS Service
	RedshiftServiceID                        = "redshift"                // Redshift Service
	ResourceGroupsTaggingAPIServiceID        = "tagging"                 // Resource Groups Tagging API
	Route53DomainsServiceID                  = "route53domains"          // Route 53 Domains Service
	Route53ResolverServiceID                 = "route53resolver"         // Route 53 Resolver Service
	Route53ServiceID                         = "route53"                 // Route 53 Service
	S3ControlServiceID                       = "s3-control"              // S3 Control Service
	SESServiceID                             = "email"                   // SES Service
	SFNServiceID                             = "states"                  // SFN Service
	SNSServiceID                             = "sns"                     // SNS Service
	SQSServiceID                             = "sqs"                     // SQS Service
	SSMServiceID                             = "ssm"                     // SSM Service
	SSOAdminServiceID                        = "sso"                     // SSO Admin Service
	STSServiceID                             = "sts"                     // STS Service
	SageMakerServiceID                       = "api.sagemaker"           // SageMaker Service
	SecretsManagerServiceID                  = "secretsmanager"          // Secrets Manager Service
	SecurityHubServiceID                     = "securityhub"             // Security Hub Service
	ServerlessApplicationRepositoryServiceID = "serverlessrepo"          // Serverless Application Repository Service
	ServiceQuotasServiceID                   = "servicequotas"           // Service Quotas Service
	WAFRegionalServiceID                     = "waf-regional"            // WAF Regional Service
	WAFServiceID                             = "waf"                     // WAF Service
	WAFV2ServiceID                           = "wafv2"                   // WAFV2 Service
	WellArchitectedServiceID                 = "wellarchitected"         // WellArchitected Service
	WorkSpacesServiceID                      = "workspaces"              // WorkSpaces Service
)

// endpoints holds the aws generated endpoints.json
type Endpoints struct {
	Partitions []Partition `json:"partitions"`
}

type Partition struct {
	PartitionName string                 `json:"partitionName"`
	PartitionCode string                 `json:"partition"`
	Regions       map[string]interface{} `json:"regions"`
	Services      map[string]Service     `json:"services"`
}

type Service struct {
	Endpoints map[string]interface{} `json:"endpoints"`
}

func GetAwsEndpoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheKey := "aws-endpoints"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*Endpoints), nil
	}
	fmt.Println("Generating AWS regions")

	resp, err := http.Get("https://raw.githubusercontent.com/aws/aws-sdk-go-v2/master/codegen/smithy-aws-go-codegen/src/main/resources/software/amazon/smithy/aws/go/codegen/endpoints.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	e := Endpoints{}
	err = json.Unmarshal(body, &e)
	if err != nil {
		return nil, err
	}

	return &e, err
}

func (e Endpoints) GetPartitionByName(partitionName string) Partition {
	var defaultPartition Partition
	for _, partition := range e.Partitions {
		if partition.PartitionCode == partitionName {
			return partition
		}
		if partition.PartitionCode == "aws" {
			defaultPartition = partition
		}
	}
	return defaultPartition
}

// Services returns a map of Service indexed by their ID. This is useful for
// enumerating over the services in a partition.
func (p Partition) Service(name string) *Service {
	service := p.Services[name]
	return &service
}

// Services returns a list of regions supported by the service.
func (s Service) Regions() []string {
	regions := make([]string, 0, len(s.Endpoints))

	for k := range s.Endpoints {
		if k != "" {
			regions = append(regions, k)
		}
	}

	return regions
}
