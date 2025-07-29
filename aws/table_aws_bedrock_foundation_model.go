package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsBedrockFoundationModel(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_bedrock_foundation_model",
		Description: "AWS Bedrock Foundation Model",
		List: &plugin.ListConfig{
			Hydrate: listBedrockFoundationModels,
			Tags:    map[string]string{"service": "bedrock", "action": "ListFoundationModels"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("model_id"),
			Hydrate:    getBedrockFoundationModel,
			Tags:       map[string]string{"service": "bedrock", "action": "GetFoundationModel"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_BEDROCK_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the foundation model.",
				Transform:   transform.FromField("ModelArn"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "model_id",
				Description: "The unique identifier of the foundation model.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "model_name",
				Description: "The name of the foundation model.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provider_name",
				Description: "The name of the model provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "input_modalities",
				Description: "The input modalities supported by the model.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "output_modalities",
				Description: "The output modalities supported by the model.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "customizations_supported",
				Description: "The customizations supported by the model.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "inference_types_supported",
				Description: "The inference types supported by the model.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "model_lifecycle",
				Description: "The lifecycle status of the model.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "response_streaming_supported",
				Description: "Whether the model supports response streaming.",
				Type:        proto.ColumnType_BOOL,
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ModelName"),
			},
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ModelArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

func listBedrockFoundationModels(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create AWS service client
	svc, err := BedrockClient(ctx, d)

	if svc == nil {
		return nil, nil
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_foundation_model.listBedrockFoundationModels", "connection_error", err)
		return nil, err
	}

	// Execute list call
	resp, err := svc.ListFoundationModels(ctx, &bedrock.ListFoundationModelsInput{})
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_foundation_model.listBedrockFoundationModels", "api_error", err)
		return nil, err
	}

	for _, model := range resp.ModelSummaries {
		d.StreamListItem(ctx, model)
	}

	return nil, nil
}

func getBedrockFoundationModel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	modelId := d.EqualsQualString("model_id")

	if modelId == "" {
		return nil, nil
	}

	// Create AWS service client
	svc, err := BedrockClient(ctx, d)

	if svc == nil {
		return nil, nil
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_foundation_model.getBedrockFoundationModel", "connection_error", err)
		return nil, err
	}

	// Execute get call
	resp, err := svc.GetFoundationModel(ctx, &bedrock.GetFoundationModelInput{
		ModelIdentifier: &modelId,
	})
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_foundation_model.getBedrockFoundationModel", "api_error", err)
		return nil, err
	}

	return resp.ModelDetails, nil
}
