package aws

import (
	"context"

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
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UnauthorizedException"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "query", Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "region", Require: plugin.AnyOf, CacheMatch: "exact"},
				{Name: "view_arn", Require: plugin.AnyOf, CacheMatch: "exact"}, // The view to be used to search resources..
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
				Name:        "resource_region",
				Description: "The AWS Region in which the resource was created and exists.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region"),
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
				Name:        "region",
				Description: "The Amazon Web Services Region to search for the resources.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("region"),
			},
			{
				Name:        "query",
				Description: "A string that includes keywords and filters that specify the resources to include in the search results.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("query"),
			},
		},
	}
}

//// LIST FUNCTION

func awsResourceExplorerSearch(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	params := &resourceexplorer2.SearchInput{
		QueryString: aws.String(""),
	}

	region := ""
	if d.KeyColumnQuals["region"] != nil {
		region = d.KeyColumnQualString("region")
	}

	if d.KeyColumnQuals["view_arn"] != nil {
		viewARN := d.KeyColumnQualString("view_arn")
		if arn.IsARN(viewARN) {
			arnData, _ := arn.Parse(viewARN)
			// API throws UnauthorizedException for cross-region and cross-account queries
			// Avoid cross-region queriying
			if region != "" && arnData.Region != region {
				return nil, nil
			}
			if region == "" {
				region = arnData.Region
			}

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
		params.ViewArn = aws.String(d.KeyColumnQualString("view_arn"))
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

	if d.KeyColumnQuals["query"] != nil {
		params.QueryString = aws.String(d.KeyColumnQualString("query"))
		plugin.Logger(ctx).Info("aws_resource_explorer_view.awsResourceExplorerSearch", "QueryString", *params.QueryString)
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

	paginator := resourceexplorer2.NewSearchPaginator(svc, params, func(o *resourceexplorer2.SearchPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_resource_explorer_view.awsResourceExplorerSearch", "api_error", err)
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
