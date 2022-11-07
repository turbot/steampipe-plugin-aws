package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudFrontOriginRequestPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudfront_origin_request_policy",
		Description: "AWS CloudFront Origin Request Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchOriginRequestPolicy", "InvalidParameter"}),
			},
			Hydrate: getCloudFrontOriginRequestPolicy,
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
	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_origin_request_policy.listCloudFrontOriginRequestPolicies", "client_error", err)
		return nil, err
	}

	// The maximum number for MaxItems parameter is not defined by the API
	// We have set the MaxItems to 1000 based on our test
	maxItems := int32(1000)

	// Reduce the basic request limit down if the user has only requested a small number of rows
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

	params := &cloudfront.ListOriginRequestPoliciesInput{
		MaxItems: &maxItems,
	}

	pagesLeft := true
	for pagesLeft {
		response, err := svc.ListOriginRequestPolicies(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudfront_origin_request_policy.listCloudFrontOriginRequestPolicies", "api_error", err)
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
	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_origin_request_policy.getCloudFrontOriginRequestPolicy", "client_error", err)
		return nil, err
	}

	var policyID string
	if h.Item != nil {
		policyID = *h.Item.(types.OriginRequestPolicySummary).OriginRequestPolicy.Id
	} else {
		policyID = d.KeyColumnQuals["id"].GetStringValue()
	}

	if strings.TrimSpace(policyID) == "" {
		return nil, nil
	}

	params := &cloudfront.GetOriginRequestPolicyInput{
		Id: aws.String(policyID),
	}

	op, err := svc.GetOriginRequestPolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_origin_request_policy.getCloudFrontOriginRequestPolicy", "api_error", err)
		return nil, err
	}

	return *op, nil
}

func getCloudFrontOriginRequestPolicyAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
	case cloudfront.GetOriginRequestPolicyOutput:
		return item.OriginRequestPolicy.Id
	case types.OriginRequestPolicySummary:
		return item.OriginRequestPolicy.Id
	}
	return nil
}
