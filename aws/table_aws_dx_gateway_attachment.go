package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directconnect"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsDxGatewayAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dx_gateway_attachment",
		Description: "AWS Direct Connect Gateway Attachment",
		List: &plugin.ListConfig{
			Hydrate: listDxGatewayAttachments,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeDirectConnectGatewayAttachments"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "direct_connect_gateway_id", Require: plugin.Optional},
				{Name: "virtual_interface_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DIRECTCONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "direct_connect_gateway_id",
				Description: "The ID of the Direct Connect gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_interface_id",
				Description: "The ID of the virtual interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_interface_type",
				Description: "The type of virtual interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_interface_name",
				Description: "The name of the virtual interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_interface_owner_account",
				Description: "The ID of the AWS account that owns the virtual interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_interface_region",
				Description: "The AWS Region where the virtual interface is located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attachment_state",
				Description: "The state of the attachment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attachment_type",
				Description: "The type of attachment.",
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
				Transform:   transform.FromField("VirtualInterfaceName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listDxGatewayAttachments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_gateway_attachment.listDxGatewayAttachments", "connection_error", err)
		return nil, err
	}

	input := &directconnect.DescribeDirectConnectGatewayAttachmentsInput{}

	// Add optional filter parameters
	if d.EqualsQuals["direct_connect_gateway_id"] != nil {
		input.DirectConnectGatewayId = aws.String(d.EqualsQuals["direct_connect_gateway_id"].GetStringValue())
	}
	if d.EqualsQuals["virtual_interface_id"] != nil {
		input.VirtualInterfaceId = aws.String(d.EqualsQuals["virtual_interface_id"].GetStringValue())
	}

	// Limiting the results
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < 100 {
			input.MaxResults = aws.Int32(limit)
		}
	}

	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := svc.DescribeDirectConnectGatewayAttachments(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dx_gateway_attachment.listDxGatewayAttachments", "api_error", err)
			return nil, err
		}

		for _, item := range output.DirectConnectGatewayAttachments {
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
