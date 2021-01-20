package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcNetworkACL(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_network_acl",
		Description: "AWS VPC Network ACL",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("network_acl_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidNetworkAclID.NotFound"}),
			ItemFromKey:       networkACLFromKey,
			Hydrate:           getVpcNetworkACL,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcNetworkACLs,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "network_acl_id",
				Description: "The ID of the network ACL",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_default",
				Description: "Indicates whether this is the default network ACL for the VPC",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC for the network ACL",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account that owns the network ACL",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associations",
				Description: "Any associations between the network ACL and one or more subnets",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "entries",
				Description: "One or more entries (rules) in the network ACL",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_raw",
				Description: "A list of tags that are attached to Network ACL",
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
				Hydrate:     getVpcNetworkACLTurbotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// ITEM FROM KEY

func networkACLFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	networkACLID := quals["network_acl_id"].GetStringValue()
	item := &ec2.NetworkAcl{
		NetworkAclId: &networkACLID,
	}
	return item, nil
}

//// LIST FUNCTION

func listVpcNetworkACLs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listVpcNetworkACLs", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := Ec2Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeNetworkAclsPages(
		&ec2.DescribeNetworkAclsInput{},
		func(page *ec2.DescribeNetworkAclsOutput, isLast bool) bool {
			for _, networkACL := range page.NetworkAcls {
				d.StreamListItem(ctx, networkACL)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcNetworkACL(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpcNetworkACL")
	networkACL := h.Item.(*ec2.NetworkAcl)
	defaultRegion := GetDefaultRegion()

	// get service
	svc, err := Ec2Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []*string{networkACL.NetworkAclId},
	}

	// Get call
	op, err := svc.DescribeNetworkAcls(params)
	if err != nil {
		logger.Debug("getVpcNetworkACL__", "ERROR", err)
		return nil, err
	}

	if op.NetworkAcls != nil && len(op.NetworkAcls) > 0 {
		return op.NetworkAcls[0], nil
	}
	return nil, nil
}

func getVpcNetworkACLTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcNetworkACLTurbotAkas")
	networkACL := h.Item.(*ec2.NetworkAcl)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":network-acl/" + *networkACL.NetworkAclId}

	return akas, nil
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
