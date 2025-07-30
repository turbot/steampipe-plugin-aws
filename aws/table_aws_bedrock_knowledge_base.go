package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagent"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagent/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsBedrockKnowledgeBase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_bedrock_knowledge_base",
		Description: "AWS Bedrock Knowledge Base",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "knowledge_base_id", Require: plugin.Required},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getBedrockKnowledgeBase,
			Tags:    map[string]string{"service": "bedrock", "action": "GetKnowledgeBase"},
		},
		List: &plugin.ListConfig{
			Hydrate: listBedrockKnowledgeBases,
			Tags:    map[string]string{"service": "bedrock", "action": "ListKnowledgeBases"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getBedrockKnowledgeBase,
				Tags: map[string]string{"service": "bedrock", "action": "GetKnowledgeBase"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_BEDROCK_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			// Key columns
			{
				Name:        "knowledge_base_id",
				Description: "The unique identifier of the knowledge base.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the knowledge base.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the knowledge base.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the knowledge base.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The time at which the knowledge base was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getBedrockKnowledgeBase,
			},
			{
				Name:        "updated_at",
				Description: "The time at which the knowledge base was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the knowledge base.",
				Transform:   transform.FromField("KnowledgeBaseArn"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockKnowledgeBase,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role with permissions to invoke API operations on the knowledge base.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockKnowledgeBase,
			},
			{
				Name:        "storage_configuration",
				Description: "The storage configuration for the knowledge base.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockKnowledgeBase,
			},
			{
				Name:        "knowledge_base_configuration",
				Description: "Contains details about the embeddings configuration of the knowledge base.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockKnowledgeBase,
			},
			{
				Name:        "failure_reasons",
				Description: "Contains reasons that the knowledge base-related API that you invoked failed.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockKnowledgeBase,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("KnowledgeBaseArn").Transform(transform.EnsureStringArray),
				Hydrate:     getBedrockKnowledgeBase,
			},
		}),
	}
}

//// LIST FUNCTION

func listBedrockKnowledgeBases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BedrockAgentClient(ctx, d)

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_knowledge_base.listBedrockKnowledgeBases", "connection_error", err)
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

	input := &bedrockagent.ListKnowledgeBasesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := bedrockagent.NewListKnowledgeBasesPaginator(svc, input, func(o *bedrockagent.ListKnowledgeBasesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// Handle the unsupported region error since the resource is not available in all the regions
			if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("ValidationException: Unknown operation")) {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_bedrock_knowledge_base.listBedrockKnowledgeBases", "api_error", err)
			return nil, err
		}

		for _, kb := range output.KnowledgeBaseSummaries {
			d.StreamListItem(ctx, kb)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getBedrockKnowledgeBase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var knowledgeBaseId string

	if h.Item != nil {
		// Retrieve knowledgeBaseId from the List call
		knowledgeBaseSummary := h.Item.(types.KnowledgeBaseSummary)
		knowledgeBaseId = *knowledgeBaseSummary.KnowledgeBaseId
	} else {
		knowledgeBaseId = d.EqualsQualString("knowledge_base_id")
	}

	// Empty check
	if knowledgeBaseId == "" {
		return nil, nil
	}

	// Create service
	svc, err := BedrockAgentClient(ctx, d)

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_knowledge_base.getBedrockKnowledgeBase", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &bedrockagent.GetKnowledgeBaseInput{
		KnowledgeBaseId: aws.String(knowledgeBaseId),
	}

	// Get call
	data, err := svc.GetKnowledgeBase(ctx, params)
	if err != nil {
		// Handle the unsupported region error since the resource is not available in all the regions
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("ValidationException: Unknown operation")) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_bedrock_knowledge_base.getBedrockKnowledgeBase", "api_error", err)
		return nil, err
	}

	return data.KnowledgeBase, nil
}
