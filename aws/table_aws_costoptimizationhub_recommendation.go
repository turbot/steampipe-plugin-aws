package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/costoptimizationhub"
	"github.com/aws/aws-sdk-go-v2/service/costoptimizationhub/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCostOptimizationHubRecommendation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_costoptimizationhub_recommendation",
		Description: "AWS Cost Optimization Hub Recommendation",
		List: &plugin.ListConfig{
			Hydrate: listCostOptimizationHubRecommendations,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "recommendation_account_id",
					Require: plugin.Optional,
				},
				{
					Name:    "action_type",
					Require: plugin.Optional,
				},
				{
					Name:    "implementation_effort",
					Require: plugin.Optional,
				},
				{
					Name:    "recommendation_id",
					Require: plugin.Optional,
				},
				{
					Name:    "resource_region",
					Require: plugin.Optional,
				},
				{
					Name:    "resource_arn",
					Require: plugin.Optional,
				},
				{
					Name:    "resource_id",
					Require: plugin.Optional,
				},
				{
					Name:    "current_resource_type",
					Require: plugin.Optional,
				},
				{
					Name:    "recommended_resource_type",
					Require: plugin.Optional,
				},
				{
					Name:    "restart_needed",
					Require: plugin.Optional,
				},
				{
					Name:    "rollback_possible",
					Require: plugin.Optional,
				},
			},
			Tags: map[string]string{"service": "cost-optimization-hub", "action": "ListRecommendations"},
		},
		Columns: awsGlobalRegionColumns(
			[]*plugin.Column{
				{
					Name:        "recommendation_id",
					Description: "The ID for the recommendation.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "resource_arn",
					Description: "The Amazon Resource Name (ARN) for the recommendation.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "resource_id",
					Description: "The resource ID for the recommendation.",
					Type:        proto.ColumnType_STRING,
				},

				// We have a common column named "account_id" for all the tables that represents current caller account ID, so renamed it to "recommendation_account_id" to avoid ambiguity.
				{
					Name:        "recommendation_account_id",
					Description: "The account that the recommendation is for.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("AccountId"),
				},
				{
					Name:        "action_type",
					Description: "The type of tasks that can be carried out by this action.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "currency_code",
					Description: "The currency code used for the recommendation.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "current_resource_summary",
					Description: "Describes the current resource.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "current_resource_type",
					Description: "The current resource type.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "estimated_monthly_cost",
					Description: "The estimated monthly cost for the recommendation.",
					Type:        proto.ColumnType_DOUBLE,
				},
				{
					Name:        "estimated_monthly_savings",
					Description: "The estimated monthly savings amount for the recommendation.",
					Type:        proto.ColumnType_DOUBLE,
				},
				{
					Name:        "estimated_savings_percentage",
					Description: "The estimated savings percentage relative to the total cost over the cost calculation lookback period.",
					Type:        proto.ColumnType_DOUBLE,
				},
				{
					Name:        "implementation_effort",
					Description: "The effort required to implement the recommendation.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "last_refresh_timestamp",
					Description: "The time when the recommendation was last generated.",
					Type:        proto.ColumnType_TIMESTAMP,
				},
				{
					Name:        "recommendation_lookback_period_in_days",
					Description: "The lookback period that's used to generate the recommendation.",
					Type:        proto.ColumnType_INT,
				},
				{
					Name:        "recommended_resource_summary",
					Description: "Describes the recommended resource.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "recommended_resource_type",
					Description: "Describes the recommended resource.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "resource_region",
					Description: "The Amazon Web Services Region of the resource.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Region"),
				},
				{
					Name:        "restart_needed",
					Description: "Whether or not implementing the recommendation requires a restart.",
					Type:        proto.ColumnType_BOOL,
				},
				{
					Name:        "rollback_possible",
					Description: "Whether or not implementing the recommendation can be rolled back.",
					Type:        proto.ColumnType_BOOL,
				},
				{
					Name:        "source",
					Description: "The source of the recommendation.",
					Type:        proto.ColumnType_STRING,
				},
				{
					Name:        "current_resource_details",
					Description: "The details for the resource.",
					Type:        proto.ColumnType_JSON,
					Hydrate:     getCostOptimizationHubRecommendations,
					Transform:   transform.FromField("CurrentResourceDetails"),
				},
				{
					Name:        "recommended_resource_details",
					Description: "The details about the recommended resource.",
					Type:        proto.ColumnType_JSON,
					Hydrate:     getCostOptimizationHubRecommendations,
					Transform:   transform.FromField("RecommendedResourceDetails"),
				},
				{
					Name:        "tags_src",
					Description: "A list of tags assigned to the recommendation.",
					Type:        proto.ColumnType_JSON,
				},

				// Steampipe standard columns
				{
					Name:        "tags",
					Description: resourceInterfaceDescription("tags"),
					Type:        proto.ColumnType_JSON,
					Transform:   transform.From(costOptimizationRecommendationTurbotTags),
				},
				{
					Name:        "title",
					Description: resourceInterfaceDescription("title"),
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("RecommendationId"),
				},
			}),
	}
}

