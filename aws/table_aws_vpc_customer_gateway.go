package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsVpcCustomerGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_customer_gateway",
		Description: "AWS VPC Customer Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("customer_gateway_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidCustomerGatewayID.NotFound", "InvalidCustomerGatewayID.Malformed"}),
			},
			Hydrate: getVpcCustomerGateway,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeCustomerGateways"},
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcCustomerGateways,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeCustomerGateways"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "ip_address", Require: plugin.Optional},
				{Name: "bgp_asn", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "customer_gateway_id",
				Description: "The ID of the customer gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of VPN connection the customer gateway supports (ipsec.1).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the customer gateway (pending | available | deleting | deleted).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bgp_asn",
				Description: "The customer gateway's Border Gateway Protocol (BGP) Autonomous System Number (ASN).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_arn",
				Description: "The Amazon Resource Name (ARN) for the customer gateway certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "device_name",
				Description: "The name of customer gateway device.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address",
				Description: "The Internet-routable IP address of the customer gateway's outside interface.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to customer gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getVpcCustomerGatewayTurbotData, "Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getVpcCustomerGatewayTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcCustomerGatewayTurbotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcCustomerGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_customer_gateway.listVpcCustomerGateways", "connection error", err)
		return nil, err
	}

	input := &ec2.DescribeCustomerGatewaysInput{}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "ip_address", FilterName: "ip-address", ColumnType: "ipaddr"},
		{ColumnName: "bgp_asn", FilterName: "bgp-asn", ColumnType: "string"},
		{ColumnName: "state", FilterName: "state", ColumnType: "string"},
		{ColumnName: "type", FilterName: "type", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// apply rate limiting
	d.WaitForListRateLimit(ctx)

	// List call
	resp, err := svc.DescribeCustomerGateways(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_customer_gateway.listVpcCustomerGateways", "api_error", err)
	}
	for _, customerGateway := range resp.CustomerGateways {
		d.StreamListItem(ctx, customerGateway)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcCustomerGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	customerGatewayID := d.EqualsQuals["customer_gateway_id"].GetStringValue()

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_customer_gateway.getVpcCustomerGateway", "connection error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeCustomerGatewaysInput{
		CustomerGatewayIds: []string{customerGatewayID},
	}

	// Get call
	op, err := svc.DescribeCustomerGateways(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_customer_gateway.getVpcCustomerGateway", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.CustomerGateways) > 0 {
		return op.CustomerGateways[0], nil
	}
	return nil, nil
}

func getVpcCustomerGatewayTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	customerGateway := h.Item.(types.CustomerGateway)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_customer_gateway.getVpcCustomerGatewayTurbotAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":customer-gateway/" + *customerGateway.CustomerGatewayId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcCustomerGatewayTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	customerGateway := d.HydrateItem.(types.CustomerGateway)
	param := d.Param.(string)

	// Get resource title
	title := customerGateway.CustomerGatewayId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if customerGateway.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range customerGateway.Tags {
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
