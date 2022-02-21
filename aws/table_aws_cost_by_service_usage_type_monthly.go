package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableAwsCostByServiceUsageTypeMonthly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_service_usage_type_monthly",
		Description: "AWS Cost Explorer - Cost by Service and Usage Type (Monthly)",
		List: &plugin.ListConfig{
			Hydrate: listCostByServiceAndUsageMonthly,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "service", Operators: []string{"="}, Require: plugin.Optional},
				{Name: "usage_type", Operators: []string{"="}, Require: plugin.Optional},
			},
		},
		Columns: awsColumns(
			costExplorerColumns([]*plugin.Column{
				{
					Name:        "service",
					Description: "The name of the AWS service.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1"),
				},
				{
					Name:        "usage_type",
					Description: "The usage type of this metric.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension2"),
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostByServiceAndUsageMonthly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Info("listCostByServiceAndUsageMonthly", "ARE WE HERE", "BEFORE PARAM")
	params := buildCostByServiceAndUsageInput(ctx, "MONTHLY", d.Quals, d.Table.List.KeyColumns)
	plugin.Logger(ctx).Info("listCostByServiceAndUsageMonthly", "ARE WE HERE", "AFTER PARAM")
	return streamCostAndUsage(ctx, d, params)
}

func buildCostByServiceAndUsageInput(ctx context.Context, granularity string, quals plugin.KeyColumnQualMap, KeyColumns plugin.KeyColumnSlice) *costexplorer.GetCostAndUsageInput {
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := getCEStartDateForGranularity(granularity).Format(timeFormat)

	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: aws.String(granularity),
		Metrics:     aws.StringSlice(AllCostMetrics()),
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("USAGE_TYPE"),
			},
		},
	}

	var filters []*costexplorer.Expression

	for _, keyQual := range KeyColumns {
		filterQual := quals[keyQual.Name]
		if filterQual == nil {
			continue
		}
		for _, qual := range filterQual.Quals {
			if qual.Value != nil {
				value := qual.Value

				filter := &costexplorer.Expression{}
				filter.Dimensions = &costexplorer.DimensionValues{}
				filter.Dimensions.Key = aws.String(strings.ToUpper(keyQual.Name))

				// Affected by the BUG - https://github.com/turbot/steampipe-plugin-sdk/issues/239
				//
				// filterVal := []string{}
				// if value.GetListValue() != nil {
				// 	for _, q := range value.GetListValue().Values {
				// 		filterVal = append(filterVal, q.GetStringValue())
				// 	}
				// } else {
				// 	plugin.Logger(ctx).Info("buildCostByServiceAndUsageInput", "SINGLE VALUE", "filterVal")
				// 	plugin.Logger(ctx).Info("buildCostByServiceAndUsageInput", "value.GetStringValue()", value.GetStringValue())
				// 	plugin.Logger(ctx).Info("buildCostByServiceAndUsageInput", "value.GetListValue()", value.GetListValue())
				// 	filterVal = append(filterVal, value.GetStringValue())
				// }
				// filter.Dimensions.Values = aws.StringSlice(filterVal)

				filter.Dimensions.Values = aws.StringSlice([]string{value.GetStringValue()})
				filters = append(filters, filter)
			}
		}
	}

	if len(filters) > 1 {
		plugin.Logger(ctx).Info("buildCostByServiceAndUsageInput", "len(filters) > 1", "LIST")
		params.Filter = &costexplorer.Expression{
			And: filters,
		}
	} else if len(filters) == 1 {
		plugin.Logger(ctx).Info("buildCostByServiceAndUsageInput", "len(filters) == 1", "SINGLE")
		params.Filter = filters[0]
	}

	return params
}
