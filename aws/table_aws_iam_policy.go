package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "is_aws_managed", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "is_attached", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "path", Require: plugin.Optional},
			},
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the iam policy.",
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
				Description: "The path to the policy.",
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
				Hydrate:     isPolicyAwsManaged,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "is_attachable",
				Description: "Specifies whether the policy can be attached to an IAM user, group, or role.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "create_date",
				Description: "The date and time, when the policy was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_date",
				Description: "The date and time, when the policy was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "attachment_count",
				Description: "The number of entities (users, groups, and roles) that the policy is attached to.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "is_attached",
				Description: "Specifies whether the policy is attached to at least one IAM user, group, or role.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AttachmentCount").Transform(attachementCountToBool),
			},
			{
				Name:        "default_version_id",
				Description: "The identifier for the version of the policy that is set as the default version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "permissions_boundary_usage_count",
				Description: "The number of entities (users and roles) for which the policy is used to set the permissions boundary.",
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
			{
				Name:        "tags_src",
				Description: "A list of tags attached with the IAM policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamPolicy,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamPolicy,
				Transform:   transform.From(handleIAMPolicyTurbotTags),
			},
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
	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_policy.listIamPolicies", "client_error", err)
		return nil, err
	}

	params := buildIamPolicyFilter(d.KeyColumnQuals, d.Quals)
	maxItems := int32(100)

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	paginator := iam.NewListPoliciesPaginator(svc, &params, func(o *iam.ListPoliciesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})
	params.MaxItems = aws.Int32(maxItems)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iam_policy.listIamPolicies", "api_error", err)
			return nil, err
		}

		for _, policy := range output.Policies {
			d.StreamListItem(ctx, policy)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		policy := h.Item.(types.Policy)
		arn = *policy.Arn
	} else {
		arn = d.KeyColumnQuals["arn"].GetStringValue()
	}

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_policy.getIamPolicy", "client_error", err)
		return nil, err
	}

	params := &iam.GetPolicyInput{
		PolicyArn: &arn,
	}

	op, err := svc.GetPolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_policy.getIamPolicy", "api_error", err)
		return nil, err
	}

	return *op.Policy, nil
}

func getPolicyVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	policy := h.Item.(types.Policy)

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_policy.getPolicyVersion", "client_error", err)
		return nil, err
	}

	params := &iam.GetPolicyVersionInput{
		PolicyArn: policy.Arn,
		VersionId: policy.DefaultVersionId,
	}

	version, err := svc.GetPolicyVersion(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_policy.getPolicyVersion", "api_error", err)
		return nil, err
	}

	return version, nil
}

// isPolicyAwsManaged returns true if policy is aws managed
func isPolicyAwsManaged(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	policy := h.Item.(types.Policy)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)

	// policy arn for aws managed policy
	// arn:aws-us-gov:iam::aws:policy/aws-service-role/AccessAnalyzerServiceRolePolicy in us gov cloud
	// arn:aws:iam::aws:policy/aws-service-role/AccessAnalyzerServiceRolePolicy in commercial cloud
	if strings.HasPrefix(*policy.Arn, "arn:"+commonColumnData.Partition+":iam::aws:policy") {
		return true, nil
	}

	return false, nil
}

//// TRANSFORM FUNCTIONS

func handleIAMPolicyTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	policy := d.HydrateItem.(types.Policy)
	var turbotTagsMap map[string]string
	if len(policy.Tags) == 0 {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range policy.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

func attachementCountToBool(_ context.Context, d *transform.TransformData) (interface{}, error) {
	attachementCount, ok := d.Value.(*int32)
	if ok && *attachementCount == 0 {
		return false, nil
	}
	return true, nil
}

//// UTILITY FUNCTIONS

func buildIamPolicyFilter(equalQuals plugin.KeyColumnEqualsQualMap, quals plugin.KeyColumnQualMap) iam.ListPoliciesInput {
	input := iam.ListPoliciesInput{}

	filterQuals := FilterQuals{
		FilterQual{"is_aws_managed", "Bool", "Scope"},
		FilterQual{"is_attached", "Bool", "OnlyAttached"},
		FilterQual{"path", "String", "PathPrefix"},
	}

	// EqualsQualMap handling
	for _, filterQual := range filterQuals {
		if equalQuals[filterQual.ColumnName] != nil {
			switch filterQual.ColumnType {
			case "String":
				input.PathPrefix = aws.String(equalQuals[filterQual.ColumnName].GetStringValue())
			case "Bool":
				if filterQual.ColumnName == "is_aws_managed" {
					input.Scope = "Local"
					if equalQuals[filterQual.ColumnName].GetBoolValue() {
						input.Scope = "AWS"
					}
				}
				if filterQual.ColumnName == "is_attached" {
					input.OnlyAttached = false
					if equalQuals[filterQual.ColumnName].GetBoolValue() {
						input.OnlyAttached = true
					}
				}
			}
		}
	}

	boolNEQuals := []string{
		"is_aws_managed",
		"is_attached",
	}
	// Non-Equals Qual Map handling
	for _, qual := range boolNEQuals {
		if quals[qual] != nil {
			for _, q := range quals[qual].Quals {
				value := q.Value.GetBoolValue()
				if q.Operator == "<>" {
					if qual == "is_aws_managed" {
						input.Scope = "Local"
						if !value {
							input.Scope = "AWS"
						}
					}
					if qual == "is_attached" {
						input.OnlyAttached = false
						if !value {
							input.OnlyAttached = true
						}
					}
				}
			}
		}
	}
	return input
}

type FilterQuals []FilterQual

type FilterQual struct {
	ColumnName   string
	ColumnType   string
	PropertyName string
}
