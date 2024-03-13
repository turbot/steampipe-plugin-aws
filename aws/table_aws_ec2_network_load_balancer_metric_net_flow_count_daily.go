package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2NetworkLoadBalancerMetricNetFlowCountDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_network_load_balancer_metric_net_flow_count_daily",
		Description: "AWS EC2 Network Load Balancer Metrics - Net Flow Count (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listEc2NetworkLoadBalancers,
			Hydrate:       listEc2NetworkLoadBalancerMetricNetFlowCountDaily,
			Tags:          map[string]string{"service": "cloudwatch", "action": "GetMetricStatistics"},
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
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
	loadBalancer := h.Item.(types.LoadBalancer)
	arn := strings.SplitN(*loadBalancer.LoadBalancerArn, "/", 2)[1]
	return listCWMetricStatistics(ctx, d, "DAILY", "AWS/NetworkELB", "NewFlowCount", "LoadBalancer", arn)
}
