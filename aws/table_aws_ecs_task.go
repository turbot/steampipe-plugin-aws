package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

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
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ClusterNotFoundException", "ServiceNotFoundException", "InvalidParameterException"}),
			},
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
					Name:       "service_name",
					Require:    plugin.Optional,
					CacheMatch: "exact",
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
				Transform:   transform.FromQual("service_name"),
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

//// LIST FUNCTION

func listEcsTasks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	equalQuals := d.KeyColumnQuals
	clusterArn := h.Item.(types.Cluster).ClusterArn

	// Create session
	svc, err := ECSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_task.listEcsTasks", "connection_error", err)
		return nil, err
	}

	var serviceName string

	// Prepare input parameters
	input := ecs.ListTasksInput{Cluster: clusterArn}

	if equalQuals["service_name"] != nil {
		serviceName = equalQuals["service_name"].GetStringValue()
		input.ServiceName = aws.String(serviceName)
	}
	if equalQuals["container_instance_arn"] != nil {
		containerInstanceArn := equalQuals["container_instance_arn"].GetStringValue()
		input.ContainerInstance = aws.String(containerInstanceArn)
	}
	if equalQuals["desired_status"] != nil {
		desiredStatus := equalQuals["desired_status"].GetStringValue()
		input.DesiredStatus = types.DesiredStatus(desiredStatus)
	}
	if equalQuals["launch_type"] != nil {
		launchType := equalQuals["launch_type"].GetStringValue()
		input.LaunchType = types.LaunchType(launchType)
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

	input.MaxResults = &maxLimit
	var taskArns [][]string

	paginator := ecs.NewListTasksPaginator(svc, &input, func(o *ecs.ListTasksPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ecs_task.listEcsTasks", "list_tasks_api_error", err)
			return nil, err
		}

		taskArns = append(taskArns, output.TaskArns)
	}

	for _, arns := range taskArns {
		if len(arns) == 0 {
			continue
		}
		input := &ecs.DescribeTasksInput{
			Cluster: clusterArn,
			Tasks:   arns,
			Include: []types.TaskField{types.TaskFieldTags},
		}

		result, err := svc.DescribeTasks(ctx, input)

		if err != nil {
			plugin.Logger(ctx).Error("aws_ecs_task.listEcsTasks", "describe_tasks_api_error", err)
			return nil, err
		}

		for _, task := range result.Tasks {
			d.StreamListItem(ctx, task)

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
	task := d.HydrateItem.(types.Task)
	clusterName := strings.Split(string(*task.ClusterArn), "/")[1]

	return clusterName, nil
}

func ecsTaskTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	task := d.HydrateItem.(types.Task)

	var turbotTagsMap map[string]string
	if len(task.Tags) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range task.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
