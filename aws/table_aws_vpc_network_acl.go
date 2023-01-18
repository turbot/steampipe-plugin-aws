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

func tableAwsVpcNetworkACL(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_network_acl",
		Description: "AWS VPC Network ACL",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("network_acl_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidNetworkAclID.NotFound"}),
			},
			Hydrate: getVpcNetworkACL,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcNetworkACLs,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "is_default", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "network_acl_id",
				Description: "The ID of the network ACL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the network ACL.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpcNetworkACLARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "is_default",
				Description: "Indicates whether this is the default network ACL for the VPC.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC for the network ACL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account that owns the network ACL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associations",
				Description: "Any associations between the network ACL and one or more subnets.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "entries",
				Description: "One or more entries (rules) in the network ACL.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to Network ACL.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getVpcNetworkACLTurbotData, "Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getVpcNetworkACLTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcNetworkACLARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcNetworkACLs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_network_acl.listVpcNetworkACLs", "connection_error", err)
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

	input := &ec2.DescribeNetworkAclsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "is_default", FilterName: "default", ColumnType: "boolean"},
		{ColumnName: "owner_id", FilterName: "owner-id", ColumnType: "string"},
		{ColumnName: "vpc_id", FilterName: "vpc-id", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeNetworkAclsPaginator(svc, input, func(o *ec2.DescribeNetworkAclsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_network_acl.listVpcNetworkACLs", "api_error", err)
			return nil, err
		}

		for _, items := range output.NetworkAcls {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getVpcNetworkACL(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	networkACLID := d.KeyColumnQuals["network_acl_id"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_network_acl.getVpcNetworkACL", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []string{networkACLID},
	}

	// Get call
	op, err := svc.DescribeNetworkAcls(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_network_acl.getVpcNetworkACL", "api_error", err)
		return nil, err
	}

	if op.NetworkAcls != nil && len(op.NetworkAcls) > 0 {
		return op.NetworkAcls[0], nil
	}
	return nil, nil
}

func getVpcNetworkACLARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	networkACL := h.Item.(types.NetworkAcl)
	region := d.KeyColumnQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_network_acl.getVpcNetworkACLARN", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":network-acl/" + *networkACL.NetworkAclId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getVpcNetworkACLTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	networkACL := d.HydrateItem.(types.NetworkAcl)
	param := d.Param.(string)

	// Get resource title
	title := networkACL.NetworkAclId

	// Get the resource tags
	turbotTagsMap := map[string]string{}
	if networkACL.Tags != nil {
		for _, i := range networkACL.Tags {
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
