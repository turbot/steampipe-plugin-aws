package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsDynamoDBTableExport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_table_export",
		Description: "AWS DynamoDB Table Export",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ExportNotFoundException"}),
			},
			Hydrate: getTableExport,
			Tags:    map[string]string{"service": "dynamodb", "action": "DescribeExport"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listDynamoDBTables,
			Hydrate:       listTableExports,
			Tags:          map[string]string{"service": "dynamodb", "action": "ListExports"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getTableExport,
				Tags: map[string]string{"service": "dynamodb", "action": "DescribeExport"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DYNAMODB_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the export.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExportArn"),
			},
			{
				Name:        "export_status",
				Description: "Export can be in one of the following states: IN_PROGRESS, COMPLETED, or FAILED.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "export_type",
				Description: "The type of export that was performed. Valid values are FULL_EXPORT or INCREMENTAL_EXPORT.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "billed_size_bytes",
				Description: "The billable size of the table export.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTableExport,
			},
			{
				Name:        "client_token",
				Description: "The client token that was provided for the export task. A client token makes calls to ExportTableToPointInTimeInput idempotent, meaning that multiple identical calls have the same effect as one single call.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "end_time",
				Description: "The time at which the export task completed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getTableExport,
			},
			{
				Name:        "export_format",
				Description: "The format of the exported data. Valid values for ExportFormat are DYNAMODB_JSON or ION.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "export_manifest",
				Description: "The name of the manifest file for the export task.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "export_time",
				Description: "Point in time from which table data was exported.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getTableExport,
			},
			{
				Name:        "failure_code",
				Description: "Status code for the result of the failed export.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "failure_message",
				Description: "Export failure reason description.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "item_count",
				Description: "The number of items exported.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTableExport,
			},
			{
				Name:        "s3_bucket",
				Description: "The name of the Amazon S3 bucket containing the export.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "s3_bucket_owner",
				Description: "The ID of the Amazon Web Services account that owns the bucket containing the export.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "s3_prefix",
				Description: "The Amazon S3 bucket prefix used as the file name and path of the exported snapshot.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "s3_sse_algorithm",
				Description: "Type of encryption used on the bucket where export data is stored.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "s3_sse_kms_key_id",
				Description: "The ID of the KMS managed key used to encrypt the S3 bucket where export data is stored (if applicable).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "start_time",
				Description: "The time at which the export task began.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getTableExport,
			},
			{
				Name:        "table_arn",
				Description: "The Amazon Resource Name (ARN) of the table that was exported.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "table_id",
				Description: "Unique ID of the table that was exported.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTableExport,
			},
			{
				Name:        "incremental_export_specification",
				Description: "Optional object containing the parameters specific to an incremental export.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTableExport,
			},

			// Steampipe standard column
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ExportArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listTableExports(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	tableName := h.Item.(types.TableDescription).TableName

	region := d.EqualsQualString(matrixKeyRegion)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table_export.listTableExports", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	tableArn := "arn:" + commonColumnData.Partition + ":dynamodb:" + region + ":" + commonColumnData.AccountId + ":table/" + *tableName

	// Create Session
	svc, err := DynamoDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table_export.listTableExports", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &dynamodb.ListExportsInput{
		MaxResults: aws.Int32(25),
		TableArn:   aws.String(tableArn),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxResults {
			if limit < 1 {
				input.MaxResults = aws.Int32(1)
			} else {
				input.MaxResults = aws.Int32(limit)
			}
		}
	}

	paginator := dynamodb.NewListExportsPaginator(svc, input, func(o *dynamodb.ListExportsPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dynamodb_table_export.listTableExports", "api_error", err)
			return nil, err
		}

		for _, export := range output.ExportSummaries {
			d.StreamListItem(ctx, export)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getTableExport(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.ExportSummary).ExportArn
	} else {
		arn = d.EqualsQuals["arn"].GetStringValue()
	}

	// Create Session
	svc, err := DynamoDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table_export.getTableExport", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &dynamodb.DescribeExportInput{
		ExportArn: &arn,
	}

	op, err := svc.DescribeExport(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table_export.getTableExport", "api_error", err)
		return nil, err
	}

	return op.ExportDescription, nil
}
