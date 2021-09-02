package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2TransitGatewayRoute(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_transit_gateway_route",
		Description: "AWS EC2 Transit Gateway Route",
		List: &plugin.ListConfig{
			ParentHydrate: listEc2TransitGatewayRouteTable,
			Hydrate:       listEc2TransitGatewayRoute,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "transit_gateway_route_table_id",
				Description: "The ID of the transit gateway route table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "destination_cidr_block",
				Description: "The CIDR block used for destination matches.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.DestinationCidrBlock"),
			},
			{
				Name:        "prefix_list_id",
				Description: "The ID of the prefix list used for destination matches.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.PrefixListId"),
			},
			{
				Name:        "state",
				Description: "The state of the route.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.State"),
			},
			{
				Name:        "type",
				Description: "The route type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.Type"),
			},
			{
				Name:        "transit_gateway_attachments",
				Description: "The attachments.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Route.TransitGatewayAttachments"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Route.DestinationCidrBlock"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2TransitGatewayRouteAka,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type RouteDetails struct {
	Route                      *ec2.TransitGatewayRoute
	TransitGatewayRouteTableId string
}

//// LIST FUNCTION

func listEc2TransitGatewayRoute(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)

	routeTableId := h.Item.(*ec2.TransitGatewayRouteTable).TransitGatewayRouteTableId

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	maxResult := int64(1000)
	filterName := "state"
	blackholeState := "blackhole"
	activeState := "active"
	pendingState := "pending"
	filterValue := []*string{&blackholeState, &pendingState, &activeState}

	// List call
	// Filter parameter is required for making the api call otherwise it is throwing the error.
	res, err := svc.SearchTransitGatewayRoutes(&ec2.SearchTransitGatewayRoutesInput{
		TransitGatewayRouteTableId: routeTableId,
		MaxResults:                 &maxResult,
		Filters:                    []*ec2.Filter{{Name: &filterName, Values: filterValue}},
	})
	if err != nil {
		return nil, err
	}

	for _, route := range res.Routes {
		d.StreamListItem(ctx, &RouteDetails{
			Route:                      route,
			TransitGatewayRouteTableId: *routeTableId,
		})
	}

	return nil, err
}

//// TRANSFORM FUNCTIONS

func getAwsEc2TransitGatewayRouteAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2TransitGatewayRouteAka")
	region := d.KeyColumnQualString(matrixKeyRegion)
	route := h.Item.(*RouteDetails)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":transit-gateway-route-table/" + route.TransitGatewayRouteTableId + ":" + *route.Route.DestinationCidrBlock}

	return akas, nil
}
