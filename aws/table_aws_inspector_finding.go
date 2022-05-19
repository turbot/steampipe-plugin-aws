package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/turbot/go-kit/helpers"
	pb "github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInspectorFinding(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector_finding",
		Description: "AWS Inspector Finding",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getInspectorFinding,
		},
		List: &plugin.ListConfig{
			Hydrate:           listInspectorFindings,
			ShouldIgnoreError: isNotFoundError([]string{"InvalidInputException"}),
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "agent_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "auto_scaling_group", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "severity", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the finding.",
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Id"),
			},
			{
				Name:        "arn",
				Description: "The ARN that specifies the finding.",
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Arn"),
			},
			{
				Name:        "agent_id",
				Description: "The ID of the agent that is installed on the EC2 instance where the finding is generated.",
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("Finding.AssetAttributes.AgentId"),
			},
			{
				Name:        "asset_type",
				Description: "The type of the host from which the finding is generated.",
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("Finding.AssetType"),
			},
			{
				Name:        "auto_scaling_group",
				Description: "The Auto Scaling group of the EC2 instance where the finding is generated.",
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("Finding.AssetAttributes.AutoScalingGroup"),
			},
			{
				Name:        "confidence",
				Description: "This data element is currently not used.",
				Type:        pb.ColumnType_INT,
				Transform:   transform.FromField("Finding.Confidence"),
			},
			{
				Name:        "created_at",
				Description: "The time when the finding was generated.",
				Type:        pb.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Finding.CreatedAt"),
			},
			{
				Name:        "updated_at",
				Description: "The time when AddAttributesToFindings is called.",
				Type:        pb.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Finding.UpdatedAt"),
			},
			{
				Name:        "description",
				Description: "The description of the finding.",
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Description"),
			},
			{
				Name:        "indicator_of_compromise",
				Description: "This data element is currently not used.",
				Type:        pb.ColumnType_BOOL,
				Transform:   transform.FromField("Finding.IndicatorOfCompromise"),
			},
			{
				Name:        "numeric_severity",
				Description: "The numeric value of the finding severity.",
				Type:        pb.ColumnType_DOUBLE,
				Transform:   transform.FromField("Finding.NumericSeverity"),
			},
			{
				Name:        "recommendation",
				Description: "The recommendation for the finding.",
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Recommendation"),
			},
			{
				Name:        "schema_version",
				Description: "The schema version of this data type.",
				Type:        pb.ColumnType_INT,
				Transform:   transform.FromField("Finding.SchemaVersion"),
			},
			{
				Name:        "service",
				Description: "The data element is set to 'Inspector'.",
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Service"),
			},
			{
				Name:        "severity",
				Description: "The finding severity. Values can be set to High, Medium, Low, and Informational.",
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Severity"),
			},

			// Json columns
			{
				Name:        "asset_attributes",
				Description: "A collection of attributes of the host from which the finding is generated.",
				Type:        pb.ColumnType_JSON,
				Transform:   transform.FromField("Finding.AssetAttributes"),
			},
			{
				Name:        "attributes",
				Description: "The system-defined attributes for the finding.",
				Type:        pb.ColumnType_JSON,
				Transform:   transform.FromField("Finding.Attributes"),
			},
			{
				Name:        "failed_items",
				Description: "Attributes details that cannot be described. An error code is provided for each failed item.",
				Type:        pb.ColumnType_JSON,
			},
			{
				Name:        "service_attributes",
				Description: "This data type is used in the Finding data type.",
				Type:        pb.ColumnType_JSON,
				Transform:   transform.FromField("Finding.ServiceAttributes"),
			},
			{
				Name:        "user_attributes",
				Description: "The user-defined attributes that are assigned to the finding.",
				Type:        pb.ColumnType_JSON,
				Transform:   transform.FromField("Finding.UserAttributes"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "The name of the finding.",
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        pb.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type InspectorFindingInfo struct {
	FailedItems map[string]*inspector.FailedItemDetails
	Finding     *inspector.Finding
}

//// LIST FUNCTION

func listInspectorFindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	serviceId := endpoints.InspectorServiceID
	region := d.KeyColumnQualString(matrixKeyRegion)
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	// Create session
	svc, err := InspectorService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &inspector.ListFindingsInput{
		MaxResults: aws.Int64(500),
	}

	filterParam := buildListInspectorFindingsParam(d.Quals)
	if filterParam != nil {
		input.Filter = filterParam
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	findingArns := []*string{}

	// List call
	err = svc.ListFindingsPages(
		input,
		func(page *inspector.ListFindingsOutput, isLast bool) bool {
			findingArns = append(findingArns, page.FindingArns...)

			// plugin.Logger(ctx).Error("Total Number of Findings ====>>>>", len(page.FindingArns))
			// if len(page.FindingArns) <= 500 {
			// 	return false
			// }
			return !isLast
		},
	)

	if err != nil {
		return nil, err
	}

	// Skip api call if there is no findings
	if len(findingArns) == 0 {
		return nil, nil
	}

	passedFindingArn := 0
	findingLeft := true
	for findingLeft {

		// DescribeFindings api can take maximum 10 number of repository name at a time.
		var arns []*string
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
		data, err := svc.DescribeFindings(params)

		if err != nil {
			return nil, err
		}

		if data != nil {
			for _, finding := range data.Findings {
				d.StreamListItem(ctx, &InspectorFindingInfo{
					FailedItems: data.FailedItems,
					Finding:     finding,
				})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					break
				}
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getInspectorFinding(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getInspectorFinding")

	serviceId := endpoints.InspectorServiceID
	region := d.KeyColumnQualString(matrixKeyRegion)
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	var findingArn string
	if h.Item != nil {
		findingArn = *h.Item.(*inspector.Finding).Arn
	} else {
		quals := d.KeyColumnQuals
		findingArn = quals["arn"].GetStringValue()
	}

	// get service
	svc, err := InspectorService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &inspector.DescribeFindingsInput{
		FindingArns: []*string{aws.String(findingArn)},
	}

	// Get call
	data, err := svc.DescribeFindings(params)
	if err != nil {
		logger.Debug("getInspectorFinding", "ERROR", err)
		return nil, err
	}
	if data.Findings != nil && len(data.Findings) > 0 {
		return &InspectorFindingInfo{
			FailedItems: data.FailedItems,
			Finding:     data.Findings[0],
		}, nil
	}
	return nil, nil
}

// Build param for findings list call
func buildListInspectorFindingsParam(quals plugin.KeyColumnQualMap) *inspector.FindingFilter {
	inspectorFindingFilter := &inspector.FindingFilter{}

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
				inspectorFindingFilter.AgentIds = []*string{aws.String(value)}
			case "auto_scaling_group":
				inspectorFindingFilter.AutoScalingGroups = []*string{aws.String(value)}
			case "severity":
				inspectorFindingFilter.Severities = []*string{aws.String(value)}
			}

		}
	}

	return inspectorFindingFilter
}
