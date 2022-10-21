package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func tableAwsDynamoDBTableExport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_table_export",
		Description: "AWS DynamoDB Table Export",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException", "ExportNotFoundException"}),
			},
			Hydrate: getTableExport,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listDynamboDbTables,
			Hydrate:       listTableExports,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	tableName := h.Item.(*dynamodb.TableDescription).TableName

	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	tableArn := "arn:" + commonColumnData.Partition + ":dynamodb:" + region + ":" + commonColumnData.AccountId + ":table/" + *tableName

	// Create Session
	svc, err := DynamoDbService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ListExportsInput{
		MaxResults: aws.Int64(25),
		TableArn:   aws.String(tableArn),
	}

	err = svc.ListExportsPages(
		input,
		func(page *dynamodb.ListExportsOutput, lastPage bool) bool {
			for _, export := range page.ExportSummaries {
				d.StreamListItem(ctx, export)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listTableExports", "list", err)
		return nil, err
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getTableExport(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		arn = *h.Item.(*dynamodb.ExportSummary).ExportArn
	} else {
		arn = d.KeyColumnQuals["arn"].GetStringValue()
	}

	// Create Session
	svc, err := DynamoDbService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DescribeExportInput{
		ExportArn: &arn,
	}

	op, err := svc.DescribeExport(input)

	if err != nil {
		plugin.Logger(ctx).Error("getTableExport", "get", err)
		return nil, err
	}

	return op.ExportDescription, nil
}
