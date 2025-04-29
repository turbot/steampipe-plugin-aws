package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"

	lightsailv1 "github.com/aws/aws-sdk-go/service/lightsail"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLightsailCertificate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lightsail_certificate",
		Description: "AWS Lightsail Certificate",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getLightsailCertificate,
			Tags:       map[string]string{"service": "lightsail", "action": "GetCertificates"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidResourceName", "DoesNotExist"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listLightsailCertificates,
			Tags:    map[string]string{"service": "lightsail", "action": "GetCertificates"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(lightsailv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateArn"),
			},
			{
				Name:        "created_at",
				Description: "The timestamp when the certificate was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CertificateDetail.CreatedAt"),
			},
			{
				Name:        "domain_name",
				Description: "The domain name of the certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainName"),
			},
			{
				Name:        "domain_validation_records",
				Description: "An array of objects that describe the domain validation records of the certificate.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CertificateDetail.DomainValidationRecords"),
			},
			{
				Name:        "eligible_to_renew",
				Description: "The renewal eligibility of the certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateDetail.EligibleToRenew"),
			},
			{
				Name:        "in_use_resource_count",
				Description: "The number of Lightsail resources that the certificate is attached to.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CertificateDetail.InUseResourceCount"),
			},
			{
				Name:        "issued_at",
				Description: "The timestamp when the certificate was issued.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CertificateDetail.IssuedAt"),
			},
			{
				Name:        "issuer_ca",
				Description: "The certificate authority that issued the certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateDetail.IssuerCA"),
			},
			{
				Name:        "key_algorithm",
				Description: "The algorithm used to generate the key pair (the public and private key).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateDetail.KeyAlgorithm"),
			},
			{
				Name:        "not_after",
				Description: "The timestamp when the certificate expires.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CertificateDetail.NotAfter"),
			},
			{
				Name:        "not_before",
				Description: "The timestamp when the certificate is first valid.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CertificateDetail.NotBefore"),
			},
			{
				Name:        "renewal_summary",
				Description: "An object that describes the status of the certificate renewal managed by Lightsail.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CertificateDetail.RenewalSummary"),
			},
			{
				Name:        "request_failure_reason",
				Description: "The validation failure reason, if any, of the certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateDetail.RequestFailureReason"),
			},
			{
				Name:        "revocation_reason",
				Description: "The reason the certificate was revoked. This value is present only when the certificate status is REVOKED.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateDetail.RevocationReason"),
			},
			{
				Name:        "revoked_at",
				Description: "The timestamp when the certificate was revoked. This value is present only when the certificate status is REVOKED.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CertificateDetail.RevokedAt"),
			},
			{
				Name:        "serial_number",
				Description: "The serial number of the certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateDetail.SerialNumber"),
			},
			{
				Name:        "status",
				Description: "The validation status of the certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateDetail.Status"),
			},
			{
				Name:        "subject_alternative_names",
				Description: "An array of strings that specify the alternate domains (e.g., example2.com) and subdomains (e.g., blog.example.com) of the certificate.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CertificateDetail.SubjectAlternativeNames"),
			},
			{
				Name:        "support_code",
				Description: "The support code. Include this code in your email to support when you have questions about your Lightsail certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateDetail.SupportCode"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the certificate.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(getLightsailCertificateTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CertificateArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listLightsailCertificates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := LightsailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_certificate.listLightsailCertificates", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &lightsail.GetCertificatesInput{
		IncludeCertificateDetails: true,
	}

	// List call
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		resp, err := svc.GetCertificates(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lightsail_certificate.listLightsailCertificates", "query_error", err)
			return nil, err
		}

		for _, cert := range resp.Certificates {
			// Stream the entire certificate summary which contains both the basic info and the details
			d.StreamListItem(ctx, cert)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextPageToken != nil {
			input.PageToken = resp.NextPageToken
		} else {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getLightsailCertificate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := LightsailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_certificate.getLightsailCertificate", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var name string
	if h.Item != nil {
		cert := h.Item.(types.CertificateSummary)
		name = *cert.CertificateName
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	params := &lightsail.GetCertificatesInput{
		CertificateName:           aws.String(name),
		IncludeCertificateDetails: true,
	}

	detail, err := svc.GetCertificates(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_certificate.getLightsailCertificate", "api_error", err)
		return nil, err
	}

	if len(detail.Certificates) > 0 {
		// Return the entire certificate summary
		return detail.Certificates[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getLightsailCertificateTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.Value.([]types.Tag)
	if tags == nil {
		return nil, nil
	}

	turbotTagsMap := make(map[string]string)
	for _, i := range tags {
		if i.Key != nil && i.Value != nil {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
