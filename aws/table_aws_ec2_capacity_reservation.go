package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2CapacityReservation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_capacity_reservation",
		Description: "AWS EC2 Capacity Reservation",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("capacity_reservation_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidCapacityReservationId.NotFound", "InvalidCapacityReservationId.Unavailable", "InvalidCapacityReservationId.Malformed"}),
			Hydrate:           getEc2CapacityReservation,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2CapacityReservations,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "instance_type", Require: plugin.Optional},
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "availability_zone_id", Require: plugin.Optional},
				{Name: "availability_zone", Require: plugin.Optional},
				{Name: "instance_platform", Require: plugin.Optional},
				{Name: "tenancy", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "start_date", Require: plugin.Optional},
				{Name: "end_date", Require: plugin.Optional},
				{Name: "end_date_type", Require: plugin.Optional},
				{Name: "instance_match_criteria", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "capacity_reservation_id",
				Description: "The ID of the capacity reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capacity_reservation_arn",
				Description: "The Amazon Resource Name (ARN) of the capacity reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_type",
				Description: "The type of instance for which the capacity reservation reserves capacity.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the capacity reservation. A capacity reservation can be in one of the following states: 'active', 'expired', 'cancelled', 'pending', 'failed'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zone",
				Description: "The availability zone in which the capacity is reserved.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zone_id",
				Description: "The availability zone ID of the capacity reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "available_instance_count",
				Description: "The remaining capacity. Indicates the number of instances that can be launched in the capacity reservation.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "create_date",
				Description: "The date and time at which the capacity reservation was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "ebs_optimized",
				Description: "Indicates whether the capacity reservation supports EBS-optimized instances.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "end_date",
				Description: "The date and time at which the capacity reservation expires.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "end_date_type",
				Description: "Indicates the way in which the capacity reservation ends. A capacity reservation can have one of the following end types: 'unlimited', 'limited'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ephemeral_storage",
				Description: "Indicates whether the capacity reservation supports instances with temporary, block-level storage.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "instance_match_criteria",
				Description: "Indicates the type of instance launches that the capacity reservation accepts. The options include: 'open', 'targeted'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_platform",
				Description: "The type of operating system for which the capacity reservation reserves capacity.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account that owns the capacity reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_date",
				Description: "The date and time at which the capacity reservation was started.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tenancy",
				Description: "Indicates the tenancy of the capacity reservation. A capacity reservation can have one of the following tenancy settings: 'default', 'dedicated'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "total_instance_count",
				Description: "The total number of instances for which the capacity reservation reserves capacity.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "tag_src",
				Description: "Any tags assigned to the capacity reservation.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CapacityReservationId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(ec2CapacityReservationTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CapacityReservationArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2CapacityReservations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listEc2CapacityReservations", "AWS_REGION", region)

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeCapacityReservationsInput{
		MaxResults: aws.Int64(500),
	}

	filters := buildEc2CapacityReservationFilter(d.KeyColumnQuals)
	if len(filters) != 0 {
		input.Filters = filters
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 5 {
				input.MaxResults = aws.Int64(5)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.DescribeCapacityReservationsPages(
		input,
		func(page *ec2.DescribeCapacityReservationsOutput, isLast bool) bool {
			for _, reservation := range page.CapacityReservations {
				d.StreamListItem(ctx, reservation)

				// Context can be cancelled due to manual cancellation or the limit has been hit
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

func getEc2CapacityReservation(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2CapacityReservation")

	region := d.KeyColumnQualString(matrixKeyRegion)
	reservationId := d.KeyColumnQuals["capacity_reservation_id"].GetStringValue()

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeCapacityReservationsInput{
		CapacityReservationIds: []*string{aws.String(reservationId)},
	}

	op, err := svc.DescribeCapacityReservations(params)
	if err != nil {
		return nil, err
	}

	if op.CapacityReservations != nil && len(op.CapacityReservations) > 0 {
		return op.CapacityReservations[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func ec2CapacityReservationTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ec2CapacityReservationTagListToTurbotTags")
	tagList := d.Value.([]*ec2.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

//// UTILITY FUNCTION
// build ec2 capacity reservation list call input filter
func buildEc2CapacityReservationFilter(equalQuals plugin.KeyColumnEqualsQualMap) []*ec2.Filter {
	filters := make([]*ec2.Filter, 0)

	filterQuals := map[string]string{
		"instance_type":           "instance-type",
		"owner_id":                "owner-id",
		"availability_zone_id":    "availability-zone-id",
		"availability_zone":       "availability-zone",
		"instance_platform":       "instance-platform",
		"tenancy":                 "tenancy",
		"state":                   "state",
		"start_date":              "start-date",
		"end_date":                "end-date",
		"end_date_type":           "end-date-type",
		"instance_match_criteria": "instance_match_criteria",
	}

	for columnName, filterName := range filterQuals {
		if equalQuals[columnName] != nil {
			filter := ec2.Filter{
				Name: aws.String(filterName),
			}
			value := equalQuals[columnName]
			if value.GetStringValue() != "" {
				filter.Values = []*string{aws.String(equalQuals[columnName].GetStringValue())}
			} else if value.GetListValue() != nil {
				filter.Values = getListValues(value.GetListValue())
			}
			filters = append(filters, &filter)
		}
	}
	return filters
}
