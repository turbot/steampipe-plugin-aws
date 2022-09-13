package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsElastiCacheParameterGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_parameter_group",
		Description: "AWS ElastiCache Parameter Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cache_parameter_group_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"CacheParameterGroupNotFound", "InvalidParameterValueException"}),
			},
			Hydrate: getElastiCacheParameterGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheParameterGroup,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cache_parameter_group_name",
				Description: "The name of the cache parameter group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) of the cache parameter group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "description",
				Description: "The description for the cache parameter group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cache_parameter_group_family",
				Description: "The name of the cache parameter group family that this cache parameter group is compatible with.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_global",
				Description: "Indicates whether the parameter group is associated with a Global Datastore.",
				Type:        proto.ColumnType_BOOL,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CacheParameterGroupName"),
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

func listElastiCacheParameterGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ElastiCacheService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &elasticache.DescribeCacheParameterGroupsInput{
		MaxRecords: aws.Int64(100),
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
	err = svc.DescribeCacheParameterGroupsPages(
		input,
		func(page *elasticache.DescribeCacheParameterGroupsOutput, isLast bool) bool {
			for _, parameterGroup := range page.CacheParameterGroups {
				d.StreamListItem(ctx, parameterGroup)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getElastiCacheParameterGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getElastiCacheParameterGroup")

	// Create service
	svc, err := ElastiCacheService(ctx, d)
	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	cacheParameterGroupName := quals["cache_parameter_group_name"].GetStringValue()

	params := &elasticache.DescribeCacheParameterGroupsInput{
		CacheParameterGroupName: aws.String(cacheParameterGroupName),
	}

	op, err := svc.DescribeCacheParameterGroups(params)
	if err != nil {
		return nil, err
	}

	if op.CacheParameterGroups != nil && len(op.CacheParameterGroups) > 0 {
		return op.CacheParameterGroups[0], nil
	}
	return nil, nil
}
