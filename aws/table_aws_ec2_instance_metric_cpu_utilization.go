package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION
func tableAwsEc2InstanceMetricCpuUtilization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_instance_metric_cpu_utilization",
		Description: "AWS EC2 Instance Cloudwatch Metrics - CPU Utilization",
		List: &plugin.ListConfig{
			ParentHydrate: listEc2Instance,
			Hydrate:       listEc2InstanceMetricCpuUtilization,
		},
		GetMatrixItem: BuildRegionList,
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

func listEc2InstanceMetricCpuUtilization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(*ec2.Instance)
	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("InstanceId"),
			Value: instance.InstanceId,
		},
	}
	return listCWMetricStatisticsTest(ctx, d, "5_MIN", "AWS/EC2", "CPUUtilization", dimensions)
}
