package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/backup/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBackupRecoveryPoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_recovery_point",
		Description: "AWS Backup Recovery Point",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"backup_vault_name", "recovery_point_arn"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"NotFoundException"}),
			},
			Hydrate: getAwsBackupRecoveryPoint,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsBackupVaults,
			Hydrate:       listAwsBackupRecoveryPoints,
			KeyColumns: []*plugin.KeyColumn{
				// API throws error Error: InvalidParameterValueException: Unsupported resource type: arn:aws:ec2:us-east-1::snapshot/snap-03ba1ca215342e331 by passing recovery point arn to list API param
				// Issue raise in AWS SDK V2 https://github.com/aws/aws-sdk-go-v2/issues/1904

				// {
				// 	Name:    "recovery_point_arn",
				// 	Require: plugin.Optional,
				// },
				{
					Name:    "resource_type",
					Require: plugin.Optional,
				},
				{
					Name:    "completion_date",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "backup_vault_name",
				Description: "The name of a logical container where backups are stored.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "recovery_point_arn",
				Description: "An ARN that uniquely identifies a recovery point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The type of Amazon Web Services resource to save as a recovery point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "A status code specifying the state of the recovery point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_size_in_bytes",
				Description: "The size, in bytes, of a backup.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "backup_vault_arn",
				Description: "An ARN that uniquely identifies a backup vault.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "The date and time that a recovery point is created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "completion_date",
				Description: "The date and time that a job to create a recovery point is completed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "encryption_key_arn",
				Description: "The server-side encryption key used to protect your backups.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_role_arn",
				Description: "Specifies the IAM role ARN used to create the target recovery point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_encrypted",
				Description: "A Boolean value that is returned as TRUE if the specified recovery point is encrypted, or FALSE if the recovery point is not encrypted.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "last_restore_time",
				Description: "The date and time that a recovery point was last restored.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "resource_arn",
				Description: "An ARN that uniquely identifies a saved resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_backup_vault_arn",
				Description: "An Amazon Resource Name (ARN) that uniquely identifies the source vault where the resource was originally backed up in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_message",
				Description: "A status message explaining the reason for the recovery point deletion failure.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_class",
				Description: "Specifies the storage class of the recovery point. Valid values are WARM or COLD.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "calculated_lifecycle",
				Description: "An object containing DeleteAt and MoveToColdStorageAt timestamps.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "created_by",
				Description: "Contains identifying information about the creation of a recovery point, including the BackupPlanArn, BackupPlanId, BackupPlanVersion, and BackupRuleId of the backup plan used to create it.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lifecycle",
				Description: "The lifecycle defines when a protected resource is transitioned to cold storage and when it expires.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsBackupRecoveryPoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vault := h.Item.(types.BackupVaultListMember)

	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_recovery_point.listAwsBackupRecoveryPoints", "connection_error", err)
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

	input := &backup.ListRecoveryPointsByBackupVaultInput{
		MaxResults: aws.Int32(maxLimit),
	}
	input.BackupVaultName = vault.BackupVaultName

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	// API throws error Error: InvalidParameterValueException: Unsupported resource type: arn:aws:ec2:us-east-1::snapshot/snap-03ba1ca215342e331 by passing recovery point arn to list API param
	// if equalQuals["recovery_point_arn"] != nil {
	// 	input.ByResourceArn = aws.String(equalQuals["recovery_point_arn"].GetStringValue())
	// }
	if equalQuals["resource_type"] != nil {
		input.ByResourceType = aws.String(equalQuals["resource_type"].GetStringValue())
	}
	if equalQuals["completion_date"] != nil {
		input.ByCreatedAfter = aws.Time(equalQuals["completion_date"].GetTimestampValue().AsTime())
	}

	paginator := backup.NewListRecoveryPointsByBackupVaultPaginator(svc, input, func(o *backup.ListRecoveryPointsByBackupVaultPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "not supported resource type") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_backup_recovery_point.listAwsBackupRecoveryPoints", "api_error", err)
			return nil, err
		}

		for _, item := range output.RecoveryPoints {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getAwsBackupRecoveryPoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_recovery_point.getAwsBackupRecoveryPoint", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var backupVaultName, recoveryPointArn string
	if h.Item != nil {
		backupVaultName = *h.Item.(types.RecoveryPointByBackupVault).BackupVaultName
		recoveryPointArn = *h.Item.(types.RecoveryPointByBackupVault).RecoveryPointArn
	} else {
		backupVaultName = d.KeyColumnQuals["backup_vault_name"].GetStringValue()
		recoveryPointArn = d.KeyColumnQuals["recovery_point_arn"].GetStringValue()
	}

	if recoveryPointArn == "" || backupVaultName == "" {
		return nil, nil
	}

	if arn.IsARN(recoveryPointArn) {
		arnData, _ := arn.Parse(recoveryPointArn)
		// Avoid cross-region queriying
		if arnData.Region != d.KeyColumnQualString(matrixKeyRegion) {
			return nil, nil
		}
	}

	params := &backup.DescribeRecoveryPointInput{
		BackupVaultName:  aws.String(backupVaultName),
		RecoveryPointArn: aws.String(recoveryPointArn),
	}

	detail, err := svc.DescribeRecoveryPoint(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_recovery_point.getAwsBackupRecoveryPoint", "api_error", err)
		return nil, err
	}

	return detail, nil
}
