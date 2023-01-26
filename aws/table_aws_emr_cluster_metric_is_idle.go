package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/emr/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEmrClusterMetricIsIdle(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_emr_cluster_metric_is_idle",
		Description: "AWS EMR Cluster Cloudwatch Metrics - IsIdle",
		List: &plugin.ListConfig{
			ParentHydrate: listEmrClusters,
			Hydrate:       listEmrClusterMetricIsIdle,
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "id",
					Description: "The unique identifier for the cluster.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listEmrClusterMetricIsIdle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(types.ClusterSummary)
	return listCWMetricStatistics(ctx, d, "5_MIN", "AWS/ElasticMapReduce", "IsIdle", "JobFlowId", *data.Id)
}
