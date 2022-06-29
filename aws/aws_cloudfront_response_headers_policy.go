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
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchFunctionExists"}),
			},
			Hydrate: getCloudFrontResponseHeadersPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudFrontResponseHeadersPolicy,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The identifier for the response headers policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResponseHeadersPolicy.Id"),
			},
			{
				Name:        "last_modified_time",
				Description: "The date and time when the response headers policy was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ResponseHeadersPolicy.LastModifiedTime"),
			},
			{
				Name:        "response_headers_policy_config",
				Description: "A response headers policy contains information about a set of HTTP response headers and their values. CloudFront adds the headers in the policy to HTTP responses that it sends for requests that match a cache behavior that’s associated with the policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResponseHeadersPolicy.ResponseHeadersPolicyConfig"),
			},

			// {
			// 	Name:        "arn",
			// 	Description: "The version identifier for the current version of the CloudFront function.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromField("FunctionMetadata.FunctionARN", "FunctionSummary.FunctionMetadata.FunctionARN"),
			// },
			// {
			// 	Name:        "status",
			// 	Description: "The status of the CloudFront function.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromField("Status", "FunctionSummary.Status"),
			// 	Hydrate:     getCloudFrontFunction,
			// },
			// {
			// 	Name:        "e_tag",
			// 	Description: "The version identifier for the current version of the CloudFront function.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromField("ETag", "FunctionSummary.ETag"),
			// 	Hydrate:     getCloudFrontFunction,
			// },
			// {
			// 	Name:        "function_config",
			// 	Description: "Contains configuration information about a CloudFront function.",
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromField("FunctionConfig", "FunctionSummary.FunctionConfig"),
			// 	Hydrate:     getCloudFrontFunction,
			// },
			// {
			// 	Name:        "function_metadata",
			// 	Description: "Contains metadata about a CloudFront function.",
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromField("FunctionMetadata", "FunctionSummary.FunctionMetadata"),
			// },
			// // Steampipe standard columns
			// {
			// 	Name:        "title",
			// 	Description: resourceInterfaceDescription("title"),
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromField("Name", "FunctionSummary.Name"),
			// },
			// {
			// 	Name:        "akas",
			// 	Description: resourceInterfaceDescription("akas"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromField("FunctionMetadata.FunctionARN", "FunctionSummary.FunctionMetadata.FunctionARN").Transform(transform.EnsureStringArray),
			// },
		}),
	}
}

//// LIST FUNCTION

func listCloudFrontResponseHeadersPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableAwsCloudFrontFunction.listCloudFrontResponseHeadersPolicy")

	// Create Session
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Set up the limit
	input := cloudfront.ListResponseHeadersPoliciesInput{
		// TODO: Restore
		// MaxItems: aws.Int64(100),
		MaxItems: aws.Int64(1),
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

func getCloudFrontResponseHeadersPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableAwsCloudFrontFunction.getCloudFrontResponseHeadersPolicy")
	var id string

	if h.Item != nil {
		summary := h.Item.(*cloudfront.ResponseHeadersPolicySummary)
		id = *summary.ResponseHeadersPolicy.Id
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

	return data, nil
}

//// TRANSFORM FUNCTION
