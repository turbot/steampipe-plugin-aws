package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/datasync"
	"github.com/aws/aws-sdk-go-v2/service/datasync/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsDataSyncTask(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_datasync_task",
		Description: "AWS DataSync Task",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidRequestException"}),
			},
			Hydrate: getDataSyncTask,
			Tags:    map[string]string{"service": "datasync", "action": "DescribeTask"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDataSyncTasks,
			Tags:    map[string]string{"service": "datasync", "action": "ListTasks"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getDataSyncTask,
				Tags: map[string]string{"service": "datasync", "action": "DescribeTask"},
			},
			{
				Func: getDataSyncTaskTags,
				Tags: map[string]string{"service": "datasync", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DATASYNC_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the task.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TaskArn"),
			},
			{
				Name:        "name",
				Description: "The name of the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "task_mode",
				Description: "The task mode that you're using.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_location_arn",
				Description: "The ARN of your transfer's source location.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "destination_location_arn",
				Description: "The ARN of your transfer's destination location.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "cloud_watch_log_group_arn",
				Description: "The ARN of an Amazon CloudWatch log group for monitoring your task.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "creation_time",
				Description: "The time that the task was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "current_task_execution_arn",
				Description: "The ARN of the most recent task execution.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "error_detail",
				Description: "If there's an issue with your task, you can use the error details to help you troubleshoot the problem.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "error_code",
				Description: "If there's an issue with your task, you can use the error code to help you troubleshoot the problem.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "excludes",
				Description: "The exclude filters that define the files, objects, and folders in your source location that you don't want DataSync to transfer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "includes",
				Description: "The include filters that define the files, objects, and folders in your source location that you want DataSync to transfer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "manifest_config",
				Description: "The configuration of the manifest that lists the files or objects that you want DataSync to transfer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "options",
				Description: "The task's settings. For example, what file metadata gets preserved, how data integrity gets verified at the end of your transfer, bandwidth limits, among other options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "schedule",
				Description: "The schedule for when you want your task to run.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "schedule_details",
				Description: "The details about your task schedule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "source_network_interface_arns",
				Description: "The ARNs of the network interfaces that DataSync created for your source location.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "destination_network_interface_arns",
				Description: "The ARNs of the network interfaces that DataSync created for your destination location.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "task_report_config",
				Description: "The configuration of your task report, which provides detailed information about your DataSync transfer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataSyncTask,
			},
			{
				Name:        "tags_src",
				Description: "The tags associated the task.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataSyncTaskTags,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataSyncTaskTags,
				Transform:   transform.From(getTaskTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TaskArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDataSyncTasks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := DataSyncClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_datasync_task.listDataSyncTasks", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(100)
	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < int64(maxItems) {
		maxItems = int32(*d.QueryContext.Limit)
	}

	input := &datasync.ListTasksInput{
		MaxResults: &maxItems,
	}

	paginator := datasync.NewListTasksPaginator(svc, input, func(o *datasync.ListTasksPaginatorOptions) {
		o.Limit = maxItems
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_datasync_task.listDataSyncTasks", "api_error", err)
			return nil, err
		}

		for _, task := range output.Tasks {
			d.StreamListItem(ctx, task)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDataSyncTask(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var taskArn string
	if h.Item != nil {
		task := h.Item.(types.TaskListEntry)
		taskArn = *task.TaskArn
	} else {
		taskArn = d.EqualsQuals["arn"].GetStringValue()
	}

	if taskArn == "" {
		return nil, nil
	}

	// Create service
	svc, err := DataSyncClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_datasync_task.getDataSyncTask", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &datasync.DescribeTaskInput{
		TaskArn: aws.String(taskArn),
	}

	// Get call
	data, err := svc.DescribeTask(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_datasync_task.getDataSyncTask", "api_error", err)
		return nil, err
	}

	return data, nil
}

func getDataSyncTaskTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var taskArn string
	if h.Item != nil {
		task := h.Item.(types.TaskListEntry)
		taskArn = *task.TaskArn
	} else {
		taskArn = d.EqualsQuals["arn"].GetStringValue()
	}

	if taskArn == "" {
		return nil, nil
	}

	// Create service
	svc, err := DataSyncClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_datasync_task.getDataSyncTask", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &datasync.ListTagsForResourceInput{
		ResourceArn: aws.String(taskArn),
	}

	// Get call
	paginator := datasync.NewListTagsForResourcePaginator(svc, params, func(o *datasync.ListTagsForResourcePaginatorOptions) {
		o.Limit = 100
	})

	var tags []types.TagListEntry
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_datasync_task.getDataSyncTaskTags", "api_error", err)
			return nil, err
		}

		tags = append(tags, output.Tags...)
	}

	return tags, nil
}

//// TRANSFORM FUNCTIONS

func getTaskTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem != nil {
		tags := d.HydrateItem.([]types.TagListEntry)
		tagsMap := make(map[string]string)
		for _, tag := range tags {
			tagsMap[*tag.Key] = *tag.Value
		}
		return tagsMap, nil
	}
	return nil, nil
}
