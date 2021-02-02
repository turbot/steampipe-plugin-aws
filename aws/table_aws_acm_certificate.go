package aws

import (
	"context"
	"log"
	"strings"

	"github.com/turbot/go-kit/types"

	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAcmCertificate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_acm_certificate",
		Description: "AWS ACM Certificate",
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.SingleColumn("certificate_arn"),
			ItemFromKey: certificateFromKey,
			Hydrate:     getAwsAcmCertificateAttributes,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsAcmCertificates,
		},
		FetchMetadata: BuildFetchMetadataList(),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "certificate_arn",
				Description: "Amazon Resource Name (ARN) of the certificate. This is of the form: arn:aws:acm:region:123456789012:certificate/12345678-1234-1234-1234-123456789012",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Certificate.CertificateArn"),
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
				Transform:   transform.FromField("Certificate.DomainName"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "created_at",
				Description: "The time at which the certificate was requested. This value exists only when the certificate type is AMAZON_ISSUED",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Certificate.CreatedAt"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "subject",
				Description: "The name of the entity that is associated with the public key contained in the certificate",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Certificate.Subject"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "issuer",
				Description: "The name of the certificate authority that issued and signed the certificate",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Certificate.Issuer"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "signature_algorithm",
				Description: "The algorithm that was used to sign the certificate",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Certificate.SignatureAlgorithm"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "failure_reason",
				Description: "The reason the certificate request failed. This value exists only when the certificate status is FAILED",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Certificate.FailureReason"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "issued_at",
				Description: "A list of ARNs for the AWS resources that are using the certificate. A certificate can be used by multiple AWS resources",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Certificate.IssuedAt"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "status",
				Description: "The status of the certificate",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Certificate.Status"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "key_algorithm",
				Description: "The algorithm that was used to generate the public-private key pair",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Certificate.KeyAlgorithm"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "not_after",
				Description: "The time after which the certificate is not valid",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Certificate.NotAfter"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "not_before",
				Description: "The time before which the certificate is not valid",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Certificate.NotBefore"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "revocation_reason",
				Description: "The reason the certificate was revoked. This value exists only when the certificate status is REVOKED",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Certificate.RevocationReason"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "revoked_at",
				Description: "The time at which the certificate was revoked. This value exists only when the certificate status is REVOKED",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Certificate.RevokedAt"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "serial",
				Description: "The serial number of the certificate",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Certificate.Serial"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "domain_validation_options",
				Description: "Contains information about the initial validation of each domain name that occurs as a result of the RequestCertificate request. This field exists only when the certificate type is AMAZON_ISSUED",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Certificate.DomainValidationOptions"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "in_use_by",
				Description: "A list of ARNs for the AWS resources that are using the certificate",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Certificate.InUseBy"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "subject_alternative_names",
				Description: "One or more domain names (subject alternative names) included in the certificate. This list contains the domain names that are bound to the public key that is contained in the certificate. The subject alternative names include the canonical domain name (CN) of the certificate and additional domain names that can be used to connect to the website",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Certificate.SubjectAlternativeNames"),
				Hydrate:     getAwsAcmCertificateAttributes,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with certificate",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAcmCertificate,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Certificate.CertificateArn").Transform(certificateArnToTitle),
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
				Transform:   transform.FromField("Certificate.CertificateArn").Transform(arnToAkas),
			},
		}),
	}
}

//// ITEM FROM KEY

func certificateFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	certificateArn := quals["certificate_arn"].GetStringValue()
	item := &acm.DescribeCertificateOutput{
		Certificate: &acm.CertificateDetail{
			CertificateArn: &certificateArn,
		},
	}
	return item, nil
}

//// LIST FUNCTION

func listAwsAcmCertificates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetFetchMetadata(ctx)[fetchMetdataKeyRegion].(string)
	logger.Trace("listAwsAcmCertificates", "AWS_REGION", region)

	// Create service
	svc, err := ACMService(ctx, d.ConnectionManager, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListCertificatesPages(
		&acm.ListCertificatesInput{},
		func(page *acm.ListCertificatesOutput, lastPage bool) bool {
			for _, certificate := range page.CertificateSummaryList {
				d.StreamListItem(ctx, &acm.DescribeCertificateOutput{
					Certificate: &acm.CertificateDetail{
						CertificateArn: certificate.CertificateArn,
					},
				})
			}
			return true
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsAcmCertificateAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetFetchMetadata(ctx)[fetchMetdataKeyRegion].(string)
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsAcmCertificateAttributes")
	item := h.Item.(*acm.DescribeCertificateOutput)
	// defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := ACMService(ctx, d.ConnectionManager, region)
	if err != nil {
		return nil, err
	}

	params := &acm.DescribeCertificateInput{
		CertificateArn: item.Certificate.CertificateArn,
	}

	detail, err := svc.DescribeCertificate(params)
	if err != nil {
		log.Println("[DEBUG] getAwsAcmCertificateAttributes__", "ERROR", err)
		return nil, err
	}
	return detail, nil
}

func getAwsAcmCertificateProperties(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetFetchMetadata(ctx)[fetchMetdataKeyRegion].(string)
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsAcmCertificateProperties")
	item := h.Item.(*acm.DescribeCertificateOutput)
	// defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := ACMService(ctx, d.ConnectionManager, region)
	if err != nil {
		return nil, err
	}

	detail, err := svc.GetCertificate(&acm.GetCertificateInput{
		CertificateArn: item.Certificate.CertificateArn,
	})

	if err != nil {
		return nil, err
	}
	return detail, nil
}

func listTagsForAcmCertificate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetFetchMetadata(ctx)[fetchMetdataKeyRegion].(string)
	logger := plugin.Logger(ctx)
	logger.Trace("listTagsForAcmCertificate")
	item := h.Item.(*acm.DescribeCertificateOutput)
	// defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := ACMService(ctx, d.ConnectionManager, region)
	if err != nil {
		return nil, err
	}

	// Build param
	param := &acm.ListTagsForCertificateInput{
		CertificateArn: item.Certificate.CertificateArn,
	}

	certificateTags, err := svc.ListTagsForCertificate(param)

	if err != nil {
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
	item := types.SafeString(d.Value)
	// Get the resource title
	title := item[strings.LastIndex(item, "/")+1:]

	return title, nil
}
