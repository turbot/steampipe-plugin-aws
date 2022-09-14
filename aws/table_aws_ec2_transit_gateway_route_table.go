package aws

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2TransitGatewayRouteTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_transit_gateway_route_table",
		Description: "AWS EC2 Transit Gateway Route Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("transit_gateway_route_table_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidRouteTableID.NotFound", "InvalidRouteTableId.Unavailable", "InvalidRouteTableId.Malformed", "InvalidAction"}),
			},
			Hydrate: getEc2TransitGatewayRouteTable,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2TransitGatewayRouteTable,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "transit_gateway_id", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "default_association_route_table", Require: plugin.Optional},
				{Name: "default_propagation_route_table", Require: plugin.Optional},
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
				Name:        "transit_gateway_id",
				Description: "The ID of the transit gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the transit gateway route table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time of transit gateway route table.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "default_association_route_table",
				Description: "Indicates whether this is the default association route table for the transit gateway.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "default_propagation_route_table",
				Description: "Indicates whether this is the default propagation route table for the transit gateway.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getEc2TransitGatewayRouteTableTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getEc2TransitGatewayRouteTableTurbotTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2TransitGatewayRouteTableTurbotData,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2TransitGatewayRouteTable(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeTransitGatewayRouteTablesInput{
		MaxResults: aws.Int64(1000),
	}

	filters := []*ec2.Filter{}
	equalQuals := d.KeyColumnQuals
	if equalQuals["transit_gateway_id"] != nil {
		filters = append(filters, &ec2.Filter{Name: aws.String("transit-gateway-id"), Values: []*string{aws.String(equalQuals["transit_gateway_id"].GetStringValue())}})
	}
	if equalQuals["state"] != nil {
		filters = append(filters, &ec2.Filter{Name: aws.String("state"), Values: []*string{aws.String(equalQuals["state"].GetStringValue())}})
	}
	if equalQuals["default_association_route_table"] != nil {
		filters = append(filters, &ec2.Filter{Name: aws.String("default-association-route-table"), Values: []*string{aws.String(fmt.Sprint(equalQuals["default_association_route_table"].GetBoolValue()))}})
	}
	if equalQuals["default_propagation_route_table"] != nil {
		filters = append(filters, &ec2.Filter{Name: aws.String("default-propagation-route-table"), Values: []*string{aws.String(fmt.Sprint(equalQuals["default_propagation_route_table"].GetBoolValue()))}})
	}

	if len(filters) > 0 {
		input.Filters = filters
	}

	// Limiting the results
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
	err = svc.DescribeTransitGatewayRouteTablesPages(
		input,
		func(page *ec2.DescribeTransitGatewayRouteTablesOutput, isLast bool) bool {
			for _, transitGatewayRouteTable := range page.TransitGatewayRouteTables {
				d.StreamListItem(ctx, transitGatewayRouteTable)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2TransitGatewayRouteTable(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	routeTableID := d.KeyColumnQuals["transit_gateway_route_table_id"].GetStringValue()

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeTransitGatewayRouteTablesInput{
		TransitGatewayRouteTableIds: []*string{aws.String(routeTableID)},
	}

	op, err := svc.DescribeTransitGatewayRouteTables(params)
	if err != nil {
		return nil, err
	}

	if op.TransitGatewayRouteTables != nil && len(op.TransitGatewayRouteTables) > 0 {
		return op.TransitGatewayRouteTables[0], nil
	}
	return nil, nil
}

func getAwsEc2TransitGatewayRouteTableTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2TransitGatewayRouteTableTurbotData")
	region := d.KeyColumnQualString(matrixKeyRegion)
	transitGatewayRouteTable := h.Item.(*ec2.TransitGatewayRouteTable)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":transit-gateway-route-table/" + *transitGatewayRouteTable.TransitGatewayRouteTableId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getEc2TransitGatewayRouteTableTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ec2.TransitGatewayRouteTable)
	return ec2TagsToMap(data.Tags)
}

func getEc2TransitGatewayRouteTableTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ec2.TransitGatewayRouteTable)
	title := data.TransitGatewayRouteTableId
	if data.Tags != nil {
		for _, i := range data.Tags {
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}
	return title, nil
}
