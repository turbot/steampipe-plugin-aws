package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEcsTaskDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_task_definition",
		Description: "AWS ECS Task Definition",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("task_definition_arn"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidParameterException"}),
			Hydrate:           getEcsTaskDefinition,
		},
		List: &plugin.ListConfig{
			Hydrate: listEcsTaskDefinitions,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "task_definition_arn",
				Description: "The Amazon Resource Name (ARN) that identifies the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TaskDefinition.TaskDefinitionArn"),
			},
			{
				Name:        "compatibilities",
				Description: "The launch type to use with your task.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Compatibilities"),
			},
			{
				Name:        "container_definitions",
				Description: "A list of container definitions in JSON format that describe the different containers that make up your task.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.ContainerDefinitions"),
			},
			{
				Name:        "cpu",
				Description: "The number of cpu units used by the task.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Cpu"),
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
				Name:        "inference_accelerators",
				Description: "The Elastic Inference accelerator associated with the task.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.InferenceAccelerators"),
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
				Type:        proto.ColumnType_STRING,
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
				Name:        "revision",
				Description: "The launch type the task requires. If no value is specified, it will default to EC2. Valid values include EC2 and FARGATE.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Revision"),
			},
			{
				Name:        "status",
				Description: "The status of the task definition.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Status"),
			},
			{
				Name:        "task_role_arn",
				Description: "The short name or full Amazon Resource Name (ARN) of the AWS Identity and Access Management (IAM) role that grants containers in the task permission to call AWS APIs on your behalf.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.TaskRoleArn"),
			},
			{
				Name:        "volumes",
				Description: "The list of volume definitions for the task.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.Volumes"),
			},
			{
				Name:        "registered_at",
				Description: "The list of volume definitions for the task.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.RegisteredAt"),
			},
			{
				Name:        "registered_by",
				Description: "The list of volume definitions for the task.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.RegisteredBy"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with target group.",
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
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromField("TaskDefinition.TaskDefinitionArn").Transform(arnToAkas),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsTaskDefinition,
				Transform:   transform.FromP(getAwsEcsTaskDefinitionTurbotData, "Tags"),
			},
		}),
	}
}

//// LIST FUNCTION

func listEcsTaskDefinitions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := EcsService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	params := &ecs.ListTaskDefinitionsInput{}
	pagesLeft := true

	for pagesLeft {
		result, err := svc.ListTaskDefinitions(params)
		if err != nil {
			return nil, err
		}

		for _, results := range result.TaskDefinitionArns {
			d.StreamListItem(ctx, &ecs.DescribeTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					TaskDefinitionArn: results,
				},
			})
		}

		if result.NextToken != nil {
			pagesLeft = true
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEcsTaskDefinition(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getEcsTaskDefinition")
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var taskDefinitionArn string
	if h.Item != nil {
		taskDefinitionArn = *h.Item.(*ecs.DescribeTaskDefinitionOutput).TaskDefinition.TaskDefinitionArn
	} else {
		quals := d.KeyColumnQuals
		taskDefinitionArn = quals["task_definition_arn"].GetStringValue()
	}

	// Create Session
	svc, err := EcsService(ctx, d, region)
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

	if op != nil {
		return op, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getAwsEcsTaskDefinitionTurbotData(ctx context.Context, d *transform.TransformData) (interface{},
	error) {

	param := d.Param.(string)
	if param == "Tags" {

		ecsTaskDefinitionTags := d.HydrateItem.(*ecs.DescribeTaskDefinitionOutput).Tags

		if ecsTaskDefinitionTags == nil {
			return nil, nil
		}

		if ecsTaskDefinitionTags != nil {
			turbotTagsMap := map[string]string{}
			for _, i := range ecsTaskDefinitionTags {
				turbotTagsMap[*i.Key] = *i.Value
			}
			return turbotTagsMap, nil
		}
	}

	taskDefinition := d.HydrateItem.(*ecs.DescribeTaskDefinitionOutput).TaskDefinition

	// Get resource title
	arn := taskDefinition.TaskDefinitionArn

	title := strings.Split(*arn, "/")[len(strings.Split(*arn, "/"))-1]
	return title, nil
}
