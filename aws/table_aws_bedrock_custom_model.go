package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
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
			KeyColumns: []*plugin.KeyColumn{
				{Name: "arn", Require: plugin.AnyOf},
				{Name: "model_name", Require: plugin.AnyOf},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getBedrockCustomModel,
		},
		List: &plugin.ListConfig{
			Hydrate: listBedrockCustomModels,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "base_model_arn", Require: plugin.Optional},
				{Name: "creation_time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
			},
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
				Transform:   transform.FromField("BaseModelName"),
			},
			{
				Name:        "creation_time",
				Description: "Creation time of the model.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the custom model.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ModelArn"),
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
		plugin.Logger(ctx).Error("aws_bedrock_custom_model.listBedrockCustomModels", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &bedrock.ListCustomModelsInput{
		MaxResults: &maxLimit,
	}

	// Apply optional quals if provided
	if d.EqualsQuals["base_model_arn"] != nil {
		input.BaseModelArnEquals = aws.String(d.EqualsQuals["base_model_arn"].GetStringValue())
	}

	if d.Quals["creation_time"] != nil {
		for _, q := range d.Quals["creation_time"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">=", ">":
				input.CreationTimeAfter = &timestamp
			case "<=", "<":
				input.CreationTimeBefore = &timestamp
			}
		}
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
	var modelIdentifier string

	if d.EqualsQualString("arn") != "" {
		modelIdentifier = d.EqualsQualString("arn")
	} else {
		modelIdentifier = d.EqualsQualString("model_name")
	}

	// Empty check
	if modelIdentifier == "" {
		return nil, nil
	}

	// Create service
	svc, err := BedrockClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_custom_model.getBedrockCustomModel", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &bedrock.GetCustomModelInput{
		ModelIdentifier: &modelIdentifier,
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
