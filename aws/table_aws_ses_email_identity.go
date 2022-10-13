package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsSESEmailIdentity(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ses_email_identity",
		Description: "AWS SES Email Identity",
		List: &plugin.ListConfig{
			Hydrate: listSESEmailIdentities,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "identity",
				Description: "The email identity.",
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
			},
			{
				Name:        "verification_token",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release. The verification token for a domain identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSESIdentityVerificationAttributes,
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

func listSESEmailIdentities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := SESClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_email_identity.listSESEmailIdentities", "get_client_error", err)
		return nil, err
	}

	// execute list call
	input := &ses.ListIdentitiesInput{
		MaxItems:     aws.Int32(1),
		IdentityType: types.IdentityTypeEmailAddress,
	}

	// Limiting the results
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxItems {
			if limit < 1 {
				input.MaxItems = aws.Int32(1)
			} else {
				input.MaxItems = aws.Int32(limit)
			}
		}
	}

	// List call
	paginator := ses.NewListIdentitiesPaginator(svc, input, func(o *ses.ListIdentitiesPaginatorOptions) {
		o.Limit = *input.MaxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ses_email_identity.listSESEmailIdentities", "api_error", err)
			return nil, err
		}

		for _, identity := range output.Identities {
			d.StreamListItem(ctx, identity)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSESIdentityVerificationAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identity := h.Item.(string)
	identities := []string{identity}

	// Create Session
	svc, err := SESClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_email_identity.getSESIdentityVerificationAttributes", "get_client_error", err)
		return nil, err
	}

	input := &ses.GetIdentityVerificationAttributesInput{
		Identities: identities,
	}
	result, err := svc.GetIdentityVerificationAttributes(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_email_identity.getSESIdentityVerificationAttributes", "api_error", err)
		return nil, err
	}
	return result.VerificationAttributes[identity], err
}

func getSESIdentityNotificationAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identity := h.Item.(string)
	identities := []string{identity}

	// Create Session
	svc, err := SESClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_email_identity.getSESIdentityNotificationAttributes", "get_client_error", err)
		return nil, err
	}

	input := &ses.GetIdentityNotificationAttributesInput{
		Identities: identities,
	}
	result, err := svc.GetIdentityNotificationAttributes(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_email_identity.getSESIdentityNotificationAttributes", "api_error", err)
		return nil, err
	}
	return result.NotificationAttributes[identity], err
}

func getSESIdentityARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identity := h.Item.(string)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_email_identity.getSESIdentityARN", "api_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":ses:" + region + ":" + commonColumnData.AccountId + ":identity/" + identity
	return arn, nil
}
