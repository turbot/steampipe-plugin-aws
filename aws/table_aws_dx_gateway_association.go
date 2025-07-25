package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directconnect"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsDxGatewayAssociation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dx_gateway_association",
		Description: "AWS Direct Connect Gateway Association",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"direct_connect_gateway_id", "associated_gateway_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DirectConnectClientException"}),
			},
			Hydrate: getDxGatewayAssociation,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeDirectConnectGatewayAssociations"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDxGatewayAssociations,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeDirectConnectGatewayAssociations"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "direct_connect_gateway_id", Require: plugin.Optional},
				{Name: "associated_gateway_id", Require: plugin.Optional},
				{Name: "association_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DIRECTCONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "association_id",
				Description: "The ID of the Direct Connect gateway association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "direct_connect_gateway_id",
				Description: "The ID of the Direct Connect gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associated_gateway_id",
				Description: "The ID of the associated gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associated_gateway_type",
				Description: "The type of associated gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associated_gateway_owner_account",
				Description: "The ID of the AWS account that owns the associated virtual private gateway or transit gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associated_gateway_region",
				Description: "The Region where the associated gateway is located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "direct_connect_gateway_owner_account",
				Description: "The ID of the AWS account that owns the Direct Connect gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "association_state",
				Description: "The state of the association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_change_error",
				Description: "The error message if the state of an object failed to advance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allowed_prefixes_to_direct_connect_gateway",
				Description: "The Amazon VPC prefixes to advertise to the Direct Connect gateway.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "virtual_gateway_id",
				Description: "The ID of the virtual private gateway. Applies only to private virtual interfaces.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_gateway_owner_account",
				Description: "The ID of the AWS account that owns the virtual private gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_gateway_region",
				Description: "The AWS Region where the virtual private gateway is located.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssociationId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listDxGatewayAssociations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_gateway_association.listDxGatewayAssociations", "connection_error", err)
		return nil, err
	}

	input := &directconnect.DescribeDirectConnectGatewayAssociationsInput{}

	// Additional filters
	if d.EqualsQuals["direct_connect_gateway_id"] != nil {
		input.DirectConnectGatewayId = aws.String(d.EqualsQuals["direct_connect_gateway_id"].GetStringValue())
	}
	if d.EqualsQuals["associated_gateway_id"] != nil {
		input.AssociatedGatewayId = aws.String(d.EqualsQuals["associated_gateway_id"].GetStringValue())
	}
	if d.EqualsQuals["association_id"] != nil {
		input.AssociationId = aws.String(d.EqualsQuals["association_id"].GetStringValue())
	}

	// Limiting the results
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		input.MaxResults = aws.Int32(limit)
	}

	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := svc.DescribeDirectConnectGatewayAssociations(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dx_gateway_association.listDxGatewayAssociations", "api_error", err)
			return nil, err
		}

		for _, item := range output.DirectConnectGatewayAssociations {
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

func getDxGatewayAssociation(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	gatewayID := d.EqualsQuals["direct_connect_gateway_id"].GetStringValue()
	associatedGatewayID := d.EqualsQuals["associated_gateway_id"].GetStringValue()

	// get service
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_gateway_association.getDxGatewayAssociation", "connection_error", err)
		return nil, err
	}

	params := &directconnect.DescribeDirectConnectGatewayAssociationsInput{
		DirectConnectGatewayId: aws.String(gatewayID),
		AssociatedGatewayId:    aws.String(associatedGatewayID),
	}

	op, err := svc.DescribeDirectConnectGatewayAssociations(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_gateway_association.getDxGatewayAssociation", "api_error", err)
		return nil, err
	}

	if len(op.DirectConnectGatewayAssociations) > 0 {
		return op.DirectConnectGatewayAssociations[0], nil
	}
	return nil, nil
}
