package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEcsCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_cluster",
		Description: "AWS ECS Cluster",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("cluster_arn"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException", "InvalidParameterException"}),
			Hydrate:           getEcsCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listEcsClusters,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cluster_arn",
				Description: "The Amazon Resource Name (ARN) that identifies the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_name",
				Description: "A user-generated string that you use to identify your cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "active_services_count",
				Description: "The number of services that are running on the cluster in an ACTIVE state.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "pending_tasks_count",
				Description: "The number of tasks in the cluster that are in the PENDING state.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "registered_container_instances_count",
				Description: "The number of container instances registered into the cluster. This includes container instances in both ACTIVE and DRAINING status.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "running_tasks_count",
				Description: "The number of tasks in the cluster that are in the RUNNING state.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "status",
				Description: "The status of the cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "attachments_status",
				Description: "The status of the capacity providers associated with the cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "attachments",
				Description: "The resources attached to a cluster. When using a capacity provider with a cluster, the Auto Scaling plan that is created will be returned as a cluster attachment.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "capacity_providers",
				Description: "The capacity providers associated with the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "default_capacity_provider_strategy",
				Description: "The default capacity provider strategy for the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "settings",
				Description: "The settings for the cluster. This parameter indicates whether CloudWatch Container Insights is enabled or disabled for a cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "statistics",
				Description: "Additional information about your clusters that are separated by launch type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsCluster,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEcsClusterTags,
				Transform:   transform.FromField("Tags"),
			},
			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsCluster,
				Transform:   transform.FromField("ClusterName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEcsClusterTags,
				Transform:   transform.From(getAwsEcsClusterTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listEcsClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

	// List call
	err = svc.ListClustersPages(
		&ecs.ListClustersInput{},
		func(page *ecs.ListClustersOutput, isLast bool) bool {
			for _, results := range page.ClusterArns {
				d.StreamListItem(ctx, &ecs.Cluster{
					ClusterArn: results,
				})
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEcsCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getEcsCluster")
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var clusterArn string
	if h.Item != nil {
		clusterArn = *h.Item.(*ecs.Cluster).ClusterArn
	} else {
		quals := d.KeyColumnQuals
		clusterArn = quals["cluster_arn"].GetStringValue()
	}

	// Create Session
	svc, err := EcsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ecs.DescribeClustersInput{
		Clusters: []*string{aws.String(clusterArn)},
	}

	op, err := svc.DescribeClusters(params)
	if err != nil {
		logger.Debug("getEcsCluster", "ERROR", err)
		return nil, err
	}

	if op.Clusters != nil && len(op.Clusters) > 0 {
		return op.Clusters[0], nil
	}

	return nil, nil
}

func getAwsEcsClusterTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEcsClusterTags")
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	clusterArn := *h.Item.(*ecs.Cluster).ClusterArn

	// Create service
	svc, err := EcsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ecs.ListTagsForResourceInput{
		ResourceArn: &clusterArn,
	}

	clusterdata, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return clusterdata, nil
}

//// TRANSFORM FUNCTIONS

func getAwsEcsClusterTurbotTags(ctx context.Context, d *transform.TransformData) (interface{},
	error) {
	ecsClusterTags := d.HydrateItem.(*ecs.ListTagsForResourceOutput)

	if ecsClusterTags.Tags == nil {
		return nil, nil
	}

	if ecsClusterTags.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range ecsClusterTags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
