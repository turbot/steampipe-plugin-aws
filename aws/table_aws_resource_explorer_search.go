package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAWSResourceExplorerSearch(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_resource_explorer_search",
		Description: "AWS Resource Explorer Search",
		List: &plugin.ListConfig{
			Hydrate: awsResourceExplorerSearch,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "query", Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "view_arn", Require: plugin.Optional, CacheMatch: "exact"},
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

func awsResourceExplorerSearch(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	region, err := getDefaultRegion(ctx, d, h)
	if err != nil {
		return nil, err
	}

	searchParams := &resourceexplorer2.SearchInput{
		QueryString: aws.String(""),
	}

	// If a view ARN is passed in, check if it's from a different account to
	// avoid an unsuccessful API call
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_resource_explorer_search.awsResourceExplorerSearchs", "common_data_error", err)
		return nil, err
	}
	accountID := commonData.(*awsCommonColumnData).AccountId

	hasViewARN := false
	if d.EqualsQuals["view_arn"] != nil {
		hasViewARN = true
		viewARN := d.EqualsQualString("view_arn")
		if arn.IsARN(viewARN) {
			arnData, _ := arn.Parse(viewARN)
			// API throws UnauthorizedException for cross-region and cross-account
			// queries
			region = arnData.Region

			// Avoid cross-account queries
			if arnData.AccountID != accountID {
				return nil, nil
			}
		}
		searchParams.ViewArn = aws.String(d.EqualsQualString("view_arn"))
	}

	svc, err := ResourceExplorerClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_resource_explorer_search.awsResourceExplorerSearch", "connnection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// If no view ARN is specified, attempt to find the account's aggregator
	// index and verify that region has a default view
	if !hasViewARN {
		// Since each account can only have 1 aggregator index, no need to page
		getAggregatorIndexParams := &resourceexplorer2.ListIndexesInput{
			Type: types.IndexTypeAggregator,
		}

		indexesOutput, err := svc.ListIndexes(ctx, getAggregatorIndexParams)
		if err != nil {
			plugin.Logger(ctx).Error("aws_resource_explorer_search.awsResourceExplorerSearch", "list_indexes_api_error", err)
			return nil, err
		}

		if len(indexesOutput.Indexes) == 0 {
			return nil, fmt.Errorf("Aggregator index not found in account %s. Please create an aggregator index or specify \"view_arn\".", accountID)
		}

		// Each account can only have 1 aggregator index
		region = *indexesOutput.Indexes[0].Region

		// Create the service connection again in the aggregator index's region
		svc, err = ResourceExplorerClient(ctx, d, region)
		if err != nil {
			plugin.Logger(ctx).Error("aws_resource_explorer_search.awsResourceExplorerSearch", "search_api_connnection_error", err)
			return nil, err
		}
		if svc == nil {
			// Unsupported region, return no data
			return nil, nil
		}

		defaultViewOutput, err := svc.GetDefaultView(ctx, &resourceexplorer2.GetDefaultViewInput{})
		if err != nil {
			plugin.Logger(ctx).Error("aws_resource_explorer_search.awsResourceExplorerSearch", "get_default_view_api_error", err)
			return nil, err
		}

		if defaultViewOutput.ViewArn == nil {
			return nil, fmt.Errorf("Default view not found in %s region in account %s. Please create a default view or specify \"view_arn\".", region, accountID)
		}
	}

	if d.EqualsQuals["query"] != nil {
		searchParams.QueryString = aws.String(d.EqualsQualString("query"))
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
			plugin.Logger(ctx).Error("aws_resource_explorer_search.awsResourceExplorerSearch", "search_api_error", err)
			return nil, err
		}

		for _, resource := range output.Resources {
			d.StreamListItem(ctx, SearchStreamItem{resource, output.ViewArn})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
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
