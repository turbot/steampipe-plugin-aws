package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/auditmanager"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAuditManagerFramework(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_auditmanager_framework",
		Description: "AWS Audit Manager Framework",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "region"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException", "ValidationException", "InternalServerException"}),
			},
			Hydrate: getAuditManagerFramework,
		},
		List: &plugin.ListConfig{
			Hydrate: listAuditManagerFrameworks,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the specified framework.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the framework.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identified for the specified framework.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The framework type, such as standard or custom.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "Specifies when the framework was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by",
				Description: "The IAM user or role that created the framework.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerFramework,
			},
			{
				Name:        "compliance_type",
				Description: "The compliance type that the new custom framework supports, such as CIS or HIPAA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "controls_count",
				Description: "The number of controls associated with the specified framework.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "control_sets_count",
				Description: "The number of control sets associated with the specified framework.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "control_sources",
				Description: "The sources from which AWS Audit Manager collects evidence for the control.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerFramework,
			},
			{
				Name:        "description",
				Description: "The description of the specified framework.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_at",
				Description: "Specifies when the framework was most recently updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_updated_by",
				Description: "The IAM user or role that most recently updated the framework.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerFramework,
			},
			{
				Name:        "logo",
				Description: "The logo associated with the framework.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "control_sets",
				Description: "The control sets associated with the framework.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAuditManagerFramework,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAuditManagerFramework,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAuditManagerFrameworks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Debug("listAuditManagerFrameworks", "REGION", region)

	// AWS Audit Manager is not supported in all regions. For unsupported regions the API throws an error, e.g.,
	// Get "https://auditmanager.ap-northeast-3.amazonaws.com/assessmentFrameworks?frameworkType=Standard": dial tcp: lookup auditmanager.ap-northeast-3.amazonaws.com: no such host
	serviceId := auditmanager.EndpointsID
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	svc, err := AuditManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List standard audit manager frameworks
	err = svc.ListAssessmentFrameworksPages(
		&auditmanager.ListAssessmentFrameworksInput{FrameworkType: aws.String("Standard")},
		func(page *auditmanager.ListAssessmentFrameworksOutput, lastPage bool) bool {
			for _, framework := range page.FrameworkMetadataList {
				d.StreamListItem(ctx, framework)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
	// for the regions where the Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
	if err != nil {
		if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("listAuditManagerFrameworks_standard", "err", err)
		return nil, err
	}

	// List custom audit manager frameworks
	err = svc.ListAssessmentFrameworksPages(
		&auditmanager.ListAssessmentFrameworksInput{FrameworkType: aws.String("Custom")},
		func(page *auditmanager.ListAssessmentFrameworksOutput, lastPage bool) bool {
			for _, framework := range page.FrameworkMetadataList {
				d.StreamListItem(ctx, framework)
			}
			return !lastPage
		},
	)

	// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
	// for the regions where the Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
	if err != nil {
		if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("listAuditManagerFrameworks_custom", "err", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAuditManagerFramework(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Debug("getAuditManagerFramework", "REGION", region)

	// AWS Audit Manager is not supported in all regions. For unsupported regions the API throws an error, e.g.,
	// Get "https://auditmanager.ap-northeast-3.amazonaws.com/assessmentFrameworks?frameworkType=Standard": dial tcp: lookup auditmanager.ap-northeast-3.amazonaws.com: no such host
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

	var id string
	if h.Item != nil {
		id = *h.Item.(*auditmanager.AssessmentFrameworkMetadata).Id
	} else {
		location := d.KeyColumnQuals["region"].GetStringValue()
		if location != region {
			return nil, nil
		}
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	params := &auditmanager.GetAssessmentFrameworkInput{
		FrameworkId: aws.String(id),
	}

	op, err := svc.GetAssessmentFramework(params)
	if err != nil {
		plugin.Logger(ctx).Error("getAuditManagerFramework", "err", err)
		return nil, err
	}

	return op.Framework, nil
}
