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

func tableAwsVpcVpnConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_vpn_connection",
		Description: "AWS VPC VPN Connection",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("vpn_connection_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidVpnConnectionID.NotFound", "InvalidVpnConnectionID.Malformed"}),
			Hydrate:           getVpcVpnConnection,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcVpnConnections,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "vpn_connection_id",
				Description: "The ID of the VPN connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "ARN of the connection.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpcVpnConnectionARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "state",
				Description: "The current state of the VPN connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of VPN connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "category",
				Description: "The category of the VPN connection. A value of VPN indicates an AWS VPN connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpn_gateway_id",
				Description: "The ID of the virtual private gateway at the AWS side of the VPN connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_gateway_id",
				Description: "The ID of the customer gateway at your end of the VPN connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_gateway_configuration",
				Description: "The configuration information for the VPN connection's customer gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "transit_gateway_id",
				Description: "The ID of the transit gateway associated with the VPN connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "options",
				Description: "The VPN connection options.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "routes",
				Description: "The static routes associated with the VPN connection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vgw_telemetry",
				Description: "Information about the VPN tunnel.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to VPN gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getVpcVpnConnectionTurbotData, "Title"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getVpcVpnConnectionTurbotData, "Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcVpnConnectionARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcVpnConnections(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listVpcVpnConnections", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	resp, err := svc.DescribeVpnConnections(&ec2.DescribeVpnConnectionsInput{})
	for _, vpnConnection := range resp.VpnConnections {
		d.StreamListItem(ctx, vpnConnection)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcVpnConnection(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpcVpnGateway")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	vpnConnectionID := d.KeyColumnQuals["vpn_connection_id"].GetStringValue()

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpnConnectionsInput{
		VpnConnectionIds: []*string{aws.String(vpnConnectionID)},
	}

	// Get call
	op, err := svc.DescribeVpnConnections(params)
	if err != nil {
		logger.Debug("getVpcVpnConnection", "ERROR", err)
		return nil, err
	}

	if op.VpnConnections != nil && len(op.VpnConnections) > 0 {
		return op.VpnConnections[0], nil
	}
	return nil, nil
}

func getVpcVpnConnectionARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcVpnConnectionARN")
	vpnConnection := h.Item.(*ec2.VpnConnection)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	arn := "arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":vpn-connection/" + *vpnConnection.VpnConnectionId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getVpcVpnConnectionTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vpnConnection := d.HydrateItem.(*ec2.VpnConnection)
	param := d.Param.(string)

	// Get resource title
	title := vpnConnection.VpnConnectionId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if vpnConnection.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range vpnConnection.Tags {
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
