package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsSSMDocument(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_document",
		Description: "AWS SSM Document",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationException", "InvalidDocument"}),
			},
			Hydrate: getAwsSSMDocument,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMDocuments,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "owner", Require: plugin.Optional},
				{Name: "document_type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Systems Manager document.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_ids",
				Description: "The account IDs that have permission to use this document.The ID can be either an AWS account or All.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMDocumentPermissionDetail,
			},
			{
				Name:        "account_sharing_info_list",
				Description: "A list of AWS accounts where the current document is shared and the version shared with each account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMDocumentPermissionDetail,
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
				Name:        "author",
				Description: "The user in your organization who created the document.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_date",
				Description: "The date when the document was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "default_version",
				Description: "The default version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "description",
				Description: "A description of the document.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
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
				Name:        "document_version",
				Description: "The document version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "hash",
				Description: "The Sha256 or Sha1 hash created by the system when the document was created.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "hash_type",
				Description: "The hash type of the document.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "latest_version",
				Description: "The latest version of the document.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "owner",
				Description: "The AWS user account that created the document.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parameters",
				Description: "A description of the parameters for a document.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "pending_review_version",
				Description: "The version of the document that is currently under review.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "platform_types",
				Description: "The operating system platform.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "requires",
				Description: "A list of SSM documents required by a document.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "review_information",
				Description: "Details about the review of a document.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "review_status",
				Description: "The current status of the review.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schema_version",
				Description: "The schema version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "sha1",
				Description: "The SHA1 hash of the document, which you can use for verification.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "status",
				Description: "The user in your organization who created the document.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "status_information",
				Description: "A message returned by AWS Systems Manager that explains the Status value.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with document",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "target_type",
				Description: "The target type which defines the kinds of resources the document can run on.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_name",
				Description: "The version of the artifact associated with the document.",
				Type:        proto.ColumnType_STRING,
			},
			// Standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMDocumentAkas,
				Transform:   transform.FromValue(),
			},
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
		}),
	}
}

//// LIST FUNCTION

func listAwsSSMDocuments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsSSMDocuments")

	// Create session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &ssm.ListDocumentsInput{
		MaxResults: aws.Int64(50),
	}

	filters := buildSsmDocumentFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.ListDocumentsPages(
		input,
		func(page *ssm.ListDocumentsOutput, isLast bool) bool {
			for _, documentIdentifier := range page.DocumentIdentifiers {
				d.StreamListItem(ctx, documentIdentifier)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
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

	var name string
	if h.Item != nil {
		name = documentName(h.Item)
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := SsmService(ctx, d)
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

func getAwsSSMDocumentPermissionDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsSSMDocumentPermissionDetail")

	var name string
	if h.Item != nil {
		name = documentName(h.Item)
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ssm.DescribeDocumentPermissionInput{
		Name:           &name,
		PermissionType: aws.String("Share"),
	}

	// Get call
	data, err := svc.DescribeDocumentPermission(params)
	if err != nil {
		logger.Debug("getAwsSSMDocumentPermissionDetail", "ERROR", err)
		return nil, err
	}

	return data, nil
}

func getAwsSSMDocumentAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsSSMDocumentAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	name := documentName(h.Item)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":document"

	if strings.HasPrefix(name, "/") {
		aka = aka + name
	} else {
		aka = aka + "/" + name
	}

	return []string{aka}, nil
}

func ssmDocumentTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ssmDocumentTagListToTurbotTags")
	data := resourceTags(d.HydrateItem)

	if data == nil {
		return nil, nil
	}
	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if data != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range data {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func documentName(item interface{}) string {
	switch item := item.(type) {
	case *ssm.DocumentDescription:
		return *item.Name
	case *ssm.DocumentIdentifier:
		return *item.Name
	}
	return ""
}

func resourceTags(item interface{}) []*ssm.Tag {
	switch item := item.(type) {
	case *ssm.DocumentDescription:
		return item.Tags
	case *ssm.DocumentIdentifier:
		return item.Tags
	}
	return nil
}

//// UTILITY FUNCTION

// Build ssm documant list call input filter
func buildSsmDocumentFilter(quals plugin.KeyColumnQualMap) []*ssm.DocumentKeyValuesFilter {
	filters := make([]*ssm.DocumentKeyValuesFilter, 0)

	filterQuals := map[string]string{
		"owner":         "Owner",
		"document_type": "DocumentType",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := ssm.DocumentKeyValuesFilter{
				Key: aws.String(filterName),
			}

			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []*string{&val}
			} else {
				filter.Values = value.([]*string)
			}
			filters = append(filters, &filter)
		}
	}
	return filters
}
