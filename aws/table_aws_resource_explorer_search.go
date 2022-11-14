package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2/types"
	"github.com/aws/aws-sdk-go/aws"
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
			// ParentHydrate: listAWSExplorerIndexes,
			Hydrate: awsExplorerSearch,
			// IgnoreConfig: &plugin.IgnoreConfig{
			// 	// ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "UnauthorizedException"}),
			// },
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "region_input", Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "query", Require: plugin.Optional, CacheMatch: "exact"},
			},
		},
		// GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "view_arn",
				Description: "The Amazon resource name (ARN) of the view that this operation used to perform the search.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The list of views available in the Amazon Web Services Region.",
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
				Name:        "region_input",
				Description: "The Amazon Web Services Region in which the index exists.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("region_input"),
			},
			{
				Name:        "region_output",
				Description: "The Amazon Web Services Region in which the index exists.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region"),
			},
			{
				Name:        "query",
				Description: "The Amazon Web Services Region in which the index exists.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("query"),
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
		},
	}
}

//// LIST FUNCTION

func awsExplorerSearch(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString("region_input")
	// query := d.KeyColumnQualString("query")
	svc, err := ResourceExplorerRegionalClient(ctx, d, region)
	// svc, err := ResourceExplorerRegionalClient(ctx, d, "ap-south-1")
	if err != nil {
		plugin.Logger(ctx).Error("aws_resource_explorer_view.listAWSExplorerViews", "connnection_error", err)
		return nil, err
	}

	params := &resourceexplorer2.SearchInput{
		// QueryString: aws.String(query),
		// QueryString: aws.String("region:us-east-1"),
		QueryString: aws.String(""),
	}

	paginator := resourceexplorer2.NewSearchPaginator(svc, params, func(o *resourceexplorer2.SearchPaginatorOptions) {
		o.Limit = 100
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_resource_explorer_view.listAWSExplorerViews", "api_error", err)
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
