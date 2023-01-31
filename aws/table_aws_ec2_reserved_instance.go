package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2ReservedInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_reserved_instance",
		Description: "AWS EC2 Reserved Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("reserved_instance_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue", "InvalidInstanceID.Unavailable", "InvalidInstanceID.Malformed"}),
			},
			Hydrate: getEc2ReservedInstance,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2ReservedInstances,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "availability_zone", Require: plugin.Optional},
				{Name: "duration", Require: plugin.Optional},
				{Name: "end_time", Require: plugin.Optional},
				{Name: "fixed_price", Require: plugin.Optional},
				{Name: "instance_type", Require: plugin.Optional},
				{Name: "scope", Require: plugin.Optional},
				{Name: "product_description", Require: plugin.Optional},
				{Name: "start_time", Require: plugin.Optional},
				{Name: "instance_state", Require: plugin.Optional},
				{Name: "usage_price", Require: plugin.Optional},
				{Name: "offering_class", Require: plugin.Optional},
				{Name: "offering_type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
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

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_reserved_instance.listEc2ReservedInstances", "connection_error", err)
		return nil, err
	}

	input := &ec2.DescribeReservedInstancesInput{}

	filters := buildEc2ReservedInstanceFilter(d.Quals)

	equalQuals := d.KeyColumnQuals
	if equalQuals["offering_class"] != nil {
		input.OfferingClass = types.OfferingClassType(equalQuals["offering_class"].GetStringValue())
	}
	if equalQuals["offering_type"] != nil {
		input.OfferingType = types.OfferingTypeValues(equalQuals["offering_type"].GetStringValue())
	}

	if len(filters) != 0 {
		input.Filters = filters
	}

	// List call
	result, err := svc.DescribeReservedInstances(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_reserved_instance.listEc2ReservedInstances", "api_error", err)
		return nil, err
	}

	for _, reservedInstance := range result.ReservedInstances {
		d.StreamListItem(ctx, reservedInstance)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2ReservedInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	instanceID := d.KeyColumnQuals["reserved_instance_id"].GetStringValue()

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_reserved_instance.getEc2ReservedInstance", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeReservedInstancesInput{
		ReservedInstancesIds: []string{instanceID},
	}

	op, err := svc.DescribeReservedInstances(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_reserved_instance.getEc2ReservedInstance", "api_error", err)
		return nil, err
	}

	if op.ReservedInstances != nil && len(op.ReservedInstances) > 0 {
		return op.ReservedInstances[0], nil
	}
	return nil, nil
}

func getEc2ReservedInstanceARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(types.ReservedInstances)
	region := d.KeyColumnQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":reserved-instances/" + *instance.ReservedInstancesId

	return arn, nil
}

func getEc2ReservedInstanceModificationDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	instance := h.Item.(types.ReservedInstances)

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_reserved_instance.getEc2ReservedInstanceModificationDetails", "connection_error", err)
		return nil, err
	}

	filterName := "reserved-instances-id"
	filterValue := *instance.ReservedInstancesId

	param := &ec2.DescribeReservedInstancesModificationsInput{
		Filters: []types.Filter{
			{
				Name:   &filterName,
				Values: []string{filterValue},
			},
		},
	}

	res, err := svc.DescribeReservedInstancesModifications(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_reserved_instance.getEc2ReservedInstanceModificationDetails", "api_error", err)
		return nil, err
	}

	if res.ReservedInstancesModifications != nil || len(res.ReservedInstancesModifications) > 0 {
		return res.ReservedInstancesModifications, nil
	}

	return nil, err
}

//// TRANSFORM FUNCTION

func getEc2ReservedInstanceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(types.ReservedInstances)
	var turbotTagsMap map[string]string

	if instance.Tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range instance.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

// // UTILITY FUNCTION
// Build ec2 reserved instance list call input filter
func buildEc2ReservedInstanceFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"availability_zone":   "availability-zone",
		"duration":            "duration",
		"end_time":            "end",
		"fixed_price":         "fixed-price",
		"instance_type":       "instance-type",
		"scope":               "scope",
		"product_description": "product-description",
		"start_time":          "start",
		"usage_price":         "usage-price",
		"instance_state":      "state",
	}

	columnsDouble := []string{"fixed_price", "usage_price"}
	columnsInt := []string{"duration"}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			if strings.Contains(fmt.Sprint(columnsDouble), columnName) { //check Double columns
				value := getQualsValueByColumn(quals, columnName, "double")
				filter.Values = []string{fmt.Sprint(value)}
			} else if strings.Contains(fmt.Sprint(columnsInt), columnName) { //check Int columns
				value := getQualsValueByColumn(quals, columnName, "int64")
				filter.Values = []string{fmt.Sprint(value)}
			} else {
				value := getQualsValueByColumn(quals, columnName, "string")
				val, ok := value.(string)
				if ok {
					filter.Values = []string{val}
				}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
