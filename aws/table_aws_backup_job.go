package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/backup/types"

	backupv1 "github.com/aws/aws-sdk-go/service/backup"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBackupJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_job",
		Description: "AWS Backup Job",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("job_id"),
			Hydrate:    getAwsBackupJob,
			Tags:       map[string]string{"service": "backup", "action": "DescribeBackupJob"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupJobs,
			Tags:    map[string]string{"service": "backup", "action": "ListBackupJobs"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(backupv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "job_id",
				Description: "The logical id of a backup job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupJobId"),
			},
			{
				Name:        "recovery_point_arn",
				Description: "An Amazon Resource Name (ARN) that uniquely identifies a recovery point.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RecoveryPointArn"),
			},
			{
				Name:        "backup_vault_arn",
				Description: "An Amazon Resource Name (ARN) that uniquely identifies the target backup vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupVaultArn"),
			},
			{
				Name:        "resource_type",
				Description: "The type of AWS resource to be backed up.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceType"),
			},
			{
				Name:        "resource_arn",
				Description: "An Amazon Resource Name (ARN) that uniquely identifies the source resource in the recovery point.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceArn"),
			},
			{
				Name:        "status",
				Description: "The current state of a backup job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State"),
			},
			{
				Name:        "status_message",
				Description: "A detailed message explaining the status of the job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StatusMessage"),
			},
			{
				Name:        "backup_size",
				Description: "The size in bytes of a backup.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("BackupSizeInBytes"),
			},
			{
				Name:        "backup_vault_name",
				Description: "The name of the target backup vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupVaultName"),
			},
			{
				Name:        "bytes_transferred",
				Description: "The size in bytes transferred to a backup vault at the time that the job status was queried.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("BytesTransferred"),
			},
			{
				Name:        "completion_date",
				Description: "The date and time a backup job is completed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CompletionDate"),
			},
			{
				Name:        "creation_date",
				Description: "The date and time a backup job is created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationDate"),
			},
			{
				Name:        "expected_completion_date",
				Description: "The date and time a backup job is expected to be completed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ExpectedCompletionDate"),
			},
			{
				Name:        "iam_role_arn",
				Description: "The ARN of the IAM role that AWS Backup uses to create the target recovery point.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IamRoleArn"),
			},
			{
				Name:        "is_parent",
				Description: "A Boolean value that is returned as TRUE if the specified job is a parent job.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IsParentJob"),
			},
			{
				Name:        "parent_job_id",
				Description: "The ID of the parent backup job, if there is one.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ParentJobId"),
			},
			{
				Name:        "percent_done",
				Description: "The percentage of job completion.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PercentDone"),
			},
			{
				Name:        "start_by",
				Description: "The date and time a backup job must be started before it is canceled.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("StartBy"),
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupJobId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceArn").Transform(arnToAkas),
			},
		}),
	}
}

func listAwsBackupJobs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_vault.listAwsBackupJobs", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	queryResultLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		queryResultLimit = min(queryResultLimit, int32(*d.QueryContext.Limit))
	}

	input := &backup.ListBackupJobsInput{
		MaxResults: aws.Int32(queryResultLimit),
	}

	paginator := backup.NewListBackupJobsPaginator(svc, input, func(o *backup.ListBackupJobsPaginatorOptions) {
		o.Limit = queryResultLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_backup_job.listAwsBackupJobs", "api_error", err)
			return nil, err
		}

		for _, items := range output.BackupJobs {
			d.StreamListItem(ctx, items)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getAwsBackupJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_vault.getAwsBackupJob", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var backupJobID string
	if h.Item != nil {
		backupJob := h.Item.(types.BackupJob)
		backupJobID = *backupJob.BackupJobId
	} else {
		backupJobID = d.EqualsQualString("job_id")
	}

	params := &backup.DescribeBackupJobInput{
		BackupJobId: aws.String(backupJobID),
	}

	op, err := svc.DescribeBackupJob(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_vault.getAwsBackupJob", "api_error", err)
		return nil, err
	}

	return op, nil
}
