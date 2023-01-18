package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInstanceAvailability(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_instance_availability",
		Description: "AWS EC2 Instance Availability",
		List: &plugin.ListConfig{
			ParentHydrate: listAwsRegions,
			Hydrate:       listAwsAvailableInstanceTypes,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "instance_type",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "instance_type",
				Description: "The instance type. For more information, see [ Instance Types ](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html) in the Amazon Elastic Compute Cloud User Guide.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The identifier for the location. This depends on the location type. For example, if the location type is region, the location is the Region code (for example, us-east-2.)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location_type",
				Description: "The type of location.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceType"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsInstanceAvailableAkas,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listAwsAvailableInstanceTypes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := h.Item.(types.Region)

	// If a region is not opted-in, we cannot list the availability zones
	if *region.OptInStatus == "not-opted-in" {
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(1000)
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

	// Create Session
	svc, err := EC2ClientForRegion(ctx, d, *region.RegionName)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance_availability.listAwsAvailableInstanceTypes", "connection_error", err)
		return nil, err
	}

	input := &ec2.DescribeInstanceTypeOfferingsInput{
		MaxResults:   aws.Int32(maxLimit),
		LocationType: types.LocationTypeRegion,
	}

	var filters []types.Filter
	filters = append(filters, types.Filter{Name: aws.String("location"), Values: []string{*region.RegionName}})

	equalQuals := d.KeyColumnQuals
	if equalQuals["instance_type"] != nil {
		filters = append(filters, types.Filter{Name: aws.String("instance-type"), Values: []string{equalQuals["instance_type"].GetStringValue()}})
	}
	input.Filters = filters

	paginator := ec2.NewDescribeInstanceTypeOfferingsPaginator(svc, input, func(o *ec2.DescribeInstanceTypeOfferingsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_instance_availability.listAwsAvailableInstanceTypes", "api_error", err)
			return nil, err
		}

		for _, items := range output.InstanceTypeOfferings {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// TRANSFORM FUNCTIONS

func getAwsInstanceAvailableAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceType := h.Item.(types.InstanceTypeOffering)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + *instanceType.Location + "::instance-type/" + string(instanceType.InstanceType)}
	return akas, nil
}
