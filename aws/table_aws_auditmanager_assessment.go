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

func tableAwsAuditManagerAssessment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_auditmanager_assessment",
		Description: "AWS Audit Manager Assessment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException", "ValidationException", "InvalidParameter"}),
			},
			Hydrate: getAwsAuditManagerAssessment,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsAuditManagerAssessments,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "Metadata.Name"),
			},
			{
				Name:        "id",
				Description: "An unique identifier for the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id", "Metadata.Id"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the assessment.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAuditManagerAssessment,
			},
			{
				Name:        "status",
				Description: "The current status of the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status", "Metadata.Status"),
			},
			{
				Name:        "compliance_type",
				Description: "The name of the compliance standard related to the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ComplianceType", "Metadata.ComplianceType"),
			},
			{
				Name:        "assessment_report_destination",
				Description: "The destination of the assessment report.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAuditManagerAssessment,
				Transform:   transform.FromField("Metadata.AssessmentReportsDestination.Destination"),
			},
			{
				Name:        "assessment_report_destination_type",
				Description: "The destination type, such as Amazon S3.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAuditManagerAssessment,
				Transform:   transform.FromField("Metadata.AssessmentReportsDestination.DestinationType"),
			},
			{
				Name:        "creation_time",
				Description: "Specifies when the assessment was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationTime", "Metadata.CreationTime"),
			},
			{
				Name:        "description",
				Description: "The description of the assessment.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAuditManagerAssessment,
				Transform:   transform.FromField("Metadata.Description"),
			},
			{
				Name:        "last_updated",
				Description: "The time of the most recent update.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastUpdated", "Metadata.LastUpdated"),
			},
			{
				Name:        "aws_account",
				Description: "The AWS account associated with the assessment.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAuditManagerAssessment,
			},
			{
				Name:        "delegations",
				Description: "The delegations associated with the assessment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Delegations", "Metadata.Delegations"),
			},
			{
				Name:        "framework",
				Description: "The framework from which the assessment was created.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAuditManagerAssessment,
			},
			{
				Name:        "scope",
				Description: "The wrapper of AWS accounts and services in scope for the assessment.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAuditManagerAssessment,
				Transform:   transform.FromField("Metadata.Scope"),
			},
			{
				Name:        "roles",
				Description: "The roles associated with the assessment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Roles", "Metadata.Roles"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "Metadata.Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAuditManagerAssessment,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAuditManagerAssessment,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsAuditManagerAssessments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := AuditManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_assessment.listAwsAuditManagerAssessments", "client_error", err)
		return nil, err
	}

	maxItems := int32(1000)
	input := &auditmanager.ListAssessmentsInput{}

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

	paginator := auditmanager.NewListAssessmentsPaginator(svc, input, func(o *auditmanager.ListAssessmentsPaginatorOptions) {
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
			plugin.Logger(ctx).Error("aws_auditmanager_assessment.listAwsAuditManagerAssessments", "api_error", err)
			return nil, err
		}

		for _, assessment := range output.AssessmentMetadata {
			d.StreamListItem(ctx, assessment)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsAuditManagerAssessment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		id = *h.Item.(types.AssessmentMetadataItem).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Get client
	svc, err := AuditManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_assessment.getAwsAuditManagerAssessment", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &auditmanager.GetAssessmentInput{
		AssessmentId: aws.String(id),
	}

	// Get call
	data, err := svc.GetAssessment(ctx, params)

	// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
	// for the regions where the  Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
	if err != nil {
		if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_auditmanager_assessment.getAwsAuditManagerAssessment", "api_error", err)
		return nil, err
	}

	return data.Assessment, nil
}
