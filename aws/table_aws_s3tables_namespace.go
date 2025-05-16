package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3tables"
	"github.com/aws/aws-sdk-go-v2/service/s3tables/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsS3tablesNamespace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3tables_namespace",
		Description: "AWS S3Tables Namespace",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"namespace", "table_bucket_arn"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getS3tablesNamespaceById,
			Tags:    map[string]string{"service": "s3tables", "action": "ListNamespaces"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listS3tablesTableBuckets,
			Hydrate:       listS3tablesNamespaces,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "table_bucket_arn", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "s3tables", "action": "ListNamespaces"},
		},
		GetMatrixItemFunc: S3TablesRegionsMatrix,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "namespace_id",
				Description: "The system-assigned unique identifier for the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Namespace.NamespaceId"),
			},
			{
				Name:        "namespace",
				Description: "The name of the namespace.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFirstIndexNamespaceValue,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "created_at",
				Description: "The date and time the namespace was created at.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Namespace.CreatedAt"),
			},
			{
				Name:        "created_by",
				Description: "The ID of the account that created the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Namespace.CreatedBy"),
			},
			{
				Name:        "owner_account_id",
				Description: "The ID of the account that owns the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Namespace.OwnerAccountId"),
			},
			{
				Name:        "table_bucket_id",
				Description: "The system-assigned unique identifier for the table bucket that contains this namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Namespace.TableBucketId"),
			},
			{
				Name:        "table_bucket_arn",
				Description: "The Amazon Resource Name (ARN) of the table bucket associated with the namespace.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe Standard Columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFirstIndexNamespaceValue,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

// NamespaceInfo holds namespace data along with its parent bucket info
type NamespaceInfo struct {
	Namespace       types.NamespaceSummary
	TableBucketName string
	TableBucketArn  string
}


//// LIST FUNCTION

func listS3tablesNamespaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get bucket details from parent hydrate
	bucket := h.Item.(types.TableBucketSummary)

	// Minimize the number of API calls
	if d.EqualsQualString("table_bucket_arn") != "" && d.EqualsQualString("table_bucket_arn") != *bucket.Arn {
		return nil, nil
	}

	// Create Session
	svc, err := S3TablesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3tables_namespace.listS3tablesNamespaces", "connection_error", err)
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

	input := &s3tables.ListNamespacesInput{
		TableBucketARN: bucket.Arn,
		MaxNamespaces:  aws.Int32(maxLimit),
	}

	paginator := s3tables.NewListNamespacesPaginator(svc, input, func(o *s3tables.ListNamespacesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3tables_namespace.listS3tablesNamespaces", "api_error", err)
			return nil, err
		}

		for _, namespace := range output.Namespaces {
			d.StreamLeafListItem(ctx, &NamespaceInfo{
				Namespace:       namespace,
				TableBucketName: *bucket.Name,
				TableBucketArn:  *bucket.Arn,
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

func getS3tablesNamespaceById(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	namespace := d.EqualsQualString("namespace")
	tableBucketArn := d.EqualsQualString("table_bucket_arn")

	// Create service
	svc, err := S3TablesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3tables_namespace.getS3tablesNamespaceById", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		return nil, nil
	}

	// List namespaces for the table bucket and filter by ID
	input := &s3tables.GetNamespaceInput{
		TableBucketARN: aws.String(tableBucketArn),
		Namespace:      aws.String(namespace),
	}

	output, err := svc.GetNamespace(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3tables_namespace.getS3tablesNamespaceById", "api_error", err)
		return nil, err
	}

	if output != nil {
		return &NamespaceInfo{
			Namespace: types.NamespaceSummary{
				NamespaceId:    output.NamespaceId,
				Namespace:      output.Namespace,
				CreatedAt:      output.CreatedAt,
				CreatedBy:      output.CreatedBy,
				OwnerAccountId: output.OwnerAccountId,
				TableBucketId:  output.TableBucketId,
			},
			TableBucketArn: tableBucketArn,
		}, nil
	}

	return nil, nil
}

func getFirstIndexNamespaceValue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	info := h.Item.(*NamespaceInfo)
	if len(info.Namespace.Namespace) > 0 {
		return info.Namespace.Namespace[0], nil
	}
	return "", nil
}
