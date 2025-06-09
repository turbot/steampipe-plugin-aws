package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCloudwatchLogMetricFilter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_metric_filter",
		Description: "AWS CloudWatch Log Metric Filter",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getCloudwatchLogMetricFilter,
			Tags:       map[string]string{"service": "logs", "action": "DescribeMetricFilters"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudwatchLogMetricFilters,
			Tags:    map[string]string{"service": "logs", "action": "DescribeMetricFilters"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
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
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_LOGS_SERVICE_ID),
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

			//// Steampipe Standard Columns

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
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_metric_filter.listCloudwatchLogMetricFilters", "client_error", err)
		return nil, err
	}

	input := &cloudwatchlogs.DescribeMetricFiltersInput{
		Limit: aws.Int32(50),
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
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

	maxItems := int32(50)

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	paginator := cloudwatchlogs.NewDescribeMetricFiltersPaginator(svc, input, func(o *cloudwatchlogs.DescribeMetricFiltersPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_log_metric_filter.listCloudwatchLogMetricFilters", "api_error", err)
			return nil, err
		}

		for _, metricFilter := range output.MetricFilters {
			d.StreamListItem(ctx, metricFilter)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudwatchLogMetricFilter(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()

	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_metric_filter.getCloudwatchLogMetricFilter", "client_error", err)
		return nil, err
	}

	params := &cloudwatchlogs.DescribeMetricFiltersInput{
		FilterNamePrefix: &name,
	}

	// execute list call
	op, err := svc.DescribeMetricFilters(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_metric_filter.getCloudwatchLogMetricFilter", "api_error", err)
		return nil, err
	}

	for _, metricFilter := range op.MetricFilters {
		if *metricFilter.FilterName == name {
			return metricFilter, nil
		}
	}
	return nil, nil
}

func getCloudwatchLogMetricFilterAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	metricFilter := h.Item.(types.MetricFilter)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_metric_filter.getCloudwatchLogMetricFilter", "api_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":logs:" + region + ":" + commonColumnData.AccountId + ":log-group:" + *metricFilter.LogGroupName + ":metric-filter:" + *metricFilter.FilterName}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func logMetricTransformationsData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	metricFilterData := d.HydrateItem.(types.MetricFilter)

	if len(metricFilterData.MetricTransformations) > 0 {
		switch param := d.Param.(string); param {
		case "MetricName":
			return metricFilterData.MetricTransformations[0].MetricName, nil
		case "MetricNamespace":
			return metricFilterData.MetricTransformations[0].MetricNamespace, nil
		case "MetricValue":
			return metricFilterData.MetricTransformations[0].MetricValue, nil
		}
	}
	return nil, nil
}
