package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsResourcePolicyAnalysis(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_resource_policy_analysis",
		Description: "AWS Resource Policy Analysis",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				//{Name: "policy", CacheMatch: "exact", Require: plugin.Optional},
				{Name: "policy", CacheMatch: "exact"},
				{Name: "account_id", CacheMatch: "exact"},
			},
			Hydrate: listResourcePolicyAnalysis,
		},
		Columns: []*plugin.Column{
			{
				Name:        "access_level",
				Type:        proto.ColumnType_STRING,
				Description: "Overall access level of the resources based off the policy. Valid values are 'public', 'shared' and 'private'.",
			},
			{
				Name:        "is_public",
				Type:        proto.ColumnType_BOOL,
				Description: "A convenient flag used to check if the access level of the policy is 'public'.",
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
				Description: "List of statement Sids that grant public access to resources. If Sid is missing then 'Statement[index]' will be used.",
			},
			{
				Name:        "shared_statement_ids",
				Type:        proto.ColumnType_JSON,
				Description: "List of statement Sids that grant shared access to resource(s). If Sid is missing then 'Statement[index]' will be used.",
			},
			{
				Name:        "public_access_levels",
				Type:        proto.ColumnType_JSON,
				Description: "Public access levels such as 'Read', 'Write', 'Tagging', etc. to describe the actions allowed by the policy",
			},
			{
				Name:        "shared_access_levels",
				Type:        proto.ColumnType_JSON,
				Description: "Shared access levels such as 'Read', 'Write', 'Tagging', etc. to describe the actions allowed by the policy",
			},
			{
				Name:        "private_access_levels",
				Type:        proto.ColumnType_JSON,
				Description: "Private access levels such as 'Read', 'Write', 'Tagging', etc. to describe the actions allowed by the policy",
			},
			{
				Name:        "allowed_organization_ids",
				Type:        proto.ColumnType_JSON,
				Description: "List of organisations allowed to access the resources.",
			},
			{
				Name:        "allowed_principals",
				Type:        proto.ColumnType_JSON,
				Description: "List of principals allowed to access the resources.",
			},
			{
				Name:        "allowed_principal_account_ids",
				Type:        proto.ColumnType_JSON,
				Description: "List of account ids allowed to access the resources.",
			},
			{
				Name:        "allowed_principal_federated_identities",
				Type:        proto.ColumnType_JSON,
				Description: "List federated identities allowed to access resources.",
			},
			{
				Name:        "allowed_principal_services",
				Type:        proto.ColumnType_JSON,
				Description: "List of AWS services allowed to access resources.",
			},
			{
				Name:        "account_id",
				Type:        proto.ColumnType_STRING,
				Description: "The id of AWS account where the policy is deployed.",
				Transform:   transform.FromQual("account_id"),
			},
		},
	}
}

//// LIST FUNCTION

func listResourcePolicyAnalysis(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountID := d.KeyColumnQuals["account_id"].GetStringValue()
	if accountID == "" {
		return nil, nil
	}

	policyVal := d.KeyColumnQuals["policy"].GetJsonbValue()
	// plugin.Logger(ctx).Trace(fmt.Sprintf("OMERO 1 - Policy: %s", policyVal))

	evaluation, err := EvaluatePolicy(policyVal, accountID)

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