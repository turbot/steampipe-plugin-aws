package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInstanceType(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_instance_type",
		Description: "AWS EC2 Instance Type",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("instance_type"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidInstanceType"}),
			Hydrate:           describeInstanceType,
			ItemFromKey:       instanceTypeOfferingFromKey,
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
				Transform:   transform.From(instanceTypeDataToAkas),
			},
		},
	}
}

// custom struct for InstanceTypeOffering
type instanceTypeOfferingInfo struct {
	Partition            string
	Region               string
	InstanceTypeOffering *ec2.InstanceTypeOffering
}

//// ITEM FROM KEY

func instanceTypeOfferingFromKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	instanceType := quals["instance_type"].GetStringValue()

	// get the primary region for aws based on its partition
	commonData, err := getCommonColumns(ctx, d, h)
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

	item := &instanceTypeOfferingInfo{
		Partition: commonColumnData.Partition,
		Region:    region,
		InstanceTypeOffering: &ec2.InstanceTypeOffering{
			InstanceType: &instanceType,
		},
	}

	return item, nil
}

//// LIST FUNCTION

func listAwsInstanceTypesOfferings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// get the primary region for aws based on its partition
	commonData, err := getCommonColumns(ctx, d, h)
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
	svc, err := Ec2Service(ctx, d.ConnectionManager, region)
	if err != nil {
		return nil, err
	}

	// First get all the types of
	params := &ec2.DescribeInstanceTypeOfferingsInput{
		LocationType: aws.String("region"),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("location"),
				Values: []*string{aws.String(region)},
			},
		},
	}

	// execute list call
	resp, err := svc.DescribeInstanceTypeOfferings(params)
	if err != nil {
		return nil, err
	}

	for _, instanceType := range resp.InstanceTypeOfferings {
		d.StreamListItem(ctx, &instanceTypeOfferingInfo{
			commonColumnData.Partition,
			region,
			instanceType,
		})
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func describeInstanceType(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceInfo := h.Item.(*instanceTypeOfferingInfo)

	// Create Session
	svc, err := Ec2Service(ctx, d.ConnectionManager, instanceInfo.Region)
	if err != nil {
		return nil, err
	}

	// First get all the types of
	params := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []*string{instanceInfo.InstanceTypeOffering.InstanceType},
	}

	// execute list call
	op, err := svc.DescribeInstanceTypes(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func instanceTypeDataToAkas(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("instanceTypeDataToAkas")
	var data *ec2.DescribeInstanceTypesOutput
	akas := []string{}

	if d.HydrateResults["describeInstanceType"] != nil {
		data = d.HydrateResults["describeInstanceType"].(*ec2.DescribeInstanceTypesOutput)
		akas = []string{"arn:aws:ec2:::instance-type/" + *data.InstanceTypes[0].InstanceType}
	}

	return akas, nil
}
