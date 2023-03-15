package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	cloudwatchv1 "github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudWatchMetricStatisticDataPoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_metric_statistic_data_point",
		Description: "AWS CloudWatch Metric Statistic Data Point",
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchMetricStatisticDataPoints,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "metric_name",
					Require: plugin.Required,
				},
				{
					Name:    "namespace",
					Require: plugin.Required,
				},
				{
					Name:      "timestamp",
					Operators: []string{">", ">=", "=", "<", "<="},
					Require:   plugin.Required,
				},
				{
					Name:    "period",
					Require: plugin.Optional,
				},
				{
					Name:    "dimensions",
					Require: plugin.Required,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchv1.EndpointsID),
		Columns: awsRegionalColumns(cwMetricColumns([]*plugin.Column{
			{
				Name:        "label",
				Description: "A label for the specified metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "period",
				Description: "The granularity, in seconds, of the returned data points.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "dimensions",
				Description: "The dimensions for the metric.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Label"),
			},
		})),
	}
}

type MetricStatistics struct {
	MetricName *string
	Namespace  *string
	Period     *int32
	Label      *string
	Dimensions []types.Dimension
	types.Datapoint
}

//// LIST FUNCTION

func listCloudWatchMetricStatisticDataPoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// set the input parameters
	params := &cloudwatch.GetMetricStatisticsInput{
		MetricName: aws.String(d.EqualsQuals["metric_name"].GetStringValue()),
		Namespace:  aws.String(d.EqualsQuals["namespace"].GetStringValue()),
		Statistics: types.Statistic.Values(types.StatisticMaximum),
	}

	//set the start and end time based on the provided timestamp
	for _, q := range d.Quals["timestamp"].Quals {
		time := q.Value.GetTimestampValue().AsTime()
		switch q.Operator {
		case "=":
			params.StartTime = aws.Time(time)
			params.EndTime = aws.Time(time)
		case ">=", ">":
			params.StartTime = aws.Time(time)
		case "<", "<=":
			params.EndTime = aws.Time(time)
		}
	}

	// set the period based on the duration between the start and end time
	duration := params.EndTime.Sub(*params.StartTime).Round(time.Hour)
	if duration.Hours() <= 1 {
		params.Period = aws.Int32(5)
	} else if duration.Hours() <= 24 {
		params.Period = aws.Int32(3600)
	} else {
		params.Period = aws.Int32(86400)
	}

	// override the period if user has provided it in query
	if d.EqualsQuals["period"] != nil {
		params.Period = aws.Int32(int32(d.EqualsQuals["period"].GetInt64Value()))
	}

	// set the dimensions
	dimensions := []types.Dimension{}
	dimensionsString := d.EqualsQuals["dimensions"].GetJsonbValue()

	if dimensionsString != "" {
		err := json.Unmarshal([]byte(dimensionsString), &dimensions)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_metric.listCloudWatchMetrics", "unmarshal_error", err)
			return nil, fmt.Errorf("failed to unmarshal dimensions filter %v: %v", dimensionsString, err)
		}
	}

	if len(dimensions) > 0 {
		params.Dimensions = dimensions
	}

	// Get client
	svc, err := CloudWatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCloudWatchMetricStatisticDataPoints", "client_error", err)
		return nil, err
	}

	statistics, err := svc.GetMetricStatistics(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("listCloudWatchMetricStatisticDataPoints", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Error("testDatapoints", len(statistics.Datapoints))
	for _, datapoints := range statistics.Datapoints {
		d.StreamListItem(ctx, &MetricStatistics{params.MetricName, params.Namespace, params.Period, statistics.Label, dimensions, datapoints})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
