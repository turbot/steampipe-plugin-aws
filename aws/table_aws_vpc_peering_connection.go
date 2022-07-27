package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsVpcPeeringConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_peering_connection",
		Description: "AWS VPC Peering Connection",
		List: &plugin.ListConfig{
			Hydrate: listVpcPeeringConnections,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "accepter_cidr_block", Require: plugin.Optional},
				{Name: "accepter_owner_id", Require: plugin.Optional},
				{Name: "accepter_vpc_id", Require: plugin.Optional},
				{Name: "expiration_time", Require: plugin.Optional},
				{Name: "requester_cidr_block", Require: plugin.Optional},
				{Name: "requester_owner_id", Require: plugin.Optional},
				{Name: "requester_vpc_id", Require: plugin.Optional},
				{Name: "status_message", Require: plugin.Optional},
				{Name: "id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the VPC peering connection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VpcPeeringConnectionId"),
			},
			{
				Name:        "status_code",
				Description: "The status of the VPC peering connection. Possible values include: 'pending-acceptance', 'failed', 'expired', 'provisioning', 'active', 'deleting', 'deleted' or 'rejected'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.Code"),
			},
			{
				Name:        "accepter_cidr_block",
				Description: "The IPv4 CIDR block for the accepter VPC.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("AccepterVpcInfo.CidrBlock"),
			},
			{
				Name:        "accepter_owner_id",
				Description: "The ID of the Amazon Web Services account that owns the accepter VPC.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccepterVpcInfo.OwnerId"),
			},
			{
				Name:        "accepter_region",
				Description: "The Region in which the accepter VPC is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccepterVpcInfo.Region"),
			},
			{
				Name:        "accepter_vpc_id",
				Description: "The ID of the accepter VPC.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccepterVpcInfo.VpcId"),
			},
			{
				Name:        "expiration_time",
				Description: "The time that an unaccepted VPC peering connection will expire.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "requester_cidr_block",
				Description: "The IPv4 CIDR block for the requester VPC.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("RequesterVpcInfo.CidrBlock"),
			},
			{
				Name:        "requester_owner_id",
				Description: "The ID of the Amazon Web Services account that owns the requester VPC.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RequesterVpcInfo.OwnerId"),
			},
			{
				Name:        "requester_region",
				Description: "The Region in which the requester VPC is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RequesterVpcInfo.Region"),
			},
			{
				Name:        "requester_vpc_id",
				Description: "The ID of the requester VPC.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RequesterVpcInfo.VpcId"),
			},
			{
				Name:        "status_message",
				Description: "A message that provides more information about the status, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.Message"),
			},
			{
				Name:        "accepter_cidr_block_set",
				Description: "Information about the IPv4 CIDR blocks for the accepter VPC.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccepterVpcInfo.CidrBlockSet"),
			},
			{
				Name:        "accepter_ipv6_cidr_block_set",
				Description: "The IPv6 CIDR block for the accepter VPC.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccepterVpcInfo.Ipv6CidrBlockSet"),
			},
			{
				Name:        "accepter_peering_options",
				Description: "Information about the VPC peering connection options for the accepter VPC.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccepterVpcInfo.PeeringOptions"),
			},
			{
				Name:        "requester_cidr_block_set",
				Description: "Information about the IPv4 CIDR blocks for the requester VPC.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RequesterVpcInfo.CidrBlockSet"),
			},
			{
				Name:        "requester_ipv6_cidr_block_set",
				Description: "The IPv6 CIDR block for the requester VPC.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RequesterVpcInfo.Ipv6CidrBlockSet"),
			},
			{
				Name:        "requester_peering_options",
				Description: "Information about the VPC peering connection options for the requester VPC.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RequesterVpcInfo.PeeringOptions"),
			},
			{
				Name:        "tags_src",
				Description: "The tags assigned to the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VpcPeeringConnectionId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(vpcPeeringConnectionTags),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcPeeringConnections(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.KeyColumnQualString(matrixKeyRegion)
	logger.Trace("listVpcPeeringConnections", "AWS_REGION", region)

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		logger.Error("listVpcPeeringConnections", "Ec2Service_error", err)
		return nil, err
	}

	input := &ec2.DescribeVpcPeeringConnectionsInput{
		MaxResults: aws.Int64(1000),
	}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "accepter_cidr_block", FilterName: "accepter-vpc-info.cidr-block", ColumnType: "cidr"},
		{ColumnName: "accepter_owner_id", FilterName: "accepter-vpc-info.owner-id", ColumnType: "string"},
		{ColumnName: "accepter_vpc_id", FilterName: "accepter-vpc-info.vpc-id", ColumnType: "string"},
		{ColumnName: "expiration_time", FilterName: "expiration-time", ColumnType: "time"},
		{ColumnName: "requester_cidr_block", FilterName: "requester-vpc-info.cidr-block", ColumnType: "cidr"},
		{ColumnName: "requester_owner_id", FilterName: "requester-vpc-info.owner-id", ColumnType: "string"},
		{ColumnName: "requester_vpc_id", FilterName: "requester-vpc-info.vpc-id", ColumnType: "string"},
		{ColumnName: "status_code", FilterName: "status-code", ColumnType: "string"},
		{ColumnName: "status_message", FilterName: "status-message", ColumnType: "string"},
		{ColumnName: "id", FilterName: "vpc-peering-connection-id", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 5 {
				input.MaxResults = aws.Int64(5)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.DescribeVpcPeeringConnectionsPages(
		input,
		func(page *ec2.DescribeVpcPeeringConnectionsOutput, isLast bool) bool {
			for _, connection := range page.VpcPeeringConnections {
				d.StreamListItem(ctx, connection)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		logger.Error("listVpcPeeringConnections", "DescribeVpcPeeringConnectionsPages_error", err)
		return nil, err
	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func vpcPeeringConnectionTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	connection := d.HydrateItem.(*ec2.VpcPeeringConnection)

	var turbotTagsMap map[string]string
	if connection.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range connection.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
