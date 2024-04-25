package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCostByServiceMonthly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_service_monthly",
		Description: "AWS Cost Explorer - Cost by Service (Monthly)",
		List: &plugin.ListConfig{
			Hydrate: listCostByServiceMonthly,
			Tags:    map[string]string{"service": "ce", "action": "GetCostAndUsage"},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:      "service",
					Operators: []string{"=", "<>"},
					Require:   plugin.Optional,
				},
				{
					Name:       "metrics",
					Require:    plugin.Optional,
					Operators:  []string{"="},
					CacheMatch: "exact",
				},
				{
					Name:       "search_start_time",
					Require:    plugin.Optional,
					Operators:  []string{">", ">=", "=", "<", "<="},
					CacheMatch: "exact",
				},
				{
					Name:       "search_end_time",
					Require:    plugin.Optional,
					Operators:  []string{">", ">=", "=", "<", "<="},
					CacheMatch: "exact",
				},
			},
		},
		Columns: awsGlobalRegionColumns(
			costExplorerColumns(
				searchByTimeAndMetricColumns([]*plugin.Column{
					{
						Name:        "service",
						Description: "The name of the AWS service.",
						Type:        proto.ColumnType_STRING,
						Transform:   transform.FromField("Dimension1"),
					},
				}),
			),
		),
	}
}

//// LIST FUNCTION

func listCostByServiceMonthly(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildCostByServiceInput("MONTHLY", d)
	return streamCostAndUsage(ctx, d, params)
}

func buildCostByServiceInput(granularity string, d *plugin.QueryData) *costexplorer.GetCostAndUsageInput {
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := getCEStartDateForGranularity(granularity).Format(timeFormat)

	st, et := getSearchStartTImeAndSearchEndTime(d, granularity)
	if st != "" {
		startTime = st
	}
	if et != "" {
		endTime = et
	}

	selectedMetrics := AllCostMetrics()
	if d.EqualsQualString("metrics") != "" {
		m := getCostMetricByMetricName(d.EqualsQualString("metrics"))
		if !(len(m) > 0) {
			panic(fmt.Sprintf("unsupported metric '%s', supported metrics are %s", d.EqualsQualString("metrics"), strings.Join(selectedMetrics, ",")))
		}

		selectedMetrics = m
	}

	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.Granularity(granularity),
		Metrics:     selectedMetrics,
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionType("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
		},
	}

	var filters []types.Expression

	for _, keyQual := range d.Table.List.KeyColumns {
		filterQual := d.Quals[keyQual.Name]
		if filterQual == nil {
			continue
		}
		for _, qual := range filterQual.Quals {
			if qual.Value != nil {
				value := qual.Value
				switch qual.Operator {
				case "=":
					filter := types.Expression{}
					filter.Dimensions = &types.DimensionValues{}
					filter.Dimensions.Key = types.Dimension(strings.ToUpper(keyQual.Name))
					filter.Dimensions.Values = []string{value.GetStringValue()}
					filters = append(filters, filter)
				case "<>":
					filter := types.Expression{}
					filter.Not = &types.Expression{}
					filter.Not.Dimensions = &types.DimensionValues{}
					filter.Not.Dimensions.Key = types.Dimension(strings.ToUpper(keyQual.Name))
					filter.Not.Dimensions.Values = []string{value.GetStringValue()}
					filters = append(filters, filter)
				}
			}
		}
	}

	if len(filters) > 1 {
		params.Filter = &types.Expression{
			And: filters,
		}
	} else if len(filters) == 1 {
		params.Filter = &(filters[0])
	}

	return params
}
