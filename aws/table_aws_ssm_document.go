package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	ssmEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAwsSSMDocument(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_document",
		Description: "AWS SSM Document",
		Get: &plugin.GetConfig{
			// To avoid the error: get call returned 23 results - the key column is not globally unique, it is recommended to use the "arn" column instead of the "name" column in the "get config" function.
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "InvalidDocument"}),
			},
			Hydrate: getAwsSSMDocument,
			Tags:    map[string]string{"service": "ssm", "action": "DescribeDocument"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMDocuments,
			Tags:    map[string]string{"service": "ssm", "action": "ListDocuments"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "document_type", Require: plugin.Optional},
				{Name: "owner_type", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsSSMDocumentPermissionDetail,
				Tags: map[string]string{"service": "ssm", "action": "DescribeDocumentPermission"},
			},
			{
				Func: getAwsSSMDocument,
				Tags: map[string]string{"service": "ssm", "action": "DescribeDocument"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ssmEndpoint.AWS_SSM_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "display_name",
				Description: "The friendly name of the SSM document.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the Systems Manager document.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the document.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMDocumentArn,
				Transform:   transform.FromValue(),
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
				Name:        "owner_type",
				Description: "The AWS user account type to filter the documents. Possible values: Self, Amazon, Public, Private, ThirdParty, All, Default.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("owner_type"),
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
				Name:        "category",
				Description: "The classification of a document to help you identify and categorize its use.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMDocument,
			},
			{
				Name:        "category_enum",
				Description: "The value that identifies a document's category.",
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

	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_document.listAwsSSMDocuments", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(50)
	input := &ssm.ListDocumentsInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	filters := buildSSMDocumentFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := ssm.NewListDocumentsPaginator(svc, input, func(o *ssm.ListDocumentsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_document.listAwsSSMDocuments", "api_error", err)
			return nil, err
		}

		for _, documentIdentifier := range output.DocumentIdentifiers {
			d.StreamListItem(ctx, documentIdentifier)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsSSMDocument(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var arn string
	if h.Item != nil {
		data, err := getAwsSSMDocumentArn(ctx, d, h)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_document.getAwsSSMDocument", "arn_formatting_error", err)
			return nil, err
		}
		arn = data.(string)
	} else {
		arn = d.EqualsQuals["arn"].GetStringValue()
	}

	matrixRegion := d.EqualsQualString(matrixKeyRegion)
	arnSplit := strings.Split(arn, ":") // Split ARN to get the region

	// Invalid ARN check
	if len(arnSplit) < 3 {
		return nil, nil
	}

	// Skip ARNs in other regions
	if matrixRegion != arnSplit[3] {
		return nil, nil
	}

	name := strings.Split(arn, "/")[1] // Split ARN to get the document name

	// Create Session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_document.getAwsSSMDocument", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Empty name input check
	if strings.TrimSpace(name) == "" {
		return nil, nil
	}

	// Build the params
	params := &ssm.DescribeDocumentInput{
		Name: &name,
	}

	// Get call
	data, err := svc.DescribeDocument(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_document.getAwsSSMDocument", "api_error", err)
		return nil, err
	}

	return data.Document, nil
}

func getAwsSSMDocumentPermissionDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		name = documentName(h.Item)
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_document.getAwsSSMDocumentPermissionDetail", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &ssm.DescribeDocumentPermissionInput{
		Name:           &name,
		PermissionType: types.DocumentPermissionType("Share"),
	}

	// Get call
	data, err := svc.DescribeDocumentPermission(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_document.getAwsSSMDocumentPermissionDetail", "api_error", err)
		return nil, err
	}

	return data, nil
}

func getAwsSSMDocumentAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	name := documentName(h.Item)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_document.getAwsSSMDocumentAkas", "common_data_error", err)
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

func getAwsSSMDocumentArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	name := documentName(h.Item)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_document.getAwsSSMDocumentArn", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":document"

	if strings.HasPrefix(name, "/") {
		arn = arn + name
	} else {
		arn = arn + "/" + name
	}

	return arn, nil
}

func ssmDocumentTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}

	var tags []types.Tag
	switch item := d.HydrateItem.(type) {
	case *types.DocumentDescription:
		tags = item.Tags
	case types.DocumentIdentifier:
		tags = item.Tags
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if len(tags) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func documentName(item interface{}) string {
	switch item := item.(type) {
	case *types.DocumentDescription:
		return *item.Name
	case types.DocumentIdentifier:
		return *item.Name
	}
	return ""
}

//// UTILITY FUNCTION

// Build ssm documant list call input filter
func buildSSMDocumentFilter(quals plugin.KeyColumnQualMap) []types.DocumentKeyValuesFilter {
	filters := make([]types.DocumentKeyValuesFilter, 0)

	filterQuals := map[string]string{
		"owner_type":    "Owner",
		"document_type": "DocumentType",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.DocumentKeyValuesFilter{
				Key: aws.String(filterName),
			}

			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []string{val}
			} else {
				filter.Values = value.([]string)
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
