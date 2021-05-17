package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcEip(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_eip",
		Description: "AWS VPC Elastic IP",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("allocation_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidAllocationID.NotFound", "InvalidAllocationID.Malformed"}),
			Hydrate:           getVpcEip,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcEips,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "allocation_id",
				Description: "Contains the ID representing the allocation of the address for use with EC2-VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the vpc eip.",
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
				Description: "Indicates whether Elastic IP address is for use with instances in EC2-Classic(standard) or instances in a VPC (vpc).",
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
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AllocationId"),
			},
			{
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

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listVpcEips", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	resp, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{})
	for _, address := range resp.Addresses {
		d.StreamListItem(ctx, address)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcEip(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcEip")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	allocationID := d.KeyColumnQuals["allocation_id"].GetStringValue()

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeAddressesInput{
		AllocationIds: []*string{aws.String(allocationID)},
	}

	// Get call
	op, err := svc.DescribeAddresses(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getVpcEip__", "ERROR", err)
		return nil, err
	}

	if op.Addresses != nil && len(op.Addresses) > 0 {
		return op.Addresses[0], nil
	}
	return nil, nil
}

func getVpcEipARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcEipARN")
	eip := h.Item.(*ec2.Address)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get resource arn
	arn := "arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":eip/" + *eip.AllocationId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getVpcEipTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	eip := d.HydrateItem.(*ec2.Address)
	return ec2TagsToMap(eip.Tags)
}
