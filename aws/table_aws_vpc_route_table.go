package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsVpcRouteTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_route_table",
		Description: "AWS VPC Route table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("route_table_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorWithContext([]string{"InvalidRouteTableID.NotFound", "InvalidRouteTableID.Malformed"}),
			},
			Hydrate: getVpcRouteTable,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcRouteTables,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
			},
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listVpcRouteTables", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeRouteTablesInput{
		MaxResults: aws.Int64(100),
	}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "owner_id", FilterName: "owner-id", ColumnType: "string"},
		{ColumnName: "vpc_id", FilterName: "vpc-id", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
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
	err = svc.DescribeRouteTablesPages(
		input,
		func(page *ec2.DescribeRouteTablesOutput, isLast bool) bool {
			for _, routeTable := range page.RouteTables {
				d.StreamListItem(ctx, routeTable)

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

func getVpcRouteTable(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcRouteTable")

	region := d.KeyColumnQualString(matrixKeyRegion)
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
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":route-table/" + *routeTable.RouteTableId}

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
