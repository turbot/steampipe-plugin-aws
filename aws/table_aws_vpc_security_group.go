package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcSecurityGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_security_group",
		Description: "AWS VPC Security Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("group_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidGroupId.Malformed", "InvalidGroupId.NotFound"}),
			Hydrate:           getVpcSecurityGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcSecurityGroups,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "group_name",
				Description: "The friendly name that identifies the security group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "group_id",
				Description: "Contains the unique ID to identify a security group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description of the security group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC for the security group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "Contains the AWS account ID of the owner of the security group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_permissions",
				Description: "A list of inbound rules associated with the security group",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ip_permissions_egress",
				Description: "A list of outbound rules associated with the security group",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the security group",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getVpcSecurityGroupTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GroupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcSecurityGroupTurbotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcSecurityGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listVpcSecurityGroups", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeSecurityGroupsPages(
		&ec2.DescribeSecurityGroupsInput{},
		func(page *ec2.DescribeSecurityGroupsOutput, isLast bool) bool {
			for _, securityGroup := range page.SecurityGroups {
				d.StreamListItem(ctx, securityGroup)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcSecurityGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcSecurityGroup")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	groupID := d.KeyColumnQuals["group_id"].GetStringValue()

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{aws.String(groupID)},
	}

	// Get call
	op, err := svc.DescribeSecurityGroups(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getVpcSecurityGroup__", "ERROR", err)
		return nil, err
	}

	if op.SecurityGroups != nil && len(op.SecurityGroups) > 0 {
		return op.SecurityGroups[0], nil
	}
	return nil, nil
}

func getVpcSecurityGroupTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcSecurityGroupTurbotAkas")
	securityGroup := h.Item.(*ec2.SecurityGroup)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":security-group/" + *securityGroup.GroupId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcSecurityGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	securityGroup := d.HydrateItem.(*ec2.SecurityGroup)

	// Get the resource tags
	if securityGroup.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range securityGroup.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
