package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsResourcePolicyAnalysis(_ context.Context) *plugin.Table {
	// permissionsData = getParliamentIamPermissions()

	return &plugin.Table{
		Name:        "aws_resource_policy_analysis",
		Description: "AWS Resource Policy Analysis",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "policy", CacheMatch: "exact"},
				{Name: "account_id", CacheMatch: "exact"},
			},
			Hydrate: listResourcePolicyAnalysis,
		},
		Columns: []*plugin.Column{
			// "Key" Columns
			{
				Name:        "access_level",
				Type:        proto.ColumnType_STRING,
				Description: "Access level of the resource based of policy. Valid values are 'public', 'private' and 'shared'.",
			},
			{
				Name:        "is_public",
				Type:        proto.ColumnType_BOOL,
				Description: "The policy status for an Amazon resource, indicating whether the resource is public.",
			},
			{
				Name:        "policy",
				Description: "The input policy to be analysed.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromQual("policy"),
			},
			{
				Name:        "public_statement_ids",
				Type:        proto.ColumnType_JSON,
				Description: "The sid of the statements that makes the resource public. If sid is not available in statement it will be of the format Statement[index].",
			},
			{
				Name:        "public_access_levels",
				Type:        proto.ColumnType_JSON,
				Description: "Public access levels (based off parliament's levels).",
			},
			{
				Name:        "allowed_organization_ids",
				Type:        proto.ColumnType_JSON,
				Description: "A list of organisations resource is accessible to.",
			},
			{
				Name:        "allowed_principals",
				Type:        proto.ColumnType_JSON,
				Description: "A list of principals resource is accessible to.",
			},
			{
				Name:        "allowed_principal_account_ids",
				Type:        proto.ColumnType_JSON,
				Description: "A list of account ids resource is accessible to.",
			},
			{
				Name:        "allowed_principal_federated_identities",
				Type:        proto.ColumnType_JSON,
				Description: "A list of federated identities resource is accessible to.",
			},
			{
				Name:        "allowed_principal_services",
				Type:        proto.ColumnType_JSON,
				Description: "A list of services resource is accessible to.",
			},
			{
				Name:        "account_id",
				Type:        proto.ColumnType_STRING,
				Description: "The id of AWS account to which resource belongs.",
				Transform:   transform.FromQual("account_id"),
			},
		},
	}
}

//// LIST FUNCTION

func listResourcePolicyAnalysis(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	policyVal := d.KeyColumnQuals["policy"].GetJsonbValue()
	if policyVal == "" {
		return nil, nil
	}

	accountID := d.KeyColumnQuals["account_id"].GetStringValue()
	if accountID == "" {
		return nil, nil
	}

	// convert input policy to canonical form for easier handling
	policy, err := canonicalPolicy(policyVal)
	if err != nil {
		plugin.Logger(ctx).Error("aws_resource_policy_analysis.listResourcePolicyAnalysis", "policy_canonicalization_error", err)
		return nil, err
	}

	policyObject, ok := policy.(Policy)
	if !ok {
		plugin.Logger(ctx).Error("aws_resource_policy_analysis.listResourcePolicyAnalysis", "policy_coercion_error", err)
		return nil, err
	}

	evaluation, err := policyObject.EvaluatePolicy(accountID)
	if err != nil {
		plugin.Logger(ctx).Error("aws_resource_policy_analysis.listResourcePolicyAnalysis", "policy_evaluation_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, evaluation)
	// Context can be cancelled due to manual cancellation or the limit has been hit
	if d.QueryStatus.RowsRemaining(ctx) == 0 {
		return nil, nil
	}

	return nil, nil
}
