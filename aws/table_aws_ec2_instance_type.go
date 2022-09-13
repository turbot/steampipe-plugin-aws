package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInstanceType(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_instance_type",
		Description: "AWS EC2 Instance Type",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("instance_type"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidInstanceType"}),
			},
			Hydrate: describeInstanceType,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsInstanceTypesOfferings,
		},
		Columns: []*plugin.Column{
			{
				Name:        "instance_type",
				Description: "The instance type. For more information, see [ Instance Types ](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html) in the Amazon Elastic Compute Cloud User Guide.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].InstanceType"),
			},
			{
				Name:        "auto_recovery_supported",
				Description: "Indicates whether auto recovery is supported.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].AutoRecoverySupported"),
			},
			{
				Name:        "bare_metal",
				Description: "Indicates whether the instance is a bare metal instance type.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].BareMetal"),
			},
			{
				Name:        "burstable_performance_supported",
				Description: "Indicates whether the instance type is a burstable performance instance type.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].BurstablePerformanceSupported"),
			},
			{
				Name:        "current_generation",
				Description: "Indicates whether the instance type is current generation.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].CurrentGeneration"),
			},
			{
				Name:        "dedicated_hosts_supported",
				Description: "Indicates whether Dedicated Hosts are supported on the instance type.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].DedicatedHostsSupported"),
			},
			{
				Name:        "free_tier_eligible",
				Description: "Indicates whether the instance type is eligible for the free tier.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].FreeTierEligible"),
			},
			{
				Name:        "hibernation_supported",
				Description: "Indicates whether On-Demand hibernation is supported.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].HibernationSupported"),
			},
			{
				Name:        "hypervisor",
				Description: "The hypervisor for the instance type.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].Hypervisor"),
			},
			{
				Name:        "instance_storage_supported",
				Description: "Describes the instance storage for the instance type.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].InstanceStorageSupported"),
			},
			{
				Name:        "ebs_info",
				Description: "Describes the Amazon EBS settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].EbsInfo"),
			},
			{
				Name:        "memory_info",
				Description: "Describes the memory for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].MemoryInfo"),
			},
			{
				Name:        "network_info",
				Description: "Describes the network settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].NetworkInfo"),
			},
			{
				Name:        "placement_group_info",
				Description: "Describes the placement group settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].PlacementGroupInfo"),
			},
			{
				Name:        "processor_info",
				Description: "Describes the processor.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].ProcessorInfo"),
			},
			{
				Name:        "supported_root_device_types",
				Description: "The supported root device types.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].SupportedRootDeviceTypes"),
			},
			{
				Name:        "supported_usage_classes",
				Description: "Indicates whether the instance type is offered for spot or On-Demand.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].SupportedUsageClasses"),
			},
			{
				Name:        "supported_virtualization_types",
				Description: "The supported virtualization types.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].SupportedVirtualizationTypes"),
			},
			{
				Name:        "v_cpu_info",
				Description: "Describes the vCPU configurations for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].VCpuInfo"),
			},
			{
				Name:        "gpu_info",
				Description: "Describes the GPU accelerator settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].GpuInfo"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceTypes[0].InstanceType"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     instanceTypeDataToAkas,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listAwsInstanceTypesOfferings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// get the primary region for aws based on its partition
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	region := "us-east-1"
	if commonColumnData.Partition == "aws-us-gov" {
		region = "us-gov-east-1"
	} else if commonColumnData.Partition == "aws-cn" {
		region = "cn-north-1"
	} else if commonColumnData.Partition == "aws-iso" {
		region = "us-iso-east-1"
	} else if commonColumnData.Partition == "aws-iso-b" {
		region = "us-isob-east-1"
	}

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// First get all the types of instance
	input := &ec2.DescribeInstanceTypeOfferingsInput{
		LocationType: aws.String("region"),
		MaxResults:   aws.Int64(1000),
	}

	var filters []*ec2.Filter
	filters = append(filters, &ec2.Filter{Name: aws.String("location"), Values: []*string{&region}})
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
		plugin.Logger(ctx).Error("listAwsInstanceTypesOfferings", "InstanceType_DescribeInstanceTypeOfferingsPages", err)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func describeInstanceType(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var instanceType string
	if h.Item != nil {
		data := h.Item.(*ec2.InstanceTypeOffering)
		instanceType = types.SafeString(data.InstanceType)
	} else {
		instanceType = d.KeyColumnQuals["instance_type"].GetStringValue()
	}

	// get the primary region for aws based on its partition
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	region := "us-east-1"
	if commonColumnData.Partition == "aws-us-gov" {
		region = "us-gov-east-1"
	} else if commonColumnData.Partition == "aws-cn" {
		region = "cn-north-1"
	}

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// First get all the types of
	params := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []*string{aws.String(instanceType)},
	}

	// execute get call
	op, err := svc.DescribeInstanceTypes(params)
	if err != nil {
		plugin.Logger(ctx).Error("describeInstanceType", "DescribeInstanceTypes", err)
		return nil, err
	}

	return op, nil
}

func instanceTypeDataToAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("instanceTypeDataToAkas")
	var instanceType string
	switch h.Item.(type) {
	case *ec2.InstanceTypeOffering:
		instanceType = *h.Item.(*ec2.InstanceTypeOffering).InstanceType
	case *ec2.DescribeInstanceTypesOutput:
		instanceType = *h.Item.(*ec2.DescribeInstanceTypesOutput).InstanceTypes[0].InstanceType
	}

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:::instance-type/" + instanceType}

	return akas, nil
}
