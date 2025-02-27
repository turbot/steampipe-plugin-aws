package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElastiCacheParameterGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_parameter_group",
		Description: "AWS ElastiCache Parameter Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cache_parameter_group_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"CacheParameterGroupNotFound", "InvalidParameterValueException"}),
			},
			Hydrate: getElastiCacheParameterGroup,
			Tags:    map[string]string{"service": "elasticache", "action": "DescribeCacheParameterGroups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheParameterGroup,
			Tags:    map[string]string{"service": "elasticache", "action": "DescribeCacheParameterGroups"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ELASTICACHE_SERVICE_ID),
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
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_parameter_group.listElastiCacheParameterGroup", "get_client_error", err)
		return nil, err
	}

	input := &elasticache.DescribeCacheParameterGroupsInput{
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

	// List call
	paginator := elasticache.NewDescribeCacheParameterGroupsPaginator(svc, input, func(o *elasticache.DescribeCacheParameterGroupsPaginatorOptions) {
		o.Limit = *input.MaxRecords
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_elasticache_parameter_group.listElastiCacheParameterGroup", "api_error", err)
			return nil, err
		}

		for _, parameterGroup := range output.CacheParameterGroups {
			d.StreamListItem(ctx, parameterGroup)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getElastiCacheParameterGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_parameter_group.getElastiCacheParameterGroup", "get_client_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	cacheParameterGroupName := quals["cache_parameter_group_name"].GetStringValue()

	params := &elasticache.DescribeCacheParameterGroupsInput{
		CacheParameterGroupName: aws.String(cacheParameterGroupName),
	}

	op, err := svc.DescribeCacheParameterGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_parameter_group.getElastiCacheParameterGroup", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.CacheParameterGroups) > 0 {
		return op.CacheParameterGroups[0], nil
	}
	return nil, nil
}
