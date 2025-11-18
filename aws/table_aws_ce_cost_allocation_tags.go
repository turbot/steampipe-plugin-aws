package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCECostAllocationTags(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ce_cost_allocation_tags",
		Description: "AWS Cost Explorer Cost Allocation Tags",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:       "status",
					Operators:  []string{"="},
					Require:    plugin.Optional,
				},
				{
					Name:       "type",
					Operators:  []string{"="},
					Require:    plugin.Optional,
				},
			},
			Hydrate: listCECostAllocationTags,
			Tags:    map[string]string{"service": "ce", "action": "ListCostAllocationTags"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "tag_key",
				Description: "The cost allocation tag key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the cost allocation tag (Active or Inactive).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the cost allocation tag (AWSGenerated or UserDefined).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_date",
				Description: "The last date that the tag was activated or deactivated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_used_date",
				Description: "The last month that the tag was used on an AWS resource.",
				Type:        proto.ColumnType_STRING,
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

func listCECostAllocationTags(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ce_cost_allocation_tags.listCECostAllocationTags", "connection_error", err)
		return nil, err
	}

	input := &costexplorer.ListCostAllocationTagsInput{
		MaxResults: aws.Int32(100),
	}

	// Add optional status filter
	if status := d.EqualsQualString("status"); status != "" {
		input.Status = types.CostAllocationTagStatus(status)
	}

	// Add optional type filter
	if tagType := d.EqualsQualString("type"); tagType != "" {
		input.Type = types.CostAllocationTagType(tagType)
	}

	// Handle pagination
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
