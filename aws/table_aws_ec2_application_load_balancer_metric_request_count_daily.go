package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2ApplicationLoadBalancerMetricRequestCountDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_application_load_balancer_metric_request_count_daily",
		Description: "AWS EC2 Application Load Balancer Metrics - Request Count (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listEc2ApplicationLoadBalancers,
			Hydrate:       listEc2ApplicationLoadBalancerMetricRequestCountDaily,
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "name",
					Description: "The friendly name of the Load Balancer that was provided during resource creation.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listEc2ApplicationLoadBalancerMetricRequestCountDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	loadBalancer := h.Item.(types.LoadBalancer)
	arn := strings.SplitN(*loadBalancer.LoadBalancerArn, "/", 2)[1]
	return listCWMetricStatistics(ctx, d, "DAILY", "AWS/ApplicationELB", "RequestCount", "LoadBalancer", arn)
}
