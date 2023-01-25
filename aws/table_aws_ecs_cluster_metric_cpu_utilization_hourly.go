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

func tableAwsEcsClusterMetricCpuUtilizationHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_cluster_metric_cpu_utilization_hourly",
		Description: "AWS ECS Cluster Cloudwatch Metrics - CPU Utilization (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listEcsClusters,
			Hydrate:       listEcsClusterMetricCpuUtilizationHourly,
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
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

func listEcsClusterMetricCpuUtilizationHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(types.Cluster)
	clusterName := strings.Split(*data.ClusterArn, "/")[1]
	return listCWMetricStatistics(ctx, d, "HOURLY", "AWS/ECS", "CPUUtilization", "ClusterName", clusterName)
}
