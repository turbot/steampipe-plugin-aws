package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableAwsCostAndUsageByCategory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_usage_by_category",
		Description: "AWS Cost Explorer - Cost and Usage By Category",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "granularity", Require: plugin.Required},
				{Name: "category_key", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "category_value", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
			Hydrate: listCostAndUsageByCategory,
		},
		Columns: awsColumns(
			costExplorerColumns([]*plugin.Column{
				{
					Name:        "category_key",
					Description: "The unique name of the Cost Category.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "category_value",
					Description: "The specific value of the Cost Category.",
					Type:        proto.ColumnType_STRING,
				},

				// Quals columns - to filter the lookups
				{
					Name:        "granularity",
					Description: "",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "search_start_time", // Optional quals for future implementation
					Description: "",
					Type:        proto.ColumnType_TIMESTAMP,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "search_end_time", // Optional quals for future implementation
					Description: "",
					Type:        proto.ColumnType_TIMESTAMP,
					Hydrate:     hydrateCostAndUsageQuals,
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostAndUsageByCategory(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return buildInputParamAndStreamCostAndUsageByCategory(ctx, d, h)
}

type CostByCategory struct {
	CategoryKey   *string
	CategoryValue *string
	CEMetricRow
}

func buildInputParamAndStreamCostAndUsageByCategory(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("buildInputParamAndStreamCostAndUsageByCategory")
	granularity := strings.ToUpper(getQualsValueByColumn(d.Quals, "granularity", "string").(string))
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := getCEStartDateForGranularity(granularity).Format(timeFormat)

	var categoryKey, categoryValue, categoryKeyOperator, categoryValueOperator string = "", "", "", ""

	optionalColumns := []string{"category_key", "category_value"}

	// Get Optional quals value and given operator
	for _, column := range optionalColumns {
		if d.Quals[column] != nil {
			for _, q := range d.Quals[column].Quals {
				switch column {
				case "category_key":
					categoryKey = q.Value.GetStringValue()
					categoryKeyOperator = q.Operator
				case "category_value":
					categoryValue = q.Value.GetStringValue()
					categoryValueOperator = q.Operator
				}

			}
		}

	}

	// List all cost category definitions
	getCostCategoryDeFinitionsCached := plugin.HydrateFunc(listCostCategoryDefinitions).WithCache()
	c, err := getCostCategoryDeFinitionsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getCostCategoryDeFinitionsCached", "api_error", err)
		return nil, err
	}
	costCategoryDefinitions := c.([]*costexplorer.CostCategoryReference)

	setCategoryKey := false
	for _, definition := range costCategoryDefinitions {

		if categoryKey != "" {
			if categoryKey != *definition.Name && categoryKeyOperator == "=" {
				continue
			}
			if categoryKey == *definition.Name && categoryKeyOperator == "<>" {
				continue
			}
		}

		// If category_key has not been provided in optional qulas, then we need to get cost and usage for all cost category
		if categoryKey == "" {
			categoryKey = *definition.Name
			categoryKeyOperator = "="
			setCategoryKey = true
		}

		// Build the common param without filter
		params := &costexplorer.GetCostAndUsageInput{
			TimePeriod: &costexplorer.DateInterval{
				Start: aws.String(startTime),
				End:   aws.String(endTime),
			},
			Granularity: aws.String(granularity),
			Metrics:     aws.StringSlice(AllCostMetrics()),
		}

		// Get cost category value based of optional quals "category_value"
		values := getCategoryValuesForFilterParam(definition.Values, categoryValue, categoryValueOperator)

		// Empty check for the category value 
		if len(values) <= 0 {
			categoryKey = ""
			continue
		}

		// Add filter in param with given optional qulals cost category value
		for _, value := range values {
			filter := &costexplorer.Expression{
				CostCategories: &costexplorer.CostCategoryValues{
					Key:          definition.Name,
					MatchOptions: []*string{aws.String("EQUALS")}, // The value will be EQUALS always because we are manipulating the data based on optional quals operator 
					Values:       []*string{value},
				},
			}
			params.SetFilter(filter)

			// If we wants to get cost by category for all category key 
			if setCategoryKey {
				categoryKey = ""
			}

			streamCostAndUsageByCategory(ctx, d, h, params, *definition.Name, *value)
		}
	}

	return nil, err
}

func streamCostAndUsageByCategory(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData, params *costexplorer.GetCostAndUsageInput, categoryKey string, categoryValue string) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("streamCostAndUsageByCategory")

	// Create session
	svc, err := CostExplorerService(ctx, d)
	if err != nil {
		logger.Error("streamCostAndUsageByCategory", "connection_error", err)
		return nil, err
	}

	// List call
	for {
		output, err := svc.GetCostAndUsage(params)
		if err != nil {
			logger.Error("streamCostAndUsageByCategory", "err", err)
			return nil, err
		}

		// stream the results...
		for _, row := range buildCEMetricRows(ctx, output, d.KeyColumnQuals) {
			d.StreamListItem(ctx, &CostByCategory{
				CategoryKey:   &categoryKey,
				CategoryValue: &categoryValue,
				CEMetricRow:   row,
			})

			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// get more pages if there are any...
		if output.NextPageToken == nil {
			break
		}
		params.SetNextPageToken(*output.NextPageToken)
	}

	return nil, nil
}

func listCostCategoryDefinitions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listCostCategoryDefinitions")

	// Create session
	svc, err := CostExplorerService(ctx, d)
	if err != nil {
		logger.Error("listCostCategoryDefinitions", "connection_error", err)
		return nil, err
	}
	var definitions []*costexplorer.CostCategoryReference
	input := &costexplorer.ListCostCategoryDefinitionsInput{}
	err = svc.ListCostCategoryDefinitionsPages(input,
		func(page *costexplorer.ListCostCategoryDefinitionsOutput, isLast bool) bool {
			definitions = append(definitions, page.CostCategoryReferences...)
			return !isLast
		},
	)

	if err != nil {
		logger.Error("listCostCategoryDefinitions", "api_error", err)
		return nil, err
	}

	return definitions, nil
}

//// UTILITY FUNCTION
func getCategoryValuesForFilterParam(values []*string, categoryValue string, categoryValueOperator string) []*string {
	var filteredValues []*string
	for _, value := range values {
		if categoryValue == *value && categoryValueOperator == "=" { // List out cost and usage for given category value 
			filteredValues = append(filteredValues, value)
			return filteredValues
		} else if categoryValue != *value && categoryValueOperator == "<>" { // List out cost and usage not for given category value
			filteredValues = append(filteredValues, value)
		} else if categoryValue == "" { // List out cost and usage for all category value
			filteredValues = append(filteredValues, value)
		}
	}

	return filteredValues
}
