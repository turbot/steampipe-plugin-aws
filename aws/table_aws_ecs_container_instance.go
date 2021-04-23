package aws

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
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
			KeyColumns:        plugin.SingleColumn("container_instance_arn"),
			Hydrate:    getEcsContainerInstance,

		},
		List: &plugin.ListConfig{
			ParentHydrate:  listEcsClusters,
			Hydrate: listEcsContainerInstances,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsS3Columns([]*plugin.Column{
			{
				Name:        "container_instance_arn",
				Description: "The namespace Amazon Resource Name (ARN) of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerInstanceArn"),
			},
			{
				Name:        "instance_id",
				Description: "The EC2 instance ID of the container instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Ec2InstanceId"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "agent_connected",
				Description: "True if the agent is connected to Amazon ECS.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AgentConnected"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "agent_update_status",
				Description: "The status of the most recent agent update.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AgentUpdateStatus"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "attachments",
				Description: "The resources attached to a container instance, such as elastic network interfaces.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Attachments"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "attributes",
				Description: "The attributes set for the container instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Attributes"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "capacity_provider_name",
				Description: "The capacity provider associated with the container instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CapacityProviderName"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "pending_tasks_count",
				Description: "The number of tasks on the container instance that are in the PENDING status.",
				Type:        proto.ColumnType_INT,
				// TODO: this is a 64bit int, make sure ColumnType_INT will work
				Transform:   transform.FromField("PendingTasksCount"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "registered_at",
				Description: "The Unix timestamp for when the container instance was registered.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("RegisteredAt"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "registered_resources",
				Description: "CPU and memory that can be allocated on this container instance to tasks.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RegisteredResources"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "remaining_resources",
				Description: "CPU and memory that is available for new tasks.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RemainingResources"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "running_tasks_count",
				Description: "CPU and memory that is available for new tasks.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("RunningTasksCount"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "status",
				Description: "The status of the container instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "status_reason",
				Description: "The reason that the container instance reached its current status.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StatusReason"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "version",
				Description: "The reason that the container instance reached its current status.",
				// TODO: This is a 64 bit int, make sure we're using the right column type
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Version"),
				Hydrate: getEcsContainerInstance,
			},
			{
				Name:        "version_info",
				Description: "Version information for the Amazon ECS container agent and Docker daemon running on the container instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VersionInfo"),
				Hydrate: getEcsContainerInstance,
			},
			//{
			//	Name:        "tags",
			//	Description: resourceInterfaceDescription("tags"),
			//	Type:        proto.ColumnType_JSON,
			//	Hydrate:     getAwsEcsClusterTags,
			//	Transform:   transform.From(getAwsEcsTurbotTags),
			//},
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

func containerInstanceArnFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	clusterArn := quals["container_instance_arn"].GetStringValue()
	item := &ecs.Cluster{
		ClusterArn: &clusterArn,
	}
	return item, nil
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
				d.StreamListItem(ctx, &ecs.ContainerInstance{ContainerInstanceArn: arn})
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

// do not have a get call for s3 bucket.
// using list api call to create get function
func getEcsContainerInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("describeEcsContainerInstance")
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

//func getContainerInstanceTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
//	plugin.Logger(ctx).Trace("getBucketTagging")
//	bucket := h.Item.(*s3.Bucket)
//	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)
//
//	// Create Session
//	svc, err := S3Service(ctx, d, *location.LocationConstraint)
//	if err != nil {
//		return nil, err
//	}
//
//	params := &s3.GetBucketTaggingInput{
//		Bucket: bucket.Name,
//	}
//
//	bucketTags, _ := svc.GetBucketTagging(params)
//	if err != nil {
//		return nil, err
//	}
//
//	return bucketTags, nil
//}
//

func getContainerInstanceARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEcsContainerInstanceArn")
	return h.Item.(*ecs.ContainerInstance).ContainerInstanceArn, nil
}

//// TRANSFORM FUNCTIONS

func containerInstanceTagsToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("s3TagsToTurbotTags")
	tags := d.Value.([]*s3.Tag)

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
