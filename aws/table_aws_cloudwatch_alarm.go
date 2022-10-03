package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudWatchAlarm(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_alarm",
		Description: "AWS CloudWatch Alarm",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getCloudWatchAlarm,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchAlarms,
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
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := CloudWatchService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &cloudwatch.DescribeAlarmsInput{
		MaxRecords: aws.Int64(100),
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.AlarmNames = []*string{aws.String(equalQuals["name"].GetStringValue())}
	}
	if equalQuals["state_value"] != nil {
		input.StateValue = aws.String(equalQuals["state_value"].GetStringValue())
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxRecords {
			if *limit < 1 {
				input.MaxRecords = aws.Int64(1)
			} else {
				input.MaxRecords = limit
			}
		}
	}

	// List call
	err = svc.DescribeAlarmsPages(
		input,
		func(page *cloudwatch.DescribeAlarmsOutput, isLast bool) bool {
			for _, alarms := range page.MetricAlarms {
				d.StreamListItem(ctx, alarms)

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

func getCloudWatchAlarm(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getCloudWatchAlarm")
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()

	// Create Session
	svc, err := CloudWatchService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &cloudwatch.DescribeAlarmsInput{
		AlarmNames: []*string{aws.String(name)},
	}

	op, err := svc.DescribeAlarms(params)
	if err != nil {
		logger.Debug("getCloudWatchAlarm", "ERROR", err)
		return nil, err
	}

	if op.MetricAlarms != nil && len(op.MetricAlarms) > 0 {
		return op.MetricAlarms[0], nil
	}

	return nil, nil
}

func getAwsCloudWatchAlarmTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsCloudWatchAlarmTags")
	alarm := h.Item.(*cloudwatch.MetricAlarm)

	// Create service
	svc, err := CloudWatchService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &cloudwatch.ListTagsForResourceInput{
		ResourceARN: alarm.AlarmArn,
	}

	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getAwsCloudWatchAlarmTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	cloudWatchAlarm := d.HydrateItem.(*cloudwatch.ListTagsForResourceOutput)

	if cloudWatchAlarm.Tags == nil {
		return nil, nil
	}

	if cloudWatchAlarm.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range cloudWatchAlarm.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
