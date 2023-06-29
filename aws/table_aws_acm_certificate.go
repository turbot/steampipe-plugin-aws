package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"

	acmv1 "github.com/aws/aws-sdk-go/service/acm"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAcmCertificate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_acm_certificate",
		Description: "AWS ACM Certificate",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("certificate_arn"),
			Hydrate:    getAwsAcmCertificateAttributes,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsAcmCertificates,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "status",
					Require: plugin.Optional,
				},
				{
					Name:    "key_algorithm",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(acmv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "certificate_arn",
				Description: "Amazon Resource Name (ARN) of the certificate. This is of the form: arn:aws:acm:region:123456789012:certificate/12345678-1234-1234-1234-123456789012",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate",
				Description: "The ACM-issued certificate corresponding to the ARN specified as input",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAcmCertificateProperties,
			},
			{
				Name:        "certificate_chain",
				Description: "The ACM-issued certificate corresponding to the ARN specified as input",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAcmCertificateProperties,
			},
			{
				Name:        "domain_name",
				Description: "Fully qualified domain name (FQDN), such as www.example.com or example.com, for the certificate",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_transparency_logging_preference",
				Description: "Indicates whether to opt in to or out of certificate transparency logging. Certificates that are not logged typically generate a browser error. Transparency makes it possible for you to detect SSL/TLS certificates that have been mistakenly or maliciously issued for your domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAcmCertificateAttributes,
				Transform:   transform.FromField("Options.CertificateTransparencyLoggingPreference"),
			},
			{
				Name:        "created_at",
				Description: "The time at which the certificate was requested. This value exists only when the certificate type is AMAZON_ISSUED",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "subject",
				Description: "The name of the entity that is associated with the public key contained in the certificate",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "imported_at",
				Description: "The name of the certificate authority that issued and signed the certificate",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "issuer",
				Description: "The name of the certificate authority that issued and signed the certificate",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "signature_algorithm",
				Description: "The algorithm that was used to sign the certificate",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "extended_key_usages",
				Description: "Specify one or more ExtendedKeyUsage extension values.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "failure_reason",
				Description: "The reason the certificate request failed. This value exists only when the certificate status is FAILED",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "issued_at",
				Description: "A list of ARNs for the AWS resources that are using the certificate. A certificate can be used by multiple AWS resources",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "status",
				Description: "The status of the certificate",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_algorithm",
				Description: "The algorithm that was used to generate the public-private key pair",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "not_after",
				Description: "The time after which the certificate is not valid",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "not_before",
				Description: "The time before which the certificate is not valid",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "renewal_eligibility",
				Description: "Specifies whether the certificate is eligible for renewal.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "revocation_reason",
				Description: "The reason the certificate was revoked. This value exists only when the certificate status is REVOKED",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "revoked_at",
				Description: "The time at which the certificate was revoked. This value exists only when the certificate status is REVOKED",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "serial",
				Description: "The serial number of the certificate",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "type",
				Description: "The source of the certificate. For certificates provided by ACM, this value is AMAZON_ISSUED.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "domain_validation_options",
				Description: "Contains information about the initial validation of each domain name that occurs as a result of the RequestCertificate request. This field exists only when the certificate type is AMAZON_ISSUED",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "in_use_by",
				Description: "A list of ARNs for the AWS resources that are using the certificate",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "subject_alternative_names",
				Description: "One or more domain names (subject alternative names) included in the certificate. This list contains the domain names that are bound to the public key that is contained in the certificate. The subject alternative names include the canonical domain name (CN) of the certificate and additional domain names that can be used to connect to the website",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with certificate",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAcmCertificate,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateArn").Transform(certificateArnToTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAcmCertificate,
				Transform:   transform.From(certificateTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CertificateArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsAcmCertificates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create service
	svc, err := ACMClient(ctx, d)
	if err != nil {
		logger.Error("listAwsAcmCertificates", "connection error", err)
		return nil, err
	}
	// key_algorithm

	keyAlgorithm := d.EqualsQualString("key_algorithm")

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

	input := &acm.ListCertificatesInput{
		MaxItems: aws.Int32(maxLimit),
		Includes: &types.Filters{},
	}

	filter := input.Includes

	if keyAlgorithm != "" {
		filter.KeyTypes = []types.KeyAlgorithm{types.KeyAlgorithm(keyAlgorithm)}
	} else {
		filter.KeyTypes = []types.KeyAlgorithm{
			types.KeyAlgorithmRsa1024,
			types.KeyAlgorithmRsa2048,
			types.KeyAlgorithmRsa3072,
			types.KeyAlgorithmRsa4096,
			types.KeyAlgorithmEcPrime256v1,
			types.KeyAlgorithmEcSecp384r1,
			types.KeyAlgorithmEcSecp521r1,
		}
	}

	input.Includes = filter

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["status"] != nil {
		input.CertificateStatuses = []types.CertificateStatus{
			types.CertificateStatus(equalQuals["status"].GetStringValue()),
		}
	}
	paginator := acm.NewListCertificatesPaginator(svc, input, func(o *acm.ListCertificatesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_acm_certificate.listAwsAcmCertificates", "api_error", err)
			return nil, err
		}

		for _, certificate := range output.CertificateSummaryList {
			d.StreamListItem(ctx, certificate)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsAcmCertificateAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := ACMClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.CertificateSummary).CertificateArn
	} else {
		arn = d.EqualsQuals["certificate_arn"].GetStringValue()
	}

	params := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(arn),
	}

	detail, err := svc.DescribeCertificate(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_acm_certificate.getAwsAcmCertificateAttributes", "api_error", err)
		return nil, err
	}

	// The API documentation (https://docs.aws.amazon.com/acm/latest/APIReference/API_CertificateSummary.html#ACM-Type-CertificateSummary-KeyAlgorithm) specifies that the API should return the response with a separator "_" between the algorithm keys. However, we have observed that it is returning the response with a "-" separator instead.
	if detail != nil && detail.Certificate != nil {
		detail.Certificate.KeyAlgorithm = types.KeyAlgorithm(strings.ReplaceAll(fmt.Sprint(detail.Certificate.KeyAlgorithm), "-", "_"))
	}

	return detail.Certificate, nil
}

func getAwsAcmCertificateProperties(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(*types.CertificateDetail)

	// Create session
	svc, err := ACMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_acm_certificate.getAwsAcmCertificateProperties", "service_client_error", err)
		return nil, err
	}

	detail, err := svc.GetCertificate(ctx, &acm.GetCertificateInput{
		CertificateArn: item.CertificateArn,
	})

	if err != nil {
		plugin.Logger(ctx).Error("aws_acm_certificate.getAwsAcmCertificateProperties", "api_error", err)
		return nil, err
	}
	return detail, nil
}

func listTagsForAcmCertificate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(*types.CertificateDetail)

	// Create session
	svc, err := ACMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_acm_certificate.listTagsForAcmCertificate", "service_client_error", err)
		return nil, err
	}

	// Build param
	param := &acm.ListTagsForCertificateInput{
		CertificateArn: item.CertificateArn,
	}

	certificateTags, err := svc.ListTagsForCertificate(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_acm_certificate.listTagsForAcmCertificate", "api_error", err)
		return nil, err
	}
	return certificateTags, nil
}

//// TRANSFORM FUNCTIONS

func certificateTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*acm.ListTagsForCertificateOutput)
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}

func certificateArnToTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	item := *d.Value.(*string)

	// Get the resource title
	title := item[strings.LastIndex(item, "/")+1:]

	return title, nil
}
