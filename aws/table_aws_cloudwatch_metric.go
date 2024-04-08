package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	cloudwatchv1 "github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudWatchMetric(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_metric",
		Description: "AWS CloudWatch Metric",
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchMetrics,
			Tags:    map[string]string{"service": "cloudwatch", "action": "ListMetrics"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "metric_name",
					Require: plugin.Optional,
				},
				{
					Name:    "namespace",
					Require: plugin.Optional,
				},
				{
					Name:       "dimensions_filter",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "metric_name",
				Description: "The name of the metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace",
				Description: "The namespace for the metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dimensions_filter",
				Description: "The dimensions to filter against.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromQual("dimensions_filter"),
			},
			{
				Name:        "dimensions",
				Description: "The dimensions for the metric.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MetricName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudWatchMetrics(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudWatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_metric.listCloudWatchMetrics", "client_error", err)
		return nil, err
	}

	input := &cloudwatch.ListMetricsInput{}

	// Additional filters
	equalQuals := d.EqualsQuals

	if equalQuals["name"] != nil {
		if equalQuals["name"].GetStringValue() != "" {
			input.MetricName = aws.String(equalQuals["name"].GetStringValue())
		}
	}
	if equalQuals["namespace"] != nil {
		if equalQuals["namespace"].GetStringValue() != "" {
			input.Namespace = aws.String(equalQuals["namespace"].GetStringValue())
		}
	}

	dimensionsFilter := []types.DimensionFilter{}
	dimensionsFilterString := equalQuals["dimensions_filter"].GetJsonbValue()

	if dimensionsFilterString != "" {
		err := json.Unmarshal([]byte(dimensionsFilterString), &dimensionsFilter)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_metric.listCloudWatchMetrics", "unmarshal_error", err)
			return nil, fmt.Errorf("failed to unmarshal dimensions filter %v: %v", dimensionsFilterString, err)
		}
	}

	if len(dimensionsFilter) > 0 {
		input.Dimensions = dimensionsFilter
	}

	paginator := cloudwatch.NewListMetricsPaginator(svc, input, func(o *cloudwatch.ListMetricsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_metric.listCloudWatchMetrics", "api_error", err)
			return nil, err
		}

		for _, metricDetail := range output.Metrics {
			d.StreamListItem(ctx, metricDetail)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
