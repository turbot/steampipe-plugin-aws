package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIoTFleetMetric(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iot_fleet_metric",
		Description: "AWS IoT Fleet Metric",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("metric_name"),
			Hydrate:    getIoTFleetMetric,
			Tags:       map[string]string{"service": "iot", "action": "DescribeFleetMetric"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listIoTFleetMetrics,
			Tags:    map[string]string{"service": "iot", "action": "ListFleetMetrics"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getIoTFleetMetric,
				Tags: map[string]string{"service": "iot", "action": "DescribeFleetMetric"},
			},
			{
				Func: getIoTFleetMetricTags,
				Tags: map[string]string{"service": "iot", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_IOT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "metric_name",
				Description: "The name of the fleet metric to describe.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the fleet metric to describe.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MetricArn"),
			},
			{
				Name:        "index_name",
				Description: "The name of the index to search.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTFleetMetric,
			},
			{
				Name:        "description",
				Description: "The fleet metric description.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTFleetMetric,
			},
			{
				Name:        "creation_date",
				Description: "The date when the fleet metric is created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getIoTFleetMetric,
				Transform:   transform.FromField("CreationDate"),
			},
			{
				Name:        "last_modified_date",
				Description: "The date when the fleet metric is last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getIoTFleetMetric,
			},
			{
				Name:        "aggregation_field",
				Description: "The field to aggregate.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTFleetMetric,
			},
			{
				Name:        "aggregation_type_name",
				Description: "The name of the aggregation type.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTFleetMetric,
				Transform:   transform.FromField("AggregationType.Name"),
			},
			{
				Name:        "period",
				Description: "The time in seconds between fleet metric emissions. Range [60(1 min), 86400(1 day)] and must be multiple of 60.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getIoTFleetMetric,
			},
			{
				Name:        "query_string",
				Description: "The search query string.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTFleetMetric,
			},
			{
				Name:        "query_version",
				Description: "The search query version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTFleetMetric,
			},
			{
				Name:        "unit",
				Description: "Used to support unit transformation such as milliseconds to seconds. The unit must be supported by CW metric (https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_MetricDatum.html)",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTFleetMetric,
			},
			{
				Name:        "version",
				Description: "The version of the fleet metric.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getIoTFleetMetric,
			},
			{
				Name:        "aggregation_type_values",
				Description: "A list of the values of aggregation types.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTFleetMetric,
				Transform:   transform.FromField("AggregationType.Values"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags currently associated with the thing type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTFleetMetricTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MetricName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTFleetMetricTags,
				Transform:   transform.From(iotFleetMetricTagListToTagsMap),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MetricArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listIoTFleetMetrics(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IoTClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_fleet_metric.listIoTFleetMetrics", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(250)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &iot.ListFleetMetricsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := iot.NewListFleetMetricsPaginator(svc, input, func(o *iot.ListFleetMetricsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iot_fleet_metric.listIoTFleetMetrics", "api_error", err)
			return nil, err
		}

		for _, item := range output.FleetMetrics {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getIoTFleetMetric(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	metricName := ""
	if h.Item != nil {
		t := h.Item.(types.FleetMetricNameAndArn)
		metricName = *t.MetricName
	} else {
		metricName = d.EqualsQualString("metric_name")
	}

	if metricName == "" {
		return nil, nil
	}

	// Create service
	svc, err := IoTClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_fleet_metric.getIoTFleetMetric", "connection_error", err)
		return nil, err
	}

	params := &iot.DescribeFleetMetricInput{
		MetricName: aws.String(metricName),
	}

	resp, err := svc.DescribeFleetMetric(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_fleet_metric.getIoTFleetMetric", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func getIoTFleetMetricTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	typeArn := ""
	switch item := h.Item.(type) {
	case *iot.DescribeThingTypeOutput:
		typeArn = *item.ThingTypeArn
	case types.ThingTypeDefinition:
		typeArn = *item.ThingTypeArn
	}

	// Create service
	svc, err := IoTClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_fleet_metric.getIoTFleetMetricTags", "connection_error", err)
		return nil, err
	}

	params := &iot.ListTagsForResourceInput{
		ResourceArn: aws.String(typeArn),
	}

	endpointTags, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_fleet_metric.getIoTFleetMetricTags", "api_error", err)
		return nil, err
	}

	return endpointTags, nil
}

//// TRANSFORM FUNCTIONS

func iotFleetMetricTagListToTagsMap(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*iot.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	if data.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range data.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
