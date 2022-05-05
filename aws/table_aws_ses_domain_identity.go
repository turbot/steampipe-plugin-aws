package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsSESDomainIdentity(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ses_domain_identity",
		Description: "AWS SES Domain Identity",
		List: &plugin.ListConfig{
			Hydrate: listSESDomainIdentities,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The domain identity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "behavior_on_mx_failure",
				Description: "The action that Amazon SES takes if it cannot successfully read the required MX record when you send an email.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityMailFromDomainAttributes,
				Transform:   transform.FromField("BehaviorOnMXFailure"),
			},
			{
				Name:        "dkim_enabled",
				Description: "Denotes if DKIM signing is enabled for email sent from the identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityDkimAttributes,
			},
			{
				Name:        "dkim_tokens",
				Description: "A set of character strings that represent the domain's identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityDkimAttributes,
			},
			{
				Name:        "dkim_verification_status",
				Description: "Describes whether Amazon SES has successfully verified the DKIM DNS records.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityDkimAttributes,
			},
			{
				Name:        "mail_from_domain",
				Description: "The custom MAIL FROM domain that the identity is configured to use.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityMailFromDomainAttributes,
			},
			{
				Name:        "mail_from_domain_status",
				Description: "The state that indicates whether Amazon SES has successfully read the MX record required for custom MAIL FROM domain setup.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityMailFromDomainAttributes,
			},
			{
				Name:        "verification_status",
				Description: "The verification status of the identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityVerificationAttributes,
			},
			{
				Name:        "verification_token",
				Description: "The verification token for a domain identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityVerificationAttributes,
			},
			{
				Name:        "notification_attributes",
				Description: "Represents the notification attributes of an identity.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIdentityNotificationAttributes,
				Transform:   transform.FromValue(),
			},

			// Standard columns for all tables
			{
				Name:        "arn",
				Description: "The ARN of the AWS SES Identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityARN,
				Transform:   transform.FromValue(),
			},
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
				Hydrate:     getIdentityAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listSESDomainIdentities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listSESDomainIdentities")

	region := d.KeyColumnQualString(matrixKeyRegion)

	// Create Session
	svc, err := SESService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &ses.ListIdentitiesInput{IdentityType: &ses.IdentityType_Values()[1]}
	identityResult, err := svc.ListIdentities(input)
	if err != nil {
		return nil, err
	}

	for _, identity := range identityResult.Identities {
		d.StreamListItem(ctx, *identity)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getIdentityDkimAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getSESIdentity")

	name := h.Item.(string)
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Create Session
	svc, err := SESService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	identities := []*string{&name}

	input := &ses.GetIdentityDkimAttributesInput{
		Identities: identities,
	}
	result, err := svc.GetIdentityDkimAttributes(input)
	if err != nil {
		return nil, err
	}
	return result.DkimAttributes[name], err
}

func getIdentityMailFromDomainAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getSESIdentity")

	name := h.Item.(string)
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Create Session
	svc, err := SESService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	identities := []*string{&name}

	input := &ses.GetIdentityMailFromDomainAttributesInput{
		Identities: identities,
	}
	result, err := svc.GetIdentityMailFromDomainAttributes(input)
	if err != nil {
		return nil, err
	}
	return result.MailFromDomainAttributes[name], err
}
