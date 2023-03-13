package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	cloudwatchv1 "github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudWatchMetricStatistics(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_metric_statistics",
		Description: "AWS CloudWatch Metric Statistics",
		List: &plugin.ListConfig{
			ParentHydrate: listCloudWatchMetrics,
			Hydrate:       listCloudWatchMetricStatistics,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "metric_name",
					Require: plugin.Optional,
				},
				{
					Name:    "namespace",
					Require: plugin.Optional,
				},
				{
					Name:       "end_time",
					Require:    plugin.Required,
					CacheMatch: "exact",
				},
				{
					Name:       "start_time",
					Require:    plugin.Required,
					CacheMatch: "exact",
				},
				{
					Name:       "period",
					Require:    plugin.Required,
					CacheMatch: "exact",
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "metric_name",
				Description: "The name of the metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace",
				Description: "The namespace for the metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_time",
				Description: "The time stamp that determines the first data point to return. Start times are evaluated relative to the time that CloudWatch receives the request. The value specified is inclusive; results include data points with the specified timestamp. In a raw HTTP query, the time stamp must be in ISO 8601 UTC format (for example, 2016-10-03T23:00:00Z).",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromQual("start_time"),
			},
			{
				Name:        "end_time",
				Description: "The time stamp that determines the last data point to return. The value specified is exclusive; results include data points up to the specified timestamp. In a raw HTTP query, the time stamp must be in ISO 8601 UTC format (for example, 2016-10-10T23:00:00Z).",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromQual("end_time"),
			},
			{
				Name:        "label",
				Description: "A label for the specified metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "period",
				Description: "The namespace for the metric.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromQual("period"),
			},
			{
				Name:        "datapoints",
				Description: "The data points for the specified metric.",
				Type:        proto.ColumnType_JSON,
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
		}),
	}
}

type MetricStatistics struct {
	MetricName string
	Namespace  string
	Dimensions []types.Dimension
	cloudwatch.GetMetricStatisticsOutput
}

//// LIST FUNCTION

func listCloudWatchMetricStatistics(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	metricName := *h.Item.(types.Metric).MetricName
	namespace := *h.Item.(types.Metric).Namespace
	dimensions := h.Item.(types.Metric).Dimensions

	// check if the metric_name provided in the where clause is identical with the metric_name from the parent hydrate
	if d.EqualsQuals["metric_name"].GetStringValue() != "" && d.EqualsQuals["metric_name"].GetStringValue() != metricName {
		return nil, nil
	}

	// check if the namespace provided in the where clause is identical with the namespace from the parent hydrate
	if d.EqualsQuals["namespace"].GetStringValue() != "" && d.EqualsQuals["namespace"].GetStringValue() != namespace {
		return nil, nil
	}

	// Get client
	svc, err := CloudWatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCloudWatchMetricStatistics", "client_error", err)
		return nil, err
	}

	params := &cloudwatch.GetMetricStatisticsInput{
		StartTime:  aws.Time(d.EqualsQuals["start_time"].GetTimestampValue().AsTime()),
		EndTime:    aws.Time(d.EqualsQuals["end_time"].GetTimestampValue().AsTime()),
		Period:     aws.Int32(int32(d.EqualsQuals["period"].GetInt64Value())),
		MetricName: aws.String(metricName),
		Namespace:  aws.String(namespace),
		Statistics: types.Statistic.Values(types.StatisticMaximum),
		Dimensions: dimensions,
	}

	statistics, err := svc.GetMetricStatistics(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("listCloudWatchMetricStatistics", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, &MetricStatistics{metricName, namespace, dimensions, *statistics})

	return nil, nil
}
