package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/savingsplans"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// https://github.com/turbot/steampipe-plugin-aws/issues/2500
// https://docs.aws.amazon.com/savingsplans/latest/APIReference/API_DescribeSavingsPlans.html
// https://docs.aws.amazon.com/savingsplans/latest/APIReference/API_DescribeSavingsPlansOfferings.html
//// TABLE DEFINITION
//
// {
// 	"nextToken": "string",
// 	"savingsPlans": [
// 	   {
// 		  "commitment": "string",
// 		  "currency": "string",
// 		  "description": "string",
// 		  "ec2InstanceFamily": "string",
// 		  "end": "string",
// 		  "offeringId": "string",
// 		  "paymentOption": "string",
// 		  "productTypes": [ "string" ],
// 		  "recurringPaymentAmount": "string",
// 		  "region": "string",
// 		  "returnableUntil": "string",
// 		  "savingsPlanArn": "string",
// 		  "savingsPlanId": "string",
// 		  "savingsPlanType": "string",
// 		  "start": "string",
// 		  "state": "string",
// 		  "tags": {
// 			 "string" : "string"
// 		  },
// 		  "termDurationInSeconds": number,
// 		  "upfrontPaymentAmount": "string"
// 	   }
// 	]
//  }

func tableAwsSavingsPlan(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_savings_plan",
		Description: "AWS Savings Plan",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("savings_plan_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				// ????
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			Hydrate: getSavingsPlan,
			Tags:    map[string]string{"service": "savingsplans", "action": "DescribeSavingsPlan"},
		},
		List: &plugin.ListConfig{
			Hydrate: listSavingsPlan,
			Tags:    map[string]string{"service": "savingsplans", "action": "DescribeSavingsPlans"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "savings_plan_type", Require: plugin.Optional},
				// {Name: "state", Require: plugin.Optional},
				// {Name: "payment_option", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				// TODO ???
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SAVINGSPLANS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "savings_plan_id",
				Description: "The ID of the Savings Plan.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SavingsPlanId"),
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
				Name:        "end",
				Description: "The end time of the Savings Plan.",
				Type:        proto.ColumnType_TIMESTAMP,
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
			// {
			// 	Name:        "region",
			// 	Description: "The AWS Region for the Savings Plan.",
			// 	Type:        proto.ColumnType_STRING,
			// },
			{
				Name:        "product_types",
				Description: "The product types supported by the Savings Plan.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the Savings Plan.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
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

func listSavingsPlan(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := SavingsPlansClient(ctx, d)

	fmt.Println("svc", svc)
	if err != nil {
		plugin.Logger(ctx).Error("aws_savings_plan.listSavingsPlan", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 20 {
				maxLimit = 20
			} else {
				maxLimit = limit
			}
		}
	}

	input := &savingsplans.DescribeSavingsPlansInput{
		MaxResults: aws.Int32(maxLimit),
	}

	// API doesn't support aws-go-sdk-v2 paginator as of date

	output, err := svc.DescribeSavingsPlans(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_savings_plan.listSavingsPlan", "api_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Info("aws_savings_plan.output", output)

	for _, item := range output.SavingsPlans {
		d.StreamListItem(ctx, item)
		fmt.Println("item", item)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSavingsPlan(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return nil, nil
}

// func getRDSReservedDBInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
// 	dbInstanceIdentifier := d.EqualsQuals["reserved_db_instance_id"].GetStringValue()

// 	// Create service
// 	svc, err := RDSClient(ctx, d)
// 	if err != nil {
// 		plugin.Logger(ctx).Error("aws_rds_reserved_db_instance.getRDSReservedDBInstance", "connection_error", err)
// 		return nil, err
// 	}

// 	params := &rds.DescribeReservedDBInstancesInput{
// 		ReservedDBInstanceId: aws.String(dbInstanceIdentifier),
// 	}

// 	op, err := svc.DescribeReservedDBInstances(ctx, params)
// 	if err != nil {
// 		plugin.Logger(ctx).Error("aws_rds_reserved_db_instance.getRDSReservedDBInstance", "api_error", err)
// 		return nil, err
// 	}

// 	if len(op.ReservedDBInstances) > 0 {
// 		return op.ReservedDBInstances[0], nil
// 	}
// 	return nil, nil
// }
