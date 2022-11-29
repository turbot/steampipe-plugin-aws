package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudFrontResponseHeadersPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudfront_response_headers_policy",
		Description: "AWS Cloudfront Response Headers Policy",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "type",
					Require: plugin.Optional,
				},
			},
			Hydrate: listCloudFrontResponseHeadersPolicies,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the response headers policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResponseHeadersPolicy.ResponseHeadersPolicyConfig.Name"),
			},
			{
				Name:        "id",
				Description: "The identifier for the response headers policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResponseHeadersPolicy.Id"),
			},
			{
				Name:        "arn",
				Description: "The version identifier for the current version of the response headers policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAccountARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "last_modified_time",
				Description: "The date and time when the response headers policy was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ResponseHeadersPolicy.LastModifiedTime"),
			},
			{
				Name:        "type",
				Description: "The type of response headers policy, either managed (created by AWS) or custom (created in this AWS account).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "The version identifier for the current version of the response headers policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getETagValue,
				Transform:   transform.FromField("ETag"),
			},
			{
				Name:        "response_headers_policy_config",
				Description: "A response headers policy contains information about a set of HTTP response headers and their values. CloudFront adds the headers in the policy to HTTP responses that it sends for requests that match a cache behavior that’s associated with the policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResponseHeadersPolicy.ResponseHeadersPolicyConfig"),
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResponseHeadersPolicy.ResponseHeadersPolicyConfig.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAccountARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudFrontResponseHeadersPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_response_headers_policy.listCloudFrontResponseHeadersPolicies", "client_error", err)
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

	input := &cloudfront.ListResponseHeadersPoliciesInput{
		MaxItems: &maxItems,
	}

	// Additonal Filter
	policyTypeColumn := d.KeyColumnQuals["type"]
	if policyTypeColumn != nil {
		input.Type = types.ResponseHeadersPolicyType(policyTypeColumn.GetStringValue())
	}

	// Paginator not avilable for the API
	pagesLeft := true
	for pagesLeft {
		data, err := svc.ListResponseHeadersPolicies(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudfront_response_headers_policy.listCloudFrontResponseHeadersPolicies", "api_error", err)
			return nil, err
		}

		for _, policy := range data.ResponseHeadersPolicyList.Items {
			d.StreamListItem(ctx, policy)
		}

		if data.ResponseHeadersPolicyList.NextMarker != nil {
			pagesLeft = true
			input.Marker = data.ResponseHeadersPolicyList.NextMarker
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getETagValue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(types.ResponseHeadersPolicySummary)
	id := *item.ResponseHeadersPolicy.Id

	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_response_headers_policy.getETagValue", "client_error", err)
		return nil, err
	}

	// Build the params
	params := &cloudfront.GetResponseHeadersPolicyInput{Id: &id}

	// Get call
	data, err := svc.GetResponseHeadersPolicy(ctx, params)

	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_response_headers_policy.getETagValue", "api_error", err)
		return nil, err
	}

	return data, nil
}

func getAccountARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get common columns which will be used to create the ARN
	response, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_response_headers_policy.getAccountARN", "common_data_error", err)
		return nil, err
	}

	commonColumnData := response.(*awsCommonColumnData)

	var id string

	item := h.Item.(types.ResponseHeadersPolicySummary)
	id = *item.ResponseHeadersPolicy.Id

	arn := "arn:" + commonColumnData.Partition + ":cloudfront::" + commonColumnData.AccountId + ":response-headers-policy/" + id

	return arn, nil
}
