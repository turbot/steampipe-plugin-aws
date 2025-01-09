package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector"
	"github.com/aws/aws-sdk-go-v2/service/inspector/types"

	inspectorEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInspectorFinding(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector_finding",
		Description: "AWS Inspector Finding",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getInspectorFinding,
			Tags:       map[string]string{"service": "inspector", "action": "ListFindings"},
		},
		List: &plugin.ListConfig{
			Hydrate: listInspectorFindings,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException", "NoSuchEntity", "InvalidParameter"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "agent_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "auto_scaling_group", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "severity", Require: plugin.Optional, Operators: []string{"="}},
			},
			Tags: map[string]string{"service": "inspector", "action": "DescribeFindings"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(inspectorEndpoint.INSPECTORServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Id"),
			},
			{
				Name:        "arn",
				Description: "The ARN that specifies the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Arn"),
			},
			{
				Name:        "agent_id",
				Description: "The ID of the agent that is installed on the EC2 instance where the finding is generated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.AssetAttributes.AgentId"),
			},
			{
				Name:        "asset_type",
				Description: "The type of the host from which the finding is generated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.AssetType"),
			},
			{
				Name:        "auto_scaling_group",
				Description: "The Auto Scaling group of the EC2 instance where the finding is generated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.AssetAttributes.AutoScalingGroup"),
			},
			{
				Name:        "confidence",
				Description: "This data element is currently not used.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Finding.Confidence"),
			},
			{
				Name:        "created_at",
				Description: "The time when the finding was generated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Finding.CreatedAt"),
			},
			{
				Name:        "updated_at",
				Description: "The time when AddAttributesToFindings is called.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Finding.UpdatedAt"),
			},
			{
				Name:        "description",
				Description: "The description of the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Description"),
			},
			{
				Name:        "indicator_of_compromise",
				Description: "This data element is currently not used.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Finding.IndicatorOfCompromise"),
			},
			{
				Name:        "numeric_severity",
				Description: "The numeric value of the finding severity.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Finding.NumericSeverity"),
			},
			{
				Name:        "recommendation",
				Description: "The recommendation for the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Recommendation"),
			},
			{
				Name:        "schema_version",
				Description: "The schema version of this data type.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Finding.SchemaVersion"),
			},
			{
				Name:        "service",
				Description: "The data element is set to 'Inspector'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Service"),
			},
			{
				Name:        "severity",
				Description: "The finding severity. Values can be set to High, Medium, Low, and Informational.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Severity"),
			},

			// Json columns
			{
				Name:        "asset_attributes",
				Description: "A collection of attributes of the host from which the finding is generated.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.AssetAttributes"),
			},
			{
				Name:        "attributes",
				Description: "The system-defined attributes for the finding.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.Attributes"),
			},
			{
				Name:        "failed_items",
				Description: "Attributes details that cannot be described. An error code is provided for each failed item.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "service_attributes",
				Description: "This data type is used in the Finding data type.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.ServiceAttributes"),
			},
			{
				Name:        "user_attributes",
				Description: "The user-defined attributes that are assigned to the finding.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.UserAttributes"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "The name of the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type InspectorFindingInfo struct {
	FailedItems map[string]types.FailedItemDetails
	Finding     types.Finding
}

//// LIST FUNCTION

func listInspectorFindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := InspectorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_finding.listInspectorFindings", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(500)
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

	input := &inspector.ListFindingsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filterParam := buildListInspectorFindingsParam(d.Quals)
	if filterParam != nil {
		input.Filter = filterParam
	}

	findingArns := []string{}

	paginator := inspector.NewListFindingsPaginator(svc, input, func(o *inspector.ListFindingsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_inspector_finding.listInspectorFindings", "api_error", err)
			return nil, err
		}
		findingArns = append(findingArns, output.FindingArns...)
	}

	// Skip api call if there is no findings
	if len(findingArns) == 0 {
		return nil, nil
	}

	passedFindingArn := 0
	findingLeft := true
	for findingLeft {

		// DescribeFindings api can take maximum 10 number of repository name at a time.
		var arns []string
		if len(findingArns) > passedFindingArn {
			if (len(findingArns) - passedFindingArn) >= 10 {
				arns = findingArns[passedFindingArn : passedFindingArn+10]
				passedFindingArn += 10
			} else {
				arns = findingArns[passedFindingArn:]
				findingLeft = false
			}
		}

		// Build the params
		params := &inspector.DescribeFindingsInput{
			FindingArns: arns,
		}

		// Get call
		data, err := svc.DescribeFindings(ctx, params)

		if err != nil {
			plugin.Logger(ctx).Error("aws_inspector_finding.listInspectorFindings.DescribeFindings", "api_error", err)
			return nil, err
		}

		if data != nil {
			for _, finding := range data.Findings {
				d.StreamListItem(ctx, &InspectorFindingInfo{
					FailedItems: data.FailedItems,
					Finding:     finding,
				})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.RowsRemaining(ctx) == 0 {
					break
				}
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getInspectorFinding(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var findingArn string
	if h.Item != nil {
		findingArn = *h.Item.(*InspectorFindingInfo).Finding.Arn
	} else {
		quals := d.EqualsQuals
		findingArn = quals["arn"].GetStringValue()
	}

	// Create Session
	svc, err := InspectorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_finding.getInspectorFinding", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &inspector.DescribeFindingsInput{
		FindingArns: []string{findingArn},
	}

	// Get call
	data, err := svc.DescribeFindings(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_finding.getInspectorFinding", "api_error", err)
		return nil, err
	}
	if data != nil && len(data.Findings) > 0 {
		return &InspectorFindingInfo{
			FailedItems: data.FailedItems,
			Finding:     data.Findings[0],
		}, nil
	}
	return nil, nil
}

// Build param for findings list call
func buildListInspectorFindingsParam(quals plugin.KeyColumnQualMap) *types.FindingFilter {
	inspectorFindingFilter := &types.FindingFilter{}

	strColumns := []string{"agent_id", "auto_scaling_group", "severity"}

	for _, s := range strColumns {
		if quals[s] == nil {
			continue
		}
		for _, q := range quals[s].Quals {
			value := q.Value.GetStringValue()
			if value == "" || q.Operator != "=" {
				continue
			}

			switch s {
			case "agent_id":
				inspectorFindingFilter.AgentIds = []string{value}
			case "auto_scaling_group":
				inspectorFindingFilter.AutoScalingGroups = []string{value}
			case "severity":
				inspectorFindingFilter.Severities = []types.Severity{
					types.Severity(value),
				}
			}

		}
	}

	return inspectorFindingFilter
}
