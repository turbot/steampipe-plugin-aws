package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCostByServiceDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_service_daily",
		Description: "AWS Cost Explorer - Cost by Service (Daily)",
		List: &plugin.ListConfig{
			Hydrate: listCostByServiceDaily,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "service", Operators: []string{"=", "<>"}, Require: plugin.Optional},
			},
		},
		Columns: awsGlobalRegionColumns(
			costExplorerColumns([]*plugin.Column{

				{
					Name:        "service",
					Description: "The name of the AWS service.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1"),
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostByServiceDaily(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildCostByServiceInput("DAILY", d)
	return streamCostAndUsage(ctx, d, params)
}
