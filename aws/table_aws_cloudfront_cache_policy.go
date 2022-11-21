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

func tableAwsCloudFrontCachePolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudfront_cache_policy",
		Description: "AWS CloudFront Cache Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchCachePolicy"}),
			},
			Hydrate: getCloudFrontCachePolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudFrontCachePolicies,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "A unique name to identify the cache policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CachePolicy.CachePolicyConfig.Name"),
			},
			{
				Name:        "id",
				Description: "The unique identifier for the cache policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CachePolicy.Id"),
			},
			{
				Name:        "comment",
				Description: "A comment to describe the cache policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CachePolicy.CachePolicyConfig.Comment"),
			},
			{
				Name:        "default_ttl",
				Description: "The default amount of time, in seconds, that you want objects to stay in the CloudFront cache before CloudFront sends another request to the origin to see if the object has been updated.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CachePolicy.CachePolicyConfig.DefaultTTL"),
			},
			{
				Name:        "etag",
				Description: "The current version of the cache policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFrontCachePolicy,
				Transform:   transform.FromField("ETag"),
			},
			{
				Name:        "max_ttl",
				Description: "The maximum amount of time, in seconds, that you want objects to stay in the CloudFront cache before CloudFront sends another request to the origin to see if the object has been updated.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CachePolicy.CachePolicyConfig.MaxTTL"),
			},
			{
				Name:        "min_ttl",
				Description: "The minimum amount of time, in seconds, that you want objects to stay in the CloudFront cache before CloudFront sends another request to the origin to see if the object has been updated.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CachePolicy.CachePolicyConfig.MinTTL"),
			},
			{
				Name:        "last_modified_time",
				Description: "The date and time when the cache policy was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CachePolicy.LastModifiedTime"),
			},
			{
				Name:        "parameters_in_cache_key_and_forwarded_to_origin",
				Description: "The HTTP headers, cookies, and URL query strings to include in the cache key. The values included in the cache key are automatically included in requests that CloudFront sends to the origin.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CachePolicy.CachePolicyConfig.ParametersInCacheKeyAndForwardedToOrigin"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CachePolicy.CachePolicyConfig.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudfrontCachePolicyAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudFrontCachePolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_cache_policy.listCloudFrontCachePolicies", "client_error", err)
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

	input := &cloudfront.ListCachePoliciesInput{
		MaxItems: &maxItems,
	}

	// Paginator not avilable for API ListCachePolicies
	pagesLeft := true
	for pagesLeft {
		result, err := svc.ListCachePolicies(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudfront_cache_policy.listCloudFrontCachePolicies", "api_error", err)
			return nil, err
		}
		for _, policy := range result.CachePolicyList.Items {
			d.StreamListItem(ctx, policy)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if result.CachePolicyList.NextMarker != nil {
			input.Marker = result.CachePolicyList.NextMarker
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudFrontCachePolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_cache_policy.getCloudFrontCachePolicy", "client_error", err)
		return nil, err
	}

	var id string
	if h.Item != nil {
		id = *h.Item.(types.CachePolicySummary).CachePolicy.Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	params := &cloudfront.GetCachePolicyInput{
		Id: aws.String(id),
	}

	op, err := svc.GetCachePolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_cache_policy.getCloudFrontCachePolicy", "api_error", err)
		return nil, err
	}

	return *op, nil
}

func getCloudfrontCachePolicyAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := cloudFrontCachePolicyAka(h.Item)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_cache_policy.getCloudfrontCachePolicyAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":cloudfront::" + commonColumnData.AccountId + ":cache-policy/" + *id}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func cloudFrontCachePolicyAka(item interface{}) *string {
	switch item := item.(type) {
	case cloudfront.GetCachePolicyOutput:
		return item.CachePolicy.Id
	case types.CachePolicySummary:
		return item.CachePolicy.Id
	}
	return nil
}
