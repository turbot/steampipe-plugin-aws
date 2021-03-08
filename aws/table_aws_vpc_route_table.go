package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsVpcRouteTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_route_table",
		Description: "AWS VPC Route table",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("route_table_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidRouteTableID.NotFound", "InvalidRouteTableID.Malformed"}),
			Hydrate:           getVpcRouteTable,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcRouteTables,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "route_table_id",
				Description: "Contains the ID of the route table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account that owns the route table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associations",
				Description: "Contains the associations between the route table and one or more subnets or a gateway.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "routes",
				Description: "A list of routes in the route table.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "propagating_vgws",
				Description: "A list of virtual private gateway (VGW) propagating routes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the route table.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getVpcRouteTableTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RouteTableId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcRouteTableAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcRouteTables(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listVpcRouteTables", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeRouteTablesPages(
		&ec2.DescribeRouteTablesInput{},
		func(page *ec2.DescribeRouteTablesOutput, isLast bool) bool {
			for _, routeTable := range page.RouteTables {
				d.StreamListItem(ctx, routeTable)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcRouteTable(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcRouteTable")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	routeTableID := d.KeyColumnQuals["route_table_id"].GetStringValue()

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{aws.String(routeTableID)},
	}

	// Get call
	op, err := svc.DescribeRouteTables(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getVpcRouteTable__", "ERROR", err)
		return nil, err
	}

	if op.RouteTables != nil && len(op.RouteTables) > 0 {
		return op.RouteTables[0], nil
	}
	return nil, nil
}

func getVpcRouteTableAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcRouteTableTurbotAkas")
	routeTable := h.Item.(*ec2.RouteTable)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":route-table/" + *routeTable.RouteTableId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcRouteTableTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	routeTable := d.HydrateItem.(*ec2.RouteTable)
	var turbotTagsMap map[string]string

	// Get the resource tags
	if routeTable.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range routeTable.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
