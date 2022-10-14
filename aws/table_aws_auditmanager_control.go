package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/auditmanager"
	"github.com/aws/aws-sdk-go-v2/service/auditmanager/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAuditManagerControl(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_auditmanager_control",
		Description: "AWS Audit Manager Control",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAuditManagerControl,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ResourceNotFoundException", "ValidationException", "InvalidParameter"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAuditManagerControls,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the specified control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the specified control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier for the specified control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of control, such as custom or standard.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerControl,
			},
			{
				Name:        "created_at",
				Description: "Specifies when the control was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by",
				Description: "The IAM user or role that created the control.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerControl,
			},
			{
				Name:        "action_plan_title",
				Description: "The title of the action plan for remediating the control.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerControl,
			},
			{
				Name:        "action_plan_instructions",
				Description: "The recommended actions to carry out if the control is not fulfilled.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerControl,
			},
			{
				Name:        "control_sources",
				Description: "The data source that determines from where AWS Audit Manager collects evidence for the control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the specified control.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerControl,
			},
			{
				Name:        "last_updated_at",
				Description: "Specifies when the control was most recently updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_updated_by",
				Description: "The IAM user or role that most recently updated the control.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerControl,
			},
			{
				Name:        "testing_information",
				Description: "The steps to follow to determine if the control has been satisfied.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerControl,
			},
			{
				Name:        "control_mapping_sources",
				Description: "The data mapping sources for the specified control.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAuditManagerControl,
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
				Hydrate:     getAuditManagerControl,
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

func listAuditManagerControls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Get client
	svc, err := AuditManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_control.listAuditManagerControls", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	maxItems := int32(100)
	params := &auditmanager.ListControlsInput{
		ControlType: types.ControlTypeStandard,
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

	paginator := auditmanager.NewListControlsPaginator(svc, params, func(o *auditmanager.ListControlsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	// List standard controls
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
			// for the regions where the  Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
			if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_auditmanager_control.listAuditManagerControls", "api_error", err)
			return nil, err
		}

		for _, items := range output.ControlMetadataList {
			d.StreamListItem(ctx, items)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	// List custom controls
	params.ControlType = types.ControlTypeCustom
	paginatorCustom := auditmanager.NewListControlsPaginator(svc, params, func(o *auditmanager.ListControlsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	// List standard controls
	for paginatorCustom.HasMorePages() {
		output, err := paginatorCustom.NextPage(ctx)
		if err != nil {
			// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
			// for the regions where the  Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
			if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_auditmanager_control.listAuditManagerControls", "api_error", err)
			return nil, err
		}

		for _, items := range output.ControlMetadataList {
			d.StreamListItem(ctx, items)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAuditManagerControl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := AuditManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_auditmanager_control.listAuditManagerControls", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	var id string
	if h.Item != nil {
		id = *h.Item.(types.ControlMetadata).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// Handle empty input id
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	params := &auditmanager.GetControlInput{ControlId: &id}

	op, err := svc.GetControl(ctx, params)

	// User with Admin access gets the error as ‘AccessDeniedException: Please complete AWS Audit Manager setup from home page to enable this action in this account’
	// for the regions where the  Audit Manager setup is not complete, this suppresses the value from the regions where the setup is completed.
	if err != nil {
		if strings.Contains(err.Error(), "Please complete AWS Audit Manager setup") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_auditmanager_control.getAuditManagerControl", "api_error", err)
		return nil, err
	}

	return op.Control, nil
}
