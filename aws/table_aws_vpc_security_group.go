package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsVpcSecurityGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_security_group",
		Description: "AWS VPC Security Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("group_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidGroupId.Malformed", "InvalidGroupId.NotFound", "InvalidGroup.NotFound"}),
				},
			Hydrate:           getVpcSecurityGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcSecurityGroups,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "description", Require: plugin.Optional},
				{Name: "group_name", Require: plugin.Optional},
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
			},
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
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the security group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpcSecurityGroupARN,
				Transform:   transform.FromValue(),
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
				Hydrate:     getVpcSecurityGroupARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcSecurityGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listVpcSecurityGroups", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeSecurityGroupsInput{
		MaxResults: aws.Int64(1000),
	}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "description", FilterName: "description", ColumnType: "string"},
		{ColumnName: "group_name", FilterName: "group-name", ColumnType: "string"},
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
	err = svc.DescribeSecurityGroupsPages(
		input,
		func(page *ec2.DescribeSecurityGroupsOutput, isLast bool) bool {
			for _, securityGroup := range page.SecurityGroups {
				d.StreamListItem(ctx, securityGroup)

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

func getVpcSecurityGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcSecurityGroup")

	region := d.KeyColumnQualString(matrixKeyRegion)
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

func getVpcSecurityGroupARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcSecurityGroupARN")
	securityGroup := h.Item.(*ec2.SecurityGroup)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":security-group/" + *securityGroup.GroupId

	return arn, nil
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
