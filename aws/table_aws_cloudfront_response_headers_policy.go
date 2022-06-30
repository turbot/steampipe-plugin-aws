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
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				// TODO: Find not found error
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchResponseHeadersPolicy"}),
			},
			Hydrate: getCloudFrontResponseHeadersPolicy,
		},
		List: &plugin.ListConfig{
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
				Transform:   transform.From(getAccountARN),
			},
			{
				Name:        "e_tag",
				Description: "The version identifier for the current version of the response headers policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag"),
				Hydrate:     getCloudFrontResponseHeadersPolicy,
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
				Transform:   transform.From(getAccountARN).Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type responseHeadersPolicyItem struct {
	Partition             string
	AccountId             string
	ResponseHeadersPolicy cloudfront.ResponseHeadersPolicy
	Type                  *string
	ETag                  *string
}

//// LIST FUNCTION

func listCloudFrontResponseHeadersPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableAwsCloudFrontFunction.listCloudFrontResponseHeadersPolicy")

	// Create Session
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get common columns which will be used to create the ARN
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Set up the limit
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

	pagesLeft := true
	for pagesLeft {
		// List CloudFront Response Headers Policies
		data, err := svc.ListResponseHeadersPolicies(&input)
		if err != nil {
			plugin.Logger(ctx).Error("ListResponseHeadersPolicies", "ERROR", err)
			return nil, err
		}

		for _, policy := range data.ResponseHeadersPolicyList.Items {
			item := responseHeadersPolicyItem{
				AccountId:             commonColumnData.AccountId,
				Partition:             commonColumnData.Partition,
				ResponseHeadersPolicy: *policy.ResponseHeadersPolicy,
				Type:                  policy.Type,
			}
			d.StreamListItem(ctx, &item)
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

func getCloudFrontResponseHeadersPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableAwsCloudFrontFunction.getCloudFrontResponseHeadersPolicy")

	// Get common columns which will be used to create the ARN
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	var id string

	if h.Item != nil {
		item := h.Item.(*responseHeadersPolicyItem)
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

	item := responseHeadersPolicyItem{
		AccountId:             commonColumnData.AccountId,
		Partition:             commonColumnData.Partition,
		ResponseHeadersPolicy: *data.ResponseHeadersPolicy,
		ETag:                  data.ETag,
	}

	return &item, nil
}

//// TRANSFORM FUNCTION
func getAccountARN(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableAwsCloudFrontFunction.getAccountARN")
	item := d.HydrateItem.(*responseHeadersPolicyItem)

	arn := "arn:" + item.Partition + ":cloudfront::" + item.AccountId + ":response-headers-policy/" + *item.ResponseHeadersPolicy.Id

	return arn, nil
}
