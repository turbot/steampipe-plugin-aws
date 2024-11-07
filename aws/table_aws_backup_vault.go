package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/backup/types"

	backupv1 "github.com/aws/aws-sdk-go/service/backup"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameter", "AccessDeniedException"}),
			},
			Hydrate: getAwsBackupVault,
			Tags:    map[string]string{"service": "backup", "action": "DescribeBackupVault"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupVaults,
			Tags:    map[string]string{"service": "backup", "action": "ListBackupVaults"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsBackupVaultNotification,
				Tags: map[string]string{"service": "backup", "action": "GetBackupVaultNotifications"},
			},
			{
				Func: getAwsBackupVaultAccessPolicy,
				Tags: map[string]string{"service": "backup", "action": "GetBackupVaultAccessPolicy"},
			},
			{
				Func: getAwsBackupVaultTags,
				Tags: map[string]string{"service": "backup", "action": "ListTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(backupv1.EndpointsID),
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
				Name:        "lock_date",
				Description: "The date and time when Backup Vault Lock configuration cannot be changed or deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "locked",
				Description: "A Boolean that indicates whether Backup Vault Lock is currently protecting the backup vault. True means that Vault Lock causes delete or update operations on the recovery points stored in the vault to fail.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "max_retention_days",
				Description: "The Backup Vault Lock setting that specifies the maximum retention period that the vault retains its recovery points.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "min_retention_days",
				Description: "The Backup Vault Lock setting that specifies the minimum retention period that the vault retains its recovery points.",
				Type:        proto.ColumnType_INT,
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
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsBackupVaultTags,
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
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_vault.listAwsBackupVaults", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &backup.ListBackupVaultsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := backup.NewListBackupVaultsPaginator(svc, input, func(o *backup.ListBackupVaultsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_backup_vault.listAwsBackupVaults", "api_error", err)
			return nil, err
		}

		for _, items := range output.BackupVaultList {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsBackupVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_vault.getAwsBackupVault", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var name string
	if h.Item != nil {
		vault := h.Item.(types.BackupVaultListMember)
		name = *vault.BackupVaultName
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	params := &backup.DescribeBackupVaultInput{
		BackupVaultName: aws.String(name),
	}

	op, err := svc.DescribeBackupVault(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_vault.getAwsBackupVault", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsBackupVaultNotification(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_vault.getAwsBackupVaultNotification", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	name := vaultID(h.Item)

	params := &backup.GetBackupVaultNotificationsInput{
		BackupVaultName: aws.String(name),
	}

	op, err := svc.GetBackupVaultNotifications(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "Failed reading notifications from database for Backup vault") {
			return &backup.GetBackupVaultNotificationsOutput{}, nil
		}

		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" || ae.ErrorCode() == "InvalidParameter" {
				return &backup.GetBackupVaultNotificationsOutput{}, nil
			}
		}
		plugin.Logger(ctx).Error("aws_backup_vault.getAwsBackupVaultNotification", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getAwsBackupVaultTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := vaultArn(h.Item)

	// Define the regex pattern for the backup vault ARN
	pattern := `arn:aws:backup:[a-z0-9\-]+:[0-9]{12}:backup-vault:.*`

	return getAwsBackupResourceTags(ctx, d, arn, pattern)
}

func getAwsBackupVaultAccessPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_vault.getAwsBackupVaultAccessPolicy", "connection_error", err)
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

	op, err := svc.GetBackupVaultAccessPolicy(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" || ae.ErrorCode() == "InvalidParameter" {
				return backup.GetBackupVaultAccessPolicyOutput{}, nil
			}
		}
		plugin.Logger(ctx).Error("aws_backup_vault.getAwsBackupVaultAccessPolicy", "api_error", err)
		return nil, err
	}

	return op, nil
}

func vaultID(item interface{}) string {
	switch item := item.(type) {
	case types.BackupVaultListMember:
		return *item.BackupVaultName
	case *backup.DescribeBackupVaultOutput:
		return *item.BackupVaultName
	}
	return ""
}

func vaultArn(item interface{}) string {
	switch item := item.(type) {
	case types.BackupVaultListMember:
		return *item.BackupVaultArn
	case *backup.DescribeBackupVaultOutput:
		return *item.BackupVaultArn
	}
	return ""
}
