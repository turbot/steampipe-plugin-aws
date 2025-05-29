package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAwsCostByResourceHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_resource_hourly",
		Description: "AWS Cost Explorer - Cost by Resource (Hourly)",
		List: &plugin.ListConfig{
			Hydrate: listCostByResourceHourly,
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

func listCostByResourceHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCostByResourceDaily", "client_error", err)
		return nil, err
	}

	params := buildCostByResourceInput("HOURLY", d)

	// We must have to provide a single filter value to make the API call
	if params.Filter == nil {
		// default filter value
		defaultFilter, err := getDefaultFilterValue(ctx, d, h)
		if err != nil {
			return nil, err
		}
		params.Filter = defaultFilter
	}

	// List call
	for {
		output, err := svc.GetCostAndUsageWithResources(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("listCostByResourceDaily", "api_error", err)
			return nil, err
		}

		// Stream the results
		for _, row := range buildCEMetricRows(ctx, (*costexplorer.GetCostAndUsageOutput)(output), nil) {
			d.StreamListItem(ctx, row)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Get more pages if there are any
		if output.NextPageToken == nil {
			break
		}
		params.NextPageToken = output.NextPageToken
	}

	return nil, nil
}

