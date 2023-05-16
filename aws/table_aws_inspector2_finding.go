package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"

	inspector2v1 "github.com/aws/aws-sdk-go/service/inspector2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInspector2Finding(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector2_finding",
		Description: "AWS Inspector2 Finding",
		List: &plugin.ListConfig{
			Hydrate: listInspector2Finding,
			KeyColumns: plugin.KeyColumnSlice{
				// The AWS CLI supports EQUALS, PREFIX, and NOT_EQUALS... we can't represent PREFIX.
				// Also, the following filter criteria aren't (easily) represented in the output columns,
				// and thus can't be input filters: componentId, componentType, ec2InstanceImageId,
				// ec2InstanceSubnetId, ec2InstanceVpcId, ecrImageArchitecture, ecrImageHash,
				// ecrImagePushedAt, ecrImageRegistry, ecrImageRepositoryName, ecrImageTags,
				// lambdaFunctionExecutionRoleArn, lambdaFunctionLastModifiedAt, lambdaFunctionLayers,
				// lambdaFunctionName, lambdaFunctionRuntime, networkProtocol,  portRange,
				// relatedVulnerabilities, resourceId, resourceTags, resourceType, vulnerabilitySource,
				// vulnerablePackages (exists, but complex!)
				{Name: "finding_account_id", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "exploit_available", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "finding_arn", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "first_observed_at", Operators: []string{"<=", ">="}, Require: plugin.Optional},
				{Name: "fix_available", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "inspector_score", Operators: []string{"<=", ">="}, Require: plugin.Optional},
				{Name: "last_observed_at", Operators: []string{"<=", ">="}, Require: plugin.Optional},
				{Name: "severity", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "status", Operators: []string{"=", "<>"}, Require: plugin.Optional}, // findingStatus
				{Name: "title", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "type", Operators: []string{"=", "<>"}, Require: plugin.Optional}, // findingType
				{Name: "updated_at", Operators: []string{"<=", ">="}, Require: plugin.Optional},
				{Name: "vendor_severity", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "vulnerability_id", Operators: []string{"=", "<>"}, Require: plugin.Optional},
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
				Description: "If a finding discovered in your environment has an exploit available.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "exploitability_details",
				Description: "The details of an exploit available for a finding discovered in your environment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "finding_arn",
				Description: "The Amazon Resource Number (ARN) of the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "first_observed_at",
				Description: "The date and time that the finding was first observed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "fix_available",
				Description: "Details on whether a fix is available through a version update.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "inspector_score",
				Description: "The Amazon Inspector score given to the finding.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "inspector_score_details",
				Description: "An object that contains details of the Amazon Inspector score.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "last_observed_at",
				Description: "The date and time that the finding was last observed.",
				Type:        proto.ColumnType_TIMESTAMP,
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
				Name:        "reference_urls",
				Description: "One or more URLs that contain details about this vulnerability type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "related_vulnerabilities",
				Description: "One or more vulnerabilities related to the one identified in this finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source",
				Description: "The source of the vulnerability information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_url",
				Description: "A URL to the source of the vulnerability information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vendor_created_at",
				Description: "The date and time that this vulnerability was first added to the vendor’s database.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "vendor_severity",
				Description: "The severity the vendor has given to this vulnerability type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vendor_updated_at",
				Description: "The date and time the vendor last updated this vulnerability in their database.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "vulnerability_id",
				Description: "The ID given to this vulnerability.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vulnerable_packages",
				Description: "The packages impacted by this vulnerability.",
				Type:        proto.ColumnType_JSON,
			},
			// {
			// 	Name:        "remediation",
			// 	Description: "An object that contains the details about how to remediate a finding.",
			// 	Type:        proto.ColumnType_JSON,
			// },
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
				Name:        "resources",
				Description: "Contains information on the resources involved in a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "severity",
				Description: "The severity of the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "title",
				Description: "The title of the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The date and time the finding was last updated at.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Steampipe standard columns
			// There's already a "title" column!
			// {
			// 	Name:        "title",
			// 	Description: resourceInterfaceDescription("title"),
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromField("FindingArn").Transform(arnToTitle),
			// },
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FindingArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

// column-to-filter mapping
type stringFilterField struct {
	columnName  string
	filterField func(f *types.FilterCriteria) *[]types.StringFilter
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
		columnName: "finding_arn",
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
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &inspector2.ListFindingsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filter := &types.FilterCriteria{}

	buildStringFilters(d, filter)
	buildNumberFilters(d, filter)
	buildDateFilters(d, filter)

	input.FilterCriteria = filter

	paginator := inspector2.NewListFindingsPaginator(svc, input, func(o *inspector2.ListFindingsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
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
