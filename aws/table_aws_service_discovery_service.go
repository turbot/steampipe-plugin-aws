package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServiceDiscoveryService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_service_discovery_service",
		Description: "AWS Service Discovery Service",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ServiceNotFound"}),
			},
			Hydrate: getServiceDiscoveryService,
			Tags:    map[string]string{"service": "servicediscovery", "action": "GetService"},
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceDiscoveryServices,
			Tags:    map[string]string{"service": "servicediscovery", "action": "ListServices"},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "namespace_id",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getServiceDirectoryServiceTags,
				Tags: map[string]string{"service": "servicediscovery", "action": "ListTagsForResource"},
			},
			{
				Func: getServiceDiscoveryService,
				Tags: map[string]string{"service": "servicediscovery", "action": "GetService"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SERVICEDISCOVERY_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID that Cloud Map assigned to the service when you created it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that Cloud Map assigns to the service when you create it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_date",
				Description: "The date and time that the service was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description for the service.",
				Type:        proto.ColumnType_STRING,
			},
			// The value for the namespace_id column will be populated the service type is DNS_HTTP.
			{
				Name:        "namespace_id",
				Description: "The ID of the namespace.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getServiceDiscoveryService,
				Transform:   transform.FromField("DnsConfig.NamespaceId"),
			},
			{
				Name:        "instance_count",
				Description: "The number of instances that are currently associated with the service.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "type",
				Description: "Describes the systems that can be used to discover the service instances. DNS_HTTP The service instances can be discovered using either DNS queries or the DiscoverInstances API operation. HTTP The service instances can only be discovered using the DiscoverInstances API operation. DNS Reserved.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "routing_policy",
				Description: "The routing policy that you want to apply to all Route 53 DNS records that Cloud Map creates when you register an instance and specify this service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DnsConfig.RoutingPolicy"),
			},
			{
				Name:        "dns_records",
				Description: "An array that contains one DnsRecord object for each Route 53 DNS record that you want Cloud Map to create when you register an instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DnsConfig.DnsRecords"),
			},
			{
				Name:        "health_check_config",
				Description: "Public DNS and HTTP namespaces only. Settings for an optional health check. If you specify settings for a health check, Cloud Map associates the health check with the records that you specify in DnsConfig.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "health_check_custom_config",
				Description: "Information about an optional custom health check.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "Information about the tags associated with the namespace.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceDirectoryServiceTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceDirectoryServiceTags,
				Transform:   transform.From(serviceTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listServiceDiscoveryServices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Client
	svc, err := ServiceDiscoveryClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_service_discovery_service.listServiceDiscoveryServices", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &servicediscovery.ListServicesInput{
		MaxResults: &maxLimit,
	}

	filters := buildServiceDiscoveryServiceFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := servicediscovery.NewListServicesPaginator(svc, input, func(o *servicediscovery.ListServicesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_service_discovery_service.listServiceDiscoveryServices", "api_error", err)
			return nil, err
		}

		for _, item := range output.Services {
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getServiceDiscoveryService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetStringValue()

	// Empty check
	if id == "" {
		return nil, nil
	}

	// Create client
	svc, err := ServiceDiscoveryClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_service_discovery_service.getServiceDiscoveryService", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &servicediscovery.GetServiceInput{
		Id: aws.String(id),
	}

	// Get call
	op, err := svc.GetService(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_service_discovery_service.getServiceDiscoveryService", "api_error", err)
		return nil, err
	}

	if op.Service != nil {
		return op.Service, nil
	}

	return nil, nil
}

func getServiceDirectoryServiceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string

	switch item := h.Item.(type) {
	case types.ServiceSummary:
		arn = *item.Arn
	case *types.Service:
		arn = *item.Arn
	}

	// Create client
	svc, err := ServiceDiscoveryClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_service_discovery_service.getServiceDirectoryServiceTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &servicediscovery.ListTagsForResourceInput{
		ResourceARN: aws.String(arn),
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_service_discovery_service.getServiceDirectoryServiceTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func serviceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*servicediscovery.ListTagsForResourceOutput)
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}

//// UTILITY FUNCTIONS

// Build service discovery service list call input filter
func buildServiceDiscoveryServiceFilter(quals plugin.KeyColumnQualMap) []types.ServiceFilter {
	filters := make([]types.ServiceFilter, 0)

	filterQuals := map[string]string{
		"namespace_id": string(types.ServiceFilterNameNamespaceId),
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.ServiceFilter{
				Name:      types.ServiceFilterName(filterName),
				Condition: types.FilterConditionEq,
			}
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []string{val}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
