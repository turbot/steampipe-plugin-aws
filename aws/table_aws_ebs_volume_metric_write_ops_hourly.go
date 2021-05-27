package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION
func tableAwsEbsVolumeMetricWriteOpsHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ebs_volume_metric_read_ops_hourly",
		Description: "AWS EBS Volume Cloudwatch Metrics - Write Ops (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listEBSVolume,
			Hydrate:       listEbsVolumeMetricWriteOpsHourly,
		},
		GetMatrixItem: BuildRegionList,
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

func listEbsVolumeMetricWriteOpsHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	volume := h.Item.(*ec2.Volume)
	return listCWMetricStatistics(ctx, d, "HOURLY", "AWS/EBS", "VolumeWriteOps", "VolumeId", *volume.VolumeId)
}
