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

//// TABLE DEFINITION

func tableAwsVpcVerifiedAccessGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_verified_access_group",
		Description: "AWS VPC Verified Access Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("verified_access_group_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue", "InvalidVerifiedAccessGroupId.NotFound", "InvalidAction"}),
			},
			Hydrate: getVpcVerifiedAccessGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcVerifiedAccessGroups,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			// DescribeVerifiedAccessGroups API accept group id as input param.
			// We are passing MaxResults value as DescribeVerifiedAccessGroups api input
			// We can not pass both MaxResults and VerifiedAccessGroupId at a time in the same input, we have to pass either one. So verified_access_group_id can not be added as optional quals and added get config for filtering out the group by their id.
			KeyColumns: []*plugin.KeyColumn{
				{Name: "verified_access_instance_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "verified_access_group_id",
				Description: "The ID of the verified access group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the verified access group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VerifiedAccessGroupArn"),
			},
			{
				Name:        "verified_access_instance_id",
				Description: "The ID of the AWS Verified Access instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "deletion_time",
				Description: "The deleteion time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description for the AWS verified access group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_time",
				Description: "The last updated time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "owner",
				Description: "The AWS account number that owns the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(verifiedAccessGroupTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(verifiedAccessGroupTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VerifiedAccessGroupArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcVerifiedAccessGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_verified_access_group.listVpcVerifiedAccessGroups", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(200)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				maxLimit = 5
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeVerifiedAccessGroupsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if d.KeyColumnQualString("verified_access_instance_id") != "" {
		input.VerifiedAccessInstanceId = aws.String(d.KeyColumnQualString("verified_access_instance_id"))
	}

	for {
		// List call
		resp, err := svc.DescribeVerifiedAccessGroups(ctx, input)

		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_verified_access_group.listVpcVerifiedAccessGroups", "api_error", err)
			return nil, nil
		}

		for _, group := range resp.VerifiedAccessGroups {
			d.StreamListItem(ctx, group)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if resp.NextToken == nil {
			break
		} else {
			input.NextToken = resp.NextToken
		}
	}

	return nil, err
}

//// HYDRATED FUNCTION

func getVpcVerifiedAccessGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	groupId := d.KeyColumnQuals["verified_access_group_id"].GetStringValue()

	// Empty check
	if groupId == "" {
		return nil, nil
	}

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_verified_access_group.getVpcVerifiedAccessGroup", "connection_error", err)
		return nil, err
	}

	// Build the params
	input := &ec2.DescribeVerifiedAccessGroupsInput{
		VerifiedAccessGroupIds: []string{groupId},
	}

	// Get call
	op, err := svc.DescribeVerifiedAccessGroups(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_verified_access_group.getVpcVerifiedAccessGroup", "api_error", err)
		return nil, err
	}

	if op.VerifiedAccessGroups != nil && len(op.VerifiedAccessGroups) > 0 {
		return op.VerifiedAccessGroups[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func verifiedAccessGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	group := d.HydrateItem.(types.VerifiedAccessGroup)

	// Get the resource tags
	var turbotTagsMap map[string]string
	if group.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range group.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func verifiedAccessGroupTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	group := d.HydrateItem.(types.VerifiedAccessGroup)
	title := group.VerifiedAccessGroupId

	if group.Tags != nil {
		for _, i := range group.Tags {
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}

	return title, nil
}
