package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/auditmanager"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAuditManagerAssessment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_auditmanager_assessment",
		Description: "AWS Audit Manager Assessment",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException", "ValidationException", "InvalidParameter"}),
			Hydrate:           getAwsAuditManagerAssessment,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsAuditManagerAssessments,
		},
		GetMatrixItem: BuildRegionList,
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listAwsAuditManagerAssessments", "AWS_REGION", region)

	// Create session
	svc, err := AuditManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	input := &auditmanager.ListAssessmentsInput{
		MaxResults: aws.Int64(1000),
	}

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
	err = svc.ListAssessmentsPages(
		input,
		func(page *auditmanager.ListAssessmentsOutput, isLast bool) bool {
			for _, assessment := range page.AssessmentMetadata {
				d.StreamListItem(ctx, assessment)

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

func getAwsAuditManagerAssessment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsAuditManagerAssessment")

	region := d.KeyColumnQualString(matrixKeyRegion)

	var id string
	if h.Item != nil {
		id = *h.Item.(*auditmanager.AssessmentMetadataItem).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// Create Session
	svc, err := AuditManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &auditmanager.GetAssessmentInput{
		AssessmentId: aws.String(id),
	}

	// Get call
	data, err := svc.GetAssessment(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAwsAuditManagerAssessment", "ERROR", err)
		return nil, err
	}

	return data.Assessment, nil
}
