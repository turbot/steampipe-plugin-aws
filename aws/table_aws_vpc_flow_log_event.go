package aws

import (
	"context"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	cloudwatchlogsv1 "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsVpcFlowLogEventListKeyColumns() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{
		{Name: "log_group_name"},
		{Name: "log_stream_name", Require: plugin.Optional},
		{Name: "filter", Require: plugin.Optional, CacheMatch: "exact"},
		{Name: "region", Require: plugin.Optional},
		{Name: "timestamp", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},

		// others
		{Name: "event_id", Require: plugin.Optional},
		{Name: "interface_id", Require: plugin.Optional},
		{Name: "src_addr", Require: plugin.Optional},
		{Name: "dst_addr", Require: plugin.Optional},
		{Name: "src_port", Require: plugin.Optional},
		{Name: "dst_port", Require: plugin.Optional},
		{Name: "action", Require: plugin.Optional},
		{Name: "log_status", Require: plugin.Optional},
	}
}

//// TABLE DEFINITION

func tableAwsVpcFlowLogEvent(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_flow_log_event",
		Description: "AWS VPC Flow Log events from CloudWatch Logs",
		List: &plugin.ListConfig{
			Hydrate:    listCloudwatchLogEvents,
			Tags:       map[string]string{"service": "logs", "action": "FilterLogEvents"},
			KeyColumns: tableAwsVpcFlowLogEventListKeyColumns(),
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchlogsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			// Top columns
			{Name: "log_group_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("log_group_name"), Description: "The name of the log group to which this event belongs."},
			{Name: "log_stream_name", Type: proto.ColumnType_STRING, Description: "The name of the log stream to which this event belongs."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Timestamp").Transform(transform.UnixMsToTimestamp), Description: "The time when the event occurred."},
			{Name: "version", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 0), Description: "The VPC Flow Logs version. If you use the default format, the version is 2. If you use a custom format, the version is the highest version among the specified fields. For example, if you specify only fields from version 2, the version is 2. If you specify a mixture of fields from versions 2, 3, and 4, the version is 4."},
			{Name: "interface_account_id", Type: proto.ColumnType_STRING, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 1), Description: "The AWS account ID of the owner of the source network interface for which traffic is recorded. If the network interface is created by an AWS service, for example when creating a VPC endpoint or Network Load Balancer, the record may display unknown for this field."},
			{Name: "interface_id", Type: proto.ColumnType_STRING, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 2), Description: "The ID of the network interface for which the traffic is recorded."},
			{Name: "src_addr", Type: proto.ColumnType_IPADDR, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 3), Description: "The source address for incoming traffic, or the IPv4 or IPv6 address of the network interface for outgoing traffic on the network interface. The IPv4 address of the network interface is always its private IPv4 address. See also pkt-srcaddr."},
			{Name: "dst_addr", Type: proto.ColumnType_IPADDR, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 4), Description: "The destination address for outgoing traffic, or the IPv4 or IPv6 address of the network interface for incoming traffic on the network interface. The IPv4 address of the network interface is always its private IPv4 address. See also pkt-dstaddr."},
			{Name: "src_port", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 5), Description: "The source port of the traffic."},
			{Name: "dst_port", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 6), Description: "The destination port of the traffic."},
			{Name: "protocol", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 7), Description: "The IANA protocol number of the traffic. For more information, see Assigned Internet Protocol Numbers."},
			{Name: "packets", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 8), Description: "The number of packets transferred during the flow."},
			{Name: "bytes", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 9), Description: "The number of bytes transferred during the flow."},
			{Name: "start", Type: proto.ColumnType_TIMESTAMP, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 10).Transform(transform.UnixToTimestamp), Description: "The time when the first packet of the flow was received within the aggregation interval. This might be up to 60 seconds after the packet was transmitted or received on the network interface."},
			{Name: "end", Type: proto.ColumnType_TIMESTAMP, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 11).Transform(transform.UnixToTimestamp), Description: "The time when the last packet of the flow was received within the aggregation interval. This might be up to 60 seconds after the packet was transmitted or received on the network interface."},
			{Name: "action", Type: proto.ColumnType_STRING, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 12), Description: "The action that is associated with the traffic: ACCEPT — The recorded traffic was permitted by the security groups and network ACLs. REJECT — The recorded traffic was not permitted by the security groups or network ACLs."},
			{Name: "log_status", Type: proto.ColumnType_STRING, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 13), Description: "The logging status of the flow log: OK — Data is logging normally to the chosen destinations. NODATA — There was no network traffic to or from the network interface during the aggregation interval. SKIPDATA — Some flow log records were skipped during the aggregation interval. This may be because of an internal capacity constraint, or an internal error."},
			// Other columns
			{Name: "event_id", Description: "The ID of the event.", Type: proto.ColumnType_STRING, Transform: transform.FromField("EventId")},
			{Name: "filter", Description: "Filter pattern for the search.", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter")},
			{Name: "ingestion_time", Description: "The time when the event was ingested.", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("IngestionTime").Transform(transform.UnixMsToTimestamp)},
		}),
	}
}

func getMessageField(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	e := h.Item.(types.FilteredLogEvent)
	fields := strings.Fields(*e.Message)
	return fields, nil
}

func getField(_ context.Context, d *transform.TransformData) (interface{}, error) {
	fields := d.Value.([]string)
	idx := d.Param.(int)
	if fields[idx] == "-" {
		return nil, nil
	}
	return fields[idx], nil
}

func buildFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := []string{"action", "log_status", "interface_id", "event_id", "src_addr", "dst_addr", "src_port", "dst_port"}

	for _, qual := range filterQuals {
		switch qual {
		case "action", "log_status", "interface_id", "event_id":
			if equalQuals[qual] != nil {
				filters = append(filters, equalQuals[qual].GetStringValue())
			}
		case "src_addr", "dst_addr":
			if equalQuals[qual] != nil {
				filters = append(filters, equalQuals[qual].GetInetValue().Addr)
			}
		case "src_port", "dst_port":
			if equalQuals[qual] != nil {
				filters = append(filters, strconv.Itoa(int(equalQuals[qual].GetInt64Value())))
			}
		}
	}

	return filters
}
