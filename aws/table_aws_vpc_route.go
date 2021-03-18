package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/ec2"
)

//// TABLE DEFINITION

func tableAwsVpcRoute(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_route",
		Description: "AWS VPC Route",
		// TODO -- get call returning a list of items

		// Get: &plugin.GetConfig{
		// 	KeyColumns:        plugin.SingleColumn("route_table_id"),
		// 	ShouldIgnoreError: isNotFoundError([]string{"InvalidRouteTableID.NotFound", "InvalidRouteTableID.Malformed"}),
		// 	Hydrate:           getAwsVpcRoute,
		// },
		List: &plugin.ListConfig{
			ParentHydrate: listVpcRouteTables,
			Hydrate:       listAwsVpcRoute,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "route_table_id",
				Description: "The ID of the route table conatining the route.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RouteTableID"),
			},
			{
				Name:        "state",
				Description: "The state of the route. The blackhole state indicates that the route's target isn't available (for example, the specified gateway isn't attached to the VPC, or the specified NAT instance has been terminated).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.State"),
			},
			{
				Name:        "destination_cidr_block",
				Description: "The IPv4 CIDR block used for the destination match.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("Route.DestinationCidrBlock"),
			},
			{
				Name:        "destination_ipv6_cidr_block",
				Description: "The IPv6 CIDR block used for the destination match.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("Route.DestinationIpv6CidrBlock"),
			},
			{
				Name:        "carrier_gateway_id",
				Description: "The ID of the carrier gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.CarrierGatewayId"),
			},
			{
				Name:        "destination_prefix_list_id",
				Description: "The prefix of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.DestinationPrefixListId"),
			},
			{
				Name:        "egress_only_internet_gateway_id",
				Description: "The ID of the egress-only internet gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.EgressOnlyInternetGatewayId"),
			},
			{
				Name:        "gateway_id",
				Description: "The ID of a gateway attached to your VPC.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.GatewayId"),
			},
			{
				Name:        "instance_id",
				Description: "The ID of a NAT instance in your VPC.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.InstanceId"),
			},
			{
				Name:        "instance_owner_id",
				Description: "The AWS account ID of the owner of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.InstanceOwnerId"),
			},
			{
				Name:        "local_gateway_id",
				Description: "The ID of the local gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.LocalGatewayId"),
			},
			{
				Name:        "nat_gateway_id",
				Description: "The ID of a NAT gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.NatGatewayId"),
			},
			{
				Name:        "network_interface_id",
				Description: "The ID of the network interface.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.NetworkInterfaceId"),
			},
			{
				Name:        "transit_gateway_id",
				Description: "The ID of a transit gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.TransitGatewayId"),
			},
			{
				Name:        "vpc_peering_connection_id",
				Description: "The ID of a VPC peering connection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.VpcPeeringConnectionId"),
			},
			{
				Name:        "origin",
				Description: "Describes how the route was created. CreateRouteTable - The route was automatically created when the route table was created. CreateRoute - The route was manually added to the route table. EnableVgwRoutePropagation - The route was propagated by route propagation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.Origin"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsVpcRouteTurbotData,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsVpcRouteTurbotData,
			},
		}),
	}
}

type routeTableRoute = struct {
	RouteTableID *string
	Route        *ec2.Route
}

//// LIST FUNCTION

func listAwsVpcRoute(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsVpcRoute")

	routeTable := h.Item.(*ec2.RouteTable)

	for _, route := range routeTable.Routes {
		d.StreamLeafListItem(ctx, &routeTableRoute{routeTable.RouteTableId, route})
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsVpcRouteTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsVpcRouteTurbotData")
	routeData := h.Item.(*routeTableRoute)
	commonColumnData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonData := commonColumnData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	var title string
	var akas []string
	if routeData.Route.DestinationCidrBlock != nil {
		title = *routeData.RouteTableID + "_" + *routeData.Route.DestinationCidrBlock
		akas = []string{"arn:" + commonData.Partition + ":ec2:" + commonData.Region + ":" + commonData.AccountId + ":route-table/" + *routeData.RouteTableID + ":" + *routeData.Route.DestinationCidrBlock}
	} else {
		title = *routeData.RouteTableID + "_" + *routeData.Route.DestinationIpv6CidrBlock
		akas = []string{"arn:" + commonData.Partition + ":ec2:" + commonData.Region + ":" + commonData.AccountId + ":route-table/" + *routeData.RouteTableID + ":" + *routeData.Route.DestinationIpv6CidrBlock}
	}

	// Mapping all turbot defined properties
	turbotData := map[string]interface{}{
		"Akas":  akas,
		"Title": title,
	}

	return turbotData, nil
}
