package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRoute53ResolverEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_resolver_endpoint",
		Description: "AWS Route53 Resolver Endpoint",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate:           getAwsRoute53ResolverEndpoint,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsRoute53ResolverEndpoint,
		},
		GetMatrixItem: BuildRegionList,
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
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsRoute53ResolverEndpoint", "AWS_REGION", region)

	// Create session
	svc, err := Route53ResolverService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListResolverEndpointsPages(
		&route53resolver.ListResolverEndpointsInput{},
		func(page *route53resolver.ListResolverEndpointsOutput, isLast bool) bool {
			for _, parameter := range page.ResolverEndpoints {
				d.StreamListItem(ctx, parameter)

			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsRoute53ResolverEndpoint(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsSSMParameter")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	svc, err := Route53ResolverService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &route53resolver.GetResolverEndpointInput{
		ResolverEndpointId: &id,
	}

	// Get call
	data, err := svc.GetResolverEndpoint(params)
	if err != nil {
		logger.Debug("getAwsRoute53ResolverEndpoint", "ERROR", err)
		return nil, err
	}

	return data.ResolverEndpoint, nil
}

func listResolverEndpointIPAddresses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listResolverEndpointIpAddresses")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	resolverEndpointData := h.Item.(*route53resolver.ResolverEndpoint)

	// Create Session
	svc, err := Route53ResolverService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &route53resolver.ListResolverEndpointIpAddressesInput{
		ResolverEndpointId: resolverEndpointData.Id,
	}

	op, err := svc.ListResolverEndpointIpAddresses(params)
	if err != nil {
		logger.Debug("listResolverEndpointIpAddresses", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getAwsRoute53ResolverEndpointTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsRoute53ResolverEndpointTags")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	resolverEndpintData := h.Item.(*route53resolver.ResolverEndpoint)

	// Create Session
	svc, err := Route53ResolverService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &route53resolver.ListTagsForResourceInput{
		ResourceArn: resolverEndpintData.Arn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Debug("getAwsRoute53ResolverEndpointTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func route53resolverTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("route53resolverTagListToTurbotTags")
	tagList := d.Value.([]*route53resolver.Tag)

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
