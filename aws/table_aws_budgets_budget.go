package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/budgets"
	"github.com/aws/aws-sdk-go-v2/service/budgets/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBudgetsBudget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_budgets_budget",
		Description: "AWS Budgets Budget",
		List: &plugin.ListConfig{
			Hydrate: listBudgetsBudgets,
			Tags:    map[string]string{"service": "budgets", "action": "DescribeBudgets"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"account_id", "name"}),
			Hydrate:    getBudgetsBudget,
			Tags:       map[string]string{"service": "budgets", "action": "DescribeBudget"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listBudgetsNotifications,
				Tags: map[string]string{"service": "budgets", "action": "DescribeNotificationsForBudget"},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of a budget.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BudgetName"),
			},
			{
				Name:        "type",
				Description: "The type of budget (COST, USAGE, RI_UTILIZATION, RI_COVERAGE, SAVINGS_PLANS_UTILIZATION, SAVINGS_PLANS_COVERAGE).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BudgetType"),
			},
			{
				Name:        "limit_amount",
				Description: "The budgeted amount for the budget.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BudgetLimit.Amount"),
			},
			{
				Name:        "limit_unit",
				Description: "The unit of measurement for the budget amount.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BudgetLimit.Unit"),
			},
			{
				Name:        "calculated_spend_actual_spend",
				Description: "The actual amount spent during the period.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CalculatedSpend.ActualSpend.Amount"),
			},
			{
				Name:        "calculated_spend_forecasted_spend",
				Description: "The forecasted amount for the period.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CalculatedSpend.ForecastedSpend.Amount"),
			},
			{
				Name:        "time_period_start",
				Description: "The start date of the budget.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimePeriod.Start"),
			},
			{
				Name:        "time_period_end",
				Description: "The end date of the budget.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimePeriod.End"),
			},
			{
				Name:        "time_unit",
				Description: "The length of the budget period (MONTHLY, QUARTERLY, ANNUALLY).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TimeUnit"),
			},
			{
				Name:        "cost_filters",
				Description: "The cost filters for the budget.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CostFilters"),
			},
			{
				Name:        "cost_types",
				Description: "The types of costs included in the budget.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CostTypes"),
			},
			{
				Name:        "notifications",
				Description: "The notifications associated with the budget.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listBudgetsNotifications,
				Transform:   transform.FromValue(),
			},
			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BudgetName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBudgetArn,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listBudgetsBudgets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := BudgetsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_budgets_budget.listBudgetsBudgets", "connection_error", err)
		return nil, err
	}

	// Get account ID from common columns
	commonColumnData, err := getCommonColumnsUncached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_budgets_budget.listBudgetsBudgets", "get_account_id_error", err)
		return nil, err
	}
	accountID := commonColumnData.(*awsCommonColumnData).AccountId

	input := &budgets.DescribeBudgetsInput{
		AccountId: aws.String(accountID),
	}

	paginator := budgets.NewDescribeBudgetsPaginator(client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_budgets_budget.listBudgetsBudgets", "api_error", err)
			return nil, err
		}

		for _, budget := range output.Budgets {
			d.StreamListItem(ctx, budget)

			// Context may get cancelled
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// GET FUNCTION

func getBudgetsBudget(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	accountID := d.EqualsQualString("account_id")
	name := d.EqualsQualString("name")

	if accountID == "" || name == "" {
		return nil, nil
	}

	client, err := BudgetsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_budgets_budget.getBudgetsBudget", "connection_error", err)
		return nil, err
	}

	input := &budgets.DescribeBudgetInput{
		AccountId:  aws.String(accountID),
		BudgetName: aws.String(name),
	}

	output, err := client.DescribeBudget(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_budgets_budget.getBudgetsBudget", "api_error", err)
		return nil, err
	}

	return output.Budget, nil
}

//// HYDRATE FUNCTIONS

func listBudgetsNotifications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	budget := h.Item.(types.Budget)

	// Get account ID from common columns
	commonColumnData, err := getCommonColumnsUncached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_budgets_budget.listBudgetsNotifications", "get_account_id_error", err)
		return nil, err
	}
	accountID := commonColumnData.(*awsCommonColumnData).AccountId

	client, err := BudgetsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_budgets_budget.listBudgetsNotifications", "connection_error", err)
		return nil, err
	}

	input := &budgets.DescribeNotificationsForBudgetInput{
		AccountId:  aws.String(accountID),
		BudgetName: budget.BudgetName,
	}

	output, err := client.DescribeNotificationsForBudget(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_budgets_budget.listBudgetsNotifications", "api_error", err)
		return nil, err
	}

	return output.Notifications, nil
}

func getBudgetArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	budget := h.Item.(types.Budget)

	// Get account ID and partition from common columns
	commonColumnData, err := getCommonColumnsUncached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_budgets_budget.getBudgetArn", "get_account_id_error", err)
		return nil, err
	}
	accountID := commonColumnData.(*awsCommonColumnData).AccountId
	partition := commonColumnData.(*awsCommonColumnData).Partition

	arn := "arn:" + partition + ":budgets::" + accountID + ":budget/" + *budget.BudgetName
	return []string{arn}, nil
}
