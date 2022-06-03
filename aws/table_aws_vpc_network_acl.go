package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsVpcNetworkACL(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_network_acl",
		Description: "AWS VPC Network ACL",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("network_acl_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidNetworkAclID.NotFound"}),
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
		GetMatrixItem: BuildRegionList,
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listVpcNetworkACLs", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeNetworkAclsInput{
		MaxResults: aws.Int64(1000),
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
	err = svc.DescribeNetworkAclsPages(
		input,
		func(page *ec2.DescribeNetworkAclsOutput, isLast bool) bool {
			for _, networkACL := range page.NetworkAcls {
				d.StreamListItem(ctx, networkACL)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcNetworkACL(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcNetworkACL")

	region := d.KeyColumnQualString(matrixKeyRegion)
	networkACLID := d.KeyColumnQuals["network_acl_id"].GetStringValue()

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []*string{aws.String(networkACLID)},
	}

	// Get call
	op, err := svc.DescribeNetworkAcls(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getVpcNetworkACL__", "ERROR", err)
		return nil, err
	}

	if op.NetworkAcls != nil && len(op.NetworkAcls) > 0 {
		return op.NetworkAcls[0], nil
	}
	return nil, nil
}

func getVpcNetworkACLARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcNetworkACLARN")
	networkACL := h.Item.(*ec2.NetworkAcl)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":network-acl/" + *networkACL.NetworkAclId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getVpcNetworkACLTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	networkACL := d.HydrateItem.(*ec2.NetworkAcl)
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
