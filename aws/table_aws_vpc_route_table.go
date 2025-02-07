package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2Endpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsVpcRouteTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_route_table",
		Description: "AWS VPC Route table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("route_table_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidRouteTableID.NotFound", "InvalidRouteTableID.Malformed"}),
			},
			Hydrate: getVpcRouteTable,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeRouteTables"},
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcRouteTables,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeRouteTables"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2Endpoint.AWS_EC2_SERVICE_ID),
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

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_route_table.listVpcRouteTables", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				maxLimit = int32(5)
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeRouteTablesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "owner_id", FilterName: "owner-id", ColumnType: "string"},
		{ColumnName: "vpc_id", FilterName: "vpc-id", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeRouteTablesPaginator(svc, input, func(o *ec2.DescribeRouteTablesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_route_table.listVpcRouteTables", "api_error", err)
			return nil, err
		}

		for _, items := range output.RouteTables {
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

func getVpcRouteTable(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	routeTableID := d.EqualsQuals["route_table_id"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_route_table.getVpcRouteTable", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeRouteTablesInput{
		RouteTableIds: []string{routeTableID},
	}

	// Get call
	op, err := svc.DescribeRouteTables(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_route_table.getVpcRouteTable", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.RouteTables) > 0 {
		return op.RouteTables[0], nil
	}
	return nil, nil
}

func getVpcRouteTableAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	routeTable := h.Item.(types.RouteTable)
	region := d.EqualsQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_route_table.getVpcRouteTableAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":route-table/" + *routeTable.RouteTableId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcRouteTableTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	routeTable := d.HydrateItem.(types.RouteTable)
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
