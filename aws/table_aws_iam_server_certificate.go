package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamServerCertificate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_server_certificate",
		Description: "AWS IAM Server Certificate",
		List: &plugin.ListConfig{
			Hydrate: listIamServerCertificates,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIamServerCertificate,
			KeyColumns: plugin.AllColumns([]string{"name"}),
		},
		Columns: awsColumns([]*plugin.Column{
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

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &iam.ListServerCertificatesInput{
		MaxItems: aws.Int64(1000),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxItems {
			if *limit < 1 {
				input.MaxItems = aws.Int64(1)
			} else {
				input.MaxItems = limit
			}
		}
	}

	err = svc.ListServerCertificatesPages(
		input,
		func(page *iam.ListServerCertificatesOutput, lastPage bool) bool {
			for _, certificate := range page.ServerCertificateMetadataList {
				d.StreamListItem(ctx, &iam.ServerCertificate{
					ServerCertificateMetadata: certificate,
				})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)
	if err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamServerCertificate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIamServerCertificate")

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	var name string
	if h.Item != nil {
		name = *h.Item.(*iam.ServerCertificate).ServerCertificateMetadata.ServerCertificateName
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	param := &iam.GetServerCertificateInput{
		ServerCertificateName: aws.String(name),
	}

	op, err := svc.GetServerCertificate(param)
	if err != nil {
		return nil, err
	}

	return op.ServerCertificate, nil
}
