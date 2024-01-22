package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/fms"
	"github.com/aws/aws-sdk-go-v2/service/fms/types"

	fmsv1 "github.com/aws/aws-sdk-go/service/fms"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsFMSPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_fms_policy",
		Description: "AWS FMS Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("policy_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getFmsPolicy,
			Tags:    map[string]string{"service": "fms", "action": "GetPolicy"},
		},
		List: &plugin.ListConfig{
			Hydrate: listFmsPolicies,
			Tags:    map[string]string{"service": "fms", "action": "ListPolicies"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getFmsPolicy,
				Tags: map[string]string{"service": "fms", "action": "GetPolicy"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(fmsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "policy_name",
				Description: "The name of the specified policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy_id",
				Description: "The ID of the specified policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicyArn"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the specified policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicyArn"),
			},
			{
				Name:        "policy_description",
				Description: "The definition of the Network Firewall firewall policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFmsPolicy,
			},
			{
				Name:        "policy_status",
				Description: "Indicates whether the policy is in or out of an admin's policy or Region scope. The possible values ACTIVE, OUT_OF_ADMIN_SCOPE.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The type of resource protected by or in scope of the policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "exclude_resource_tags",
				Description: "If set to True , resources with the tags that are specified in the ResourceTag array are not in scope of the policy.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getFmsPolicy,
			},
			{
				Name:        "remediation_enabled",
				Description: "Indicates if the policy should be automatically applied to new resources.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "security_service_policy_data",
				Description: "Details about the security service that is being used to protect the resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "delete_unused_fm_managed_resources",
				Description: "The AWS account that created the file system.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DeleteUnusedFMManagedResources"),
			},
			{
				Name:        "policy_update_token",
				Description: "A unique identifier for each update to the policy. When issuing a PutPolicy request, the PolicyUpdateToken in the request must match the PolicyUpdateToken of the current policy version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFmsPolicy,
			},
			{
				Name:        "exclude_map",
				Description: "Specifies the Amazon Web Services account IDs and Organizations organizational units (OUs) to exclude from the policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFmsPolicy,
			},
			{
				Name:        "include_map",
				Description: "Specifies the Amazon Web Services account IDs and Organizations organizational units (OUs) to include in the policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFmsPolicy,
			},
			{
				Name:        "resource_set_ids",
				Description: "The unique identifiers of the resource sets used by the policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFmsPolicy,
			},
			{
				Name:        "resource_tags",
				Description: "An array of ResourceTag objects.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFmsPolicy,
			},
			{
				Name:        "resource_type_list",
				Description: "An array of ResourceType objects. Use this only to specify multiple resource types.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFmsPolicy,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicyName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type PolicyInfo struct {
	ExcludeResourceTags            bool
	PolicyArn                      *string
	PolicyName                     *string
	RemediationEnabled             bool
	ResourceType                   *string
	SecurityServicePolicyData      *types.SecurityServicePolicyData
	DeleteUnusedFMManagedResources bool
	ExcludeMap                     map[string][]string
	IncludeMap                     map[string][]string
	PolicyDescription              *string
	PolicyId                       *string
	PolicyStatus                   types.CustomerPolicyStatus
	PolicyUpdateToken              *string
	ResourceSetIds                 []string
	ResourceTags                   []types.ResourceTag
	ResourceTypeList               []string
}

//// LIST FUNCTION

func listFmsPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := FMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fms_policy.listFmsPolicies", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxItems := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	input := fms.ListPoliciesInput{
		MaxResults: aws.Int32(maxItems),
	}

	paginator := fms.NewListPoliciesPaginator(svc, &input, func(o *fms.ListPoliciesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_fms_policy.listFmsPolicies", "api_error", err)
			return nil, err
		}

		for _, policy := range output.PolicyList {
			d.StreamListItem(ctx, PolicyInfo{
				PolicyArn:                      policy.PolicyArn,
				PolicyName:                     policy.PolicyName,
				PolicyId:                       policy.PolicyId,
				PolicyStatus:                   policy.PolicyStatus,
				DeleteUnusedFMManagedResources: policy.DeleteUnusedFMManagedResources,
				RemediationEnabled:             policy.RemediationEnabled,
				ResourceType:                   policy.ResourceType,
				SecurityServicePolicyData: &types.SecurityServicePolicyData{
					Type: policy.SecurityServiceType,
				},
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getFmsPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	policyId := ""

	if h.Item != nil {
		data := h.Item.(types.PolicySummary)
		policyId = *data.PolicyId
	} else {
		policyId = d.EqualsQualString("policy_id")
	}

	if policyId == "" {
		return nil, nil
	}
	// Create service
	svc, err := FMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fms_policy.getFmsPolicy", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &fms.GetPolicyInput{
		PolicyId: &policyId,
	}

	op, err := svc.GetPolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fms_policy.getFmsPolicy", "api_error", err)
		return nil, err
	}

	if op != nil && op.Policy != nil {
		return PolicyInfo{
			PolicyArn:                      op.PolicyArn,
			PolicyName:                     op.Policy.PolicyName,
			PolicyId:                       op.Policy.PolicyId,
			PolicyStatus:                   op.Policy.PolicyStatus,
			DeleteUnusedFMManagedResources: op.Policy.DeleteUnusedFMManagedResources,
			RemediationEnabled:             op.Policy.RemediationEnabled,
			ResourceType:                   op.Policy.ResourceType,
			SecurityServicePolicyData:      op.Policy.SecurityServicePolicyData,
			ExcludeResourceTags:            op.Policy.ExcludeResourceTags,
			ExcludeMap:                     op.Policy.ExcludeMap,
			IncludeMap:                     op.Policy.IncludeMap,
			PolicyDescription:              op.Policy.PolicyDescription,
			PolicyUpdateToken:              op.Policy.PolicyUpdateToken,
			ResourceSetIds:                 op.Policy.ResourceSetIds,
			ResourceTags:                   op.Policy.ResourceTags,
			ResourceTypeList:               op.Policy.ResourceTypeList,
		}, nil
	}

	return op, nil
}
