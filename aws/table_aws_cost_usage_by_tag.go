package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	goKitTypes "github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsCostAndUsageByTag(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_usage_by_tags",
		Description: "AWS Cost Explorer - Cost and Usage By Tags",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "granularity", Require: plugin.Required},
				{Name: "tag_key_1", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "tag_key_2", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: "exact"},
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
					Description: "The tag key to group by.",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "tag_key_2",
					Description: "A secondary tag key to group by.",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "tag_value_1",
					Description: "The primary tag value grouped by",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1").Transform(splitCETagValue),
				},
				{
					Name:        "tag_value_2",
					Description: "A secondary tag value grouped by",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension2").Transform(splitCETagValue),
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostAndUsageByTags(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildInputFromTagKeyAndTagValueQuals(ctx, d)
	// If any of the optional quals are not given, the streamCostAndUsage function still returns the result, indicating that we have a param check for optional quals.
	// If user does not provide any of the tag key then return empty row.
	if len(params.GroupBy) <= 0 {
		return nil, nil
	}
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

//// TRANSFORM FUNCTIONS

func splitCETagValue(_ context.Context, d *transform.TransformData) (interface{}, error) {

	// get the value of policy safely
	tagString := goKitTypes.SafeString(d.Value)

	tag := strings.Split(tagString, "$")

	if len(tag) == 1 {
		return nil, nil
	}

	return strings.Join(tag[1:], "$"), nil
}
