package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
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
			Hydrate: listCloudFrontResponseHeadersPolicy,
		},
		GetMatrixItem: BuildRegionList,
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
				Transform:   transform.FromField("Type"),
			},
			{
				Name:        "e_tag",
				Description: "The version identifier for the current version of the response headers policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getETagValue,
				Transform:   transform.FromValue(),
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

func listCloudFrontResponseHeadersPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableAwsCloudFrontFunction.listCloudFrontResponseHeadersPolicy")

	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := cloudfront.ListResponseHeadersPoliciesInput{
		MaxItems: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxItems {
			if *limit < 1 {
				input.MaxItems = aws.Int64(1)
			} else {
				input.MaxItems = limit
			}
		}
	}

	// Additonal Filter
	policyTypeColumn := d.KeyColumnQuals["type"]
	if policyTypeColumn != nil {
		input.Type = aws.String(policyTypeColumn.GetStringValue())
	}

	pagesLeft := true
	for pagesLeft {
		data, err := svc.ListResponseHeadersPolicies(&input)
		if err != nil {
			plugin.Logger(ctx).Error("ListResponseHeadersPolicies", "ERROR", err)
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
	plugin.Logger(ctx).Trace("tableAwsCloudFrontFunction.getCloudFrontResponseHeadersPolicy")

	var id string

	if h.Item != nil {
		item := h.Item.(*cloudfront.ResponseHeadersPolicySummary)
		id = *item.ResponseHeadersPolicy.Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	if id == "" {
		plugin.Logger(ctx).Trace("Id is null, ignoring")
		return nil, nil
	}

	// Create service
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := cloudfront.GetResponseHeadersPolicyInput{
		Id: aws.String(id),
	}

	// Get call
	data, err := svc.GetResponseHeadersPolicy(&params)

	if err != nil {
		plugin.Logger(ctx).Error("GetResponseHeadersPolicy", "ERROR", err)
		return nil, err
	}

	return data.ETag, nil
}

func getAccountARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableAwsCloudFrontFunction.getAccountARN")

	// Get common columns which will be used to create the ARN
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	response, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := response.(*awsCommonColumnData)

	var id string

	if h.Item != nil {
		item := h.Item.(*cloudfront.ResponseHeadersPolicySummary)
		id = *item.ResponseHeadersPolicy.Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	arn := "arn:" + commonColumnData.Partition + ":cloudfront::" + commonColumnData.AccountId + ":response-headers-policy/" + id

	return arn, nil
}
