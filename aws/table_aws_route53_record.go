package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
)

func tableAwsRoute53Record(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_record",
		Description: "AWS Route53 Record",
		// Get: &plugin.GetConfig{
		// 	KeyColumns:  plugin.SingleColumn("zone_id"),
		// 	ItemFromKey: hostedZone,
		// 	Hydrate:     getHostedZone,
		// },
		List: &plugin.ListConfig{
			ParentHydrate: listHostedZones,
			Hydrate:       listRoute53Records,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the record.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Record.Name"),
			},
			{
				Name:        "zone_id",
				Description: "The ID of the hosted zone to contain this record.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The record type. Valid values are A, AAAA, CAA, CNAME, MX, NAPTR, NS, PTR, SOA, SPF, SRV and TXT.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Record.Type"),
			},
			{
				Name:        "alias_target",
				Description: "Alias resource record sets only: Information about the AWS resource, such as a CloudFront distribution or an Amazon S3 bucket, that you want to route traffic to.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Record.AliasTarget"),
			},
			{
				Name:        "failover",
				Description: "Failover resource record sets only: To configure failover, you add the Failover element to two resource record sets. For one resource record set, you specify PRIMARY as the value for Failover; for the other resource record set, you specify SECONDARY. In addition, you include the HealthCheckId element and specify the health check that you want Amazon Route 53 to perform for each resource record set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Record.Failover"),
			},
			{
				Name:        "geolocation_routing_policy",
				Description: "Geolocation resource record sets only: A complex type that lets you control how Amazon Route 53 responds to DNS queries based on the geographic origin of the query. For example, if you want all queries from Africa to be routed to a web server with an IP address of 192.0.2.111, create a resource record set with a Type of A and a ContinentCode of AF.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Record.GeoLocation"),
			},
			{
				Name:        "health_check_id",
				Description: "The health check the record should be associated with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Record.HealthCheckId"),
			},
			{
				Name:        "multi_value_answer",
				Description: "Multivalue answer resource record sets only: To route traffic approximately randomly to multiple resources, such as web servers, create one multivalue answer record for each resource and specify true for MultiValueAnswer.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Record.MultiValueAnswer"),
			},
			{
				Name:        "latency_region",
				Description: "An AWS region from which to measure latency",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Record.Region"),
			},
			{
				Name:        "records",
				Description: "If the health check or hosted zone was created by another service, an optional description that can be provided by the other service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(flattenResourceRecords),
			},
			{
				Name:        "set_identifier",
				Description: "Unique identifier to differentiate records with routing policies from one another.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Record.SetIdentifier"),
			},
			{
				Name:        "ttl",
				Description: "The resource record cache time to live (TTL), in seconds.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Record.TTL"),
			},
			{
				Name:        "traffic_policy_instance_id",
				Description: "The ID of the traffic policy instance that Route 53 created this resource record set for.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Record.TrafficPolicyInstanceId"),
			},
			{
				Name:        "weight",
				Description: "Weighted resource record sets only: Among resource record sets that have the same combination of DNS name and type, a value that determines the proportion of DNS queries that Amazon Route 53 responds to using the current resource record set. Route 53 calculates the sum of the weights for the resource record sets that have the same combination of DNS name and type. Route 53 then responds to queries based on the ratio of a resource's weight to the total.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Record.Weight"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Record.Name"),
			},
			// {
			// 	Name:        "akas",
			// 	Description: resourceInterfaceDescription("akas"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     getRoute53HostedZoneTurbotAkas,
			// 	Transform:   transform.FromValue(),
			// },
		}),
	}
}

type recordInfo struct {
	ZoneID *string
	Record *route53.ResourceRecordSet
}

//// LIST FUNCTION

func listRoute53Records(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listRoute53Records", "AWS_REGION", defaultRegion)
	hostedZone := h.Item.(*route53.HostedZone)

	// Create session
	svc, err := Route53Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	err = svc.ListResourceRecordSetsPages(
		&route53.ListResourceRecordSetsInput{
			HostedZoneId: hostedZone.Id,
		},
		func(page *route53.ListResourceRecordSetsOutput, isLast bool) bool {
			for _, record := range page.ResourceRecordSets {
				d.StreamLeafListItem(ctx, &recordInfo{hostedZone.Id, record})
			}
			return true
		},
	)

	return nil, err
}

//// TRANSFORM FUNCTION

func flattenResourceRecords(_ context.Context, d *transform.TransformData) (interface{}, error) {
	hostedZone := d.HydrateItem.(*recordInfo)
	typeStr := types.SafeString(hostedZone.Record.Type)

	strs := make([]string, 0, len(hostedZone.Record.ResourceRecords))
	for _, r := range hostedZone.Record.ResourceRecords {
		if r.Value != nil {
			s := *r.Value
			if typeStr == "TXT" || typeStr == "SPF" {
				s = fmt.Sprintf(`"%s"`, s)
			}
			strs = append(strs, s)
		}
	}
	return strs, nil
}
