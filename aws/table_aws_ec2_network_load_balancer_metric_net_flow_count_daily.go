package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION
func tableAwsEc2NetworkLoadBalancerMetricNetFlowCountDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_network_load_balancer_metric_net_flow_count_daily",
		Description: "AWS EC2 Network Load Balancer Metrics - Net Flow Count (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listEc2NetworkLoadBalancers,
			Hydrate:       listEc2NetworkLoadBalancerMetricNetFlowCountDaily,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "name",
					Description: "The friendly name of the Load Balancer.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listEc2NetworkLoadBalancerMetricNetFlowCountDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	loadBalancer := h.Item.(*elbv2.LoadBalancer)
	arn := strings.SplitN(*loadBalancer.LoadBalancerArn, "/", 2)[1]
	return listCWMetricStatistics(ctx, d, "DAILY", "AWS/NetworkELB", "NewFlowCount", "LoadBalancer", arn)
}
