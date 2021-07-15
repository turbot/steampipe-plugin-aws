package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcNatGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_nat_gateway",
		Description: "AWS VPC Network Address Translation Gateway",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("nat_gateway_id"),
			ShouldIgnoreError: isNotFoundError([]string{"NatGatewayMalformed", "NatGatewayNotFound"}),
			Hydrate:           getVpcNatGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcNatGateways,
		},
		GetMatrixItem: BuildRegionList,
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

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listVpcNatGateways", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeNatGatewaysPages(
		&ec2.DescribeNatGatewaysInput{},
		func(page *ec2.DescribeNatGatewaysOutput, isLast bool) bool {
			for _, securityGroup := range page.NatGateways {
				d.StreamListItem(ctx, securityGroup)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcNatGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcNatGateway")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	natGatewayID := d.KeyColumnQuals["nat_gateway_id"].GetStringValue()

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []*string{aws.String(natGatewayID)},
	}

	// Get call
	op, err := svc.DescribeNatGateways(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getVpcNatGateway__", "ERROR", err)
		return nil, err
	}

	if op.NatGateways != nil && len(op.NatGateways) > 0 {
		return op.NatGateways[0], nil
	}
	return nil, nil
}

func getVpcNatGatewayARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcNatGatewayARN")

	natGateway := h.Item.(*ec2.NatGateway)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Build ARN
	arn := "arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":natgateway/" + *natGateway.NatGatewayId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getVpcNatGatewayTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	natGateway := d.HydrateItem.(*ec2.NatGateway)
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
