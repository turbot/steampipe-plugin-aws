package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcsClusterMetricCpuUtilizationDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_cluster_metric_cpu_utilization_daily",
		Description: "AWS ECS Cluster Cloudwatch Metrics - CPU Utilization (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listEcsClusters,
			Hydrate:       listEcsClusterMetricCpuUtilizationDaily,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "cluster_name",
					Description: "A user-generated string that you use to identify your cluster.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listEcsClusterMetricCpuUtilizationDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(types.Cluster)
	clusterName := strings.Split(*data.ClusterArn, "/")[1]
	return listCWMetricStatistics(ctx, d, "DAILY", "AWS/ECS", "CPUUtilization", "ClusterName", clusterName)
}
