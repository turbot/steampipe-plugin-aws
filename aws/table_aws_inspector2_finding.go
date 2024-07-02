package aws

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"

	inspector2v1 "github.com/aws/aws-sdk-go/service/inspector2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

//// TABLE DEFINITION

func tableAwsInspector2Finding(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector2_finding",
		Description: "AWS Inspector2 Finding",
		List: &plugin.ListConfig{
			Hydrate: listInspector2Finding,
			Tags:    map[string]string{"service": "inspector2", "action": "ListFindings"},
			KeyColumns: plugin.KeyColumnSlice{
				// The AWS CLI supports EQUALS, PREFIX, and NOT_EQUALS... we can't represent PREFIX.
				{Name: "finding_account_id", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "exploit_available", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "arn", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "first_observed_at", Operators: []string{"<=", ">="}, Require: plugin.Optional},
				{Name: "fix_available", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "inspector_score", Operators: []string{"<=", ">="}, Require: plugin.Optional},
				{Name: "last_observed_at", Operators: []string{"<=", ">="}, Require: plugin.Optional},
				{Name: "severity", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "component_id", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "component_type", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "ec2_instance_image_id", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "ec2_instance_subnet_id", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "ec2_instance_vpc_id", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "ecr_image_architecture", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "ecr_image_hash", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "ecr_image_pushed_at", Operators: []string{"<=", ">="}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "lambda_function_last_modified_at", Operators: []string{"<=", ">="}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "ecr_image_registry", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "ecr_image_repository_name", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "lambda_function_execution_role_arn", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "lambda_function_layers", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "lambda_function_name", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "lambda_function_runtime", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "network_protocol", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "related_vulnerabilitie", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "resource_id", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "resource_type", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "ecr_image_tags", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "source", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "status", Operators: []string{"=", "<>"}, Require: plugin.Optional}, // findingStatus
				{Name: "title", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "type", Operators: []string{"=", "<>"}, Require: plugin.Optional}, // findingType
				{Name: "updated_at", Operators: []string{"<=", ">="}, Require: plugin.Optional},
				{Name: "vendor_severity", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "vulnerability_id", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "resource_tags", Operators: []string{"="}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "vulnerable_package", Operators: []string{"="}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(inspector2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				// Technically, this would be "aws_account_id", but "aws" is
				// redundant, and doesn't add any semantic meaning. As the
				// description says, this is the account ID of the finding,
				// hence "finding_account_id".
				Name:        "finding_account_id",
				Type:        proto.ColumnType_STRING,
				Description: "The Amazon Web Services account ID associated with the finding.",
				Transform:   transform.FromField("AwsAccountId"),
			},
			{
				Name:        "description",
				Description: "The description of the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "exploit_available",
				Description: "If a finding discovered in your environment has an exploit available. Valid values are: YES | NO.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Number (ARN) of the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FindingArn"),
			},
			{
				Name:        "status",
				Description: "The status of the finding. Valid values are: ACTIVE | SUPPRESSED | CLOSED.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the finding. Valid values are: NETWORK_REACHABILITY | PACKAGE_VULNERABILITY.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "first_observed_at",
				Description: "The date and time that the finding was first observed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "fix_available",
				Description: "Details on whether a fix is available through a version update. Valid values are: YES | NO | PARTIAL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "inspector_score",
				Description: "The Amazon Inspector score given to the finding.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "resource_id",
				Description: "The ID of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("resource_id"),
			},
			{
				Name:        "resource_type",
				Description: "The resource type supported by AWS.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("resource_type"),
			},
			{
				Name:        "component_type",
				Description: "The component type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("component_type"),
			},
			{
				Name:        "component_id",
				Description: "The component ID of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("component_id"),
			},
			{
				Name:        "ec2_instance_image_id",
				Description: "The Amazon EC2 instance image ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("ec2_instance_image_id"),
			},
			{
				Name:        "ec2_instance_subnet_id",
				Description: "The Amazon EC2 instance subnet ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("ec2_instance_subnet_id"),
			},
			{
				Name:        "ec2_instance_vpc_id",
				Description: "The Amazon EC2 instance VPC ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("ec2_instance_vpc_id"),
			},
			{
				Name:        "ecr_image_architecture",
				Description: "The Amazon ECR image architecture.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("ecr_image_architecture"),
			},
			{
				Name:        "ecr_image_hash",
				Description: "The Amazon ECR image hash.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("ecr_image_hash"),
			},
			{
				Name:        "ecr_image_pushed_at",
				Description: "The Amazon ECR image push date and time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromQual("ecr_image_pushed_at"),
			},
			{
				Name:        "ecr_image_registry",
				Description: "The Amazon ECR registry.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("ecr_image_registry"),
			},
			{
				Name:        "ecr_image_repository_name",
				Description: "The name of the Amazon ECR repository.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("ecr_image_repository_name"),
			},
			{
				Name:        "ecr_image_tags",
				Description: "The tags attached to the Amazon ECR container image.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("ecr_image_tags"),
			},
			{
				Name:        "lambda_function_execution_role_arn",
				Description: "The AWS Lambda function execution role ARN.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("lambda_function_execution_role_arn"),
			},
			{
				Name:        "lambda_function_last_modified_at",
				Description: "The AWS Lambda functions the date and time that a user last updated the configuration.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromQual("lambda_function_last_modified_at"),
			},
			{
				Name:        "lambda_function_layers",
				Description: "The AWS Lambda function layer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("lambda_function_layers"),
			},
			{
				Name:        "lambda_function_name",
				Description: "The AWS Lambda function name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("lambda_function_name"),
			},
			{
				Name:        "lambda_function_runtime",
				Description: "The AWS Lambda function runtime environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("lambda_function_runtime"),
			},
			{
				Name:        "network_protocol",
				Description: "The ingress source addresse.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("network_protocol"),
			},
			{
				Name:        "related_vulnerabilitie",
				Description: "The related vulnerabilitie.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("related_vulnerabilitie"),
			},
			{
				Name:        "last_observed_at",
				Description: "The date and time that the finding was last observed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "remediation_recommendation_text",
				Description: "The recommended course of action to remediate the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Remediation.Recommendation.Text"),
			},
			{
				Name:        "remediation_recommendation_url",
				Description: "The URL address to the CVE remediation recommendations.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Remediation.Recommendation.Url"),
			},
			{
				Name:        "severity",
				Description: "The severity of the finding. Valid values are: INFORMATIONAL | LOW | MEDIUM | HIGH | CRITICAL | UNTRIAGED.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The date and time the finding was last updated at.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "source",
				Description: "The source of the vulnerability information.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PackageVulnerabilityDetails.Source"),
			},
			{
				Name:        "source_url",
				Description: "A URL to the source of the vulnerability information.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PackageVulnerabilityDetails.SourceUrl"),
			},
			{
				Name:        "vendor_created_at",
				Description: "The date and time that this vulnerability was first added to the vendor’s database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("PackageVulnerabilityDetails.VendorCreatedAt"),
			},
			{
				Name:        "vendor_severity",
				Description: "The severity the vendor has given to this vulnerability type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PackageVulnerabilityDetails.VendorSeverity"),
			},
			{
				Name:        "vendor_updated_at",
				Description: "The date and time the vendor last updated this vulnerability in their database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("PackageVulnerabilityDetails.vendorUpdatedAt"),
			},
			{
				Name:        "vulnerability_id",
				Description: "The ID given to this vulnerability.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PackageVulnerabilityDetails.VulnerabilityId"),
			},
			{
				Name:        "exploitability_details",
				Description: "The details of an exploit available for a finding discovered in your environment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "inspector_score_details",
				Description: "An object that contains details of the Amazon Inspector score.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_reachability_details",
				Description: "An object that contains the details of a network reachability finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "package_vulnerability_details",
				Description: "An object that contains the details of a package vulnerability finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cvss",
				Description: "An object that contains details about the CVSS score of a finding.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PackageVulnerabilityDetails.Cvss"),
			},
			{
				Name:        "reference_urls",
				Description: "One or more URLs that contain details about this vulnerability type.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PackageVulnerabilityDetails.ReferenceUrls"),
			},
			{
				Name:        "related_vulnerabilities",
				Description: "One or more vulnerabilities related to the one identified in this finding.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PackageVulnerabilityDetails.RelatedVulnerabilities"),
			},
			{
				Name:        "vulnerable_package",
				Description: "The package impacted by this vulnerability.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromQual("vulnerable_package"),
			},
			{
				Name:        "vulnerable_packages",
				Description: "The packages impacted by this vulnerability.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PackageVulnerabilityDetails.VulnerablePackages"),
			},
			{
				Name:        "resources",
				Description: "Contains information on the resources involved in a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_tags",
				Description: "Details on the resource tags used to filter findings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromQual("resource_tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "The title of the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FindingArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listInspector2Finding(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := Inspector2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector2_finding.listInspector2Finding", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, etc., return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &inspector2.ListFindingsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filter := &types.FilterCriteria{}

	buildStringFilters(d, filter)
	buildNumberFilters(d, filter)
	buildDateFilters(d, filter)
	buildResourceTagFilter(d, filter)
	buildVulnerablePackageFilter(d, filter)

	input.FilterCriteria = filter

	paginator := inspector2.NewListFindingsPaginator(svc, input, func(o *inspector2.ListFindingsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_inspector2_finding.listInspector2Finding", "api_error", err)
			return nil, err
		}

		for _, item := range output.Findings {
			// item := item
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

// column-to-filter mapping
type stringFilterField struct {
	columnName  string
	filterField func(f *types.FilterCriteria) *[]types.StringFilter
}

type resourceTagFilter struct {
	columnName  string
	filterField func(f *types.FilterCriteria) *[]types.MapFilter
}

type vulnerablePackageFilter struct {
	columnName  string
	filterField func(f *types.FilterCriteria) *[]types.PackageFilter
}

type numberFilterField struct {
	columnName  string
	filterField func(f *types.FilterCriteria) *[]types.NumberFilter
}

type dateFilterField struct {
	columnName  string
	filterField func(f *types.FilterCriteria) *[]types.DateFilter
}

var findingStringFilters = []stringFilterField{
	{
		columnName: "finding_account_id",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.AwsAccountId)
		},
	},
	{
		columnName: "exploit_available",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.ExploitAvailable)
		},
	},
	{
		columnName: "arn",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.FindingArn)
		},
	},
	{
		columnName: "fix_available",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.FixAvailable)
		},
	},
	{
		columnName: "severity",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.Severity)
		},
	},
	{
		columnName: "source",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.VulnerabilitySource)
		},
	},
	{
		columnName: "resource_id",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.ResourceId)
		},
	},
	{
		columnName: "resource_type",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.ResourceType)
		},
	},
	{
		columnName: "component_id",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.ComponentId)
		},
	},
	{
		columnName: "component_type",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.ComponentType)
		},
	},
	{
		columnName: "ec2_instance_image_id",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.Ec2InstanceImageId)
		},
	},
	{
		columnName: "ec2_instance_subnet_id",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.Ec2InstanceSubnetId)
		},
	},
	{
		columnName: "ec2_instance_vpc_id",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.Ec2InstanceVpcId)
		},
	},
	{
		columnName: "ecr_image_architecture",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.EcrImageArchitecture)
		},
	},
	{
		columnName: "ecr_image_hash",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.EcrImageHash)
		},
	},
	{
		columnName: "ecr_image_registry",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.EcrImageRegistry)
		},
	},
	{
		columnName: "ecr_image_repository_name",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.EcrImageRepositoryName)
		},
	},
	{
		columnName: "ecr_image_tags",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.EcrImageTags)
		},
	},
	{
		columnName: "lambda_function_execution_role_arn",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.LambdaFunctionExecutionRoleArn)
		},
	},
	{
		columnName: "lambda_function_layers",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.LambdaFunctionLayers)
		},
	},
	{
		columnName: "lambda_function_name",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.LambdaFunctionName)
		},
	},
	{
		columnName: "lambda_function_name",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.LambdaFunctionName)
		},
	},
	{
		columnName: "lambda_function_runtime",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.LambdaFunctionRuntime)
		},
	},
	{
		columnName: "network_protocol",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.NetworkProtocol)
		},
	},
	{
		columnName: "related_vulnerabilitie",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.RelatedVulnerabilities)
		},
	},
	{
		columnName: "status",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.FindingStatus)
		},
	},
	{
		columnName: "title",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.Title)
		},
	},
	{
		columnName: "type",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.FindingType)
		},
	},
	{
		columnName: "vendor_severity",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.VendorSeverity)
		},
	},
	{
		columnName: "vulnerability_id",
		filterField: func(f *types.FilterCriteria) *[]types.StringFilter {
			return &(f.VulnerabilityId)
		},
	},
}

