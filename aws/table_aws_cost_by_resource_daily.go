package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAwsCostByResourceDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_resource_daily",
		Description: "AWS Cost Explorer - Cost by Resource (Daily)",
		List: &plugin.ListConfig{
			Hydrate: listCostByResourceDaily,
			Tags:    map[string]string{"service": "ce", "action": "GetCostAndUsageWithResources"},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:      "resource_id",
					Operators: []string{"=", "<>"},
					Require:   plugin.Optional,
				},
				{
					Name:      "dimension_key",
					Operators: []string{"="},
					Require:   plugin.Optional,
				},
				{
					Name:      "dimension_value",
					Operators: []string{"="},
					Require:   plugin.Optional,
				},
				{
					Name:       "period_start",
					Require:    plugin.Optional,
					Operators:  []string{">", ">=", "=", "<", "<="},
					CacheMatch: query_cache.CacheMatchExact,
				},
				{
					Name:       "period_end",
					Require:    plugin.Optional,
					Operators:  []string{">", ">=", "=", "<", "<="},
					CacheMatch: query_cache.CacheMatchExact,
				},
			},
		},
		Columns: awsGlobalRegionColumns(
			costExplorerColumns([]*plugin.Column{
				{
					Name:        "resource_id",
					Description: "The unique identifier for the resource.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1"),
				},
				{
					Name:        "dimension_key",
					Description: "The name of the dimension key.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromQual("dimension_key"),
					Default:     "LINKED_ACCOUNT",
				},
				{
					Name:        "dimension_value",
					Description: "The value of the dimension key.",
					Type:        proto.ColumnType_STRING,
					Hydrate:     getDimensionValue,
					Transform:   transform.FromValue(),
				},
			}),
		),
	}
}

func listCostByResourceDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	params := buildCostByResourceInput("DAILY", d)

	// We must have to provide a single filter value to make the API call
	if params.Filter == nil {
		// default filter value
		defaultFilter, err := getDefaultFilterValue(ctx, d, h)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cost_by_resource_daily.listCostByResourceDaily", "getDefaultFilterValue", err)
			return nil, err
		}
		params.Filter = defaultFilter
	}

	return streamCostAndUsageByResource(ctx, d, params)
}
