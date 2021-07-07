package aws

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsDynamoDBMetricAccountWriteThroughput(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_metric_account_write_throughput",
		Description: "AWS DynamoDB Metric Account Write Throughput",
		List: &plugin.ListConfig{
			Hydrate: listDynamboDbMetricAccountWriteThroughput,
		},
		GetMatrixItem: BuildRegionList,
		Columns:       awsRegionalColumns(cwMetricColumns([]*plugin.Column{})),
	}
}

//// LIST FUNCTION

func listDynamboDbMetricAccountWriteThroughput(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listCWMetricStatistics(ctx, d, "5_MIN", "AWS/DynamoDB", "AccountProvisionedWriteCapacityUtilization", "", "")
}
