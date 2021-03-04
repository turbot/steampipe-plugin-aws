package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//type awsIamPolicySimulatorResult struct {
//	PrincipalArn string
//	Action       string
//	ResourceArn  string
//	Decision     *string
//	Result       *iam.EvaluationResult
//}

func tableAwsIamAccountPasswordPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_account_password_policy",
		Description: "AWS IAM Account Password Policy",
		List: &plugin.ListConfig{
			Hydrate: listAccountPasswordPolicies,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "allow_users_to_change_password",
				Description: "Specifies whether IAM users are allowed to change their own password",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "expire_passwords",
				Description: "Indicates whether passwords in the account expire. Returns true if MaxPasswordAge contains a value greater than 0. Returns false if MaxPasswordAge is 0 or not present",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "hard_expiry",
				Description: "Specifies whether IAM users are prevented from setting a new password after",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "max_password_age",
				Description: "The number of days that an IAM user password is valid",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "minimum_password_length",
				Description: "Minimum length to require for IAM user passwords",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "password_reuse_prevention",
				Description: "Specifies the number of previous passwords that IAM users are prevented from reusing",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "require_lowercase_characters",
				Description: "Specifies whether to require lowercase characters for IAM user passwords",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "require_numbers",
				Description: "Specifies whether to require numbers for IAM user passwords",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "require_symbols",
				Description: "Specifies whether to require symbols for IAM user passwords",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "require_uppercase_characters",
				Description: "Specifies whether to require uppercase characters for IAM user passwords",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromGo(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAccountPasswordPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAccountPasswordPolicies")

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	resp, err := svc.GetAccountPasswordPolicy(&iam.GetAccountPasswordPolicyInput{})
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, resp.PasswordPolicy)
	return nil, nil
}
