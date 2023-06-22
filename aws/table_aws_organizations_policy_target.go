package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsOrganizationsPolicyTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_organizations_policy_target",
		Description: "AWS Organizations Policy Target",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"PolicyNotFoundException", "InvalidInputException"}),
			},
			Hydrate: getOrganizationsPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate:    listOrganizationsPolicyTragets,
			KeyColumns: plugin.AllColumns([]string{"type", "target_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException", "TargetNotFoundException"}),
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name of the policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicySummary.Name"),
			},
			{
				Name:        "id",
				Description: "The unique identifier (ID) of the policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicySummary.Id"),
			},
			{
				Name:        "target_id",
				Description: "The unique identifier (ID) of the root, organizational unit, or account whose policies you want to list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("target_id"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicySummary.Arn"),
			},
			{
				Name:        "type",
				Description: "The type of policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicySummary.Type"),
			},
			{
				Name:        "aws_managed",
				Description: "A boolean value that indicates whether the specified policy is an Amazon Web Services managed policy. If true, then you can attach the policy to roots, OUs, or accounts, but you cannot edit it.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PolicySummary.AwsManaged"),
			},
			{
				Name:        "description",
				Description: "The description of the policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicySummary.Description"),
			},
			{
				Name:        "content",
				Description: "The text content of the policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOrganizationsPolicy,
				Transform:   transform.FromField("Content"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicySummary.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicySummary.Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listOrganizationsPolicyTragets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_policy_target.listOrganizationsPolicyTragets", "client_error", err)
		return nil, err
	}

	policyType := d.EqualsQualString("type")
	targetId := d.EqualsQualString("target_id")

	// Empty Check
	if policyType == "" || targetId == "" {
		return nil, nil
	}

	// Limiting the result
	maxItems := int32(20)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	params := &organizations.ListPoliciesForTargetInput{
		Filter:     types.PolicyType(policyType),
		TargetId:   &targetId,
		MaxResults: &maxItems,
	}

	paginator := organizations.NewListPoliciesForTargetPaginator(svc, params, func(o *organizations.ListPoliciesForTargetPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_organizations_policy_target.listOrganizationsPolicies", "api_error", err)
			return nil, err
		}

		for _, policy := range output.Policies {
			d.StreamListItem(ctx, types.Policy{
				PolicySummary: &types.PolicySummary{
					Arn:         policy.Arn,
					AwsManaged:  policy.AwsManaged,
					Description: policy.Description,
					Id:          policy.Id,
					Name:        policy.Name,
					Type:        policy.Type,
				},
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
