package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsEcsContainerInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_container",
		Description: "AWS ECS Container",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("container_instance_arn"),
			Hydrate:    getEcsContainerInstance,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEcsClusters,
			Hydrate:       listEcsContainerInstances,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "container_instance_arn",
				Description: "The namespace Amazon Resource Name (ARN) of the cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("ContainerInstanceArn"),
			},
			{
				Name:        "instance_id",
				Description: "The EC2 instance ID of the container instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("Ec2InstanceId"),
			},
			{
				Name:        "agent_connected",
				Description: "True if the agent is connected to Amazon ECS.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("AgentConnected"),
			},
			{
				Name:        "agent_update_status",
				Description: "The status of the most recent agent update.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("AgentUpdateStatus"),
			},
			{
				Name:        "attachments",
				Description: "The resources attached to a container instance, such as elastic network interfaces.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("Attachments"),
			},
			{
				Name:        "attributes",
				Description: "The attributes set for the container instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("Attributes"),
			},
			{
				Name:        "capacity_provider_name",
				Description: "The capacity provider associated with the container instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("CapacityProviderName"),
			},
			{
				Name:        "pending_tasks_count",
				Description: "The number of tasks on the container instance that are in the PENDING status.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEcsContainerInstance,
				Transform: transform.FromField("PendingTasksCount"),
			},
			{
				Name:        "registered_at",
				Description: "The Unix timestamp for when the container instance was registered.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("RegisteredAt"),
			},
			{
				Name:        "registered_resources",
				Description: "CPU and memory that can be allocated on this container instance to tasks.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("RegisteredResources"),
			},
			{
				Name:        "remaining_resources",
				Description: "CPU and memory that is available for new tasks.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("RemainingResources"),
			},
			{
				Name:        "running_tasks_count",
				Description: "CPU and memory that is available for new tasks.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("RunningTasksCount"),
			},
			{
				Name:        "status",
				Description: "The status of the container instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("Status"),
			},
			{
				Name:        "status_reason",
				Description: "The reason that the container instance reached its current status.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("StatusReason"),
			},
			{
				Name:        "version",
				Description: "The reason that the container instance reached its current status.",
				Type:      proto.ColumnType_INT,
				Hydrate:   getEcsContainerInstance,
				Transform: transform.FromField("Version"),
			},
			{
				Name:        "version_info",
				Description: "Version information for the Amazon ECS container agent and Docker daemon running on the container instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("VersionInfo"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.From(containerInstanceTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("ContainerInstanceArn"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsContainerInstance,
				Transform:   transform.FromField("ContainerInstanceArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listEcsContainerInstances(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listEcsContainerInstances")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := EcsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// EcsContainerInstance is a sub resource of an EcsCluster, we need the cluster ARN to list these.
	var clusterArn string
	if h.Item != nil {
		clusterArn = *h.Item.(*ecs.Cluster).ClusterArn
	} else {
		clusterArn = d.KeyColumnQuals["cluster_arn"].GetStringValue()
	}

	// execute list call
	err = svc.ListContainerInstancesPages(
		&ecs.ListContainerInstancesInput{
			Cluster: aws.String(clusterArn),
		},
		func(page *ecs.ListContainerInstancesOutput, isLast bool) bool {
			for _, arn := range page.ContainerInstanceArns {
				d.StreamListItem(ctx, &ecs.ContainerInstance{
					ContainerInstanceArn: arn,
				})
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getEcsContainerInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEcsContainerInstance")
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var containerInstanceArn string
	if h.Item != nil {
		containerInstanceArn = *h.Item.(*ecs.ContainerInstance).ContainerInstanceArn
	} else {
		containerInstanceArn = d.KeyColumnQuals["container_instances_arn"].GetStringValue()
	}

	// This isn't returned from the ListContainerInstances API and is needed to run DescribeContainerInstancesInput.
	// We can however construct it manually from the containerInstanceArn
	clusterArn := clusterArnFromContainerInstanceArn(containerInstanceArn)

	// Create Session
	svc, err := EcsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &ecs.DescribeContainerInstancesInput{
		Cluster: aws.String(clusterArn),
		ContainerInstances: []*string{aws.String(containerInstanceArn)},
		Include: []*string{aws.String("TAGS")},
	}
	result, err := svc.DescribeContainerInstances(input)
	if err != nil {
		return nil, err
	}

	if result.ContainerInstances != nil && len(result.ContainerInstances) > 0 {
		return result.ContainerInstances[0], nil
	}

	return nil, nil
}

// clusterArnFromContainerInstanceArn returns the ClusterArn associated with a given container instance.
func clusterArnFromContainerInstanceArn(arn string) string {
	// Example ARN: arn:aws:ecs:us-east-1:111111111111:container-instance/test
	parts := strings.Split(arn, ":")
	resource := strings.Split(parts[5], "/")
	resource[0] = "cluster"

	parts[5] = strings.Join(resource[0:2], "/")
	return strings.Join(parts, ":")
}

//// TRANSFORM FUNCTIONS

func containerInstanceTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("containerInstanceTagsToTurbotTags")
	tags := d.HydrateItem.(*ecs.ContainerInstance).Tags

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
