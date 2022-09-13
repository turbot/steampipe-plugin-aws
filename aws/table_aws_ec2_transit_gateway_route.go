package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
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

	input := &ec2.SearchTransitGatewayRoutesInput{
		MaxResults:                 aws.Int64(1000),
		TransitGatewayRouteTableId: routeTableId,
	}

	filterName := "state"
	blackholeState := "blackhole"
	activeState := "active"
	pendingState := "pending"
	filterValue := []*string{&blackholeState, &pendingState, &activeState}

	filters := buildEc2TransitGatewayRouteFilter(d.Quals)
	filters = append(filters, &ec2.Filter{Name: &filterName, Values: filterValue})

	input.Filters = filters

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 5 {
				input.MaxResults = aws.Int64(5)
			} else {
				input.MaxResults = limit
			}
		}
	}
	// List call
	// Filter parameter is required for making the api call otherwise it is throwing the error.
	res, err := svc.SearchTransitGatewayRoutes(input)
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

//// UTILITY FUNCTION
// Build ec2 transit gateway route list call input filter
func buildEc2TransitGatewayRouteFilter(quals plugin.KeyColumnQualMap) []*ec2.Filter {
	filters := make([]*ec2.Filter, 0)

	filterQuals := map[string]string{
		"prefix_list_id": "prefix-list-id",
		"state":          "state",
		"type":           "type",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := ec2.Filter{
				Name: aws.String(filterName),
			}
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []*string{aws.String(val)}
			} else {
				v := value.([]*string)
				filter.Values = v
			}
			filters = append(filters, &filter)
		}
	}
	return filters
}
