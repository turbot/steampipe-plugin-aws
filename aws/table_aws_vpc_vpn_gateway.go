package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcVpnGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_vpn_gateway",
		Description: "AWS VPC VPN Gateway",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("vpn_gateway_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidVpnGatewayID.NotFound", "InvalidVpnGatewayID.Malformed"}),
			ItemFromKey:       vpnGatewayFromKey,
			Hydrate:           getVpcVpnGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcVpnGateways,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "vpn_gateway_id",
				Description: "The ID of the virtual private gateway",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the virtual private gateway",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of VPN connection the virtual private gateway supports",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "amazon_side_asn",
				Description: "The private Autonomous System Number (ASN) for the Amazon side of a BGP session",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "availability_zone",
				Description: "The Availability Zone where the virtual private gateway was created, if applicable. This field may be empty or not returned",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_attachments",
				Description: "Any VPCs attached to the virtual private gateway",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to VPN gateway",
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

//// ITEM FROM KEY

func vpnGatewayFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	vpnGatewayID := quals["vpn_gateway_id"].GetStringValue()
	item := &ec2.VpnGateway{
		VpnGatewayId: &vpnGatewayID,
	}
	return item, nil
}

//// LIST FUNCTION

func listVpcVpnGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listVpcVpnGateways", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := Ec2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	resp, err := svc.DescribeVpnGateways(&ec2.DescribeVpnGatewaysInput{})
	for _, vpnGateway := range resp.VpnGateways {
		d.StreamListItem(ctx, vpnGateway)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcVpnGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpcVpnGateway")
	vpnGateway := h.Item.(*ec2.VpnGateway)
	defaultRegion := GetDefaultRegion()

	// get service
	svc, err := Ec2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpnGatewaysInput{
		VpnGatewayIds: []*string{vpnGateway.VpnGatewayId},
	}

	// Get call
	op, err := svc.DescribeVpnGateways(params)
	if err != nil {
		logger.Debug("getVpcVpnGateway__", "ERROR", err)
		return nil, err
	}

	if op.VpnGateways != nil && len(op.VpnGateways) > 0 {
		return op.VpnGateways[0], nil
	}
	return nil, nil
}

func getVpcVpnGatewayTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcVpnGatewayTurbotAkas")
	vpnGateway := h.Item.(*ec2.VpnGateway)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":vpn-gateway/" + *vpnGateway.VpnGatewayId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcVpnGatewayTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vpnGateway := d.HydrateItem.(*ec2.VpnGateway)
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
