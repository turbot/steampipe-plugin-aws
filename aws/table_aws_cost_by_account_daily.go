package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCostByLinkedAccountDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_account_daily",
		Description: "AWS Cost Explorer - Cost by Linked Account (Daily)",
		List: &plugin.ListConfig{
			Hydrate: listCostByLinkedAccountDaily,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:      "metrics",
					Require:   plugin.Optional,
					Operators: []string{"="},
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

func listCostByLinkedAccountDaily(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildCostByLinkedAccountInput(d, "DAILY")
	return streamCostAndUsage(ctx, d, params)
}
