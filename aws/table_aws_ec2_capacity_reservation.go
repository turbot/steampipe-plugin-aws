package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2Endpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2CapacityReservation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_capacity_reservation",
		Description: "AWS EC2 Capacity Reservation",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("capacity_reservation_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidCapacityReservationId.NotFound", "InvalidCapacityReservationId.Unavailable", "InvalidCapacityReservationId.Malformed"}),
			},
			Hydrate: getEc2CapacityReservation,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeCapacityReservations"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2CapacityReservations,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeCapacityReservations"},
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
		GetMatrixItemFunc: SupportedRegionMatrix(ec2Endpoint.AWS_EC2_SERVICE_ID),
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
				Name:        "capacity_reservation_fleet_id",
				Description: "The ID of the Capacity Reservation Fleet to which the Capacity Reservation belongs. Only valid for Capacity Reservations that were created by a Capacity Reservation Fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "outpost_arn",
				Description: "The Amazon Resource Name (ARN) of the Outpost on which the Capacity Reservation was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "placement_group_arn",
				Description: "The Amazon Resource Name (ARN) of the cluster placement group in which the Capacity Reservation was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reservation_type",
				Description: "The type of Capacity Reservation.",
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
				Name:        "capacity_allocations",
				Description: "Information about instance capacity usage.",
				Type:        proto.ColumnType_JSON,
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

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_capacity_reservation.listEc2CapacityReservations", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				maxLimit = 5
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeCapacityReservationsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filters := buildEc2CapacityReservationFilter(d.Quals)
	if len(filters) != 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeCapacityReservationsPaginator(svc, input, func(o *ec2.DescribeCapacityReservationsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_capacity_reservation.listEc2CapacityReservations", "api_error", err)
			return nil, err
		}

		for _, items := range output.CapacityReservations {
			d.StreamListItem(ctx, items)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2CapacityReservation(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	reservationId := d.EqualsQuals["capacity_reservation_id"].GetStringValue()

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_capacity_reservation.getEc2CapacityReservation", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeCapacityReservationsInput{
		CapacityReservationIds: []string{reservationId},
	}

	op, err := svc.DescribeCapacityReservations(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_capacity_reservation.getEc2CapacityReservation", "api_error", err)
		return nil, err
	}

	if len(op.CapacityReservations) > 0 {
		return op.CapacityReservations[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func ec2CapacityReservationTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.Value.([]types.Tag)

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

// Build ec2 capacity reservation list call input filter
func buildEc2CapacityReservationFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

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
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
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
