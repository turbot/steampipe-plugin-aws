package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElasticacheRedisEngineCPUUtilizationDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_redis_metric_engine_cpu_utilization_daily",
		Description: "AWS Elasticache Redis EngineCPUUtilization metric (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listElastiCacheClusters,
			Hydrate:       listElastiCacheMetricEngineCPUUtilizationDaily,
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

func listElastiCacheMetricEngineCPUUtilizationDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheClusterConfiguration := h.Item.(types.CacheCluster)
	return listCWMetricStatistics(ctx, d, "DAILY", "AWS/ElastiCache", "EngineCPUUtilization", "CacheClusterId", *cacheClusterConfiguration.CacheClusterId)
}
