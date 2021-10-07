package aws

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
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
				Transform:   transform.From(iamPolicyTurbotTags),
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
	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := buildIamPolicyFilter(d.KeyColumnQuals, d.Quals)
	input.MaxItems = types.Int64(100)

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxItems {
			input.MaxItems = limit
		}
	}

	// List call
	err = svc.ListPoliciesPages(&input, func(page *iam.ListPoliciesOutput, lastPage bool) bool {
		for _, policy := range page.Policies {
			d.StreamListItem(ctx, policy)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return false
			}
		}

		return !lastPage
	},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIamPolicy")

	var arn string
	if h.Item != nil {
		policy := h.Item.(*iam.Policy)
		arn = *policy.Arn
	} else {
		arn = d.KeyColumnQuals["arn"].GetStringValue()
	}

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

// isPolicyAwsManaged returns true if policy is aws managed
func isPolicyAwsManaged(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("isPolicyAwsManaged")

	policy := h.Item.(*iam.Policy)

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

func iamPolicyTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	policy := d.HydrateItem.(*iam.Policy)
	var turbotTagsMap map[string]string
	if policy.Tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range policy.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

func attachementCountToBool(_ context.Context, d *transform.TransformData) (interface{}, error) {
	attachementCount := types.Int64Value((d.Value.(*int64)))
	if attachementCount == 0 {
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
				input.PathPrefix = types.String(equalQuals[filterQual.ColumnName].GetStringValue())
			case "Bool":
				if filterQual.ColumnName == "is_aws_managed" {
					input.SetScope("Local")
					if equalQuals[filterQual.ColumnName].GetBoolValue() {
						input.SetScope("AWS")
					}
				}
				if filterQual.ColumnName == "is_attached" {
					input.SetOnlyAttached(false)
					if equalQuals[filterQual.ColumnName].GetBoolValue() {
						input.SetOnlyAttached(true)
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
						input.SetScope("Local")
						if !value {
							input.SetScope("AWS")
						}
					}
					if qual == "is_attached" {
						input.SetOnlyAttached(false)
						if !value {
							input.SetOnlyAttached(true)
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
