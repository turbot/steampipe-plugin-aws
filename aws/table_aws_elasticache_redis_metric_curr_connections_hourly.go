package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElasticacheRedisMetricCurrConnectionsHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_redis_metric_curr_connections_hourly",
		Description: "AWS Elasticache Redis CurrConnections metric (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listElastiCacheClusters,
			Hydrate:       listElastiCacheMetricCurrConnectionsHourly,
		},
		GetMatrixItemFunc: BuildRegionList,
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

func listElastiCacheMetricCurrConnectionsHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheClusterConfiguration := h.Item.(*elasticache.CacheCluster)
	return listCWMetricStatistics(ctx, d, "Hourly", "AWS/ElastiCache", "CurrConnections", "CacheClusterId", *cacheClusterConfiguration.CacheClusterId)
}