var findingNumberFilters = []numberFilterField{
	{
		columnName: "inspector_score",
		filterField: func(f *types.FilterCriteria) *[]types.NumberFilter {
			return &(f.InspectorScore)
		},
	},
}

var findingVulnerablePackageFilter = []vulnerablePackageFilter{
	{
		columnName: "vulnerable_package",
		filterField: func(f *types.FilterCriteria) *[]types.PackageFilter {
			return &(f.VulnerablePackages)
		},
	},
}

var findingResourceTagFilters = []resourceTagFilter{
	{
		columnName: "resource_tags",
		filterField: func(f *types.FilterCriteria) *[]types.MapFilter {
			return &(f.ResourceTags)
		},
	},
}

var findingDateFilters = []dateFilterField{
	{
		columnName: "first_observed_at",
		filterField: func(f *types.FilterCriteria) *[]types.DateFilter {
			return &(f.FirstObservedAt)
		},
	},
	{
		columnName: "last_observed_at",
		filterField: func(f *types.FilterCriteria) *[]types.DateFilter {
			return &(f.LastObservedAt)
		},
	},
	{
		columnName: "updated_at",
		filterField: func(f *types.FilterCriteria) *[]types.DateFilter {
			return &(f.UpdatedAt)
		},
	},
	{
		columnName: "ecr_image_pushed_at",
		filterField: func(f *types.FilterCriteria) *[]types.DateFilter {
			return &(f.EcrImagePushedAt)
		},
	},
	{
		columnName: "lambda_function_last_modified_at",
		filterField: func(f *types.FilterCriteria) *[]types.DateFilter {
			return &(f.LambdaFunctionLastModifiedAt)
		},
	},
}

