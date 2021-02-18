package aws

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsIamPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_policy",
		Description: "AWS IAM Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"arn"}),
			Hydrate:    getIamPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listIamPolicies,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the iam policy",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicyName"),
			},
			{
				Name:        "policy_id",
				Description: "The stable and unique string identifying the policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path",
				Description: "The path to the policy",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the iam policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_aws_managed",
				Description: "Specifies whether the policy is AWS Managed or Customer Managed. If true policy is aws managed otherwise customer managed.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.From(isPolicyAwsManaged),
			},
			{
				Name:        "is_attachable",
				Description: "Specifies whether the policy can be attached to an IAM user, group, or role.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "create_date",
				Description: "The date and time, when the policy was created",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_date",
				Description: "The date and time, when the policy was last updated",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "attachment_count",
				Description: "The number of entities (users, groups, and roles) that the policy is attached to",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "default_version_id",
				Description: "The identifier for the version of the policy that is set as the default version",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "permissions_boundary_usage_count",
				Description: "The number of entities (users and roles) for which the policy is used to set the permissions boundary",
				Type:        proto.ColumnType_INT,
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

			// Standard columns for all tables
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
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listIamPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	err = svc.ListPoliciesPages(
		&iam.ListPoliciesInput{},
		func(page *iam.ListPoliciesOutput, lastPage bool) bool {
			for _, policy := range page.Policies {
				d.StreamListItem(ctx, policy)
			}
			return true
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIamPolicy")
	arn := d.KeyColumnQuals["arn"].GetStringValue()

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

	return op.Policy, nil
}

func getPolicyVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPolicyVersion")
	policy := h.Item.(*iam.Policy)

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &iam.GetPolicyVersionInput{
		PolicyArn: policy.Arn,
		VersionId: policy.DefaultVersionId,
	}

	version, err := svc.GetPolicyVersion(params)
	if err != nil {
		return nil, err
	}

	return version, nil
}

//// Transform Function

// isPolicyAwsManaged returns true if policy is aws managed
func isPolicyAwsManaged(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("isPolicyAwsManaged")

	policy := d.HydrateItem.(*iam.Policy)

	// TODO arn:{{ partition }}:iam::aws:policy
	if strings.HasPrefix(*policy.Arn, "arn:aws:iam::aws:policy") {
		return true, nil
	}

	return false, nil
}
