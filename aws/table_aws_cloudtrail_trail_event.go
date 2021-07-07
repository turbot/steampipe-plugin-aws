package aws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/ettle/strcase"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudtrailTrailEvent(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudtrail_trail_event",
		Description: "CloudTrail evnts from CloudWatch Logs",
		// DefaultTransform: transform.FromCamel(),
		List: &plugin.ListConfig{
			Hydrate:    listCloudTrailEvents,
			KeyColumns: tableAwsCloudtrailEventsListKeyColumns(),
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			// CloudTrail
			{Name: "event_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time the request was made, in coordinated universal time (UTC)."},
			// {Name: "event_version", Type: proto.ColumnType_STRING, Description: "The version of the log event format."},
			{Name: "event_source", Type: proto.ColumnType_STRING, Description: "The AWS service that the request was made to."},
			{Name: "event_name", Type: proto.ColumnType_STRING, Description: "The name of the event returned."},
			{Name: "resource_type", Type: proto.ColumnType_STRING, Description: "The name of the event returned."},
			{Name: "resource_name", Type: proto.ColumnType_STRING, Description: "The name of the event returned."},
			{Name: "username", Type: proto.ColumnType_STRING, Description: "The name of the event returned."},
			{Name: "read_only", Type: proto.ColumnType_BOOL, Description: "Information about whether the event is a write event or a read event."},
			{Name: "access_key_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("AccessKeyId"), Description: "The AWS access key ID that was used to sign the request. If the request was made with temporary security credentials, this is the access key ID of the temporary credentials."},
			{Name: "cloudtrail_event", Type: proto.ColumnType_JSON, Transform: transform.FromField("CloudTrailEvent"), Description: "The CloudTrail event."},

			{Name: "resources", Type: proto.ColumnType_JSON, Description: "A list of resources referenced by the event returned."},
			// {Name: "aws_region", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "source_ip_address", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "event_category", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "event_type", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "request_id", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "shared_event_id", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "recipient_account_id", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "vpc_endpoint_id", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "user_agent", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "user_identity", Type: proto.ColumnType_JSON, Description: "Information about the user that made a request."},
			// {Name: "request_parameters", Type: proto.ColumnType_JSON, Description: ""},
			// {Name: "response_elements", Type: proto.ColumnType_JSON, Description: ""},
			// {Name: "additional_event_data", Type: proto.ColumnType_JSON, Description: ""},
			// {Name: "tls_details", Type: proto.ColumnType_JSON, Description: ""},
			// // Other columns
			// {Name: "parsed_message", Description: "", Type: proto.ColumnType_STRING, Transform: transform.FromValue()},
			// {Name: "message", Description: "", Type: proto.ColumnType_STRING, Transform: transform.FromField("Message")},
			{Name: "event_id", Description: "The ID of the event.", Type: proto.ColumnType_STRING, Transform: transform.FromField("EventId")},
			// {Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter pattern for the search."},
			// {Name: "ingestion_time", Description: "The time when the event was ingested.", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("IngestionTime").Transform(transform.UnixMsToTimestamp)},
		}),
	}
}

type cloudtrailEvent struct {
	// TODO - apiVersion
	// TODO - errorCode
	// TODO - errorMessage
	// TODO - serviceEventDetails
	// TODO - addendum
	// TODO - sessionCredentialFromConsole
	// TODO - edgeDeviceDetails

	AwsRegion           *string      `type:"string"`
	SourceIpAddress     *string      `json:"sourceIPAddress" type:"string"`
	EventCategory       *string      `type:"string"`
	SharedEventId       *string      `json:"sharedEventID" type:"string"`
	RecipientAccountId  *string      `type:"string"`
	VpcEndpointId       *string      `json:"vpcEndpointId" type:"string"`
	UserAgent           *string      `type:"string"`
	UserIdentity        *interface{} `type:"map"`
	RequestParameters   *interface{} `type:"map"`
	ResponseElements    *interface{} `type:"map"`
	AdditionalEventData *interface{} `type:"map"`
	TlsDetails          *interface{} `type:"map"`
	ManagementEveent    *bool        `type:"bool"`
	EventType           *string      `type:"string"`
	RequestId           *string      `type:"string"`

	// The AWS access key ID that was used to sign the request. If the request was made with temporary security credentials, this is the access key ID of the temporary credentials.
	AccessKeyId *string `type:"string"`
	// A JSON string that contains a representation of the event returned.
	CloudTrailEvent *string `type:"string"`
	// The CloudTrail ID of the event returned.
	EventId *string `type:"string"`
	// The name of the event returned.
	EventVersion *string `type:"string"`
	// The name of the event returned.
	EventName *string `type:"string"`
	// The AWS service that the request was made to.
	EventSource *string `type:"string"`
	// The date and time of the event returned.
	EventTime *time.Time `type:"timestamp"`
	// Information about whether the event is a write event or a read event.
	ReadOnly *bool `type:"bool"`
	// A list of resources referenced by the event returned.
	Resources []*interface{} `type:"list"`
	// A user name or role name of the requester that called the API in the event
	// returned.
	Username *string `type:"string"`
	// contains filtered or unexported fields
}

