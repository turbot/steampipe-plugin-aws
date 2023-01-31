package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/backup/types"

	backupv1 "github.com/aws/aws-sdk-go/service/backup"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBackupSelection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_selection",
		Description: "AWS Backup Selection",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"backup_plan_id", "selection_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue", "InvalidParameterValueException"}),
			},
			Hydrate: getBackupSelection,
		},
		List: &plugin.ListConfig{
			Hydrate:       listBackupSelections,
			ParentHydrate: listAwsBackupPlans,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(backupv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "selection_name",
				Description: "The display name of a resource selection document.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SelectionName", "BackupSelection.SelectionName"),
			},
			{
				Name:        "selection_id",
				Description: "Uniquely identifies a request to assign a set of resources to a backup plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_plan_id",
				Description: "An ID that uniquely identifies a backup plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the backup selection.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBackupSelectionARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "creation_date",
				Description: "The date and time a resource backup plan is created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "creator_request_id",
				Description: "An unique string that identifies the request and allows failed requests to be retried without the risk of running the operation twice.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_role_arn",
				Description: "Specifies the IAM role Amazon Resource Name (ARN) to create the target recovery point.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IamRoleArn", "BackupSelection.IamRoleArn"),
			},
			{
				Name:        "list_of_tags",
				Description: "An array of conditions used to specify a set of resources to assign to a backup plan.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBackupSelection,
				Transform:   transform.FromField("BackupSelection.ListOfTags"),
			},
			{
				Name:        "resources",
				Description: "Contains a list of BackupOptions for a resource type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBackupSelection,
				Transform:   transform.FromField("BackupSelection.Resources"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SelectionName", "BackupSelection.SelectionName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBackupSelectionARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listBackupSelections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_selection.listBackupSelections", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Get backup plan details
	plan := h.Item.(types.BackupPlansListMember)

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

	input := &backup.ListBackupSelectionsInput{
		MaxResults:   aws.Int32(maxLimit),
		BackupPlanId: aws.String(*plan.BackupPlanId),
	}

	paginator := backup.NewListBackupSelectionsPaginator(svc, input, func(o *backup.ListBackupSelectionsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_backup_selection.listBackupSelections", "api_error", err)
			return nil, err
		}

		for _, items := range output.BackupSelectionsList {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getBackupSelection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_selection.getBackupSelection", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var backupPlanID, selectionID string
	if h.Item != nil {
		backupPlanID = *h.Item.(types.BackupSelectionsListMember).BackupPlanId
		selectionID = *h.Item.(types.BackupSelectionsListMember).SelectionId
	} else {
		backupPlanID = d.KeyColumnQuals["backup_plan_id"].GetStringValue()
		selectionID = d.KeyColumnQuals["selection_id"].GetStringValue()
	}

	// Return nil, if no input provided
	if backupPlanID == "" || selectionID == "" {
		return nil, nil
	}

	params := &backup.GetBackupSelectionInput{
		BackupPlanId: aws.String(backupPlanID),
		SelectionId:  aws.String(selectionID),
	}

	op, err := svc.GetBackupSelection(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_selection.getBackupSelection", "api_error", err)

		// API returns error if Backup Plan is not available in a region
		if strings.Contains(err.Error(), "Cannot find Backup plan with ID") {
			return nil, nil
		}
		return nil, err
	}
	return op, nil
}

func getBackupSelectionARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := selectionID(h.Item)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Build ARN
	arn := "arn:" + commonColumnData.Partition + ":backup:" + region + ":" + commonColumnData.AccountId + ":backup-plan:" + data["PlanID"] + "/selection/" + data["SelectionID"]

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func selectionID(item interface{}) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case types.BackupSelectionsListMember:
		data["PlanID"] = *item.BackupPlanId
		data["SelectionID"] = *item.SelectionId
	case *backup.GetBackupSelectionOutput:
		data["PlanID"] = *item.BackupPlanId
		data["SelectionID"] = *item.SelectionId
	}
	return data
}
