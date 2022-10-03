package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsCloudwatchLogMetricFilter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_metric_filter",
		Description: "AWS CloudWatch Log Metric Filter",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getCloudwatchLogMetricFilter,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudwatchLogMetricFilters,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "log_group_name",
					Require: plugin.Optional,
				},
				{
					Name:    "metric_transformation_name",
					Require: plugin.Optional,
				},
				{
					Name:    "metric_transformation_namespace",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the metric filter",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FilterName"),
			},
			{
				Name:        "log_group_name",
				Description: "The name of the log group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time of the metric filter",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationTime").Transform(convertTimestamp),
			},
			{
				Name:        "filter_pattern",
				Description: "A symbolic description of how CloudWatch Logs should interpret the data in each log event",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metric_transformation_name",
				Description: "The name of the CloudWatch metric",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(logMetricTransformationsData, "MetricName"),
			},
			{
				Name:        "metric_transformation_namespace",
				Description: "A custom namespace to contain metric in CloudWatch. Namespaces are used to group together metrics that are similar",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(logMetricTransformationsData, "MetricNamespace"),
			},
			{
				Name:        "metric_transformation_value",
				Description: "The value to publish to the CloudWatch metric when a filter pattern matches a log event",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(logMetricTransformationsData, "MetricValue"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FilterName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudwatchLogMetricFilterAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudwatchLogMetricFilters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &cloudwatchlogs.DescribeMetricFiltersInput{
		Limit: aws.Int64(50),
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.FilterNamePrefix = aws.String(equalQuals["name"].GetStringValue())
	}
	if equalQuals["log_group_name"] != nil {
		input.LogGroupName = aws.String(equalQuals["log_group_name"].GetStringValue())
	}
	if equalQuals["metric_transformation_name"] != nil {
		input.MetricName = aws.String(equalQuals["metric_transformation_name"].GetStringValue())
	}
	if equalQuals["metric_transformation_namespace"] != nil {
		input.MetricNamespace = aws.String(equalQuals["metric_transformation_namespace"].GetStringValue())
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.Limit {
			if *limit < 1 {
				input.Limit = aws.Int64(1)
			} else {
				input.Limit = limit
			}
		}
	}

	err = svc.DescribeMetricFiltersPages(
		input,
		func(page *cloudwatchlogs.DescribeMetricFiltersOutput, isLast bool) bool {
			for _, metricFilter := range page.MetricFilters {
				d.StreamListItem(ctx, metricFilter)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudwatchLogMetricFilter(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudwatchLogMetricFilter")

	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &cloudwatchlogs.DescribeMetricFiltersInput{
		FilterNamePrefix: &name,
	}

	// execute list call
	op, err := svc.DescribeMetricFilters(params)
	if err != nil {
		return nil, err
	}

	for _, metricFilter := range op.MetricFilters {
		if *metricFilter.FilterName == name {
			plugin.Logger(ctx).Trace("getCloudwatchLogMetricFilter", "FilterName", metricFilter)
			return metricFilter, nil
		}
	}
	return nil, nil
}

func getCloudwatchLogMetricFilterAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudwatchLogGroup")
	region := d.KeyColumnQualString(matrixKeyRegion)
	metricFilter := h.Item.(*cloudwatchlogs.MetricFilter)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":logs:" + region + ":" + commonColumnData.AccountId + ":log-group:" + *metricFilter.LogGroupName + ":metric-filter:" + *metricFilter.FilterName}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func logMetricTransformationsData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("logMetricTransformationsData")
	metricFilterData := d.HydrateItem.(*cloudwatchlogs.MetricFilter)

	if metricFilterData.MetricTransformations != nil && len(metricFilterData.MetricTransformations) > 0 {
		if d.Param.(string) == "MetricName" {
			return metricFilterData.MetricTransformations[0].MetricName, nil
		} else if d.Param.(string) == "MetricNamespace" {
			return metricFilterData.MetricTransformations[0].MetricNamespace, nil
		} else {
			return metricFilterData.MetricTransformations[0].MetricValue, nil
		}
	}
	return nil, nil
}
