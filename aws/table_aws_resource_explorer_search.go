package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAWSResourceExplorerSearch(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_resource_explorer_search",
		Description: "AWS Resource Explorer Search",
		List: &plugin.ListConfig{
			Hydrate: awsResourceExplorerSearch,
			IgnoreConfig: &plugin.IgnoreConfig{
				// UnauthorizedException error thrown for below cases in Resource Explorer
				// 1. Default view is not present in the region queried
				// 2. Credentials doesn't have access to the view used for searching
				// 3. Cross-account or cross-region view is used for searching
				// ValidationException error thrown for below cases in Resource Explorer
				// 1. If the `query` uses a filter properties that aren't included in the view.
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UnauthorizedException", "ValidationException"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "query", Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "view_arn", Require: plugin.Optional, CacheMatch: "exact"}, // The view to be used to search resources..
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon resource name (ARN) of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service",
				Description: "The Amazon Web Service that owns the resource and is responsible for creating and updating it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owning_account_id",
				Description: "The Amazon Web Services account that owns the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_reported_at",
				Description: "The date and time that Resource Explorer last queried this resource and updated the index with the latest information about the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "properties",
				Description: "Additional type-specific details about the resource.",
				Type:        proto.ColumnType_JSON,
			},
			// Inputs to the table
			{
				Name:        "view_arn",
				Description: "The Amazon resource name (ARN) of the view that this table uses to perform the search.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "query",
				Description: "A string that includes keywords and filters that specify the resources to include in the search results.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("query"),
			},
			{
				Name:        "region",
				Description: "The AWS Region in which the resource was created and exists.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION

/*
In an effort to make this table a bit easier, I think we should update the search table to work as follows:

The table should have no required key columns, and we remove region as a key column entirely (but keep query and view_arn

1. Queries for this table that do not specify any key columns, e.g., select * from aws_resource_explorer_search, should still work and follow this sequence

2. List all indexes and look for an aggregator index (connect to the default region to get these)

  - If there's no aggregator index in the account, return an error letting the user know there's no aggregator
  - Connect to the region where the aggregator index exists
  - Look for a default view in the aggregator region
  - If there's no default view, return an error letting the user know there's no default view
    Note: If this isn't possible due to Resource explorer ListViews always erroring out with error ValidationException aws/aws-sdk-go-v2#1916, we should skip this step for now and just make the search call
  - Make the search API call

This design has been largely modeled off of how AWS handles activating unified search:
  - To enable including your account's resources in the search results for unified search from any Amazon console, you must complete the following steps:
  - Activate Amazon Resource Explorer in one or more Amazon Web Services Regions in your account.
  - Register one Region to contain the aggregator index.
  - Create a default view in the Region with the aggregator index.
  - If a user does pass a view_arn, then we should extract the region out of the view ARN, connect to that region, and pass that view_arn in the input params.

Other notes:
  - AWS has a 500 requests/month limit on single region indexes and 10,000 requests/month limit on aggregator region indexes, so we strongly prefer searching on the aggregator indexes.
  - AWS in general returns 401 UnauthorizedException error codes for any error, e.g., Resource Explorer disabled, no index, no default view, insufficient IAM permissions, so we should do our best in the code to wrap these errors and provide better error messages when possible.
*/
func awsResourceExplorerSearch(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	searchParams := &resourceexplorer2.SearchInput{
		QueryString: aws.String(""),
	}

	region := getDefaultAwsRegion(d)
	hasViewARN := false
	if d.KeyColumnQuals["view_arn"] != nil {
		hasViewARN = true
		viewARN := d.KeyColumnQualString("view_arn")
		if arn.IsARN(viewARN) {
			arnData, _ := arn.Parse(viewARN)
			// API throws UnauthorizedException for cross-region and cross-account queries
			region = arnData.Region

			// Avoid cross-account queriying
			commonData, err := getCommonColumns(ctx, d, h)
			if err != nil {
				plugin.Logger(ctx).Error("aws_resource_explorer_view.awsResourceExplorerSearchs", "common_data_error", err)
				return nil, err
			}
			if arnData.AccountID != commonData.(*awsCommonColumnData).AccountId {
				return nil, nil
			}
		}
		searchParams.ViewArn = aws.String(d.KeyColumnQualString("view_arn"))
	}

	svc, err := ResourceExplorerClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_resource_explorer_view.awsResourceExplorerSearch", "connnection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Make an API Error to get aggregator index
	if !hasViewARN {
		getAggregatorIndexParams := &resourceexplorer2.ListIndexesInput{
			Type: types.IndexTypeAggregator,
		}

		indexesOutput, err := svc.ListIndexes(ctx, getAggregatorIndexParams)
		if err != nil {
			plugin.Logger(ctx).Error("aws_resource_explorer_view.awsResourceExplorerSearch", "list_indexes_api_error", err)
			return nil, err
		}

		if len(indexesOutput.Indexes) == 0 {
			return nil, fmt.Errorf("AGGREGATOR index is not found in found in the Account. Please use \"view_arn\" to serach resources or create aggreagtor index in account with default view and try again.")
		}

		region = *indexesOutput.Indexes[0].Region
	}

	// Create the service connection again with the aggregator region
	svc, err = ResourceExplorerClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_resource_explorer_view.awsResourceExplorerSearch", "search_api_connnection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	//
	if d.KeyColumnQuals["query"] != nil {
		searchParams.QueryString = aws.String(d.KeyColumnQualString("query"))
	}

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

	paginator := resourceexplorer2.NewSearchPaginator(svc, searchParams, func(o *resourceexplorer2.SearchPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_resource_explorer_view.awsResourceExplorerSearch", "serach_api_error", err)
			return nil, err
		}

		for _, resource := range output.Resources {
			d.StreamListItem(ctx, SearchStreamItem{resource, output.ViewArn})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

type SearchStreamItem struct {
	types.Resource
	ViewArn *string
}
