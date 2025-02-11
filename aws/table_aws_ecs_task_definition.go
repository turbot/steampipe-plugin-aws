package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcsTaskDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_task_definition",
		Description: "AWS ECS Task Definition",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("task_definition_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterException", "ClientException"}),
			},
			Hydrate: getEcsTaskDefinition,
			Tags:    map[string]string{"service": "ecs", "action": "DescribeTaskDefinition"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEcsTaskDefinitions,
			Tags:    map[string]string{"service": "ecs", "action": "ListTaskDefinitions"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "family", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getEcsTaskDefinition,
				Tags: map[string]string{"service": "ecs", "action": "DescribeTaskDefinition"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ECS_SERVICE_ID),
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
				Name:        "deregistered_at",
				Description: "The Unix timestamp for the time when the task definition was deregistered.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.DeregisteredAt").Transform(transform.UnixToTimestamp).Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "ephemeral_storage_size_in_gib",
				Description: "The total amount, in GiB, of ephemeral storage to set for the task. The minimum supported value is 21 GiB and the maximum supported value is 200 GiB.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.EphemeralStorage.SizeInGiB"),
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
				Name:        "runtime_platform",
				Description: "The operating system that your task definitions are running on.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.RuntimePlatform"),
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
	svc, err := ECSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_task_definition.listEcsTaskDefinitions", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ecs.ListTaskDefinitionsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	equalQuala := d.EqualsQuals
	if equalQuala["family"] != nil {
		input.FamilyPrefix = aws.String(equalQuala["family"].GetStringValue())
	}
	if equalQuala["status"] != nil {
		input.Status = types.TaskDefinitionStatus(equalQuala["status"].GetStringValue())
	}

	paginator := ecs.NewListTaskDefinitionsPaginator(svc, input, func(o *ecs.ListTaskDefinitionsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ecs_task_definition.listEcsTaskDefinitions", "api_error", err)
			return nil, err
		}

		for _, items := range output.TaskDefinitionArns {
			d.StreamListItem(ctx, &ecs.DescribeTaskDefinitionOutput{
				TaskDefinition: &types.TaskDefinition{TaskDefinitionArn: aws.String(items)},
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEcsTaskDefinition(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var taskDefinitionArn string
	if h.Item != nil {
		taskDefinitionArn = *h.Item.(*ecs.DescribeTaskDefinitionOutput).TaskDefinition.TaskDefinitionArn
	} else {
		quals := d.EqualsQuals
		taskDefinitionArn = quals["task_definition_arn"].GetStringValue()
	}

	// Create Session
	svc, err := ECSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_task_definition.getEcsTaskDefinition", "connection_error", err)
		return nil, err
	}

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefinitionArn,
		Include: []types.TaskDefinitionField{
			types.TaskDefinitionFieldTags,
		},
	}

	op, err := svc.DescribeTaskDefinition(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_task_definition.getEcsTaskDefinition", "api_error", err)
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
