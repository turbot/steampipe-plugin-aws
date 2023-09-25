package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"

	backupv1 "github.com/aws/aws-sdk-go/service/backup"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBackupReportPlan(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_report_plan",
		Description: "AWS Backup Report Plan",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("report_plan_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValueException", "ResourceNotFoundException"}),
			},
			Hydrate: getAwsBackupReportPlan,
			Tags:    map[string]string{"service": "backup", "action": "DescribeReportPlan"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupReportPlans,
			Tags:    map[string]string{"service": "backup", "action": "ListReportPlans"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(backupv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "An Amazon Resource Name (ARN) that uniquely identifies a resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReportPlanArn"),
			},
			{
				Name:        "report_plan_name",
				Description: "The unique name of the report plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "An optional description of the report plan with a maximum 1,024 characters.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReportPlanDescription"),
			},
			{
				Name:        "creation_time",
				Description: "The date and time that a report plan is created, in Unix format and Coordinated Universal Time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "deployment_status",
				Description: "The deployment status of a report plan. The statuses are CREATE_IN_PROGRESS, UPDATE_IN_PROGRESS, DELETE_IN_PROGRESS, and COMPLETED.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_attempted_execution_time",
				Description: "The date and time that a report job associated with this report plan last attempted to run, in Unix format and Coordinated Universal Time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_successful_execution_time",
				Description: "The date and time that a report job associated with this report plan last successfully ran, in Unix format and Coordinated Universal Time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// JSON Columns
			{
				Name:        "report_delivery_channel",
				Description: "Contains information about where and how to deliver your reports, specifically your Amazon S3 bucket name, S3 key prefix, and the formats of your reports.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "report_setting",
				Description: "Identifies the report template for the report. Reports are built using a report template.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReportPlanName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ReportPlanArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsBackupReportPlans(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_report_plan.listAwsBackupReportPlans", "connection_error", err)
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
			maxLimit = limit
		}
	}

	input := &backup.ListReportPlansInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := backup.NewListReportPlansPaginator(svc, input, func(o *backup.ListReportPlansPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_backup_report_plan.listAwsBackupReportPlans", "api_error", err)
			return nil, err
		}

		for _, items := range output.ReportPlans {
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

func getAwsBackupReportPlan(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_report_plan.getAwsBackupReportPlan", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	reportPlanName := d.EqualsQualString("report_plan_name")

	// check if name is empty
	if reportPlanName == "" {
		return nil, nil
	}

	params := &backup.DescribeReportPlanInput{
		ReportPlanName: aws.String(reportPlanName),
	}

	op, err := svc.DescribeReportPlan(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_report_plan.getAwsBackupReportPlan", "api_error", err)
	}

	if op != nil {
		return *op.ReportPlan, nil
	}

	return nil, nil
}
