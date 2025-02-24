package aws

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSDBMaintenanceAction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_maintenance_action",
		Description: "Lists pending maintenance actions for Amazon RDS instances and clusters.",
		List: &plugin.ListConfig{
			Hydrate: getRDSDBClusterPendingMaintenanceActionCopy,
		},
		Columns: []*plugin.Column{
			{
				Name:        "resource_identifier",
				Description: "The ARN of the resource with pending maintenance actions.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceIdentifier"),
			},
			{
				Name:        "action",
				Description: "The type of pending maintenance action.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "opt_in_status",
				Description: "The opt-in status for the maintenance action.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_applied_after_date",
				Description: "The effective date when the pending maintenance action will be automatically applied.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "current_apply_date",
				Description: "The current application date for the maintenance action.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description of the maintenance action.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "forced_apply_date",
				Description: "The date when the maintenance action will be forcibly applied.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
		},
	}
}

func getRDSDBClusterPendingMaintenanceActionCopy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_maintenance_action.", "connection_error", err)
		return nil, err
	}

	op, err := svc.DescribePendingMaintenanceActions(ctx, nil)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_maintenance_action.getRDSDBClusterPendingMaintenanceAction", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Warn("get result", len(op.PendingMaintenanceActions))

	return nil, nil
}
