package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsSESEmailIdentity(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ses_email_identity",
		Description: "AWS SES Email Identity",
		List: &plugin.ListConfig{
			Hydrate: listSESIdentities,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsS3Columns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The user friendly name of the bucket.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "behavior_on_mx_failure",
				Description: "The user friendly name of the bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityMailFromDomainAttributes,
				Transform:   transform.FromField("BehaviorOnMXFailure"),
			},
			{
				Name:        "dkim_enabled",
				Description: "The user friendly name of the bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityDkimAttributes,
			},
			{
				Name:        "dkim_tokens",
				Description: "The user friendly name of the bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityDkimAttributes,
			},
			{
				Name:        "dkim_verification_status",
				Description: "The user friendly name of the bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityDkimAttributes,
			},
			{
				Name:        "mail_from_domain",
				Description: "The user friendly name of the bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityMailFromDomainAttributes,
			},
			{
				Name:        "mail_from_domain_status",
				Description: "The user friendly name of the bucket.",
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
				Description: "The user friendly name of the bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIdentityNotificationAttributes,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "arn",
				Description: "The ARN of the AWS SES Identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityARN,
				Transform:   transform.FromValue(),
			},
			// {
			// 	Name:        "tags_src",
			// 	Description: "A list of tags assigned to bucket.",
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     getBucketTagging,
			// 	Transform:   transform.FromField("TagSet"),
			// },
			// {
			// 	Name:        "tags",
			// 	Description: resourceInterfaceDescription("tags"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     getBucketTagging,
			// 	Transform:   transform.FromField("TagSet").Transform(s3TagsToTurbotTags),
			// },
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
			{
				Name:        "region",
				Description: "The AWS Region in which the resource is located.",
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

// type IdentityInfo struct {
// 	Name   *string
// 	Region *string
// }

//// LIST FUNCTION

func listSESIdentities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listSESIdentities")

	region := d.KeyColumnQualString(matrixKeyRegion)

	// Create Session
	svc, err := SESService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &ses.ListIdentitiesInput{IdentityType: &ses.IdentityType_Values()[0]}
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

//// HYDRATE FUNCTIONS

func getIdentityVerificationAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	// execute list call
	input := &ses.GetIdentityVerificationAttributesInput{
		Identities: identities,
	}
	result, err := svc.GetIdentityVerificationAttributes(input)
	if err != nil {
		return nil, err
	}
	// logger.Trace("==========>>>>>>>>>>>>>>>>>>>", result.VerificationAttributes[name], name, result)
	return result.VerificationAttributes[name], err
}

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

	// execute list call
	input := &ses.GetIdentityDkimAttributesInput{
		Identities: identities,
	}
	result, err := svc.GetIdentityDkimAttributes(input)
	if err != nil {
		return nil, err
	}
	// logger.Trace("==========>>>>>>>>>>>>>>>>>>>", result.VerificationAttributes[name], name, result)
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

	// execute list call
	input := &ses.GetIdentityMailFromDomainAttributesInput{
		Identities: identities,
	}
	result, err := svc.GetIdentityMailFromDomainAttributes(input)
	if err != nil {
		return nil, err
	}
	// logger.Trace("==========>>>>>>>>>>>>>>>>>>>", result.VerificationAttributes[name], name, result)
	return result.MailFromDomainAttributes[name], err
}

func getIdentityNotificationAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	// execute list call
	input := &ses.GetIdentityNotificationAttributesInput{
		Identities: identities,
	}
	result, err := svc.GetIdentityNotificationAttributes(input)
	if err != nil {
		return nil, err
	}
	// logger.Trace("==========>>>>>>>>>>>>>>>>>>>", result.VerificationAttributes[name], name, result)
	return result.NotificationAttributes[name], err
}

//// TRANSFORM FUNCTIONS

func getIdentityARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIdentityARN")

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

func getIdentityAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	arn, err := getIdentityARN(ctx, d, h)
	if err != nil {
		return nil, nil
	}

	return []string{types.SafeString(arn)}, nil
}

/* func sesTagsToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("s3TagsToTurbotTags")
	tags := d.Value.([]*s3.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
*/
