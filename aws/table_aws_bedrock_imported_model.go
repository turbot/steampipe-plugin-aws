package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrock/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsBedrockImportedModel(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_bedrock_imported_model",
		Description: "AWS Bedrock Imported Model",
		List: &plugin.ListConfig{
			Hydrate: listBedrockImportedModels,
			Tags:    map[string]string{"service": "bedrock", "action": "ListImportedModels"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "creation_time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "arn", Require: plugin.AnyOf},
				{Name: "model_name", Require: plugin.AnyOf},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getBedrockImportedModel,
			Tags:    map[string]string{"service": "bedrock", "action": "GetImportedModel"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_BEDROCK_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "model_name",
				Description: "The name of the imported model.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the imported model.",
				Transform:   transform.FromField("ModelArn"),
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
			{
				Name:        "job_arn",
				Description: "Job Amazon Resource Name (ARN) associated with the imported model.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockImportedModel,
			},
			{
				Name:        "job_name",
				Description: "Job name associated with the imported model.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockImportedModel,
			},
			{
				Name:        "model_data_source",
				Description: "The data source for this imported model.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockImportedModel,
			},
			{
				Name:        "model_kms_key_arn",
				Description: "The imported model is encrypted at rest using this key.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockImportedModel,
			},
			{
				Name:        "custom_model_units",
				Description: "Information about the hardware utilization for a single copy of the model.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockImportedModel,
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

func listBedrockImportedModels(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := BedrockClient(ctx, d)
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_imported_model.listBedrockImportedModels", "client_error", err)
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

	input := &bedrock.ListImportedModelsInput{
		MaxResults: &maxLimit,
	}

	// Apply creation_time qual if provided
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
			// Handle the unsupported region error since the resource is not available in all the regions: ValidationException: Unknown operation
			if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("ValidationException: Unknown operation")) {
				return nil, nil
			}
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
	var modelIdentifier string

	if h.Item != nil {
		modelIdentifier = *h.Item.(types.ImportedModelSummary).ModelArn
	} else {
		if d.EqualsQualString("arn") != "" {
			modelIdentifier = d.EqualsQualString("arn")
		} else {
			modelIdentifier = d.EqualsQualString("model_name")
		}
	}

	// Empty check
	if modelIdentifier == "" {
		return nil, nil
	}

	// Create client
	svc, err := BedrockClient(ctx, d)

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_imported_model.getBedrockImportedModel", "client_error", err)
		return nil, err
	}

	// Build the params
	params := &bedrock.GetImportedModelInput{
		ModelIdentifier: &modelIdentifier,
	}

	// Get call
	data, err := svc.GetImportedModel(ctx, params)
	if err != nil {
		// Handle the unsupported region error since the resource is not available in all the regions: ValidationException: Unknown operation
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("ValidationException: Unknown operation")) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_bedrock_imported_model.getBedrockImportedModel", "api_error", err)
		return nil, err
	}

	return data, nil
}
