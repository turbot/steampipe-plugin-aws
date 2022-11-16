package aws

// import (
// 	"context"

// 	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
// 	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
// 	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
// 	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
// )

// //// TABLE DEFINITION

// func tableAWSResourceExplorerView(_ context.Context) *plugin.Table {
// 	return &plugin.Table{
// 		Name:        "aws_resource_explorer_view",
// 		Description: "AWS Resource Explorer View",
// 		List: &plugin.ListConfig{
// 			// ParentHydrate: listAWSExplorerIndexes,
// 			Hydrate: listAWSExplorerViews,
// 			// IgnoreConfig: &plugin.IgnoreConfig{
// 			// 	// ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException"}),
// 			// },
// 			KeyColumns: plugin.KeyColumnSlice{
// 				{Name: "region", Require: plugin.Optional, CacheMatch: "exact"},
// 			},
// 		},
// 		// GetMatrixItemFunc: BuildRegionList,
// 		Columns: []*plugin.Column{
// 			{
// 				Name:        "view",
// 				Description: "The list of views available in the Amazon Web Services Region.",
// 				Type:        proto.ColumnType_STRING,
// 				Transform:   transform.FromValue(),
// 			},
// 			{
// 				Name:        "region",
// 				Description: "The Amazon Web Services Region in which the index exists.",
// 				Type:        proto.ColumnType_STRING,
// 				Transform:   transform.FromQual("region"),
// 			},
// 			// {
// 			// 	Name:        "type",
// 			// 	Description: "The type of index. It can be one of the following values: LOCAL, AGGREGATOR.",
// 			// 	Type:        proto.ColumnType_STRING,
// 			// },
// 			// {
// 			// 	Name:        "partition",
// 			// 	Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
// 			// 	Type:        proto.ColumnType_STRING,
// 			// 	Hydrate:     getCommonColumns,
// 			// },
// 			// {
// 			// 	Name:        "account_id",
// 			// 	Description: "The AWS Account ID in which the resource is located.",
// 			// 	Type:        proto.ColumnType_STRING,
// 			// 	Hydrate:     getCommonColumns,
// 			// 	Transform:   transform.FromCamel(),
// 			// },
// 		},
// 	}
// }

// //// LIST FUNCTION

// func listAWSExplorerViews(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	// explorerIndex := h.Item.(types.Index)
// 	// Create Session
// 	region := d.KeyColumnQualString("region")
// 	svc, err := ResourceExplorerClient(ctx, d, region)
// 	if err != nil {
// 		plugin.Logger(ctx).Error("aws_resource_explorer_view.listAWSExplorerViews", "connnection_error", err)
// 		return nil, err
// 	}

// 	params := &resourceexplorer2.ListViewsInput{}

// 	paginator := resourceexplorer2.NewListViewsPaginator(svc, params, func(o *resourceexplorer2.ListViewsPaginatorOptions) {
// 		o.Limit = 100
// 		o.StopOnDuplicateToken = true
// 	})

// 	for paginator.HasMorePages() {
// 		output, err := paginator.NextPage(ctx)
// 		if err != nil {
// 			plugin.Logger(ctx).Error("aws_resource_explorer_view.listAWSExplorerViews", "api_error", err)
// 			return nil, err
// 		}

// 		for _, view := range output.Views {
// 			d.StreamListItem(ctx, view)

// 			// Context may get cancelled due to manual cancellation or if the limit has been reached
// 			if d.QueryStatus.RowsRemaining(ctx) == 0 {
// 				return nil, nil
// 			}
// 		}
// 	}

// 	return nil, nil
// }
