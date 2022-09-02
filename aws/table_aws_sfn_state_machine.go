package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsStepFunctionsStateMachine(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sfn_state_machine",
		Description: "AWS Step Functions State Machine",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException", "StateMachineDoesNotExist", "InvalidArn"}),
			},
			Hydrate: getStepFunctionsStateMachine,
		},
		List: &plugin.ListConfig{
			Hydrate: listStepFunctionsStateManchines,
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Name:        "tags_src",
				Description: "The list of tags associated with the state machine.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getStepFunctionStateMachineTags,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tracing_configuration",
				Description: "Selects whether AWS X-Ray tracing is enabled.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getStepFunctionsStateMachine,
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getStepFunctionStateMachineTags,
				Transform:   transform.From(stateMachineTagsToTurbotTags),
			},
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

	input := &sfn.ListStateMachinesInput{
		MaxResults: aws.Int64(1000),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			input.MaxResults = limit
		}
	}

	err = svc.ListStateMachinesPages(
		input,
		func(page *sfn.ListStateMachinesOutput, isLast bool) bool {
			for _, stateMachine := range page.StateMachines {
				d.StreamListItem(ctx, stateMachine)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
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

	if arn == "" {
		return nil, nil
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

func getStepFunctionStateMachineTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketTagging")
	stateMachineArn := getStateMachineArn(h.Item)

	// Empty Check
	if stateMachineArn == nil {
		return nil, nil
	}

	// Create Session
	svc, err := StepFunctionsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getStepFunctionStateMachineTags", "connection_error", err)
		return nil, err
	}

	params := &sfn.ListTagsForResourceInput{
		ResourceArn: stateMachineArn,
	}

	tags, err := svc.ListTagsForResource(params)
	if err != nil {
		plugin.Logger(ctx).Error("getStepFunctionStateMachineTags", err)
		return nil, err
	}

	return tags.Tags, nil
}

//// TRANSFORM FUNCTIONS

func stateMachineTagsToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("stateMachineTagsToTurbotTags")
	tags := d.HydrateItem.([]*sfn.Tag)

	if tags == nil {
		return nil, nil
	}
	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func getStateMachineArn(item interface{}) *string {
	switch item := item.(type) {
	case *sfn.StateMachineListItem:
		return item.StateMachineArn
	case *sfn.DescribeStateMachineOutput:
		return item.StateMachineArn
	}
	return nil
}
