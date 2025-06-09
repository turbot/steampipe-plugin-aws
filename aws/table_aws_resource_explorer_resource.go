package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAwsResourceExplorerResource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_resource_explorer_resource",
		Description: "AWS Resource Explorer Resource provides information about resources across regions in your AWS account.",
		List: &plugin.ListConfig{
			Hydrate: listResourceExplorerResources,
			Tags:    map[string]string{"service": "resource-explorer-2", "action": "ListResources"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "view_arn", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "view_arn",
				Description: "The Amazon resource name (ARN) of the view that this operation used to perform the search.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_reported_at",
				Description: "The date and time that Resource Explorer last queried this resource and updated the index with the latest information about the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "owning_account_id",
				Description: "The AWS account that owns the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "filter",
				Description: "The string that contains the search keywords, prefixes, and operators to control the results that can be returned by a Search operation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("filter"),
			},
			{
				Name:        "properties",
				Description: "A structure with additional type-specific details about the resource. These properties can be added by turning on integration between Resource Explorer and other AWS services.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "region",
				Description: "The AWS Region in which the resource was created and exists.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service",
				Description: "The AWS service that owns the resource and is responsible for creating and updating it.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe common columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Name"),
			},
		},
	}
}

type resourceInfo struct {
	types.Resource
	ViewArn *string
}

func listResourceExplorerResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region, err := getDefaultRegion(ctx, d, h)
	if err != nil {
		return nil, err
	}

	// Create client
	svc, err := ResourceExplorerClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_resource_explorer_resource.listResourceExplorerResources", "client_error", err)
		return nil, err
	}

	filter := d.EqualsQualString("filter")

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	// Build the params
	input := &resourceexplorer2.ListResourcesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if filter != "" {
		input.Filters = &types.SearchFilter{
			FilterString: aws.String(filter),
		}
	}

	viewArn := d.EqualsQualString("view_arn")
	if viewArn != "" {
		input.ViewArn = aws.String(viewArn)
	}

	// Get call
	paginator := resourceexplorer2.NewListResourcesPaginator(svc, input)
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_resource_explorer_resource.listResourceExplorerResources", "api_error", err)
			return nil, err
		}

		for _, resource := range output.Resources {
			d.StreamListItem(ctx, &resourceInfo{resource, output.ViewArn})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
