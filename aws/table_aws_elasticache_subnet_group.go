package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsElastiCacheSubnetGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_subnet_group",
		Description: "AWS ElastiCache Subnet Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("cache_subnet_group_name"),
			ShouldIgnoreError: isNotFoundError([]string{"CacheSubnetGroupNotFoundFault"}),
			Hydrate:           getElastiCacheSubnetGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheSubnetGroups,
		},
		GetMatrixItem: BuildRegionList,
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
	svc, err := ElastiCacheService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeCacheSubnetGroupsPages(
		&elasticache.DescribeCacheSubnetGroupsInput{},
		func(page *elasticache.DescribeCacheSubnetGroupsOutput, isLast bool) bool {
			for _, cacheSubnetGroup := range page.CacheSubnetGroups {
				d.StreamListItem(ctx, cacheSubnetGroup)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getElastiCacheSubnetGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getElastiCacheSubnetGroup")

	// Create service
	svc, err := ElastiCacheService(ctx, d)
	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	cacheSubnetGroupName := quals["cache_subnet_group_name"].GetStringValue()

	params := &elasticache.DescribeCacheSubnetGroupsInput{
		CacheSubnetGroupName: aws.String(cacheSubnetGroupName),
	}

	op, err := svc.DescribeCacheSubnetGroups(params)
	if err != nil {
		return nil, err
	}

	if op.CacheSubnetGroups != nil && len(op.CacheSubnetGroups) > 0 {
		return op.CacheSubnetGroups[0], nil
	}
	return nil, nil
}
