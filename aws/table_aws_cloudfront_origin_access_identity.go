package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudFrontOriginAccessIdentity(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudfront_origin_access_identity",
		Description: "AWS CloudFront Origin Access Identity",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchCloudFrontOriginAccessIdentity"}),
			},
			Hydrate: getCloudFrontOriginAccessIdentity,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudFrontOriginAccessIdentities,
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID for the origin access identity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id", "CloudFrontOriginAccessIdentity.Id"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the origin access identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFrontOriginAccessIdentityARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "s3_canonical_user_id",
				Description: "The Amazon S3 canonical user ID for the origin access identity, which you use when giving the origin access identity read permission to an object in Amazon S3.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("S3CanonicalUserId", "CloudFrontOriginAccessIdentity.S3CanonicalUserId"),
			},
			{
				Name:        "caller_reference",
				Description: "A unique value that ensures that the request can't be replayed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CloudFrontOriginAccessIdentity.CloudFrontOriginAccessIdentityConfig.CallerReference"),
				Hydrate:     getCloudFrontOriginAccessIdentity,
			},
			{
				Name:        "comment",
				Description: "The comment for this origin access identity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Comment", "CloudFrontOriginAccessIdentity.CloudFrontOriginAccessIdentityConfig.Comment"),
			},
			{
				Name:        "etag",
				Description: "The current version of the origin access identity's information.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFrontOriginAccessIdentity,
				Transform:   transform.FromField("ETag"),
			},

			//  Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id", "CloudFrontOriginAccessIdentity.Id"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFrontOriginAccessIdentityARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudFrontOriginAccessIdentities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_origin_access_identity.listCloudFrontOriginAccessIdentities", "client_error", err)
		return nil, err
	}

	// The maximum number for MaxItems parameter is not defined by the API
	// We have set the MaxItems to 1000 based on our test
	maxItems := int32(1000)

	// Reduce the basic request limit down if the user has only requested a small number
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input := &cloudfront.ListCloudFrontOriginAccessIdentitiesInput{
		MaxItems: &maxItems,
	}

	paginator := cloudfront.NewListCloudFrontOriginAccessIdentitiesPaginator(svc, input, func(o *cloudfront.ListCloudFrontOriginAccessIdentitiesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudfront_origin_access_identity.listCloudFrontOriginAccessIdentities", "api_error", err)
			return nil, err
		}

		for _, identity := range output.CloudFrontOriginAccessIdentityList.Items {
			d.StreamListItem(ctx, identity)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudFrontOriginAccessIdentity(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_origin_access_identity.getCloudFrontOriginAccessIdentity", "client_error", err)
		return nil, err
	}

	var identityID string
	if h.Item != nil {
		identityID = *h.Item.(types.CloudFrontOriginAccessIdentitySummary).Id
	} else {
		identityID = d.KeyColumnQuals["id"].GetStringValue()
	}

	if strings.TrimSpace(identityID) == "" {
		return nil, nil
	}

	params := &cloudfront.GetCloudFrontOriginAccessIdentityInput{
		Id: &identityID,
	}

	op, err := svc.GetCloudFrontOriginAccessIdentity(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_origin_access_identity.getCloudFrontOriginAccessIdentity", "api_error", err)
		return nil, err
	}

	return *op, nil
}

func getCloudFrontOriginAccessIdentityARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	originAccessIdentityData := *originAccessIdentityID(h.Item)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_origin_access_identity.getCloudFrontOriginAccessIdentityARN", "common_data_error", err)
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":cloudfront::" + commonColumnData.AccountId + ":origin-access-identity/" + originAccessIdentityData

	return arn, nil
}

func originAccessIdentityID(item interface{}) *string {
	switch item := item.(type) {
	case cloudfront.GetCloudFrontOriginAccessIdentityOutput:
		return item.CloudFrontOriginAccessIdentity.Id
	case types.CloudFrontOriginAccessIdentitySummary:
		return item.Id
	}
	return nil
}
