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

func tableAwsCloudWatchAlarm(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_alarm",
		Description: "AWS CloudWatch Alarm",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getCloudWatchAlarm,
			Tags:       map[string]string{"service": "cloudwatch", "action": "DescribeAlarms"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchAlarms,
			Tags:    map[string]string{"service": "cloudwatch", "action": "DescribeAlarms"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "state_value",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsCloudWatchAlarmTags,
				Tags: map[string]string{"service": "cloudwatch", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the alarm.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlarmName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the alarm.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlarmArn"),
			},
			{
				Name:        "state_value",
				Description: "The state value for the alarm.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actions_enabled",
				Description: "Indicates whether actions should be executed during any changes to the alarm state.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "alarm_configuration_updated_timestamp",
				Description: "The time stamp of the last update to the alarm configuration.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "alarm_description",
				Description: "The description of the alarm.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "comparison_operator",
				Description: "The arithmetic operation to use when comparing the specified statistic and threshold. The specified statistic value is used as the first operand.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "datapoints_to_alarm",
				Description: "The number of data points that must be breaching to trigger the alarm.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "evaluate_low_sample_count_percentile",
				Description: "Used only for alarms based on percentiles.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "evaluation_periods",
				Description: "The number of periods over which data is compared to the specified threshold.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "extended_statistic",
				Description: "The percentile statistic for the metric associated with the alarm. Specify a value between p0.0 and p100.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metric_name",
				Description: "The name of the metric associated with the alarm, if this is an alarm based on a single metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace",
				Description: "The namespace of the metric associated with the alarm.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "period",
				Description: "The period, in seconds, over which the statistic is applied.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "state_reason",
				Description: "An explanation for the alarm state, in text format.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_reason_data",
				Description: "An explanation for the alarm state, in JSON format.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_updated_timestamp",
				Description: "The time stamp of the last update to the alarm state.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "statistic",
				Description: "The statistic for the metric associated with the alarm, other than percentile.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "threshold",
				Description: "The value to compare with the specified statistic.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "threshold_metric_id",
				Description: "In an alarm based on an anomaly detection model, this is the ID of the ANOMALY_DETECTION_BAND function used as the threshold for the alarm.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "treat_missing_data",
				Description: "Sets how this alarm is to handle missing data points. If this parameter is omitted, the default behavior of missing is used.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "unit",
				Description: "The unit of the metric associated with the alarm.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "alarm_actions",
				Description: "The actions to execute when this alarm transitions to the ALARM state from any other state. Each action is specified as an Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dimensions",
				Description: "The dimensions for the metric associated with the alarm.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "insufficient_data_actions",
				Description: "The actions to execute when this alarm transitions to the INSUFFICIENT_DATA state from any other state. Each action is specified as an Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metrics",
				Description: "An array of MetricDataQuery structures, used in an alarm based on a metric math expression.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ok_actions",
				Description: "The actions to execute when this alarm transitions to the OK state from any other state. Each action is specified as an Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("OKActions"),
			},
			{
				Name:        "tags_src",
				Description: "The list of tag keys and values associated with alarm.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsCloudWatchAlarmTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlarmName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsCloudWatchAlarmTags,
				Transform:   transform.From(getAwsCloudWatchAlarmTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AlarmArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudWatchAlarms(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CloudWatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_alarm.listCloudWatchAlarms", "get_client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	params := &cloudwatch.DescribeAlarmsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["name"] != nil {
		params.AlarmNames = []string{(equalQuals["name"].GetStringValue())}
	}
	if equalQuals["state_value"] != nil {
		params.StateValue = types.StateValue(equalQuals["state_value"].GetStringValue())
	}

	paginator := cloudwatch.NewDescribeAlarmsPaginator(svc, params, func(o *cloudwatch.DescribeAlarmsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_alarm.listCloudWatchAlarms", "api_error", err)
			return nil, err
		}
		for _, alarms := range output.MetricAlarms {
			d.StreamListItem(ctx, alarms)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudWatchAlarm(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	quals := d.EqualsQuals
	name := quals["name"].GetStringValue()

	// Create session
	svc, err := CloudWatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_alarm.getCloudWatchAlarm", "get_client_error", err)
		return nil, err
	}

	params := &cloudwatch.DescribeAlarmsInput{
		AlarmNames: []string{name},
	}

	// execute list call
	item, err := svc.DescribeAlarms(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_alarm.getCloudWatchAlarm", "api_error", err)
		return nil, err
	}

	if item.MetricAlarms != nil && len(item.MetricAlarms) > 0 {
		return item.MetricAlarms[0], nil
	}

	return nil, nil
}

func getAwsCloudWatchAlarmTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	alarm := h.Item.(types.MetricAlarm)

	// Create session
	svc, err := CloudWatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_alarm.getAwsCloudWatchAlarmTags", "client_error", err)
		return nil, err
	}

	params := &cloudwatch.ListTagsForResourceInput{
		ResourceARN: alarm.AlarmArn,
	}

	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_alarm.getAwsCloudWatchAlarmTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getAwsCloudWatchAlarmTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.HydrateItem.(*cloudwatch.ListTagsForResourceOutput)

	if len(tagList.Tags) == 0 {
		return nil, nil
	}
	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
