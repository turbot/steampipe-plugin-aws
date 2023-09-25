package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudFrontFunction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudfront_function",
		Description: "AWS CloudFront Function",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchFunctionExists"}),
			},
			Hydrate: getCloudFrontFunction,
			Tags:    map[string]string{"service": "cloudfront", "action": "GetFunction"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchFunctions,
			Tags:    map[string]string{"service": "cloudfront", "action": "ListFunctions"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCloudFrontFunction,
				Tags: map[string]string{"service": "cloudfront", "action": "GetFunction"},
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the CloudFront function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "FunctionSummary.Name"),
			},
			{
				Name:        "arn",
				Description: "The version identifier for the current version of the CloudFront function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FunctionMetadata.FunctionARN", "FunctionSummary.FunctionMetadata.FunctionARN"),
			},
			{
				Name:        "status",
				Description: "The status of the CloudFront function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status", "FunctionSummary.Status"),
				Hydrate:     getCloudFrontFunction,
			},
			{
				Name:        "e_tag",
				Description: "The version identifier for the current version of the CloudFront function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag", "FunctionSummary.ETag"),
				Hydrate:     getCloudFrontFunction,
			},
			{
				Name:        "function_config",
				Description: "Contains configuration information about a CloudFront function.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FunctionConfig", "FunctionSummary.FunctionConfig"),
				Hydrate:     getCloudFrontFunction,
			},
			{
				Name:        "function_metadata",
				Description: "Contains metadata about a CloudFront function.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FunctionMetadata", "FunctionSummary.FunctionMetadata"),
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "FunctionSummary.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FunctionMetadata.FunctionARN", "FunctionSummary.FunctionMetadata.FunctionARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudWatchFunctions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_function.listCloudWatchFunctions", "client_error", err)
		return nil, err
	}

	maxItems := int32(100)

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

	input := &cloudfront.ListFunctionsInput{
		MaxItems: &maxItems,
	}

	// Paginator not available for the API
	pagesLeft := true
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		// List CloudFront functions
		data, err := svc.ListFunctions(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudfront_function.listCloudWatchFunctions", "api_error", err)
			return nil, err
		}

		for _, function := range data.FunctionList.Items {
			d.StreamListItem(ctx, function)
		}

		if data.FunctionList.NextMarker != nil {
			pagesLeft = true
			input.Marker = data.FunctionList.NextMarker
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudFrontFunction(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var name string

	if h.Item != nil {
		function_summary := h.Item.(types.FunctionSummary)
		name = *function_summary.Name
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	if strings.TrimSpace(name) == "" {
		return nil, nil
	}

	// Create service
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_function.getCloudFrontFunction", "client_error", err)
		return nil, err
	}

	// Build the params
	params := &cloudfront.DescribeFunctionInput{
		Name: &name,
	}

	// Get call
	data, err := svc.DescribeFunction(ctx, params)

	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_function.getCloudFrontFunction", "api_error", err)
		return nil, err
	}

	return *data, nil
}

//// TRANSFORM FUNCTION
