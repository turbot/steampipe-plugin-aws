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

func tableAwsBackupPlan(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_plan",
		Description: "AWS Backup Plan",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"backup_plan_id", "version_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValueException"}),
			},
			Hydrate: getAwsBackupPlan,
			Tags:    map[string]string{"service": "backup", "action": "GetBackupPlan"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupPlans,
			Tags:    map[string]string{"service": "backup", "action": "ListBackupPlans"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsBackupPlan,
				Tags: map[string]string{"service": "backup", "action": "GetBackupPlan"},
			},
			{
				Func: getAwsBackupPlanTags,
				Tags: map[string]string{"service": "backup", "action": "ListTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(backupv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The display name of a saved backup plan.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupPlanName", "BackupPlan.BackupPlanName"),
			},
			{
				Name:        "arn",
				Description: "An Amazon Resource Name (ARN) that uniquely identifies a backup plan.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupPlanArn"),
			},
			{
				Name:        "backup_plan_id",
				Description: "Specifies the id to identify a backup plan uniquely.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "The date and time a resource backup plan is created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "deletion_date",
				Description: "The date and time a backup plan is deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_execution_date",
				Description: "The last time a job to back up resources was run with this rule.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "creator_request_id",
				Description: "An unique string that identifies the request and allows failed requests to be retried without the risk of running the operation twice.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_id",
				Description: "Unique, randomly generated, Unicode, UTF-8 encoded strings that are at most 1,024 bytes long. Version IDs cannot be edited.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_plan",
				Description: "Specifies the body of a backup plan.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsBackupPlan,
			},
			{
				Name:        "advanced_backup_settings",
				Description: "Contains a list of BackupOptions for a resource type.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupPlanName", "BackupPlan.BackupPlanName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsBackupPlanTags,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackupPlanArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsBackupPlans(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_plan.listAwsBackupPlans", "connection_error", err)
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

	input := &backup.ListBackupPlansInput{
		MaxResults:     aws.Int32(maxLimit),
		IncludeDeleted: aws.Bool(true),
	}

	paginator := backup.NewListBackupPlansPaginator(svc, input, func(o *backup.ListBackupPlansPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_backup_plan.listAwsBackupPlans", "api_error", err)
			return nil, err
		}

		for _, items := range output.BackupPlansList {
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

func getAwsBackupPlan(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_plan.getAwsBackupPlan", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var id, versionId string
	if h.Item != nil {
		plan := h.Item.(types.BackupPlansListMember)
		id = *plan.BackupPlanId
		versionId = *plan.VersionId
	} else {
		id = d.EqualsQuals["backup_plan_id"].GetStringValue()
		versionId = d.EqualsQuals["version_id"].GetStringValue()
	}

	// check if id is empty
	if id == "" || versionId == "" {
		return nil, nil
	}

	params := &backup.GetBackupPlanInput{
		BackupPlanId: aws.String(id),
		VersionId:    aws.String(versionId),
	}

	op, err := svc.GetBackupPlan(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_plan.getAwsBackupPlan", "api_error", err)
	}

	return op, nil
}

func getAwsBackupPlanTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := backupPlanArn(h.Item)

	// Define the regex pattern for the recovery point ARN
	pattern := `arn:aws:backup:[a-z0-9\-]+:[0-9]{12}:backup-plan:.*`

	return getAwsBackupResourceTags(ctx, d, arn, pattern)
}

func backupPlanArn(item interface{}) string {
	switch item := item.(type) {
	case types.BackupPlansListMember:
		return *item.BackupPlanArn
	case backup.GetBackupPlanOutput:
		return *item.BackupPlanArn
	}
	return ""
}
