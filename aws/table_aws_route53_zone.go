package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
)

func tableAwsRoute53Zone(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_zone",
		Description: "AWS Route53 Zone",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getHostedZone,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchHostedZone"}),
			},
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
				Name:        "comment",
				Description: "A comment for the zone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Config.Comment"),
			},
			{
				Name:        "private_zone",
				Description: "If true, the zone is Private hosted Zone, otherwise it is public.",
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
				Description: "The number of resource record sets in the hosted zone.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "query_logging_configs",
				Description: "A list of configuration for DNS query logging that is associated with the current AWS account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getHostedZoneQueryLoggingConfigs,
			},
			{
				Name:        "dnssec_key_signing_keys",
				Description: "The key-signing keys (KSKs) in AWS account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getHostedZoneDNSSEC,
				Transform:   transform.FromField("KeySigningKeys"),
			},
			{
				Name:        "dnssec_status",
				Description: "The status of DNSSEC.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getHostedZoneDNSSEC,
				Transform:   transform.FromField("Status"),
			},
			{
				Name:        "tags_src",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getHostedZoneTags,
				Transform:   transform.FromField("ResourceTagSet.Tags"),
			},
			{
				Name:        "vpcs",
				Description: "The list of VPCs that are authorized to be associated with the specified hosted zone.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getHostedZone,
				Transform:   transform.FromField("VPCs"),
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

type HostedZoneResult struct {
	route53.HostedZone
	VPCs []*route53.VPC
}

//// LIST FUNCTION

func listHostedZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listhostedZone")

	// Create session
	svc, err := Route53Service(ctx, d)
	if err != nil {
		return nil, err
	}

	// https://docs.aws.amazon.com/Route53/latest/APIReference/API_ListHostedZones.html
	// The maximum/minimum record set per page is not mentioned in doc, so it has been set 1000 to max and 1 to min
	input := &route53.ListHostedZonesInput{
		MaxItems: aws.String("1000"),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < 1000 {
			if *limit < 1 {
				input.MaxItems = aws.String("1")
			} else {
				input.MaxItems = aws.String(fmt.Sprint(*limit))
			}
		}
	}

	err = svc.ListHostedZonesPages(
		input,
		func(page *route53.ListHostedZonesOutput, isLast bool) bool {
			for _, hostedZone := range page.HostedZones {
				d.StreamListItem(ctx, &HostedZoneResult{HostedZone: *hostedZone})

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

func getHostedZone(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHostedZone")

	// Create session
	svc, err := Route53Service(ctx, d)
	if err != nil {
		return nil, err
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	if h.Item != nil && id == "" {
		hostedZone := h.Item.(*HostedZoneResult)
		id = *hostedZone.Id
	}

	// Error: pq: rpc error: code = Unknown desc = InvalidParameter: 1 validation error(s) found.
	// - minimum field size of 1, GetHostedZoneInput.Id.
	if len(id) < 1 {
		return nil, nil
	}

	params := &route53.GetHostedZoneInput{
		Id: aws.String(id),
	}

	// execute list call
	item, err := svc.GetHostedZone(params)
	if err != nil {
		return nil, err
	}

	return &HostedZoneResult{
		HostedZone: *item.HostedZone,
		VPCs:       item.VPCs,
	}, nil
}

func getHostedZoneTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHostedZone")
	hostedZone := h.Item.(*HostedZoneResult)

	// Create session
	svc, err := Route53Service(ctx, d)
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

func getHostedZoneQueryLoggingConfigs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHostedZoneQueryLoggingConfigs")
	hostedZone := h.Item.(*HostedZoneResult)

	// Create session
	svc, err := Route53Service(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &route53.ListQueryLoggingConfigsInput{
		HostedZoneId: &strings.Split(*hostedZone.Id, "/")[2],
	}
	resp, err := svc.ListQueryLoggingConfigs(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NoSuchHostedZone" {
				return nil, nil
			}
			return nil, err
		}
	}

	return resp, nil
}

func getHostedZoneDNSSEC(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHostedZoneDNSSEC")
	hostedZone := h.Item.(*HostedZoneResult)

	// Operation is unsupported for private hosted zones.
	if *hostedZone.Config.PrivateZone {
		return nil, nil
	}

	// Create session
	svc, err := Route53Service(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &route53.GetDNSSECInput{
		HostedZoneId: hostedZone.Id,
	}

	// execute list call
	resp, err := svc.GetDNSSEC(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getRoute53HostedZoneTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRoute53HostedZoneTurbotAkas")
	hostedZone := h.Item.(*HostedZoneResult)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
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
	hostedZone := d.HydrateItem.(*HostedZoneResult)
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
