package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/backup"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBackupVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_vault",
		Description: "AWS Backup Vault",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			// DescribeBackupVault API returns AccessDeniedException instead of a not found error when it is called for vaults that do not exist
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidParameter", "AccessDeniedException"}),
			},
			Hydrate: getAwsBackupVault,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupVaults,
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Hydrate:     getAwsBackupVaultNotification,
				Transform:   transform.FromField("SNSTopicArn"),
			},
			{
				Name:        "policy",
				Description: "The backup vault access policy document in JSON format.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsBackupVaultAccessPolicy,
			},
			{
				Name:        "policy_std",
				Description: "Contains the backup vault access policy document in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsBackupVaultAccessPolicy,
				Transform:   transform.FromField("Policy").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "backup_vault_events",
				Description: "An array of events that indicate the status of jobs to back up resources to the backup vault.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsBackupVaultNotification,
			},

			// Steampipe standard columns
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

func listAwsBackupVaults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &backup.ListBackupVaultsInput{
		MaxResults: aws.Int64(1000),
	}

	// Limiting the results per page
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = types.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	err = svc.ListBackupVaultsPages(
		input,
		func(page *backup.ListBackupVaultsOutput, lastPage bool) bool {
			for _, vault := range page.BackupVaultList {
				d.StreamListItem(ctx, vault)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsBackupVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
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
	// Create Session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}
	name := vaultID(h.Item)

	params := &backup.GetBackupVaultNotificationsInput{
		BackupVaultName: aws.String(name),
	}

	op, err := svc.GetBackupVaultNotifications(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ResourceNotFoundException" || a.Code() == "InvalidParameter" {
				return backup.GetBackupVaultNotificationsOutput{}, nil
			}
			return nil, err
		}
	}
	return op, nil
}

func getAwsBackupVaultAccessPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	name := vaultID(h.Item)
	params := &backup.GetBackupVaultAccessPolicyInput{
		BackupVaultName: aws.String(name),
	}

	op, err := svc.GetBackupVaultAccessPolicy(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ResourceNotFoundException" || a.Code() == "InvalidParameter" {
				return backup.GetBackupVaultAccessPolicyOutput{}, nil
			}
			return nil, err
		}
	}
	return op, nil
}

func vaultID(item interface{}) string {
	switch item := item.(type) {
	case *backup.VaultListMember:
		return *item.BackupVaultName
	case *backup.DescribeBackupVaultOutput:
		return *item.BackupVaultName
	}
	return ""
}
