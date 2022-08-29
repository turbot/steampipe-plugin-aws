package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsEcsTask(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_task",
		Description: "AWS ECS Task",
		List: &plugin.ListConfig{
			Hydrate:       listEcsTasks,
			ParentHydrate: listEcsClusters,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "container_instance_arn",
					Require: plugin.Optional,
				},
				{
					Name:    "desired_status",
					Require: plugin.Optional,
				},
				{
					Name:    "launch_type",
					Require: plugin.Optional,
				},
				{
					Name:    "service_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "task_arn",
				Description: "The Amazon Resource Name (ARN) of the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "container_instance_arn",
				Description: "The ARN of the container instances that host the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_name",
				Description: "A user-generated string that you use to identify your cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(extractClusterName),
			},
			{
				Name:        "desired_status",
				Description: "The desired status of the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "launch_type",
				Description: "The infrastructure on which your task is running.",
				Type:        proto.ColumnType_STRING,
			},

			{
				Name:        "availability_zone",
				Description: "The availability zone of the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capacity_provider_name",
				Description: "The capacity provider associated with the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_arn",
				Description: "The ARN of the cluster that hosts the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connectivity",
				Description: "The connectivity status of a task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connectivity_at",
				Description: "The Unix timestamp for when the task last went into CONNECTED status.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "cpu",
				Description: "The number of CPU units used by the task as expressed in a task definition.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "created_at",
				Description: "The Unix timestamp for when the task was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "enable_execute_command",
				Description: "Whether or not execute command functionality is enabled for this task. If true, this enables execute command functionality on all containers in the task.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "execution_stopped_at",
				Description: "The Unix timestamp for when the task execution stopped.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "group",
				Description: "The name of the task group associated with the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_status",
				Description: "The health status for the task, which is determined by the health of the essential containers in the task. If all essential containers in the task are reporting as HEALTHY, then the task status also reports as HEALTHY.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_status",
				Description: "The last known status of the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "memory",
				Description: "The amount of memory (in MiB) used by the task as expressed in a task definition.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "platform_version",
				Description: "The platform version on which your task is running.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pull_started_at",
				Description: "The Unix timestamp for when the container image pull began.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "pull_stopped_at",
				Description: "The Unix timestamp for when the container image pull completed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "service_name",
				Description: "The name of the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "started_at",
				Description: "The Unix timestamp for when the task started.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "started_by",
				Description: "The tag specified when a task is started.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stop_code",
				Description: "The stop code indicating why a task was stopped.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stopped_at",
				Description: "The Unix timestamp for when the task was stopped.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "stopped_reason",
				Description: "The reason that the task was stopped.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stopping_at",
				Description: "The Unix timestamp for when the task stops.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "task_definition_arn",
				Description: "The ARN of the task definition that creates the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The version counter for the task.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "attachments",
				Description: "The Elastic Network Adapter associated with the task if the task uses the awsvpc network mode.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "attributes",
				Description: "The attributes of the task.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "containers",
				Description: "The containers associated with the task.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ephemeral_storage",
				Description: "The ephemeral storage settings for the task.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "inference_accelerators",
				Description: "The Elastic Inference accelerator associated with the task.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "overrides",
				Description: "One or more container overrides.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with task.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(ecsTaskTags),
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

type tasksInfo struct {
	ecs.Task
	ServiceName string
}

//// LIST FUNCTION

func listEcsTasks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listEcsTasks")
	equalQuals := d.KeyColumnQuals

	// Create session
	svc, err := EcsService(ctx, d)
	if err != nil {
		return nil, err
	}

	var serviceName string

	clusterArn := h.Item.(*ecs.Cluster).ClusterArn

	// Prepare input parameters
	input := ecs.ListTasksInput{
		MaxResults: types.Int64(100),
		Cluster:    clusterArn,
	}

	if equalQuals["service_name"] != nil {
		serviceName = equalQuals["service_name"].GetStringValue()
		input.ServiceName = types.String(serviceName)
	}
	if equalQuals["container_instance_arn"] != nil {
		containerInstanceArn := equalQuals["container_instance_arn"].GetStringValue()
		input.ContainerInstance = types.String(containerInstanceArn)
	}
	if equalQuals["desired_status"] != nil {
		desiredStatus := equalQuals["desired_status"].GetStringValue()
		input.DesiredStatus = types.String(desiredStatus)
	}
	if equalQuals["launch_type"] != nil {
		launchType := equalQuals["launch_type"].GetStringValue()
		input.LaunchType = types.String(launchType)
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	var taskArns [][]*string

	// execute list call
	err = svc.ListTasksPages(
		&input,
		func(page *ecs.ListTasksOutput, isLast bool) bool {
			if len(page.TaskArns) != 0 {
				taskArns = append(taskArns, page.TaskArns)
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listECSTasks", "ListTasksPages_error", err)
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ServiceNotFoundException" {
				return nil, nil
			} else if a.Code() == "InvalidParameterException" {
				return nil, nil
			} else if a.Code() == "ClusterNotFoundException" {
				return nil, nil
			}
		}
		return nil, err
	}

	for _, arn := range taskArns {
		input := &ecs.DescribeTasksInput{
			Cluster: clusterArn,
			Tasks:   arn,
			Include: []*string{aws.String("TAGS")},
		}

		result, err := svc.DescribeTasks(input)

		if err != nil {
			plugin.Logger(ctx).Error("listECSTasks", "DescribeTasks_error", err)
			return nil, err
		}

		for _, task := range result.Tasks {
			d.StreamListItem(ctx, tasksInfo{*task, serviceName})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func extractClusterName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	task := d.HydrateItem.(tasksInfo).Task
	clusterName := strings.Split(string(*task.ClusterArn), "/")[1]

	return clusterName, nil
}

func ecsTaskTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	task := d.HydrateItem.(tasksInfo).Task

	if task.Tags == nil {
		return nil, nil
	}

	turbotTagsMap := map[string]string{}
	for _, i := range task.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return turbotTagsMap, nil
}
