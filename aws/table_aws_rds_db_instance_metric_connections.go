package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION
func tableAwsRdsInstanceMetricConnections(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_instance_metric_connections",
		Description: "AWS RDS DB Instance Cloudwatch Metrics - DB Connections",
		List: &plugin.ListConfig{
			ParentHydrate: listRDSDBInstances,
			Hydrate:       listRdsInstanceMetricConnections,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "db_instance_identifier",
					Description: "The friendly name to identify the DB Instance.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listRdsInstanceMetricConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(*rds.DBInstance)
	return listCWMetricStatistics(ctx, d, "5_MIN", "AWS/RDS", "DatabaseConnections", "DBInstanceIdentifier", *instance.DBInstanceIdentifier)
}
