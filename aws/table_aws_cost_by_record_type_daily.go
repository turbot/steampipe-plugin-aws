package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCostByRecordTypeDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_record_type_daily",
		Description: "AWS Cost Explorer - Cost by Record Type (Daily)",
		List: &plugin.ListConfig{
			Hydrate: listCostByRecordTypeDaily,
			Tags:    map[string]string{"service": "ce", "action": "GetCostAndUsage"},
		},
		Columns: awsGlobalRegionColumns(
			costExplorerColumns([]*plugin.Column{

				{
					Name:        "linked_account_id",
					Description: "The linked AWS Account ID.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1"),
				},
				{
					Name:        "record_type",
					Description: "The different types of charges such as RI fees, usage, costs, tax refunds, and credits.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension2"),
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostByRecordTypeDaily(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildCostByRecordTypeInput("DAILY")
	return streamCostAndUsage(ctx, d, params)
}
