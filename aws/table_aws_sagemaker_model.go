package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"

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
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ValidationException", "NotFoundException", "RecordNotFound"}),
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
	// Create Session
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_model.listAwsSageMakerModels", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Check for unsupported region
		return nil, nil
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

	input := &sagemaker.ListModelsInput{
		MaxResults: aws.Int32(maxLimit),
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

	paginator := sagemaker.NewListModelsPaginator(svc, input, func(o *sagemaker.ListModelsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sagemaker_model.listAwsSageMakerModels", "api_error", err)
			return nil, err
		}

		for _, items := range output.Models {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

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
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_model.getAwsSageMakerModel", "connection_error", err)
		return nil, err
	}
	
	if svc == nil {
		// Check for unsupported region
		return nil, nil
	}

	// Build the params
	params := &sagemaker.DescribeModelInput{
		ModelName: aws.String(name),
	}

	// Get call
	data, err := svc.DescribeModel(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_model.getAwsSageMakerModel", "api_error", err)
		return nil, err
	}
	return data, nil
}

func listAwsSageMakerModelTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var modelArn string
	switch h.Item.(type) {
	case types.ModelSummary:
		modelArn = *h.Item.(types.ModelSummary).ModelArn
	case *sagemaker.DescribeModelOutput:
		modelArn = *h.Item.(*sagemaker.DescribeModelOutput).ModelArn
	}

	// Create Session
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_model.listAwsSageMakerModelTags", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(modelArn),
	}

	pagesLeft := true
	tags := []types.Tag{}
	for pagesLeft {
		keyTags, err := svc.ListTags(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sagemaker_model.listAwsSageMakerModelTags", "api_error", err)
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
	case types.ModelSummary:
		return *item.ModelName
	case *sagemaker.DescribeModelOutput:
		return *item.ModelName
	}
	return ""
}
