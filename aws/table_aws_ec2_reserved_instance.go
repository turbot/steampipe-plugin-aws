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

func tableAwsEc2ReservedInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_reserved_instance",
		Description: "AWS EC2 Reserved Instance",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("reserved_instance_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidParameterValue", "InvalidInstanceID.Unavailable", "InvalidInstanceID.Malformed"}),
			Hydrate:           getEc2ReservedInstance,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2ReservedInstances,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "reserved_instance_id",
				Description: "The ID of the Reserved instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReservedInstancesId"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEc2ReservedInstanceARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "instance_type",
				Description: "The instance type on which the Reserved Instance can be used.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_state",
				Description: "The state of the Reserved Instance purchase.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State"),
			},
			{
				Name:        "availability_zone",
				Description: "The Availability Zone in which the Reserved Instance can be used.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "currency_code",
				Description: "The currency of the Reserved Instance. It's specified using ISO 4217 standard currency codes. At this time, the only supported currency is USD.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "duration",
				Description: "The duration of the Reserved Instance, in seconds.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "end_time",
				Description: "The time when the Reserved Instance expires.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("End"),
			},
			{
				Name:        "fixed_price",
				Description: "The purchase price of the Reserved Instance.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "instance_count",
				Description: "The number of reservations purchased.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "instance_tenancy",
				Description: "The tenancy of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "offering_class",
				Description: "The offering class of the Reserved Instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "offering_type",
				Description: "The Reserved Instance offering type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "product_description",
				Description: "The Reserved Instance product platform description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scope",
				Description: "The scope of the Reserved Instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_time",
				Description: "The date and time the Reserved Instance started.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Start"),
			},
			{
				Name:        "usage_price",
				Description: "The usage price of the Reserved Instance, per hour.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "reserved_instances_modifications",
				Description: "The Reserved Instance modification information.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEc2ReservedInstanceModificationDetails,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the reserved instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReservedInstancesId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getEc2ReservedInstanceTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEc2ReservedInstanceARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2ReservedInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listEc2ReservedInstances", "AWS_REGION", region)

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	param := &ec2.DescribeReservedInstancesInput{}
	// List call
	result, err := svc.DescribeReservedInstances(param)
	if err != nil {
		return nil, err
	}

	for _, reservedInstance := range result.ReservedInstances {
		d.StreamListItem(ctx, reservedInstance)
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2ReservedInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2ReservedInstance")

	region := d.KeyColumnQualString(matrixKeyRegion)
	instanceID := d.KeyColumnQuals["reserved_instance_id"].GetStringValue()

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeReservedInstancesInput{
		ReservedInstancesIds: []*string{aws.String(instanceID)},
	}

	op, err := svc.DescribeReservedInstances(params)
	if err != nil {
		return nil, err
	}

	if op.ReservedInstances != nil && len(op.ReservedInstances) > 0 {
		return op.ReservedInstances[0], nil
	}
	return nil, nil
}

func getEc2ReservedInstanceARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2ReservedInstanceARN")
	instance := h.Item.(*ec2.ReservedInstances)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":reserved-instances/" + *instance.ReservedInstancesId

	return arn, nil
}

func getEc2ReservedInstanceModificationDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2ReservedInstanceModificationDetails")

	instance := h.Item.(*ec2.ReservedInstances)
	region := d.KeyColumnQualString(matrixKeyRegion)

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	filterName := "reserved-instances-id"
	filterValue := *instance.ReservedInstancesId

	param := &ec2.DescribeReservedInstancesModificationsInput{
		Filters: []*ec2.Filter{
			{
				Name:   &filterName,
				Values: []*string{&filterValue},
			},
		},
	}

	res, err := svc.DescribeReservedInstancesModifications(param)
	if err != nil {
		return nil, err
	}

	if res.ReservedInstancesModifications != nil || len(res.ReservedInstancesModifications) > 0 {
		return res.ReservedInstancesModifications, nil
	}

	return nil, err
}

//// TRANSFORM FUNCTION

func getEc2ReservedInstanceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(*ec2.ReservedInstances)
	return ec2TagsToMap(instance.Tags)
}
