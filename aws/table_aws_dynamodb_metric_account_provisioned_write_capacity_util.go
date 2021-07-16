package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsDynamoDBMetricAccountProvisionedWriteCapacityUtilization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_metric_account_provisioned_write_capacity_util",
		Description: "AWS DynamoDB Metric Account Provisioned Write Capacity Utilization",
		List: &plugin.ListConfig{
			Hydrate: listDynamboDbMetricAccountProvisionedWriteCapacityUtilization,
		},
		GetMatrixItem: BuildRegionList,
		Columns:       awsRegionalColumns(cwMetricColumns([]*plugin.Column{})),
	}
}

//// LIST FUNCTION

func listDynamboDbMetricAccountProvisionedWriteCapacityUtilization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listCWMetricStatistics(ctx, d, "5_MIN", "AWS/DynamoDB", "AccountProvisionedWriteCapacityUtilization", "", "")
}
