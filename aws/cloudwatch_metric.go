package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

// append the common cloudwatch metric columns onto the column list
func cwMetricColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonCwMetricColumns()...)
}

func commonCwMetricColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "metric_name",
			Description: "The name of the metric.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "namespace",
			Description: "The metric namespace.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "average",
			Description: "The average of the metric values that correspond to the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "maximum",
			Description: "The maximum metric value for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "minimum",
			Description: "The minimum metric value for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "sample_count",
			Description: "The number of metric values that contributed to the aggregate value of this data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "sum",
			Description: "The sum of the metric values for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "unit",
			Description: "The standard unit for the data point.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "timestamp",
			Description: "The time stamp used for the data point.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

type CWMetricRow struct {
	// The (single) metric Dimension name
	DimensionName *string

	// The value for the (single) metric Dimension
	DimensionValue *string

	// The namespace of the metric
	Namespace *string

	// The name of the metric
	MetricName *string

	// The average of the metric values that correspond to the data point.
	Average *float64

	// The percentile statistic for the data point.
	//ExtendedStatistics map[string]*float64 `type:"map"`

	// The maximum metric value for the data point.
	Maximum *float64

	// The minimum metric value for the data point.
	Minimum *float64

	// The number of metric values that contributed to the aggregate value of this
	// data point.
	SampleCount *float64

	// The sum of the metric values for the data point.
	Sum *float64

	// The time stamp used for the data point.
	Timestamp *time.Time

	// The standard unit for the data point.
	Unit *string
}

func getCWStartDateForGranularity(granularity string) time.Time {
	switch strings.ToUpper(granularity) {
	case "DAILY":
		// 1 year
		return time.Now().AddDate(-1, 0, 0)
	case "HOURLY":
		// 60 days
		return time.Now().AddDate(0, 0, -60)
	}
	// else 5 days
	return time.Now().AddDate(0, 0, -5)
}

func getCWPeriodForGranularity(granularity string) int32 {
	switch strings.ToUpper(granularity) {
	case "DAILY":
		// 24 hours
		return 86400
	case "HOURLY":
		// 1 hour
		return 3600
	}
	// else 5 minutes
	return 300
}

func listCWMetricStatistics(ctx context.Context, d *plugin.QueryData, granularity string, namespace string, metricName string, dimensionName string, dimensionValue string) (*cloudwatch.GetMetricStatisticsOutput, error) {
	// Create Session
	svc, err := CloudWatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCWMetricStatistics", "connection_error", err)
		return nil, err
	}

	endTime := time.Now()
	startTime := getCWStartDateForGranularity(granularity)
	period := getCWPeriodForGranularity(granularity)

	params := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String(namespace),
		MetricName: aws.String(metricName),
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(endTime),
		Period:     aws.Int32(period),
		Statistics: []types.Statistic{
			types.StatisticAverage,
			types.StatisticSampleCount,
			types.StatisticSum,
			types.StatisticMinimum,
			types.StatisticMaximum,
		},
	}

	if dimensionName != "" && dimensionValue != "" {
		params.Dimensions = []types.Dimension{
			{
				Name:  aws.String(dimensionName),
				Value: aws.String(dimensionValue),
			},
		}
	}

	stats, err := svc.GetMetricStatistics(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("listCWMetricStatistics", "api_error", err)
		return nil, err
	}

	for _, datapoint := range stats.Datapoints {
		d.StreamLeafListItem(ctx, &CWMetricRow{
			DimensionValue: aws.String(dimensionValue),
			DimensionName:  aws.String(dimensionName),
			Namespace:      aws.String(namespace),
			MetricName:     aws.String(metricName),
			Average:        datapoint.Average,
			Maximum:        datapoint.Maximum,
			Minimum:        datapoint.Minimum,
			Timestamp:      datapoint.Timestamp,
			SampleCount:    datapoint.SampleCount,
			Sum:            datapoint.Sum,
			Unit:           aws.String(fmt.Sprint(datapoint.Unit)),
		})

		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
