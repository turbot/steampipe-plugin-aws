package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/auditmanager"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAuditManagerEvidenceFolder(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_auditmanager_evidence_folder",
		Description: "AWS Audit Manager Evidence Folder",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"id", "assessment_id", "control_set_id"}),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException", "InvalidParameter"}),
			Hydrate:           getAuditManagerEvidenceFolder,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsAuditManagerAssessments,
			Hydrate:       listAuditManagerEvidenceFolders,
		},
		GetMatrixItem: BuildRegionList,
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listAuditManagerEvidenceFolders", "AWS_REGION", region)

	// Create session
	svc, err := AuditManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Get assessment details
	assessmentID := *h.Item.(*auditmanager.AssessmentMetadataItem).Id

	input := &auditmanager.GetEvidenceFoldersByAssessmentInput{
		MaxResults: aws.Int64(1000),
	}

	input.AssessmentId = &assessmentID
	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 5 {
				input.MaxResults = types.Int64(5)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.GetEvidenceFoldersByAssessmentPages(
		input,
		func(page *auditmanager.GetEvidenceFoldersByAssessmentOutput, isLast bool) bool {
			for _, folder := range page.EvidenceFolders {
				d.StreamListItem(ctx, folder)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getAuditManagerEvidenceFolder(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAuditManagerEvidenceFolder")

	region := d.KeyColumnQualString(matrixKeyRegion)

	// Create Session
	svc, err := AuditManagerService(ctx, d, region)
	if err != nil {
		return nil, err
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
	data, err := svc.GetEvidenceFolder(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAuditManagerEvidenceFolder", "ERROR", err)
		return nil, err
	}

	return data.EvidenceFolder, nil
}

func getAuditManagerEvidenceFolderARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAuditManagerEvidenceFolderARN")
	region := d.KeyColumnQualString(matrixKeyRegion)
	evidenceFolderID := *h.Item.(*auditmanager.AssessmentEvidenceFolder).Id

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":auditmanager:" + region + ":" + commonColumnData.AccountId + ":evidence-folder/" + evidenceFolderID

	return arn, nil
}
