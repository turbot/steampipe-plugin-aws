package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsVpcSecurityGroupVpcAssociation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_security_group_vpc_association",
		Description: "AWS VPC Security Group VPC Association",
		List: &plugin.ListConfig{
			Hydrate: listVpcSecurityGroupVpcAssociations,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeSecurityGroupVpcAssociations"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "group_id", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidFilter"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_EC2_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "group_id",
				Description: "The ID of the security group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_owner_id",
				Description: "The ID of the AWS account that owns the VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the association.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getVpcSecurityGroupVpcAssociationTitle),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcSecurityGroupVpcAssociations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group_vpc_association.listVpcSecurityGroupVpcAssociations", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
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

	// Build the input parameters
	input := &ec2.DescribeSecurityGroupVpcAssociationsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	// Add filters if provided
	filters := []types.Filter{}

	if d.EqualsQuals["vpc_id"] != nil {
		filters = append(filters, types.Filter{
			Name:   aws.String("vpc-id"),
			Values: []string{d.EqualsQuals["vpc_id"].GetStringValue()},
		})
	}

	if d.EqualsQuals["group_id"] != nil {
		filters = append(filters, types.Filter{
			Name:   aws.String("group-id"),
			Values: []string{d.EqualsQuals["group_id"].GetStringValue()},
		})
	}

	if d.EqualsQuals["state"] != nil {
		filters = append(filters, types.Filter{
			Name:   aws.String("state"),
			Values: []string{d.EqualsQuals["state"].GetStringValue()},
		})
	}

	if len(filters) > 0 {
		input.Filters = filters
	}

	// Create paginator for DescribeSecurityGroupVpcAssociations
	paginator := ec2.NewDescribeSecurityGroupVpcAssociationsPaginator(svc, input, func(o *ec2.DescribeSecurityGroupVpcAssociationsPaginatorOptions) {
		o.Limit = maxLimit
	})

	// Iterate through the paginator
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_security_group_vpc_association.listVpcSecurityGroupVpcAssociations", "api_error", err)
			return nil, err
		}

		for _, association := range output.SecurityGroupVpcAssociations {
			d.StreamListItem(ctx, association)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getVpcSecurityGroupVpcAssociationTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	association := d.HydrateItem.(types.SecurityGroupVpcAssociation)
	return *association.GroupId + " - " + *association.VpcId, nil
}
