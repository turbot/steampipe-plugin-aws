package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/redshift/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRedshiftClusterMetricCpuUtilizationDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshift_cluster_metric_cpu_utilization_daily",
		Description: "AWS Redshift Cluster Cloudwatch Metrics - CPU Utilization (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listRedshiftClusters,
			Hydrate:       listRedshiftClusterMetricCpuUtilizationDaily,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "cluster_identifier",
					Description: "The friendly name to identify the DB Instance.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listRedshiftClusterMetricCpuUtilizationDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(types.Cluster)
	return listCWMetricStatistics(ctx, d, "DAILY", "AWS/Redshift", "CPUUtilization", "ClusterIdentifier", *cluster.ClusterIdentifier)
}
