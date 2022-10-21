package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEbsVolumeMetricReadOpsDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ebs_volume_metric_read_ops_daily",
		Description: "AWS EBS Volume Cloudwatch Metrics - Read Ops (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listEBSVolume,
			Hydrate:       listEbsVolumeMetricReadOpsDaily,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "volume_id",
					Description: "The EBS Volume ID.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listEbsVolumeMetricReadOpsDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	volume := h.Item.(types.Volume)
	return listCWMetricStatistics(ctx, d, "DAILY", "AWS/EBS", "VolumeReadOps", "VolumeId", *volume.VolumeId)
}
