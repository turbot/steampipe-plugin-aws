package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudFrontOriginAccessIdentity(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudfront_origin_access_identity",
		Description: "AWS CloudFront Origin Access Identity",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchCloudFrontOriginAccessIdentity"}),
			},
			Hydrate: getCloudFrontOriginAccessIdentity,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudFrontOriginAccessIdentities,
		},
		Columns: awsColumns([]*plugin.Column{
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
	plugin.Logger(ctx).Trace("listCloudFrontOriginAccessIdentities")

	// Create session
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The maximum number for MaxItems parameter is not defined by the API
	// We have set the MaxItems to 1000 based on our test
	input := &cloudfront.ListCloudFrontOriginAccessIdentitiesInput{
		MaxItems: aws.Int64(1000),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxItems {
			if *limit < 1 {
				input.MaxItems = types.Int64(1)
			} else {
				input.MaxItems = limit
			}
		}
	}

	// List call
	err = svc.ListCloudFrontOriginAccessIdentitiesPages(
		input,
		func(page *cloudfront.ListCloudFrontOriginAccessIdentitiesOutput, isLast bool) bool {
			for _, identity := range page.CloudFrontOriginAccessIdentityList.Items {
				d.StreamListItem(ctx, identity)

				// Context can be cancelled due to manual cancellation or the limit has been hit
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

func getCloudFrontOriginAccessIdentity(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudFrontOriginAccessIdentity")

	// Create session
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	var identityID string
	if h.Item != nil {
		identityID = *h.Item.(*cloudfront.OriginAccessIdentitySummary).Id
	} else {
		identityID = d.KeyColumnQuals["id"].GetStringValue()
	}

	params := &cloudfront.GetCloudFrontOriginAccessIdentityInput{
		Id: aws.String(identityID),
	}

	op, err := svc.GetCloudFrontOriginAccessIdentity(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getCloudFrontOriginAccessIdentityARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudFrontOriginAccessIdentityARN")
	originAccessIdentityData := *originAccessIdentityID(h.Item)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":cloudfront::" + commonColumnData.AccountId + ":origin-access-identity/" + originAccessIdentityData

	return arn, nil
}

func originAccessIdentityID(item interface{}) *string {
	switch item := item.(type) {
	case *cloudfront.GetCloudFrontOriginAccessIdentityOutput:
		return item.CloudFrontOriginAccessIdentity.Id
	case *cloudfront.OriginAccessIdentitySummary:
		return item.Id
	}
	return nil
}
