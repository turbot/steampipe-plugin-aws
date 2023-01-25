package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsVpcVpnGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_vpn_gateway",
		Description: "AWS VPC VPN Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("vpn_gateway_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidVpnGatewayID.NotFound", "InvalidVpnGatewayID.Malformed"}),
			},
			Hydrate: getVpcVpnGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcVpnGateways,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "amazon_side_asn", Require: plugin.Optional},
				{Name: "availability_zone", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "vpn_gateway_id",
				Description: "The ID of the virtual private gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the virtual private gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of VPN connection the virtual private gateway supports.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "amazon_side_asn",
				Description: "The private Autonomous System Number (ASN) for the Amazon side of a BGP session.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "availability_zone",
				Description: "The Availability Zone where the virtual private gateway was created, if applicable. This field may be empty or not returned.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_attachments",
				Description: "Any VPCs attached to the virtual private gateway.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to VPN gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getVpcVpnGatewayTurbotData, "Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getVpcVpnGatewayTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcVpnGatewayTurbotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcVpnGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_vpn_gateway.listVpcVpnGateways", "connection_error", err)
		return nil, err
	}

	input := &ec2.DescribeVpnGatewaysInput{}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "amazon_side_asn", FilterName: "amazon-side-asn", ColumnType: "int64"},
		{ColumnName: "availability_zone", FilterName: "availability-zone", ColumnType: "string"},
		{ColumnName: "state", FilterName: "state", ColumnType: "string"},
		{ColumnName: "type", FilterName: "type", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// List call
	resp, err := svc.DescribeVpnGateways(ctx, input)

	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_vpn_gateway.listVpcVpnGateways", "api_error", err)
	}

	for _, vpnGateway := range resp.VpnGateways {
		d.StreamListItem(ctx, vpnGateway)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcVpnGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	vpnGatewayID := d.KeyColumnQuals["vpn_gateway_id"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_vpn_gateway.getVpcVpnGateway", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpnGatewaysInput{
		VpnGatewayIds: []string{vpnGatewayID},
	}

	// Get call
	op, err := svc.DescribeVpnGateways(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_vpn_gateway.getVpcVpnGateway", "api_error", err)
		return nil, err
	}

	if op.VpnGateways != nil && len(op.VpnGateways) > 0 {
		return op.VpnGateways[0], nil
	}
	return nil, nil
}

func getVpcVpnGatewayTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vpnGateway := h.Item.(types.VpnGateway)
	region := d.KeyColumnQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_vpn_gateway.getVpcVpnGatewayTurbotAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":vpn-gateway/" + *vpnGateway.VpnGatewayId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcVpnGatewayTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vpnGateway := d.HydrateItem.(types.VpnGateway)
	param := d.Param.(string)

	// Get resource title
	title := vpnGateway.VpnGatewayId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if vpnGateway.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range vpnGateway.Tags {
			turbotTagsMap[*i.Key] = *i.Value
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}

	if param == "Tags" {
		return turbotTagsMap, nil
	}

	return title, nil
}
