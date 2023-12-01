package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice/types"

	databasemigrationservicev1 "github.com/aws/aws-sdk-go/service/databasemigrationservice"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsDmsCertificate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dms_certificate",
		Description: "AWS DMS Certificate",
		List: &plugin.ListConfig{
			Hydrate: listDmsCertificates,
			// If the ARN provided as an input parameter refers to a resource that is unavailable in the specified region, the API throws an InvalidParameterValueException exception.
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValueException", "ResourceNotFoundFault"}),
			},
			Tags: map[string]string{"service": "dms", "action": "DescribeCertificates"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "certificate_identifier",
					Require: plugin.Optional,
				},
				{
					Name:    "arn",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getDmsReplicationInstanceTags,
				Tags: map[string]string{"service": "dms", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(databasemigrationservicev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "certificate_identifier",
				Description: "A customer-assigned name for the certificate. Identifiers must begin with a letter and must contain only ASCII letters, digits, and hyphens. They can't end with a hyphen or contain two consecutive hyphens.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateArn"),
			},
			{
				Name:        "certificate_creation_date",
				Description: "The date that the certificate was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "certificate_owner",
				Description: "The owner of the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_pem",
				Description: "The contents of a .pem file, which contains an X.509 certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_wallet",
				Description: "The location of an imported Oracle Wallet certificate for use with SSL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_length",
				Description: "The key length of the cryptographic algorithm being used.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "signing_algorithm",
				Description: "The signing algorithm for the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "valid_from_date",
				Description: "The beginning date that the certificate is valid.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "valid_to_date",
				Description: "The beginning date that the certificate is valid.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags currently associated with the certificate.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDmsCertificateTags,
				Transform:   transform.FromField("TagList"),
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateIdentifier"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDmsCertificateTags,
				Transform:   transform.From(dmsCertificateTagListToTagsMap),
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

func listDmsCertificates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := DatabaseMigrationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_certificate.listDmsCertificates", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	// Build the params
	input := &databasemigrationservice.DescribeCertificatesInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	var filter []types.Filter

	// Additonal Filter
	if d.EqualsQualString("arn") != "" {
		paramFilter := types.Filter{
			Name:   aws.String("certificate-arn"),
			Values: []string{d.EqualsQualString("arn")},
		}
		filter = append(filter, paramFilter)
	}
	if d.EqualsQualString("certificate_identifier") != "" {
		paramFilter := types.Filter{
			Name:   aws.String("certificate-id"),
			Values: []string{d.EqualsQualString("certificate_identifier")},
		}
		filter = append(filter, paramFilter)
	}
	input.Filters = filter

	paginator := databasemigrationservice.NewDescribeCertificatesPaginator(svc, input, func(o *databasemigrationservice.DescribeCertificatesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dms_replication_instance.listDmsReplicationInstances", "api_error", err)
			return nil, err
		}

		for _, items := range output.Certificates {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

func getDmsCertificateTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	certificateArn := h.Item.(types.Certificate).CertificateArn

	// Create service
	svc, err := DatabaseMigrationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_certificate.getDmsCertificateTags", "connection_error", err)
		return nil, err
	}

	params := &databasemigrationservice.ListTagsForResourceInput{
		ResourceArn: certificateArn,
	}

	replicationInstanceTags, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_certificate.getDmsCertificateTags", "api_error", err)
		return nil, err
	}

	return replicationInstanceTags, nil
}

//// TRANSFORM FUNCTIONS

func dmsCertificateTagListToTagsMap(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*databasemigrationservice.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	if data.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range data.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
