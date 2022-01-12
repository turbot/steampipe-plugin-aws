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

	input := &ec2.DescribeInstanceTypeOfferingsInput{
		MaxResults:   aws.Int64(1000),
		LocationType: aws.String("region"),
	}

	var filters []*ec2.Filter
	filters = append(filters, &ec2.Filter{Name: aws.String("location"), Values: []*string{region.RegionName}})

	equalQuals := d.KeyColumnQuals
	if equalQuals["instance_type"] != nil {
		filters = append(filters, &ec2.Filter{Name: aws.String("instance-type"), Values: []*string{aws.String(equalQuals["instance_type"].GetStringValue())}})
	}
	input.Filters = filters

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
	err = svc.DescribeInstanceTypeOfferingsPages(
		input,
		func(page *ec2.DescribeInstanceTypeOfferingsOutput, isLast bool) bool {
			for _, instanceTypeOffering := range page.InstanceTypeOfferings {
				d.StreamListItem(ctx, instanceTypeOffering)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listAwsAvailableInstanceTypes", "DescribeInstanceTypeOfferingsPages", err)
	}

	return nil, err
}

//// TRANSFORM FUNCTIONS

func getAwsInstanceAvailableAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsInstanceAvailableAkas")
	instanceType := h.Item.(*ec2.InstanceTypeOffering)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + *instanceType.Location + "::instance-type/" + *instanceType.InstanceType}
	return akas, nil
}
