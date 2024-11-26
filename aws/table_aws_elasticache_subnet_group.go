package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"

	elasticachev1 "github.com/aws/aws-sdk-go/service/elasticache"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElastiCacheSubnetGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_subnet_group",
		Description: "AWS ElastiCache Subnet Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cache_subnet_group_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"CacheSubnetGroupNotFoundFault"}),
			},
			Hydrate: getElastiCacheSubnetGroup,
			Tags:    map[string]string{"service": "elasticache", "action": "DescribeCacheSubnetGroups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheSubnetGroups,
			Tags:    map[string]string{"service": "elasticache", "action": "DescribeCacheSubnetGroups"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(elasticachev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cache_subnet_group_name",
				Description: "The name of the cache subnet group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) of the cache subnet group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "cache_subnet_group_description",
				Description: "The description of the cache subnet group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The Amazon Virtual Private Cloud identifier (VPC ID) of the cache subnet group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnets",
				Description: "A list of subnets associated with the cache subnet group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "supported_network_types",
				Description: "Either ipv4 | ipv6 | dual_stack . IPv6 is supported for workloads using Redis engine version 6.2 onward or Memcached engine version 1.6.6 on all instances built on the Nitro system.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CacheSubnetGroupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listElastiCacheSubnetGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_subnet_group.listElastiCacheSubnetGroups", "get_client_error", err)
		return nil, err
	}

	input := &elasticache.DescribeCacheSubnetGroupsInput{
		MaxRecords: aws.Int32(100),
	}

	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxRecords {
			if limit < 20 {
				input.MaxRecords = aws.Int32(20)
			} else {
				input.MaxRecords = aws.Int32(limit)
			}
		}
	}

	paginator := elasticache.NewDescribeCacheSubnetGroupsPaginator(svc, input, func(o *elasticache.DescribeCacheSubnetGroupsPaginatorOptions) {
		o.Limit = *input.MaxRecords
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_elasticache_subnet_group.listElastiCacheSubnetGroups", "api_error", err)
			return nil, err
		}

		for _, cacheSubnetGroup := range output.CacheSubnetGroups {
			d.StreamListItem(ctx, cacheSubnetGroup)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getElastiCacheSubnetGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_subnet_group.getElastiCacheSubnetGroup", "get_client_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	cacheSubnetGroupName := quals["cache_subnet_group_name"].GetStringValue()

	params := &elasticache.DescribeCacheSubnetGroupsInput{
		CacheSubnetGroupName: aws.String(cacheSubnetGroupName),
	}

	op, err := svc.DescribeCacheSubnetGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_subnet_group.getElastiCacheSubnetGroup", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.CacheSubnetGroups) > 0 {
		return op.CacheSubnetGroups[0], nil
	}
	return nil, nil
}
