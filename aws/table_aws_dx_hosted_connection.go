package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	"github.com/aws/aws-sdk-go-v2/service/directconnect/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsDxHostedConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dx_hosted_connection",
		Description: "AWS Direct Connect Hosted Connection",
		List: &plugin.ListConfig{
			ParentHydrate: listDxConnections,
			Hydrate:       listDxHostedConnections,
			Tags:          map[string]string{"service": "directconnect", "action": "DescribeHostedConnections"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DIRECTCONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "connection_id",
				Description: "The ID of the hosted connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connection_name",
				Description: "The name of the hosted connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connection_state",
				Description: "The state of the hosted connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location of the hosted connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bandwidth",
				Description: "The bandwidth of the hosted connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vlan",
				Description: "The ID of the VLAN.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "partner_name",
				Description: "The name of the AWS Direct Connect service provider associated with the hosted connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lag_id",
				Description: "The ID of the LAG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "aws_device",
				Description: "The Direct Connect endpoint on which the physical connection terminates.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "aws_device_v2",
				Description: "The Direct Connect endpoint that terminates the logical connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "aws_logical_device_id",
				Description: "The Direct Connect endpoint that terminates the logical connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "has_logical_redundancy",
				Description: "Indicates whether the connection supports a secondary BGP peer in the same address family (IPv4/IPv6).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "jumbo_frame_capable",
				Description: "Indicates whether jumbo frames (9001 MTU) are supported.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "macsec_capable",
				Description: "Indicates whether the connection supports MAC Security (MACsec).",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "port_encryption_status",
				Description: "The MAC Security (MACsec) port link status of the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encryption_mode",
				Description: "The MAC Security (MACsec) connection encryption mode.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provider_name",
				Description: "The name of the service provider associated with the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_account",
				Description: "The ID of the AWS account that owns the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "loa_issue_time",
				Description: "The time of the most recent call to DescribeConnectionLoa for this connection.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "macsec_keys",
				Description: "The MAC Security (MACsec) security keys associated with the connection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the hosted connection.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectionName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getDxHostedConnectionTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDxHostedConnectionARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}


//// LIST FUNCTION

func listDxHostedConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	connection := h.Item.(types.Connection)

	// Create session
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_hosted_connection.listDxHostedConnections", "connection_error", err)
		return nil, err
	}

	input := &directconnect.DescribeHostedConnectionsInput{
		ConnectionId: connection.ConnectionId,
	}

	// Execute list call
	// apply rate limiting
	d.WaitForListRateLimit(ctx)

	output, err := svc.DescribeHostedConnections(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_hosted_connection.listDxHostedConnections", "api_error", err)
		return nil, err
	}

	for _, item := range output.Connections {
		d.StreamListItem(ctx, item)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDxHostedConnectionARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	connection := h.Item.(types.Connection)
	region := d.EqualsQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":directconnect:" + region + ":" + *connection.OwnerAccount + ":dxcon/" + *connection.ConnectionId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getDxHostedConnectionTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	connection := d.HydrateItem.(types.Connection)

	if connection.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range connection.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
