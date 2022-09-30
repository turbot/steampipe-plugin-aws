package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsRoute53Record(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_record",
		Description: "AWS Route53 Record",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "zone_id", Require: plugin.Required},
				{Name: "name", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
			Hydrate: listRoute53Records,
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
				Name:        "geo_location",
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
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Record.Weight"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Record.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53RecordSetAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type recordInfo struct {
	ZoneID *string
	Record types.ResourceRecordSet
}

//// LIST FUNCTION

func listRoute53Records(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	hostedZoneID := d.KeyColumnQuals["zone_id"].GetStringValue()

	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_record.listRoute53Records", "client_error", err)
		return nil, err
	}

	input := route53.ListResourceRecordSetsInput{}
	maxItems := int32(1000)

	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		if equalQuals["name"].GetStringValue() != "" {
			input.StartRecordName = aws.String(equalQuals["name"].GetStringValue())
		}
	}
	if equalQuals["type"] != nil {
		// StartRecordType has a constraint that it must be used with StartRecordName
		if equalQuals["type"].GetStringValue() != "" && input.StartRecordName != nil {
			// input.StartRecordType = aws.String(equalQuals["type"].GetStringValue())
			input.StartRecordType.Values()
		}
	}

	// https://docs.aws.amazon.com/Route53/latest/APIReference/API_ListResourceRecordSets.html
	// The maximum/minimum record set per page is not mentioned in doc, so it has been set 1000 to max and 1 to min
	// Reduce the basic request limit down if the user has only requested a small number of rows
	// limit := d.QueryContext.Limit
	// if d.QueryContext.Limit != nil {
	// 	if *limit < 1000 {
	// 		if *limit < 1 {
	// 			input.MaxItems = aws.String("1")
	// 		} else {
	// 			input.MaxItems = aws.String(fmt.Sprint(*limit))
	// 		}
	// 	}
	// }

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

op, err := svc.ListResourceRecordSets(ctx, &input)
if err != nil{
	plugin.Logger(ctx).Error("aws_route53_record.listRoute53Records", "api_error", err)
		return nil, err
}

for _, record := range op.ResourceRecordSets {
				// The StartRecordName and StartRecordType input parameters only tell
				// the API where to start when returning results, so any records/types
				// that are greater in lexicographic order will also be returned.
				// Since Postgres will filter on exact matches anyway, check for exact
				// matches as an optimization to reduce the number of requests.

				if input.StartRecordName != nil && *record.Name != *input.StartRecordName {
					plugin.Logger(ctx).Debug("aws_route53_record.listRoute53Records mismatched record name", "input.StartRecordName", *input.StartRecordName, "record.Name", *record.Name)
					return nil, nil
				}

				if *record.Ty != *input.StartRecordType {
					plugin.Logger(ctx).Debug("aws_route53_record.listRoute53Records mismatched record type", "input.StartRecordType", *input.StartRecordType, "record.Type", *record.Type)
					return nil, nil
				}
				d.StreamListItem(ctx, &recordInfo{&hostedZoneID, record})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
		return nil, err
}

}

///
	// err = svc.ListResourceRecordSetsPages(
	// 	input,
	// 	func(page *route53.ListResourceRecordSetsOutput, isLast bool) bool {
	// 		for _, record := range page.ResourceRecordSets {
	// 			// The StartRecordName and StartRecordType input parameters only tell
	// 			// the API where to start when returning results, so any records/types
	// 			// that are greater in lexicographic order will also be returned.
	// 			// Since Postgres will filter on exact matches anyway, check for exact
	// 			// matches as an optimization to reduce the number of requests.

	// 			if input.StartRecordName != nil && *record.Name != *input.StartRecordName {
	// 				plugin.Logger(ctx).Debug("aws_route53_record.listRoute53Records mismatched record name", "input.StartRecordName", *input.StartRecordName, "record.Name", *record.Name)
	// 				return false
	// 			}

	// 			if input.StartRecordType != nil && *record.Type != *input.StartRecordType {
	// 				plugin.Logger(ctx).Debug("aws_route53_record.listRoute53Records mismatched record type", "input.StartRecordType", *input.StartRecordType, "record.Type", *record.Type)
	// 				return false
	// 			}

	// 			d.StreamListItem(ctx, &recordInfo{&hostedZoneID, record})

	// 			// Context may get cancelled due to manual cancellation or if the limit has been reached
	// 			if d.QueryStatus.RowsRemaining(ctx) == 0 {
	// 				return false
	// 			}
			// }
		// 	return !isLast
		// },

// 	notFoundErrors := []string{"InvalidParameter", "NoSuchHostedZone"}
// 	if err != nil {
// 		if helpers.StringSliceContains(notFoundErrors, err.(awserr.Error).Code()) {
// 			return nil, nil
// 		}
// 	}
// 	return nil, err
// }

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

func getRoute53RecordSetAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	recordData := h.Item.(*recordInfo)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Trace("aws_route53_record.getRoute53RecordSetAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":route53:::" +
		"hostedzone/" + *recordData.ZoneID +
		"/recordset/" + *recordData.Record.Name +
		"/" + *recordData.Record.Type

	if recordData.Record.SetIdentifier != nil {
		arn += "/" + *recordData.Record.SetIdentifier
	}

	// Get data for turbot defined properties
	akas := []string{arn}

	return akas, nil
}
