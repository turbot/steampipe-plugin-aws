package aws

import (
	"context"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentcorecontrol"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentcorecontrol/types"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableAwsBedrockAgentCoreAgentRuntime(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_bedrock_agentcore_agent_runtime",
		Description: "AWS Bedrock AgentCore Agent Runtime",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("agent_runtime_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getBedrockAgentCoreAgentRuntime,
			Tags:    map[string]string{"service": "bedrock-agentcore", "action": "GetAgentRuntime"},
		},
		List: &plugin.ListConfig{
			Hydrate: listBedrockAgentCoreAgentRuntimes,
			Tags:    map[string]string{"service": "bedrock-agentcore", "action": "ListAgentRuntimes"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getBedrockAgentCoreAgentRuntime,
				Tags: map[string]string{"service": "bedrock-agentcore", "action": "GetAgentRuntime"},
			},
			{
				Func: getBedrockAgentCoreAgentRuntimeTags,
				// The runtime can be deleted between the List call and this hydrate
				IgnoreConfig: &plugin.IgnoreConfig{
					ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "NotFoundException"}),
				},
				Tags: map[string]string{"service": "bedrock-agentcore", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_BEDROCK_AGENTCORE_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			// Columns from ListAgentRuntimes (AgentRuntime)
			{
				Name:        "agent_runtime_id",
				Description: "The unique identifier of the agent runtime.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "agent_runtime_name",
				Description: "The name of the agent runtime.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the agent runtime.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AgentRuntimeArn"),
			},
			{
				Name:        "agent_runtime_version",
				Description: "The version of the agent runtime.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current status of the agent runtime. Possible values are CREATING, CREATE_FAILED, UPDATING, UPDATE_FAILED, READY, DELETING, and DELETE_FAILED.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the agent runtime.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_at",
				Description: "The timestamp when the agent runtime was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Columns from GetAgentRuntime
			{
				Name:        "created_at",
				Description: "The timestamp when the agent runtime was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role that provides permissions for the agent runtime.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
			},
			{
				Name:        "platform_version",
				Description: "The platform version of the agent runtime.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
			},
			{
				Name:        "failure_reason",
				Description: "The reason for failure if the agent runtime is in a failed state.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
			},
			{
				Name:        "network_configuration",
				Description: "The network configuration of the agent runtime, including the network mode and VPC settings.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
			},
			{
				Name:        "agent_runtime_artifact",
				Description: "The artifact configuration of the agent runtime. Contains either a ContainerConfiguration key (container image URI) or a CodeConfiguration key (S3 code location, runtime, and entry point).",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
				Transform:   transform.FromField("AgentRuntimeArtifact").Transform(flattenSmithyUnionsTransform),
			},
			{
				Name:        "authorizer_configuration",
				Description: "The inbound authorizer configuration of the agent runtime, keyed by authorizer type, for example CustomJWTAuthorizer. Null when the runtime uses AWS IAM authorization.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
				Transform:   transform.FromField("AuthorizerConfiguration").Transform(flattenSmithyUnionsTransform),
			},
			{
				Name:        "capacity_provider_configuration",
				Description: "The capacity provider configuration of the agent runtime.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
			},
			{
				Name:        "environment_variables",
				Description: "Environment variables set in the agent runtime environment.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
			},
			{
				Name:        "filesystem_configurations",
				Description: "The file systems mounted into the agent runtime. Each entry is keyed by type: SessionStorage, S3FilesAccessPoint, EfsAccessPoint, or CapacityProviderVolume.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
				Transform:   transform.FromField("FilesystemConfigurations").Transform(flattenSmithyUnionsTransform),
			},
			{
				Name:        "lifecycle_configuration",
				Description: "The lifecycle configuration of the agent runtime, such as idle session timeout.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
			},
			{
				Name:        "metadata_configuration",
				Description: "The metadata configuration of the agent runtime.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
			},
			{
				Name:        "protocol_configuration",
				Description: "The protocol configuration of the agent runtime, such as HTTP, MCP, or A2A.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
			},
			{
				Name:        "request_header_configuration",
				Description: "The request header configuration of the agent runtime, keyed by type, for example RequestHeaderAllowlist with the list of headers passed through to the agent.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
				Transform:   transform.FromField("RequestHeaderConfiguration").Transform(flattenSmithyUnionsTransform),
			},
			{
				Name:        "workload_identity_details",
				Description: "The workload identity details for the agent runtime, including the ARN of the workload identity it uses to obtain access tokens.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntime,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AgentRuntimeName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBedrockAgentCoreAgentRuntimeTags,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AgentRuntimeArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listBedrockAgentCoreAgentRuntimes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BedrockAgentCoreControlClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_agentcore_agent_runtime.listBedrockAgentCoreAgentRuntimes", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &bedrockagentcorecontrol.ListAgentRuntimesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := bedrockagentcorecontrol.NewListAgentRuntimesPaginator(svc, input, func(o *bedrockagentcorecontrol.ListAgentRuntimesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_bedrock_agentcore_agent_runtime.listBedrockAgentCoreAgentRuntimes", "api_error", err)
			return nil, err
		}

		for _, runtime := range output.AgentRuntimes {
			d.StreamListItem(ctx, runtime)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getBedrockAgentCoreAgentRuntime(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var runtimeId string

	if h.Item != nil {
		// Retrieve the runtime ID from the List call
		runtime := h.Item.(types.AgentRuntime)
		runtimeId = *runtime.AgentRuntimeId
	} else {
		runtimeId = d.EqualsQualString("agent_runtime_id")
	}

	// The API returns AccessDeniedException rather than ResourceNotFoundException
	// for IDs that do not match its pattern, so reject those before calling it
	if !isValidBedrockAgentCoreAgentRuntimeId(runtimeId) {
		return nil, nil
	}

	// Create service
	svc, err := BedrockAgentCoreControlClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_agentcore_agent_runtime.getBedrockAgentCoreAgentRuntime", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &bedrockagentcorecontrol.GetAgentRuntimeInput{
		AgentRuntimeId: aws.String(runtimeId),
	}

	// Get call
	data, err := svc.GetAgentRuntime(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_agentcore_agent_runtime.getBedrockAgentCoreAgentRuntime", "api_error", err)
		return nil, err
	}

	return data, nil
}

func getBedrockAgentCoreAgentRuntimeTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := bedrockAgentCoreAgentRuntimeArn(h.Item)
	if arn == "" {
		return nil, nil
	}

	// Create service
	svc, err := BedrockAgentCoreControlClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_agentcore_agent_runtime.getBedrockAgentCoreAgentRuntimeTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &bedrockagentcorecontrol.ListTagsForResourceInput{
		ResourceArn: aws.String(arn),
	}

	data, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_agentcore_agent_runtime.getBedrockAgentCoreAgentRuntimeTags", "api_error", err)
		return nil, err
	}

	return data, nil
}

// Pattern for AgentRuntimeId from the Bedrock AgentCore Control API model.
var bedrockAgentCoreAgentRuntimeIdRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{0,99}-[a-zA-Z0-9]{10}$`)

// isValidBedrockAgentCoreAgentRuntimeId reports whether id matches the API's AgentRuntimeId pattern.
func isValidBedrockAgentCoreAgentRuntimeId(id string) bool {
	return bedrockAgentCoreAgentRuntimeIdRe.MatchString(id)
}

// bedrockAgentCoreAgentRuntimeArn returns the runtime ARN from either a list item or a get result.
func bedrockAgentCoreAgentRuntimeArn(item interface{}) string {
	switch v := item.(type) {
	case types.AgentRuntime:
		if v.AgentRuntimeArn != nil {
			return *v.AgentRuntimeArn
		}
	case *bedrockagentcorecontrol.GetAgentRuntimeOutput:
		if v.AgentRuntimeArn != nil {
			return *v.AgentRuntimeArn
		}
	}
	return ""
}
