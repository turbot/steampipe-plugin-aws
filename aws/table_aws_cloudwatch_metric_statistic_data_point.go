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
					Name:       "timestamp",
					Operators:  []string{">", ">=", "=", "<", "<="},
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "period",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "dimensions",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "unit",
					Require:    plugin.Optional,
					CacheMatch: "exact",
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
		Statistics: []types.Statistic{types.StatisticSampleCount, types.StatisticAverage, types.StatisticSum, types.StatisticMinimum, types.StatisticMaximum},
	}

	if d.EqualsQuals["unit"] != nil {
		params.Unit = types.StandardUnit(d.EqualsQuals["unit"].GetStringValue())
	}

	//set the start and end time based on the provided timestamp
	if d.Quals["timestamp"] != nil {
		for _, q := range d.Quals["timestamp"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case "=":
				params.StartTime = aws.Time(timestamp)
				params.EndTime = aws.Time(timestamp)
			case ">=", ">":
				params.StartTime = aws.Time(timestamp)
			case "<", "<=":
				params.EndTime = aws.Time(timestamp)
			}
		}
	}

	if params.StartTime == nil {
		params.StartTime = aws.Time(time.Now().AddDate(0, 0, -1))
	}
	if params.EndTime == nil {
		params.EndTime = aws.Time(time.Now())
	}

	// set the period based on the duration between the start and end time
	// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudwatch@v1.25.1#GetMetricStatisticsInput.Period
	// here we have tried setting the period in such a way that it could provide a good spread under 1440 datapoints

	// for example with 5 days duration the maximum datapoints could be (5 * 24 * 3600) = 432000
	// now due to API limitation of 1440, as per the below calculation, period will be 432000/1440 = 300 and with this period we will get upto 1440 datapoints

	// another example, for a 5 days 15 hours duration the maximum datapoints could be ((5 * 24 + 15) * 3600) = 486000
	// now due to API limitation of 1440, as per the below calculation, period will be ((486000/1440)/60 + 1)*60 = 360
	// in this case 486000/1440 = 337, which is not multiple of 60, so the closest multiple of 60 after 337 is 360
	// with this period we will get upto 1350 datapoints

	// 1 hour - period is 60 sec
	// 6 hours - period is 60 sec
	// 1 day  - period is 60 sec
	// 5 days  - period is 300 sec
	// 7 days - period is 420 sec
	// 15 days - period is 900 sec
	// 30 days - period is 1800 sec
	// 60 days - period is 3600 sec
	// 63 days - period is 3780 sec
	// 90 days - period is 5400 sec

	duration := params.EndTime.Sub(*params.StartTime).Hours()
	durationSec := int32(duration) * 3600
	defaultPeriod := (int32(duration) * 3600) / 1440

	if duration <= 360 { // if the duration is under 15 days
		if int32(durationSec)%1440 == 0 {
			if defaultPeriod < 60 {
				params.Period = aws.Int32(60)
			} else {
				params.Period = aws.Int32(defaultPeriod)
			}
		} else {
			params.Period = aws.Int32((defaultPeriod/60 + 1) * 60)
		}
	} else if duration <= 1512 { // if the duration is between 15 and 63 days
		if int32(durationSec)%1440 == 0 {
			if defaultPeriod < 300 {
				params.Period = aws.Int32(300)
			} else {
				params.Period = aws.Int32(defaultPeriod)
			}
		} else {
			params.Period = aws.Int32((defaultPeriod/300 + 1) * 300)
		}
	} else { // if the duration is greater than 63 days
		if int32(durationSec)%1440 == 0 {
			if defaultPeriod < 3600 {
				params.Period = aws.Int32(3600)
			} else {
				params.Period = aws.Int32(defaultPeriod)
			}
		} else {
			params.Period = aws.Int32((defaultPeriod/3600 + 1) * 3600)
		}
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
			plugin.Logger(ctx).Error("listCloudWatchMetricStatisticDataPoints", "unmarshal_error", err)
			return nil, fmt.Errorf("failed to unmarshal dimensions %v: %v", dimensionsString, err)
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

	for _, datapoints := range statistics.Datapoints {
		d.StreamListItem(ctx, &MetricStatistics{params.MetricName, params.Namespace, params.Period, statistics.Label, dimensions, datapoints})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
