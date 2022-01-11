package aws

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsIamPolicyAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_policy_attachment",
		Description: "AWS IAM Policy Attachment",
		List: &plugin.ListConfig{
			ParentHydrate: listIamPolicies,
			Hydrate:       listIamPolicyAttachments,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "is_attached", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "policy_arn",
				Description: "The Amazon Resource Name (ARN) specifying the IAM policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_attached",
				Description: "Specifies whether the policy is attached to at least one IAM user, group, or role.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AttachmentCount").Transform(attachementCountToBool),
			},
			{
				Name:        "policy_groups",
				Description: "A list of IAM groups that the policy is attached to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy_roles",
				Description: "A list of IAM roles that the policy is attached to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy_users",
				Description: "A list of IAM users that the policy is attached to.",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

type PolicyAttachment struct {
	PolicyArn       string
	AttachmentCount *int64
	PolicyGroups    []*iam.PolicyGroup
	PolicyRoles     []*iam.PolicyRole
	PolicyUsers     []*iam.PolicyUser
}

//// LIST FUNCTION

func listIamPolicyAttachments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}
	policy := h.Item.(*iam.Policy)

	params := &iam.ListEntitiesForPolicyInput{
		PolicyArn: policy.Arn,
		MaxItems:  types.Int64(1000),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxItems {
			if *limit < 1 {
				params.MaxItems = types.Int64(1)
			} else {
				params.MaxItems = limit
			}
		}
	}

	// List call
	err = svc.ListEntitiesForPolicyPages(params, func(page *iam.ListEntitiesForPolicyOutput, lastPage bool) bool {
		d.StreamListItem(ctx, PolicyAttachment{*policy.Arn, policy.AttachmentCount, page.PolicyGroups, page.PolicyRoles, page.PolicyUsers})

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return false
		}
		return !lastPage
	},
	)
	return nil, err
}
