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

func tableAwsRDSDBMaintenanceAction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_maintenance_action",
		Description: "Lists pending maintenance actions for Amazon RDS instances and clusters.",
		List: &plugin.ListConfig{
			Hydrate: listRDSMaintenanceActions,
			Tags:    map[string]string{"service": "rds", "action": "DescribePendingMaintenanceActions"},
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
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceIdentifier"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceIdentifier").Transform(arnToAkas),
			},
		}),
	}
}

func listRDSMaintenanceActions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_maintenance_action.listMaintenanceAction", "connection_error", err)
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
	type result struct {
		ResourceIdentifier string
		types.PendingMaintenanceAction
	}
	paginator := rds.NewDescribePendingMaintenanceActionsPaginator(client, &rds.DescribePendingMaintenanceActionsInput{
		MaxRecords: &maxLimit,
	})
	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_maintenance_action.listMaintenanceAction", "api_error", err)
			return nil, err
		}

		for _, action := range output.PendingMaintenanceActions {
			for _, detail := range action.PendingMaintenanceActionDetails {
				r := &result{
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
