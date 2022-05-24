package aws

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableAwsResourcePolicyAnalysis(_ context.Context) *plugin.Table {
	// permissionsData = getParliamentIamPermissions()

	return &plugin.Table{
		Name:        "aws_resource_policy_analysis",
		Description: "AWS Resource Policy Analysis",
		List: &plugin.ListConfig{
			// KeyColumns: plugin.KeyColumnSlice{
			// 	{Name: "policy"},
			// },
			Hydrate: listResourcePolicyAnalysis,
		},
		Columns: []*plugin.Column{
			// "Key" Columns
			{
				Name:        "access_level",
				Type:        proto.ColumnType_STRING,
				Description: "Public access levels (based off parliament's levels).",
			},
			{
				Name:        "is_public",
				Type:        proto.ColumnType_BOOL,
				Description: "Public access levels (based off parliament's levels).",
			},
			{
				Name:        "policy",
				Description: "The input policy to be analysed.",
				Type:        proto.ColumnType_JSON,
				// Transform:   transform.FromQual("policy"),
			},
			{
				Name:        "public_statement_ids",
				Type:        proto.ColumnType_JSON,
				Description: "The description for this action.",
			},
			{
				Name:        "public_access_levels",
				Type:        proto.ColumnType_JSON,
				Description: "The description for this action.",
			},
			{
				Name:        "allowed_organization_ids",
				Type:        proto.ColumnType_JSON,
				Description: "The privilege for this action.",
				// Transform:   transform.FromGo(),
			},
			{
				Name:        "allowed_principals",
				Type:        proto.ColumnType_JSON,
				Description: "The access level for this action.",
			},
			{
				Name:        "allowed_principal_account_ids",
				Type:        proto.ColumnType_JSON,
				Description: "The description for this action.",
				// Transform:   transform.FromGo(),
			},
			{
				Name:        "allowed_principal_federated_identities",
				Type:        proto.ColumnType_JSON,
				Description: "The description for this action.",
			},
			{
				Name:        "allowed_principal_services",
				Type:        proto.ColumnType_JSON,
				Description: "The description for this action.",
			},
		},
	}
}

//// LIST FUNCTION

func listResourcePolicyAnalysis(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// policy := d.KeyColumnQuals["policy"].GetJsonbValue()
	policy, err := canonicalPolicy(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"OrganizationAccess","Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test","Condition":{"StringEquals":{"aws:PrincipalOrgID":["o-123456"]}}},{"Sid":"AccountPrincipals","Effect":"Allow","Principal":{"AWS":["arn:aws:iam::123456789012:user/victor@xyz.com","arn:aws:iam::111122223333:root"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"FederatedPrincipals","Effect":"Allow","Principal":{"Federated":"cognito-identity.amazonaws.com"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"ServicePrincipals","Effect":"Allow","Principal":{"Service":["ecs.amazonaws.com","elasticloadbalancing.amazonaws.com"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"PublicAccess","Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"}]}`)

	if err != nil {
		return nil, err
	}

	policyObject, ok := policy.(Policy)
	if !ok {
		return nil, fmt.Errorf("Unable to parse input as policy")
	}

	evaluation, err := policyObject.EvaluatePolicy()
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, evaluation)

	// Context can be cancelled due to manual cancellation or the limit has been hit
	if d.QueryStatus.RowsRemaining(ctx) == 0 {
		return nil, nil

	}
	return nil, nil
}
