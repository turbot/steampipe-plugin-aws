package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsVpcVerifiedAccessInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_verified_access_instance",
		Description: "AWS VPC Verified Access Instance",
		List: &plugin.ListConfig{
			Hydrate: listVpcVerifiedAccessInstances,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "verified_access_instance_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
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
				Name:        "description",
				Description: "A description for the AWS Verified Access instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_time",
				Description: "The last updated time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "verified_access_trust_providers",
				Description: "The IDs of the AWS Verified Access trusted providers.",
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.From(verifiedAccessInstanceTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(verifiedAccessInstanceTurbotTags),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcVerifiedAccessInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_verified_access_instance.listVpcVerifiedAccessInstances", "connection_error", err)
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

	input := &ec2.DescribeVerifiedAccessInstancesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if d.EqualsQualString("verified_access_instance_id") != "" {
		input.VerifiedAccessInstanceIds = []string{d.EqualsQualString("verified_access_instance_id")}
	}

	for {
		// List call
		resp, err := svc.DescribeVerifiedAccessInstances(ctx, input)

		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_verified_access_instance.listVpcVerifiedAccessInstances", "api_error", err)
			return nil, nil
		}

		for _, instance := range resp.VerifiedAccessInstances {
			d.StreamListItem(ctx, instance)

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

//// TRANSFORM FUNCTIONS

func verifiedAccessInstanceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	accessPoint := d.HydrateItem.(types.VerifiedAccessInstance)

	// Get the resource tags
	var turbotTagsMap map[string]string
	if accessPoint.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range accessPoint.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func verifiedAccessInstanceTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	accessPoint := d.HydrateItem.(types.VerifiedAccessInstance)
	title := accessPoint.VerifiedAccessInstanceId

	if accessPoint.Tags != nil {
		for _, i := range accessPoint.Tags {
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}

	return title, nil
}
