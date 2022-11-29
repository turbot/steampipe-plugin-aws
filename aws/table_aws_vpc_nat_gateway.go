package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsVpcNatGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_nat_gateway",
		Description: "AWS VPC Network Address Translation Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("nat_gateway_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NatGatewayMalformed", "NatGatewayNotFound"}),
			},
			Hydrate: getVpcNatGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcNatGateways,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "state", Require: plugin.Optional},
				{Name: "subnet_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "nat_gateway_id",
				Description: "The ID of the NAT gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the NAT gateway.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpcNatGatewayARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "nat_gateway_addresses",
				Description: "Information about the IP addresses and network interface associated with the NAT gateway.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "state",
				Description: "The current state of the NAT gateway (pending | failed | available | deleting | deleted).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The date and time the NAT gateway was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC in which the NAT gateway is located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnet_id",
				Description: "The ID of the subnet in which the NAT gateway is located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delete_time",
				Description: "The date and time the NAT gateway was deleted, if applicable.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "failure_code",
				Description: "If the NAT gateway could not be created, specifies the error code for the failure. (InsufficientFreeAddressesInSubnet | Gateway.NotAttached | InvalidAllocationID.NotFound | Resource.AlreadyAssociated | InternalError | InvalidSubnetID.NotFound).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "failure_message",
				Description: "If the NAT gateway could not be created, specifies the error message for the failure.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioned_bandwidth",
				Description: "Reserved. If you need to sustain traffic greater than the documented limits (https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat-gateway.html).",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ProvisionedBandwidth"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to NAT gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getVpcNatGatewayTurbotData, "Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getVpcNatGatewayTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcNatGatewayARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcNatGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_nat_gateway.listVpcNatGateways", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
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

	input := &ec2.DescribeNatGatewaysInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "state", FilterName: "state", ColumnType: "string"},
		{ColumnName: "subnet_id", FilterName: "subnet-id", ColumnType: "string"},
		{ColumnName: "vpc_id", FilterName: "vpc-id", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filter = filters
	}

	paginator := ec2.NewDescribeNatGatewaysPaginator(svc, input, func(o *ec2.DescribeNatGatewaysPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_nat_gateway.listVpcNatGateways", "api_error", err)
			return nil, err
		}

		for _, items := range output.NatGateways {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcNatGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	natGatewayID := d.KeyColumnQuals["nat_gateway_id"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_nat_gateway.getVpcNatGateway", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []string{natGatewayID},
	}

	// Get call
	op, err := svc.DescribeNatGateways(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_nat_gateway.getVpcNatGateway", "api_error", err)
		return nil, err
	}

	if op.NatGateways != nil && len(op.NatGateways) > 0 {
		return op.NatGateways[0], nil
	}
	return nil, nil
}

func getVpcNatGatewayARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	natGateway := h.Item.(types.NatGateway)
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_nat_gateway.getVpcNatGatewayARN", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Build ARN
	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":natgateway/" + *natGateway.NatGatewayId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getVpcNatGatewayTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	natGateway := d.HydrateItem.(types.NatGateway)
	param := d.Param.(string)

	// Get resource title
	title := natGateway.NatGatewayId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if natGateway.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range natGateway.Tags {
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
