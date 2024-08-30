package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAwsCostForecastMonthly(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_forecast_monthly",
		Description: "AWS Cost Explorer - Cost Forecast (Monthly)",
		List: &plugin.ListConfig{
			Hydrate: listCostForecastMonthly,
			Tags:    map[string]string{"service": "ce", "action": "GetCostForecast"},
			KeyColumns: plugin.KeyColumnSlice{
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
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "period_start",
				Description: "Start timestamp for this cost metric",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimePeriod.Start"),
			},
			{
				Name:        "period_end",
				Description: "End timestamp for this cost metric",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimePeriod.End"),
			},
			{
				Name:        "mean_value",
				Description: "Average forecasted value",
				Type:        proto.ColumnType_DOUBLE,
			},
		}),
	}
}

//// LIST FUNCTION

func listCostForecastMonthly(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Get client
	svc, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cost_forecast_monthly.listCostForecastMonthly", "client_error", err)
		return nil, err
	}

	params := buildCostForecastInput(d, "MONTHLY")

	output, err := svc.GetCostForecast(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cost_forecast_monthly.listCostForecastMonthly", "api_error", err)
		return nil, err
	}

	// stream the results...
	for _, r := range output.ForecastResultsByTime {
		d.StreamListItem(ctx, r)

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
