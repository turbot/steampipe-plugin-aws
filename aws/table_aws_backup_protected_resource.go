package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/backup/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBackupProtectedResource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_protected_resource",
		Description: "AWS Backup Protected Resource",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("resource_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidParameter", "InvalidParameterValueException"}),
			},
			Hydrate: getAwsBackupProtectedResource,
			Tags:    map[string]string{"service": "backup", "action": "DescribeProtectedResource"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupProtectedResources,
			Tags:    map[string]string{"service": "backup", "action": "ListProtectedResources"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsBackupProtectedResource,
				Tags: map[string]string{"service": "backup", "action": "DescribeProtectedResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_BACKUP_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "resource_name",
				Description: "This is the non-unique name of the resource that belongs to the specified backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_arn",
				Description: "An Amazon Resource Name (ARN) that uniquely identifies a resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_backup_vault_arn",
				Description: "This is the ARN (Amazon Resource Name) of the backup vault that contains the most recent backup recovery point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_recovery_point_arn",
				Description: "This is the ARN (Amazon Resource Name) of the most recent recovery point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The type of Amazon Web Services resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_backup_time",
				Description: "The date and time a resource was last backed up.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "latest_restore_execution_time_minutes",
				Description: "This is the time in minutes the most recent restore job took to complete.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsBackupProtectedResource,
			},
			{
				Name:        "latest_restore_job_creation_date",
				Description: "This is the creation date of the most recent restore job.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsBackupProtectedResource,
			},
			{
				Name:        "latest_restore_recovery_point_creation_date",
				Description: "This is the date the most recent recovery point was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsBackupProtectedResource,
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

func listAwsBackupProtectedResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_protected_resource.listAwsBackupProtectedResources", "connection_error", err)
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

	input := &backup.ListProtectedResourcesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := backup.NewListProtectedResourcesPaginator(svc, input, func(o *backup.ListProtectedResourcesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_backup_protected_resource.listAwsBackupProtectedResources", "api_error", err)
			return nil, err
		}

		for _, items := range output.Results {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getAwsBackupProtectedResource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_protected_resource.getAwsBackupProtectedResource", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.ProtectedResource).ResourceArn
	} else {
		arn = d.EqualsQuals["resource_arn"].GetStringValue()
	}

	params := &backup.DescribeProtectedResourceInput{
		ResourceArn: aws.String(arn),
	}

	detail, err := svc.DescribeProtectedResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_protected_resource.getAwsBackupProtectedResource", "api_error", err)
		return nil, err
	}

	return detail, nil
}
