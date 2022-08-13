package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
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
				Name:        "identity",
				Description: "The domain identity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "arn",
				Description: "The ARN of the AWS SES identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDomainIdentityARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "verification_status",
				Description: "The verification status of the identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDomainIdentityVerificationAttributes,
				Transform:   transform.FromField("VerificationStatus"),
			},
			{
				Name:        "verification_token",
				Description: "The verification token for a domain identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDomainIdentityVerificationAttributes,
				Transform:   transform.FromField("VerificationToken"),
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
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDomainIdentityARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
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

	input := &ses.ListIdentitiesInput{
		MaxItems:     aws.Int64(1000),
		IdentityType: aws.String(ses.IdentityTypeDomain),
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

func getDomainIdentityVerificationAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getDomainIdentityVerificationAttributes")

	identity := h.Item.(string)
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Create Session
	svc, err := SESService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ses.GetIdentityVerificationAttributesInput{
		Identities: []*string{&identity},
	}
	result, err := svc.GetIdentityVerificationAttributes(input)
	if err != nil {
		return nil, err
	}
	return result.VerificationAttributes[identity], err
}

func getDomainIdentityNotificationAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getDomainIdentityNotificationAttributes")

	identity := h.Item.(string)
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Create Session
	svc, err := SESService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ses.GetIdentityNotificationAttributesInput{
		Identities: []*string{&identity},
	}
	result, err := svc.GetIdentityNotificationAttributes(input)
	if err != nil {
		return nil, err
	}
	return result.NotificationAttributes[identity], err
}

func getDomainIdentityARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDomainIdentityARN")

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
