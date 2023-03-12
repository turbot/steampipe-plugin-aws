package aws

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchv1 "github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func listDynamicCloudWatchMetricNames(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	// set matrix regions
	regionMatrix := SupportedRegionMatrix(cloudwatchv1.EndpointsID)
	matrixRegions := regionMatrix(ctx, d)
	metricNames := []string{}

	var wg sync.WaitGroup
	metricCh := make(chan *cloudwatch.ListMetricsOutput, len(matrixRegions))
	errorCh := make(chan error, len(matrixRegions))
	for _, matrixRegion := range matrixRegions {
		wg.Add(1)
		go getDynamicCloudWatchMetricAsync(ctx, d, matrixRegion["region"].(string), &wg, metricCh, errorCh)
	}

	// wait for all metrics to be processed
	wg.Wait()

	// NOTE: close channel before ranging over results
	close(metricCh)
	close(errorCh)

	for err := range errorCh {
		// return the first error
		plugin.Logger(ctx).Error("listDynamicCloudWatchMetricNames", "channel_error", err)
		return nil, err
	}

	for metricDetail := range metricCh {
		for _, metric := range metricDetail.Metrics {
			if !helpers.StringSliceContains(metricNames, *metric.MetricName) {
				metricNames = append(metricNames, *metric.MetricName)
			}
		}
	}

	return metricNames, nil
}

func getDynamicCloudWatchMetricAsync(ctx context.Context, d *plugin.QueryData, region string, wg *sync.WaitGroup, metricCh chan *cloudwatch.ListMetricsOutput, errorCh chan error) {
	defer wg.Done()

	rowData, err := getDynamicCloudWatchMetricNames(ctx, d, region)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		metricCh <- rowData
	}
}

func getDynamicCloudWatchMetricNames(ctx context.Context, d *plugin.QueryData, region string) (*cloudwatch.ListMetricsOutput, error) {
	// Get client
	svc, err := CloudWatchDynamicClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("getDynamicCloudWatchMetricNames", "client_error", err)
		return nil, err
	}

	// execute list call
	input := &cloudwatch.ListMetricsInput{}
	output, err := svc.ListMetrics(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("getDynamicCloudWatchMetricNames", "api_error", err)
		return nil, err
	}

	return output, nil
}

func CloudWatchDynamicClient(ctx context.Context, d *plugin.QueryData, region string) (*cloudwatch.Client, error) {
	cfg, err := getClient(ctx, d, region)
	if err != nil {
		return nil, err
	}
	return cloudwatch.NewFromConfig(*cfg), nil
}
