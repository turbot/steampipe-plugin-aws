package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElasticacheRedisMetricCurrConnectionsHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_redis_metric_curr_connections_hourly",
		Description: "AWS Elasticache Redis CurrConnections metric (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listElastiCacheClusters,
			Hydrate:       listElastiCacheMetricCurrConnectionsHourly,
			Tags:          map[string]string{"service": "cloudwatch", "action": "GetMetricStatistics"},
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

func listElastiCacheMetricCurrConnectionsHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheClusterConfiguration := h.Item.(types.CacheCluster)
	return listCWMetricStatistics(ctx, d, "Hourly", "AWS/ElastiCache", "CurrConnections", "CacheClusterId", *cacheClusterConfiguration.CacheClusterId)
}
