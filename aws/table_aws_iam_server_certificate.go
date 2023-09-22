package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamServerCertificate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_server_certificate",
		Description: "AWS IAM Server Certificate",
		List: &plugin.ListConfig{
			Hydrate: listIamServerCertificates,
			Tags:    map[string]string{"service": "iam", "action": "ListServerCertificates"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "path", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIamServerCertificate,
			KeyColumns: plugin.AllColumns([]string{"name"}),
			Tags:    map[string]string{"service": "iam", "action": "GetServerCertificate"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name that identifies the server certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerCertificateMetadata.ServerCertificateName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the server certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerCertificateMetadata.Arn"),
			},
			{
				Name:        "server_certificate_id",
				Description: "The stable and unique string identifying the server certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerCertificateMetadata.ServerCertificateId"),
			},
			{
				Name:        "expiration",
				Description: "The date on which the certificate is set to expire.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ServerCertificateMetadata.Expiration"),
			},
			{
				Name:        "certificate_body",
				Description: "The contents of the public key certificate.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIamServerCertificate,
			},
			{
				Name:        "certificate_chain",
				Description: "The contents of the public key certificate chain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIamServerCertificate,
			},
			{
				Name:        "path",
				Description: "The path to the server certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerCertificateMetadata.Path"),
			},
			{
				Name:        "upload_date",
				Description: "The Amazon Resource Name (ARN) of the account that is designated as the management account for the organization",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ServerCertificateMetadata.UploadDate"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamServerCertificate,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerCertificateMetadata.ServerCertificateName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamServerCertificate,
				Transform:   transform.From(getIamServerCertificateTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServerCertificateMetadata.Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listIamServerCertificates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_server_certificate.listIamServerCertificates", "client_error", err)
		return nil, err
	}

	maxItems := int32(100)

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	params := &iam.ListServerCertificatesInput{MaxItems: &maxItems}

	equalQual := d.EqualsQuals
	if equalQual["path"] != nil {
		params.PathPrefix = aws.String(equalQual["path"].GetStringValue())
	}

	paginator := iam.NewListServerCertificatesPaginator(svc, params, func(o *iam.ListServerCertificatesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iam_server_certificate.listIamServerCertificates", "api_error", err)
			return nil, err
		}

		for _, certificate := range output.ServerCertificateMetadataList {
			d.StreamListItem(ctx, ServerCertificate{
				ServerCertificateMetadata: certificate,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamServerCertificate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIamServerCertificate")

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_server_certificate.getIamServerCertificate", "client_error", err)
		return nil, err
	}

	var name string
	if h.Item != nil {
		name = *h.Item.(ServerCertificate).ServerCertificateMetadata.ServerCertificateName
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	param := &iam.GetServerCertificateInput{ServerCertificateName: aws.String(name)}

	op, err := svc.GetServerCertificate(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_server_certificate.getIamServerCertificate", "api_error", err)
		return nil, err
	}

	return op.ServerCertificate, nil
}

type ServerCertificate struct {
	CertificateBody           *string
	ServerCertificateMetadata types.ServerCertificateMetadata
	CertificateChain          *string
	Tags                      []types.Tag
}

//// TRANSFORM FUNCTIONS

func getIamServerCertificateTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	cert, ok := d.HydrateItem.(*types.ServerCertificate)

	if !ok || len(cert.Tags) == 0 {
		return nil, nil
	}

	turbotTagsMap := map[string]string{}
	for _, i := range cert.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return turbotTagsMap, nil
}
