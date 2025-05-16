package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3tables"
	"github.com/aws/aws-sdk-go-v2/service/s3tables/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsS3tablesTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3tables_table",
		Description: "AWS S3Tables Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace", "table_bucket_arn"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getS3tablesTable,
			Tags:    map[string]string{"service": "s3tables", "action": "GetTable"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listS3tablesTableBuckets,
			Hydrate:       listS3tablesTables,
			Tags:          map[string]string{"service": "s3tables", "action": "ListTables"},
		},
		GetMatrixItemFunc: S3TablesRegionsMatrix,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the table.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TableARN"),
			},
			{
				Name:        "created_at",
				Description: "The date and time the table was created at.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "modified_at",
				Description: "The date and time the table was last modified at.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by",
				Description: "The ID of the account that created the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getS3tablesTable,
			},
			{
				Name:        "modified_by",
				Description: "The ID of the account that last modified the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getS3tablesTable,
			},
			{
				Name:        "namespace",
				Description: "The namespace associated with the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     extractNamespaceNameFromStringArray,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "namespace_id",
				Description: "The unique identifier for the namespace that contains this table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getS3tablesTable,
			},
			{
				Name:        "table_bucket_id",
				Description: "The unique identifier for the table bucket that contains this table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "table_bucket_arn",
				Description: "The Amazon Resource Name (ARN) of the table bucket associated with the table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_account_id",
				Description: "The ID of the account that owns the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getS3tablesTable,
			},
			{
				Name:        "version_token",
				Description: "The version token of the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getS3tablesTable,
			},
			{
				Name:        "metadata_location",
				Description: "The metadata location of the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getS3tablesTable,
			},
			{
				Name:        "format",
				Description: "The format of the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getS3tablesTable,
			},
			{
				Name:        "warehouse_location",
				Description: "The warehouse location of the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getS3tablesTable,
			},
			{
				Name:        "managed_by_service",
				Description: "The service that manages the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getS3tablesTable,
			},

			// Steampipe Standard Columns
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
				Transform:   transform.FromField("TableARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

// TableInfo is a structure to hold S3tables data along with its bucket information
type TableInfo struct {
	// Common fields for both GetTable and ListTables responses
	Name          *string
	TableARN      *string
	TableBucketId *string
	NamespaceId   *string
	Namespace     []string
	CreatedAt     *time.Time
	ModifiedAt    *time.Time
	Type          types.TableType

	// Fields only available in GetTable API response
	CreatedBy         *string
	ModifiedBy        *string
	OwnerAccountId    *string
	Format            *string
	WarehouseLocation *string
	MetadataLocation  *string
	VersionToken      *string
	ManagedByService  *string

	// Parent info
	TableBucketArn string
}

//// LIST FUNCTION

func listS3tablesTables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get bucket details from parent hydrate
	bucket := h.Item.(types.TableBucketSummary)

	// Create Session
	svc, err := S3TablesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3tables_table.listS3tablesTables", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &s3tables.ListTablesInput{
		TableBucketARN: bucket.Arn,
		MaxTables:      aws.Int32(maxLimit),
	}

	paginator := s3tables.NewListTablesPaginator(svc, input, func(o *s3tables.ListTablesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3tables_table.listS3tablesTables", "api_error", err)
			return nil, err
		}

		for _, table := range output.Tables {
			// Create unified TableInfo from TableSummary
			tableInfo := &TableInfo{
				Name:           table.Name,
				TableARN:       table.TableARN,
				TableBucketId:  table.TableBucketId,
				NamespaceId:    table.NamespaceId,
				Namespace:      table.Namespace,
				CreatedAt:      table.CreatedAt,
				ModifiedAt:     table.ModifiedAt,
				Type:           table.Type,
				TableBucketArn: *bucket.Arn,
			}

			d.StreamLeafListItem(ctx, tableInfo)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getS3tablesTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name, namespace, bucketARN string
	var tableInfo *TableInfo

	if h.Item != nil {
		// When called as a hydrate function, get details from the list call
		tableInfo = h.Item.(*TableInfo)
		if tableInfo.Name != nil {
			name = *tableInfo.Name
		}
		if tableInfo.Namespace != nil && len(tableInfo.Namespace) > 0 {
			namespace = tableInfo.Namespace[0]
		}
		bucketARN = tableInfo.TableBucketArn
	} else {
		// When called directly from Get, use query parameters
		name = d.EqualsQualString("name")
		namespace = d.EqualsQualString("namespace")
		bucketARN = d.EqualsQualString("table_bucket_arn")

		// Initialize a new TableInfo for direct Get calls
		tableInfo = &TableInfo{
			TableBucketArn: bucketARN,
		}
	}

	// Create service
	svc, err := S3TablesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3tables_table.getS3tablesTable", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		return nil, nil
	}

	// Build the params
	params := &s3tables.GetTableInput{
		Name:           aws.String(name),
		Namespace:      aws.String(namespace),
		TableBucketARN: aws.String(bucketARN),
	}

	// Get call
	data, err := svc.GetTable(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3tables_table.getS3tablesTable", "api_error", err)
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	// Update the TableInfo with data from GetTable response
	tableInfo.Name = data.Name
	tableInfo.TableARN = data.TableARN
	tableInfo.TableBucketId = data.TableBucketId
	tableInfo.NamespaceId = data.NamespaceId
	tableInfo.Namespace = data.Namespace
	tableInfo.CreatedAt = data.CreatedAt
	tableInfo.ModifiedAt = data.ModifiedAt
	tableInfo.Type = data.Type
	tableInfo.CreatedBy = data.CreatedBy
	tableInfo.ModifiedBy = data.ModifiedBy
	tableInfo.OwnerAccountId = data.OwnerAccountId
	// Convert the Format enum to a string since it's a string type
	tableInfo.Format = aws.String(string(data.Format))
	tableInfo.WarehouseLocation = data.WarehouseLocation
	tableInfo.MetadataLocation = data.MetadataLocation
	tableInfo.VersionToken = data.VersionToken
	tableInfo.ManagedByService = data.ManagedByService

	return tableInfo, nil
}

func extractNamespaceNameFromStringArray(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	info := h.Item.(*TableInfo)
	if len(info.Namespace) > 0 {
		return info.Namespace[0], nil
	}
	return "", nil
}
