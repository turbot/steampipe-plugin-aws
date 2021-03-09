package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
		},
		GetMatrixItem: BuildRegionList,
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
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listCloudwatchLogMetricFilters", "AWS_REGION", region)

	// Create session
	svc, err := CloudWatchLogsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	err = svc.DescribeMetricFiltersPages(
		&cloudwatchlogs.DescribeMetricFiltersInput{},
		func(page *cloudwatchlogs.DescribeMetricFiltersOutput, isLast bool) bool {
			for _, metricFilter := range page.MetricFilters {
				d.StreamListItem(ctx, metricFilter)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudwatchLogMetricFilter(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudwatchLogMetricFilter")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create session
	svc, err := CloudWatchLogsService(ctx, d, region)
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
	metricFilter := h.Item.(*cloudwatchlogs.MetricFilter)

	commonColumnData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonData := commonColumnData.(*awsCommonColumnData)
	// Get data for turbot defined properties
	akas := []string{"arn:" + commonData.Partition + ":logs:" + commonData.Region + ":" + commonData.AccountId + ":log-group:" + *metricFilter.LogGroupName + ":metric-filter:" + *metricFilter.FilterName}

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
