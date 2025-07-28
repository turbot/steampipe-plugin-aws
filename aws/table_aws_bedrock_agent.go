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

func tableAwsBedrockAgent(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_bedrock_agent",
		Description: "AWS Bedrock Agent",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "agent_id", Require: plugin.Required},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getBedrockAgent,
			Tags:    map[string]string{"service": "bedrock", "action": "GetAgent"},
		},
		List: &plugin.ListConfig{
			Hydrate: listBedrockAgents,
			Tags:    map[string]string{"service": "bedrock", "action": "ListAgents"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getBedrockAgent,
				Tags: map[string]string{"service": "bedrock", "action": "GetAgent"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_BEDROCK_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			// Columns from ListAgents (AgentSummary)
			{
				Name:        "agent_id",
				Description: "The unique identifier of the agent.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "agent_name",
				Description: "The name of the agent.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "agent_status",
				Description: "The status of the agent and whether it is ready for use.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AgentStatus"),
			},
			{
				Name:        "updated_at",
				Description: "The time at which the agent was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description of the agent.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "latest_agent_version",
				Description: "The latest version of the agent.",
				Type:        proto.ColumnType_STRING,
			},

			// Columns from GetAgent (Agent)
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the agent.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockAgent,
			},
			{
				Name:        "agent_resource_role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role with permissions to invoke API operations on the agent.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockAgent,
			},
			{
				Name:        "agent_version",
				Description: "The version of the agent.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockAgent,
			},
			{
				Name:        "created_at",
				Description: "The time at which the agent was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getBedrockAgent,
			},
			{
				Name:        "prepared_at",
				Description: "The time at which the agent was last prepared.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getBedrockAgent,
			},
			{
				Name:        "idle_session_ttl_in_seconds",
				Description: "The number of seconds for which Amazon Bedrock keeps information about a user's conversation with the agent.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("IdleSessionTTLInSeconds"),
				Hydrate:     getBedrockAgent,
			},
			{
				Name:        "foundation_model",
				Description: "The foundation model used for orchestration by the agent.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockAgent,
			},
			{
				Name:        "instruction",
				Description: "Instructions that tell the agent what it should do and how it should interact with users.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockAgent,
			},
			{
				Name:        "orchestration_type",
				Description: "Specifies the orchestration strategy for the agent.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OrchestrationType"),
				Hydrate:     getBedrockAgent,
			},
			{
				Name:        "customer_encryption_key_arn",
				Description: "The Amazon Resource Name (ARN) of the KMS key that encrypts the agent.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockAgent,
			},
			{
				Name:        "failure_reasons",
				Description: "Contains reasons that the agent-related API that you invoked failed.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgent,
			},
			{
				Name:        "recommended_actions",
				Description: "Contains recommended actions to take for the agent-related API that you invoked to succeed.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgent,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AgentName"),
			},
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AgentArn").Transform(transform.EnsureStringArray),
				Hydrate:     getBedrockAgent,
			},
		}),
	}
}

//// LIST FUNCTION

func listBedrockAgents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BedrockAgentClient(ctx, d)

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_agent.listBedrockAgents", "connection_error", err)
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

	input := &bedrockagent.ListAgentsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := bedrockagent.NewListAgentsPaginator(svc, input, func(o *bedrockagent.ListAgentsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// Handle the unsupported region error since the resource is not available in all the regions: ValidationException: Unknown operation
			if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("ValidationException: Unknown operation")) {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_bedrock_agent.listBedrockAgents", "api_error", err)
			return nil, err
		}

		for _, agent := range output.AgentSummaries {
			d.StreamListItem(ctx, agent)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getBedrockAgent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var agentId string

	if h.Item != nil {
		// Retrieve agentId from the List call
		agentSummary := h.Item.(types.AgentSummary)
		agentId = *agentSummary.AgentId
	} else {
		agentId = d.EqualsQualString("agent_id")
	}

	// Empty check
	if agentId == "" {
		return nil, nil
	}

	// Create service
	svc, err := BedrockAgentClient(ctx, d)

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_agent.getBedrockAgent", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &bedrockagent.GetAgentInput{
		AgentId: aws.String(agentId),
	}

	// Get call
	data, err := svc.GetAgent(ctx, params)
	if err != nil {
		// Handle the unsupported region error since the resource is not available in all the regions: ValidationException: Unknown operation
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("ValidationException: Unknown operation")) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_bedrock_agent.getBedrockAgent", "api_error", err)
		return nil, err
	}

	return data.Agent, nil
}
