package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCostForecastDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_forecast_daily",
		Description: "AWS Cost Explorer - Cost Forecast (Daily)",
		List: &plugin.ListConfig{
			Hydrate: listCostForecastDaily,
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
		},
		),
	}
}

//// LIST FUNCTION

func listCostForecastDaily(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Get client
	svc, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cost_forecast_daily.listCostForecastDaily", "client_error", err)
		return nil, err
	}

	params := buildCostForecastInput(d.EqualsQuals, "DAILY")

	output, err := svc.GetCostForecast(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cost_forecast_daily.listCostForecast", "api_error", err)
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

func buildCostForecastInput(_ map[string]*proto.QualValue, granularity string) *costexplorer.GetCostForecastInput {

	// TO DO - specify metric as qual?   get all cost metrics in parallel?
	//metric := strings.ToUpper(keyQuals["metric"].GetStringValue())

	// As the response of the api call doesn't return metric value so we do not have a column for it,
	// If we will have a column for it, then we need to get the value from quals only(we may get the null value if it has not been passed), so we have not added it as optional quals.
	// We can add it as required param, but there is a bug with "in" clause so we cann't iterate the value properly in param.

	metric := "UNBLENDED_COST"

	timeFormat := "2006-01-02"
	startTime := time.Now().UTC().Format(timeFormat)
	endTime := getForecastEndDateForGranularity(granularity).Format(timeFormat)

	params := &costexplorer.GetCostForecastInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.Granularity(granularity),
		Metric:      types.Metric(metric),
	}

	return params
}

func getForecastEndDateForGranularity(granularity string) time.Time {
	switch granularity {
	case "MONTHLY":
		return lastDayOfMonth(12) // 1 year
	case "DAILY":
		return lastDayOfMonth(3) // 3 months
	}
	return lastDayOfMonth(12) // 1 year
}

func lastDayOfMonth(numMonths int) time.Time {
	today := time.Now()
	goneDaysOfMonth := today.Day()

	if goneDaysOfMonth == 1 {
		return today.AddDate(0, numMonths, 0)
	}
	return today.AddDate(0, numMonths+1, -goneDaysOfMonth+1)
}
