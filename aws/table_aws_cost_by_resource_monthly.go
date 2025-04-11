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
)

func tableAwsCostByResourceMonthly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_resource_monthly",
		Description: "AWS Cost Explorer - Cost by Resource (Monthly)",
		List: &plugin.ListConfig{
			Hydrate: listCostByResourceMonthly,
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

func listCostByResourceMonthly(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCostByResourceMonthly", "client_error", err)
		return nil, err
	}

	params := buildCostByResourceInput("MONTHLY", d)

	// List call
	for {
		output, err := svc.GetCostAndUsageWithResources(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("listCostByResourceMonthly", "api_error", err)
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

func buildCostByResourceInput(granularity string, d *plugin.QueryData) *costexplorer.GetCostAndUsageWithResourcesInput {
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := getCEStartDateForGranularityWithResources(granularity).Format(timeFormat)

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
			{
				Type: types.GroupDefinitionType("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
		},
	}

	var filters []types.Expression

	for _, keyQual := range d.Table.List.KeyColumns {
		filterQual := d.Quals[keyQual.Name]
		if filterQual == nil {
			continue
		}
		for _, qual := range filterQual.Quals {
			if qual.Value != nil {
				value := qual.Value
				switch qual.Operator {
				case "=":
					filter := types.Expression{}
					filter.Dimensions = &types.DimensionValues{}
					filter.Dimensions.Key = types.Dimension(strings.ToUpper(keyQual.Name))
					filter.Dimensions.Values = []string{value.GetStringValue()}
					filters = append(filters, filter)
				case "<>":
					filter := types.Expression{}
					filter.Not = &types.Expression{}
					filter.Not.Dimensions = &types.DimensionValues{}
					filter.Not.Dimensions.Key = types.Dimension(strings.ToUpper(keyQual.Name))
					filter.Not.Dimensions.Values = []string{value.GetStringValue()}
					filters = append(filters, filter)
				}
			}
		}
	}

	// Make sure at least RESOURCE_ID is in either Filter or GroupBy as required by the API
	// hasResourceIdFilter := false
	// for _, filter := range filters {
	// 	if filter.Dimensions != nil && filter.Dimensions.Key == "RESOURCE_ID" {
	// 		hasResourceIdFilter = true
	// 		break
	// 	}
	// 	if filter.Not != nil && filter.Not.Dimensions != nil && filter.Not.Dimensions.Key == "RESOURCE_ID" {
	// 		hasResourceIdFilter = true
	// 		break
	// 	}
	// }

	if len(filters) > 1 {
		params.Filter = &types.Expression{
			And: filters,
		}
	} else if len(filters) == 1 {
		params.Filter = &(filters[0])
	}

	return params
}
