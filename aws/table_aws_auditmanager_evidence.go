package aws

import (
	"context"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/auditmanager"
	"github.com/aws/aws-sdk-go-v2/service/auditmanager/types"

	auditmanagerv1 "github.com/aws/aws-sdk-go/service/auditmanager"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type evidenceInfo struct {
	Evidence     types.Evidence
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
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidParameter"}),
			},
			Hydrate: getAuditManagerEvidence,
			Tags:    map[string]string{"service": "auditmanager", "action": "GetEvidence"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsAuditManagerAssessments,
			Hydrate:       listAuditManagerEvidences,
			Tags:          map[string]string{"service": "auditmanager", "action": "GetEvidenceByEvidenceFolder"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(auditmanagerv1.EndpointsID),
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

	// Get assessment details
	assessmentID := *h.Item.(types.AssessmentMetadataItem).Id

	// Create session
	svc, err := AuditManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_evidence.listAuditManagerEvidences", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var evidenceFolders []types.AssessmentEvidenceFolder
	maxItems := int32(100)

	// Get assessment details
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
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
			// for the regions where the  Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
			if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_auditmanager_evidence.listAuditManagerEvidences", "api_error", err)
			return nil, err
		}

		evidenceFolders = append(evidenceFolders, output.EvidenceFolders...)
	}

	var wg sync.WaitGroup
	evidenceCh := make(chan []evidenceInfo, len(evidenceFolders))
	errorCh := make(chan error, len(evidenceFolders))

	// Iterating all the available evidence folder
	for _, item := range evidenceFolders {
		wg.Add(1)
		go getRowDataForEvidenceAsync(ctx, svc, d, item, &wg, evidenceCh, errorCh)
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
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

func getRowDataForEvidenceAsync(ctx context.Context, svc *auditmanager.Client, d *plugin.QueryData, item types.AssessmentEvidenceFolder, wg *sync.WaitGroup, subnetCh chan []evidenceInfo, errorCh chan error) {
	defer wg.Done()

	rowData, err := getRowDataForEvidence(ctx, svc, d, item)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		subnetCh <- rowData
	}
}

func getRowDataForEvidence(ctx context.Context, svc *auditmanager.Client, d *plugin.QueryData, item types.AssessmentEvidenceFolder) ([]evidenceInfo, error) {

	params := &auditmanager.GetEvidenceByEvidenceFolderInput{
		AssessmentId:     item.AssessmentId,
		ControlSetId:     item.ControlSetId,
		EvidenceFolderId: item.Id,
	}

	var items []evidenceInfo

	listEvidence, err := svc.GetEvidenceByEvidenceFolder(ctx, params)

	// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
	// for the regions where the Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
	if err != nil {
		if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_auditmanager_evidence.getRowDataForEvidence", "api_error", err)
		return nil, err
	}

	for _, evidence := range listEvidence.Evidence {
		items = append(items, evidenceInfo{evidence, item.AssessmentId, item.ControlSetId})
	}

	return items, nil
}

//// HYDRATE FUNCTIONS

func getAuditManagerEvidence(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := AuditManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_evidence.getAuditManagerEvidence", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	assessmentID := d.EqualsQuals["assessment_id"].GetStringValue()
	controlSetID := d.EqualsQuals["control_set_id"].GetStringValue()
	evidenceFolderID := d.EqualsQuals["evidence_folder_id"].GetStringValue()
	evidenceID := d.EqualsQuals["id"].GetStringValue()

	// Build params
	params := &auditmanager.GetEvidenceInput{
		AssessmentId:     aws.String(assessmentID),
		ControlSetId:     aws.String(controlSetID),
		EvidenceFolderId: aws.String(evidenceFolderID),
		EvidenceId:       aws.String(evidenceID),
	}

	// Get call
	data, err := svc.GetEvidence(ctx, params)

	// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
	// for the regions where the Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
	if err != nil {
		if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_auditmanager_evidence.getAuditManagerEvidence", "api_error", err)
		return nil, err
	}

	return evidenceInfo{*data.Evidence, &assessmentID, &controlSetID}, nil
}

func getAuditManagerEvidenceARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAuditManagerEvidenceARN")
	region := d.EqualsQualString(matrixKeyRegion)
	evidenceID := *h.Item.(evidenceInfo).Evidence.Id

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_evidence.getAuditManagerEvidenceARN", "common_data_error", err)
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":auditmanager:" + region + ":" + commonColumnData.AccountId + ":evidence/" + evidenceID

	return arn, nil
}
