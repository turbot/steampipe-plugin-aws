package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/turbot/go-kit/types"
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
	logger := plugin.Logger(ctx)
	logger.Trace("listSESEmailIdentities")

	// Create Session
	svc, err := SESService(ctx, d)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &ses.ListIdentitiesInput{
		MaxItems:     aws.Int64(1000),
		IdentityType: aws.String(ses.IdentityTypeEmailAddress),
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxItems {
			if *limit < 1 {
				input.MaxItems = types.Int64(1)
			} else {
				input.MaxItems = limit
			}
		}
	}

	// List call
	err = svc.ListIdentitiesPages(
		input,
		func(page *ses.ListIdentitiesOutput, lastPage bool) bool {
			for _, identity := range page.Identities {
				d.StreamListItem(ctx, *identity)

				// Context can be cancelled due to manual cancellation or the limit has been hit
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

func getSESIdentityVerificationAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getSESIdentityVerificationAttributes")

	identity := h.Item.(string)
	identities := []*string{&identity}

	// Create Session
	svc, err := SESService(ctx, d)
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
	return result.VerificationAttributes[identity], err
}

func getSESIdentityNotificationAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getSESIdentityNotificationAttributes")

	identity := h.Item.(string)
	identities := []*string{&identity}

	// Create Session
	svc, err := SESService(ctx, d)
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
	return result.NotificationAttributes[identity], err
}

func getSESIdentityARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSESIdentityARN")

	identity := h.Item.(string)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":ses:" + region + ":" + commonColumnData.AccountId + ":identity/" + identity
	return arn, nil
}
