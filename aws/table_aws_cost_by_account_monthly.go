package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCostByLinkedAccountMonthly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_account_monthly",
		Description: "AWS Cost Explorer - Cost by Linked Account (Monthly)",
		List: &plugin.ListConfig{
			Hydrate: listCostByLinkedAccountMonthly,
			Tags:    map[string]string{"service": "ce", "action": "GetCostAndUsage"},
		},
		Columns: awsGlobalRegionColumns(
			costExplorerColumns([]*plugin.Column{

				{
					Name:        "linked_account_id",
					Description: "The AWS Account ID.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1"),
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostByLinkedAccountMonthly(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildCostByLinkedAccountInput("MONTHLY")

	return streamCostAndUsage(ctx, d, params)
}

func buildCostByLinkedAccountInput(granularity string) *costexplorer.GetCostAndUsageInput {
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := getCEStartDateForGranularity(granularity).Format(timeFormat)

	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.Granularity(granularity),
		Metrics:     AllCostMetrics(),
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionType("DIMENSION"),
				Key:  aws.String("LINKED_ACCOUNT"),
			},
		},
	}

	return params
}
