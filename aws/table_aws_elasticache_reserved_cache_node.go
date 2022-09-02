package aws

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsElastiCacheReservedCacheNode(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_reserved_cache_node",
		Description: "AWS ElastiCache Reserved Cache Node",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("reserved_cache_node_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ReservedCacheNodeNotFound"}),
			},
			Hydrate: getElastiCacheReservedCacheNode,
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheReservedCacheNodes,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cache_node_type", Require: plugin.Optional},
				{Name: "duration", Require: plugin.Optional},
				{Name: "offering_type", Require: plugin.Optional},
				{Name: "reserved_cache_nodes_offering_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "reserved_cache_node_id",
				Description: "The unique identifier for the reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the reserved cache node.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReservationARN"),
			},
			{
				Name:        "reserved_cache_nodes_offering_id",
				Description: "The offering identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the reserved cache node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cache_node_type",
				Description: "The cache node type for the reserved cache nodes.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cache_node_count",
				Description: "The number of cache nodes that have been reserved.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "duration",
				Description: "The duration of the reservation in seconds.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "fixed_price",
				Description: "The fixed price charged for this reserved cache node.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "offering_type",
				Description: "The offering type of this reserved cache node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "product_description",
				Description: "The description of the reserved cache node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_time",
				Description: "The time the reservation started.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "usage_price",
				Description: "The hourly price charged for this reserved cache node.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "recurring_charges",
				Description: "The recurring price charged to run this reserved cache node.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReservedCacheNodeId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ReservationARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listElastiCacheReservedCacheNodes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ElastiCacheService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &elasticache.DescribeReservedCacheNodesInput{
		MaxRecords: aws.Int64(100),
	}

	if d.KeyColumnQuals["cache_node_type"] != nil {
		input.CacheNodeType = aws.String(d.KeyColumnQuals["cache_node_type"].GetStringValue())
	}

	if d.KeyColumnQuals["duration"] != nil {
		input.Duration = aws.String(fmt.Sprintf("%v", d.KeyColumnQuals["duration"].GetInt64Value()))
	}

	if d.KeyColumnQuals["offering_type"] != nil {
		input.OfferingType = aws.String(d.KeyColumnQuals["offering_type"].GetStringValue())
	}

	if d.KeyColumnQuals["reserved_cache_nodes_offering_id"] != nil {
		input.ReservedCacheNodesOfferingId = aws.String(d.KeyColumnQuals["reserved_cache_nodes_offering_id"].GetStringValue())
	}
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxRecords {
			if *limit < 20 {
				input.MaxRecords = aws.Int64(20)
			} else {
				input.MaxRecords = limit
			}
		}
	}

	// List call
	err = svc.DescribeReservedCacheNodesPages(
		input,
		func(page *elasticache.DescribeReservedCacheNodesOutput, isLast bool) bool {
			for _, reservedCacheNode := range page.ReservedCacheNodes {
				d.StreamListItem(ctx, reservedCacheNode)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listElastiCacheReservedCacheNodes", "DescribeReservedCacheNodesPages", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getElastiCacheReservedCacheNode(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	reservedCacheNodeId := quals["reserved_cache_node_id"].GetStringValue()

	// check if reservedCacheNodeId is empty
	if reservedCacheNodeId == "" {
		return nil, nil
	}

	// Create service
	svc, err := ElastiCacheService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &elasticache.DescribeReservedCacheNodesInput{
		ReservedCacheNodeId: aws.String(reservedCacheNodeId),
	}

	op, err := svc.DescribeReservedCacheNodes(params)
	if err != nil {
		plugin.Logger(ctx).Error("getElastiCacheReservedCacheNode", "DescribeReservedCacheNodes", err)
		return nil, err
	}

	if len(op.ReservedCacheNodes) > 0 {
		return op.ReservedCacheNodes[0], nil
	}
	return nil, nil
}
