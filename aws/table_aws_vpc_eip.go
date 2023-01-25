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

func tableAwsVpcEip(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_eip",
		Description: "AWS VPC Elastic IP",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("allocation_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAllocationID.NotFound", "InvalidAllocationID.Malformed"}),
			},
			Hydrate: getVpcEip,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcEips,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "association_id", Require: plugin.Optional},
				{Name: "domain", Require: plugin.Optional},
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "network_border_group", Require: plugin.Optional},
				{Name: "network_interface_id", Require: plugin.Optional},
				{Name: "network_interface_owner_id", Require: plugin.Optional},
				{Name: "private_ip_address", Require: plugin.Optional},
				{Name: "public_ip", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "allocation_id",
				Description: "Contains the ID representing the allocation of the address for use with EC2-VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				// EIPs in EC2-Classic have no valid ARN due to no allocation ID
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the VPC EIP.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpcEipARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "public_ip",
				Description: "Contains the Elastic IP address.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "public_ipv4_pool",
				Description: "The ID of an address pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain",
				Description: "Indicates whether Elastic IP address is for use with instances in EC2-Classic (standard) or instances in a VPC (vpc).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "association_id",
				Description: "Contains the ID representing the association of the address with an instance in a VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "carrier_ip",
				Description: "The carrier IP address associated. This option is only available for network interfaces which reside in a subnet in a Wavelength Zone (for example an EC2 instance).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_owned_ip",
				Description: "The customer-owned IP address.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "customer_owned_ipv4_pool",
				Description: "The ID of the customer-owned address pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_id",
				Description: "Contains the ID of the instance that the address is associated with.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_border_group",
				Description: "The name of the unique set of Availability Zones, Local Zones, or Wavelength Zones from which AWS advertises IP addresses.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_interface_id",
				Description: "The ID of the network interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_interface_owner_id",
				Description: "The ID of the AWS account that owns the network interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_ip_address",
				Description: "The private IP address associated with the Elastic IP address.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the vpc.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getVpcEipTurbotTags),
			},
			{
				// Fallback to public IP for EIPs in EC2-Classic
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AllocationId", "PublicIp"),
			},
			{
				// EIPs in EC2-Classic have no valid ARN, so no valid AKAs either
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcEipARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcEips(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_eip.listVpcEips", "connection_error", err)
		return nil, err
	}

	input := &ec2.DescribeAddressesInput{}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "association_id", FilterName: "association-id", ColumnType: "string"},
		{ColumnName: "domain", FilterName: "domain", ColumnType: "string"},
		{ColumnName: "instance_id", FilterName: "instance-id", ColumnType: "string"},
		{ColumnName: "network_border_group", FilterName: "network-border-group", ColumnType: "string"},
		{ColumnName: "network_interface_id", FilterName: "network-interface-id", ColumnType: "string"},
		{ColumnName: "network_interface_owner_id", FilterName: "network-interface-owner-id", ColumnType: "string"},
		{ColumnName: "private_ip_address", FilterName: "private-ip-address", ColumnType: "ipaddr"},
		{ColumnName: "public_ip", FilterName: "public-ip", ColumnType: "ipaddr"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// List call
	resp, err := svc.DescribeAddresses(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_eip.listVpcEips", "api_error", err)
	}
	for _, address := range resp.Addresses {
		d.StreamListItem(ctx, address)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcEip(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	allocationID := d.KeyColumnQuals["allocation_id"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_eip.getVpcEip", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeAddressesInput{
		AllocationIds: []string{allocationID},
	}

	// Get call
	op, err := svc.DescribeAddresses(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_eip.getVpcEip", "api_error", err)
		return nil, err
	}

	if op.Addresses != nil && len(op.Addresses) > 0 {
		return op.Addresses[0], nil
	}
	return nil, nil
}

func getVpcEipARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)

	eip := h.Item.(types.Address)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_eip.getVpcEipARN", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// EIPs in EC2-Classic do not have an allocation ID, therefore no valid ARN
	if eip.AllocationId == nil {
		return nil, nil
	}

	// Get resource ARN
	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":eip/" + *eip.AllocationId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getVpcEipTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(types.Address).Tags
	var turbotTagsMap map[string]string
	if tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}
