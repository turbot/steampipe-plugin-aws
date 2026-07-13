package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	"github.com/aws/aws-sdk-go-v2/service/directconnect/types"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsDirectConnectConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_directconnect_connection",
		Description: "AWS Direct Connect Connection",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("connection_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DirectConnectClientException"}),
			},
			Hydrate: getDirectConnectConnection,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeConnections"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDirectConnectConnections,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeConnections"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DIRECTCONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "connection_id",
				Description: "The ID of the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connection_name",
				Description: "The name of the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the connection.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDirectConnectConnectionARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "connection_state",
				Description: "The state of the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bandwidth",
				Description: "The bandwidth of the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location of the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_account",
				Description: "The ID of the Amazon Web Services account that owns the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "partner_name",
				Description: "The name of the Direct Connect service provider associated with the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provider_name",
				Description: "The name of the service provider associated with the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "aws_device_v2",
				Description: "The Direct Connect endpoint that terminates the physical connection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AwsDeviceV2"),
			},
			{
				Name:        "aws_logical_device_id",
				Description: "The Direct Connect endpoint that terminates the logical connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encryption_mode",
				Description: "The MAC Security (MACsec) connection encryption mode.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "has_logical_redundancy",
				Description: "Indicates whether the connection supports a secondary BGP peer in the same address family (IPv4/IPv6).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "jumbo_frame_capable",
				Description: "Indicates whether jumbo frames are supported.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "lag_id",
				Description: "The ID of the LAG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "loa_issue_time",
				Description: "The time of the most recent call to DescribeLoa for this connection.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "mac_sec_capable",
				Description: "Indicates whether the connection supports MAC Security (MACsec).",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "mac_sec_keys",
				Description: "The MAC Security (MACsec) security keys associated with the connection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "port_encryption_status",
				Description: "The MAC Security (MACsec) port link status of the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vlan",
				Description: "The ID of the VLAN.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the connection.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(directConnectConnectionTagsToTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectionName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDirectConnectConnectionARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDirectConnectConnections(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directconnect_connection.listDirectConnectConnections", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &directconnect.DescribeConnectionsInput{}

	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := svc.DescribeConnections(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_directconnect_connection.listDirectConnectConnections", "api_error", err)
			return nil, err
		}

		for _, connection := range output.Connections {
			d.StreamListItem(ctx, connection)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if output.NextToken == nil {
			break
		}
		input.NextToken = output.NextToken
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDirectConnectConnection(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	connectionID := d.EqualsQualString("connection_id")

	// Empty check
	if connectionID == "" {
		return nil, nil
	}

	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directconnect_connection.getDirectConnectConnection", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	params := &directconnect.DescribeConnectionsInput{
		ConnectionId: &connectionID,
	}

	output, err := svc.DescribeConnections(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directconnect_connection.getDirectConnectConnection", "api_error", err)
		return nil, err
	}

	if len(output.Connections) > 0 {
		return output.Connections[0], nil
	}

	return nil, nil
}

func getDirectConnectConnectionARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	connection := h.Item.(types.Connection)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directconnect_connection.getDirectConnectConnectionARN", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// arn:${Partition}:directconnect:${Region}:${Account}:dxcon/${ConnectionId}
	arn := "arn:" + commonColumnData.Partition + ":directconnect:" + region + ":" + commonColumnData.AccountId + ":dxcon/" + *connection.ConnectionId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func directConnectConnectionTagsToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.Value.([]types.Tag)

	if tags == nil {
		return nil, nil
	}

	turbotTags := make(map[string]string)
	for _, tag := range tags {
		if tag.Key != nil && tag.Value != nil {
			turbotTags[*tag.Key] = *tag.Value
		}
	}
	return turbotTags, nil
}
