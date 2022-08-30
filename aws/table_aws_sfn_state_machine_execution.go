package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsStepFunctionsStateMachineExecution(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sfn_state_machine_execution",
		Description: "AWS Step Functions State Machine Execution",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("execution_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidParameter", "ExecutionDoesNotExist", "InvalidArn"}),
			},
			Hydrate: getStepFunctionsStateMachineExecution,
		},
		List: &plugin.ListConfig{
			Hydrate:       listStepFunctionsStateMachineExecutions,
			ParentHydrate: listStepFunctionsStateManchines,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "status", Require: plugin.Optional},
				{Name: "state_machine_arn", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the execution.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "execution_arn",
				Description: "The Amazon Resource Name (ARN) that identifies the execution.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current status of the execution.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "input",
				Description: "The string that contains the JSON input data of the execution.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getStepFunctionsStateMachineExecution,
			},
			{
				Name:        "output",
				Description: "The JSON output data of the execution.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getStepFunctionsStateMachineExecution,
			},
			{
				Name:        "start_date",
				Description: "The date the execution started.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "state_machine_arn",
				Description: "The Amazon Resource Name (ARN) of the executed state machine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stop_date",
				Description: "If the execution already ended, the date the execution stopped.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "trace_header",
				Description: "The AWS X-Ray trace header that was passed to the execution.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getStepFunctionsStateMachineExecution,
			},
			{
				Name:        "input_details",
				Description: "Provides details about execution input or output.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getStepFunctionsStateMachineExecution,
			},
			{
				Name:        "output_details",
				Description: "Provides details about execution input or output.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getStepFunctionsStateMachineExecution,
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
				Transform:   transform.FromField("ExecutionArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listStepFunctionsStateMachineExecutions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := StepFunctionsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listStepFunctionsStateMachineExecutions", "connection_error", err)
		return nil, err
	}

	arn := h.Item.(*sfn.StateMachineListItem).StateMachineArn

	equalQuals := d.KeyColumnQuals
	// Minimize the API call with the given layer name
	if equalQuals["state_machine_arn"] != nil {
		if equalQuals["state_machine_arn"].GetStringValue() != "" {
			if equalQuals["state_machine_arn"].GetStringValue() != "" && equalQuals["state_machine_arn"].GetStringValue() != *arn {
				return nil, nil
			}
		} else if len(getListValues(equalQuals["state_machine_arn"].GetListValue())) > 0 {
			if !strings.Contains(fmt.Sprint(getListValues(equalQuals["state_machine_arn"].GetListValue())), *arn) {
				return nil, nil
			}
		}
	}

	input := &sfn.ListExecutionsInput{
		StateMachineArn: arn,
		MaxResults:      aws.Int64(1000),
	}

	if equalQuals["status"] != nil {
		input.StatusFilter = aws.String(equalQuals["status"].GetStringValue())
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			input.MaxResults = limit
		}
	}

	err = svc.ListExecutionsPages(
		input,
		func(page *sfn.ListExecutionsOutput, isLast bool) bool {
			for _, execution := range page.Executions {
				d.StreamListItem(ctx, execution)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listStepFunctionsStateMachineExecutions", "ListExecutionsPages_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getStepFunctionsStateMachineExecution(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getStepFunctionsStateMachineExecution")

	var arn string
	if h.Item != nil {
		arn = *h.Item.(*sfn.ExecutionListItem).ExecutionArn
	} else {
		arn = d.KeyColumnQuals["execution_arn"].GetStringValue()
	}

	// Create Session
	svc, err := StepFunctionsService(ctx, d)
	if err != nil {
		logger.Error("getStepFunctionsStateMachineExecution", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &sfn.DescribeExecutionInput{
		ExecutionArn: &arn,
	}

	// Get call
	data, err := svc.DescribeExecution(params)
	if err != nil {
		logger.Error("getStepFunctionsStateMachineExecution", "DescribeExecution_error", err)
		return nil, err
	}

	return data, nil
}
