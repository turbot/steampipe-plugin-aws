package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/auditmanager"
	"github.com/aws/aws-sdk-go-v2/service/auditmanager/types"

	auditmanagerEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAuditManagerFramework(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_auditmanager_framework",
		Description: "AWS Audit Manager Framework",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "region"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException", "InternalServerException"}),
			},
			Hydrate: getAuditManagerFramework,
			Tags:    map[string]string{"service": "auditmanager", "action": "GetAssessmentFramework"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAuditManagerFrameworks,
			Tags:    map[string]string{"service": "auditmanager", "action": "ListAssessmentFrameworks"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAuditManagerFramework,
				Tags: map[string]string{"service": "auditmanager", "action": "GetAssessmentFramework"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(auditmanagerEndpoint.AWS_AUDITMANAGER_SERVICE_ID),
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
	svc, err := AuditManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_framework.listAuditManagerFrameworks", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(100)
	params := &auditmanager.ListAssessmentFrameworksInput{
		FrameworkType: types.FrameworkTypeStandard,
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
	paginator := auditmanager.NewListAssessmentFrameworksPaginator(svc, params, func(o *auditmanager.ListAssessmentFrameworksPaginatorOptions) {
		o.Limit = 32
		o.StopOnDuplicateToken = true
	})

	// List standard audit manager frameworks
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
			plugin.Logger(ctx).Error("aws_auditmanager_framework.listAuditManagerFrameworks", "api_error", err)
			return nil, err
		}

		for _, framework := range output.FrameworkMetadataList {
			d.StreamListItem(ctx, framework)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	params.FrameworkType = types.FrameworkTypeCustom
	paginatorCustom := auditmanager.NewListAssessmentFrameworksPaginator(svc, params, func(o *auditmanager.ListAssessmentFrameworksPaginatorOptions) {
		o.Limit = 32
		o.StopOnDuplicateToken = true
	})

	// List standard audit manager frameworks
	for paginatorCustom.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginatorCustom.NextPage(ctx)
		if err != nil {
			// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
			// for the regions where the  Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
			if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_auditmanager_framework.listAuditManagerFrameworks", "api_error", err)
			return nil, err
		}

		for _, framework := range output.FrameworkMetadataList {
			d.StreamListItem(ctx, framework)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAuditManagerFramework(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := AuditManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_framework.getAuditManagerFramework", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var id string
	if h.Item != nil {
		id = *h.Item.(types.AssessmentFrameworkMetadata).Id
	} else {
		region := d.EqualsQualString(matrixKeyRegion)
		location := d.EqualsQuals["region"].GetStringValue()
		if location != region {
			return nil, nil
		}
		id = d.EqualsQuals["id"].GetStringValue()
	}

	params := &auditmanager.GetAssessmentFrameworkInput{
		FrameworkId: aws.String(id),
	}

	op, err := svc.GetAssessmentFramework(ctx, params)

	// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
	// for the regions where the Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
	if err != nil {
		if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_auditmanager_framework.getAuditManagerFramework", "api_error", err)
		return nil, err
	}

	return *op.Framework, nil
}
