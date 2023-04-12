package aws

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/aws/smithy-go"

	sfnv1 "github.com/aws/aws-sdk-go/service/sfn"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsStepFunctionsStateMachineExecutionHistory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sfn_state_machine_execution_history",
		Description: "AWS Step Functions State Machine Execution History",
		List: &plugin.ListConfig{
			Hydrate:       listStepFunctionsStateMachineExecutionHistories,
			ParentHydrate: listStepFunctionsStateMachines,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "execution_arn", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(sfnv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The id of the event.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "execution_arn",
				Description: "The Amazon Resource Name (ARN) that identifies the execution.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "previous_event_id",
				Description: "The id of the previous event.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timestamp",
				Description: "The date and time the event occurred.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "type",
				Description: "The type of the event.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "activity_failed_event_details",
				Description: "Contains details about an activity that failed during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "activity_schedule_failed_event_details",
				Description: "Contains details about an activity schedule event that failed during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "activity_scheduled_event_details",
				Description: "Contains details about an activity scheduled during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "activity_started_event_details",
				Description: "Contains details about the start of an activity during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "activity_succeeded_event_details",
				Description: "Contains details about an activity that successfully terminated during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "activity_timed_out_event_details",
				Description: "Contains details about an activity timeout that occurred during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "execution_aborted_event_details",
				Description: "Contains details about an abort of an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "execution_failed_event_details",
				Description: "Contains details about an execution failure event.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "execution_started_event_details",
				Description: "Contains details about the start of the execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "execution_succeeded_event_details",
				Description: "Contains details about the successful termination of the execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "execution_timed_out_event_details",
				Description: "Contains details about the execution timeout that occurred during the execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lambda_function_failed_event_details",
				Description: "Contains details about a lambda function that failed during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lambda_function_schedule_failed_event_details",
				Description: "Contains details about a failed lambda function schedule event that occurred during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lambda_function_scheduled_event_details",
				Description: "Contains details about a lambda function scheduled during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lambda_function_start_failed_event_details",
				Description: "Contains details about a lambda function that failed to start during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lambda_function_succeeded_event_details",
				Description: "Contains details about a lambda function that terminated successfully during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lambda_function_timed_out_event_details",
				Description: "Contains details about a lambda function timeout that occurred during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "map_iteration_aborted_event_details",
				Description: "Contains details about an iteration of a Map state that was aborted.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "map_iteration_failed_event_details",
				Description: "Contains details about an iteration of a Map state that failed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "map_iteration_started_event_details",
				Description: "Contains details about an iteration of a Map state that was started.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "map_iteration_succeeded_event_details",
				Description: "Contains details about an iteration of a Map state that succeeded.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "map_state_started_event_details",
				Description: "Contains details about Map state that was started.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "state_entered_event_details",
				Description: "Contains details about a state entered during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "state_exited_event_details",
				Description: "Contains details about an exit from a state during an execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "task_failed_event_details",
				Description: "Contains details about the failure of a task.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "task_scheduled_event_details",
				Description: "Contains details about a task that was scheduled.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "task_start_failed_event_details",
				Description: "Contains details about a task that failed to start.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "task_started_event_details",
				Description: "Contains details about a task that was started.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "task_submit_failed_event_details",
				Description: "Contains details about a task that where the submit failed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "task_submitted_event_details",
				Description: "Contains details about a submitted task.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "task_succeeded_event_details",
				Description: "Contains details about a task that succeeded.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "task_timed_out_event_details",
				Description: "Contains details about a task that timed out.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(executionHistoryArn),
			},
		}),
	}
}

type historyInfo struct {
	types.HistoryEvent
	ExecutionArn string
}

//// LIST FUNCTION

