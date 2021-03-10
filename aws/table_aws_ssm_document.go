package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsSSMDocument(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_document",
		Description: "AWS SSM Document",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"ValidationException", "InvalidDocument"}),
			Hydrate:           getAwsSSMDocument,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMDocuments,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Systems Manager document.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "document_version",
				Description: "The document version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The user in your organization who created the document.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "default_version",
				Description: "The default version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "schema_version",
				Description: "The schema version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "owner",
				Description: "The AWS user account that created the document.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_name",
				Description: "The version of the artifact associated with the document.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "document_format",
				Description: "The document format, either JSON or YAML.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "document_type",
				Description: "The type of document.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "platform_types",
				Description: "The operating system platform.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "review_status",
				Description: "The current status of the review.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_type",
				Description: "The target type which defines the kinds of resources the document can run on.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "approved_version",
				Description: "The version of the document currently approved for use in the organization.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "attachments_information",
				Description: "Details about the document attachments, including names, locations, sizes,and so on.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "created_date",
				Description: "The date when the document was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "latest_version",
				Description: "The latest version of the document.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with document",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(ssmDocumentTagListToTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMDocumentAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSSMDocuments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsSSMDocuments", "AWS_REGION", region)

	// Create session
	svc, err := SsmService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListDocumentsPages(
		&ssm.ListDocumentsInput{},
		func(page *ssm.ListDocumentsOutput, isLast bool) bool {
			for _, documentIdentifier := range page.DocumentIdentifiers {
				d.StreamListItem(ctx, &ssm.DocumentDescription{
					Name:            documentIdentifier.Name,
					TargetType:      documentIdentifier.TargetType,
					DocumentVersion: documentIdentifier.DocumentVersion,
					Owner:           documentIdentifier.Owner,
					VersionName:     documentIdentifier.VersionName,
					DocumentFormat:  documentIdentifier.DocumentFormat,
					PlatformTypes:   documentIdentifier.PlatformTypes,
					ReviewStatus:    documentIdentifier.ReviewStatus,
					Tags:            documentIdentifier.Tags,
					DocumentType:    documentIdentifier.DocumentType,
				})

			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsSSMDocument(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsSSMDocument")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	var name string
	if h.Item != nil {
		name = *h.Item.(*ssm.DocumentDescription).Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := SsmService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ssm.DescribeDocumentInput{
		Name: &name,
	}

	// Get call
	data, err := svc.DescribeDocument(params)
	if err != nil {
		logger.Debug("getAwsSSMDocument", "ERROR", err)
		return nil, err
	}

	return data.Document, nil
}

func getAwsSSMDocumentAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsSSMDocumentAkas")
	name := *h.Item.(*ssm.DocumentDescription).Name
	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":ssm:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":document"

	if strings.HasPrefix(name, "/") {
		aka = aka + name
	} else {
		aka = aka + "/" + name
	}

	return []string{aka}, nil
}

func ssmDocumentTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ssmDocumentTagListToTurbotTags")
	document := d.HydrateItem.(*ssm.DocumentDescription)

	if document.Tags == nil {
		return nil, nil
	}
	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if document != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range document.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
