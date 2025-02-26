package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	elasticachev1 "github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsElasticacheUpdateAction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_update_action",
		Description: "AWS ElastiCache Update Action",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"cache_cluster_id", "replication_group_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"CacheClusterNotFound", "InvalidParameterValue"}),
			},
			Hydrate: listElastiCacheUpdateActions,
			Tags:    map[string]string{"service": "elasticache", "action": "DescribeCacheClusters"},
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheUpdateActions,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(elasticachev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cache_cluster_id",
				Description: "The ID of the cache cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_group_id",
				Description: "The ID of the replication group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine",
				Description: "The ElastiCache engine (Redis or Memcached).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "estimated_update_time",
				Description: "The estimated length of time for the update to complete.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "nodes_updated",
				Description: "The progress of the service update on the replication group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_update_name",
				Description: "The unique ID of the service update.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_update_recommended_apply_by_date",
				Description: "Recommended date to apply the update for compliance.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "service_update_release_date",
				Description: "The date the update was first available.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "service_update_severity",
				Description: "Severity level of the service update.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_update_status",
				Description: "Current status of the service update.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_update_type",
				Description: "Type of service update (security update, etc).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sla_met",
				Description: "Indicates if nodes were updated by recommended date.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_action_available_date",
				Description: "Date the update became available to the replication group.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_action_status",
				Description: "Current status of the update action.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_action_status_modified_date",
				Description: "Last modification date of the update action status.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(extracElasticacheUpdateActionId),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extracElasticacheUpdateActionId).Transform(arnToAkas),
			},
		}),
	}
}

// 核心数据获取逻辑
func listElastiCacheUpdateActions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_update_action.listElastiCacheUpdateActions", "client_error", err)
		return nil, err
	}
	input := &elasticache.DescribeUpdateActionsInput{}
	if v, ok := d.EqualsQuals["cache_cluster_id"]; ok {
		input.CacheClusterIds = []string{v.GetStringValue()}
	}
	if v, ok := d.EqualsQuals["replication_group_id"]; ok {
		input.ReplicationGroupIds = []string{v.GetStringValue()}
	}

	paginator := elasticache.NewDescribeUpdateActionsPaginator(client, input)
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_elasticache_update_action.listElastiCacheUpdateActions", "api_error", err)
			return nil, err
		}

		for _, action := range page.UpdateActions {
			d.StreamListItem(ctx, action)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func extracElasticacheUpdateActionId(ctx context.Context, data *transform.TransformData) (interface{}, error) {
	rs := data.HydrateItem.(types.UpdateAction)
	if rs.CacheClusterId != nil {
		return *rs.CacheClusterId, nil
	}
	return *rs.ReplicationGroupId, nil
}
