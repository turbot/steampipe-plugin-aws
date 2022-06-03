package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsCostByServiceMonthly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_by_service_monthly",
		Description: "AWS Cost Explorer - Cost by Service (Monthly)",
		List: &plugin.ListConfig{
			Hydrate: listCostByServiceMonthly,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "service", Operators: []string{"=", "<>"}, Require: plugin.Optional},
			},
		},
		Columns: awsColumns(
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

func listCostByServiceMonthly(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildCostByServiceInput("MONTHLY", d)
	return streamCostAndUsage(ctx, d, params)
}

func buildCostByServiceInput(granularity string, d *plugin.QueryData) *costexplorer.GetCostAndUsageInput {
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := getCEStartDateForGranularity(granularity).Format(timeFormat)

	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: aws.String(granularity),
		Metrics:     aws.StringSlice(AllCostMetrics()),
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
		},
	}

	var filters []*costexplorer.Expression

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
					filter := &costexplorer.Expression{}
					filter.Dimensions = &costexplorer.DimensionValues{}
					filter.Dimensions.Key = aws.String(strings.ToUpper(keyQual.Name))
					filter.Dimensions.Values = aws.StringSlice([]string{value.GetStringValue()})
					filters = append(filters, filter)
				case "<>":
					filter := &costexplorer.Expression{}
					filter.Not = &costexplorer.Expression{}
					filter.Not.Dimensions = &costexplorer.DimensionValues{}
					filter.Not.Dimensions.Key = aws.String(strings.ToUpper(keyQual.Name))
					filter.Not.Dimensions.Values = aws.StringSlice([]string{value.GetStringValue()})
					filters = append(filters, filter)
				}
			}
		}
	}

	if len(filters) > 1 {
		params.Filter = &costexplorer.Expression{
			And: filters,
		}
	} else if len(filters) == 1 {
		params.Filter = filters[0]
	}

	return params
}
