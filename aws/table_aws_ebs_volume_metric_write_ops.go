package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEbsVolumeMetricWriteOps(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ebs_volume_metric_write_ops",
		Description: "AWS EBS Volume Cloudwatch Metrics - Write Ops",
		List: &plugin.ListConfig{
			ParentHydrate: listEBSVolume,
			Hydrate:       listEbsVolumeMetricWriteOps,
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
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

func listEbsVolumeMetricWriteOps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	volume := h.Item.(types.Volume)
	return listCWMetricStatistics(ctx, d, "5_MIN", "AWS/EBS", "VolumeWriteOps", "VolumeId", *volume.VolumeId)
}
