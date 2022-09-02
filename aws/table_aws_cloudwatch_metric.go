package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudWatchMetric(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_metric",
		Description: "AWS CloudWatch Metric",
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchMetrics,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidParameterValue"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "namespace",
					Require: plugin.Optional,
				},
				{
					Name:    "dimension_name",
					Require: plugin.Optional,
				},
				{
					Name:    "dimension_value",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the metric.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MetricName"),
			},
			{
				Name:        "namespace",
				Description: "The namespace for the metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dimension_name",
				Description: "The dimension name for the metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dimension_value",
				Description: "The dimension value for the metric.",
				Type:        proto.ColumnType_STRING,
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

type MetricDetails struct {
	MetricName     string
	Namespace      string
	DimensionName  string
	DimensionValue string
}

//// LIST FUNCTION

func listCloudWatchMetrics(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CloudWatchService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &cloudwatch.ListMetricsInput{}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	dimensionFilters := []*cloudwatch.DimensionFilter{}
	dimensionFilter := cloudwatch.DimensionFilter{}
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

	if d.KeyColumnQualString("dimension_name") != "" && d.KeyColumnQualString("dimension_value") != "" {
		dimensionFilter.Name = aws.String(equalQuals["dimension_name"].GetStringValue())
		dimensionFilter.Value = aws.String(equalQuals["dimension_value"].GetStringValue())
		dimensionFilters = append(dimensionFilters, &dimensionFilter)
	}

	if len(dimensionFilters) > 0 {
		input.Dimensions = dimensionFilters
	}

	// List call
	err = svc.ListMetricsPages(
		input,
		func(page *cloudwatch.ListMetricsOutput, isLast bool) bool {
			for _, metricDetail := range page.Metrics {
				if metricDetail.Dimensions == nil {
					d.StreamListItem(ctx, &MetricDetails{
						MetricName: *metricDetail.MetricName,
						Namespace:  *metricDetail.Namespace,
					})
				} else {
					for _, dimension := range metricDetail.Dimensions {
						d.StreamListItem(ctx, &MetricDetails{
							MetricName:     *metricDetail.MetricName,
							Namespace:      *metricDetail.Namespace,
							DimensionName:  *dimension.Name,
							DimensionValue: *dimension.Value,
						})
					}

					// Context can be cancelled due to manual cancellation or the limit has been hit
					if d.QueryStatus.RowsRemaining(ctx) == 0 {
						return false
					}
				}

			}
			return !isLast
		},
	)

	return nil, err
}
