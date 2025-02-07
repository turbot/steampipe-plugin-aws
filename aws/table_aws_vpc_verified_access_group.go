package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2Endpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			Tags:    map[string]string{"service": "ec2", "action": "DescribeVerifiedAccessGroups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcVerifiedAccessGroups,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeVerifiedAccessGroups"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			// DescribeVerifiedAccessGroups API accept group id as input param.
			// We are passing MaxResults value as DescribeVerifiedAccessGroups api input
			// We cannot pass both MaxResults and VerifiedAccessGroupId at a time in the same input, we have to pass either one. So verified_access_group_id cannot be added as optional quals. Added get config for filtering out the group by its id.
			KeyColumns: []*plugin.KeyColumn{
				{Name: "verified_access_instance_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2Endpoint.AWS_EC2_SERVICE_ID),
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

	if d.EqualsQualString("verified_access_instance_id") != "" {
		input.VerifiedAccessInstanceId = aws.String(d.EqualsQualString("verified_access_instance_id"))
	}

	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		// List call
		resp, err := svc.DescribeVerifiedAccessGroups(ctx, input)

		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_verified_access_group.listVpcVerifiedAccessGroups", "api_error", err)
			return nil, nil
		}

		for _, group := range resp.VerifiedAccessGroups {
			d.StreamListItem(ctx, group)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
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

	groupId := d.EqualsQuals["verified_access_group_id"].GetStringValue()

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

	if len(op.VerifiedAccessGroups) > 0 {
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
