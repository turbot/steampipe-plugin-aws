package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/auditmanager"
	"github.com/aws/aws-sdk-go-v2/service/auditmanager/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAuditManagerEvidenceFolder(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_auditmanager_evidence_folder",
		Description: "AWS Audit Manager Evidence Folder",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "assessment_id", "control_set_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidParameter"}),
			},
			Hydrate: getAuditManagerEvidenceFolder,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsAuditManagerAssessments,
			Hydrate:       listAuditManagerEvidenceFolders,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the specified evidence folder.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The identifier for the folder in which evidence is stored.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the evidence folder.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerEvidenceFolderARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "assessment_id",
				Description: "The identifier for the specified assessment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "control_set_id",
				Description: "The identifier for the control set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "assessment_report_selection_count",
				Description: "The total count of evidence included in the assessment report.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "author",
				Description: "The name of the user who created the evidence folder.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "control_id",
				Description: "The unique identifier for the specified control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "control_name",
				Description: "The name of the control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_source",
				Description: "The AWS service from which the evidence was collected.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "date",
				Description: "The date when the first evidence was added to the evidence folder.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "evidence_aws_service_source_count",
				Description: "The total number of AWS resources assessed to generate the evidence.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "evidence_by_type_compliance_check_count",
				Description: "The number of evidence that falls under the compliance check category.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "evidence_by_type_compliance_check_issues_count",
				Description: "The total number of issues that were reported directly from AWS Security Hub, AWS Config, or both.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "evidence_by_type_configuration_data_count",
				Description: "The number of evidence that falls under the configuration data category.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "evidence_by_type_manual_count",
				Description: "The number of evidence that falls under the manual category.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "evidence_by_type_user_activity_count",
				Description: "The number of evidence that falls under the user activity category.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "evidence_resources_included_count",
				Description: "The amount of evidence included in the evidence folder.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "total_evidence",
				Description: "The total amount of evidence in the evidence folder.",
				Type:        proto.ColumnType_INT,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAuditManagerEvidenceFolderARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAuditManagerEvidenceFolders(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	svc, err := AuditManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_evidence_folder.listAuditManagerEvidenceFolders", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(100)
	// Get assessment details
	assessmentID := *h.Item.(types.AssessmentMetadataItem).Id
	params := &auditmanager.GetEvidenceFoldersByAssessmentInput{
		AssessmentId: &assessmentID,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	params.MaxResults = &maxItems

	paginator := auditmanager.NewGetEvidenceFoldersByAssessmentPaginator(svc, params, func(o *auditmanager.GetEvidenceFoldersByAssessmentPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
			// for the regions where the  Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
			if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_auditmanager_evidence_folder.listAuditManagerEvidenceFolders", "api_error", err)
			return nil, err
		}

		for _, folder := range output.EvidenceFolders {
			d.StreamListItem(ctx, folder)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAuditManagerEvidenceFolder(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := AuditManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_evidence_folder.getAuditManagerEvidenceFolder", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	assessmentID := d.KeyColumnQuals["assessment_id"].GetStringValue()
	controlSetID := d.KeyColumnQuals["control_set_id"].GetStringValue()
	evidenceFolderID := d.KeyColumnQuals["id"].GetStringValue()

	// Build params
	params := &auditmanager.GetEvidenceFolderInput{
		AssessmentId:     aws.String(assessmentID),
		ControlSetId:     aws.String(controlSetID),
		EvidenceFolderId: aws.String(evidenceFolderID),
	}

	// Get call
	data, err := svc.GetEvidenceFolder(ctx, params)

	// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
	// for the regions where the Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
	if err != nil {
		if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_auditmanager_evidence_folder.getAuditManagerEvidenceFolder", "api_error", err)
		return nil, err
	}

	return *data.EvidenceFolder, nil
}

func getAuditManagerEvidenceFolderARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	evidenceFolderID := *h.Item.(types.AssessmentEvidenceFolder).Id

	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_evidence_folder.getAuditManagerEvidenceFolderARN", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":auditmanager:" + region + ":" + commonColumnData.AccountId + ":evidence-folder/" + evidenceFolderID

	return arn, nil
}
