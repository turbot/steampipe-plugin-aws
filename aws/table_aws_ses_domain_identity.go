package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"

	sesEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSESDomainIdentity(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ses_domain_identity",
		Description: "AWS SES Domain Identity",
		List: &plugin.ListConfig{
			Hydrate: listSESDomainIdentities,
			Tags:    map[string]string{"service": "ses", "action": "ListIdentities"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getSESDomainIdentityDkimAttributes,
				Tags: map[string]string{"service": "ses", "action": "GetIdentityDkimAttributes"},
			},
			{
				Func: getSESDomainIdentityMailFromDomainAttributes,
				Tags: map[string]string{"service": "ses", "action": "GetIdentityMailFromDomainAttributes"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(sesEndpoint.EMAILServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "identity",
				Description: "The domain identity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "arn",
				Description: "The ARN of the AWS SES identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSESIdentityARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "verification_status",
				Description: "The verification status of the identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSESIdentityVerificationAttributes,
				Transform:   transform.FromField("VerificationStatus"),
			},
			{
				Name:        "verification_token",
				Description: "The verification token for a domain identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSESIdentityVerificationAttributes,
				Transform:   transform.FromField("VerificationToken"),
			},
			{
				Name:        "dkim_attributes",
				Description: "The DKIM attributes for an email address or a domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSESDomainIdentityDkimAttributes,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "identity_mail_from_domain_attributes",
				Description: "The custom MAIL FROM attributes for a list of identities.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSESDomainIdentityMailFromDomainAttributes,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "notification_attributes",
				Description: "Represents the notification attributes of an identity.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSESIdentityNotificationAttributes,
				Transform:   transform.FromValue(),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSESIdentityARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSESDomainIdentities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := SESClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_domain_identity.listSESDomainIdentities", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(1000)
	// Limiting the results
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = 1
			} else {
				maxItems = limit
			}
		}
	}

	input := &ses.ListIdentitiesInput{
		MaxItems:     &maxItems,
		IdentityType: types.IdentityTypeDomain,
	}

	// List call
	paginator := ses.NewListIdentitiesPaginator(svc, input, func(o *ses.ListIdentitiesPaginatorOptions) {
		o.Limit = *input.MaxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ses_domain_identity.listSESDomainIdentities", "api_error", err)
			return nil, err
		}

		for _, identity := range output.Identities {
			d.StreamListItem(ctx, identity)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATED CALLS

func getSESDomainIdentityDkimAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identity := h.Item.(string)

	// Create Client
	svc, err := SESClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_domain_identity.getSESDomainIdentityDkimAttributes", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &ses.GetIdentityDkimAttributesInput{
		Identities: []string{identity},
	}

	op, err := svc.GetIdentityDkimAttributes(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_domain_identity.getSESDomainIdentityDkimAttributes", "api_error", err)
		return nil, err
	}

	if op != nil && op.DkimAttributes != nil {
		return op.DkimAttributes, nil
	}

	return nil, nil
}

func getSESDomainIdentityMailFromDomainAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identity := h.Item.(string)

	// Create Client
	svc, err := SESClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_domain_identity.getSESDomainIdentityMailFromDomainAttributes", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &ses.GetIdentityMailFromDomainAttributesInput{
		Identities: []string{identity},
	}

	op, err := svc.GetIdentityMailFromDomainAttributes(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_domain_identity.getSESDomainIdentityMailFromDomainAttributes", "api_error", err)
		return nil, err
	}

	if op != nil && op.MailFromDomainAttributes != nil {
		return op.MailFromDomainAttributes, nil
	}

	return nil, nil
}
