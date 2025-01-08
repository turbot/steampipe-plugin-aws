package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2Endpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsVpcSubnet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_subnet",
		Description: "AWS VPC Subnet",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("subnet_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidSubnetID.Malformed", "InvalidSubnetID.NotFound"}),
			},
			Hydrate: getVpcSubnet,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeSubnets"},
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcSubnets,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeSubnets"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "availability_zone", Require: plugin.Optional},
				{Name: "availability_zone_id", Require: plugin.Optional},
				{Name: "available_ip_address_count", Require: plugin.Optional},
				{Name: "cidr_block", Require: plugin.Optional},
				{Name: "default_for_az", Require: plugin.Optional},
				{Name: "outpost_arn", Require: plugin.Optional},
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "subnet_arn", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2Endpoint.EC2ServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "subnet_id",
				Description: "Contains the unique ID to specify a subnet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnet_arn",
				Description: "Contains the Amazon Resource Name (ARN) of the subnet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "ID of the VPC, the subnet is in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cidr_block",
				Description: "Contains the IPv4 CIDR block assigned to the subnet.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "state",
				Description: "Current state of the subnet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "Contains the AWS account that own the subnet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "assign_ipv6_address_on_creation",
				Description: "Indicates whether a network interface created in this subnet (including a network interface created by RunInstances) receives an IPv6 address.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "available_ip_address_count",
				Description: "The number of unused private IPv4 addresses in the subnet. The IPv4 addresses for any stopped instances are considered unavailable.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "availability_zone",
				Description: "The Availability Zone of the subnet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zone_id",
				Description: "The AZ ID of the subnet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_owned_ipv4_pool",
				Description: "The customer-owned IPv4 address pool associated with the subnet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_for_az",
				Description: "Indicates whether this is the default subnet for the Availability Zone.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "map_customer_owned_ip_on_launch",
				Description: "Indicates whether a network interface created in this subnet (including a network interface created by RunInstances) receives a customer-owned IPv4 address.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "map_public_ip_on_launch",
				Description: "Indicates whether instances launched in this subnet receive a public IPv4 address.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "outpost_arn",
				Description: "The Amazon Resource Name (ARN) of the Outpost. Available only if subnet is on an outpost.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ipv6_cidr_block_association_set",
				Description: "A list of IPv6 CIDR blocks associated with the subnet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the subnet.",
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

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_subnet.listVpcSubnets", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				maxLimit = int32(5)
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeSubnetsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "availability_zone", FilterName: "availability-zone", ColumnType: "string"},
		{ColumnName: "availability_zone_id", FilterName: "availability-zone-id", ColumnType: "string"},
		{ColumnName: "available_ip_address_count", FilterName: "available-ip-address-count", ColumnType: "int64"},
		{ColumnName: "cidr_block", FilterName: "cidr-block", ColumnType: "cidr"},
		{ColumnName: "default_for_az", FilterName: "default-for-az", ColumnType: "boolean"},
		{ColumnName: "outpost_arn", FilterName: "outpost-arn", ColumnType: "string"},
		{ColumnName: "owner_id", FilterName: "owner-id", ColumnType: "string"},
		{ColumnName: "state", FilterName: "state", ColumnType: "string"},
		{ColumnName: "subnet_arn", FilterName: "subnet-arn", ColumnType: "string"},
		{ColumnName: "vpc_id", FilterName: "vpc-id", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeSubnetsPaginator(svc, input, func(o *ec2.DescribeSubnetsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_subnet.listVpcSubnets", "api_error", err, "connection_name", d.Connection.Name, "region", d.EqualsQualString(matrixKeyRegion))
			return nil, err
		}

		for _, items := range output.Subnets {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcSubnet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	subnetID := d.EqualsQuals["subnet_id"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_subnet.getVpcSubnet", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeSubnetsInput{
		SubnetIds: []string{subnetID},
	}

	// Get call
	op, err := svc.DescribeSubnets(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_subnet.getVpcSubnet", "api_error", err, "connection_name", d.Connection.Name, "region", d.EqualsQualString(matrixKeyRegion))
		return nil, err
	}

	if op != nil && len(op.Subnets) > 0 {
		return op.Subnets[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getVpcSubnetTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(types.Subnet).Tags

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

func getSubnetTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	subnet := d.HydrateItem.(types.Subnet)
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
			title = *subnetData["getVpcSubnet"].(types.Subnet).SubnetId
		} else {
			title = *subnetData["listVpcSubnets"].(types.Subnet).SubnetId
		}
	}
	return title, nil
}
