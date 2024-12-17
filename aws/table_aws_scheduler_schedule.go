package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"

	schedulerv1 "github.com/aws/aws-sdk-go/service/scheduler"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSchedulerSchedule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_scheduler_schedule",
		Description: "AWS EventBridge Scheduler Schedule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "group_name"}),
			Hydrate:    getAwsSchedulerSchedule,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags: map[string]string{"service": "scheduler", "action": "GetSchedule"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSchedulerSchedules,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "group_name", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "scheduler", "action": "ListSchedules"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(schedulerv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the schedule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the schedule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "Specifies whether the schedule is enabled or disabled.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State"),
			},
			{
				Name:        "description",
				Description: "The description of the schedule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSchedulerSchedule,
			},
			{
				Name:        "schedule_expression",
				Description: "The expression that defines when the schedule runs.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSchedulerSchedule,
			},
			{
				Name:        "schedule_expression_timezone",
				Description: "The timezone in which the scheduling expression is evaluated.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSchedulerSchedule,
			},
			{
				Name:        "creation_date",
				Description: "The time at which the schedule was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_modification_date",
				Description: "The time at which the schedule was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "start_date",
				Description: "The date, in UTC, after which the schedule can begin invoking its target.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsSchedulerSchedule,
			},
			{
				Name:        "end_date",
				Description: "The date, in UTC, before which the schedule can invoke its target.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsSchedulerSchedule,
			},
			{
				Name:        "group_name",
				Description: "The name of the schedule group associated with this schedule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_arn",
				Description: "The ARN for a customer managed KMS Key used for encryption.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSchedulerSchedule,
			},
			{
				Name:        "action_after_completion",
				Description: "Indicates the action that EventBridge Scheduler applies after completion.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSchedulerSchedule,
				Transform:   transform.FromField("ActionAfterCompletion"),
			},
			{
				Name:        "flexible_time_window",
				Description: "Allows you to configure a time window during which the schedule invokes the target.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSchedulerSchedule,
			},
			{
				Name:        "target",
				Description: "The schedule target.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSchedulerSchedule,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSchedulerSchedules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := SchedulerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_scheduler_schedule.listAwsSchedulerSchedules", "client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	// Build params
	params := &scheduler.ListSchedulesInput{}
	if d.EqualsQuals["group_name"] != nil {
		params.GroupName = aws.String(d.EqualsQuals["group_name"].GetStringValue())
	}

	paginator := scheduler.NewListSchedulesPaginator(svc, params, func(o *scheduler.ListSchedulesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// Iterate and stream results
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_scheduler_schedule.listAwsSchedulerSchedules", "api_error", err)
			return nil, err
		}

		for _, schedule := range page.Schedules {
			d.StreamListItem(ctx, schedule)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// GET FUNCTION

func getAwsSchedulerSchedule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	groupName := d.EqualsQuals["group_name"].GetStringValue()

	if h.Item != nil {
		scheduler := h.Item.(types.ScheduleSummary)
		name = *scheduler.Name
		groupName = *scheduler.GroupName
	}

	// Empty Check
	if name == "" || groupName == "" {
		return nil, nil
	}

	// Create session
	svc, err := SchedulerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_scheduler_schedule.getAwsSchedulerSchedule", "client_error", err)
		return nil, err
	}

	// Build params
	params := &scheduler.GetScheduleInput{
		Name:      aws.String(name),
		GroupName: aws.String(groupName),
	}

	// Call API
	result, err := svc.GetSchedule(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_scheduler_schedule.getAwsSchedulerSchedule", "api_error", err)
		return nil, err
	}

	return result, nil
}
