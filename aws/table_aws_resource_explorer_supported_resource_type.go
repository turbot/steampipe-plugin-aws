package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAWSResourceExplorerSupportedResourceType(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_resource_explorer_supported_resource_type",
		Description: "AWS Resource Explorer Supported Resource Type",
		List: &plugin.ListConfig{
			Hydrate: listAWSExplorerSupportedTypes,
			Tags:    map[string]string{"service": "resource-explorer-2", "action": "ListSupportedResourceTypes"},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "resource_type",
				Description: "The unique identifier of the resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service",
				Description: "The Amazon Web Service that is associated with the resource type. This is the primary service that lets you create and interact with resources of this type.",
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

//// LIST FUNCTION

func listAWSExplorerSupportedTypes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	region, err := getDefaultRegion(ctx, d, h)
	if err != nil {
		return nil, err
	}

	// Create Session
	svc, err := ResourceExplorerClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_resource_explorer_supported_resource_type.listAWSExplorerSupportedTypes", "connnection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	paginator := resourceexplorer2.NewListSupportedResourceTypesPaginator(svc, &resourceexplorer2.ListSupportedResourceTypesInput{}, func(o *resourceexplorer2.ListSupportedResourceTypesPaginatorOptions) {
		o.Limit = 100
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {

		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_resource_explorer_supported_resource_type.listAWSExplorerSupportedTypes", "api_error", err)
			return nil, err
		}

		for _, resourceType := range output.ResourceTypes {
			d.StreamListItem(ctx, resourceType)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
