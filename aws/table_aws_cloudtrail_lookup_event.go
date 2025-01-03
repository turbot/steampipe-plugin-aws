package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"

	cloudtrailv1 "github.com/aws/aws-sdk-go/service/cloudtrail"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

//// TABLE DEFINITION

func tableAwsCloudtrailLookupEvent(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudtrail_lookup_event",
		Description: "AWS CloudTrail Lookup Event",
		List: &plugin.ListConfig{
			Hydrate: listCloudtrailLookupEvents,
			Tags:    map[string]string{"service": "cloudtrail", "action": "LookupEvents"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidLookupAttributesException"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "event_id", Require: plugin.Optional},
				{Name: "event_name", Require: plugin.Optional},
				{Name: "event_source", Require: plugin.Optional},
				{Name: "read_only", Require: plugin.Optional},
				{Name: "end_time", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "start_time", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "resource_name", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "resource_type", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "username", Require: plugin.Optional},
				{Name: "access_key_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudtrailv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "event_id",
				Description: "The CloudTrail ID of the event returned.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_name",
				Description: "The name of the event returned.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_source",
				Description: "The Amazon Web Services service to which the request was made.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "read_only",
				Description: "Information about whether the event is a write event or a read event.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "access_key_id",
				Description: "The AWS access key ID that was used to sign the request. If the request was made with temporary security credentials, this is the access key ID of the temporary credentials.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_time",
				Description: "The date and time of the event returned.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "end_time",
				Description: "Specifies that only events that occur before or at the specified time are returned. If the specified end time is before the specified start time, an error is returned.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromQual("end_time"),
			},
			{
				Name:        "start_time",
				Description: "Specifies that only events that occur after or at the specified time are returned. If the specified start time is after the specified end time, an error is returned.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromQual("start_time"),
			},
			{
				Name:        "resource_name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("resource_name"),
			},
			{
				Name:        "resource_type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("resource_type"),
			},
			{
				Name:        "username",
				Description: "A user name or role name of the requester that called the API in the event returned.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resources",
				Description: "A list of resources referenced by the event returned.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cloudtrail_event",
				Description: "A JSON string that contains a representation of the event returned.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cloud_trail_event",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release, use the cloudtrail_event column instead. A JSON string that contains a representation of the event returned.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EventName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudtrailLookupEvents(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_lookup_event.listCloudtrailLookupEvents", "client_error", err)
		return nil, err
	}

	input := buildCloudtrailLookupEventFilter(ctx, d.Quals)
	maxItems := int32(50)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}
	input.MaxResults = &maxItems

	pageLeft := true
	for pageLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		resp, err := svc.LookupEvents(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudtrail_lookup_event.listCloudtrailLookupEvents", "api_error", err)
			return nil, err
		}

		for _, event := range resp.Events {
			d.StreamListItem(ctx, event)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if resp.NextToken != nil {
			input.NextToken = resp.NextToken
		} else {
			pageLeft = false
		}
	}

	return nil, err
}

//// UTILITY FUNCTION

// Build Cloudtrail Lookup Event list call input filter
func buildCloudtrailLookupEventFilter(ctx context.Context, quals plugin.KeyColumnQualMap) *cloudtrail.LookupEventsInput {

	input := &cloudtrail.LookupEventsInput{
		MaxResults: aws.Int32(50),
	}
	attributeKeyMap := map[string]types.LookupAttributeKey{
		"event_id":      types.LookupAttributeKeyEventId,
		"event_name":    types.LookupAttributeKeyEventName,
		"read_only":     types.LookupAttributeKeyReadOnly,
		"username":      types.LookupAttributeKeyUsername,
		"event_source":  types.LookupAttributeKeyEventSource,
		"resource_name": types.LookupAttributeKeyResourceName,
		"resource_type": types.LookupAttributeKeyResourceType,
		"access_key_id": types.LookupAttributeKeyAccessKeyId,
	}

	var lookupAttributes []types.LookupAttribute
	for columnName, attributeKey := range attributeKeyMap {
		if quals[columnName] != nil {
			var value interface{}
			if columnName == "read_only" {
				value = getQualsValueByColumn(quals, columnName, "bool")
			} else {
				value = getQualsValueByColumn(quals, columnName, "string")
			}
			lookupAttribute := types.LookupAttribute{
				AttributeKey:   attributeKey,
				AttributeValue: aws.String(value.(string)),
			}
			lookupAttributes = append(lookupAttributes, lookupAttribute)
		}
	}
	input.LookupAttributes = lookupAttributes

	if quals["start_time"] != nil {
		value := getQualsValueByColumn(quals, "start_time", "time")
		input.StartTime = aws.Time(value.(time.Time))
	}

	if quals["end_time"] != nil {
		value := getQualsValueByColumn(quals, "end_time", "time")
		input.EndTime = aws.Time(value.(time.Time))
	}

	return input
}