func listStepFunctionsStateMachineExecutionHistories(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := StepFunctionsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sfn_state_machine_execution_history.listStepFunctionsStateMachineExecutionHistories", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	var executions []string

	plugin.Logger(ctx).Trace("aws_sfn_state_machine_execution_history.listStepFunctionsStateMachineExecutionHistories", fmt.Sprintf("d.Quals=%#v", d.Quals))
	executionArnQuals := getQualsValueByColumn(d.Quals, "execution_arn", "string") // FIXME: this does not filter on operators at all, so issues with other operators than `=` or `in` would occur, e.g. with `like`
	plugin.Logger(ctx).Debug("aws_sfn_state_machine_execution_history.listStepFunctionsStateMachineExecutionHistories", "execution_arn quals", executionArnQuals)

	// Minimize the API call with the given execution ARN
	if executionArnQuals != nil {
		if executionArnQualsStr, ok := executionArnQuals.(string); ok && executionArnQualsStr != "" {
			executions = []string{executionArnQualsStr}

		} else if executionArnQualsList, ok := executionArnQuals.([]string); ok && len(executionArnQualsList) > 0 {
			executions = executionArnQualsList
		}
	} else {
		stateMachineArn := h.Item.(types.StateMachineListItem).StateMachineArn
		maxLimit := int32(1000)
		input := &sfn.ListExecutionsInput{
			MaxResults:      maxLimit,
			StateMachineArn: stateMachineArn,
		}
		paginator := sfn.NewListExecutionsPaginator(svc, input, func(o *sfn.ListExecutionsPaginatorOptions) {
			o.Limit = maxLimit
			o.StopOnDuplicateToken = true
		})
		// List call
		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				plugin.Logger(ctx).Error("aws_sfn_state_machine_execution_history.listStepFunctionsStateMachineExecutionHistories", "api_error", err)
				return nil, err
			}
			for _, execution := range output.Executions {
				executions = append(executions, *execution.ExecutionArn)
			}
		}
		if err != nil {
			plugin.Logger(ctx).Error("aws_sfn_state_machine_execution_history.listStepFunctionsStateMachineExecutionHistories", "api_error", err)
			return nil, err
		}
	}

	// Iterating all the available executions matching the query quals, if any
	for _, executionArn := range executions {
		historyEvents, err := getRowDataForExecutionHistory(ctx, d, executionArn)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sfn_state_machine_execution_history.listStepFunctionsStateMachineExecutionHistories", "api_error", err)
			return nil, err
		}
		for _, event := range historyEvents {
			d.StreamLeafListItem(ctx, historyInfo{event.HistoryEvent, event.ExecutionArn})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getRowDataForExecutionHistory(ctx context.Context, d *plugin.QueryData, arn string) ([]historyInfo, error) {
	// Create session
	svc, err := StepFunctionsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sfn_state_machine_execution_history.getRowDataForExecutionHistory", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxLimit := int32(1000)
	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(maxLimit) {
			maxLimit = int32(*limit)
		}
	}

	var items []historyInfo

	input := &sfn.GetExecutionHistoryInput{
		MaxResults:   maxLimit,
		ExecutionArn: aws.String(arn),
	}
	paginator := sfn.NewGetExecutionHistoryPaginator(svc, input, func(o *sfn.GetExecutionHistoryPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})
	// List call
	for paginator.HasMorePages() {
		plugin.Logger(ctx).Trace("aws_sfn_state_machine_execution_history.getRowDataForExecutionHistory", "api_call GetExecutionHistory", arn)
		output, err := paginator.NextPage(ctx)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				switch apiErr.(type) {
				case *types.ExecutionDoesNotExist:
					// Ignore expired executions for which history is no longer available
					plugin.Logger(ctx).Trace("aws_sfn_state_machine_execution_history.getRowDataForExecutionHistory", "api_error ignore_expired", err)
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_sfn_state_machine_execution_history.getRowDataForExecutionHistory", "api_error", err)
			return nil, err
		}

		for _, event := range output.Events {
			items = append(items, historyInfo{event, arn})
		}
	}

	return items, nil
}

//// Transform Function

func executionHistoryArn(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	history := d.HydrateItem.(historyInfo)

	// For State Machine, ARN format is arn:aws:states:us-east-1:632902152528:stateMachine:HelloWorld
	// For State Machine Execution, ARN format is arn:aws:states:us-east-1:632902152528:execution:HelloWorld:a44bc846-3601-fd75-63f7-60ac06a4ef97
	akas := []string{strings.Replace(history.ExecutionArn, "execution", "executionHistory", 1) + ":" + strconv.Itoa(int(history.Id))}

	return akas, nil
}
