package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	"github.com/aws/aws-sdk-go-v2/service/directconnect/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Custom struct to represent customer metadata with individual agreement rows
type DxCustomerMetadataResult struct {
	NniPartnerType types.NniPartnerType `json:"nni_partner_type"`
	AgreementName  *string              `json:"agreement_name"`
	Status         *string              `json:"status"`
}

func tableAwsDxCustomerMetadata(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dx_customer_metadata",
		Description: "AWS Direct Connect Customer Metadata",
		List: &plugin.ListConfig{
			Hydrate: listDxCustomerMetadata,
			Tags:    map[string]string{"service": "directconnect", "action": "DescribeCustomerMetadata"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DIRECTCONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "nni_partner_type",
				Description: "The type of network-to-network interface (NNI) partner. Can be V1, V2, or nonPartner.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NniPartnerType"),
			},
			{
				Name:        "agreement_name",
				Description: "The name of the customer agreement.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AgreementName"),
			},
			{
				Name:        "status",
				Description: "The status of the customer agreement.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AgreementName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listDxCustomerMetadata(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := DirectConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_customer_metadata.listDxCustomerMetadata", "connection_error", err)
		return nil, err
	}

	input := &directconnect.DescribeCustomerMetadataInput{}

	// Execute call
	output, err := svc.DescribeCustomerMetadata(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dx_customer_metadata.listDxCustomerMetadata", "api_error", err)
		return nil, err
	}

	// Stream each agreement as a separate row
	for _, agreement := range output.Agreements {
		result := DxCustomerMetadataResult{
			NniPartnerType: output.NniPartnerType,
			AgreementName:  agreement.AgreementName,
			Status:         agreement.Status,
		}

		d.StreamListItem(ctx, result)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
