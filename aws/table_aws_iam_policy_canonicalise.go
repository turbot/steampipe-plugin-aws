package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/iam"
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
					Name:    "policy_arn",
					Require: plugin.Required,
				},
			},
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "policy_name",
				Description: "The name of the policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy_arn",
				Description: "The arn of the policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Arn"),
			},
			{
				Name:        "policy",
				Description: "Contains the details about the policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPolicyVersion,
				Transform:   transform.FromField("PolicyVersion.Document").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPolicyVersion,
				Transform:   transform.FromField("PolicyVersion.Document").Transform(unescape).Transform(policyToCanonical),
			},
		}),
	}
}

//// LIST FUNCTION

func getIAMPolicyCanonicalise(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIAMPolicyCanonicalise")

	arn := d.KeyColumnQuals["policy_arn"].GetStringValue()

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &iam.GetPolicyInput{
		PolicyArn: &arn,
	}

	op, err := svc.GetPolicy(params)
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, op.Policy)

	return nil, nil
}
