package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCostAllocationTags(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_allocation_tags",
		Description: "AWS Cost Allocation Tags",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:       "status",
					Operators:  []string{"="},
					Require:    plugin.Optional,
				},
			},
			Hydrate: listCostAllocationTags,
			Tags:    map[string]string{"service": "ce", "action": "ListCostAllocationTags"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "tag_key",
				Description: "The tag key for the cost allocation tag.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TagKey"),
			},
			{
				Name:        "tag_type",
				Description: "The type of the cost allocation tag (AWSGenerated or UserDefined).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TagType"),
			},
			{
				Name:        "status",
				Description: "The status of the cost allocation tag (Active or Inactive).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status"),
			},
			{
				Name:        "last_updated_date",
				Description: "The last updated date of the cost allocation tag.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LastUpdatedDate"),
			},
			{
				Name:        "last_used_date",
				Description: "The last used date of the cost allocation tag.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LastUsedDate"),
			},
			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TagKey"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCostAllocationTags(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cost_allocation_tags.listCostAllocationTags", "connection_error", err)
		return nil, err
	}

	// Get the status filter if provided
	var statusFilter *types.TagStatus
	if d.EqualsQualString("status") != "" {
		status := d.EqualsQualString("status")
		statusFilter = (*types.TagStatus)(&status)
	}

	input := &costexplorer.ListCostAllocationTagsInput{
		Status: statusFilter,
	}

	paginator := costexplorer.NewListCostAllocationTagsPaginator(client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ce_cost_allocation_tags.listCECostAllocationTags", "api_error", err)
			return nil, err
		}

		for _, tag := range output.CostAllocationTags {
			d.StreamListItem(ctx, tag)

			// Context may get cancelled
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

