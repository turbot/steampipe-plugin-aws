package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsVpcCustomerGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_customer_gateway",
		Description: "AWS VPC Customer Gateway",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("customer_gateway_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidCustomerGatewayID.NotFound", "InvalidCustomerGatewayID.Malformed"}),
			Hydrate:           getVpcCustomerGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcCustomerGateways,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "ip_address", Require: plugin.Optional},
				{Name: "bgp_asn", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
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
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
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

	// List call
	resp, err := svc.DescribeCustomerGateways(input)
	for _, customerGateway := range resp.CustomerGateways {
		d.StreamListItem(ctx, customerGateway)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcCustomerGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcCustomerGateway")

	region := d.KeyColumnQualString(matrixKeyRegion)
	customerGatewayID := d.KeyColumnQuals["customer_gateway_id"].GetStringValue()

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeCustomerGatewaysInput{
		CustomerGatewayIds: []*string{aws.String(customerGatewayID)},
	}

	// Get call
	op, err := svc.DescribeCustomerGateways(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getVpcCustomerGateway__", "ERROR", err)
		return nil, err
	}

	if op.CustomerGateways != nil && len(op.CustomerGateways) > 0 {
		return op.CustomerGateways[0], nil
	}
	return nil, nil
}

func getVpcCustomerGatewayTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcCustomerGatewayTurbotAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	customerGateway := h.Item.(*ec2.CustomerGateway)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":customer-gateway/" + *customerGateway.CustomerGatewayId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcCustomerGatewayTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
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
