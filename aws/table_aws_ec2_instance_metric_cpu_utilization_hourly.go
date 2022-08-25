package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION
func tableAwsEc2InstanceMetricCpuUtilizationHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_instance_metric_cpu_utilization_hourly",
		Description: "AWS EC2 Instance Cloudwatch Metrics - CPU Utilization (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listEc2Instance,
			Hydrate:       listEc2InstanceMetricCpuUtilizationHourly,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "instance_id",
					Description: "The ID of the instance.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listEc2InstanceMetricCpuUtilizationHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(*ec2.Instance)
	return listCWMetricStatistics(ctx, d, "HOURLY", "AWS/EC2", "CPUUtilization", "InstanceId", *instance.InstanceId)
}
