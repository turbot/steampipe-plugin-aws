package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsIamPolicyCanonicalise(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_policy_canonicalise",
		Description: "AWS IAM Policy Canonicalise",
		List: &plugin.ListConfig{
			Hydrate: getIAMPolicyCanonicalise,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "policy",
					Require: plugin.Required,
				},
			},
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "policy",
				Description: "Contains the details about the policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "policy_jsonb",
				Description: "Contains the details about the policy in jsonb format.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue().Transform(unescape).Transform(policyToCanonical),
			},
		}),
	}
}

//// LIST FUNCTION

func getIAMPolicyCanonicalise(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIAMPolicyCanonicalise")

	policy := d.KeyColumnQuals["policy"].GetStringValue()

	d.StreamListItem(ctx, policy)

	return nil, nil
}
