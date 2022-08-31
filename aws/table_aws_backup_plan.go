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

func tableAwsBackupPlan(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_plan",
		Description: "AWS Backup Plan",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("backup_plan_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidParameterValueException"}),
			},
			Hydrate: getAwsBackupPlan,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupPlans,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &backup.ListBackupPlansInput{
		MaxResults: aws.Int64(1000),
	}
	input.IncludeDeleted = aws.Bool(true)

	// Limiting the results
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

	err = svc.ListBackupPlansPages(
		input,
		func(page *backup.ListBackupPlansOutput, lastPage bool) bool {
			for _, plan := range page.BackupPlansList {
				d.StreamListItem(ctx, plan)

				// Context can be cancelled due to manual cancellation or the limit has been hit
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

func getAwsBackupPlan(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}

	var id string
	if h.Item != nil {
		plan := h.Item.(*backup.PlansListMember)
		id = *plan.BackupPlanId
	} else {
		id = d.KeyColumnQuals["backup_plan_id"].GetStringValue()
	}

	// check if id is empty
	if id == "" {
		return nil, nil
	}

	params := &backup.GetBackupPlanInput{
		BackupPlanId: aws.String(id),
	}

	op, err := svc.GetBackupPlan(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		return nil, err
	}
	return op, nil
}
