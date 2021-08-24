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
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "capacity_reservation_arn",
				Description: "The Amazon Resource Name (ARN) of the capacity reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capacity_reservation_id",
				Description: "The ID of the capacity reservation.",
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
				Description: "The total number of instances for which the capacity reservation reserves capacity",
				Type:        proto.ColumnType_INT,
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
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2CapacityReservationAkas,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
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
	
	// List call
	err = svc.DescribeCapacityReservationsPages(
		&ec2.DescribeCapacityReservationsInput{},
		func(page *ec2.DescribeCapacityReservationsOutput, isLast bool) bool {
			for _, reservation := range page.CapacityReservations {
				d.StreamListItem(ctx, reservation)
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

func getAwsEc2CapacityReservationAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2CapacityReservationAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	reservation := h.Item.(*ec2.CapacityReservation)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + *reservation.OwnerId + ":capacity-reservation/" + *reservation.CapacityReservationId}

	return akas, nil
}
