package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsBedrockCustomModel(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_bedrock_custom_model",
		Description: "AWS Bedrock Custom Model",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("model_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getBedrockCustomModel,
		},
		List: &plugin.ListConfig{
			Hydrate: listBedrockCustomModels,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_BEDROCK_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "base_model_arn",
				Description: "The base model Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "base_model_name",
				Description: "The base model name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "Creation time of the model.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "model_arn",
				Description: "The Amazon Resource Name (ARN) of the custom model.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "model_name",
				Description: "The name of the custom model.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customization_type",
				Description: "Specifies whether to carry out continued pre-training of a model or whether to fine-tune it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "model_status",
				Description: "The current status of the custom model.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_account_id",
				Description: "The unique identifier of the account that owns the model.",
				Type:        proto.ColumnType_STRING,
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

//// LIST FUNCTION

func listBedrockCustomModels(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BedrockClient(ctx, d)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("ValidationException: Unknown operation")) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_bedrock_custom_model.listBedrockCustomModels", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &bedrock.ListCustomModelsInput{
		MaxResults: &maxLimit,
	}

	paginator := bedrock.NewListCustomModelsPaginator(svc, input, func(o *bedrock.ListCustomModelsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("ValidationException: Unknown operation")) {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_bedrock_custom_model.listBedrockCustomModels", "api_error", err)
			return nil, err
		}

		for _, model := range output.ModelSummaries {
			d.StreamListItem(ctx, model)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getBedrockCustomModel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	modelId := d.EqualsQualString("model_arn")

	// Empty check
	if modelId == "" {
		return nil, nil
	}

	// Create service
	svc, err := BedrockClient(ctx, d)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("ValidationException: Unknown operation")) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_bedrock_custom_model.getBedrockCustomModel", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &bedrock.GetCustomModelInput{
		ModelIdentifier: &modelId,
	}

	// Get call
	data, err := svc.GetCustomModel(ctx, params)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("ValidationException: Unknown operation")) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_bedrock_custom_model.getBedrockCustomModel", "api_error", err)
		return nil, err
	}

	return data, nil
}
