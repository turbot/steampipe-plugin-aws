package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableAwsCostAndUsage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_usage",
		Description: "AWS Cost Explorer - Cost and Usage",
		List: &plugin.ListConfig{
			//KeyColumns: plugin.AllColumns([]string{"search_start_time", "search_end_time", "granularity", "dimension_type_1", "dimension_type_2"}),
			KeyColumns: plugin.AllColumns([]string{"granularity", "dimension_type_1", "dimension_type_2"}),
			Hydrate:    listCostAndUsage,
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
	params := buildInputFromQuals(d.KeyColumnQuals)
	return streamCostAndUsage(ctx, d, params)
}

func buildInputFromQuals(keyQuals map[string]*proto.QualValue) *costexplorer.GetCostAndUsageInput {
	granularity := strings.ToUpper(keyQuals["granularity"].GetStringValue())
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := getCEStartDateForGranularity(granularity).Format(timeFormat)

	//dim1 := strings.ToUpper(keyQuals["dimension_type_1"].GetStringValue())
	//dim2 := strings.ToUpper(keyQuals["dimension_type_2"].GetStringValue())

	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.Granularity(granularity),
		Metrics:     AllCostMetrics(),
	}
	var groupings []types.GroupDefinition

	dim1 := keyQuals["dimension_type_1"].GetStringValue()
	if dim1 != "" {
		groupings = append(groupings, buildGroupingDefinition(dim1))
	}

	dim2 := keyQuals["dimension_type_2"].GetStringValue()
	if dim2 != "" {
		groupings = append(groupings, buildGroupingDefinition(dim2))
	}

	// if dim1 != "" {
	// 	groupings = append(groupings, types.GroupDefinition{
	// 		Type: types.GroupDefinitionType("DIMENSION"), 
	// 		Key:  aws.String(dim1),
	// 	})
	// }
	// if dim2 != "" {
	// 	groupings = append(groupings, types.GroupDefinition{
	// 		Type: types.GroupDefinitionType("DIMENSION"),
	// 		Key:  aws.String(dim2),
	// 	})
	// }
	params.GroupBy = groupings

	return params
}

//// HYDRATE FUNCTIONS

// func hydrateKeyQuals(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("hydrateKeyQuals")
// 	plugin.Logger(ctx).Warn("hydrateKeyQuals", "d.KeyColumnQuals", d.KeyColumnQuals)
// 	quals := make(map[string]interface{})

// 	for k, v := range d.KeyColumnQuals {
// 		quals[k] = v.Value
// 	}
// 	plugin.Logger(ctx).Warn("hydrateKeyQuals", "quals", quals)

// 	return &quals, nil
// }

///////////

func getDimKeyAndValue(dimensionString string) (string, string) {
	dim := strings.Split(dimensionString, ":")
	if len(dim) == 1 {
		return strings.ToUpper(dim[0]), ""
	}
	return strings.ToUpper(dim[0]), dim[1]
}

func buildGroupingDefinition(dimensionString string) types.GroupDefinition {
	k, v :=  getDimKeyAndValue(dimensionString)
	if k == "TAG" {
		return types.GroupDefinition{
			Type: types.GroupDefinitionType("TAG"),
			Key:  aws.String(v),
		}
	}
	return types.GroupDefinition{
		Type: types.GroupDefinitionType("DIMENSION"),
		Key:  aws.String(k),
	}
}