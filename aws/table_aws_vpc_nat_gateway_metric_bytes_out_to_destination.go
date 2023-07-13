package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsVpcNatGatewayMetricBytesOutToDestination(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_nat_gateway_metric_bytes_out_to_destination",
		Description: "AWS VPC Nat Gateway Cloudwatch Metrics - BytesOutToDestination",
		List: &plugin.ListConfig{
			ParentHydrate: listVpcNatGateways,
			Hydrate:       listVpcNatGatewayMetricBytesOutToDestination,
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "nat_gateway_id",
					Description: "The ID of the NAT gateway.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listVpcNatGatewayMetricBytesOutToDestination(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	natGateway := h.Item.(types.NatGateway)
	return listCWMetricStatistics(ctx, d, "5_MIN", "AWS/NATGateway", "BytesOutToDestination", "NatGatewayId", *natGateway.NatGatewayId)
}
