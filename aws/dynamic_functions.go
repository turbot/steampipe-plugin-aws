package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchv1 "github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func listDynamicCloudWatchMetricNames(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	// set matrix regions
	_ = &plugin.Table{GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchv1.EndpointsID)}

	// Get client
	svc, err := CloudWatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listDynamicCloudWatchMetricNames", "client_error", err)
		return nil, err
	}

	// execute list call
	input := &cloudwatch.ListMetricsInput{}
	output, err := svc.ListMetrics(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("listDynamicCloudWatchMetricNames", "api_error", err)
		return nil, err
	}
	metricNames := []string{}
	for _, metricDetail := range output.Metrics {
		if !helpers.StringSliceContains(metricNames, *metricDetail.MetricName) {
			metricNames = append(metricNames, *metricDetail.MetricName)
		}
	}
	return metricNames, nil
}
