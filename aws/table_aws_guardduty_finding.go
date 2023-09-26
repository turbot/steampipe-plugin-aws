package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"

	guarddutyv1 "github.com/aws/aws-sdk-go/service/guardduty"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type findingInfo struct {
	types.Finding
	DetectorId string
}

//// TABLE DEFINITION

func tableAwsGuardDutyFinding(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_finding",
		Description: "AWS GuardDuty Finding",
		List: &plugin.ListConfig{
			ParentHydrate: listGuardDutyDetectors,
			Hydrate:       listGuardDutyFindings,
			Tags:          map[string]string{"service": "guardduty", "action": "ListFindings"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "detector_id", Require: plugin.Optional},
				{Name: "id", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "type", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(guarddutyv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The title of the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Title"),
			},
			{
				Name:        "id",
				Description: "The ID of the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "detector_id",
				Description: "The ID of the detector.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "severity",
				Description: "The severity of the finding.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "created_at",
				Description: "The time and date when the finding was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "confidence",
				Description: "The confidence score for the finding.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "description",
				Description: "The description of the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schema_version",
				Description: "The version of the schema used for the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The time and date when the finding was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "resource",
				Description: "Contains information about the AWS resource associated with the activity that prompted GuardDuty to generate a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "service",
				Description: "Contains additional information about the generated finding.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

// listGuardDutyFindings handles both listing and get the details of the findings.
func listGuardDutyFindings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_finding.listGuardDutyFindings", "connection_error", err)
		return nil, err
	}

	detectorId := h.Item.(detectorInfo).DetectorID
	equalQuals := d.EqualsQuals
	// Minimize the API call with the given detector_id
	if equalQuals["detector_id"] != nil {
		if equalQuals["detector_id"].GetStringValue() != detectorId {
			return nil, nil
		}
	}

	// var findingIds [][]*string
	maxItems := int32(50)
	input := &guardduty.ListFindingsInput{
		DetectorId: aws.String(detectorId),
	}

	filterCriteria := buildGuarddutyFindingFilterCriteria(d.Quals, ctx)
	if filterCriteria != nil {
		input.FindingCriteria = filterCriteria
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			input.MaxResults = limit
		}
	}

	paginator := guardduty.NewListFindingsPaginator(svc, input, func(o *guardduty.ListFindingsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	var findingIds [][]string
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_guardduty_finding.listGuardDutyFindings", "api_error", err)
			return nil, err
		}
		findingIds = append(findingIds, output.FindingIds)
	}

	// Using this pattern as the GetFindings API supports an array of findings
	for _, ids := range findingIds {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		param := &guardduty.GetFindingsInput{
			DetectorId: aws.String(detectorId),
			FindingIds: ids,
		}
		result, err := svc.GetFindings(ctx, param)
		if err != nil {
			return nil, err
		}

		for _, finding := range result.Findings {
			d.StreamListItem(ctx, findingInfo{finding, detectorId})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

//// UTILITY FUNCTION

// Build guardduty finding filter param
func buildGuarddutyFindingFilterCriteria(quals plugin.KeyColumnQualMap, ctx context.Context) *types.FindingCriteria {

	type FilterKeyMap struct {
		ColumnName string
		FilterName string
		ColumnType string
	}

	filterCtiteria := make(map[string]types.Condition)
	filterQuals := []FilterKeyMap{
		{"id", "id", "string"},
		{"type", "type", "string"},
	}

	filterValue := types.Condition{}

	for _, filterMap := range filterQuals {
		if quals[filterMap.ColumnName] != nil {
			for _, q := range quals[filterMap.ColumnName].Quals {
				value := getQualsValueByColumn(quals, filterMap.ColumnName, "string")
				val, ok := value.(string)
				switch q.Operator {
				case "=":
					if ok {
						filterValue.Equals = []string{val}
					} else {
						filterValue.Equals = value.([]string)
					}
					filterCtiteria[filterMap.FilterName] = filterValue
				case "<>":
					if ok {
						filterValue.NotEquals = []string{val}
					} else {
						filterValue.NotEquals = value.([]string)
					}
					filterCtiteria[filterMap.FilterName] = filterValue
				}

			}
		}
	}
	return &types.FindingCriteria{Criterion: filterCtiteria}
}
