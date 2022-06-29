package aws

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAwsIamSamlProvider(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_saml_provider",
		Description: "AWS IAM Saml Provider",
		List: &plugin.ListConfig{
			ParentHydrate: listIamPolicies,
			Hydrate:       listIamSamlProviders,
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

//// LIST FUNCTION

func listIamSamlProviders(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_saml_provider.listIamSamlProviders", "service_creation_error", err)
		return nil, err
	}

	params := &iam.ListSAMLProvidersInput{}

	// List call
	result, err := svc.ListSAMLProviders(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_saml_provider.listIamSamlProviders", "api_error", err)
		return nil, err
	}

	for _, row := range result.SAMLProviderList {
		d.StreamListItem(ctx, row)

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			break
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTION

func getIamSamlProvider(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		entry := h.Item.(*iam.SAMLProviderListEntry)
		arn = *entry.Arn
	} else {
		arn = d.KeyColumnQuals["arn"].GetStringValue()
	}

	if strings.TrimSpace(arn) == "" {
		return nil, nil
	}

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_saml_provider.getIamSamlProvider", "service_creation_error", err)
		return nil, err
	}

	params := &iam.GetSAMLProviderInput{
		SAMLProviderArn: aws.String(arn),
	}

	// List call
	result, err := svc.GetSAMLProvider(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_saml_provider.getIamSamlProvider", "api_error", err)
		return nil, err
	}
	return result, nil
}
