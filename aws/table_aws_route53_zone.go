package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
)

func tableAwsRoute53Zone(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_zone",
		Description: "AWS Route53 Zone",
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.SingleColumn("id"),
			ItemFromKey: hostedZone,
			Hydrate:     getHostedZone,
		},
		List: &plugin.ListConfig{
			Hydrate: listHostedZones,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the domain. For public hosted zones, this is the name that is registered with your DNS registrar.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID that Amazon Route 53 assigned to the hosted zone when it was created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(route53ZoneID),
			},
			{
				Name:        "caller_reference",
				Description: "The value that you specified for CallerReference when you created the hosted zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "config_comment",
				Description: "A complex type that includes the Comment and PrivateZone elements",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Config.Comment"),
			},
			{
				Name:        "config_private_zone",
				Description: "A complex type that includes the Comment and PrivateZone elements",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Config.PrivateZone"),
			},
			{
				Name:        "linked_service_principal",
				Description: "If the health check or hosted zone was created by another service, the service that created the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LinkedService.ServicePrincipal"),
			},
			{
				Name:        "linked_service_description",
				Description: "If the health check or hosted zone was created by another service, an optional description that can be provided by the other service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LinkedService.Description"),
			},
			{
				Name:        "resource_record_set_count",
				Description: "The number of resource record sets in the hosted zone",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "tags_src",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getHostedZoneTags,
				Transform:   transform.FromField("ResourceTagSet.Tags"),
			},
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
				Hydrate:     getHostedZoneTags,
				Transform:   transform.FromField("ResourceTagSet.Tags").Transform(route53HostedZoneTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53HostedZoneTurbotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// BUILD HYDRATE INPUT

func hostedZone(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	id := quals["id"].GetStringValue()
	item := &route53.HostedZone{
		Id: &id,
	}
	return item, nil
}

//// LIST FUNCTION

func listHostedZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listhostedZone", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := Route53Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	err = svc.ListHostedZonesPages(
		&route53.ListHostedZonesInput{},
		func(page *route53.ListHostedZonesOutput, isLast bool) bool {
			for _, hostedZone := range page.HostedZones {
				d.StreamListItem(ctx, hostedZone)
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getHostedZone(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHostedZone")

	defaultRegion := GetDefaultRegion()
	hostedZone := h.Item.(*route53.HostedZone)

	// Create session
	svc, err := Route53Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &route53.GetHostedZoneInput{
		Id: hostedZone.Id,
	}

	// execute list call
	item, err := svc.GetHostedZone(params)
	if err != nil {
		return nil, err
	}

	return item.HostedZone, nil
}

func getHostedZoneTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHostedZone")

	defaultRegion := GetDefaultRegion()
	hostedZone := h.Item.(*route53.HostedZone)

	// Create session
	svc, err := Route53Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &route53.ListTagsForResourceInput{
		ResourceId:   hostedZone.Id,
		ResourceType: types.String("hostedzone"),
	}

	// execute list call
	resp, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getRoute53HostedZoneTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRoute53HostedZoneTurbotAkas")
	hostedZone := h.Item.(*route53.HostedZone)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	id := strings.Split(string(*hostedZone.Id), "/")

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":route53:::" + "hostedzone/" + id[2]}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func route53ZoneID(_ context.Context, d *transform.TransformData) (interface{}, error) {
	hostedZone := d.HydrateItem.(*route53.HostedZone)
	id := strings.Split(string(*hostedZone.Id), "/")[2]

	return id, nil
}

func route53HostedZoneTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("route53HostedZoneTurbotTags")
	tags := d.Value.([]*route53.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
