package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsEcsTasks(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_tasks",
		Description: "AWS ECS Tasks",
		List: &plugin.ListConfig{
			Hydrate:           listEcsTasks,
			ParentHydrate:     listEcsClusters,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "cluster_arn",
					Require: plugin.Optional,
				},
				{
					Name:    "cluster_name",
					Require: plugin.Optional,
				},
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
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
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
				Transform:   transform.FromField("Task.DesiredStatus"),
			},
			{
				Name:        "launch_type",
				Description: "The infrastructure on which your task is running.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.LaunchType"),
			},
			{
				Name:        "task_arn",
				Description: "The Amazon Resource Name (ARN) of the task.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.TaskArn"),
			},
			{
				Name:        "availability_zone",
				Description: "The availability zone of the task.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.AvailabilityZone"),
			},
			{
				Name:        "capacity_provider_name",
				Description: "The capacity provider associated with the task.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.CapacityProviderName"),
			},
			{
				Name:        "cluster_arn",
				Description: "The ARN of the cluster that hosts the task.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.ClusterArn"),
			},
			{
				Name:        "connectivity",
				Description: "The connectivity status of a task.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.Connectivity"),
			},
			{
				Name:        "connectivity_at",
				Description: "The Unix timestamp for when the task last went into CONNECTED status.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Task.ConnectivityAt"),
			},
			{
				Name:        "container_instance_arn",
				Description: "The ARN of the container instances that host the task.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.ContainerInstanceArn"),
			},
			{
				Name:        "cpu",
				Description: "The number of CPU units used by the task as expressed in a task definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.Cpu"),
			},
			{
				Name:        "created_at",
				Description: "The Unix timestamp for when the task was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Task.CreatedAt"),
			},
			{
				Name:        "enable_execute_command",
				Description: "Whether or not execute command functionality is enabled for this task. If true, this enables execute command functionality on all containers in the task.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("Task.EnableExecuteCommand"),
			},
			{
				Name:        "execution_stopped_at",
				Description: "The Unix timestamp for when the task execution stopped.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Task.ExecutionStoppedAt"),
			},
			{
				Name:        "group",
				Description: "The name of the task group associated with the task.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.Group"),
			},
			{
				Name:        "health_status",
				Description: "The health status for the task, which is determined by the health of the essential containers in the task. If all essential containers in the task are reporting as HEALTHY, then the task status also reports as HEALTHY.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.HealthStatus"),
			},
			{
				Name:        "last_status",
				Description: "The last known status of the task.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.LastStatus"),
			},
			{
				Name:        "memory",
				Description: "The amount of memory (in MiB) used by the task as expressed in a task definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.Memory"),
			},
			{
				Name:        "platform_version",
				Description: "The platform version on which your task is running.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.PlatformVersion"),
			},
			{
				Name:        "pull_started_at",
				Description: "The Unix timestamp for when the container image pull began.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Task.PullStartedAt"),
			},
			{
				Name:        "pull_stopped_at",
				Description: "The Unix timestamp for when the container image pull completed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Task.PullStoppedAt"),
			},
			{
				Name:        "service_name",
				Description: "The name of the service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceName"),
			},
			{
				Name:        "started_at",
				Description: "The Unix timestamp for when the task started.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Task.StartedAt"),
			},
			{
				Name:        "started_by",
				Description: "The tag specified when a task is started.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.StartedBy"),
			},
			{
				Name:        "stop_code",
				Description: "The stop code indicating why a task was stopped.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.StopCode"),
			},
			{
				Name:        "stopped_at",
				Description: "The Unix timestamp for when the task was stopped.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Task.StoppedAt"),
			},
			{
				Name:        "stopped_reason",
				Description: "The reason that the task was stopped.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.StoppedReason"),
			},
			{
				Name:        "stopping_at",
				Description: "The Unix timestamp for when the task stops.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Task.StoppingAt"),
			},
			{
				Name:        "task_definition_arn",
				Description: "The ARN of the task definition that creates the task.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Task.TaskDefinitionArn"),
			},
			{
				Name:        "version",
				Description: "The version counter for the task.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Task.Version"),
			},
			{
				Name:        "attachments",
				Description: "The Elastic Network Adapter associated with the task if the task uses the awsvpc network mode.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Task.Attachments"),
			},
			{
				Name:        "attributes",
				Description: "The attributes of the task.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Task.Attributes"),
			},
			{
				Name:        "containers",
				Description: "The containers associated with the task.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Task.Containers"),
			},
			{
				Name:        "ephemeral_storage",
				Description: "The ephemeral storage settings for the task.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Task.EphemeralStorage"),
			},
			{
				Name:        "inference_accelerators",
				Description: "The Elastic Inference accelerator associated with the task.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Task.InferenceAccelerators"),
			},
			{
				Name:        "overrides",
				Description: "One or more container overrides.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Task.Overrides"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Task.Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Task.TaskArn").Transform(arnToAkas),
			},
		}),
	}
}

type tasksInfo struct {
	Task *ecs.Task
	ServiceName string
}

//// LIST FUNCTION

func listEcsTasks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listAwsEBSSnapshots", "AWS_REGION", region)
	equalQuals := d.KeyColumnQuals

	// Create session
	svc, err := EcsService(ctx, d)
	if err != nil {
		return nil, err
	}
	
	var cluster, serviceName string

	if equalQuals["cluster_arn"] != nil {
		cluster = equalQuals["cluster_arn"].GetStringValue()
	} else if equalQuals["cluster_name"] != nil {
		cluster = equalQuals["cluster_name"].GetStringValue()
	} else if h.Item != nil {
		cluster = *h.Item.(*ecs.Cluster).ClusterArn
	}

	// Prepare input parameters
	input := ecs.ListTasksInput{
		MaxResults: types.Int64(100),
	}

	if types.String(cluster) != nil {
		input.Cluster = types.String(cluster)
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
			input.MaxResults = limit
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
			Tasks:              arn,
			Include:            []*string{aws.String("TAGS")},
		}

		if types.String(cluster) != nil {
			input.Cluster = types.String(cluster)
		}

		result, err := svc.DescribeTasks(input)

		if err != nil {
			plugin.Logger(ctx).Error("listECSTasks", "DescribeTasks_error", err)
			return nil, err
		}

		for _, task := range result.Tasks {
			d.StreamListItem(ctx, tasksInfo{task, serviceName})
		}

	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func extractClusterName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	task := d.HydrateItem.(tasksInfo).Task
	clusterName := strings.Split(string(*task.ClusterArn), "/")[1]

	return clusterName, nil
}
