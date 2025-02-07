package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"

	timestreamwriteEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsTimestreamwriteDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_timestreamwrite_database",
		Description: "AWS Timestreamwrite Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("database_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getAwsTimestreamwriteDatabase,
			Tags:    map[string]string{"service": "timestream-write", "action": "DescribeDatabase"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsTimestreamwriteDatabases,
			Tags:    map[string]string{"service": "timestream-write", "action": "ListDatabases"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(timestreamwriteEndpoint.AWS_INGEST_TIMESTREAM_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "database_name",
				Description: "The name of the Timestream database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name that uniquely identifies this database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time when the database was created, calculated from the Unix epoch time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_updated_time",
				Description: "The last time that this database was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "kms_key_id",
				Description: "The identifier of the KMS key used to encrypt the data stored in the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "table_count",
				Description: "The total number of tables found within a Timestream database.",
				Type:        proto.ColumnType_INT,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseName"),
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

func listAwsTimestreamwriteDatabases(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := TimestreamwriteClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_timestreamwrite_database.listAwsTimestreamwriteDatabases", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil // Unsupported region check
	}

	// Limiting the results
	maxLimit := int32(20)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &timestreamwrite.ListDatabasesInput{
		MaxResults: &maxLimit,
	}

	paginator := timestreamwrite.NewListDatabasesPaginator(svc, input, func(o *timestreamwrite.ListDatabasesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_timestreamwrite_database.listAwsTimestreamwriteDatabases", "api_error", err)
			return nil, err
		}

		for _, item := range output.Databases {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsTimestreamwriteDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	svc, err := TimestreamwriteClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_timestreamwrite_database.getAwsTimestreamwriteDatabase", "connection_error", err)
		return nil, err
	}

	dbName := d.EqualsQualString("database_name")

	// Empty Check
	if dbName == "" {
		return nil, nil
	}

	// Build the params
	params := &timestreamwrite.DescribeDatabaseInput{
		DatabaseName: &dbName,
	}

	// Get call
	op, err := svc.DescribeDatabase(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_timestreamwrite_database.getAwsTimestreamwriteDatabase", "api_error", err)
		return nil, err
	}

	if op.Database != nil {
		return op.Database, nil
	}
	return nil, nil
}
