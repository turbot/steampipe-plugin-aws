package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	ecsEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcsCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_cluster",
		Description: "AWS ECS Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cluster_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidParameterException"}),
			},
			Hydrate: getEcsCluster,
			Tags:    map[string]string{"service": "ecs", "action": "DescribeClusters"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEcsClusters,
			Tags:    map[string]string{"service": "ecs", "action": "ListClusters"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getEcsCluster,
				Tags: map[string]string{"service": "ecs", "action": "DescribeClusters"},
			},
			{
				Func: getAwsEcsClusterTags,
				Tags: map[string]string{"service": "ecs", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ecsEndpoint.AWS_ECS_SERVICE_ID),
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
				Name:        "service_connect_defaults_namespace",
				Description: "The namespace name or full Amazon Resource Name (ARN) of the Cloud Map namespace. When you create a service and don't specify a Service Connect configuration, this namespace is used.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsCluster,
				Transform:   transform.FromField("ServiceConnectDefaults.Namespace"),
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
				Name:        "execute_command_configuration",
				Description: "The execute command configuration for the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsCluster,
				Transform:   transform.FromField("Configuration.ExecuteCommandConfiguration"),
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
	// Create Session
	svc, err := ECSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_cluster.listEcsClusters", "connection_error", err)
		return nil, err
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

	input := &ecs.ListClustersInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := ecs.NewListClustersPaginator(svc, input, func(o *ecs.ListClustersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ecs_cluster.listEcsClusters", "api_error", err)
			return nil, err
		}

		for _, items := range output.ClusterArns {
			d.StreamListItem(ctx, types.Cluster{
				ClusterArn: aws.String(items),
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEcsCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var clusterArn string
	if h.Item != nil {
		clusterArn = *h.Item.(types.Cluster).ClusterArn
	} else {
		quals := d.EqualsQuals
		clusterArn = quals["cluster_arn"].GetStringValue()
	}

	// Create Session
	svc, err := ECSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_cluster.getEcsCluster", "connection_error", err)
		return nil, err
	}

	params := &ecs.DescribeClustersInput{
		Clusters: []string{clusterArn},
		Include: []types.ClusterField{
			types.ClusterFieldAttachments,
			types.ClusterFieldSettings,
			types.ClusterFieldStatistics,
		},
	}

	op, err := svc.DescribeClusters(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_cluster.getEcsCluster", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.Clusters) > 0 {
		return op.Clusters[0], nil
	}

	return nil, nil
}

func getAwsEcsClusterTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	clusterArn := *h.Item.(types.Cluster).ClusterArn

	// Create service
	svc, err := ECSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_cluster.getAwsEcsClusterTags", "connection_error", err)
		return nil, err
	}

	params := &ecs.ListTagsForResourceInput{
		ResourceArn: &clusterArn,
	}

	clusterdata, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_cluster.getAwsEcsClusterTags", "api_error", err)
		return nil, err
	}

	return clusterdata, nil
}

//// TRANSFORM FUNCTIONS

func getAwsEcsClusterTurbotTags(_ context.Context, d *transform.TransformData) (interface{},
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
