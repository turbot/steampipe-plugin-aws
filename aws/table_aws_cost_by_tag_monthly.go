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

func tableAwsCostByTagMonthly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_tag_monthly",
		Description: "AWS Cost Explorer - Cost by Tag (Monthly)",
		List: &plugin.ListConfig{
			Hydrate: listCostByTagMonthly,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "tag_key_1", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "tag_key_2", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: "exact"},
			},
		},
		Columns: awsColumns(
			costExplorerColumns([]*plugin.Column{
				{
					Name:        "tag_key_1",
					Description: "The tag key to group by",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromQual("tag_key_1"),
				},
				{
					Name:        "tag_key_2",
					Description: "A secondary tag key to group by",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromQual("tag_key_2"),
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

func listCostByTagMonthly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	params := buildCostByTagInput("MONTHLY", d)
	return streamCostAndUsage(ctx, d, params)
}

func buildCostByTagInput(granularity string, d *plugin.QueryData) *costexplorer.GetCostAndUsageInput {
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := getCEStartDateForGranularity(granularity).Format(timeFormat)

	//
	var groupings []types.GroupDefinition

	tag1 := d.KeyColumnQuals["tag_key_1"].GetStringValue()
	if tag1 != "" {
		groupings = append(groupings, buildTagGroupingDefinition(tag1))
	}

	tag2 := d.KeyColumnQuals["tag_key_2"].GetStringValue()
	if tag2 != "" {
		groupings = append(groupings, buildTagGroupingDefinition(tag2))
	}

	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.Granularity(granularity),
		Metrics:     AllCostMetrics(),
		GroupBy:     groupings,
	}

	//
	// params := &costexplorer.GetCostAndUsageInput{
	// 	TimePeriod: &types.DateInterval{
	// 		Start: aws.String(startTime),
	// 		End:   aws.String(endTime),
	// 	},
	// 	Granularity: types.Granularity(granularity),
	// 	Metrics:     AllCostMetrics(),
	// 	GroupBy: []types.GroupDefinition{
	// 		{
	// 			Type: types.GroupDefinitionType("TAG"),
	// 			Key:  aws.String("Cost Center"),
	// 		},
	// 		{
	// 			Type: types.GroupDefinitionType("DIMENSION"),
	// 			Key:  aws.String("USAGE_TYPE"),
	// 		},
	// 	},
	// }

	// var filters []types.Expression

	// for _, keyQual := range d.Table.List.KeyColumns {
	// 	filterQual := d.Quals[keyQual.Name]
	// 	if filterQual == nil {
	// 		continue
	// 	}
	// 	for _, qual := range filterQual.Quals {
	// 		if qual.Value != nil {
	// 			value := qual.Value

	// 			filter := &types.Expression{}
	// 			filter.Dimensions = &types.DimensionValues{}
	// 			filter.Dimensions.Key = types.Dimension(strings.ToUpper(keyQual.Name))

	// 			switch qual.Operator {
	// 			case "=":
	// 				filter := types.Expression{}
	// 				filter.Dimensions = &types.DimensionValues{}
	// 				filter.Dimensions.Key = types.Dimension(strings.ToUpper(keyQual.Name))
	// 				filter.Dimensions.Values = []string{value.GetStringValue()}
	// 				filters = append(filters, filter)
	// 			case "<>":
	// 				filter := types.Expression{}
	// 				filter.Not = &types.Expression{}
	// 				filter.Not.Dimensions = &types.DimensionValues{}
	// 				filter.Not.Dimensions.Key = types.Dimension(strings.ToUpper(keyQual.Name))
	// 				filter.Not.Dimensions.Values = []string{value.GetStringValue()}
	// 				filters = append(filters, filter)
	// 			}
	// 		}
	// 	}
	// }

	// if len(filters) > 1 {
	// 	params.Filter = &types.Expression{
	// 		And: filters,
	// 	}
	// } else if len(filters) == 1 {
	// 	params.Filter = &(filters[0])
	// }

	return params
}

func buildTagGroupingDefinition(tagKey string) types.GroupDefinition {
	return types.GroupDefinition{
		Type: types.GroupDefinitionType("TAG"),
		Key:  aws.String(tagKey),
	}
}

//// TRANSFORM FUNCTIONS

func splitCETagValue(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("splitCETagValue")

	// get the value of policy safely
	tagString := goKitTypes.SafeString(d.Value)

	tag := strings.Split(tagString, "$")

	if len(tag) == 1 {
		return nil, nil
	}

	return strings.Join(tag[1:], "$"), nil //tag[1], nil
}
