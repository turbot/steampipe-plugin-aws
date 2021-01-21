package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func tableAwsDynamoDBBackup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_backup",
		Description: "AWS DynamoDB Backup",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("arn"),
			ShouldIgnoreError: isNotFoundError([]string{"ValidationException"}),
			ItemFromKey:       tableBackupFromKey,
			Hydrate:           getDynamodbBackup,
		},
		List: &plugin.ListConfig{
			Hydrate: listDynamodbBackups,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the backup",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupName"),
			},
			{
				Name:        "arn",
				Description: "Amazon Resource Name associated with the backup",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupArn"),
			},
			{
				Name:        "table_name",
				Description: "Unique identifier for the table to which backup belongs",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "table_arn",
				Description: "Name of the table to which backup belongs",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "table_id",
				Description: "ARN associated with the table to which backup belongs",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_status",
				Description: "Current status of the backup. Backup can be in one of the following states: CREATING, ACTIVE, DELETED",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_type",
				Description: "Backup type (USER | SYSTEM | AWS_BACKUP)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_creation_datetime",
				Description: "Time at which the backup was created",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("BackupCreationDateTime"),
			},
			{
				Name:        "backup_expiry_datetime",
				Description: "Time at which the automatic on-demand backup created by DynamoDB will expire",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("BackupExpiryDateTime"),
			},
			{
				Name:        "backup_size_bytes",
				Description: "Size of the backup in bytes",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackupArn").Transform(arnToAkas),
			},
		}),
	}
}

//// ITEM FROM KEY

func tableBackupFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	arn := quals["arn"].GetStringValue()
	item := &dynamodb.BackupSummary{
		BackupArn: &arn,
	}

	return item, nil
}

//// LIST FUNCTION

func listDynamodbBackups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listDynamodbBackups", "AWS_REGION", defaultRegion)

	// Create Session
	svc, err := DynamoDbService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	results, err := svc.ListBackups(&dynamodb.ListBackupsInput{})
	if err != nil {
		return nil, err
	}

	if results.BackupSummaries != nil {
		for _, backup := range results.BackupSummaries {
			d.StreamListItem(ctx, backup)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDynamodbBackup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getDynamodbBackup")
	backup := h.Item.(*dynamodb.BackupSummary)
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := DynamoDbService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &dynamodb.DescribeBackupInput{
		BackupArn: backup.BackupArn,
	}

	item, err := svc.DescribeBackup(params)
	if err != nil {
		logger.Debug("getDynamodbBackup__", "ERROR", err)
		return nil, err
	}

	var rowData *dynamodb.BackupSummary

	if item.BackupDescription != nil {
		rowData = &dynamodb.BackupSummary{
			BackupName:             item.BackupDescription.BackupDetails.BackupName,
			BackupArn:              item.BackupDescription.BackupDetails.BackupArn,
			BackupStatus:           item.BackupDescription.BackupDetails.BackupStatus,
			BackupType:             item.BackupDescription.BackupDetails.BackupType,
			BackupSizeBytes:        item.BackupDescription.BackupDetails.BackupSizeBytes,
			BackupCreationDateTime: item.BackupDescription.BackupDetails.BackupCreationDateTime,
			BackupExpiryDateTime:   item.BackupDescription.BackupDetails.BackupExpiryDateTime,
			TableName:              item.BackupDescription.SourceTableDetails.TableName,
			TableArn:               item.BackupDescription.SourceTableDetails.TableArn,
			TableId:                item.BackupDescription.SourceTableDetails.TableId,
		}
	}

	return rowData, nil
}
