package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directconnect"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsDxGatewayAssociationProposal(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dx_gateway_association_proposal",
		Description: "AWS Direct Connect Gateway Association Proposal",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("proposal_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DirectConnectClientException"}),
			},
			Hydrate: getDxGatewayAssociationProposal,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeDirectConnectGatewayAssociationProposals"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDxGatewayAssociationProposals,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeDirectConnectGatewayAssociationProposals"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "direct_connect_gateway_id", Require: plugin.Optional},
				{Name: "associated_gateway_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DIRECTCONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "proposal_id",
				Description: "The ID of the association proposal.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "direct_connect_gateway_id",
				Description: "The ID of the Direct Connect gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "direct_connect_gateway_owner_account",
				Description: "The ID of the AWS account that owns the Direct Connect gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "proposal_state",
				Description: "The state of the proposal.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associated_gateway_id",
				Description: "The ID of the associated gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssociatedGateway.Id"),
			},
			{
				Name:        "associated_gateway_type",
				Description: "The type of associated gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssociatedGateway.Type"),
			},
			{
				Name:        "associated_gateway_owner_account",
				Description: "The ID of the AWS account that owns the associated virtual private gateway or transit gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssociatedGateway.OwnerAccount"),
			},
			{
				Name:        "associated_gateway_region",
				Description: "The Region where the associated gateway is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssociatedGateway.Region"),
			},
			{
				Name:        "existing_allowed_prefixes_to_direct_connect_gateway",
				Description: "The existing Amazon VPC prefixes advertised to the Direct Connect gateway.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "requested_allowed_prefixes_to_direct_connect_gateway",
				Description: "The Amazon VPC prefixes to advertise to the Direct Connect gateway.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProposalId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listDxGatewayAssociationProposals(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_gateway_association_proposal.listDxGatewayAssociationProposals", "connection_error", err)
		return nil, err
	}

	input := &directconnect.DescribeDirectConnectGatewayAssociationProposalsInput{}

	// Add optional filter parameters
	if d.EqualsQuals["direct_connect_gateway_id"] != nil {
		input.DirectConnectGatewayId = aws.String(d.EqualsQuals["direct_connect_gateway_id"].GetStringValue())
	}
	if d.EqualsQuals["associated_gateway_id"] != nil {
		input.AssociatedGatewayId = aws.String(d.EqualsQuals["associated_gateway_id"].GetStringValue())
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

		output, err := svc.DescribeDirectConnectGatewayAssociationProposals(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dx_gateway_association_proposal.listDxGatewayAssociationProposals", "api_error", err)
			return nil, err
		}

		for _, item := range output.DirectConnectGatewayAssociationProposals {
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

func getDxGatewayAssociationProposal(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	proposalID := d.EqualsQuals["proposal_id"].GetStringValue()

	// get service
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_gateway_association_proposal.getDxGatewayAssociationProposal", "connection_error", err)
		return nil, err
	}

	params := &directconnect.DescribeDirectConnectGatewayAssociationProposalsInput{
		ProposalId: aws.String(proposalID),
	}

	op, err := svc.DescribeDirectConnectGatewayAssociationProposals(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_gateway_association_proposal.getDxGatewayAssociationProposal", "api_error", err)
		return nil, err
	}

	if len(op.DirectConnectGatewayAssociationProposals) > 0 {
		return op.DirectConnectGatewayAssociationProposals[0], nil
	}
	return nil, nil
}
