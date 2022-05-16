package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dlm"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAwsDlmLifecyclePolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dlm_lifecycle_policy",
		Description: "AWS DLM Lifecycle Policy",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("policy_id"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate:           getDlmLifecyclePolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listDlmLifecyclePolicies,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "policy_id",
				Description: "The identifier of the lifecycle policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicyArn"),
				Hydrate:     getDlmLifecyclePolicy,
			},
			{
				Name:        "description",
				Description: "The description of the lifecycle policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "date_created",
				Description: "The local date and time when the lifecycle policy was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getDlmLifecyclePolicy,
			},
			{
				Name:        "date_modified",
				Description: "The local date and time when the lifecycle policy was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getDlmLifecyclePolicy,
			},
			{
				Name:        "execution_role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role used to run the operations specified by the lifecycle policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDlmLifecyclePolicy,
			},
			{
				Name:        "policy_type",
				Description: "The type of policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The activation state of the lifecycle policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_message",
				Description: "The description of the status.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDlmLifecyclePolicy,
			},
			{
				Name:        "policy_details",
				Description: "The configuration of the lifecycle policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDlmLifecyclePolicy,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicyId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyArn").Transform(transform.EnsureStringArray),
				Hydrate:     getDlmLifecyclePolicy,
			},
		}),
	}
}

//// LIST FUNCTION

func listDlmLifecyclePolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDlmLifecyclePolicies")

	// Create Session
	svc, err := DLMService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &dlm.GetLifecyclePoliciesInput{}

	policies, err := svc.GetLifecyclePolicies(input)
	if err != nil {
		plugin.Logger(ctx).Error("listDlmLifecyclePolicies", "list", err)
		return nil, err
	}
	if policies.Policies == nil {
		return nil, nil
	}

	for _, policy := range policies.Policies {
		d.StreamListItem(ctx, policy)

		// Check if context has been cancelled or if the limit has been reached (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDlmLifecyclePolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDlmLifecyclePolicy")

	var id string
	if h.Item != nil {
		id = *policyId(h.Item)
	} else {
		id = d.KeyColumnQuals["policy_id"].GetStringValue()
	}

	// Empty check
	if len(id) < 1 {
		return nil, nil
	}

	// Create service
	svc, err := DLMService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getDlmLifecyclePolicy", "connection", err)
		return nil, err
	}

	params := &dlm.GetLifecyclePolicyInput{
		PolicyId: aws.String(id),
	}

	op, err := svc.GetLifecyclePolicy(params)
	if err != nil {
		plugin.Logger(ctx).Error("getDlmLifecyclePolicy", "get", err)
		return nil, err
	}
	return op.Policy, nil
}

func policyId(item interface{}) *string {
	switch item := item.(type) {
	case *dlm.LifecyclePolicy:
		return item.PolicyId
	case *dlm.LifecyclePolicySummary:
		return item.PolicyId
	}
	return nil
}
