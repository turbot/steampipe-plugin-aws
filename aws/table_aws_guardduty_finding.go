package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type FindingInfo = struct {
	guardduty.Finding
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
			KeyColumns: []*plugin.KeyColumn{
				{Name: "detector_id", Require: plugin.Optional},
				{Name: "id", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "type", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := GuardDutyService(ctx, d)
	if err != nil {
		return nil, err
	}

	detectorId := h.Item.(detectorInfo).DetectorID
	equalQuals := d.KeyColumnQuals

	// Minimize the API call with the given detector_id
	if equalQuals["detector_id"] != nil {
		if equalQuals["detector_id"].GetStringValue() != "" {
			if equalQuals["detector_id"].GetStringValue() != "" && equalQuals["detector_id"].GetStringValue() != detectorId {
				return nil, nil
			}
		} else if len(getListValues(equalQuals["detector_id"].GetListValue())) > 0 {
			if !strings.Contains(fmt.Sprint(getListValues(equalQuals["detector_id"].GetListValue())), detectorId) {
				return nil, nil
			}
		}
	}

	var findingIds [][]*string

	input := &guardduty.ListFindingsInput{
		DetectorId: aws.String(detectorId),
		MaxResults: aws.Int64(50),
	}

	filterCriteria := buildGuarddutyFindingFilterCriteria(d.Quals, ctx)
	if filterCriteria != nil {
		input.FindingCriteria = filterCriteria
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

	// execute list call
	err = svc.ListFindingsPages(
		input,
		func(page *guardduty.ListFindingsOutput, isLast bool) bool {
			if len(page.FindingIds) != 0 {
				findingIds = append(findingIds, page.FindingIds)
			}
			return !isLast
		},
	)
	if err != nil {
		return nil, err
	}

	for _, ids := range findingIds {
		param := &guardduty.GetFindingsInput{
			DetectorId: aws.String(detectorId),
			FindingIds: ids,
		}
		result, err := svc.GetFindings(param)
		if err != nil {
			return nil, err
		}

		for _, finding := range result.Findings {
			d.StreamListItem(ctx, FindingInfo{*finding, detectorId})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// UTILITY FUNCTION

// Build guardduty finding filter param
func buildGuarddutyFindingFilterCriteria(quals plugin.KeyColumnQualMap, ctx context.Context) *guardduty.FindingCriteria {

	type FilterKeyMap struct {
		ColumnName string
		FilterName string
		ColumnType string
	}

	filterCtiteria := make(map[string]*guardduty.Condition)
	filterQuals := []FilterKeyMap{
		{"id", "id", "string"},
		{"type", "type", "string"},
	}

	filterValue := &guardduty.Condition{}

	for _, filterMap := range filterQuals {
		if quals[filterMap.ColumnName] != nil {
			for _, q := range quals[filterMap.ColumnName].Quals {
				value := getQualsValueByColumn(quals, filterMap.ColumnName, "string")
				val, ok := value.(string)
				switch q.Operator {
				case "=":
					if ok {
						filterValue.Equals = []*string{aws.String(val)}
					} else {
						filterValue.Equals = value.([]*string)
					}
					filterCtiteria[filterMap.FilterName] = filterValue
				case "<>":
					if ok {
						filterValue.NotEquals = []*string{aws.String(val)}
					} else {
						filterValue.NotEquals = value.([]*string)
					}
					filterCtiteria[filterMap.FilterName] = filterValue
				}

			}
		}
	}
	return &guardduty.FindingCriteria{Criterion: filterCtiteria}
}
