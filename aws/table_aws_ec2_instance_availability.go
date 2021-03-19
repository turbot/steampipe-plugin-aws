package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInstanceAvailability(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_instance_availability",
		Description: "AWS EC2 Instance Availability",
		List: &plugin.ListConfig{
			ParentHydrate: listAwsRegions,
			Hydrate:       listAwsAvailableInstanceTypes,
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
	region := h.Item.(*ec2.Region)
	plugin.Logger(ctx).Trace("listAwsAvailableInstanceTypes", "region", *region.RegionName)

	// If a region is not opted-in, we cannot list the availability zones
	if types.SafeString(region.OptInStatus) == "not-opted-in" {
		return nil, nil
	}

	// Create Session
	svc, err := Ec2Service(ctx, d, *region.RegionName)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeInstanceTypeOfferingsInput{
		LocationType: aws.String("region"),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("location"),
				Values: []*string{region.RegionName},
			},
		},
	}

	// execute list call
	resp, err := svc.DescribeInstanceTypeOfferings(params)
	if err != nil {
		return nil, err
	}

	for _, zone := range resp.InstanceTypeOfferings {
		d.StreamLeafListItem(ctx, zone)
	}

	return nil, err
}

//// TRANSFORM FUNCTIONS

func getAwsInstanceAvailableAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsInstanceAvailableAkas")
	instanceType := h.Item.(*ec2.InstanceTypeOffering)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + *instanceType.Location + "::instance-type/" + *instanceType.InstanceType}
	return akas, nil
}
