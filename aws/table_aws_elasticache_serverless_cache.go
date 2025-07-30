package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElastiCacheServerlessCache(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_serverless_cache",
		Description: "AWS ElastiCache Serverless Cache",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("serverless_cache_name"),
			Hydrate:    getElastiCacheServerlessCache,
			Tags:       map[string]string{"service": "elasticache", "action": "DescribeServerlessCaches"},
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheServerlessCaches,
			Tags:    map[string]string{"service": "elasticache", "action": "DescribeServerlessCaches"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listTagsForElastiCacheServerlessCache,
				Tags: map[string]string{"service": "elasticache", "action": "ListTagsForResource"},
			},
		},

		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ELASTICACHE_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "serverless_cache_name",
				Description: "The unique identifier of the serverless cache.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) of the serverless cache.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "status",
				Description: "The current status of the serverless cache.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine",
				Description: "The name of the cache engine (e.g., redis, valkey) used by the serverless cache.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "full_engine_version",
				Description: "The version of the cache engine that is used in this serverless cache.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The date and time when the serverless cache was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "subnet_ids",
				Description: "The subnet IDs for the serverless cache.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_group_ids",
				Description: "The security group IDs for the serverless cache.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "user_group_id",
				Description: "The ID of the user group associated with the serverless cache.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cache_usage_limits",
				Description: "The cache usage limits for the serverless cache.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "daily_snapshot_time",
				Description: "The daily time range during which ElastiCache takes a snapshot of the serverless cache.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "snapshot_retention_limit",
				Description: "The number of days for which ElastiCache retains automatic serverless cache snapshots before deleting them.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "The description of the serverless cache.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The ID of the KMS key used to encrypt the serverless cache.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "endpoint",
				Description: "The endpoint information for the serverless cache.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "reader_endpoint",
				Description: "The reader endpoint information for the serverless cache.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the serverless cache.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForElastiCacheServerlessCache,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerlessCacheName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForElastiCacheServerlessCache,
				Transform:   transform.From(serverlessCacheTagListToTurbotTags),
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

func listElastiCacheServerlessCaches(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_serverless_cache.listElastiCacheServerlessCaches", "connection_error", err)
		return nil, err
	}

	input := &elasticache.DescribeServerlessCachesInput{
		MaxResults: aws.Int32(100),
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
	paginator := elasticache.NewDescribeServerlessCachesPaginator(svc, input, func(o *elasticache.DescribeServerlessCachesPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_elasticache_serverless_cache.listElastiCacheServerlessCaches", "api_error", err)
			return nil, err
		}

		for _, serverlessCache := range output.ServerlessCaches {
			d.StreamListItem(ctx, serverlessCache)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getElastiCacheServerlessCache(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_serverless_cache.getElastiCacheServerlessCache", "connection_error", err)
		return nil, err
	}

	serverlessCacheName := d.EqualsQuals["serverless_cache_name"].GetStringValue()

	// Return nil, if no input provided
	if serverlessCacheName == "" {
		return nil, nil
	}

	params := &elasticache.DescribeServerlessCachesInput{
		ServerlessCacheName: aws.String(serverlessCacheName),
	}

	op, err := svc.DescribeServerlessCaches(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_serverless_cache.getElastiCacheServerlessCache", "api_error", err)
		return nil, err
	}

	if len(op.ServerlessCaches) > 0 {
		return op.ServerlessCaches[0], nil
	}

	return nil, nil
}

func listTagsForElastiCacheServerlessCache(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	serverlessCache := h.Item.(types.ServerlessCache)

	// Create session
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_serverless_cache.listTagsForElastiCacheServerlessCache", "connection_error", err)
		return nil, err
	}

	// Build param
	param := &elasticache.ListTagsForResourceInput{
		ResourceName: serverlessCache.ARN,
	}

	serverlessCacheTags, err := svc.ListTagsForResource(ctx, param)

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ServerlessCacheNotFoundFault" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_elasticache_serverless_cache.listTagsForElastiCacheServerlessCache", "api_error", err)
		return nil, err
	}

	return serverlessCacheTags, nil
}

//// TRANSFORM FUNCTIONS

func serverlessCacheTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	serverlessCacheTag := d.HydrateItem.(*elasticache.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if len(serverlessCacheTag.TagList) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range serverlessCacheTag.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
