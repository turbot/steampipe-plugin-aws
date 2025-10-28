package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrock/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsBedrockGuardrail(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_bedrock_guardrail",
		Description: "Amazon Bedrock Guardrail.",
		List: &plugin.ListConfig{
			Hydrate: listBedrockGuardrails,
			Tags:    map[string]string{"service": "bedrock", "action": "ListGuardrails"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("guardrail_id"),
			Hydrate:    getBedrockGuardrail,
			Tags:       map[string]string{"service": "bedrock", "action": "GetGuardrail"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getBedrockGuardrail,
				Tags: map[string]string{"service": "bedrock", "action": "GetGuardrail"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_BEDROCK_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			// String columns (available in List)
			{Name: "arn", Type: proto.ColumnType_STRING, Description: "ARN of the guardrail.", Transform: transform.From(guardrailArn)},
			{Name: "guardrail_id", Type: proto.ColumnType_STRING, Description: "ID of the guardrail.", Transform: transform.From(guardrailID)},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the guardrail."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the guardrail."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the guardrail."},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "Version (DRAFT or a number)."},

			// Timestamp columns (available in List)
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time at which the guardrail was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time at which the guardrail was last updated."},

			// Additional fields from Get call
			{Name: "blocked_input_messaging", Type: proto.ColumnType_STRING, Description: "The message that the guardrail returns when it blocks a prompt.", Hydrate: getBedrockGuardrail},
			{Name: "blocked_outputs_messaging", Type: proto.ColumnType_STRING, Description: "The message that the guardrail returns when it blocks a model response.", Hydrate: getBedrockGuardrail},
			{Name: "kms_key_arn", Type: proto.ColumnType_STRING, Description: "The ARN of the KMS key that encrypts the guardrail.", Hydrate: getBedrockGuardrail},

			// JSON columns (from Get call)
			{Name: "content_policy", Type: proto.ColumnType_JSON, Description: "The content policy configuration for the guardrail.", Hydrate: getBedrockGuardrail},
			{Name: "contextual_grounding_policy", Type: proto.ColumnType_JSON, Description: "The contextual grounding policy settings.", Hydrate: getBedrockGuardrail},
			{Name: "sensitive_information_policy", Type: proto.ColumnType_JSON, Description: "The policy for handling sensitive information.", Hydrate: getBedrockGuardrail},
			{Name: "topic_policy", Type: proto.ColumnType_JSON, Description: "The topic-based policy configuration.", Hydrate: getBedrockGuardrail},
			{Name: "word_policy", Type: proto.ColumnType_JSON, Description: "The word-based policy settings.", Hydrate: getBedrockGuardrail},
			{Name: "cross_region_details", Type: proto.ColumnType_JSON, Description: "Details about system-defined guardrail profile across regions.", Hydrate: getBedrockGuardrail},
			{Name: "failure_recommendations", Type: proto.ColumnType_JSON, Description: "List of recommendations if guardrail creation/update failed.", Hydrate: getBedrockGuardrail},
			{Name: "status_reasons", Type: proto.ColumnType_JSON, Description: "Reasons for failure status if applicable.", Hydrate: getBedrockGuardrail},

			// Steampipe standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: resourceInterfaceDescription("title"), Transform: transform.FromField("Name")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: resourceInterfaceDescription("akas"), Transform: transform.From(guardrailArn).Transform(transform.EnsureStringArray)},
		}),
	}
}

func listBedrockGuardrails(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BedrockClient(ctx, d)

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_guardrail.listBedrockGuardrails", "connection_error", err)
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

	input := &bedrock.ListGuardrailsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := bedrock.NewListGuardrailsPaginator(svc, input, func(o *bedrock.ListGuardrailsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_bedrock_guardrail.listBedrockGuardrails", "api_error", err)
			return nil, err
		}

		for _, guardrail := range output.Guardrails {
			d.StreamListItem(ctx, guardrail)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

func getBedrockGuardrail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create service
	svc, err := BedrockClient(ctx, d)

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_guardrail.getBedrockGuardrail", "connection_error", err)
		return nil, err
	}

	var guardrailId string

	if h.Item != nil {
		// Retrieve guardrailId from the List call
		guardrail := h.Item.(types.GuardrailSummary)
		guardrailId = *guardrail.Id
	} else {
		guardrailId = d.EqualsQuals["guardrail_id"].GetStringValue()
	}

	if strings.TrimSpace(guardrailId) == "" {
		return nil, nil
	}

	// Build the params
	params := &bedrock.GetGuardrailInput{
		GuardrailIdentifier: aws.String(guardrailId),
	}

	data, err := svc.GetGuardrail(ctx, params)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("ValidationException")) {
			plugin.Logger(ctx).Debug("aws_bedrock_guardrail.getBedrockGuardrail", "validation_exception", err)
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_bedrock_guardrail.getBedrockGuardrail", "api_error", err)
		return nil, err
	}

	return data, nil
}

// Transform functions to handle field name differences between GuardrailSummary (List) and GetGuardrailOutput (Get)
func guardrailID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	switch item := d.HydrateItem.(type) {
	case types.GuardrailSummary:
		return item.Id, nil
	case *bedrock.GetGuardrailOutput:
		return item.GuardrailId, nil
	}
	return nil, nil
}

func guardrailArn(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	switch item := d.HydrateItem.(type) {
	case types.GuardrailSummary:
		return item.Arn, nil
	case *bedrock.GetGuardrailOutput:
		return item.GuardrailArn, nil
	}
	return nil, nil
}
