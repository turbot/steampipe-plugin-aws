package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsRoute53Zone(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_zone",
		Description: "AWS Route53 Zone",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getHostedZone,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchHostedZone"}),
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
				Transform:   transform.FromField("QueryLoggingConfigs"),
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
	types.HostedZone
	VPCs []types.VPC
}

//// LIST FUNCTION

func listHostedZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_zone.listHostedZones", "client_error", err)
		return nil, err
	}

	// https://docs.aws.amazon.com/Route53/latest/APIReference/API_ListHostedZones.html
	// The maximum/minimum record set per page is not mentioned in doc, so it has been set 1000 to max and 1 to min
	maxItems := int32(1000)
	input := route53.ListHostedZonesInput{}

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

	input.MaxItems = aws.Int32(maxItems)
	paginator := route53.NewListHostedZonesPaginator(svc, &input, func(o *route53.ListHostedZonesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_route53_zone.listHostedZones", "api_error", err)
			return nil, err
		}

		for _, hostedZone := range output.HostedZones {
			d.StreamListItem(ctx, HostedZoneResult{HostedZone: hostedZone})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getHostedZone(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_zone.getHostedZone", "client_error", err)
		return nil, err
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	if h.Item != nil && id == "" {
		hostedZone := h.Item.(HostedZoneResult)
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
	item, err := svc.GetHostedZone(ctx, params)
	if err != nil {
		return nil, err
	}

	return HostedZoneResult{
		HostedZone: *item.HostedZone,
		VPCs:       item.VPCs,
	}, nil
}

func getHostedZoneTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	hostedZone := h.Item.(HostedZoneResult)

	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_zone.getHostedZoneTags", "client_error", err)
		return nil, err
	}

	params := &route53.ListTagsForResourceInput{
		ResourceId:   &strings.Split(*hostedZone.Id, "/")[2],
		ResourceType: "hostedzone",
	}

	// execute list call
	resp, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NoSuchHostedZone" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_route53_zone.getHostedZoneTags", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func getHostedZoneQueryLoggingConfigs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	hostedZone := h.Item.(HostedZoneResult)

	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_zone.getHostedZoneQueryLoggingConfigs", "client_error", err)
		return nil, err
	}

	params := &route53.ListQueryLoggingConfigsInput{
		HostedZoneId: &strings.Split(*hostedZone.Id, "/")[2],
	}
	resp, err := svc.ListQueryLoggingConfigs(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NoSuchHostedZone" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_route53_zone.getHostedZoneQueryLoggingConfigs", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func getHostedZoneDNSSEC(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	hostedZone := h.Item.(HostedZoneResult)

	// Operation is unsupported for private hosted zones.
	if hostedZone.Config.PrivateZone {
		return nil, nil
	}

	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_zone.getHostedZoneDNSSEC", "client_error", err)
		return nil, err
	}

	params := &route53.GetDNSSECInput{
		HostedZoneId: hostedZone.Id,
	}

	// execute list call
	resp, err := svc.GetDNSSEC(ctx, params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getRoute53HostedZoneTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	hostedZone := h.Item.(HostedZoneResult)

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
	hostedZone := d.HydrateItem.(HostedZoneResult)
	id := strings.Split(string(*hostedZone.Id), "/")[2]

	return id, nil
}

func route53HostedZoneTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.Value.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	if len(tags) > 0 {
		turbotTagsMap := map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}

	return nil, nil
}
