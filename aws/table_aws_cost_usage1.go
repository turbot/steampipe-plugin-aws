package aws

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsCostAndUsage1(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_usage1",
		Description: "AWS Cost Explorer - Cost and Usage",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "granularity", Require: plugin.Required},
				{Name: "dimension_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "dimension_values", Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "cost_category_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "cost_category_values", Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "dimension_type_1", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "dimension_type_2", Require: plugin.Optional, Operators: []string{"="}},
			},
			Hydrate: listCostAndUsage1,
		},
		Columns: awsColumns(
			costExplorerColumns([]*plugin.Column{
				{
					Name:        "dimension_1",
					Description: "Valid values are AZ, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, PURCHASE_TYPE, SERVICE, USAGE_TYPE, PLATFORM, TENANCY, RECORD_TYPE, LEGAL_ENTITY_NAME, DEPLOYMENT_OPTION, DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, REGION, BILLING_ENTITY, RESERVATION_ID, SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, OPERATING_SYSTEM",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "dimension_2",
					Description: "Valid values are AZ, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, PURCHASE_TYPE, SERVICE, USAGE_TYPE, PLATFORM, TENANCY, RECORD_TYPE, LEGAL_ENTITY_NAME, DEPLOYMENT_OPTION, DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, REGION, BILLING_ENTITY, RESERVATION_ID, SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, OPERATING_SYSTEM",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "dimension_key",
					Description: "TODO",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromQual("dimension_key"),
				},
				{
					Name:        "dimension_values",
					Description: "TODO",
					Type:        proto.ColumnType_JSON,
					Transform:   transform.FromQual("dimension_values"),
				},
				{
					Name:        "cost_category_key",
					Description: "TODO",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromQual("cost_category_key"),
				},
				{
					Name:        "cost_category_values",
					Description: "TODO",
					Type:        proto.ColumnType_JSON,
					Transform:   transform.FromQual("cost_category_values"),
				},
				{
					Name:        "dimension_type_1",
					Description: "",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "dimension_type_2",
					Description: "",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},

				// Quals columns - to filter the lookups
				{
					Name:        "granularity",
					Description: "",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "search_start_time",
					Description: "",
					Type:        proto.ColumnType_TIMESTAMP,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "search_end_time",
					Description: "",
					Type:        proto.ColumnType_TIMESTAMP,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "cost_cost_category_key",
					Description: "The unique name of the Cost Category.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "cost_category_value",
					Description: "The specific value of the Cost Category.",
					Type:        proto.ColumnType_STRING,
				},
				// {
				// 	Name:        "dimension_type_1",
				// 	Description: "",
				// 	Type:        proto.ColumnType_STRING,
				// 	// Hydrate:     hydrateCostAndUsageQuals,
				// 	Transform: transform.FromField("Dimension1"),
				// },
				// {
				// 	Name:        "dimension_type_2",
				// 	Description: "",
				// 	Type:        proto.ColumnType_STRING,
				// 	// Hydrate:     hydrateCostAndUsageQuals,
				// 	Transform: transform.FromField("Dimension2"),
				// },
				// {
				// 	Name:        "raw_quals",
				// 	Description: "",
				// 	Type:        proto.ColumnType_STRING,
				// 	Hydrate:     hydrateKeyQuals,
				// 	Transform:   transform.FromValue(),
				// },
				// {
				// 	Name:        "raw",
				// 	Description: "raw data",
				// 	Type:        proto.ColumnType_JSON,
				// 	Transform:   transform.FromValue(),
				// },

				//Standard columns for all tables
				// {
				// 	Name:        "tags",
				// 	Description: resourceInterfaceDescription("tags"),
				// 	Type:        proto.ColumnType_JSON,
				// 	Transform:   transform.FromConstant(nil),
				// },
				// {
				// 	Name:        "title",
				// 	Description: resourceInterfaceDescription("title"),
				// 	Type:        proto.ColumnType_STRING,
				// 	Transform:   transform.FromField("ServiceCode"),
				// },
				// {
				// 	Name:        "akas",
				// 	Description: resourceInterfaceDescription("akas"),
				// 	Type:        proto.ColumnType_JSON,
				// 	Hydrate:     getAwsVpcTurbotData,
				// 	Transform:   transform.FromValue(),
				// },
			}),
		),
	}
}

//// LIST FUNCTION

func listCostAndUsage1(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params, err := buildInputFromQuals1(d.KeyColumnQuals, d.KeyColumnQuals, ctx)
	if err != nil {
		return nil, err
	}
	// var rows []CEMetricRow
	plugin.Logger(ctx).Debug("PARAMS +++++>>>>", params)
	_, err = streamCostAndUsage(ctx, d, params)
	if err != nil {
		plugin.Logger(ctx).Error("listCostAndUsage1", "streamCostAndUsage", err)
		return nil, err
	}

	return nil, nil
}

