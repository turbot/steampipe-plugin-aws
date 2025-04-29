package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/fsx"
	"github.com/aws/aws-sdk-go-v2/service/fsx/types"

	fsxv1 "github.com/aws/aws-sdk-go/service/fsx"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsFsxBackup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_fsx_backup",
		Description: "AWS FSx Backup",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("backup_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"BackupNotFound"}),
			},
			Hydrate: getFsxBackup,
			Tags:    map[string]string{"service": "fsx", "action": "DescribeBackups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listFsxBackups,
			Tags:    map[string]string{"service": "fsx", "action": "DescribeBackups"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "file_system_id", Require: plugin.Optional},
				{Name: "backup_type", Require: plugin.Optional},
				{Name: "lifecycle", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(fsxv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "backup_id",
				Description: "The ID of the backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the backup resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceARN"),
			},
			{
				Name:        "creation_time",
				Description: "The time when a particular backup was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "file_system_id",
				Description: "The ID of the file system from which the backup was taken.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FileSystem.FileSystemId"),
			},
			{
				Name:        "lifecycle",
				Description: "The lifecycle status of the backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_type",
				Description: "The type of the backup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Type"),
			},
			{
				Name:        "directory_information",
				Description: "The configuration of the self-managed Microsoft Active Directory (AD) to which the Windows File Server instance is joined.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "failure_details",
				Description: "Details explaining any failures that occurred when creating a backup.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "kms_key_id",
				Description: "The ID of the AWS Key Management Service (AWS KMS) key used to encrypt the backup of the Amazon FSx file system's data at rest.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The AWS account ID that owns the backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "progress_percent",
				Description: "The current percent of progress of an asynchronous task.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "resource_arn",
				Description: "The Amazon Resource Name (ARN) for the backup resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_backup_id",
				Description: "The ID of the source backup. Specifies the backup that you are copying.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_backup_region",
				Description: "The source Region of the backup. Specifies the Region from where this backup is copied.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the backup.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: "The tags associated with the backup.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getFsxBackupTags),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listFsxBackups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := FSxClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fsx_backup.listFsxBackups", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(1000)
	input := &fsx.DescribeBackupsInput{
		MaxResults: &maxItems,
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["file_system_id"] != nil {
		input.Filters = append(input.Filters, types.Filter{
			Name:   "file-system-id",
			Values: []string{equalQuals["file_system_id"].GetStringValue()},
		})
	}
	if equalQuals["backup_type"] != nil {
		input.Filters = append(input.Filters, types.Filter{
			Name:   "backup-type",
			Values: []string{equalQuals["backup_type"].GetStringValue()},
		})
	}
	if equalQuals["lifecycle"] != nil {
		input.Filters = append(input.Filters, types.Filter{
			Name:   "lifecycle",
			Values: []string{equalQuals["lifecycle"].GetStringValue()},
		})
	}

	paginator := fsx.NewDescribeBackupsPaginator(svc, input, func(o *fsx.DescribeBackupsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_fsx_backup.listFsxBackups", "api_error", err)
			return nil, err
		}

		for _, backup := range output.Backups {
			d.StreamListItem(ctx, backup)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getFsxBackup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	backupId := d.EqualsQuals["backup_id"].GetStringValue()

	// Create session
	svc, err := FSxClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fsx_backup.getFsxBackup", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &fsx.DescribeBackupsInput{
		BackupIds: []string{backupId},
	}

	op, err := svc.DescribeBackups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fsx_backup.getFsxBackup", "api_error", err)
		return nil, err
	}

	if op.Backups != nil && len(op.Backups) > 0 {
		return op.Backups[0], nil
	}

	return nil, nil
}

func getFsxBackupTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	backup := d.HydrateItem.(types.Backup)

	// Convert tags to a map
	tagsMap := make(map[string]string)
	for _, tag := range backup.Tags {
		if tag.Key != nil && tag.Value != nil {
			tagsMap[*tag.Key] = *tag.Value
		}
	}

	return tagsMap, nil
}
