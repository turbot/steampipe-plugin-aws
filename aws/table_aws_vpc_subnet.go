package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcSubnet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_subnet",
		Description: "AWS VPC Subnet",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("subnet_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidSubnetID.Malformed", "InvalidSubnetID.NotFound"}),
			Hydrate:           getVpcSubnet,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcSubnets,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "subnet_id",
				Description: "Contains the unique ID to specify a subnet",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnet_arn",
				Description: "Contains the Amazon Resource Name (ARN) of the subnet",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "ID of  the VPC, the subnet is in",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cidr_block",
				Description: "Contains the IPv4 CIDR block assigned to the subnet",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "state",
				Description: "Current state of the subnet",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "Contains the AWS account that own the subnet",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "assign_ipv6_address_on_creation",
				Description: "Indicates whether a network interface created in this subnet (including a network interface created by RunInstances) receives an IPv6 address",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "available_ip_address_count",
				Description: "The number of unused private IPv4 addresses in the subnet. The IPv4 addresses for any stopped instances are considered unavailable",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "availability_zone",
				Description: "The Availability Zone of the subnet",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zone_id",
				Description: "The AZ ID of the subnet",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_owned_ipv4_pool",
				Description: "The customer-owned IPv4 address pool associated with the subnet",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_for_az",
				Description: "Indicates whether this is the default subnet for the Availability Zone.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "map_customer_owned_ip_on_launch",
				Description: "Indicates whether a network interface created in this subnet (including a network interface created by RunInstances) receives a customer-owned IPv4 address",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "map_public_ip_on_launch",
				Description: "Indicates whether instances launched in this subnet receive a public IPv4 address",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "outpost_arn",
				Description: "The Amazon Resource Name (ARN) of the Outpost. Available only if subnet is on an outpost",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ipv6_cidr_block_association_set",
				Description: "A list of IPv6 CIDR blocks associated with the subnet",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the subnet",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getVpcSubnetTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getSubnetTurbotTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SubnetArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcSubnets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listVpcSubnets", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeSubnetsPages(
		&ec2.DescribeSubnetsInput{},
		func(page *ec2.DescribeSubnetsOutput, isLast bool) bool {
			for _, subnet := range page.Subnets {
				d.StreamListItem(ctx, subnet)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcSubnet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcSubnet")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	subnetID := d.KeyColumnQuals["subnet_id"].GetStringValue()

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeSubnetsInput{
		SubnetIds: []*string{aws.String(subnetID)},
	}

	// Get call
	op, err := svc.DescribeSubnets(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getVpcSubnet__", "ERROR", err)
		return nil, err
	}

	if op.Subnets != nil && len(op.Subnets) > 0 {
		return op.Subnets[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getVpcSubnetTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	subnet := d.HydrateItem.(*ec2.Subnet)
	return ec2TagsToMap(subnet.Tags)
}

func getSubnetTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	subnet := d.HydrateItem.(*ec2.Subnet)
	subnetData := d.HydrateResults
	var title string
	if subnet.Tags != nil {
		for _, i := range subnet.Tags {
			if *i.Key == "Name" {
				title = *i.Value
			}
		}
	}

	if title == "" {
		if subnetData["getVpcSubnet"] != nil {
			title = *subnetData["getVpcSubnet"].(*ec2.Subnet).SubnetId
		} else {
			title = *subnetData["listVpcSubnets"].(*ec2.Subnet).SubnetId
		}
	}
	return title, nil
}
