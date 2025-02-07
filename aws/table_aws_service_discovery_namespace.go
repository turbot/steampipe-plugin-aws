package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"

	servicediscoveryEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServiceDiscoveryNamespace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_service_discovery_namespace",
		Description: "AWS Service Discovery Namespace",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NamespaceNotFound"}),
			},
			Hydrate: getServiceDiscoveryNamespace,
			Tags:    map[string]string{"service": "servicediscovery", "action": "GetNamespace"},
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceDiscoveryNamespaces,
			Tags:    map[string]string{"service": "servicediscovery", "action": "ListNamespaces"},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "type",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getServiceDirectoryNamespaceTags,
				Tags: map[string]string{"service": "servicediscovery", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(servicediscoveryEndpoint.AWS_SERVICEDISCOVERY_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID of the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that Cloud Map assigns to the namespace when you create it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_date",
				Description: "The date and time that the namespace was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description for the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_count",
				Description: "The number of services that were created using the namespace.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "type",
				Description: "The type of the namespace, either public or private.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dns_properties",
				Description: "A complex type that contains the ID for the Route 53 hosted zone that Cloud Map creates when you create a namespace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.DnsProperties"),
			},
			{
				Name:        "http_properties",
				Description: "A complex type that contains the name of an HTTP namespace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.HttpProperties"),
			},
			{
				Name:        "tags_src",
				Description: "Information about the tags associated with the namespace.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceDirectoryNamespaceTags,
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
				Hydrate:     getServiceDirectoryNamespaceTags,
				Transform:   transform.From(namespaceTurbotTags),
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

func listServiceDiscoveryNamespaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Client
	svc, err := ServiceDiscoveryClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_service_discovery_namespace.listServiceDiscoveryNamespaces", "connection_error", err)
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

	input := &servicediscovery.ListNamespacesInput{
		MaxResults: &maxLimit,
	}

	filters := buildServiceDiscoveryNamespaceFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := servicediscovery.NewListNamespacesPaginator(svc, input, func(o *servicediscovery.ListNamespacesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_service_discovery_namespace.listServiceDiscoveryNamespaces", "api_error", err)
			return nil, err
		}

		for _, item := range output.Namespaces {
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

func getServiceDiscoveryNamespace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	id := d.EqualsQualString("id")

	// Empty check
	if id == "" {
		return nil, nil
	}

	// Create client
	svc, err := ServiceDiscoveryClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_service_discovery_namespace.getServiceDiscoveryNamespace", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &servicediscovery.GetNamespaceInput{
		Id: aws.String(id),
	}

	// Get call
	op, err := svc.GetNamespace(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_service_discovery_namespace.getServiceDiscoveryNamespace", "api_error", err)
		return nil, err
	}

	if op.Namespace != nil {
		return op.Namespace, nil
	}

	return nil, nil
}

func getServiceDirectoryNamespaceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	var arn string

	switch item := h.Item.(type) {
	case types.NamespaceSummary:
		arn = *item.Arn
	case *types.Namespace:
		arn = *item.Arn
	}

	// Create client
	svc, err := ServiceDiscoveryClient(ctx, d)
	if err != nil {
		logger.Error("aws_service_discovery_namespace.getServiceDirectoryNamespaceTags", "connection_error", err)
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
		logger.Error("aws_service_discovery_namespace.getServiceDirectoryNamespaceTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func namespaceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
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

// Build service discovery namespace list call input filter
func buildServiceDiscoveryNamespaceFilter(quals plugin.KeyColumnQualMap) []types.NamespaceFilter {
	filters := make([]types.NamespaceFilter, 0)

	filterQuals := map[string]string{
		"name": string(types.NamespaceFilterNameName),
		"type": string(types.NamespaceFilterNameType),
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.NamespaceFilter{
				Name:      types.NamespaceFilterName(filterName),
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
