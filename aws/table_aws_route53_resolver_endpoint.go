package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53resolver"
	"github.com/aws/aws-sdk-go-v2/service/route53resolver/types"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRoute53ResolverEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_resolver_endpoint",
		Description: "AWS Route53 Resolver Endpoint",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getAwsRoute53ResolverEndpoint,
			Tags:    map[string]string{"service": "route53resolver", "action": "GetResolverEndpoint"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsRoute53ResolverEndpoint,
			Tags:    map[string]string{"service": "route53resolver", "action": "ListResolverEndpoints"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "creator_request_id", Require: plugin.Optional},
				{Name: "direction", Require: plugin.Optional},
				{Name: "host_vpc_id", Require: plugin.Optional},
				{Name: "ip_address_count", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listResolverEndpointIPAddresses,
				Tags: map[string]string{"service": "route53resolver", "action": "ListResolverEndpointIpAddresses"},
			},
			{
				Func: getAwsRoute53ResolverEndpointTags,
				Tags: map[string]string{"service": "route53resolver", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ROUTE53RESOLVER_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name that you assigned to the Resolver endpoint when you submitted a CreateResolverEndpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID of the Resolver endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) for the Resolver endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The date and time that the endpoint was created, in Unix time format and Coordinated Universal Time (UTC).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creator_request_id",
				Description: "A unique string that identifies the request that created the Resolver endpoint.The CreatorRequestId allows failed requests to be retried without the risk of executing the operation twice.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "preferred_instance_type",
				Description: "The Amazon EC2 instance type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resolver_endpoint_type",
				Description: "The Resolver endpoint IP address type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "outpost_arn",
				Description: "The ARN (Amazon Resource Name) for the Outpost.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "direction",
				Description: "Indicates whether the Resolver endpoint allows inbound or outbound DNS queries.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "host_vpc_id",
				Description: "The ID of the VPC that you want to create the Resolver endpoint in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("HostVPCId"),
			},
			{
				Name:        "ip_address_count",
				Description: "The number of IP addresses that the Resolver endpoint can use for DNS queries.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "modification_time",
				Description: "The date and time that the endpoint was last modified, in Unix time format and Coordinated Universal Time (UTC).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "A code that specifies the current status of the Resolver endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_message",
				Description: "A detailed description of the status of the Resolver endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_addresses",
				Description: "Information about the IP addresses in your VPC that DNS queries originate from (for outbound endpoints) or that you forward DNS queries to (for inbound endpoints).",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listResolverEndpointIPAddresses,
			},
			{
				Name:        "security_group_ids",
				Description: "The ID of one or more security groups that control access to this VPC.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "protocols",
				Description: "Protocols used for the endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Resolver endpoint.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRoute53ResolverEndpointTags,
				Transform:   transform.FromField("Tags"),
			},
			// Standard columns for all tables
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
				Hydrate:     getAwsRoute53ResolverEndpointTags,
				Transform:   transform.FromField("Tags").Transform(route53resolverTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsRoute53ResolverEndpoint(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := Route53ResolverClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_resolver_endpoint.listAwsRoute53ResolverEndpoint", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(100)
	input := route53resolver.ListResolverEndpointsInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	filter := buildRoute53ResolverEndpointFilter(d.Quals)
	if len(filter) > 0 {
		input.Filters = filter
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := route53resolver.NewListResolverEndpointsPaginator(svc, &input, func(o *route53resolver.ListResolverEndpointsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_route53_resolver_endpoint.listAwsRoute53ResolverEndpoint", "api_error", err)
			return nil, err
		}

		for _, resolverEndpoint := range output.ResolverEndpoints {
			d.StreamListItem(ctx, resolverEndpoint)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsRoute53ResolverEndpoint(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	id := d.EqualsQuals["id"].GetStringValue()

	// Create session
	svc, err := Route53ResolverClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_resolver_endpoint.getAwsRoute53ResolverEndpoint", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &route53resolver.GetResolverEndpointInput{
		ResolverEndpointId: &id,
	}

	// Execute get call
	data, err := svc.GetResolverEndpoint(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_resolver_endpoint.getAwsRoute53ResolverEndpoint", "api_error", err)
		return nil, err
	}
	return data.ResolverEndpoint, nil
}

func listResolverEndpointIPAddresses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	resolverEndpointData := h.Item.(types.ResolverEndpoint)

	// Create session
	svc, err := Route53ResolverClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_resolver_endpoint.listResolverEndpointIPAddresses", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &route53resolver.ListResolverEndpointIpAddressesInput{
		ResolverEndpointId: resolverEndpointData.Id,
	}

	op, err := svc.ListResolverEndpointIpAddresses(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_resolver_endpoint.listResolverEndpointIPAddresses", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsRoute53ResolverEndpointTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	route53resolverEndpointArn := ""
	switch h.Item.(type) {
	case types.ResolverEndpoint:
		route53resolverEndpointArn = *h.Item.(types.ResolverEndpoint).Arn
	case *types.ResolverEndpoint:
		route53resolverEndpointArn = *h.Item.(*types.ResolverEndpoint).Arn
	}

	// Create session
	svc, err := Route53ResolverClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_resolver_endpoint.listResolverEndpointIPAddresses", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &route53resolver.ListTagsForResourceInput{
		ResourceArn: aws.String(route53resolverEndpointArn),
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_resolver_endpoint.listResolverEndpointIPAddresses", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func route53resolverTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.Value.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	} else {
		return nil, nil
	}

	return turbotTagsMap, nil
}

//// UTILITY FUNCTION

// Build route53resolver endpoint list call input filter
func buildRoute53ResolverEndpointFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"creator_request_id": "CreatorRequestId",
		"direction":          "Direction",
		"host_vpc_id":        "HostVPCId",
		"ip_address_count":   "IpAddressCount",
		"status":             "Status",
		"name":               "Name",
	}

	columnsInt := []string{"ip_address_count"}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			if helpers.StringSliceContains(columnsInt, columnName) { //check Int columns
				value := getQualsValueByColumn(quals, columnName, "int64")
				val, ok := value.(int64)
				if ok {
					filter.Values = []string{fmt.Sprint(val)}
				}
			} else {
				value := getQualsValueByColumn(quals, columnName, "string")
				val, ok := value.(string)
				if ok {
					filter.Values = []string{fmt.Sprint(val)}
				}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
