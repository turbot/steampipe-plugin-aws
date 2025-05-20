package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsTimestreamwriteTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_timestreamwrite_table",
		Description: "AWS Timestreamwrite Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"table_name", "database_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getAwsTimestreamwriteTable,
			Tags:    map[string]string{"service": "timestream-write", "action": "DescribeTable"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsTimestreamwriteTables,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "database_name", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "timestream-write", "action": "ListTables"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_INGEST_TIMESTREAM_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "table_name",
				Description: "The name of the Timestream table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name that uniquely identifies this table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "table_status",
				Description: "The current state of the table. Possible values are: 'ACTIVE', 'DELETING', or 'RESTORING'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time when the Timestream table was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_updated_time",
				Description: "The time when the Timestream table was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "database_name",
				Description: "The name of the Timestream database that contains this table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schema",
				Description: "The schema of the table.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "magnetic_store_write_properties",
				Description: "Contains properties to set on the table when enabling magnetic store writes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "retention_properties",
				Description: "The retention duration for the memory store and magnetic store.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TableName"),
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

func listAwsTimestreamwriteTables(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := TimestreamwriteClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_timestreamwrite_table.listAwsTimestreamwriteTables", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}
	// Limiting the results
	maxLimit := int32(20)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	dbName := d.EqualsQualString("database_name")

	input := &timestreamwrite.ListTablesInput{
		MaxResults: &maxLimit,
	}

	if dbName != "" {
		input.DatabaseName = &dbName
	}

	paginator := timestreamwrite.NewListTablesPaginator(svc, input, func(o *timestreamwrite.ListTablesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_timestreamwrite_table.listAwsTimestreamwriteTables", "api_error", err)
			return nil, err
		}

		for _, item := range output.Tables {
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

func getAwsTimestreamwriteTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	svc, err := TimestreamwriteClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_timestreamwrite_table.getAwsTimestreamwriteTable", "connection_error", err)
		return nil, err
	}

	tableName := d.EqualsQualString("table_name")
	dbName := d.EqualsQualString("database_name")

	// Empty Check
	if tableName == "" || dbName == "" {
		return nil, nil
	}

	// Build the params
	params := &timestreamwrite.DescribeTableInput{
		DatabaseName: &dbName,
		TableName:    &tableName,
	}

	// Get call
	op, err := svc.DescribeTable(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_timestreamwrite_table.getAwsTimestreamwriteTable", "api_error", err)
		return nil, err
	}

	if op.Table != nil {
		return op.Table, nil
	}
	return nil, nil
}
