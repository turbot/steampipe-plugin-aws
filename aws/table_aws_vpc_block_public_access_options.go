package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAwsVpcBlockPublicAccessOptions(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_block_public_access_options",
		Description: "AWS VPC Block Public Access Options",
		List: &plugin.ListConfig{
			Hydrate: listVpcBlockPublicAccessOptions,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeVpcBlockPublicAccessOptions"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_EC2_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "exclusions_allowed",
				Description: "Determines if exclusions are allowed. If you have [enabled VPC BPA at the Organization level], exclusions may be not-allowed. Otherwise, they are allowed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "internet_gateway_block_mode",
				Description: "The current mode of VPC BPA. Possible values are: 'off', 'block-bidirectional', and 'block-ingress'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_update_timestamp",
				Description: "The last time the VPC BPA mode was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "managed_by",
				Description: "The entity that manages the state of VPC BPA. Possible values include: 'account' and 'declarative-policy'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reason",
				Description: "The reason for the current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of VPC BPA.",
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

func listVpcBlockPublicAccessOptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_block_public_access_options.listVpcBlockPublicAccessOptions", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &ec2.DescribeVpcBlockPublicAccessOptionsInput{}

	op, err := svc.DescribeVpcBlockPublicAccessOptions(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_block_public_access_options.listVpcBlockPublicAccessOptions", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, op.VpcBlockPublicAccessOptions)

	return nil, nil
}
