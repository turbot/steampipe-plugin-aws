package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
	AttachmentCount *int32
	PolicyGroups    []types.PolicyGroup
	PolicyRoles     []types.PolicyRole
	PolicyUsers     []types.PolicyUser
}

//// LIST FUNCTION

func listIamPolicyAttachments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_policy_attachment.listIamPolicyAttachments", "api_error", err)
		return nil, err
	}
	policy := h.Item.(types.Policy)
	maxItems := int32(100)

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	params := &iam.ListEntitiesForPolicyInput{
		PolicyArn: policy.Arn,
		MaxItems:  aws.Int32(maxItems),
	}

	paginator := iam.NewListEntitiesForPolicyPaginator(svc, params, func(o *iam.ListEntitiesForPolicyPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iam_policy_attachment.listIamPolicyAttachments", "api_error", err)
			return nil, err
		}

		policyAttachment := PolicyAttachment{
			*policy.Arn,
			policy.AttachmentCount,
			output.PolicyGroups,
			output.PolicyRoles,
			output.PolicyUsers,
		}

		d.StreamListItem(ctx, policyAttachment)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
