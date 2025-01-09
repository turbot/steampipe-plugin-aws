package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/health"
	"github.com/aws/aws-sdk-go-v2/service/health/types"

	healthEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsHealthEvent(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_health_event",
		Description: "AWS Health Event",
		List: &plugin.ListConfig{
			Hydrate: listHealthEvents,
			Tags:    map[string]string{"service": "health", "action": "DescribeEvents"},
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
		GetMatrixItemFunc: SupportedRegionMatrix(healthEndpoint.HEALTHServiceID),
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the HealthEvent.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_region",
				Description: "The Amazon Web Services Region name of the event.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region"),
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
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listHealthEvents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := HealthClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_health_event.listHealthEvents", "client error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 10 {
				maxLimit = 10
			} else {
				maxLimit = limit
			}
		}
	}

	filter := buildHealthEventFilter(d)

	input := &health.DescribeEventsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if filter != nil {
		input.Filter = filter
	}

	paginator := health.NewDescribeEventsPaginator(svc, input, func(o *health.DescribeEventsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_health_event.listHealthEvents", "api_error", err)
			return nil, err
		}

		for _, item := range output.Events {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, err
}

// / UTILITY FUNCTION
// Build health event list call input filter
func buildHealthEventFilter(d *plugin.QueryData) *types.EventFilter {
	filter := &types.EventFilter{}

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
		if dataType == "string" && d.EqualsQualString(columnName) != "" {
			value := d.EqualsQualString(columnName)
			switch columnName {
			case "arn":
				filter.EntityArns = ([]string{value})
			case "availability_zone":
				filter.AvailabilityZones = []string{value}
			case "status_code":
				filter.EventStatusCodes = []types.EventStatusCode{
					types.EventStatusCode(value),
				}
			case "event_type_category":
				filter.EventTypeCategories = []types.EventTypeCategory{
					types.EventTypeCategory(value),
				}
			case "event_type_code":
				filter.EventTypeCodes = []string{value}
			case "service":
				filter.Services = []string{value}
			}
		}
		if dataType == "time" {
			if d.Quals[columnName] != nil {
				for _, q := range d.Quals[columnName].Quals {
					if q.Value.GetTimestampValue() != nil {
						t := &types.DateTimeRange{
							From: aws.Time(q.Value.GetTimestampValue().AsTime()),
							To:   aws.Time(q.Value.GetTimestampValue().AsTime()),
						}
						switch columnName {
						case "last_updated_time":
							filter.LastUpdatedTimes = []types.DateTimeRange{*t}
						case "start_time":
							filter.StartTimes = []types.DateTimeRange{*t}
						case "end_time":
							filter.EndTimes = []types.DateTimeRange{*t}
						}
					}
				}
			}
		}
	}

	return filter
}
