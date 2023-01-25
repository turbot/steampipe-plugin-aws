package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsVpcSecurityGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_security_group",
		Description: "AWS VPC Security Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("group_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidGroupId.Malformed", "InvalidGroupId.NotFound", "InvalidGroup.NotFound"}),
			},
			Hydrate: getVpcSecurityGroup,
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
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
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

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group.listVpcSecurityGroups", "connection_error", err)
		return nil, err
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

	input := &ec2.DescribeSecurityGroupsInput{
		MaxResults: aws.Int32(maxLimit),
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

	paginator := ec2.NewDescribeSecurityGroupsPaginator(svc, input, func(o *ec2.DescribeSecurityGroupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_security_group.listVpcSecurityGroups", "api_error", err)
			return nil, err
		}

		for _, items := range output.SecurityGroups {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcSecurityGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	groupID := d.KeyColumnQuals["group_id"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group.getVpcSecurityGroup", "api_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{groupID},
	}

	// Get call
	op, err := svc.DescribeSecurityGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group.getVpcSecurityGroup", "api_error", err)
		return nil, err
	}

	if op.SecurityGroups != nil && len(op.SecurityGroups) > 0 {
		return op.SecurityGroups[0], nil
	}
	return nil, nil
}

func getVpcSecurityGroupARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	securityGroup := h.Item.(types.SecurityGroup)
	region := d.KeyColumnQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group.getVpcSecurityGroupARN", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":security-group/" + *securityGroup.GroupId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getVpcSecurityGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	securityGroup := d.HydrateItem.(types.SecurityGroup)

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
