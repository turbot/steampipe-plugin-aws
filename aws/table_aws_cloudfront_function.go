package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudFrontFunction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudfront_function",
		Description: "AWS CloudFront Function",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchFunctionExists"}),
			},
			Hydrate: getCloudFrontFunction,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchFunctions,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	plugin.Logger(ctx).Trace("tableAwsCloudFrontFunction.listCloudWatchFunctions")

	// Create Session
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Set up the limit
	input := cloudfront.ListFunctionsInput{
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

	// Prepare for the slice of functions
	pagesLeft := true
	for pagesLeft {
		// List CloudFront functions
		data, err := svc.ListFunctions(&input)
		if err != nil {
			plugin.Logger(ctx).Error("ListFunctions", "ERROR", err)
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
	plugin.Logger(ctx).Trace("tableAwsCloudFrontFunction.getCloudFrontFunction")

	var name string

	if h.Item != nil {
		function_summary := h.Item.(*cloudfront.FunctionSummary)
		name = *function_summary.Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	if name == "" {
		plugin.Logger(ctx).Trace("Name is null, ignoring")
		return nil, nil
	}

	// Create service
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := cloudfront.DescribeFunctionInput{
		Name: aws.String(name),
	}

	// Get call
	data, err := svc.DescribeFunction(&params)

	if err != nil {
		plugin.Logger(ctx).Error("DescribeFunction", "ERROR", err)
		return nil, err
	}

	return data, nil
}

//// TRANSFORM FUNCTION
