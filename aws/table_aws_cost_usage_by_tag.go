package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableAwsCostAndUsageByTag(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_usage_by_tags",
		Description: "AWS Cost Explorer - Cost and Usage By Tags",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "granularity", Require: plugin.Required},
				{Name: "tag_key_1", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "tag_key_2", Require: plugin.Optional, Operators: []string{"="}},
			},
			Hydrate: listCostAndUsageByTags,
		},
		Columns: awsColumns(
			costExplorerColumns([]*plugin.Column{

				// Quals columns - to filter the lookups
				{
					Name:        "granularity",
					Description: "",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "tag_key_1",
					Description: "",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "tag_key_2",
					Description: "",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostAndUsageByTags(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildInputFromTagKeyAndTagValueQuals(ctx, d)
	return streamCostAndUsage(ctx, d, params)
}

func buildInputFromTagKeyAndTagValueQuals(ctx context.Context, d *plugin.QueryData) *costexplorer.GetCostAndUsageInput {
	granularity := strings.ToUpper(d.KeyColumnQualString("granularity"))
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := getCEStartDateForGranularity(granularity).Format(timeFormat)

	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.Granularity(granularity),
		Metrics:     AllCostMetrics(),
	}
	var groupsBy []types.GroupDefinition

	for _, keyQual := range d.Table.List.KeyColumns {
		filterQual := d.Quals[keyQual.Name]
		if filterQual == nil {
			continue
		}
		if keyQual.Name == "tag_key_1" || keyQual.Name == "tag_key_2" {
			for _, qual := range filterQual.Quals {
				if qual.Value != nil {
					value := qual.Value
					groupBy := types.GroupDefinition{
						Type: "TAG",
						Key:  aws.String(value.GetStringValue()),
					}
					groupsBy = append(groupsBy, groupBy)
				}
			}
		}
	}
	if len(groupsBy) > 0 {
		params.GroupBy = groupsBy
	}

	return params
}
