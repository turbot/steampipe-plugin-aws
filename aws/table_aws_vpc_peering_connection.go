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

func tableAwsVpcPeeringConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_peering_connection",
		Description: "AWS VPC Peering Connection",
		List: &plugin.ListConfig{
			Hydrate: listVpcPeeringConnections,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the VPC peering connection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VpcPeeringConnectionId"),
			},
			{
				Name:        "accepter_cidr_block",
				Description: "The IPv4 CIDR block for the accepter VPC.",
				Type:        proto.ColumnType_STRING,
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
				Type:        proto.ColumnType_STRING,
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
				Name:        "status_code",
				Description: "The status of the VPC peering connection. Possible values include: 'pending-acceptance', 'failed', 'expired', 'provisioning', 'active', 'deleting', 'deleted' or 'rejected'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.Code"),
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

	params := getRequestParameters(d)

	// List call
	err = svc.DescribeVpcPeeringConnectionsPages(
		params,
		func(page *ec2.DescribeVpcPeeringConnectionsOutput, isLast bool) bool {
			for _, connection := range page.VpcPeeringConnections {
				d.StreamListItem(ctx, connection)
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

func getRequestParameters(d *plugin.QueryData) (*ec2.DescribeVpcPeeringConnectionsInput) {
	filters := []*ec2.Filter{}
	equalQuals := d.KeyColumnQuals
	params := &ec2.DescribeVpcPeeringConnectionsInput{
		MaxResults: aws.Int64(100),
	}

	if equalQuals["accepter_cidr_block"] != nil {
		filter := ec2.Filter{
			Name:   aws.String("accepter-vpc-info.cidr-block"),
			Values: []*string{aws.String(equalQuals["accepter_cidr_block"].GetStringValue())},
		}
		filters = append(filters, &filter)
	}

	if equalQuals["accepter_owner_id"] != nil {
		filter := ec2.Filter{
			Name:   aws.String("accepter-vpc-info.owner-id"),
			Values: []*string{aws.String(equalQuals["accepter_owner_id"].GetStringValue())},
		}
		filters = append(filters, &filter)
	}

	if equalQuals["accepter_vpc_id"] != nil {
		filter := ec2.Filter{
			Name:   aws.String("accepter-vpc-info.vpc-id"),
			Values: []*string{aws.String(equalQuals["accepter_vpc_id"].GetStringValue())},
		}
		filters = append(filters, &filter)
	}

	if equalQuals["expiration_time"] != nil {
		filter := ec2.Filter{
			Name:   aws.String("expiration-time"),
			Values: []*string{aws.String(equalQuals["expiration_time"].GetStringValue())},
		}
		filters = append(filters, &filter)
	}

	if equalQuals["requester_cidr_block"] != nil {
		filter := ec2.Filter{
			Name:   aws.String("requester-vpc-info.cidr-block"),
			Values: []*string{aws.String(equalQuals["requester_cidr_block"].GetStringValue())},
		}
		filters = append(filters, &filter)
	}

	if equalQuals["requester_owner_id"] != nil {
		filter := ec2.Filter{
			Name:   aws.String("requester-vpc-info.owner-id"),
			Values: []*string{aws.String(equalQuals["requester_owner_id"].GetStringValue())},
		}
		filters = append(filters, &filter)
	}

	if equalQuals["requester_vpc_id"] != nil {
		filter := ec2.Filter{
			Name:   aws.String("requester-vpc-info.vpc-id"),
			Values: []*string{aws.String(equalQuals["requester_vpc_id"].GetStringValue())},
		}
		filters = append(filters, &filter)
	}

	if equalQuals["status_code"] != nil {
		filter := ec2.Filter{
			Name:   aws.String("status-code"),
			Values: []*string{aws.String(equalQuals["status_code"].GetStringValue())},
		}
		filters = append(filters, &filter)
	}

	if equalQuals["status_message"] != nil {
		filter := ec2.Filter{
			Name:   aws.String("status-message"),
			Values: []*string{aws.String(equalQuals["status_message"].GetStringValue())},
		}
		filters = append(filters, &filter)
	}

	if equalQuals["id"] != nil {
		filter := ec2.Filter{
			Name:   aws.String("vpc-peering-connection-id"),
			Values: []*string{aws.String(equalQuals["id"].GetStringValue())},
		}
		filters = append(filters, &filter)
	}

	// Add filters as request parameter when at least one filter is present
	if len(filters) > 0 {
		params.Filters = filters
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxResults {
			params.MaxResults = limit
		}
	}

	return params;
}
