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

func tableAwsVpcVerifiedAccessTrustProvider(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_verified_access_trust_provider",
		Description: "AWS VPC Verified Access Trust Provider",
		List: &plugin.ListConfig{
			Hydrate: listVpcVerifiedAccessTrustProviders,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "verified_access_trust_provider_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "verified_access_trust_provider_id",
				Description: "The ID of the AWS Verified Access trust provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description for the AWS Verified Access trust provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "device_trust_provider_type",
				Description: "The type of device-based trust provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_time",
				Description: "The last updated time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "policy_reference_name",
				Description: "The identifier to be used when working with policy rules.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "trust_provider_type",
				Description: "The type of Verified Access trust provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_trust_provider_type",
				Description: "The type of user-based trust provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "oidc_options",
				Description: "The OpenID Connect details for an oidc-type, user-identity based trust provider.",
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
				Transform:   transform.FromP(trustProviderTurbotData, "Title"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(trustProviderTurbotData, "Tags"),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcVerifiedAccessTrustProviders(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_verified_access_trust_provider.listVpcVerifiedAccessTrustProviders", "connection_error", err)
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

	input := &ec2.DescribeVerifiedAccessTrustProvidersInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if d.KeyColumnQualString("verified_access_trust_provider_id") != "" {
		input.VerifiedAccessTrustProviderIds = []string{d.KeyColumnQualString("verified_access_trust_provider_id")}
	}

	for {
		// List call
		resp, err := svc.DescribeVerifiedAccessTrustProviders(ctx, input)

		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_verified_access_trust_provider.listVpcVerifiedAccessTrustProviders", "api_error", err)
			return nil, nil
		}

		for _, provider := range resp.VerifiedAccessTrustProviders {
			d.StreamListItem(ctx, provider)

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

func trustProviderTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	accessPoint := d.HydrateItem.(types.VerifiedAccessTrustProvider)
	param := d.Param.(string)
	title := accessPoint.VerifiedAccessTrustProviderId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if accessPoint.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range accessPoint.Tags {
			turbotTagsMap[*i.Key] = *i.Value
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}

	if param == "Tags" {
		if accessPoint.Tags == nil {
			return nil, nil
		}
		return turbotTagsMap, nil
	}

	return title, nil
}
