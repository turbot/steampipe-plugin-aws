package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/memorydb"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsMemoryDBReservedNodes(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_memorydb_reserved_nodes",
		Description: "AWS MemoryDB Reserved Nodes",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("reservation_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ReservedNodeNotFoundFault"}),
			},
			Hydrate: getMemoryDBReservedNode,
			Tags:    map[string]string{"service": "memorydb", "action": "DescribeReservedNodes"},
		},
		List: &plugin.ListConfig{
			Hydrate: listMemoryDBReservedNodes,
			Tags:    map[string]string{"service": "memorydb", "action": "DescribeReservedNodes"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "reservation_id", Require: plugin.Optional},
				{Name: "reserved_nodes_offering_id", Require: plugin.Optional},
				{Name: "node_type", Require: plugin.Optional},
				{Name: "duration", Require: plugin.Optional},
				{Name: "offering_type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_MEMORY_DB_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "reservation_id",
				Description: "A customer-specified identifier to track this reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the reserved node.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "reserved_nodes_offering_id",
				Description: "The ID of the reserved node offering to purchase.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the reserved node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_type",
				Description: "The node type for the reserved nodes.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_count",
				Description: "The number of nodes that have been reserved.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "duration",
				Description: "The duration of the reservation in seconds.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "fixed_price",
				Description: "The fixed price charged for this reserved node.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "offering_type",
				Description: "The offering type of this reserved node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_time",
				Description: "The time the reservation started.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "recurring_charges",
				Description: "The recurring price charged to run this reserved node.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReservationId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listMemoryDBReservedNodes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := MemoryDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_memorydb_reserved_nodes.listMemoryDBReservedNodes", "get_client_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &memorydb.DescribeReservedNodesInput{
		MaxResults: aws.Int32(100),
	}

	if d.EqualsQuals["reservation_id"] != nil {
		input.ReservationId = aws.String(d.EqualsQuals["reservation_id"].GetStringValue())
	}

	if d.EqualsQuals["reserved_nodes_offering_id"] != nil {
		input.ReservedNodesOfferingId = aws.String(d.EqualsQuals["reserved_nodes_offering_id"].GetStringValue())
	}

	if d.EqualsQuals["node_type"] != nil {
		input.NodeType = aws.String(d.EqualsQuals["node_type"].GetStringValue())
	}

	if d.EqualsQuals["duration"] != nil {
		input.Duration = aws.String(fmt.Sprintf("%v", d.EqualsQuals["duration"].GetInt64Value()))
	}

	if d.EqualsQuals["offering_type"] != nil {
		input.OfferingType = aws.String(d.EqualsQuals["offering_type"].GetStringValue())
	}

	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxResults {
			if limit < 20 {
				input.MaxResults = aws.Int32(20)
			} else {
				input.MaxResults = aws.Int32(limit)
			}
		}
	}

	// List call
	paginator := memorydb.NewDescribeReservedNodesPaginator(svc, input, func(o *memorydb.DescribeReservedNodesPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_memorydb_reserved_nodes.listMemoryDBReservedNodes", "api_error", err)
			return nil, err
		}

		for _, reservedNode := range output.ReservedNodes {
			d.StreamListItem(ctx, reservedNode)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getMemoryDBReservedNode(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	reservationId := quals["reservation_id"].GetStringValue()

	// check if reservationId is empty
	if reservationId == "" {
		return nil, nil
	}

	// Create service
	svc, err := MemoryDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_memorydb_reserved_nodes.getMemoryDBReservedNode", "get_client_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &memorydb.DescribeReservedNodesInput{
		ReservationId: aws.String(reservationId),
	}

	op, err := svc.DescribeReservedNodes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_memorydb_reserved_nodes.getMemoryDBReservedNode", "api_error", err)
		return nil, err
	}

	if len(op.ReservedNodes) > 0 {
		return op.ReservedNodes[0], nil
	}
	return nil, nil
}
