package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsEcsTaskDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_task_definition",
		Description: "AWS ECS Task Definition",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("task_definition_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidParameterException", "ClientException"}),
			},
			Hydrate: getEcsTaskDefinition,
		},
		List: &plugin.ListConfig{
			Hydrate: listEcsTaskDefinitions,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "family", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "task_definition_arn",
				Description: "The Amazon Resource Name (ARN) that identifies the task definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TaskDefinition.TaskDefinitionArn"),
			},
			{
				Name:        "cpu",
				Description: "The number of cpu units used by the task.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Cpu"),
			},
			{
				Name:        "status",
				Description: "The status of the task definition.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Status"),
			},
			{
				Name:        "execution_role_arn",
				Description: "The Amazon Resource Name (ARN) of the task execution role that grants the Amazon ECS container agent permission to make AWS API calls on your behalf.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.ExecutionRoleArn"),
			},
			{
				Name:        "family",
				Description: "The name of a family that this task definition is registered to.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Family"),
			},
			{
				Name:        "ipc_mode",
				Description: "The IPC resource namespace to use for the containers in the task.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.IpcMode"),
			},
			{
				Name:        "memory",
				Description: "The amount (in MiB) of memory used by the task.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Memory"),
			},
			{
				Name:        "network_mode",
				Description: "The Docker networking mode to use for the containers in the task.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.NetworkMode"),
			},
			{
				Name:        "pid_mode",
				Description: "The process namespace to use for the containers in the task.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.PidMode"),
			},
			{
				Name:        "revision",
				Description: "The revision of the task in a particular family.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Revision"),
			},
			{
				Name:        "task_role_arn",
				Description: "The short name or full Amazon Resource Name (ARN) of the AWS Identity and Access Management (IAM) role that grants containers in the task permission to call AWS APIs on your behalf.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.TaskRoleArn"),
			},
			{
				Name:        "registered_at",
				Description: "The Unix timestamp for when the task definition was registered.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.RegisteredAt"),
			},
			{
				Name:        "registered_by",
				Description: "The principal that registered the task definition.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.RegisteredBy"),
			},
			{
				Name:        "container_definitions",
				Description: "A list of container definitions in JSON format that describe the different containers that make up your task.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.ContainerDefinitions"),
			},
			{
				Name:        "compatibilities",
				Description: "The launch type to use with your task.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Compatibilities"),
			},
			{
				Name:        "inference_accelerators",
				Description: "The Elastic Inference accelerator associated with the task.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.InferenceAccelerators"),
			},
			{
				Name:        "placement_constraints",
				Description: "An array of placement constraint objects to use for tasks.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.PlacementConstraints"),
			},
			{
				Name:        "proxy_configuration",
				Description: "The configuration details for the App Mesh proxy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.ProxyConfiguration"),
			},
			{
				Name:        "requires_attributes",
				Description: "The container instance attributes required by your task.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.RequiresAttributes"),
			},
			{
				Name:        "requires_compatibilities",
				Description: "The launch type the task requires. If no value is specified, it will default to EC2. Valid values include EC2 and FARGATE.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.RequiresCompatibilities"),
			},
			{
				Name:        "volumes",
				Description: "The list of volume definitions for the task.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Volumes"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with task.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromP(getAwsEcsTaskDefinitionTurbotData, "Title"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromP(getAwsEcsTaskDefinitionTurbotData, "Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.TaskDefinitionArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listEcsTaskDefinitions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EcsService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &ecs.ListTaskDefinitionsInput{
		MaxResults: aws.Int64(100),
	}

	equalQuala := d.KeyColumnQuals
	if equalQuala["family"] != nil {
		input.FamilyPrefix = aws.String(equalQuala["family"].GetStringValue())
	}
	if equalQuala["status"] != nil {
		input.Status = aws.String(equalQuala["status"].GetStringValue())
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.ListTaskDefinitionsPages(
		input,
		func(page *ecs.ListTaskDefinitionsOutput, isLast bool) bool {
			for _, result := range page.TaskDefinitionArns {
				d.StreamListItem(ctx, &ecs.DescribeTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						TaskDefinitionArn: result,
					},
				})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	return nil, err

}

//// HYDRATE FUNCTIONS

func getEcsTaskDefinition(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getEcsTaskDefinition")

	var taskDefinitionArn string
	if h.Item != nil {
		taskDefinitionArn = *h.Item.(*ecs.DescribeTaskDefinitionOutput).TaskDefinition.TaskDefinitionArn
	} else {
		quals := d.KeyColumnQuals
		taskDefinitionArn = quals["task_definition_arn"].GetStringValue()
	}

	// Create Session
	svc, err := EcsService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefinitionArn,
		Include:        []*string{aws.String("TAGS")},
	}

	op, err := svc.DescribeTaskDefinition(params)
	if err != nil {
		logger.Debug("getEcsTaskDefinition", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getAwsEcsTaskDefinitionTurbotData(_ context.Context, d *transform.TransformData) (interface{},
	error) {
	param := d.Param.(string)
	ecsTaskDefinition := d.HydrateItem.(*ecs.DescribeTaskDefinitionOutput)

	// Get resource title
	arn := ecsTaskDefinition.TaskDefinition.TaskDefinitionArn
	splitArn := strings.Split(*arn, "/")
	title := splitArn[len(splitArn)-1]

	if param == "Tags" {
		if ecsTaskDefinition.Tags == nil {
			return nil, nil
		}

		// Get the resource tags
		if ecsTaskDefinition.Tags != nil {
			turbotTagsMap := map[string]string{}
			for _, i := range ecsTaskDefinition.Tags {
				turbotTagsMap[*i.Key] = *i.Value
			}
			return turbotTagsMap, nil
		}
	}

	return title, nil
}
