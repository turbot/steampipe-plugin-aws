package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRdsInstanceMetricReadIopsHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_instance_metric_read_iops_hourly",
		Description: "AWS RDS DB Instance Cloudwatch Metrics - Read IOPS (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listRDSDBInstances,
			Hydrate:       listRdsInstanceMetricReadIopsHourly,
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

func listRdsInstanceMetricReadIopsHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(types.DBInstance)
	return listCWMetricStatistics(ctx, d, "HOURLY", "AWS/RDS", "ReadIOPS", "DBInstanceIdentifier", *instance.DBInstanceIdentifier)
}
