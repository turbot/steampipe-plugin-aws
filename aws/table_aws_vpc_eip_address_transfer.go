package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsVpcEipAddressTransfer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_eip_address_transfer",
		Description: "AWS VPC Elastic IP Address Transfer",
		List: &plugin.ListConfig{
			ParentHydrate: listVpcEips,
			Hydrate:       listVpcEipAddressTransfers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "allocation_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "allocation_id",
				Description: "The allocation ID of an Elastic IP address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "address_transfer_status",
				Description: "The Elastic IP address transfer status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_ip",
				Description: "The Elastic IP address being transferred.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "transfer_account_id",
				Description: "The ID of the account that you want to transfer the Elastic IP address to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "transfer_offer_accepted_timestamp",
				Description: "The timestamp when the Elastic IP address transfer was accepted.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "transfer_offer_expiration_timestamp",
				Description: "The timestamp when the Elastic IP address transfer expired.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AllocationId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcEipAddressTransfers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	address := h.Item.(types.Address)
	allocationId := d.EqualsQualString("allocation_id")

	// Avoid api call if user specified allocation_id doesn't match the hydrate id.
	if allocationId != "" && allocationId != *address.AllocationId {
		return nil, nil
	}

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_eip_address_transfer.listVpcEipAddressTransfers", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	// According to the docs MaxResult can be set between 5-1000, but the API throws error if we pass value >10
	// api error InvalidMaxResults: Value ( 1000 ) for parameter maxResults is invalid. Expecting a value less than or equal to 10.
	maxLimit := int32(10)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				maxLimit = int32(5)
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeAddressTransfersInput{
		AllocationIds: []string{*address.AllocationId},
		MaxResults:    &maxLimit,
	}

	paginator := ec2.NewDescribeAddressTransfersPaginator(svc, input, func(o *ec2.DescribeAddressTransfersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_eip_address_transfer.listVpcEipAddressTransfers", "api_error", err)
			return nil, err
		}

		for _, items := range output.AddressTransfers {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}
