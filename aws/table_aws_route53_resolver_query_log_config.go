package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53resolver"
	"github.com/aws/aws-sdk-go-v2/service/route53resolver/types"

	route53resolverv1 "github.com/aws/aws-sdk-go/service/route53resolver"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAwsRoute53ResolverQueryLogConfig(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_resolver_query_log_config",
		Description: "AWS Route53 Resolver Query Logging Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getRoute53ResolverQueryLogConfig,
			Tags:       map[string]string{"service": "route53resolver", "action": "GetResolverQueryLogConfig"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(route53resolverv1.EndpointsID),
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "creator_request_id", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
				{Name: "ip_address_count", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "status", Require: plugin.Optional},
			},
			Hydrate: listRoute53ResolverQueryLogConfigs,
			Tags:    map[string]string{"service": "route53resolver", "action": "ListResolverQueryLogConfigs"},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID for the query logging configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the query logging configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) for the query logging configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The date and time that the query logging configuration was created, in Unix time format and Coordinated Universal Time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status",
				Description: "The status of the specified query logging configuration. Valid values include CREATING|CREATED|DELETING|FAILED.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "association_count",
				Description: "The number of VPCs that are associated with the query logging configuration.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "ip_address_count",
				Description: "The number of IP addresses that you have associated with the Resolver endpoint.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromQual("ip_address_count"),
			},
			{
				Name:        "creator_request_id",
				Description: "A unique string that identifies the request that created the query logging configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "destination_arn",
				Description: "The ARN of the resource that you want Resolver to send query logs: an Amazon S3 bucket, a CloudWatch Logs log group, or a Kinesis Data Firehose delivery stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The Amazon Web Services account ID for the account that created the query logging configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "share_status",
				Description: "An indication of whether the query logging configuration is shared with other Amazon Web Services accounts, or was shared with the current account by another Amazon Web Services account. Sharing is configured through Resource Access Manager (RAM).",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns

			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53ResolverQueryLogConfigAkas,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listRoute53ResolverQueryLogConfigs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := Route53ResolverClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_resolver_query_log_config.listRoute53ResolverQueryLogConfigs", "client_error", err)
		return nil, err
	}

	maxItems := int32(100)
	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	input := &route53resolver.ListResolverQueryLogConfigsInput{
		MaxResults: aws.Int32(maxItems),
	}

	filters := buildListRoute53ResolverQueryLogConfigInputParam(d.EqualsQuals)

	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := route53resolver.NewListResolverQueryLogConfigsPaginator(svc, input, func(o *route53resolver.ListResolverQueryLogConfigsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_route53_resolver_query_log_config.listRoute53ResolverQueryLogConfigs", "api_error", err)
			return nil, err
		}

		for _, config := range output.ResolverQueryLogConfigs {
			d.StreamListItem(ctx, config)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRoute53ResolverQueryLogConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := Route53ResolverClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_resolver_query_log_config.getRoute53ResolverQueryLogConfig", "client_error", err)
		return nil, err
	}

	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	input := &route53resolver.GetResolverQueryLogConfigInput{
		ResolverQueryLogConfigId: &id,
	}

	op, err := svc.GetResolverQueryLogConfig(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_resolver_query_log_config.getRoute53ResolverQueryLogConfig", "api_error", err)
		return nil, err
	}
	return *op.ResolverQueryLogConfig, nil
}

//// TRANSFORM FUNCTION

func getRoute53ResolverQueryLogConfigAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string

	switch item := h.Item.(type) {
	case *types.ResolverQueryLogConfig:
		arn = *item.Arn
	case types.ResolverQueryLogConfig:
		arn = *item.Arn
	}

	return arn, nil
}

func buildListRoute53ResolverQueryLogConfigInputParam(equalQuals plugin.KeyColumnEqualsQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"creator_request_id": "CreatorRequestId",
		"ip_address_count":   "IpAddressCount",
		"name":               "Name",
		"status":             "Status",
	}

	for columnName, filterName := range filterQuals {
		if equalQuals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			value := equalQuals[columnName]
			if columnName == "ip_address_count" {
				if value.GetInt64Value() != 0 {
					filter.Values = []string{fmt.Sprint(equalQuals[columnName].GetInt64Value())}
				}
			} else {

				if value.GetStringValue() != "" {
					filter.Values = []string{equalQuals[columnName].GetStringValue()}
				}
				filters = append(filters, filter)
			}
		}
	}

	return filters
}
