package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSageMakerModel(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_model",
		Description: "AWS Sagemaker Model",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"ValidationException", "NotFoundException", "RecordNotFound"}),
			Hydrate:           getAwsSageMakerModel,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSageMakerModels,
		},
		GetMatrixItem: BuildRegionList,
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
				Transform:   transform.FromField("Tags"),
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
				Transform:   transform.FromField("Tags").Transform(sageMakerModelTurbotTags),
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
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsSageMakerModels", "AWS_REGION", region)

	// Create Session
	svc, err := SageMakerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListModelsPages(
		&sagemaker.ListModelsInput{},
		func(page *sagemaker.ListModelsOutput, isLast bool) bool {
			for _, model := range page.Models {
				d.StreamListItem(ctx, model)
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsSageMakerModel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var name string
	if h.Item != nil {
		name = modelName(h.Item)
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create service
	svc, err := SageMakerService(ctx, d, region)
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

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var modelArn string
	switch h.Item.(type) {
	case *sagemaker.ModelSummary:
		modelArn = *h.Item.(*sagemaker.ModelSummary).ModelArn
	case *sagemaker.DescribeModelOutput:
		modelArn = *h.Item.(*sagemaker.DescribeModelOutput).ModelArn
	}

	// Create Session
	svc, err := SageMakerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(modelArn),
	}

	// Get call
	op, err := svc.ListTags(params)
	if err != nil {
		logger.Debug("listAwsSageMakerModelTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTION

func sageMakerModelTurbotTags(_ context.Context, d *transform.TransformData) (interface{},
	error) {
	data := d.HydrateItem.(*sagemaker.ListTagsOutput)

	if data.Tags == nil {
		return nil, nil
	}

	if data.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range data.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}

func modelName(item interface{}) string {
	switch item.(type) {
	case *sagemaker.ModelSummary:
		return *item.(*sagemaker.ModelSummary).ModelName
	case *sagemaker.DescribeModelOutput:
		return *item.(*sagemaker.DescribeModelOutput).ModelName
	}
	return ""
}
