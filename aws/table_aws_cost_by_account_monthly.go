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

func tableAwsCostByLinkedAccountMonthly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_account_monthly",
		Description: "AWS Cost Explorer - Cost by Linked Account (Monthly)",
		List: &plugin.ListConfig{
			Hydrate: listCostByLinkedAccountMonthly,
			KeyColumns: plugin.KeyColumnSlice{
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
			Tags: map[string]string{"service": "ce", "action": "GetCostAndUsage"},
		},
		Columns: awsGlobalRegionColumns(
			costExplorerColumns(
				searchByTimeAndMetricColumns([]*plugin.Column{
					{
						Name:        "linked_account_id",
						Description: "The AWS Account ID.",
						Type:        proto.ColumnType_STRING,
						Transform:   transform.FromField("Dimension1"),
					},
				}),
			),
		),
	}
}

//// LIST FUNCTION

func listCostByLinkedAccountMonthly(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildCostByLinkedAccountInput(d, "MONTHLY")
	return streamCostAndUsage(ctx, d, params)
}

func buildCostByLinkedAccountInput(d *plugin.QueryData, granularity string) *costexplorer.GetCostAndUsageInput {
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
				Key:  aws.String("LINKED_ACCOUNT"),
			},
		},
	}

	return params
}
