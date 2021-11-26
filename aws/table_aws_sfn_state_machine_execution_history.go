package aws

import (
	"context"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsStepFunctionsStateMachineExecutionHistory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sfn_state_machine_execution_history",
		Description: "AWS Step Functions State Machine Execution History",
		List: &plugin.ListConfig{
			Hydrate:           listStepFunctionsStateMachineExecutionHistories,
			ShouldIgnoreError: isNotFoundError([]string{"ExecutionDoesNotExist", "InvalidParameter", "InvalidArn"}),
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "execution_arn",
					Require: plugin.Required,
				},
			},
		},
		GetMatrixItem: BuildRegionList,
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
				Transform:   transform.From(executionHistoryAkas),
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
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := d.KeyColumnQuals["execution_arn"].GetStringValue()
	accountId := commonColumnData.AccountId

	// check if the arn is empty or it contains a valid accountId
	if arn == "" || accountId != strings.Split(arn, ":")[4] {
		return nil, nil
	}

	params := &sfn.GetExecutionHistoryInput{
		ExecutionArn: aws.String(arn),
		MaxResults:   aws.Int64(1000),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxResults {
			params.MaxResults = limit
		}
	}

	err = svc.GetExecutionHistoryPages(
		params,
		func(page *sfn.GetExecutionHistoryOutput, isLast bool) bool {
			for _, events := range page.Events {
				d.StreamListItem(ctx, historyInfo{*events, arn})
				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listStepFunctionsStateMachineExecutionHistories", "ListExecutionsPages_error", err)
		return nil, err
	}

	return nil, nil
}

//// Transform Function

func executionHistoryAkas(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("executionHistoryAkas")
	history := d.HydrateItem.(historyInfo)

	// For State Machine, ARN format is arn:aws:states:us-east-1:632902152528:stateMachine:HelloWorld
	// For State Machine Execution, ARN format is arn:aws:states:us-east-1:632902152528:execution:HelloWorld:a44bc846-3601-fd75-63f7-60ac06a4ef97
	akas := []string{strings.Replace(history.ExecutionArn, "execution", "executionHistory", 1) + ":" + strconv.Itoa(int(*history.Id))}

	return akas, nil
}
