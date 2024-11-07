package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/backup/types"

	backupv1 "github.com/aws/aws-sdk-go/service/backup"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBackupRecoveryPoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_recovery_point",
		Description: "AWS Backup Recovery Point",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"backup_vault_name", "recovery_point_arn"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getAwsBackupRecoveryPoint,
			Tags:    map[string]string{"service": "backup", "action": "DescribeRecoveryPoint"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsBackupVaults,
			Hydrate:       listAwsBackupRecoveryPoints,
			Tags:          map[string]string{"service": "backup", "action": "ListRecoveryPointsByBackupVault"},
			KeyColumns: []*plugin.KeyColumn{
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
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsBackupRecoveryPoint,
				Tags: map[string]string{"service": "backup", "action": "DescribeRecoveryPoint"},
			},
			{
				Func: getAwsBackupRecoveryPointTags,
				Tags: map[string]string{"service": "backup", "action": "ListTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(backupv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "backup_vault_name",
				Description: "The name of a logical container where backups are stored.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "composite_member_identifier",
				Description: "This is the identifier of a resource within a composite group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parent_recovery_point_arn",
				Description: "This is the Amazon Resource Name (ARN) of the parent (composite) recovery point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_name",
				Description: "This is the non-unique name of the resource that belongs to the specified backup.",
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
				Name:        "is_parent",
				Description: "This is a boolean value indicating this is a parent (composite) recovery point.",
				Type:        proto.ColumnType_BOOL,
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
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsBackupRecoveryPointTags,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceArn").Transform(transform.EnsureStringArray),
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
	equalQuals := d.EqualsQuals
	// The ListRecoveryPointsByBackupVault returns results with ARNs like arn:aws:ec2:us-east-1::snapshot/snap-03ba1ca215342e331, but when trying
	// to pass this value in, the API throws "Error: InvalidParameterValueException: Unsupported resource type: arn:aws:ec2:us-east-1::snapshot/snap-03ba1ca215342e331"
	// Raised https://github.com/aws/aws-sdk-go-v2/issues/1904 to better understand what to pass in
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
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

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
			if d.RowsRemaining(ctx) == 0 {
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
		backupVaultName = d.EqualsQuals["backup_vault_name"].GetStringValue()
		recoveryPointArn = d.EqualsQuals["recovery_point_arn"].GetStringValue()
	}

	if recoveryPointArn == "" || backupVaultName == "" {
		return nil, nil
	}

	if arn.IsARN(recoveryPointArn) {
		arnData, _ := arn.Parse(recoveryPointArn)
		// Avoid cross-region queriying
		if arnData.Region != d.EqualsQualString(matrixKeyRegion) {
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

func getAwsBackupRecoveryPointTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := recoveryPointArn(h.Item)

	// Define the regex pattern for the backup recovery point ARN
	pattern := `arn:aws:backup:[a-z0-9\-]+:[0-9]{12}:recovery-point:.*`

	return getAwsBackupResourceTags(ctx, d, arn, pattern)
}

func recoveryPointArn(item interface{}) string {
	switch item := item.(type) {
	case types.RecoveryPointByBackupVault:
		return *item.RecoveryPointArn
	case *backup.DescribeRecoveryPointOutput:
		return *item.RecoveryPointArn
	}
	return ""
}
