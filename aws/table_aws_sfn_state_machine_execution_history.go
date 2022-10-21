package aws

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsStepFunctionsStateMachineExecutionHistory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sfn_state_machine_execution_history",
		Description: "AWS Step Functions State Machine Execution History",
		List: &plugin.ListConfig{
			Hydrate:       listStepFunctionsStateMachineExecutionHistories,
			ParentHydrate: listStepFunctionsStateManchines,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	sfn.HistoryEvent
	ExecutionArn string
}

//// LIST FUNCTION

func listStepFunctionsStateMachineExecutionHistories(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := StepFunctionsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listStepFunctionsStateMachineExecutionHistories", "connection_error", err)
		return nil, err
	}

	stateMachineArn := h.Item.(*sfn.StateMachineListItem).StateMachineArn
	var executions []sfn.ExecutionListItem

	input := &sfn.ListExecutionsInput{
		MaxResults:      types.Int64(1000),
		StateMachineArn: stateMachineArn,
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			input.MaxResults = limit
		}
	}

	// List call
	err = svc.ListExecutionsPages(
		input,
		func(page *sfn.ListExecutionsOutput, isLast bool) bool {
			for _, execution := range page.Executions {
				executions = append(executions, *execution)
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listStepFunctionsStateMachineExecutionHistories", "ListExecutionsPages_error", err)
		return nil, err
	}

	var wg sync.WaitGroup
	executionCh := make(chan []historyInfo, len(executions))
	errorCh := make(chan error, len(executions))

	// Iterating all the available executions
	for _, item := range executions {
		wg.Add(1)
		go getRowDataForExecutionHistoryAsync(ctx, d, *item.ExecutionArn, &wg, executionCh, errorCh)
	}

	// wait for all executions to be processed
	wg.Wait()
	close(executionCh)
	close(errorCh)

	for err := range errorCh {
		plugin.Logger(ctx).Error("listStepFunctionsStateMachineExecutionHistories", "getRowDataForExecutionHistoryAsync_error", err)
		return nil, err
	}

	for item := range executionCh {
		for _, data := range item {
			d.StreamLeafListItem(ctx, historyInfo{data.HistoryEvent, data.ExecutionArn})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getRowDataForExecutionHistoryAsync(ctx context.Context, d *plugin.QueryData, arn string, wg *sync.WaitGroup, executionCh chan []historyInfo, errorCh chan error) {
	defer wg.Done()

	rowData, err := getRowDataForExecutionHistory(ctx, d, arn)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		executionCh <- rowData
	}
}

func getRowDataForExecutionHistory(ctx context.Context, d *plugin.QueryData, arn string) ([]historyInfo, error) {
	// Create session
	svc, err := StepFunctionsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getRowDataForExecutionHistory", "connection_error", err)
		return nil, err
	}

	params := &sfn.GetExecutionHistoryInput{
		ExecutionArn: types.String(arn),
	}

	var items []historyInfo

	listHistory, err := svc.GetExecutionHistory(params)
	if err != nil {
		plugin.Logger(ctx).Error("getRowDataForExecutionHistory", "GetExecutionHistory_error", err)
		return nil, err
	}

	for _, event := range listHistory.Events {
		items = append(items, historyInfo{*event, arn})
	}

	return items, nil
}

//// Transform Function

func executionHistoryArn(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("executionHistoryArn")
	history := d.HydrateItem.(historyInfo)

	// For State Machine, ARN format is arn:aws:states:us-east-1:632902152528:stateMachine:HelloWorld
	// For State Machine Execution, ARN format is arn:aws:states:us-east-1:632902152528:execution:HelloWorld:a44bc846-3601-fd75-63f7-60ac06a4ef97
	akas := []string{strings.Replace(history.ExecutionArn, "execution", "executionHistory", 1) + ":" + strconv.Itoa(int(*history.Id))}

	return akas, nil
}