//// LIST FUNCTION

func listCostOptimizationHubRecommendations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildCostOptimizationHubRecommendationInputFromQuals(d.Quals)

	// Create Client
	svc, err := CostOptimizationHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_costoptimizationhub_recommendation.listCostOptimizationHubRecommendations", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &costoptimizationhub.ListRecommendationsInput{
		MaxResults: &maxLimit,
	}

	if params != nil {
		input.Filter = params
	}

	paginator := costoptimizationhub.NewListRecommendationsPaginator(svc, input, func(o *costoptimizationhub.ListRecommendationsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_costoptimizationhub_recommendation.listCostOptimizationHubRecommendations", "api_error", err)
			return nil, err
		}

		for _, item := range output.Items {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCostOptimizationHubRecommendations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	recommendation := h.Item.(types.Recommendation)

	// Create Client
	svc, err := CostOptimizationHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_costoptimizationhub_recommendation.getCostOptimizationHubRecommendations", "connection_error", err)
		return nil, err
	}

	input := &costoptimizationhub.GetRecommendationInput{
		RecommendationId: recommendation.RecommendationId,
	}

	result, err := svc.GetRecommendation(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_costoptimizationhub_recommendation.getCostOptimizationHubRecommendations", "api_error", err)
		return nil, err
	}

	return result, nil
}

//// TRANSFORM FUNCTIONS

func costOptimizationRecommendationTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	r := d.HydrateItem.(types.Recommendation)
	var turbotTagsMap map[string]string
	if r.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range r.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}

////  Build input parameter fot list API call

func buildCostOptimizationHubRecommendationInputFromQuals(quals plugin.KeyColumnQualMap) *types.Filter {
	param := &types.Filter{}
	filterQuals := []string{"recommendation_account_id", "action_type", "implementation_effort", "recommendation_id", "resource_region", "resource_arn", "resource_id", "current_resource_type", "recommended_resource_type", "restart_needed", "rollback_possible"}

	for _, columnName := range filterQuals {
		if quals[columnName] != nil {
			switch columnName {
			case "restart_needed", "rollback_possible":
				value := getQualsValueByColumn(quals, columnName, "boolean")
				val := value.(bool)
				if columnName == "restart_needed" {
					param.RestartNeeded = &val
				}
				if columnName == "rollback_possible" {
					param.RollbackPossible = &val
				}
			default:
				value := getQualsValueByColumn(quals, columnName, "string")
				switch columnName {
				case "recommendation_account_id":
					param.AccountIds = []string{fmt.Sprint(value)}
				case "recommendation_id":
					param.RecommendationIds = []string{fmt.Sprint(value)}
				case "resource_region":
					param.Regions = []string{fmt.Sprint(value)}
				case "resource_arn":
					param.ResourceArns = []string{fmt.Sprint(value)}
				case "resource_id":
					param.ResourceIds = []string{fmt.Sprint(value)}
				case "implementation_effort":
					param.ImplementationEfforts = []types.ImplementationEffort{types.ImplementationEffort(value.(string))}
				case "action_type":
					param.ActionTypes = []types.ActionType{types.ActionType(value.(string))}
				case "current_resource_type", "recommended_resource_type":
					param.ResourceTypes = []types.ResourceType{types.ResourceType(value.(string))}
				}
			}
		}
	}

	return param
}
