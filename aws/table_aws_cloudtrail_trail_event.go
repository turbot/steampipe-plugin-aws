package aws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudtrailTrailEvent(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudtrail_trail_event",
		Description: "CloudTrail evnts from CloudWatch Logs",
		List: &plugin.ListConfig{
			Hydrate:    listCloudwatchLogEvents,
			KeyColumns: tableAwsCloudwatchLogEventListKeyColumns(),
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			// Top columns
			{Name: "log_group_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("log_group_name"), Description: "The name of the log group to which this event belongs."},
			{Name: "log_stream_name", Type: proto.ColumnType_STRING, Description: "The name of the log stream to which this event belongs."},
			{Name: "timestamp_ms", Type: proto.ColumnType_INT, Transform: transform.FromField("Timestamp"), Description: "The time when the event occurred."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Timestamp").Transform(transform.UnixMsToTimestamp), Description: "The time when the event occurred."},
			// CloudTrail
			{Name: "event_time", Type: proto.ColumnType_TIMESTAMP, Hydrate: getCloudtrailMessageField, Description: "The date and time the request was made, in coordinated universal time (UTC)."},
			{Name: "event_version", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: "The version of the log event format."},
			{Name: "event_source", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: "The AWS service that the request was made to."},
			{Name: "event_name", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: "The name of the event returned."},
			{Name: "read_only", Type: proto.ColumnType_BOOL, Hydrate: getCloudtrailMessageField, Description: "Information about whether the event is a write event or a read event."},
			{Name: "access_key_id", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Transform: transform.FromField("AccessKeyId"), Description: "The AWS access key ID that was used to sign the request. If the request was made with temporary security credentials, this is the access key ID of the temporary credentials."},
			{Name: "cloudtrail_event", Type: proto.ColumnType_JSON, Hydrate: getCloudtrailMessageField, Transform: transform.FromField("CloudTrailEvent"), Description: "The CloudTrail event."},
			{Name: "resources", Type: proto.ColumnType_JSON, Hydrate: getCloudtrailMessageField, Description: "A list of resources referenced by the event returned."},
			{Name: "aws_region", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "source_ip_address", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "event_category", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "event_type", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "request_id", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "shared_event_id", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "recipient_account_id", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "vpc_endpoint_id", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "user_agent", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "user_identity", Type: proto.ColumnType_JSON, Hydrate: getCloudtrailMessageField, Description: "Information about the user that made a request."},
			{Name: "request_parameters", Type: proto.ColumnType_JSON, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "response_elements", Type: proto.ColumnType_JSON, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "additional_event_data", Type: proto.ColumnType_JSON, Hydrate: getCloudtrailMessageField, Description: ""},
			{Name: "tls_details", Type: proto.ColumnType_JSON, Hydrate: getCloudtrailMessageField, Description: ""},
			// Other columns
			{Name: "parsed_message", Description: "", Type: proto.ColumnType_STRING, Hydrate: getCloudtrailMessageField, Transform: transform.FromValue()},
			{Name: "message", Description: "", Type: proto.ColumnType_STRING, Transform: transform.FromField("Message")},
			{Name: "event_id", Description: "The ID of the event.", Type: proto.ColumnType_STRING, Transform: transform.FromField("EventId")},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter pattern for the search."},
			{Name: "ingestion_time", Description: "The time when the event was ingested.", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("IngestionTime").Transform(transform.UnixMsToTimestamp)},
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
