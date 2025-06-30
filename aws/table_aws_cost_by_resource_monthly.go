package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAwsCostByResourceMonthly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_resource_monthly",
		Description: "AWS Cost Explorer - Cost by Resource (Monthly)",
		List: &plugin.ListConfig{
			Hydrate: listCostByResourceMonthly,
			Tags:    map[string]string{"service": "ce", "action": "GetCostAndUsageWithResources"},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:      "resource_id",
					Operators: []string{"=", "<>"},
					Require:   plugin.Optional,
				},
				{
					Name:      "dimension_key",
					Operators: []string{"="},
					Require:   plugin.Optional,
				},
				{
					Name:      "dimension_value",
					Operators: []string{"="},
					Require:   plugin.Optional,
				},
				{
					Name:       "period_start",
					Require:    plugin.Optional,
					Operators:  []string{">", ">=", "=", "<", "<="},
					CacheMatch: query_cache.CacheMatchExact,
				},
				{
					Name:       "period_end",
					Require:    plugin.Optional,
					Operators:  []string{">", ">=", "=", "<", "<="},
					CacheMatch: query_cache.CacheMatchExact,
				},
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
					Name:        "dimension_key",
					Description: "The name of the dimension key.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromQual("dimension_key"),
					Default:     "LINKED_ACCOUNT",
				},
				{
					Name:        "dimension_value",
					Description: "The value of the dimension key.",
					Type:        proto.ColumnType_STRING,
					Hydrate:     getDimensionValue,
					Transform:   transform.FromValue(),
				},
			}),
		),
	}
}

func listCostByResourceMonthly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	params := buildCostByResourceInput("MONTHLY", d)

	// We must have to provide a single filter value to make the API call
	if params.Filter == nil {
		// default filter value
		defaultFilter, err := getDefaultFilterValue(ctx, d, h)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cost_by_resource_monthly.listCostByResourceMonthly", "getDefaultFilterValue", err)
			return nil, err
		}
		params.Filter = defaultFilter
	}

	return streamCostAndUsageByResource(ctx, d, params)
}

//// Common Functions used by aws_cost_by_resource_* tables ////

// dimension_value is not available in the response, so we need to get it from the query data else the default dimension value(caller account id) will be used.
func getDimensionValue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dimensionValue := d.EqualsQualString("dimension_value")

	if dimensionValue != "" {
		return dimensionValue, nil
	}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	return commonColumnData.AccountId, nil
}

func buildCostByResourceInput(granularity string, d *plugin.QueryData) *costexplorer.GetCostAndUsageWithResourcesInput {
	if d == nil {
		return nil
	}

	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}

	now := time.Now()
	endTime := now.Format(timeFormat)

	// Set the default start time to 14 days ago.
	//
	// AWS Cost Explorer supports resource-level cost data with the following granularities:
	//
	// DAILY:
	//   - Available for the most recent 14 days only.
	//   - There is no option to extend this retention period.
	//   - Enables fine-grained cost visibility, but within a short timeframe.
	//
	// MONTHLY:
	//   - Also defaults to 14 days of data unless historical resource-level export is enabled.
	//   - With historical export enabled, monthly data can be retained for up to 38 months.
	//
	// HOURLY:
	//   - Always limited to the past 14 days.
	//   - Incurs additional cost when enabled.
	//   - Retention cannot be extended beyond 14 days.
	//
	// These restrictions apply specifically to resource-level cost data.
	// Non-resource-level cost data has broader retention: up to 12 months (daily) and 38 months (monthly).
	startTime := time.Now().AddDate(0, 0, -13).Format(timeFormat)

	// Get `period_start` and `period_end` from the quals if provided.
	st, et := getSearchStartTimeAndSearchEndTime(d, granularity)
	if st != "" {
		startTime = st
	}
	if et != "" {
		endTime = et
	}

	params := &costexplorer.GetCostAndUsageWithResourcesInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.Granularity(granularity),
		Metrics:     AllCostMetrics(),
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionType("DIMENSION"),
				Key:  aws.String("RESOURCE_ID"),
			},
		},
	}

	var filters []types.Expression

	// Only add resource_id filter if it exists
	if d.Quals != nil && d.Quals["resource_id"] != nil {
		filter := types.Expression{
			Dimensions: &types.DimensionValues{
				Key:          "RESOURCE_ID",
				MatchOptions: []types.MatchOption{types.MatchOption(types.MatchOptionEquals)},
				Values:       []string{d.EqualsQualString("resource_id")},
			},
		}
		filters = append(filters, filter)
	}

	dimKey := d.EqualsQualString("dimension_key")
	dimValue := d.EqualsQualString("dimension_value")
	if dimKey != "" && dimValue != "" {
		filter := types.Expression{
			Dimensions: &types.DimensionValues{
				Key:    types.Dimension(strings.ToUpper(dimKey)),
				Values: []string{dimValue},
			},
		}
		filters = append(filters, filter)
	}

	// Add filters to params if we have any
	if len(filters) > 1 {
		params.Filter = &types.Expression{
			And: filters,
		}
	} else if len(filters) == 1 {
		params.Filter = &filters[0]
	}

	return params
}

// Get default filter value with Dimension "LinkedAccount" and value is the current account ID
func getDefaultFilterValue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*types.Expression, error) {

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	filter := &types.Expression{
		Dimensions: &types.DimensionValues{
			Key:    types.DimensionLinkedAccount,
			Values: []string{commonColumnData.AccountId},
		},
	}

	return filter, nil
}

func streamCostAndUsageByResource(ctx context.Context, d *plugin.QueryData, params *costexplorer.GetCostAndUsageWithResourcesInput) (interface{}, error) {
	// Create session
	svc, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("streamCostAndUsageByResource", "client_error", err)
		return nil, err
	}

	// List call
	for {
		output, err := svc.GetCostAndUsageWithResources(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("streamCostAndUsageByResource", "api_error", err)
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
