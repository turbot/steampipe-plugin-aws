package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acmpca"
	"github.com/aws/aws-sdk-go-v2/service/acmpca/types"

	acmpcav1 "github.com/aws/aws-sdk-go/service/acmpca"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAcmPcaCertificateAuthority(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_acmpca_certificate_authority",
		Description: "AWS ACM Private certificate authorities",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getAwsAcmPcaCertificateAuthority,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags: map[string]string{"service": "acm-pca", "action": "DescribeCertificateAuthority"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsAcmPcaCertificateAuthorities,
			Tags:    map[string]string{"service": "acm-pca", "action": "ListCertificateAuthorities"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsAcmPcaCertificateAuthority,
				Tags: map[string]string{"service": "acm-pca", "action": "DescribeCertificateAuthority"},
			},
			{
				Func: listTagsForAcmPcaAuthority,
				Tags: map[string]string{"service": "acm-pca", "action": "ListTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(acmpcav1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "Amazon Resource Name (ARN) for your private certificate authority (CA). The format is 12345678-1234-1234-1234-123456789012.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "Date and time at which your private CA was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "failure_reason",
				Description: "Reason the request to create your private CA failed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_storage_security_standard",
				Description: "Defines a cryptographic key management compliance standard used for handling CA keys. Default: FIPS_140_2_LEVEL_3_OR_HIGHER Note: Amazon Web Services Region ap-northeast-3 supports only FIPS_140_2_LEVEL_2_OR_HIGHER. You must explicitly specify this parameter and value when creating a CA in that Region. Specifying a different value (or no value) results in an InvalidArgsException with the message 'A certificate authority cannot be created in this region with the specified security standard.'",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_state_change_at",
				Description: "Date and time at which your private CA was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "not_after",
				Description: "Date and time after which your private CA certificate is not valid.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "not_before",
				Description: "Date and time before which your private CA certificate is not valid.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "owner_account",
				Description: "The Amazon Web Services account ID that owns the certificate authority.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "restorable_until",
				Description: "The period during which a deleted CA can be restored. For more information, see the PermanentDeletionTimeInDays parameter of the DeleteCertificateAuthorityRequest action.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "serial",
				Description: "Serial number of your private CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Status of your private CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Type of your private CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "usage_mode",
				Description: "Specifies whether the CA issues general-purpose certificates that typically require a revocation mechanism, or short-lived certificates that may optionally omit revocation because they expire quickly. Short-lived certificate validity is limited to seven days. The default value is GENERAL_PURPOSE.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_authority_configuration",
				Description: "Your private CA configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "revocation_configuration",
				Description: "Information about the Online Certificate Status Protocol (OCSP) configuration or certificate revocation list (CRL) created and maintained by your private CA.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with private certificate authority (CA).",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAcmPcaAuthority,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe Standard Columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Arn").Transform(authorityArnToTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAcmPcaAuthority,
				Transform:   transform.From(authorityTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsAcmPcaCertificateAuthorities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create service
	svc, err := ACMPCAClient(ctx, d)
	if err != nil {
		logger.Error("aws_acmpca_certificate_authority.listAwsAcmPcaCertificateAuthorities", "connection error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &acmpca.ListCertificateAuthoritiesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := acmpca.NewListCertificateAuthoritiesPaginator(svc, input, func(o *acmpca.ListCertificateAuthoritiesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			logger.Error("aws_acmpca_certificate_authority.listAwsAcmPcaCertificateAuthorities", "api_error", err)
			return nil, err
		}

		for _, item := range output.CertificateAuthorities {
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsAcmPcaCertificateAuthority(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := d.EqualsQualString("arn")

	// Empty check
	if arn == "" {
		return nil, nil
	}

	// Create session
	svc, err := ACMPCAClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_acmpca_certificate_authority.getAwsAcmPcaCertificateAuthority", "connection error", err)
		return nil, err
	}

	params := &acmpca.DescribeCertificateAuthorityInput{
		CertificateAuthorityArn: aws.String(arn),
	}

	detail, err := svc.DescribeCertificateAuthority(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_acmpca_certificate_authority.getAwsAcmPcaCertificateAuthority", "api_error", err)
		return nil, err
	}

	return detail.CertificateAuthority, nil
}

func listTagsForAcmPcaAuthority(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn, err := getAcmPcaAuthorityArn(ctx, d, h)

	if err != nil {
		plugin.Logger(ctx).Error("aws_acmpca_certificate_authority.listTagsForAcmPcaAuthority", "type_error", err)
		return nil, err
	}

	if arn == nil {
		return nil, nil
	}

	// Create session
	svc, err := ACMPCAClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_acmpca_certificate_authority.listTagsForAcmPcaAuthority", "connection_error", err)
		return nil, err
	}

	// Build param
	param := &acmpca.ListTagsInput{
		CertificateAuthorityArn: arn,
	}

	authorityTags, err := svc.ListTags(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_acmpca_certificate_authority.listTagsForAcmPcaAuthority", "api_error", err)
		return nil, err
	}
	return authorityTags, nil
}

func getAcmPcaAuthorityArn(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (*string, error) {
	if h.Item != nil {
		switch item := h.Item.(type) {
		case *types.CertificateAuthority:
			return item.Arn, nil
		case types.CertificateAuthority:
			return item.Arn, nil
		}
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func authorityTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*acmpca.ListTagsOutput)
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}

func authorityArnToTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	item := *d.Value.(*string)

	// Get the resource title
	title := item[strings.LastIndex(item, "/")+1:]

	return title, nil
}