// there are a bunch of keys with identically-structured syntax/types... we
// can't easily include the filter field without reflection, but we *can* at
// least construct the filter itself.
func buildStringFilters(d *plugin.QueryData, filter *types.FilterCriteria) {
	for _, info := range findingStringFilters {
		if d.Quals[info.columnName] != nil {
			field := info.filterField(filter)
			for _, q := range d.Quals[info.columnName].Quals {
				val := aws.String(q.Value.GetStringValue())
				var comp types.StringComparison
				switch q.Operator {
				case "=":
					comp = types.StringComparisonEquals
				case "<>":
					comp = types.StringComparisonNotEquals
				}
				*field = append(*field, types.StringFilter{
					Comparison: comp,
					Value:      val,
				})
			}
		}
	}
}

func buildResourceTagFilter(d *plugin.QueryData, filter *types.FilterCriteria) {
	for _, info := range findingResourceTagFilters {
		if d.Quals[info.columnName] != nil {
			field := info.filterField(filter)
			tagValue := make([]map[string]string, 0)
			for _, q := range d.Quals[info.columnName].Quals {
				val := q.Value.GetJsonbValue()
				if val != "" && q.Operator == "=" {
					_ = json.Unmarshal([]byte(val), &tagValue)
					for _, v := range tagValue {
						if v["key"] == "" || v["value"] == "" {
							tagfilter := types.MapFilter{
								Comparison: types.MapComparisonEquals,
								Key:        aws.String(v["key"]),
								Value:      aws.String(v["value"]),
							}
							*field = append(*field, tagfilter)
						}
					}
				}
			}
		}
	}
}

