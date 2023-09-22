package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAwsDynamoDBMetricAccountProvisionedReadCapacityUtilization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_metric_account_provisioned_read_capacity_util",
		Description: "AWS DynamoDB Metric Account Provisioned Read Capacity Utilization",
		List: &plugin.ListConfig{
			Hydrate: listDynamoDBMetricAccountProvisionedReadCapacityUtilization,
			Tags:    map[string]string{"service": "dynamodb", "action": "GetMetricStatistics"},
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
		Columns:           awsRegionalColumns(cwMetricColumns([]*plugin.Column{})),
	}
}

//// LIST FUNCTION

func listDynamoDBMetricAccountProvisionedReadCapacityUtilization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listCWMetricStatistics(ctx, d, "5_MIN", "AWS/DynamoDB", "AccountProvisionedReadCapacityUtilization", "", "")
}
