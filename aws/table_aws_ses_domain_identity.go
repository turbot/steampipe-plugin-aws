package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/turbot/go-kit/types"
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
				Hydrate:     getDomainIdentityMailFromDomainAttributes,
				Transform:   transform.FromField("BehaviorOnMXFailure"),
			},
			{
				Name:        "dkim_enabled",
				Description: "Denotes if DKIM signing is enabled for email sent from the identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDomainIdentityDkimAttributes,
			},
			{
				Name:        "dkim_tokens",
				Description: "A set of character strings that represent the domain's identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDomainIdentityDkimAttributes,
			},
			{
				Name:        "dkim_verification_status",
				Description: "Describes whether Amazon SES has successfully verified the DKIM DNS records.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDomainIdentityDkimAttributes,
			},
			{
				Name:        "mail_from_domain",
				Description: "The custom MAIL FROM domain that the identity is configured to use.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDomainIdentityMailFromDomainAttributes,
			},
			{
				Name:        "mail_from_domain_status",
				Description: "The state that indicates whether Amazon SES has successfully read the MX record required for custom MAIL FROM domain setup.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDomainIdentityMailFromDomainAttributes,
			},
			{
				Name:        "verification_status",
				Description: "The verification status of the identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDomainIdentityVerificationAttributes,
			},
			{
				Name:        "verification_token",
				Description: "The verification token for a domain identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDomainIdentityVerificationAttributes,
			},
			{
				Name:        "notification_attributes",
				Description: "Represents the notification attributes of an identity.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDomainIdentityNotificationAttributes,
				Transform:   transform.FromValue(),
			},

			// Standard columns for all tables
			{
				Name:        "arn",
				Description: "The ARN of the AWS SES Identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDomainIdentityARN,
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
				Hydrate:     getDomainIdentityAkas,
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

func getDomainIdentityVerificationAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getDomainIdentityVerificationAttributes")

	name := h.Item.(string)
	region := d.KeyColumnQualString(matrixKeyRegion)
	identities := []*string{&name}

	// Create Session
	svc, err := SESService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ses.GetIdentityVerificationAttributesInput{
		Identities: identities,
	}
	result, err := svc.GetIdentityVerificationAttributes(input)
	if err != nil {
		return nil, err
	}
	return result.VerificationAttributes[name], err
}

func getDomainIdentityNotificationAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getDomainIdentityNotificationAttributes")

	name := h.Item.(string)
	region := d.KeyColumnQualString(matrixKeyRegion)
	identities := []*string{&name}

	// Create Session
	svc, err := SESService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ses.GetIdentityNotificationAttributesInput{
		Identities: identities,
	}
	result, err := svc.GetIdentityNotificationAttributes(input)
	if err != nil {
		return nil, err
	}
	return result.NotificationAttributes[name], err
}

func getDomainIdentityDkimAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getDomainIdentityDkimAttributes")

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

func getDomainIdentityMailFromDomainAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getDomainIdentityMailFromDomainAttributes")

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

//// TRANSFORM FUNCTIONS

func getDomainIdentityARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDomainIdentityARN")

	name := h.Item.(string)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":ses:" + region + ":" + commonColumnData.AccountId + ":identity/" + name
	return arn, nil
}

func getDomainIdentityAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn, err := getEmailIdentityARN(ctx, d, h)
	if err != nil {
		return nil, nil
	}
	return []string{types.SafeString(arn)}, nil
}
