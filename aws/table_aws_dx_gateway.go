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

func tableAwsDxGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dx_gateway",
		Description: "AWS Direct Connect Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("direct_connect_gateway_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DirectConnectClientException"}),
			},
			Hydrate: getDxGateway,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeDirectConnectGateways"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDxGateways,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeDirectConnectGateways"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DIRECTCONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "direct_connect_gateway_id",
				Description: "The ID of the Direct Connect gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "direct_connect_gateway_name",
				Description: "The name of the Direct Connect gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "direct_connect_gateway_state",
				Description: "The state of the Direct Connect gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "amazon_side_asn",
				Description: "The autonomous system number (ASN) for the Amazon side of the connection.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "owner_account",
				Description: "The ID of the AWS account that owns the Direct Connect gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_change_error",
				Description: "The error message if the state of an object failed to advance.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DirectConnectGatewayName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDxGatewayARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDxGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_gateway.listDxGateways", "connection_error", err)
		return nil, err
	}

	input := &directconnect.DescribeDirectConnectGatewaysInput{}

	// Limiting the results
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		input.MaxResults = aws.Int32(limit)
	}

	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)
		
		// Execute list call
		output, err := svc.DescribeDirectConnectGateways(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dx_gateway.listDxGateways", "api_error", err)
			return nil, err
		}
	
		for _, item := range output.DirectConnectGateways {
			d.StreamListItem(ctx, item)
	
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

func getDxGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	gatewayID := d.EqualsQuals["direct_connect_gateway_id"].GetStringValue()

	// get service
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_gateway.getDxGateway", "connection_error", err)
		return nil, err
	}

	params := &directconnect.DescribeDirectConnectGatewaysInput{
		DirectConnectGatewayId: aws.String(gatewayID),
	}

	op, err := svc.DescribeDirectConnectGateways(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_gateway.getDxGateway", "api_error", err)
		return nil, err
	}

	if len(op.DirectConnectGateways) > 0 {
		return op.DirectConnectGateways[0], nil
	}
	return nil, nil
}

func getDxGatewayARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	gateway := h.Item.(types.DirectConnectGateway)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Direct Connect Gateway is a global resource, so no region in ARN
	arn := "arn:" + commonColumnData.Partition + ":directconnect::" + *gateway.OwnerAccount + ":dx-gateway/" + *gateway.DirectConnectGatewayId

	return arn, nil
}
