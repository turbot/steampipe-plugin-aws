package aws

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsDynamoDBMetricAccountReadThroughput(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_metric_account_read_throughput",
		Description: "AWS DynamoDB Metric Account Read Throughput",
		List: &plugin.ListConfig{
			Hydrate: listDynamboDbMetricAccountReadThroughput,
		},
		GetMatrixItem: BuildRegionList,
		Columns:       awsRegionalColumns(cwMetricColumns([]*plugin.Column{})),
	}
}

//// LIST FUNCTION

func listDynamboDbMetricAccountReadThroughput(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listCWMetricStatistics(ctx, d, "5_MIN", "AWS/DynamoDB", "AccountProvisionedReadCapacityUtilization", "", "")
}
