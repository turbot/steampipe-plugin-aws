package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3tables"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsS3tablesTableBucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3tables_table_bucket",
		Description: "AWS S3Tables Table Bucket",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getS3tablesTableBucket,
			Tags:    map[string]string{"service": "s3tables", "action": "GetTableBucket"},
		},
		List: &plugin.ListConfig{
			Hydrate: listS3tablesTableBuckets,
			Tags: map[string]string{"service": "s3tables", "action": "ListTableBuckets"},
		},
		GetMatrixItemFunc: S3TablesRegionsMatrix,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the table bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the table bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The date and time the table bucket was created at.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "table_bucket_id",
				Description: "The system-assigned unique identifier for the table bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_account_id",
				Description: "The ID of the account that owns the table bucket.",
				Type:        proto.ColumnType_STRING,
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
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listS3tablesTableBuckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := S3TablesClient(ctx, d, d.EqualsQualString("region"))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3tables_table_bucket.listS3tablesTableBuckets", "connection_error", err)
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

	input := &s3tables.ListTableBucketsInput{
		MaxBuckets: aws.Int32(maxLimit),
	}

	paginator := s3tables.NewListTableBucketsPaginator(svc, input, func(o *s3tables.ListTableBucketsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3tables_table_bucket.listS3tablesTableBuckets", "api_error", err)
			return nil, err
		}

		for _, bucket := range output.TableBuckets {
			d.StreamListItem(ctx, bucket)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getS3tablesTableBucket(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	arn := d.EqualsQualString("arn")

	// Create service
	svc, err := S3TablesClient(ctx, d, d.EqualsQualString("region"))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3tables_table_bucket.getS3tablesTableBucket", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		return nil, nil
	}

	// Build the params
	params := &s3tables.GetTableBucketInput{
		TableBucketARN: aws.String(arn),
	}

	// Get call
	data, err := svc.GetTableBucket(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3tables_table_bucket.getS3tablesTableBucket", "api_error", err)
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	return data, nil
}
