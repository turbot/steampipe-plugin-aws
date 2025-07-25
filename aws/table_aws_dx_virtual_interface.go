package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	"github.com/aws/aws-sdk-go-v2/service/directconnect/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsDxVirtualInterface(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dx_virtual_interface",
		Description: "AWS Direct Connect Virtual Interface",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("virtual_interface_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DirectConnectClientException"}),
			},
			Hydrate: getDxVirtualInterface,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeVirtualInterfaces"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDxVirtualInterfaces,
			KeyColumns: plugin.OptionalColumns([]string{"connection_id"}),
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeVirtualInterfaces"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DIRECTCONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "virtual_interface_id",
				Description: "The ID of the virtual interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_interface_name",
				Description: "The name of the virtual interface assigned by the customer network.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_interface_type",
				Description: "The type of virtual interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_interface_state",
				Description: "The state of the virtual interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connection_id",
				Description: "The ID of the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vlan",
				Description: "The ID of the VLAN.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "location",
				Description: "The location of the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bgp_asn",
				Description: "The autonomous system (AS) number for Border Gateway Protocol (BGP) configuration.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Asn"),
			},
			{
				Name:        "amazon_address",
				Description: "The IP address assigned to the Amazon interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_address",
				Description: "The IP address assigned to the customer interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "address_family",
				Description: "The address family for the BGP peer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "amazon_side_asn",
				Description: "The autonomous system number (ASN) for the Amazon side of the connection.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "auth_key",
				Description: "The authentication key for BGP configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "aws_device_v2",
				Description: "The Direct Connect endpoint that terminates the physical connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "aws_logical_device_id",
				Description: "The Direct Connect endpoint that terminates the logical connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_router_config",
				Description: "The customer router configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mtu",
				Description: "The maximum transmission unit (MTU), in bytes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "jumbo_frame_capable",
				Description: "Indicates whether jumbo frames (9001 MTU) are supported.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "owner_account",
				Description: "The ID of the AWS account that owns the virtual interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "direct_connect_gateway_id",
				Description: "The ID of the Direct Connect gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_gateway_id",
				Description: "The ID of the virtual private gateway. Applies only to private virtual interfaces.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "site_link_enabled",
				Description: "Indicates whether SiteLink is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "route_filter_prefixes",
				Description: "The routes to be filtered for this virtual interface.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "bgp_peers",
				Description: "The BGP peers configured on this virtual interface.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "router_config_interface_id",
				Description: "The virtual interface ID that can be used with DescribeRouterConfiguration API.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDxVirtualInterfaceRouterConfig,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the virtual interface.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualInterfaceName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getDxVirtualInterfaceTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDxVirtualInterfaceARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDxVirtualInterfaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_virtual_interface.listDxVirtualInterfaces", "connection_error", err)
		return nil, err
	}

	input := &directconnect.DescribeVirtualInterfacesInput{}

	if d.EqualsQualString("connection_id") != "" {
		input.ConnectionId = aws.String(d.EqualsQualString("connection_id"))
	}

	// Execute list call
	// apply rate limiting
	d.WaitForListRateLimit(ctx)

	output, err := svc.DescribeVirtualInterfaces(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_virtual_interface.listDxVirtualInterfaces", "api_error", err)
		return nil, err
	}

	for _, item := range output.VirtualInterfaces {
		d.StreamListItem(ctx, item)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDxVirtualInterface(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	virtualInterfaceID := d.EqualsQuals["virtual_interface_id"].GetStringValue()

	// get service
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_virtual_interface.getDxVirtualInterface", "connection_error", err)
		return nil, err
	}

	params := &directconnect.DescribeVirtualInterfacesInput{
		VirtualInterfaceId: aws.String(virtualInterfaceID),
	}

	op, err := svc.DescribeVirtualInterfaces(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_virtual_interface.getDxVirtualInterface", "api_error", err)
		return nil, err
	}

	if len(op.VirtualInterfaces) > 0 {
		return op.VirtualInterfaces[0], nil
	}
	return nil, nil
}

func getDxVirtualInterfaceARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vif := h.Item.(types.VirtualInterface)
	region := d.EqualsQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":directconnect:" + region + ":" + *vif.OwnerAccount + ":dxvif/" + *vif.VirtualInterfaceId

	return arn, nil
}

func getDxVirtualInterfaceRouterConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vif := h.Item.(types.VirtualInterface)

	// Only attempt to get router config if we have a virtual interface ID
	if vif.VirtualInterfaceId == nil {
		return nil, nil
	}

	// get service
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_virtual_interface.getDxVirtualInterfaceRouterConfig", "connection_error", err)
		return nil, err
	}

	params := &directconnect.DescribeRouterConfigurationInput{
		VirtualInterfaceId: vif.VirtualInterfaceId,
	}

	op, err := svc.DescribeRouterConfiguration(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_virtual_interface.getDxVirtualInterfaceRouterConfig", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getDxVirtualInterfaceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vif := d.HydrateItem.(types.VirtualInterface)

	if vif.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range vif.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
