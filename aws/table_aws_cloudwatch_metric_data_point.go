package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	cloudwatchv1 "github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudWatchMetricDataPoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_metric_data_point",
		Description: "AWS CloudWatch Metric Data Point",
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchMetricDataPoints,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "id",
					Require: plugin.Required,
				},
				{
					Name:       "account_id",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "expression",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "metric_stat",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "period",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "return_data",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "scan_by",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "timestamp",
					Operators:  []string{">", ">=", "=", "<", "<="},
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "timezone",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The short name you specified to represent this metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "label",
				Description: "The human-readable label associated with the data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "period",
				Description: "The granularity, in seconds, of the returned data points.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "status_code",
				Description: "The status of the returned data. Complete indicates that all data points in the requested time range were returned. PartialData means that an incomplete set of data points were returned.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timestamp",
				Description: "The timestamp for the data points, formatted in Unix timestamp format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "value",
				Description: "The data point for the metric corresponding to Timestamp.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "expression",
				Description: "This field can contain either a Metrics Insights query, or a metric math expression to be performed on the returned data.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("expression"),
			},
			{
				Name:        "return_data",
				Description: "When used in GetMetricData, this option indicates whether to return the timestamps and raw data values of this metric. If you are performing this call just to do math expressions and do not also need the raw data returned, you can specify false. If you omit this, the default of true is used. When used in PutMetricAlarm, specify true for the one expression result to use as the alarm. For all other metrics and expressions in the same PutMetricAlarm operation, specify ReturnData as False.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromQual("return_data"),
			},
			{
				Name:        "scan_by",
				Description: "The order in which data points should be returned. TimestampDescending returns the newest data first and paginates when the MaxDatapoints limit is reached. TimestampAscending returns the oldest data first and paginates when the MaxDatapoints limit is reached.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("scan_by"),
			},
			{
				Name:        "timezone",
				Description: "You can use timezone to specify your time zone so that the labels of returned data display the correct time for your time zone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("timezone"),
			},
			{
				Name:        "metric_stat",
				Description: "The metric to be returned, along with statistics, period, and units. Use this parameter only if this object is retrieving a metric and not performing a math expression on returned data.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromQual("metric_stat"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Label"),
			},
		}),
	}
}

type MetricDataPoint struct {
	Id         *string
	Label      *string
	StatusCode types.StatusCode
	Period     *int32
	Timestamp  time.Time
	Value      float64
}

//// LIST FUNCTION

func listCloudWatchMetricDataPoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// check if expression & metric_stat both are available, or both are empty
	if (d.EqualsQuals["expression"] != nil && d.EqualsQuals["metric_stat"] != nil) || (d.EqualsQuals["expression"] == nil && d.EqualsQuals["metric_stat"] == nil) {
		return nil, errors.New("please provide either expression or metric_stat in where clause to use this table")
	}

	metricDataQueries := types.MetricDataQuery{}
	metricDataQueries.Id = aws.String(d.EqualsQuals["id"].GetStringValue())

	// Limiting the results
	maxLimit := int32(100800)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	params := &cloudwatch.GetMetricDataInput{
		MaxDatapoints: aws.Int32(maxLimit),
	}
	if d.EqualsQuals["scan_by"] != nil {
		params.ScanBy = types.ScanBy(d.EqualsQuals["scan_by"].GetStringValue())
	}

	if d.EqualsQuals["timezone"] != nil {
		labelOptions := types.LabelOptions{}
		labelOptions.Timezone = aws.String(d.EqualsQuals["timezone"].GetStringValue())
		params.LabelOptions = &labelOptions
	}

	//set the start and end time based on the provided timestamp
	if d.Quals["timestamp"] != nil {
		for _, q := range d.Quals["timestamp"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case "=":
				params.StartTime = aws.Time(timestamp)
				params.EndTime = aws.Time(timestamp)
			case ">=", ">":
				params.StartTime = aws.Time(timestamp)
			case "<", "<=":
				params.EndTime = aws.Time(timestamp)
			}
		}
	} else {
		params.StartTime = aws.Time(time.Now().AddDate(0, 0, -1))
		params.EndTime = aws.Time(time.Now())
	}

	if params.StartTime == nil {
		params.StartTime = aws.Time(time.Now().AddDate(0, 0, -1))
	}
	if params.EndTime == nil {
		params.EndTime = aws.Time(time.Now())
	}

	// set the period based on the duration between the start and end time
	// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudwatch/types#MetricStat.Period
	// * Start time between 3 hours and 15 days ago - Use a multiple of 60 seconds (1 minute).
	// * Start time between 15 and 63 days ago - Use a multiple of 300 seconds (5 minutes).
	// * Start time greater than 63 days ago - Use a multiple of 3600 seconds (1 hour).

	var period int32
	duration := params.EndTime.Sub(*params.StartTime).Round(time.Hour)
	// if the duration is less than 3 hours
	if duration.Hours() <= 3 {
		period = int32(5)
	} else if duration.Hours() <= 360 { // if the duration is between 3 hours and 15 days
		period = int32(60)
	} else if duration.Hours() <= 1512 { // if the duration is between 15 and 63 days
		period = int32(300)
	} else { // if the duration is greater than 63 days
		period = int32(3600)
	}

	// override the period if user has provided it in query
	if d.EqualsQuals["period"] != nil {
		period = int32(d.EqualsQuals["period"].GetInt64Value())
	}

	// set Expression or MetricStat
	if d.EqualsQuals["expression"] != nil {
		metricDataQueries.Expression = aws.String(d.EqualsQuals["expression"].GetStringValue())
		metricDataQueries.Period = aws.Int32(period)
	} else {
		metric_stat := types.MetricStat{}
		metric_stat_string := d.EqualsQuals["metric_stat"].GetJsonbValue()

		if metric_stat_string != "" {
			err := json.Unmarshal([]byte(metric_stat_string), &metric_stat)
			if err != nil {
				plugin.Logger(ctx).Error("listCloudWatchMetricDataPoints", "unmarshal_error", err)
				return nil, fmt.Errorf("failed to unmarshal metric_stat %v: %v", metric_stat_string, err)
			}
		}
		metric_stat.Period = aws.Int32(period)
		metricDataQueries.MetricStat = &metric_stat
	}

	if d.EqualsQuals["account_id"] != nil {
		metricDataQueries.AccountId = aws.String(d.EqualsQuals["account_id"].GetStringValue())
	}

	if d.EqualsQuals["return_data"] != nil {
		metricDataQueries.ReturnData = aws.Bool(d.EqualsQuals["return_data"].GetBoolValue())
	}

	params.MetricDataQueries = []types.MetricDataQuery{metricDataQueries}

	// Get client
	svc, err := CloudWatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCloudWatchMetricDataPoints", "client_error", err)
		return nil, err
	}

	data, err := svc.GetMetricData(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("listCloudWatchMetricDataPoints", "api_error", err)
		return nil, err
	}

	for _, result := range data.MetricDataResults {
		for item := 0; item < len(result.Timestamps); item++ {
			d.StreamListItem(ctx, &MetricDataPoint{result.Id, result.Label, result.StatusCode, aws.Int32(period), result.Timestamps[item], result.Values[item]})
		}
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
