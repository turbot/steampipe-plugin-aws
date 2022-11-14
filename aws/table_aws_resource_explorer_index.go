package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAWSResourceExplorerIndex(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_resource_explorer_index",
		Description: "AWS Resource Explorer Index",
		List: &plugin.ListConfig{
			Hydrate: listAWSExplorerIndexes,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException"}),
			},
		},
		Columns: awsDefaultColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon resource name (ARN) of the index.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of index. It can be one of the following values: LOCAL, AGGREGATOR.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "The AWS Region in which the index exists.",
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

//// LIST FUNCTION

func listAWSExplorerIndexes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ResourceExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_resource_explorer_index.listAWSExplorerIndexes", "connnection_error", err)
		return nil, err
	}

	paginator := resourceexplorer2.NewListIndexesPaginator(svc, &resourceexplorer2.ListIndexesInput{}, func(o *resourceexplorer2.ListIndexesPaginatorOptions) {
		o.Limit = 100
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_resource_explorer_index.listAWSExplorerIndexes", "api_error", err)
			return nil, err
		}

		for _, index := range output.Indexes {
			d.StreamListItem(ctx, index)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
