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
		List: &plugin.ListConfig{
			Hydrate: listVpcVerifiedAccessGroups,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "verified_access_group_id", Require: plugin.Optional},
				{Name: "verified_access_instance_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "verified_access_group_id",
				Description: "The ID of the Verified Access group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the Verified Access group.",
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
				Description: "A description for the AWS Verified Access group.",
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
				Transform:   transform.FromP(verifiedAccessGroupTurbotData, "Title"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(verifiedAccessGroupTurbotData, "Tags"),
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

	if d.KeyColumnQualString("verified_access_group_id") != "" {
		input.VerifiedAccessGroupIds = []string{d.KeyColumnQualString("verified_access_group_id")}
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

		for _, instance := range resp.VerifiedAccessGroups {
			d.StreamListItem(ctx, instance)

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

//// TRANSFORM FUNCTIONS

func verifiedAccessGroupTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	group := d.HydrateItem.(types.VerifiedAccessGroup)
	param := d.Param.(string)
	title := group.VerifiedAccessInstanceId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if group.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range group.Tags {
			turbotTagsMap[*i.Key] = *i.Value
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}

	if param == "Tags" {
		if group.Tags == nil {
			return nil, nil
		}
		return turbotTagsMap, nil
	}

	return title, nil
}
