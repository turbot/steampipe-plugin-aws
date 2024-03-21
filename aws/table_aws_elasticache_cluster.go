package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"

	elasticachev1 "github.com/aws/aws-sdk-go/service/elasticache"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElastiCacheCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_cluster",
		Description: "AWS ElastiCache Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cache_cluster_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"CacheClusterNotFound", "InvalidParameterValue"}),
			},
			Hydrate: getElastiCacheCluster,
			Tags:    map[string]string{"service": "elasticache", "action": "DescribeCacheClusters"},
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheClusters,
			Tags:    map[string]string{"service": "elasticache", "action": "DescribeCacheClusters"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listTagsForElastiCacheCluster,
				Tags: map[string]string{"service": "elasticache", "action": "ListTagsForResource"},
			},
		},

		GetMatrixItemFunc: SupportedRegionMatrix(elasticachev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cache_cluster_id",
				Description: "An unique identifier for ElastiCache cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) of the cache cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "cache_node_type",
				Description: "The name of the compute and memory capacity node type for the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cache_cluster_status",
				Description: "The current state of this cluster, one of the following values: available, creating, deleted, deleting, incompatible-network, modifying, rebooting cluster nodes, restore-failed, or snapshotting.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "at_rest_encryption_enabled",
				Description: "A flag that enables encryption at-rest when set to true.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "auth_token_enabled",
				Description: "A flag that enables using an AuthToken (password) when issuing Redis commands.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "auto_minor_version_upgrade",
				Description: "This parameter is currently disabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "cache_cluster_create_time",
				Description: "The date and time when the cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "auth_token_last_modified_date",
				Description: "The date the auth token was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "ip_discovery",
				Description: "The network type associated with the cluster, either ipv4 | ipv6.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_type",
				Description: "Must be either ipv4 | ipv6 | dual_stack.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "transit_encryption_mode",
				Description: "A setting that allows you to migrate your clients to use in-transit encryption, with no downtime.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "preferred_outpost_arn",
				Description: "The outpost ARN in which the cache cluster is created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_group_log_delivery_enabled",
				Description: "A boolean value indicating whether log delivery is enabled for the replication group.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "cache_subnet_group_name",
				Description: "The name of the cache subnet group associated with the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "client_download_landing_page",
				Description: "The URL of the web page where you can download the latest ElastiCache client library.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "configuration_endpoint",
				Description: "Represents a Memcached cluster endpoint which can be used by an application to connect to any node in the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine",
				Description: "The name of the cache engine (memcached or redis) to be used for this cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "The version of the cache engine that is used in this cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "num_cache_nodes",
				Description: "The number of cache nodes in the cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "preferred_availability_zone",
				Description: "The name of the Availability Zone in which the cluster is located or 'Multiple' if the cache nodes are located in different Availability Zones.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "preferred_maintenance_window",
				Description: "Specifies the weekly time range during which maintenance on the cluster is performed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_group_id",
				Description: "The replication group to which this cluster belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "snapshot_retention_limit",
				Description: "The number of days for which ElastiCache retains automatic cluster snapshots before deleting them.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "snapshot_window",
				Description: "The daily time range (in UTC) during which ElastiCache begins taking a daily snapshot of your cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "transit_encryption_enabled",
				Description: "A flag that enables in-transit encryption when set to true.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "cache_parameter_group",
				Description: "Status of the cache parameter group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cache_nodes",
				Description: "A list of cache nodes that are members of the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cache_security_groups",
				Description: "A list of cache security group elements, composed of name and status sub-elements.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "log_delivery_configurations",
				Description: "Returns the destination, format and type of the logs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "notification_configuration",
				Description: "Describes a notification topic and its status.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "pending_modified_values",
				Description: "A group of settings that are applied to the cluster in the future, or that are currently being applied.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_groups",
				Description: "A list of VPC Security Groups associated with the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForElastiCacheCluster,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CacheClusterId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForElastiCacheCluster,
				Transform:   transform.From(clusterTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listElastiCacheClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_cluster.listElastiCacheClusters", "connection_error", err)
		return nil, err
	}

	input := &elasticache.DescribeCacheClustersInput{
		MaxRecords: aws.Int32(100),
	}

	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxRecords {
			if limit < 20 {
				input.MaxRecords = aws.Int32(20)
			} else {
				input.MaxRecords = aws.Int32(limit)
			}
		}
	}

	// List call
	paginator := elasticache.NewDescribeCacheClustersPaginator(svc, input, func(o *elasticache.DescribeCacheClustersPaginatorOptions) {
		o.Limit = *input.MaxRecords
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_elasticache_cluster.listElastiCacheClusters", "api_error", err)
			return nil, err
		}

		for _, cacheCluster := range output.CacheClusters {
			d.StreamListItem(ctx, cacheCluster)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getElastiCacheCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_cluster.getElastiCacheCluster", "connection_error", err)
		return nil, err
	}

	cacheClusterID := d.EqualsQuals["cache_cluster_id"].GetStringValue()

	// Return nil, if no input provided
	if cacheClusterID == "" {
		return nil, nil
	}

	params := &elasticache.DescribeCacheClustersInput{
		CacheClusterId: aws.String(cacheClusterID),
	}

	op, err := svc.DescribeCacheClusters(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_cluster.getElastiCacheCluster", "api_error", err)
		return nil, err
	}

	if len(op.CacheClusters) > 0 {
		return op.CacheClusters[0], nil
	}

	return nil, nil
}

func listTagsForElastiCacheCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(types.CacheCluster)

	// Create session
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_cluster.listTagsForElastiCacheCluster", "connection_error", err)
		return nil, err
	}

	// Build param
	param := &elasticache.ListTagsForResourceInput{
		ResourceName: cluster.ARN,
	}

	clusterTags, err := svc.ListTagsForResource(ctx, param)

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "CacheClusterNotFound" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_elasticache_cluster.listTagsForElastiCacheCluster", "api_error", err)
		return nil, err
	}

	return clusterTags, nil
}

//// TRANSFORM FUNCTIONS

func clusterTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	clusterTag := d.HydrateItem.(*elasticache.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if len(clusterTag.TagList) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range clusterTag.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
