package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcCustomerGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_customer_gateway",
		Description: "AWS VPC Customer Gateway",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("customer_gateway_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidCustomerGatewayID.NotFound", "InvalidCustomerGatewayID.Malformed"}),
			ItemFromKey:       customerGatewayFromKey,
			Hydrate:           getVpcCustomerGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcCustomerGateways,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "customer_gateway_id",
				Description: "The ID of the customer gateway",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of VPN connection the customer gateway supports (ipsec.1)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the customer gateway (pending | available | deleting | deleted)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bgp_asn",
				Description: "The customer gateway's Border Gateway Protocol (BGP) Autonomous System Number (ASN)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_arn",
				Description: "The Amazon Resource Name (ARN) for the customer gateway certificate",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "device_name",
				Description: "The name of customer gateway device",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address",
				Description: "The Internet-routable IP address of the customer gateway's outside interface",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "tags_raw",
				Description: "A list of tags that are attached to customer gateway",
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

//// ITEM FROM KEY

func customerGatewayFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	customerGatewayID := quals["customer_gateway_id"].GetStringValue()
	item := &ec2.CustomerGateway{
		CustomerGatewayId: &customerGatewayID,
	}
	return item, nil
}

//// LIST FUNCTION

func listVpcCustomerGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listVpcCustomerGateways", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := Ec2Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	resp, err := svc.DescribeCustomerGateways(&ec2.DescribeCustomerGatewaysInput{})
	for _, customerGateway := range resp.CustomerGateways {
		d.StreamListItem(ctx, customerGateway)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcCustomerGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpcCustomerGateway")
	customerGateway := h.Item.(*ec2.CustomerGateway)
	defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := Ec2Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeCustomerGatewaysInput{
		CustomerGatewayIds: []*string{customerGateway.CustomerGatewayId},
	}

	// Get call
	op, err := svc.DescribeCustomerGateways(params)
	if err != nil {
		logger.Debug("getVpcCustomerGateway__", "ERROR", err)
		return nil, err
	}

	if op.CustomerGateways != nil && len(op.CustomerGateways) > 0 {
		return op.CustomerGateways[0], nil
	}
	return nil, nil
}

func getVpcCustomerGatewayTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcCustomerGatewayTurbotAkas")
	customerGateway := h.Item.(*ec2.CustomerGateway)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":customer-gateway/" + *customerGateway.CustomerGatewayId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcCustomerGatewayTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	customerGateway := d.HydrateItem.(*ec2.CustomerGateway)
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
