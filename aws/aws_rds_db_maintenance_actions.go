package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	rdsv1 "github.com/aws/aws-sdk-go/service/rds"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"strings"
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
				Name:        "name",
				Description: "The name for the resource that the pending maintenance action applies to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_cluster",
				Description: "Indicates whether the resource is a cluster.",
				Type:        proto.ColumnType_BOOL,
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
		}),
	}
}

func listRDSMaintenanceActions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_maintenance_action.listMaintenanceAction", "connection_error", err)
		return nil, err
	}
	var maxItems = calculateMaxLimit[int32](100, d)
	type result struct {
		Name      string
		IsCluster bool
		types.PendingMaintenanceAction
	}
	paginator := rds.NewDescribePendingMaintenanceActionsPaginator(client, &rds.DescribePendingMaintenanceActionsInput{
		MaxRecords: &maxItems,
	})
	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_maintenance_action.listMaintenanceAction", "api_error", err)
			return nil, err
		}

		for _, action := range output.PendingMaintenanceActions {
			splitIdentifier := strings.Split(*action.ResourceIdentifier, ":")
			n := len(splitIdentifier)
			for _, detail := range action.PendingMaintenanceActionDetails {
				r := &result{
					Name:                     splitIdentifier[n-1],
					IsCluster:                splitIdentifier[n-2] == "cluster",
					PendingMaintenanceAction: detail,
				}
				d.StreamListItem(ctx, r)
			}
		}
	}
	return nil, nil
}
