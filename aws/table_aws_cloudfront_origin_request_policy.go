package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudFrontOriginRequestPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudfront_origin_request_policy",
		Description: "AWS CloudFront Origin Request Policy",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"NoSuchOriginRequestPolicy", "InvalidParameter"}),
			Hydrate:           getCloudFrontOriginRequestPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudFrontOriginRequestPolicies,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "A unique name to identify the origin request policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "OriginRequestPolicy.OriginRequestPolicyConfig.Name"),
			},
			{
				Name:        "id",
				Description: "The ID for the origin request policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id", "OriginRequestPolicy.Id"),
			},
			{
				Name:        "comment",
				Description: "The comment for this origin request policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Comment", "OriginRequestPolicy.OriginRequestPolicyConfig.Comment"),
			},
			{
				Name:        "etag",
				Description: "The current version of the origin request policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFrontOriginRequestPolicy,
				Transform:   transform.FromField("ETag"),
			},
			{
				Name:        "last_modified_time",
				Description: "The date and time when the origin request policy was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastModifiedTime", "OriginRequestPolicy.LastModifiedTime"),
			},
			{
				Name:        "cookies_config",
				Description: "The cookies from viewer requests to include in origin requests.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CookiesConfig", "OriginRequestPolicy.OriginRequestPolicyConfig.CookiesConfig"),
			},
			{
				Name:        "headers_config",
				Description: "The HTTP headers to include in origin requests.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HeadersConfig", "OriginRequestPolicy.OriginRequestPolicyConfig.HeadersConfig"),
			},
			{
				Name:        "query_strings_config",
				Description: "The URL query strings from viewer requests to include in origin requests.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("QueryStringsConfig", "OriginRequestPolicy.OriginRequestPolicyConfig.QueryStringsConfig"),
			},

			//  Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "OriginRequestPolicy.OriginRequestPolicyConfig.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFrontOriginRequestPolicyAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudFrontOriginRequestPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	plugin.Logger(ctx).Trace("listCloudFrontOriginRequestPolicies")

	// Create session
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	params := &cloudfront.ListOriginRequestPoliciesInput{
		MaxItems: aws.Int64(1000),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxItems {
			if *limit < 1 {
				params.MaxItems = types.Int64(1)
			} else {
				params.MaxItems = limit
			}
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := svc.ListOriginRequestPolicies(params)
		if err != nil {
			return nil, err
		}
		for _, policy := range response.OriginRequestPolicyList.Items {
			d.StreamListItem(ctx, policy)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if response.OriginRequestPolicyList.NextMarker != nil {
			pagesLeft = true
			params.Marker = response.OriginRequestPolicyList.NextMarker
		} else {
			pagesLeft = false
		}

	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudFrontOriginRequestPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudFrontOriginRequestPolicy")

	// Create session
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	var policyID string
	if h.Item != nil {
		policyID = *h.Item.(*cloudfront.OriginRequestPolicySummary).OriginRequestPolicy.Id
	} else {
		policyID = d.KeyColumnQuals["id"].GetStringValue()
	}

	params := &cloudfront.GetOriginRequestPolicyInput{
		Id: aws.String(policyID),
	}

	op, err := svc.GetOriginRequestPolicy(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getCloudFrontOriginRequestPolicyAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudFrontOriginRequestPolicyAkas")
	policyID := *originRequestPolicyID(h.Item)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":cloudfront::" + commonColumnData.AccountId + ":origin-request-policy/" + policyID

	//return arn, nil
	return []string{aka}, nil
}

func originRequestPolicyID(item interface{}) *string {
	switch item := item.(type) {
	case *cloudfront.GetOriginRequestPolicyOutput:
		return item.OriginRequestPolicy.Id
	case *cloudfront.OriginRequestPolicySummary:
		return item.OriginRequestPolicy.Id
	}
	return nil
}
