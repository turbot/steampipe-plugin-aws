package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/sfn"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsStepFunctionsStateMachine(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sfn_state_machine",
		Description: "AWS Step Functions State Machine",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("arn"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException", "StateMachineDoesNotExist"}),
			Hydrate:           getStepFunctionsStateMachine,
		},
		List: &plugin.ListConfig{
			Hydrate: listStepFunctionsStateManchines,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the state machine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that identifies the state machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StateMachineArn"),
			},
			{
				Name:        "status",
				Description: "The current status of the state machine.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getStepFunctionsStateMachine,
			},
			{
				Name:        "type",
				Description: "The type of the state machine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "The date the state machine is created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "definition",
				Description: "The Amazon States Language definition of the state machine.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getStepFunctionsStateMachine,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role used when creating this state machine.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getStepFunctionsStateMachine,
			},
			{
				Name:        "logging_configuration",
				Description: "The LoggingConfiguration data type is used to set CloudWatch Logs options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getStepFunctionsStateMachine,
			},
			{
				Name:        "tracing_configuration",
				Description: "Selects whether AWS X-Ray tracing is enabled.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getStepFunctionsStateMachine,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("StateMachineArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listStepFunctionsStateManchines(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := StepFunctionsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listStepFunctionsStateManchines", "connection_error", err)
		return nil, err
	}

	err = svc.ListStateMachinesPages(
		&sfn.ListStateMachinesInput{},
		func(page *sfn.ListStateMachinesOutput, isLast bool) bool {
			for _, stateMachine := range page.StateMachines {
				d.StreamListItem(ctx, stateMachine)
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listStepFunctionsStateManchines", "Error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getStepFunctionsStateMachine(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getStepFunctionsStateMachine")

	var arn string
	if h.Item != nil {
		arn = *h.Item.(*sfn.StateMachineListItem).StateMachineArn
	} else {
		arn = d.KeyColumnQuals["arn"].GetStringValue()
	}

	// Create Session
	svc, err := StepFunctionsService(ctx, d)
	if err != nil {
		logger.Error("getStepFunctionsStateMachine", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &sfn.DescribeStateMachineInput{
		StateMachineArn: &arn,
	}

	// Get call
	data, err := svc.DescribeStateMachine(params)
	if err != nil {
		logger.Error("getStepFunctionsStateMachine", "ERROR", err)
		return nil, err
	}

	return data, nil
}
