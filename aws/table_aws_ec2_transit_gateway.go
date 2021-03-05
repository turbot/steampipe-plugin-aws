package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2TransitGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_transit_gateway",
		Description: "AWS EC2 Transit Gateway",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("network_interface_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidTransitGatewayID.NotFound", "InvalidTransitGatewayID.Unavailable", "InvalidTransitGatewayID.Malformed"}),
			Hydrate:           getEc2TransitGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2TransitGateways,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "transit_gateway_id",
				Description: "The ID of the transit gateway",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "transit_gateway_arn",
				Description: "The Amazon Resource Name (ARN) of the transit gateway",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the transit gateway",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account ID that owns the transit gateway",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the transit gateway",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The date and time when transit gateway was created",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "amazon_side_asn",
				Description: "A private Autonomous System Number (ASN) for the Amazon side of a BGP session. The range is 64512 to 65534 for 16-bit ASNs and 4200000000 to 4294967294 for 32-bit ASNs",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Options.AmazonSideAsn"),
			},
			{
				Name:        "association_default_route_table_id",
				Description: "The ID of the default association route table",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.AssociationDefaultRouteTableId"),
			},
			{
				Name:        "auto_accept_shared_attachments",
				Description: "Indicates whether attachment requests are automatically accepted",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.AutoAcceptSharedAttachments"),
			},
			{
				Name:        "default_route_table_association",
				Description: "Indicates whether resource attachments are automatically associated with the default association route table",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.DefaultRouteTableAssociation"),
			},
			{
				Name:        "default_route_table_propagation",
				Description: "Indicates whether resource attachments are automatically associated with the default association route table",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.DefaultRouteTablePropagation"),
			},
			{
				Name:        "dns_support",
				Description: "Indicates whether DNS support is enabled",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.DnsSupport"),
			},
			{
				Name:        "multicast_support",
				Description: "Indicates whether multicast is enabled on the transit gateway",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.MulticastSupport"),
			},
			{
				Name:        "propagation_default_route_table_id",
				Description: "The ID of the default propagation route table",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.PropagationDefaultRouteTableId"),
			},
			{
				Name:        "vpn_ecmp_support",
				Description: "Indicates whether Equal Cost Multipath Protocol support is enabled",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.VpnEcmpSupport"),
			},
			{
				Name:        "cidr_blocks",
				Description: "A list of transit gateway CIDR blocks",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Options.TransitGatewayCidrBlocks"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are assigned to the transit gateway",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getEc2TransitGatewayTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getEc2TransitGatewayTurbotTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TransitGatewayArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2TransitGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeTransitGatewaysPages(
		&ec2.DescribeTransitGatewaysInput{},
		func(page *ec2.DescribeTransitGatewaysOutput, isLast bool) bool {
			for _, transitGateway := range page.TransitGateways {
				d.StreamListItem(ctx, transitGateway)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2TransitGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	transitGatewayID := d.KeyColumnQuals["transit_gateway_id"].GetStringValue()

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeTransitGatewaysInput{
		TransitGatewayIds: []*string{aws.String(transitGatewayID)},
	}

	op, err := svc.DescribeTransitGateways(params)
	if err != nil {
		return nil, err
	}

	if op.TransitGateways != nil && len(op.TransitGateways) > 0 {
		return op.TransitGateways[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getEc2TransitGatewayTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ec2.TransitGateway)
	return ec2TagsToMap(data.Tags)
}

func getEc2TransitGatewayTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ec2.TransitGateway)
	title := data.TransitGatewayId
	if data.Tags != nil {
		for _, i := range data.Tags {
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}
	return title, nil
}
