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

func tableAwsElastiCacheParameterGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_parameter_group",
		Description: "AWS ElastiCache Parameter Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("cache_parameter_group_name"),
			ShouldIgnoreError: isNotFoundError([]string{"CacheParameterGroupNotFound", "InvalidParameterValueException"}),
			Hydrate:           getElastiCacheParameterGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheParameterGroup,
		},
		GetMatrixItem: BuildRegionList,
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
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := ElastiCacheService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeCacheParameterGroupsPages(
		&elasticache.DescribeCacheParameterGroupsInput{},
		func(page *elasticache.DescribeCacheParameterGroupsOutput, isLast bool) bool {
			for _, parameterGroup := range page.CacheParameterGroups {
				d.StreamListItem(ctx, parameterGroup)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getElastiCacheParameterGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getElastiCacheParameterGroup")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// create service
	svc, err := ElastiCacheService(ctx, d, region)
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
