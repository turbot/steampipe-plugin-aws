package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/health"
)

func tableAwsHealthEvent(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_health_event",
		Description: "AWS Health Event",
		List: &plugin.ListConfig{
			Hydrate: listHealthEvents,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "arn", Require: plugin.Optional},
				{Name: "availability_zone", Require: plugin.Optional},
				{Name: "end_time", Require: plugin.Optional},
				{Name: "event_type_category", Require: plugin.Optional},
				{Name: "event_type_code", Require: plugin.Optional},
				{Name: "last_updated_time", Require: plugin.Optional},
				{Name: "service", Require: plugin.Optional},
				{Name: "start_time", Require: plugin.Optional},
				{Name: "status_code", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the HealthEvent.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zone",
				Description: "The Amazon Web Services Availability Zone of the event.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_time",
				Description: "The date and time that the event began.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "end_time",
				Description: "The date and time that the event ended.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "event_scope_code",
				Description: "This parameter specifies if the Health event is a public Amazon Web Services service event or an account-specific event.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_type_category",
				Description: "A list of event type category codes. Possible values are issue, accountNotification, or scheduledChange.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_type_code",
				Description: "The unique identifier for the event type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_time",
				Description: "The most recent date and time that the event was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "service",
				Description: "The Amazon Web Services service that is affected by the event. For example, EC2, RDS.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_code",
				Description: "The most recent status of the event. Possible values are open, closed, and upcoming.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn"),
			},
		}),
	}
}

//// LIST FUNCTION

func listHealthEvents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listHealthEvents")

	// Create Session
	svc, err := HealthService(ctx, d)
	if err != nil {
		return nil, err
	}

	filter := buildHealthEventFilter(d)

	params := &health.DescribeEventsInput{
		MaxResults: aws.Int64(100),
	}

	if filter != nil {
		params.Filter = filter
	}
	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxResults {
			if *limit < 10 {
				params.MaxResults = aws.Int64(10)
			} else {
				params.MaxResults = limit
			}
		}
	}

	// List IAM user access keys
	err = svc.DescribeEventsPages(
		params,
		func(page *health.DescribeEventsOutput, isLast bool) bool {
			for _, event := range page.Events {
				d.StreamListItem(ctx, event)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listHealthEvents", "listHealthEventsPages", err)
	}

	return nil, err
}

/// UTILITY FUNCTION
// Build health event list call input filter
func buildHealthEventFilter(d *plugin.QueryData) *health.EventFilter {
	filter := &health.EventFilter{}

	filterQuals := map[string]string{
		"arn":                 "string",
		"availability_zone":   "string",
		"status_code":         "string",
		"event_type_category": "string",
		"event_type_code":     "string",
		"last_updated_time":   "time",
		"service":             "string",
		"start_time":          "time",
		"end_time":            "time",
	}

	for columnName, dataType := range filterQuals {
		if dataType == "string" && d.KeyColumnQualString(columnName) != "" {
			value := d.KeyColumnQualString(columnName)
			switch columnName {
			case "arn":
				filter.SetEventArns([]*string{aws.String(value)})
			case "availability_zone":
				filter.SetAvailabilityZones([]*string{aws.String(value)})
			case "status_code":
				filter.SetEventStatusCodes([]*string{aws.String(value)})
			case "event_type_category":
				filter.SetEventTypeCategories([]*string{aws.String(value)})
			case "event_type_code":
				filter.SetEventTypeCodes([]*string{aws.String(value)})
			case "service":
				filter.SetServices([]*string{aws.String(value)})
			}
		}
		if dataType == "time" {
			if d.Quals[columnName] != nil {
				for _, q := range d.Quals[columnName].Quals {
					if q.Value.GetTimestampValue() != nil {
						t := &health.DateTimeRange{
							From: aws.Time(q.Value.GetTimestampValue().AsTime()),
							To:   aws.Time(q.Value.GetTimestampValue().AsTime()),
						}
						switch columnName {
						case "last_updated_time":
							filter.SetLastUpdatedTimes([]*health.DateTimeRange{t})
						case "start_time":
							filter.SetStartTimes([]*health.DateTimeRange{t})
						case "end_time":
							filter.SetEndTimes([]*health.DateTimeRange{t})
						}
					}
				}
			}
		}
	}

	return filter
}
