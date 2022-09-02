package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSageMakerModel(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_model",
		Description: "AWS Sagemaker Model",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationException", "NotFoundException", "RecordNotFound"}),
			},
			Hydrate: getAwsSageMakerModel,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSageMakerModels,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "creation_time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the model.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ModelName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the model.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ModelArn"),
			},
			{
				Name:        "creation_time",
				Description: "A timestamp that indicates when the model was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "enable_network_isolation",
				Description: "If True, no inbound or outbound network calls can be made to or from the model container.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsSageMakerModel,
			},
			{
				Name:        "execution_role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role that you specified for the model.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerModel,
			},
			{
				Name:        "containers",
				Description: "The containers in the inference pipeline.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerModel,
			},
			{
				Name:        "inference_execution_config",
				Description: "Specifies details of how containers in a multi-container endpoint are called.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerModel,
			},
			{
				Name:        "primary_container",
				Description: "The location of the primary inference code, associated artifacts, and custom environment map that the inference code uses when it is deployed in production.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerModel,
			},
			{
				Name:        "vpc_config",
				Description: "A VpcConfig object that specifies the VPC that this model has access to.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerModel,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the model.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSageMakerModelTags,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ModelName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSageMakerModelTags,
				Transform:   transform.FromValue().Transform(sageMakerTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ModelArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSageMakerModels(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsSageMakerModels")

	// Create Session
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &sagemaker.ListModelsInput{
		MaxResults: aws.Int64(100),
	}

	quals := d.Quals
	if quals["creation_time"] != nil {
		for _, q := range quals["creation_time"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">=", ">":
				input.CreationTimeAfter = aws.Time(timestamp)
			case "<", "<=":
				input.CreationTimeBefore = aws.Time(timestamp)
			}
		}
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
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
	err = svc.ListModelsPages(
		input,
		func(page *sagemaker.ListModelsOutput, isLast bool) bool {
			for _, model := range page.Models {
				d.StreamListItem(ctx, model)

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

func getAwsSageMakerModel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		name = modelName(h.Item)
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create service
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &sagemaker.DescribeModelInput{
		ModelName: aws.String(name),
	}

	// Get call
	data, err := svc.DescribeModel(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAwsSageMakerModel", "ERROR", err)
		return nil, err
	}
	return data, nil
}

func listAwsSageMakerModelTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listAwsSageMakerModelTags")

	var modelArn string
	switch h.Item.(type) {
	case *sagemaker.ModelSummary:
		modelArn = *h.Item.(*sagemaker.ModelSummary).ModelArn
	case *sagemaker.DescribeModelOutput:
		modelArn = *h.Item.(*sagemaker.DescribeModelOutput).ModelArn
	}

	// Create Session
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(modelArn),
	}

	pagesLeft := true
	tags := []*sagemaker.Tag{}
	for pagesLeft {
		keyTags, err := svc.ListTags(params)
		if err != nil {
			plugin.Logger(ctx).Error("listAwsSageMakerModelTags", "ListTags_error", err)
			return nil, err
		}
		tags = append(tags, keyTags.Tags...)

		if keyTags.NextToken != nil {
			params.NextToken = keyTags.NextToken
		} else {
			pagesLeft = false
		}
	}

	return tags, nil
}

//// TRANSFORM FUNCTION

func modelName(item interface{}) string {
	switch item := item.(type) {
	case *sagemaker.ModelSummary:
		return *item.ModelName
	case *sagemaker.DescribeModelOutput:
		return *item.ModelName
	}
	return ""
}
