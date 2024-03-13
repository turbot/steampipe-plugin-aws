package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAvailabilityZone(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_availability_zone",
		Description: "AWS Availability Zone",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "region_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			Hydrate: getAwsAvailabilityZone,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeAvailabilityZones"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsRegions,
			Hydrate:       listAwsAvailabilityZones,
			Tags:          map[string]string{"service": "ec2", "action": "DescribeAvailabilityZones"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "zone_id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Availability Zone, Local Zone, or Wavelength Zone",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ZoneName"),
			},
			{
				Name:        "zone_id",
				Description: "The ID of the Availability Zone, Local Zone, or Wavelength Zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "zone_type",
				Description: "The type of zone. The valid values are availability-zone, local-zone, and wavelength-zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "opt_in_status",
				Description: "For Availability Zones, this parameter always has the value of opt-in-not-required. For Local Zones and Wavelength Zones, this parameter is the opt-in status. The possible values are opted-in, and not-opted-in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "group_name",
				Description: "For Availability Zones, this parameter has the same value as the Region name. For Local Zones, the name of the associated group, for example us-west-2-lax-1. For Wavelength Zones, the name of the associated group, for example us-east-1-wl1-bos-wlz-1.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region_name",
				Description: "The name of the Region.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parent_zone_name",
				Description: "The name of the zone that handles some of the Local Zone or Wavelength Zone control plane operations, such as API calls.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parent_zone_id",
				Description: "The ID of the zone that handles some of the Local Zone or Wavelength Zone control plane operations, such as API calls",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "messages",
				Description: "Any messages about the Availability Zone, Local Zone, or Wavelength Zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ZoneName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAvailabilityZoneAkas,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listAwsAvailabilityZones(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := h.Item.(types.Region)

	// If a region is not opted-in, we cannot list the availability zones
	if *region.OptInStatus == "not-opted-in" {
		return nil, nil
	}

	// Create Session
	svc, err := EC2ClientForRegion(ctx, d, *region.RegionName)
	if err != nil {
		plugin.Logger(ctx).Error("aws_availability_zone.listAwsAvailabilityZones", "connection_error", err)
		return nil, err
	}

	input := &ec2.DescribeAvailabilityZonesInput{
		AllAvailabilityZones: aws.Bool(true),
		Filters: []types.Filter{
			{
				Name:   aws.String("region-name"),
				Values: []string{*region.RegionName},
			},
		},
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["zone_id"] != nil {
		input.ZoneIds = []string{equalQuals["zone_id"].GetStringValue()}
	}
	if equalQuals["name"] != nil {
		input.ZoneNames = []string{equalQuals["name"].GetStringValue()}
	}

	// execute list call
	resp, err := svc.DescribeAvailabilityZones(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_availability_zone.listAwsAvailabilityZones", "api_error", err)
		return nil, err
	}

	for _, zone := range resp.AvailabilityZones {
		d.StreamLeafListItem(ctx, zone)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsAvailabilityZone(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	regionName := d.EqualsQuals["region_name"].GetStringValue()

	// Create Session
	svc, err := EC2ClientForRegion(ctx, d, regionName)
	if err != nil {
		plugin.Logger(ctx).Error("aws_availability_zone.getAwsAvailabilityZone", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeAvailabilityZonesInput{
		AllAvailabilityZones: aws.Bool(true),
		ZoneNames:            []string{name},
	}

	// execute get call
	op, err := svc.DescribeAvailabilityZones(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_availability_zone.getAwsAvailabilityZone", "api_error", err)
		return nil, err
	}

	if len(op.AvailabilityZones) > 0 {
		return op.AvailabilityZones[0], nil
	}

	return nil, nil
}

func getAwsAvailabilityZoneAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zone := h.Item.(types.AvailabilityZone)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + "::" + *zone.RegionName + "::availability-zone/" + *zone.ZoneName}
	return akas, nil
}