func getCloudtrailMessageField(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	e := h.Item.(*cloudwatchlogs.FilteredLogEvent)
	cte := cloudtrailEvent{}
	err := json.Unmarshal([]byte(*e.Message), &cte)
	if err != nil {
		return nil, err
	}
	return cte, nil
}

func listCloudTrailEvents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create session
	svc, err := CloudTrailService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := cloudtrail.LookupEventsInput{}

	equalQuals := d.KeyColumnQuals // Key Quals
	quals := d.Quals               // Other Quals

	lookupAttributes1 := getLookupAttributes(equalQuals)
	plugin.Logger(ctx).Error("listCloudTrailEvents", "lookupAttributes1", lookupAttributes1)

	// var lookupAttributes []*cloudtrail.LookupAttribute
	// lookupAttributes = append(lookupAttributes, &cloudtrail.LookupAttribute{
	// 	AttributeKey:   aws.String(toPascalCase("event_name")),
	// 	AttributeValue: aws.String(equalQuals["event_name"].GetStringValue()),
	// })

	if lookupAttributes1 != nil {
		input.LookupAttributes = lookupAttributes1
	}

	if equalQuals["event_category"] != nil {
		input.EventCategory = aws.String(equalQuals["event_category"].GetStringValue())
	}

	if quals["start_time"] != nil {
		for _, q := range quals["start_time"].Quals {
			startTime := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case "=", ">=", ">", "<", "<=":
				input.StartTime = aws.Time(startTime)
			}
		}
	}
	if quals["end_time"] != nil {
		for _, q := range quals["start_time"].Quals {
			startTime := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case "=", ">=", ">", "<", "<=":
				input.EndTime = aws.Time(startTime)
			}
		}
	}

	plugin.Logger(ctx).Error("listCloudTrailEvents", "Input", input)

	err = svc.LookupEventsPages(
		&input,
		func(page *cloudtrail.LookupEventsOutput, _ bool) bool {
			for _, trailEvent := range page.Events {
				d.StreamListItem(ctx, trailEvent)
			}
			// Abort if we've been cancelled, which probably means we've reached the requested limit
			select {
			case <-ctx.Done():
				return false
			default:
				return true
			}
		},
	)

	// Handle log group not found errors gracefully
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == "ResourceNotFoundException" {
			return nil, nil
		}
	}

	return nil, err
}

// https://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/CloudTrail.html#lookupEvents-property
func tableAwsCloudtrailEventsListKeyColumns() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{
		{Name: "event_category", Require: plugin.Optional},
		// {Name: "end_time", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
		// {Name: "start_time", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},

		// LookupAttributes
		{Name: "event_id", Require: plugin.Optional},
		{Name: "event_name", Require: plugin.Optional},
		{Name: "read_only", Require: plugin.Optional},
		{Name: "username", Require: plugin.Optional},
		{Name: "resource_type", Require: plugin.Optional},
		{Name: "resource_name", Require: plugin.Optional},
		{Name: "event_source", Require: plugin.Optional},
		{Name: "access_key_id", Require: plugin.Optional},
	}
}

func getLookupAttributes(equalQuals plugin.KeyColumnEqualsQualMap) []*cloudtrail.LookupAttribute {
	var lookupAttributes []*cloudtrail.LookupAttribute

	lookupKeys := []string{"event_id", "event_name", "read_only", "username", "resource_type", "resource_name", "event_source", "access_key_id"}

	for _, key := range lookupKeys {
		value := getLookupAttribute(equalQuals, key)
		if value != nil {
			lookupAttributes = append(lookupAttributes, value)
		}
	}
	return lookupAttributes
}

func getLookupAttribute(equalQuals plugin.KeyColumnEqualsQualMap, key string) *cloudtrail.LookupAttribute {
	if equalQuals[key] != nil {
		return &cloudtrail.LookupAttribute{
			AttributeKey:   aws.String(toPascalCase(key)),
			AttributeValue: aws.String(equalQuals[key].GetStringValue()),
		}
	}
	return nil
}

func toPascalCase(key string) string {
	return strcase.ToPascal(key)
}
