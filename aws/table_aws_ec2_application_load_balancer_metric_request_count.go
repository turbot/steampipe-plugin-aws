package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION
func tableAwsEc2ApplicationLoadBalancerMetricRequestCount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_application_load_balancer_metric_request_count",
		Description: "AWS EC2 Application Load Balancer Metrics - Request Count",
		List: &plugin.ListConfig{
			ParentHydrate: listEc2ApplicationLoadBalancers,
			Hydrate:       listEc2ApplicationLoadBalancerMetricRequestCount,
		},
		GetMatrixItemFunc: BuildRegionList,
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

func listEc2ApplicationLoadBalancerMetricRequestCount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	loadBalancer := h.Item.(*elbv2.LoadBalancer)
	arn := strings.SplitN(*loadBalancer.LoadBalancerArn, "/", 2)[1]
	return listCWMetricStatistics(ctx, d, "5_MIN", "AWS/ApplicationELB", "RequestCount", "LoadBalancer", arn)
}
