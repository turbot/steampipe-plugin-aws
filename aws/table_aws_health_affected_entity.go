package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	healthEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/health"
	"github.com/aws/aws-sdk-go-v2/service/health/types"
)

func tableAwsHealthAffectedEntity(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_health_affected_entity",
		Description: "AWS Health Affected Entity",
		List: &plugin.ListConfig{
			ParentHydrate: listHealthEvents,
			Hydrate:       listHealthAffectedEntities,
			Tags:          map[string]string{"service": "health", "action": "DescribeAffectedEntities"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"SubscriptionRequiredException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "arn", Require: plugin.Optional},
				{Name: "event_arn", Require: plugin.Optional},
				{Name: "entity_value", Require: plugin.Optional},
				{Name: "status_code", Require: plugin.Optional},
				{Name: "last_updated_time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<=", "="}},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(healthEndpoint.AWS_HEALTH_SERVICE_ID),
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the health entity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EntityArn"),
			},
			{
				Name:        "entity_url",
				Description: "The URL of the affected entity.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "entity_value",
				Description: "The ID of the affected entity.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_arn",
				Description: "The Amazon Resource Name (ARN) of the health event.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_time",
				Description: "The most recent time that the entity was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status_code",
				Description: "The most recent status of the entity affected by the event. The possible values are IMPAIRED, UNIMPAIRED, and UNKNOWN.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EntityArn").Transform(transform.NullIfZeroValue).Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listHealthAffectedEntities(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	event := h.Item.(types.Event)

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

	// Validate if user provided input matches hydrate value
	if d.EqualsQuals["event_arn"] != nil && d.EqualsQualString("event_arn") != *event.Arn {
		return nil, nil
	}

	filter := buildHealthAffectedEntityFilter(d)
	filter.EventArns = []string{*event.Arn}
	input := &health.DescribeAffectedEntitiesInput{
		MaxResults: aws.Int32(maxLimit),
		Filter:     filter,
	}

	svc, err := HealthClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_health_affected_entity.listHealthAffectedEntities", "client_error", err)
		return nil, err
	}

	paginator := health.NewDescribeAffectedEntitiesPaginator(svc, input, func(o *health.DescribeAffectedEntitiesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_health_affected_entity.listHealthAffectedEntities", "api_error", err)
			return nil, err
		}

		for _, item := range output.Entities {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, err
}

//// UTILITY FUNCTION
// Build health affected entity list call input filter
func buildHealthAffectedEntityFilter(d *plugin.QueryData) *types.EntityFilter {
	filter := &types.EntityFilter{}

	filterQuals := map[string]string{
		"arn":               "string",
		"entity_value":      "string",
		"last_updated_time": "time",
		"status_code":       "string",
	}

	for columnName, dataType := range filterQuals {
		if dataType == "string" && d.EqualsQualString(columnName) != "" {
			value := d.EqualsQualString(columnName)
			switch columnName {
			case "arn":
				filter.EntityArns = ([]string{value})
			case "status_code":
				filter.StatusCodes = []types.EntityStatusCode{
					types.EntityStatusCode(value),
				}
			}
		}
		if dataType == "time" {
			if d.Quals[columnName] != nil {
				for _, q := range d.Quals[columnName].Quals {
					if q.Value.GetTimestampValue() != nil {
						t := &types.DateTimeRange{}
						switch q.Operator {
						case ">=", ">":
							t.From = aws.Time(q.Value.GetTimestampValue().AsTime())
						case "<=", "<":
							t.To = aws.Time(q.Value.GetTimestampValue().AsTime())
						case "=":
							t.From = aws.Time(q.Value.GetTimestampValue().AsTime())
							t.To = aws.Time(q.Value.GetTimestampValue().AsTime())
						}
						switch columnName {
						case "last_updated_time":
							filter.LastUpdatedTimes = []types.DateTimeRange{*t}
						}
					}
				}
			}
		}
	}

	return filter
}
