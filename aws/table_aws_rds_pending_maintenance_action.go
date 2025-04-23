package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	rdsv1 "github.com/aws/aws-sdk-go/service/rds"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSPendingMaintenanceAction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_pending_maintenance_action",
		Description: "AWS RDS Pending Maintenance Action",
		List: &plugin.ListConfig{
			Hydrate: listRDSMaintenanceActions,
			Tags:    map[string]string{"service": "rds", "action": "DescribePendingMaintenanceActions"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "resource_identifier",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "resource_identifier",
				Description: "The Amazon Resource Name (ARN) of the resource that the pending maintenance action applies to.",
				Type:        proto.ColumnType_STRING,
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

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceIdentifier"),
			},
		}),
	}
}

type rdsMaintenanceActionResult struct {
	ResourceIdentifier string
	types.PendingMaintenanceAction
}

func listRDSMaintenanceActions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_pending_maintenance_action.listMaintenanceAction", "connection_error", err)
		return nil, err
	}

	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 20 {
				maxLimit = 20
			} else {
				maxLimit = limit
			}
		}
	}
	input := &rds.DescribePendingMaintenanceActionsInput{
		MaxRecords: &maxLimit,
	}
	if resourceIdentifier := d.EqualsQualString("resource_identifier"); resourceIdentifier != "" {
		input.ResourceIdentifier = &resourceIdentifier
	}

	paginator := rds.NewDescribePendingMaintenanceActionsPaginator(client, input)
	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_pending_maintenance_action.listMaintenanceAction", "api_error", err)
			return nil, err
		}

		for _, action := range output.PendingMaintenanceActions {
			for _, detail := range action.PendingMaintenanceActionDetails {
				r := &rdsMaintenanceActionResult{
					ResourceIdentifier:       *action.ResourceIdentifier,
					PendingMaintenanceAction: detail,
				}
				d.StreamListItem(ctx, r)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}
	return nil, nil
}