func buildVulnerablePackageFilter(d *plugin.QueryData, filter *types.FilterCriteria) {
	for _, info := range findingVulnerablePackageFilter {
		if d.Quals[info.columnName] != nil {
			field := info.filterField(filter)
			packageAttributes := make([]map[string]string, 0)
			for _, q := range d.Quals[info.columnName].Quals {
				val := q.Value.GetJsonbValue()
				if val != "" && q.Operator == "=" {
					_ = json.Unmarshal([]byte(val), &packageAttributes)
					for _, v := range packageAttributes {
						packageFilter := types.PackageFilter{}
						if v["architecture"] != "" {
							packageFilter.Architecture = &types.StringFilter{
								Comparison: types.StringComparisonEquals,
								Value:      aws.String(v["architecture"]),
							}
						}
						if v["epoch"] != "" {
							value, _ := strconv.ParseFloat(v["epoch"], 64)
							packageFilter.Epoch = &types.NumberFilter{
								LowerInclusive: aws.Float64(value),
							}
						}
						if v["name"] != "" {
							packageFilter.Name = &types.StringFilter{
								Comparison: types.StringComparisonEquals,
								Value:      aws.String(v["name"]),
							}
						}
						if v["release"] != "" {
							packageFilter.Release = &types.StringFilter{
								Comparison: types.StringComparisonEquals,
								Value:      aws.String(v["release"]),
							}
						}
						if v["sourceLambdaLayerArn"] != "" {
							packageFilter.SourceLambdaLayerArn = &types.StringFilter{
								Comparison: types.StringComparisonEquals,
								Value:      aws.String(v["sourceLambdaLayerArn"]),
							}
						}
						if v["sourceLayerHash"] != "" {
							packageFilter.SourceLayerHash = &types.StringFilter{
								Comparison: types.StringComparisonEquals,
								Value:      aws.String(v["sourceLayerHash"]),
							}
						}
						if v["version"] != "" {
							packageFilter.Version = &types.StringFilter{
								Comparison: types.StringComparisonEquals,
								Value:      aws.String(v["version"]),
							}
						}
						*field = append(*field, packageFilter)
					}
				}
			}
		}
	}
}

func buildNumberFilters(d *plugin.QueryData, filter *types.FilterCriteria) {
	for _, info := range findingNumberFilters {
		if d.Quals[info.columnName] != nil {
			field := info.filterField(filter)
			for _, q := range d.Quals[info.columnName].Quals {
				val := aws.Float64(q.Value.GetDoubleValue())
				var f types.NumberFilter
				switch q.Operator {
				case "<=":
					f = types.NumberFilter{LowerInclusive: val}
				case ">=":
					f = types.NumberFilter{UpperInclusive: val}
				}
				*field = append(*field, f)
			}
		}
	}
}

func buildDateFilters(d *plugin.QueryData, filter *types.FilterCriteria) {
	for _, info := range findingDateFilters {
		if d.Quals[info.columnName] != nil {
			field := info.filterField(filter)
			for _, q := range d.Quals[info.columnName].Quals {
				val := aws.Time(q.Value.GetTimestampValue().AsTime())
				var f types.DateFilter
				switch q.Operator {
				case "<=":
					f = types.DateFilter{StartInclusive: val}
				case ">=":
					f = types.DateFilter{EndInclusive: val}
				}
				*field = append(*field, f)
			}
		}
	}
}
