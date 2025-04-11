package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCostByResourceDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_resource_daily",
		Description: "AWS Cost Explorer - Cost by Resource (Daily)",
		List: &plugin.ListConfig{
			Hydrate: listCostByResourceDaily,
			Tags:    map[string]string{"service": "ce", "action": "GetCostAndUsageWithResources"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "resource_id", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "service", Operators: []string{"=", "<>"}, Require: plugin.Required},
			},
		},
		Columns: awsGlobalRegionColumns(
			costExplorerColumns([]*plugin.Column{
				{
					Name:        "resource_id",
					Description: "The unique identifier for the resource.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1"),
				},
				{
					Name:        "service",
					Description: "The name of the AWS service.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension2"),
				},
			}),
		),
	}
}

func listCostByResourceDaily(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCostByResourceDaily", "client_error", err)
		return nil, err
	}

	params := buildCostByResourceInput("DAILY", d)

	// List call
	for {
		output, err := svc.GetCostAndUsageWithResources(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("listCostByResourceDaily", "api_error", err)
			return nil, err
		}

		// Stream the results
		for _, row := range buildCEMetricRows(ctx, (*costexplorer.GetCostAndUsageOutput)(output), nil) {
			d.StreamListItem(ctx, row)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Get more pages if there are any
		if output.NextPageToken == nil {
			break
		}
		params.NextPageToken = output.NextPageToken
	}

	return nil, nil
}
