package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAwsBedrockImportedModel(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_bedrock_imported_model",
		Description: "AWS Bedrock Imported Model",
		List: &plugin.ListConfig{
			Hydrate: listBedrockImportedModels,
			Tags:    map[string]string{"service": "bedrock", "action": "ListImportedModels"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("model_arn"),
			Hydrate:    getBedrockImportedModel,
			Tags:       map[string]string{"service": "bedrock", "action": "GetImportedModel"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_BEDROCK_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "model_name",
				Description: "The name of the imported model.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "model_arn",
				Description: "The Amazon Resource Name (ARN) of the imported model.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time when the model was imported.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "instruct_supported",
				Description: "Specifies if the imported model supports converse.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "model_architecture",
				Description: "The architecture of the imported model.",
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

//// LIST FUNCTION

func listBedrockImportedModels(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := BedrockClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_imported_model.listBedrockImportedModels", "client_error", err)
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

	input := &bedrock.ListImportedModelsInput{
		MaxResults: &maxLimit,
	}

	paginator := bedrock.NewListImportedModelsPaginator(svc, input, func(o *bedrock.ListImportedModelsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_bedrock_imported_model.listBedrockImportedModels", "api_error", err)
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

func getBedrockImportedModel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	modelArn := d.EqualsQualString("model_arn")

	// Empty check
	if modelArn == "" {
		return nil, nil
	}

	// Create client
	svc, err := BedrockClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_imported_model.getBedrockImportedModel", "client_error", err)
		return nil, err
	}

	// Build the params
	params := &bedrock.GetImportedModelInput{
		ModelIdentifier: &modelArn,
	}

	// Get call
	data, err := svc.GetImportedModel(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_imported_model.getBedrockImportedModel", "api_error", err)
		return nil, err
	}

	return data, nil
}
