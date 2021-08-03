package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEmrClusterMetricIsidle(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_emr_cluster_metric_isidle",
		Description: "AWS EMR Cluster Cloudwatch Metrics - Isidle",
		List: &plugin.ListConfig{
			ParentHydrate: listEmrClusters,
			Hydrate:       listEmrClusterMetricIsidle,
		},
		GetMatrixItem: BuildRegionList,
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

func listEmrClusterMetricIsidle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*emr.ClusterSummary)
	return listCWMetricStatistics(ctx, d, "5_MIN", "AWS/ElasticMapReduce", "IsIdle", "JobFlowId", *data.Id)
}
