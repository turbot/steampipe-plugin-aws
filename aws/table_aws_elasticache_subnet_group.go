package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsElastiCacheSubnetGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_subnet_group",
		Description: "AWS ElastiCache Subnet Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cache_subnet_group_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"CacheSubnetGroupNotFoundFault"}),
			},
			Hydrate: getElastiCacheSubnetGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheSubnetGroups,
		},
		GetMatrixItemFunc: BuildRegionList,
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
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_elasticache_subnet_group.listElastiCacheSubnetGroups", "api_error", err)
			return nil, err
		}

		for _, cacheSubnetGroup := range output.CacheSubnetGroups {
			d.StreamListItem(ctx, cacheSubnetGroup)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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

	quals := d.KeyColumnQuals
	cacheSubnetGroupName := quals["cache_subnet_group_name"].GetStringValue()

	params := &elasticache.DescribeCacheSubnetGroupsInput{
		CacheSubnetGroupName: aws.String(cacheSubnetGroupName),
	}

	op, err := svc.DescribeCacheSubnetGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_subnet_group.getElastiCacheSubnetGroup", "api_error", err)
		return nil, err
	}

	if op.CacheSubnetGroups != nil && len(op.CacheSubnetGroups) > 0 {
		return op.CacheSubnetGroups[0], nil
	}
	return nil, nil
}
