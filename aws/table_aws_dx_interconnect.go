package aws

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	"github.com/aws/aws-sdk-go-v2/service/directconnect/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsDxInterconnect(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dx_interconnect",
		Description: "AWS Direct Connect Interconnect",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("interconnect_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DirectConnectClientException"}),
			},
			Hydrate: getDxInterconnect,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeInterconnects"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDxInterconnects,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeInterconnects"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DIRECTCONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "interconnect_id",
				Description: "The ID of the interconnect.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "interconnect_name",
				Description: "The name of the interconnect.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "interconnect_state",
				Description: "The state of the interconnect.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location of the interconnect.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bandwidth",
				Description: "The bandwidth of the interconnect.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "loa_issue_time",
				Description: "The time of the most recent call to DescribeInterconnectLoa for this interconnect.",
				Type:        proto.ColumnType_TIMESTAMP,
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
				Description: "Indicates whether the interconnect supports a secondary BGP peer in the same address family (IPv4/IPv6).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "jumbo_frame_capable",
				Description: "Indicates whether jumbo frames (9001 MTU) are supported.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "provider_name",
				Description: "The name of the service provider associated with the interconnect.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "loa_content",
				Description: "The binary contents of the LOA-CFA document (base64 encoded).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDxInterconnectLoa,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the interconnect.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterconnectName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getDxInterconnectTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDxInterconnectARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDxInterconnects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_interconnect.listDxInterconnects", "connection_error", err)
		return nil, err
	}

	input := &directconnect.DescribeInterconnectsInput{}

	// Execute list call
	// apply rate limiting
	d.WaitForListRateLimit(ctx)

	output, err := svc.DescribeInterconnects(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_interconnect.listDxInterconnects", "api_error", err)
		return nil, err
	}

	for _, item := range output.Interconnects {
		d.StreamListItem(ctx, item)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDxInterconnect(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	interconnectID := d.EqualsQuals["interconnect_id"].GetStringValue()

	// get service
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_interconnect.getDxInterconnect", "connection_error", err)
		return nil, err
	}

	params := &directconnect.DescribeInterconnectsInput{
		InterconnectId: aws.String(interconnectID),
	}

	op, err := svc.DescribeInterconnects(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_interconnect.getDxInterconnect", "api_error", err)
		return nil, err
	}

	if len(op.Interconnects) > 0 {
		return op.Interconnects[0], nil
	}
	return nil, nil
}

func getDxInterconnectARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	interconnect := h.Item.(types.Interconnect)
	region := d.EqualsQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Note: Interconnects don't have an owner account field, using the current account
	arn := "arn:" + commonColumnData.Partition + ":directconnect:" + region + ":" + commonColumnData.AccountId + ":dxi/" + *interconnect.InterconnectId

	return arn, nil
}

func getDxInterconnectLoa(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	interconnect := h.Item.(types.Interconnect)

	// get service
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_interconnect.getDxInterconnectLoa", "connection_error", err)
		return nil, err
	}

	params := &directconnect.DescribeLoaInput{
		ConnectionId: interconnect.InterconnectId,
	}

	// Default to PDF content type
	params.LoaContentType = types.LoaContentTypePdf

	op, err := svc.DescribeLoa(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_interconnect.getDxInterconnectLoa", "api_error", err)
		return nil, err
	}

	// Return base64 encoded content similar to connections table
	if len(op.LoaContent) > 0 {
		return base64.StdEncoding.EncodeToString(op.LoaContent), nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getDxInterconnectTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	interconnect := d.HydrateItem.(types.Interconnect)

	if interconnect.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range interconnect.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
