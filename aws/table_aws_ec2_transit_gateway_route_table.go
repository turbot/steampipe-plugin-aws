package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2TransitGatewayRouteTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_transit_gateway_route_table",
		Description: "AWS EC2 Transit Gateway Route Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("transit_gateway_route_table_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidRouteTableID.NotFound", "InvalidRouteTableId.Unavailable", "InvalidRouteTableId.Malformed"}),
			},
			Hydrate: getEc2TransitGatewayRouteTable,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeTransitGatewayRouteTables"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2TransitGatewayRouteTable,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeTransitGatewayRouteTables"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "transit_gateway_id", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "default_association_route_table", Require: plugin.Optional},
				{Name: "default_propagation_route_table", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAction"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
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

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_transit_gateway_route_table.listEc2TransitGatewayRouteTable", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				maxLimit = 5
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeTransitGatewayRouteTablesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filters := []types.Filter{}
	equalQuals := d.EqualsQuals
	if equalQuals["transit_gateway_id"] != nil {
		filters = append(filters, types.Filter{Name: aws.String("transit-gateway-id"), Values: []string{equalQuals["transit_gateway_id"].GetStringValue()}})
	}
	if equalQuals["state"] != nil {
		filters = append(filters, types.Filter{Name: aws.String("state"), Values: []string{equalQuals["state"].GetStringValue()}})
	}
	if equalQuals["default_association_route_table"] != nil {
		filters = append(filters, types.Filter{Name: aws.String("default-association-route-table"), Values: []string{fmt.Sprint(equalQuals["default_association_route_table"].GetBoolValue())}})
	}
	if equalQuals["default_propagation_route_table"] != nil {
		filters = append(filters, types.Filter{Name: aws.String("default-propagation-route-table"), Values: []string{fmt.Sprint(equalQuals["default_propagation_route_table"].GetBoolValue())}})
	}

	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeTransitGatewayRouteTablesPaginator(svc, input, func(o *ec2.DescribeTransitGatewayRouteTablesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_transit_gateway_route_table.listEc2TransitGatewayRouteTable", "api_error", err)
			return nil, err
		}

		for _, items := range output.TransitGatewayRouteTables {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2TransitGatewayRouteTable(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	routeTableID := d.EqualsQuals["transit_gateway_route_table_id"].GetStringValue()

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_transit_gateway_route_table.getEc2TransitGatewayRouteTable", "api_error", err)
		return nil, err
	}

	params := &ec2.DescribeTransitGatewayRouteTablesInput{
		TransitGatewayRouteTableIds: []string{routeTableID},
	}

	op, err := svc.DescribeTransitGatewayRouteTables(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_transit_gateway_route_table.getEc2TransitGatewayRouteTable", "api_error", err)
		return nil, err
	}

	if op.TransitGatewayRouteTables != nil && len(op.TransitGatewayRouteTables) > 0 {
		return op.TransitGatewayRouteTables[0], nil
	}
	return nil, nil
}

func getAwsEc2TransitGatewayRouteTableTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	transitGatewayRouteTable := h.Item.(types.TransitGatewayRouteTable)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_transit_gateway_route_table.getAwsEc2TransitGatewayRouteTableTurbotData", "api_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":transit-gateway-route-table/" + *transitGatewayRouteTable.TransitGatewayRouteTableId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getEc2TransitGatewayRouteTableTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.TransitGatewayRouteTable)
	var turbotTagsMap map[string]string
	if data.Tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range data.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

func getEc2TransitGatewayRouteTableTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.TransitGatewayRouteTable)
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
