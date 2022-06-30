package aws

import (
	"context"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/auditmanager"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

type evidenceInfo struct {
	Evidence     *auditmanager.Evidence
	AssessmentID *string
	ControlSetID *string
}

//// TABLE DEFINITION

func tableAwsAuditManagerEvidence(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_auditmanager_evidence",
		Description: "AWS Audit Manager Evidence",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "evidence_folder_id", "assessment_id", "control_set_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException", "InvalidParameter"}),
			},
			Hydrate: getAuditManagerEvidence,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsAuditManagerAssessments,
			Hydrate:       listAuditManagerEvidences,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The identifier for the evidence.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.Id"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the evidence.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerEvidenceARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "assessment_id",
				Description: "An unique identifier for the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentID"),
			},
			{
				Name:        "control_set_id",
				Description: "The identifier for the control set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ControlSetID"),
			},
			{
				Name:        "evidence_folder_id",
				Description: "The identifier for the folder in which the evidence is stored.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.EvidenceFolderId"),
			},
			{
				Name:        "assessment_report_selection",
				Description: "Specifies whether the evidence is included in the assessment report.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.AssessmentReportSelection"),
			},
			{
				Name:        "aws_account_id",
				Description: "The identifier for the specified AWS account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.AwsAccountId"),
			},
			{
				Name:        "aws_organization",
				Description: "The AWS account from which the evidence is collected, and its AWS organization path.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.AwsOrganization"),
			},
			{
				Name:        "compliance_check",
				Description: "The evaluation status for evidence that falls under the compliance check category.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.ComplianceCheck"),
			},
			{
				Name:        "data_source",
				Description: "The data source from which the specified evidence was collected.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.DataSource"),
			},
			{
				Name:        "event_name",
				Description: "The name of the specified evidence event.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.EventName"),
			},
			{
				Name:        "event_source",
				Description: "The AWS service from which the evidence is collected.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.EventSource"),
			},
			{
				Name:        "evidence_aws_account_id",
				Description: "The identifier for the specified AWS account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.EvidenceAwsAccountId"),
			},
			{
				Name:        "evidence_by_type",
				Description: "The type of automated evidence.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.EvidenceByType"),
			},
			{
				Name:        "iam_id",
				Description: "The unique identifier for the IAM user or role associated with the evidence.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.IamId"),
			},
			{
				Name:        "time",
				Description: "The timestamp that represents when the evidence was collected.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Evidence.Time"),
			},
			{
				Name:        "attributes",
				Description: "The names and values used by the evidence event",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Evidence.Attributes"),
			},
			{
				Name:        "resources_included",
				Description: "The list of resources assessed to generate the evidence.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Evidence.ResourcesIncluded"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Evidence.Id"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAuditManagerEvidenceARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAuditManagerEvidences(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listAuditManagerEvidences", "AWS_REGION", region)

	// AWS Audit Manager is not supported in all regions. For unsupported regions the API throws an error, e.g.,
	// Get "https://auditmanager.sa-east-1.amazonaws.com/": dial tcp: lookup auditmanager.sa-east-1.amazonaws.com: no such host
	serviceId := auditmanager.EndpointsID
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	// Get assessment details
	assessmentID := *h.Item.(*auditmanager.AssessmentMetadataItem).Id

	// Create session
	svc, err := AuditManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	var evidenceFolders []auditmanager.AssessmentEvidenceFolder
	input := &auditmanager.GetEvidenceFoldersByAssessmentInput{
		MaxResults: aws.Int64(1000),
	}
	input.AssessmentId = aws.String(assessmentID)

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

	// List call
	err = svc.GetEvidenceFoldersByAssessmentPages(
		input,
		func(page *auditmanager.GetEvidenceFoldersByAssessmentOutput, isLast bool) bool {
			for _, evidenceFolder := range page.EvidenceFolders {
				evidenceFolders = append(evidenceFolders, *evidenceFolder)
			}
			return !isLast
		},
	)

	var wg sync.WaitGroup
	evidenceCh := make(chan []evidenceInfo, len(evidenceFolders))
	errorCh := make(chan error, len(evidenceFolders))

	// Iterating all the available evidence folder
	for _, item := range evidenceFolders {
		wg.Add(1)
		go getRowDataForEvidenceAsync(ctx, d, item, &wg, evidenceCh, errorCh, region)
	}

	// wait for all evidence folder to be processed
	wg.Wait()
	close(evidenceCh)
	close(errorCh)

	for err := range errorCh {
		return nil, err
	}

	for item := range evidenceCh {
		for _, data := range item {
			d.StreamLeafListItem(ctx, evidenceInfo{data.Evidence, data.AssessmentID, data.ControlSetID})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

func getRowDataForEvidenceAsync(ctx context.Context, d *plugin.QueryData, item auditmanager.AssessmentEvidenceFolder, wg *sync.WaitGroup, subnetCh chan []evidenceInfo, errorCh chan error, region string) {
	defer wg.Done()

	rowData, err := getRowDataForEvidence(ctx, d, item, region)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		subnetCh <- rowData
	}
}

func getRowDataForEvidence(ctx context.Context, d *plugin.QueryData, item auditmanager.AssessmentEvidenceFolder, region string) ([]evidenceInfo, error) {
	svc, err := AuditManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &auditmanager.GetEvidenceByEvidenceFolderInput{
		AssessmentId:     item.AssessmentId,
		ControlSetId:     item.ControlSetId,
		EvidenceFolderId: item.Id,
	}

	var items []evidenceInfo

	listEvidence, err := svc.GetEvidenceByEvidenceFolder(params)

	// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
	// for the regions where the Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
	if err != nil {
		if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("getRowDataForEvidence", "err", err)
		return nil, err
	}

	for _, evidence := range listEvidence.Evidence {
		items = append(items, evidenceInfo{evidence, item.AssessmentId, item.ControlSetId})
	}

	return items, nil
}

//// HYDRATE FUNCTIONS

func getAuditManagerEvidence(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAuditManagerEvidence")
	region := d.KeyColumnQualString(matrixKeyRegion)

	// AWS Audit Manager is not supported in all regions. For unsupported regions the API throws an error, e.g.,
	// Get "https://auditmanager.sa-east-1.amazonaws.com/": dial tcp: lookup auditmanager.sa-east-1.amazonaws.com: no such host
	serviceId := auditmanager.EndpointsID
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	// Create Session
	svc, err := AuditManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	assessmentID := d.KeyColumnQuals["assessment_id"].GetStringValue()
	controlSetID := d.KeyColumnQuals["control_set_id"].GetStringValue()
	evidenceFolderID := d.KeyColumnQuals["evidence_folder_id"].GetStringValue()
	evidenceID := d.KeyColumnQuals["id"].GetStringValue()

	// Build params
	params := &auditmanager.GetEvidenceInput{
		AssessmentId:     aws.String(assessmentID),
		ControlSetId:     aws.String(controlSetID),
		EvidenceFolderId: aws.String(evidenceFolderID),
		EvidenceId:       aws.String(evidenceID),
	}

	// Get call
	data, err := svc.GetEvidence(params)
	if err != nil {
		plugin.Logger(ctx).Error("getAuditManagerEvidence", "ERROR", err)
		return nil, err
	}

	return evidenceInfo{data.Evidence, &assessmentID, &controlSetID}, nil
}

func getAuditManagerEvidenceARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAuditManagerEvidenceARN")
	region := d.KeyColumnQualString(matrixKeyRegion)
	evidenceID := *h.Item.(evidenceInfo).Evidence.Id

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":auditmanager:" + region + ":" + commonColumnData.AccountId + ":evidence/" + evidenceID

	return arn, nil
}