func buildInputFromQuals1(keyQuals map[string]*proto.QualValue, queryCols plugin.KeyColumnEqualsQualMap, ctx context.Context) (*costexplorer.GetCostAndUsageInput, error ){
	granularity := strings.ToUpper(keyQuals["granularity"].GetStringValue())

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
	}

	filter := &costexplorer.Expression{}

	var dims []interface{}
	var catg []interface{}

	dimensionKey := queryCols["dimension_key"].GetStringValue()
	dimensionValues := queryCols["dimension_values"].GetJsonbValue()
	categoryKey := queryCols["cost_category_key"].GetStringValue()
	categoryValues := queryCols["cost_category_values"].GetJsonbValue()

	var dimKey, dimValues, catgKey, catgValues, dimensionType1, dimensionType2 = "", []string{}, "", []string{}, "", ""

	var catgAndDimFilter []*costexplorer.Expression
	// Get quals value for setting the param
	// And expression must have at least 2 operands
	if dimensionKey != "" || dimensionValues != "" {
		json.Unmarshal([]byte(dimensionValues), &dims)
		dimKey = *aws.String(dimensionKey)
		for _, d := range dims {
			dimValues = append(dimValues, d.(string))
		}

		if dimKey == "" || len(dimValues) == 0 { // Check for dimension key and dimension values
			err := errors.New(`You must have to pass both dimension_key and dimension_values with following format "dimension_key = 'string' and dimension_values = '["string", "string"...]'"`)
			return nil, err
		} else {
			dimensionFilter := &costexplorer.DimensionValues{
				Key:          aws.String(dimKey),
				MatchOptions: []*string{aws.String("EQUALS")},
				Values:       aws.StringSlice(dimValues),
			}
			if categoryKey != "" && categoryValues != "" {
				catgAndDimFilter = append(catgAndDimFilter, &costexplorer.Expression{
					Dimensions: dimensionFilter,
				})
			} else {
				filter.Dimensions = dimensionFilter
				params.SetFilter(filter)
			}
		}

	}
	if categoryKey != "" || categoryValues != "" {
		json.Unmarshal([]byte(categoryValues), &catg)
		catgKey = categoryKey
		for _, v := range catg {
			catgValues = append(catgValues, v.(string))
		}

		if catgKey == "" || len(catgValues) <= 0 { // Check for category key and category values
			err := errors.New(`You must have to pass both cost_category_key and cost_category_values with following format "cost_category_key = 'string' and cost_category_values = '["string", "string",...]'"`)
			return nil, err
		} else {
			categoryFilter := &costexplorer.CostCategoryValues{
				Key:          aws.String(catgKey),
				MatchOptions: []*string{aws.String("EQUALS")},
				Values:       aws.StringSlice(catgValues),
			}
			if dimensionKey != "" && dimensionValues != "" {
				catgAndDimFilter = append(catgAndDimFilter, &costexplorer.Expression{
					CostCategories: categoryFilter,
				})
			} else {
				filter.CostCategories = categoryFilter
				params.SetFilter(filter)
			}
		}

	}

	if len(catgAndDimFilter) > 0 {
		params.SetFilter(filter)
		params.Filter.And = catgAndDimFilter
	}
	if queryCols["dimension_type_1"].GetStringValue() != "" || queryCols["dimension_type_2"].GetStringValue() != "" {
		dimensionType1 = queryCols["dimension_type_1"].GetStringValue()
		dimensionType2 = queryCols["dimension_type_2"].GetStringValue()

		// Check for dimension type
		if dimensionType1 != "" && dimensionType2 == "" {
			err := errors.New(`You must have to pass dimension_type_2, if you are passing dimension_type_1`)
			return nil, err
		} else if dimensionType1 == "" && dimensionType2 != "" {
			err := errors.New(`You must have to pass dimension_type_1, if you are passing dimension_type_2`)
			return nil, err
		} else {
			var groupings []*costexplorer.GroupDefinition
			groupings = append(groupings, &costexplorer.GroupDefinition{Type: aws.String("DIMENSION"), Key: aws.String(dimensionType1)})
			groupings = append(groupings, &costexplorer.GroupDefinition{Type: aws.String("DIMENSION"), Key: aws.String(dimensionType2)})

			params.SetGroupBy(groupings)
		}

	}

	if dimensionKey == "" && dimensionValues == "" && dimensionType1 == "" && dimensionType2 == "" && categoryKey == "" && categoryValues == "" {
		err := errors.New("You must have to pass at least either of the column pairs 'dimension_key and dimension_values' or 'cost_category_key and cost_category_values' or 'dimension_type_1 and dimension_type_2'")
			return nil, err
	}

	return params, nil
}
