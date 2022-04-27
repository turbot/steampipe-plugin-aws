package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsCostForecastMonthly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_forecast_monthly",
		Description: "AWS Cost Explorer - Cost Forecast (Monthly)",
		List: &plugin.ListConfig{
			Hydrate: listCostForecastMonthly,
		},
		Columns: awsColumns([]*plugin.Column{
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
		},
		),
	}
}

//// LIST FUNCTION

func listCostForecastMonthly(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	logger := plugin.Logger(ctx)
	logger.Trace("listCostForecast")

	// Create session
	svc, err := CostExplorerService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := buildCostForecastInput(d.KeyColumnQuals, "MONTHLY")

	output, err := svc.GetCostForecast(params)
	if err != nil {
		logger.Error("listCostForecast", "err", err)
		return nil, err
	}

	// stream the results...
	for _, r := range output.ForecastResultsByTime {
		d.StreamListItem(ctx, r)

		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
