package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsOrganizationsServiceControlPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_organizations_service_control_policy",
		Description: "AWS Organizations Service Control Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"PolicyNotFoundException", "InvalidInputException"}),
			},
			Hydrate: getOrganizationsServiceControlPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listOrganizationsServiceControlPolicies,
			// Specifies the type of policy that you want to include in the response. You must
			// specify one of the following values:
			//
			// * AISERVICES_OPT_OUT_POLICY
			// (https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_ai-opt-out.html)
			//
			// *
			// BACKUP_POLICY
			// (https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_backup.html)
			//
			// *
			// SERVICE_CONTROL_POLICY
			// (https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_scp.html)
			//
			// *
			// TAG_POLICY
			// (https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_tag-policies.html)
			KeyColumns: plugin.SingleColumn("type"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException"}),
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
				Hydrate:     getOrganizationsServiceControlPolicy,
				Transform:   transform.FromField("Content"),
			},

			// Standard columns for all tables
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

func listOrganizationsServiceControlPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_account.listOrganizationsServiceControlPolicies", "client_error", err)
		return nil, err
	}

	policyType := d.EqualsQualString("type")

	// Empty Check
	if policyType == "" {
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

	params := &organizations.ListPoliciesInput{
		Filter:     types.PolicyType(policyType),
		MaxResults: &maxItems,
	}

	paginator := organizations.NewListPoliciesPaginator(svc, params, func(o *organizations.ListPoliciesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_organizations_account.listOrganizationsServiceControlPolicies", "api_error", err)
			return nil, err
		}

		for _, policy := range output.Policies {
			d.StreamListItem(ctx, &types.Policy{
				PolicySummary: &policy,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOrganizationsServiceControlPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var policyId string

	if h.Item != nil {
		policyId = *h.Item.(*types.Policy).PolicySummary.Id
	} else {
		policyId = d.EqualsQuals["id"].GetStringValue()
	}

	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_account.getOrganizationsServiceControlPolicy", "client_error", err)
		return nil, err
	}

	params := &organizations.DescribePolicyInput{
		PolicyId: aws.String(policyId),
	}

	op, err := svc.DescribePolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_account.getOrganizationsServiceControlPolicy", "api_error", err)
		return nil, err
	}

	return *op.Policy, nil
}
