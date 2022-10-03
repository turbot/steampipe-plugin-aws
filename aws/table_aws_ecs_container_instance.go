package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsEcsContainerInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_container_instance",
		Description: "AWS ECS Container Instance",
		List: &plugin.ListConfig{
			ParentHydrate: listEcsClusters,
			Hydrate:       listEcsContainerInstances,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The namespace Amazon Resource Name (ARN) of the container instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerInstanceArn"),
			},
			{
				Name:        "ec2_instance_id",
				Description: "The EC2 instance ID of the container instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_arn",
				Description: "The ARN of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "agent_connected",
				Description: "True if the agent is connected to Amazon ECS.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "agent_update_status",
				Description: "The status of the most recent agent update.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attachments",
				Description: "The resources attached to a container instance, such as elastic network interfaces.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "attributes",
				Description: "The attributes set for the container instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "capacity_provider_name",
				Description: "The capacity provider associated with the container instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pending_tasks_count",
				Description: "The number of tasks on the container instance that are in the PENDING status.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "registered_at",
				Description: "The Unix timestamp for when the container instance was registered.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "registered_resources",
				Description: "CPU and memory that can be allocated on this container instance to tasks.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "remaining_resources",
				Description: "CPU and memory that is available for new tasks.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "running_tasks_count",
				Description: "CPU and memory that is available for new tasks.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "status",
				Description: "The status of the container instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_reason",
				Description: "The reason that the container instance reached its current status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The reason that the container instance reached its current status.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "version_info",
				Description: "Version information for the Amazon ECS container agent and Docker daemon running on the container instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(containerInstanceTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerInstanceArn"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerInstanceArn").Transform(arnToAkas),
			},
		}),
	}
}

type containerInstanceData = struct {
	ecs.ContainerInstance
	ClusterArn string
}

//// LIST FUNCTION

// listEcsContainerInstances handles both listing and describing of the instances.
//
// The reason for this is the DescribeContainerInstance call can accept up to 100 ARNs. If we moved it out to another
// hydrate functions we may save a request or two if we only wanted to retrieve the ARNs but the tradeoff is we need
// to get any other info an API call per container instance would need to be made. So in the case where we need to get
// all info for less then 100 instances including the Describe request here, and batching requests means only making
// two API calls as opposed to 101.
func listEcsContainerInstances(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listEcsContainerInstances")

	// Create Session
	svc, err := EcsService(ctx, d)
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

	// DescribeContainerInstances can accept up to 100 ARNs at a time, so make sure
	// ListContainerInstances returns the same and append to this in chunks not more then 100.
	var containerInstanceArns [][]*string

	input := &ecs.ListContainerInstancesInput{
		Cluster:    aws.String(clusterArn),
		MaxResults: aws.Int64(100),
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

	// execute list call
	err = svc.ListContainerInstancesPages(
		input,
		func(page *ecs.ListContainerInstancesOutput, isLast bool) bool {
			if len(page.ContainerInstanceArns) != 0 {
				containerInstanceArns = append(containerInstanceArns, page.ContainerInstanceArns)
			}
			return !isLast
		},
	)
	if err != nil {
		return nil, err
	}

	for _, arns := range containerInstanceArns {
		input := &ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(clusterArn),
			ContainerInstances: arns,
			Include:            []*string{aws.String("TAGS")},
		}
		result, err := svc.DescribeContainerInstances(input)
		if err != nil {
			return nil, err
		}

		for _, inst := range result.ContainerInstances {
			d.StreamListItem(ctx, containerInstanceData{*inst, clusterArn})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func containerInstanceTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("containerInstanceTagsToTurbotTags")
	tags := d.HydrateItem.(containerInstanceData).Tags

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
