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
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAwsCostAndUsage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_usage",
		Description: "AWS Cost Explorer - Cost and Usage",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "granularity",
					Require: plugin.Required,
				},
				{
					Name:    "dimension_type_1",
					Require: plugin.Required,
				},
				{
					Name:    "dimension_type_2",
					Require: plugin.Required,
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

			Hydrate: listCostAndUsage,
			Tags:    map[string]string{"service": "ce", "action": "GetCostAndUsage"},
		},
		Columns: awsGlobalRegionColumns(
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
					Name:        "search_start_time",
					Description: "[Deprecated] The beginning of the time period.",
					Type:        proto.ColumnType_TIMESTAMP,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "search_end_time",
					Description: "[Deprecated] The end of the time period.",
					Type:        proto.ColumnType_TIMESTAMP,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				// Quals columns - to filter the lookups
				{
					Name:        "granularity",
					Description: "The AWS cost granularity. Valid values are DAILY, MONTHLY, or HOURLY.",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "dimension_type_1",
					Description: "The first dimension to group results by. Valid values include AZ, INSTANCE_TYPE, LINKED_ACCOUNT, LINKED_ACCOUNT_NAME, OPERATION, PURCHASE_TYPE, REGION, SERVICE, SERVICE_CODE, USAGE_TYPE, USAGE_TYPE_GROUP, RECORD_TYPE, OPERATING_SYSTEM, TENANCY, SCOPE, PLATFORM, SUBSCRIPTION_ID, LEGAL_ENTITY_NAME, DEPLOYMENT_OPTION, DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, BILLING_ENTITY, RESERVATION_ID, RESOURCE_ID, RIGHTSIZING_TYPE, SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, PAYMENT_OPTION.",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},
				{
					Name:        "dimension_type_2",
					Description: "The second dimension to group results by. Valid values include AZ, INSTANCE_TYPE, LINKED_ACCOUNT, LINKED_ACCOUNT_NAME, OPERATION, PURCHASE_TYPE, REGION, SERVICE, SERVICE_CODE, USAGE_TYPE, USAGE_TYPE_GROUP, RECORD_TYPE, OPERATING_SYSTEM, TENANCY, SCOPE, PLATFORM, SUBSCRIPTION_ID, LEGAL_ENTITY_NAME, DEPLOYMENT_OPTION, DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, BILLING_ENTITY, RESERVATION_ID, RESOURCE_ID, RIGHTSIZING_TYPE, SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, PAYMENT_OPTION.",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostAndUsageQuals,
				},
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

func listCostAndUsage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildInputFromQuals(ctx, d)
	return streamCostAndUsage(ctx, d, params)
}

func buildInputFromQuals(ctx context.Context, keyQuals *plugin.QueryData) *costexplorer.GetCostAndUsageInput {
	granularity := strings.ToUpper(keyQuals.EqualsQuals["granularity"].GetStringValue())
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := getCEStartDateForGranularity(granularity).Format(timeFormat)

	st, et := getSearchStartTimeAndSearchEndTime(keyQuals, granularity)
	if st != "" {
		startTime = st
	}
	if et != "" {
		endTime = et
	}

	selectedMetrics := AllCostMetrics()
	if len(getMetricsByQueryContext(keyQuals.QueryContext)) > 0 {
		selectedMetrics = getMetricsByQueryContext(keyQuals.QueryContext)
	}

	dim1 := strings.ToUpper(keyQuals.EqualsQuals["dimension_type_1"].GetStringValue())
	dim2 := strings.ToUpper(keyQuals.EqualsQuals["dimension_type_2"].GetStringValue())

	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.Granularity(granularity),
		Metrics:     selectedMetrics,
	}
	var groupings []types.GroupDefinition
	if dim1 != "" {
		groupings = append(groupings, types.GroupDefinition{
			Type: types.GroupDefinitionType("DIMENSION"),
			Key:  aws.String(dim1),
		})
	}
	if dim2 != "" {
		groupings = append(groupings, types.GroupDefinition{
			Type: types.GroupDefinitionType("DIMENSION"),
			Key:  aws.String(dim2),
		})
	}
	params.GroupBy = groupings

	return params
}

func getSearchStartTimeAndSearchEndTime(keyQuals *plugin.QueryData, granularity string) (string, string) {
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}

	st, et := "", ""

	if keyQuals.Quals["period_start"] != nil && len(keyQuals.Quals["period_start"].Quals) <= 1 {
		for _, q := range keyQuals.Quals["period_start"].Quals {
			t := q.Value.GetTimestampValue().AsTime().Format(timeFormat)
			switch q.Operator {
			case "=", ">=", ">":
				st = t
			case "<", "<=":
				et = t
			}
		}
	}

	// The API supports a single value with the '=' operator.
	// For queries like: "period_end BETWEEN current_timestamp - interval '31d' AND current_timestamp - interval '1d'", the FDW parses the query parameters with multiple qualifiers.
	// In this case, we will have multiple qualifiers with operators such as:
	// 1. The length of keyQuals.Quals["period_end"].Quals will be 2.
	// 2. The qualifier values would be "2024-05-10" with the '>=' operator and "2024-06-09" with the '<=' operator.
	// Plugin Log:
	// 2024-06-10 11:17:39.071 UTC [DEBUG] steampipe-plugin-aws.plugin: [ERROR] 1718018259212: Period end Scan Length ===>>> : EXTRA_VALUE_AT_END=2
	// 2024-06-10 11:17:39.071 UTC [DEBUG] steampipe-plugin-aws.plugin: [ERROR] 1718018259212: Period End => : >=2024-05-10
	// 2024-06-10 11:17:39.071 UTC [DEBUG] steampipe-plugin-aws.plugin: [ERROR] 1718018259212: Period End => : <=2024-06-09
	// In this scenario, manipulating the start and end time is a bit difficult and challenging.
	// Let the API fetch all the rows, and filtering will occur at the Steampipe level.

	if keyQuals.Quals["period_end"] != nil && len(keyQuals.Quals["period_end"].Quals) <= 1 {
		for _, q := range keyQuals.Quals["period_end"].Quals {
			t := q.Value.GetTimestampValue().AsTime().Format(timeFormat)
			switch q.Operator {
			case "=", ">=", ">":
				if st == "" {
					st = t
				}
			case "<", "<=":
				if et == "" {
					et = t
				}
			}
		}
	}
	return st, et
}

//// HYDRATE FUNCTIONS

// func hydrateKeyQuals(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("hydrateKeyQuals")
// 	plugin.Logger(ctx).Warn("hydrateKeyQuals", "d.EqualsQuals", d.EqualsQuals)
// 	quals := make(map[string]interface{})

// 	for k, v := range d.EqualsQuals {
// 		quals[k] = v.Value
// 	}
// 	plugin.Logger(ctx).Warn("hydrateKeyQuals", "quals", quals)

// 	return &quals, nil
// }

///////////
