package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/savingsplans"
	"github.com/aws/aws-sdk-go-v2/service/savingsplans/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSavingsPlan(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_savings_plan",
		Description: "AWS Savings Plan",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("savings_plan_id"),
			Hydrate:    getSavingsPlan,
			Tags:       map[string]string{"service": "savingsplans", "action": "DescribeSavingsPlans"},
		},
		List: &plugin.ListConfig{
			Hydrate: listSavingsPlans,
			Tags:    map[string]string{"service": "savingsplans", "action": "DescribeSavingsPlans"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "state", Require: plugin.Optional},
				{Name: "region", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "ec2_instance_family", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "commitment", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "term_duration_in_seconds", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "savings_plan_type", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "payment_option", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "start", Require: plugin.Optional, Operators: []string{">="}},
				{Name: "end_time", Require: plugin.Optional, Operators: []string{"<="}},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "savings_plan_id",
				Description: "The ID of the Savings Plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the Savings Plan.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SavingsPlanArn"),
			},
			{
				Name:        "offering_id",
				Description: "The ID of the offering.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "savings_plan_type",
				Description: "The type of Savings Plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "payment_option",
				Description: "The payment option for the Savings Plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the Savings Plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "currency",
				Description: "The currency of the Savings Plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "commitment",
				Description: "The hourly commitment amount for the Savings Plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "upfront_payment_amount",
				Description: "The up-front payment amount for the Savings Plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "recurring_payment_amount",
				Description: "The recurring payment amount for the Savings Plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "term_duration_in_seconds",
				Description: "The duration of the Savings Plan term in seconds.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "start",
				Description: "The start time of the Savings Plan.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "end_time",
				Description: "The end time of the Savings Plan.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("End"), // Renamed from 'end' to 'end_time' to avoid SQL reserved keyword conflicts
			},
			{
				Name:        "returnable_until",
				Description: "The time until which the Savings Plan can be returned.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description of the Savings Plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ec2_instance_family",
				Description: "The instance family of the EC2 Savings Plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "product_types",
				Description: "The product types supported by the Savings Plan.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: "A list of tags associated with the Savings Plan.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SavingsPlanId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SavingsPlanArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSavingsPlans(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := SavingsPlansClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_savings_plan.listSavingsPlans", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	// https://docs.aws.amazon.com/savingsplans/latest/APIReference/API_DescribeSavingsPlans.html#API_DescribeSavingsPlans_RequestSyntax
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &savingsplans.DescribeSavingsPlansInput{
		MaxResults: aws.Int32(maxLimit),
	}

	// Add filters based on optional key columns
	if d.EqualsQuals["state"] != nil {
		input.States = []types.SavingsPlanState{types.SavingsPlanState(d.EqualsQuals["state"].GetStringValue())}
	}

	filters := getSavingsPlanFilter(d)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// Handle pagination manually since AWS SDK v2 doesn't have paginator for this API
	var nextToken *string
	for {
		if nextToken != nil {
			input.NextToken = nextToken
		}

		// Apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := svc.DescribeSavingsPlans(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_savings_plan.listSavingsPlans", "api_error", err)
			return nil, err
		}

		for _, item := range output.SavingsPlans {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Check if there are more pages
		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSavingsPlan(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	savingsPlanId := d.EqualsQuals["savings_plan_id"].GetStringValue()

	// Create service
	svc, err := SavingsPlansClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_savings_plan.getSavingsPlan", "connection_error", err)
		return nil, err
	}

	params := &savingsplans.DescribeSavingsPlansInput{
		SavingsPlanIds: []string{savingsPlanId},
	}

	op, err := svc.DescribeSavingsPlans(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_savings_plan.getSavingsPlan", "api_error", err)
		return nil, err
	}

	if len(op.SavingsPlans) > 0 {
		return op.SavingsPlans[0], nil
	}
	return nil, nil
}

//// UTILITY FUNCTIONS

func getSavingsPlanFilter(d *plugin.QueryData) []types.SavingsPlanFilter {
	savingsPlanFilters := make([]types.SavingsPlanFilter, 0)

	for q, v := range d.EqualsQuals {
		savingPlan := &types.SavingsPlanFilter{}
		addFilter := true

		switch q {
		case "region":
			savingPlan.Name = types.SavingsPlansFilterNameRegion
			savingPlan.Values = []string{v.GetStringValue()}
		case "ec2_instance_family":
			savingPlan.Name = types.SavingsPlansFilterNameEc2InstanceFamily
			savingPlan.Values = []string{v.GetStringValue()}
		case "commitment":
			savingPlan.Name = types.SavingsPlansFilterNameCommitment
			savingPlan.Values = []string{v.GetStringValue()}
		case "payment_option":
			savingPlan.Name = types.SavingsPlansFilterNamePaymentOption
			savingPlan.Values = []string{v.GetStringValue()}
		case "term_duration_in_seconds":
			savingPlan.Name = types.SavingsPlansFilterNameTerm
			savingPlan.Values = []string{fmt.Sprint(v.GetInt64Value())}
		case "start":
			savingPlan.Name = types.SavingsPlansFilterNameStart
			val := v.GetTimestampValue().AsTime()
			savingPlan.Values = []string{val.Format(time.RFC3339)}
		case "end_time":
			savingPlan.Name = types.SavingsPlansFilterNameEnd
			val := v.GetTimestampValue().AsTime()
			savingPlan.Values = []string{val.Format(time.RFC3339)}
		case "savings_plan_type":
			savingPlan.Name = types.SavingsPlansFilterNameSavingsPlanType
			savingPlan.Values = []string{v.GetStringValue()}
		default:
			addFilter = false
		}

		if addFilter {
			savingsPlanFilters = append(savingsPlanFilters, *savingPlan)
		}
	}

	return savingsPlanFilters
}
