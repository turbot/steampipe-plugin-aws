package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2TransitGatewayRoute(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_transit_gateway_route",
		Description: "AWS EC2 Transit Gateway Route",
		List: &plugin.ListConfig{
			ParentHydrate: listEc2TransitGatewayRouteTable,
			Hydrate:       listEc2TransitGatewayRoute,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "prefix_list_id", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidAction"}),
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "transit_gateway_route_table_id",
				Description: "The ID of the transit gateway route table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "destination_cidr_block",
				Description: "The CIDR block used for destination matches.",
				Type:        proto.ColumnType_CIDR,
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
	Route                      types.TransitGatewayRoute
	TransitGatewayRouteTableId string
}

//// LIST FUNCTION

func listEc2TransitGatewayRoute(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	routeTableId := h.Item.(types.TransitGatewayRouteTable).TransitGatewayRouteTableId

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_transit_gateway_route.listEc2TransitGatewayRoute", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.SearchTransitGatewayRoutesInput{
		MaxResults:                 aws.Int32(maxLimit),
		TransitGatewayRouteTableId: routeTableId,
	}

	filterName := "state"
	blackholeState := "blackhole"
	activeState := "active"
	pendingState := "pending"
	filterValue := []string{blackholeState, pendingState, activeState}

	filters := buildEc2TransitGatewayRouteFilter(d.Quals)
	filters = append(filters, types.Filter{Name: &filterName, Values: filterValue})

	input.Filters = filters

	// List call
	// Filter parameter is required for making the api call otherwise it is throwing the error.
	res, err := svc.SearchTransitGatewayRoutes(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, route := range res.Routes {
		d.StreamListItem(ctx, &RouteDetails{
			Route:                      route,
			TransitGatewayRouteTableId: *routeTableId,
		})

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// TRANSFORM FUNCTIONS

func getAwsEc2TransitGatewayRouteAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

// // UTILITY FUNCTION
// Build ec2 transit gateway route list call input filter
func buildEc2TransitGatewayRouteFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"prefix_list_id": "prefix-list-id",
		"state":          "state",
		"type":           "type",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []string{val}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
