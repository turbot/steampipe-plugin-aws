package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/auditmanager"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAuditManagerControl(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_audit_manager_control",
		Description: "AWS Audit Manager Control",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAuditManagerControl,
		},
		List: &plugin.ListConfig{
			Hydrate: listAuditManagerControls,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the specified control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the specified control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the specified control.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerControl,
			},
			{
				Name:        "type",
				Description: "The type of control, such as custom or standard.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerControl,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the specified control.",
				Type:        proto.ColumnType_STRING,
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
				Name:        "control_sources",
				Description: "The data source that determines from where AWS Audit Manager collects evidence for the control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "action_plan_instructions",
				Description: "The recommended actions to carry out if the control is not fulfilled.",
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
				Name:        "testing_information",
				Description: "The steps to follow to determine if the control has been satisfied.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAuditManagerControl,
			},

			// JSON Column

			{
				Name:        "control_mapping_sources",
				Description: "The data mapping sources for the specified control.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ControlMappingSources"),
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
				Transform:   transform.FromField("Tags"),
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

func listAuditManagerControls(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAuditManagerControls", "AWS_REGION", region)
	// restAPI := h.Item.(*apigateway.RestApi)

	// Create Session
	svc, err := AuditManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	standardControltype := "Standard"

	err = svc.ListControlsPages(
		&auditmanager.ListControlsInput{
			ControlType: &standardControltype,
		},
		func(page *auditmanager.ListControlsOutput, lastPage bool) bool {
			for _, items := range page.ControlMetadataList {
				d.StreamListItem(ctx, items)
			}
			return !lastPage
		},
	)

	customControltype := "Custom"

	err = svc.ListControlsPages(
		&auditmanager.ListControlsInput{
			ControlType: &customControltype,
		},
		func(page *auditmanager.ListControlsOutput, lastPage bool) bool {
			for _, items := range page.ControlMetadataList {
				d.StreamListItem(ctx, items)
			}
			return !lastPage
		},
	)

	return nil, err
}

//// GET FUNCTION

func getAuditManagerControl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAuditManagerControl")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := AuditManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	var id string
	if h.Item != nil {
		id = *h.Item.(*auditmanager.ControlMetadata).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	params := &auditmanager.GetControlInput{
		ControlId: &id,
	}

	op, err := svc.GetControl(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAuditManagerControl", "ERROR", err)
		return nil, err
	}

	plugin.Logger(ctx).Trace("BalBabu", op.Control)

	return op.Control, nil
}
