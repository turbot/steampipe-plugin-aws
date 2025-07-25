package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	"github.com/aws/aws-sdk-go-v2/service/directconnect/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsDxLag(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dx_lag",
		Description: "AWS Direct Connect LAG (Link Aggregation Group)",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("lag_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DirectConnectClientException"}),
			},
			Hydrate: getDxLag,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeLags"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDxLags,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeLags"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DIRECTCONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "lag_id",
				Description: "The ID of the LAG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lag_name",
				Description: "The name of the LAG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lag_state",
				Description: "The state of the LAG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location of the LAG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "minimum_links",
				Description: "The minimum number of physical dedicated connections that must be operational for the LAG itself to be operational.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "connections_bandwidth",
				Description: "The individual bandwidth of the physical connections bundled by the LAG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "number_of_connections",
				Description: "The number of physical dedicated connections bundled by the LAG, up to a maximum of 10.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "allows_hosted_connections",
				Description: "Indicates whether the LAG can host other connections.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "jumbo_frame_capable",
				Description: "Indicates whether jumbo frames (9001 MTU) are supported.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "has_logical_redundancy",
				Description: "Indicates whether the LAG supports a secondary BGP peer in the same address family (IPv4/IPv6).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "aws_device",
				Description: "The AWS Direct Connect endpoint that hosts the LAG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "aws_device_v2",
				Description: "The AWS Direct Connect endpoint that hosts the LAG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "aws_logical_device_id",
				Description: "The AWS Direct Connect endpoint that terminates the logical connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_account",
				Description: "The ID of the AWS account that owns the LAG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provider_name",
				Description: "The name of the service provider associated with the LAG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "macsec_capable",
				Description: "Indicates whether the LAG supports MAC Security (MACsec).",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "encryption_mode",
				Description: "The MAC Security (MACsec) encryption mode.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "macsec_keys",
				Description: "The MAC Security (MACsec) security keys associated with the LAG.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "connections",
				Description: "The connections bundled by the LAG.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the LAG.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LagName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getDxLagTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDxLagARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDxLags(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_lag.listDxLags", "connection_error", err)
		return nil, err
	}

	input := &directconnect.DescribeLagsInput{}

	// Execute list call
	// apply rate limiting
	d.WaitForListRateLimit(ctx)

	output, err := svc.DescribeLags(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_lag.listDxLags", "api_error", err)
		return nil, err
	}

	for _, item := range output.Lags {
		d.StreamListItem(ctx, item)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDxLag(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	lagID := d.EqualsQuals["lag_id"].GetStringValue()

	// get service
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_lag.getDxLag", "connection_error", err)
		return nil, err
	}

	params := &directconnect.DescribeLagsInput{
		LagId: aws.String(lagID),
	}

	op, err := svc.DescribeLags(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_lag.getDxLag", "api_error", err)
		return nil, err
	}

	if len(op.Lags) > 0 {
		return op.Lags[0], nil
	}
	return nil, nil
}

func getDxLagARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	lag := h.Item.(types.Lag)
	region := d.EqualsQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":directconnect:" + region + ":" + *lag.OwnerAccount + ":dxlag/" + *lag.LagId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getDxLagTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	lag := d.HydrateItem.(types.Lag)

	if lag.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range lag.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
