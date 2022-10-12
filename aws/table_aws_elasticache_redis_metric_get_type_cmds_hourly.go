package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElasticacheRedisMetricGetTypeCmdsHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_redis_metric_get_type_cmds_hourly",
		Description: "AWS Elasticache Redis GetTypeCmds metric(Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listElastiCacheClusters,
			Hydrate:       listElastiCacheMetricGetTypeCmdsHourly,
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

func listElastiCacheMetricGetTypeCmdsHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheClusterConfiguration := h.Item.(types.CacheCluster)
	return listCWMetricStatistics(ctx, d, "Hourly", "AWS/ElastiCache", "GetTypeCmds", "CacheClusterId", *cacheClusterConfiguration.CacheClusterId)
}
