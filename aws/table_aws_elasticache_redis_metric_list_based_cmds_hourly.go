package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElasticacheRedisMetricListBasedCmdsHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_redis_metric_list_based_cmds_hourly",
		Description: "AWS Elasticache Redis ListBasedCmds metric (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listElastiCacheClusters,
			Hydrate:       listElastiCacheMetricListBasedCmdsHourly,
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "cache_cluster_id",
					Description: "The cache cluster id.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listElastiCacheMetricListBasedCmdsHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheClusterConfiguration := h.Item.(types.CacheCluster)
	return listCWMetricStatistics(ctx, d, "Hourly", "AWS/ElastiCache", "ListBasedCmds", "CacheClusterId", *cacheClusterConfiguration.CacheClusterId)
}
