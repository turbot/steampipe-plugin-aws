package aws

import (
	"context"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/backup"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBackupVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_vault",
		Description: "AWS Backup Vault",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AnyColumn([]string{"name"}),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidParameterValue"}),
			Hydrate:           getAwsBackupVault,
		},
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupVaults,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of a logical container where backups are stored.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupVaultName"),
			},
			{
				Name:        "arn",
				Description: "An Amazon Resource Name (ARN) that uniquely identifies a backup vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupVaultArn"),
			},
			{
				Name:        "creation_date",
				Description: "The date and time a resource backup is created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "creator_request_id",
				Description: "An unique string that identifies the request and allows failed requests to be retried without the risk of running the operation twice.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encryption_key_arn",
				Description: "The server-side encryption key that is used to protect your backups.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "number_of_recovery_points",
				Description: "The number of recovery points that are stored in a backup vault.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "sns_topic_arn",
				Description: "An ARN that uniquely identifies an Amazon Simple Notification Service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SNSTopicArn"),
			},
			{
				Name:        "policy",
				Description: "The backup vault access policy document in JSON format.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsBackupVaultAccessPolicy,
				Transform:   transform.FromField("Policy"),
			},
			{
				Name:        "backup_vault_events",
				Description: "An array of events that indicate the status of jobs to back up resources to the backup vault.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsBackupVaultNotification,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupVaultName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackupVaultArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsBackupVaults(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	plugin.Logger(ctx).Trace("listAwsBackupVaults", "AWS_BACKUP", region)

	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}

	err = svc.ListBackupVaultsPages(
		&backup.ListBackupVaultsInput{},
		func(page *backup.ListBackupVaultsOutput, lastPage bool) bool {
			for _, vault := range page.BackupVaultList {
				d.StreamListItem(ctx, vault)
			}
			return !lastPage
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsBackupVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsBackupVault")

	// Create Session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}

	var name string
	if h.Item != nil {
		vault := h.Item.(*backup.VaultListMember)
		name = *vault.BackupVaultName
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	params := &backup.DescribeBackupVaultInput{
		BackupVaultName: aws.String(name),
	}

	op, err := svc.DescribeBackupVault(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAwsBackupVault", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getAwsBackupVaultNotification(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsBackupVaultNotification")

	// Create Session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}
	name := vaultID(h.Item)

	params := &backup.GetBackupVaultNotificationsInput{
		BackupVaultName: aws.String(name),
	}

	op, err := svc.GetBackupVaultNotifications(params)
	if serverErr, ok := err.(*errors.ServerError); ok {
		if serverErr.ErrorCode() == "InvalidParameterValue" {
			plugin.Logger(ctx).Warn("getAwsBackupVaultNotification", "not_found_error", serverErr, "request", params)
			return nil, nil
		}
		plugin.Logger(ctx).Debug("getAwsBackupVaultNotification", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getAwsBackupVaultAccessPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsBackupVaultAccessPolicy")

	// Create Session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}
	name := vaultID(h.Item)
	params := &backup.GetBackupVaultAccessPolicyInput{
		BackupVaultName: aws.String(name),
	}

	op, err := svc.GetBackupVaultAccessPolicy(params)

	if serverErr, ok := err.(*errors.ServerError); ok {
		if serverErr.ErrorCode() == "InvalidParameterValue" {
			plugin.Logger(ctx).Warn("getAwsBackupVaultAccessPolicy", "not_found_error", serverErr, "request", params)
			return nil, nil
		}
		plugin.Logger(ctx).Debug("getAwsBackupVaultAccessPolicy", "ERROR", err)
		return nil, err
	}
	return op, nil
}

func vaultID(item interface{}) string {
	switch item.(type) {
	case *backup.VaultListMember:
		return *item.(*backup.VaultListMember).BackupVaultName
	case *backup.DescribeBackupVaultOutput:
		return *item.(*backup.DescribeBackupVaultOutput).BackupVaultName
	}
	return ""
}
