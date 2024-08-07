package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInstanceType(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_instance_type",
		Description: "AWS EC2 Instance Type",
		Get: &plugin.GetConfig{
			// We must have to include the region in the query parameter to make the gate API call.
			// Otherwise we will get an Error: get call returned 9 results - the key column is not globally unique (SQLSTATE HV000)
			KeyColumns: plugin.AllColumns([]string{"instance_type", "region"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInstanceType"}),
			},
			Hydrate: describeInstanceType,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeInstanceTypes"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsInstanceTypesOfferings,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "instance_type", Require: plugin.Optional, Operators: []string{"="}},
			},
			Tags:    map[string]string{"service": "ec2", "action": "DescribeInstanceTypeOfferings"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: describeInstanceType,
				Tags: map[string]string{"service": "ec2", "action": "DescribeInstanceTypes"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "instance_type",
				Description: "The instance type. For more information, see [ Instance Types ](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html) in the Amazon Elastic Compute Cloud User Guide.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "auto_recovery_supported",
				Description: "Indicates whether auto recovery is supported.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "bare_metal",
				Description: "Indicates whether the instance is a bare metal instance type.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "burstable_performance_supported",
				Description: "Indicates whether the instance type is a burstable performance instance type.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "current_generation",
				Description: "Indicates whether the instance type is current generation.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "dedicated_hosts_supported",
				Description: "Indicates whether Dedicated Hosts are supported on the instance type.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "free_tier_eligible",
				Description: "Indicates whether the instance type is eligible for the free tier.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "nitro_enclaves_support",
				Description: "Indicates whether Nitro Enclaves is supported.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "nitro_tpm_support",
				Description: "Indicates whether NitroTPM is supported.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "hibernation_supported",
				Description: "Indicates whether On-Demand hibernation is supported.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "nitro_enclaves_support",
				Description: "Indicates whether instance storage is supported.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "nitro_tpm_support",
				Description: "Indicates whether NitroTPM is supported.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "hypervisor",
				Description: "The hypervisor for the instance type.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "instance_storage_supported",
				Description: "Describes the instance storage for the instance type.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "ebs_info",
				Description: "Describes the Amazon EBS settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "memory_info",
				Description: "Describes the memory for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "network_info",
				Description: "Describes the network settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "placement_group_info",
				Description: "Describes the placement group settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "processor_info",
				Description: "Describes the processor.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "supported_root_device_types",
				Description: "The supported root device types.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "supported_usage_classes",
				Description: "Indicates whether the instance type is offered for spot or On-Demand.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "supported_virtualization_types",
				Description: "The supported virtualization types.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "v_cpu_info",
				Description: "Describes the vCPU configurations for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("VCpuInfo"),
			},
			{
				Name:        "gpu_info",
				Description: "Describes the GPU accelerator settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "fpga_info",
				Description: "Describes the FPGA accelerator settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "inference_accelerator_info",
				Description: "Describes the Inference accelerator settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "instance_storage_info",
				Description: "Describes the instance storage for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "media_accelerator_info",
				Description: "Describes the media accelerator settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "neuron_info",
				Description: "Describes the Neuron accelerator settings for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "nitro_tpm_info",
				Description: "Describes the supported NitroTPM versions for the instance type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "supported_boot_modes",
				Description: "The supported boot modes.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeInstanceType,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeInstanceType,
				Transform:   transform.FromField("InstanceType"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     instanceTypeDataToAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsInstanceTypesOfferings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance_type.listAwsInstanceTypesOfferings", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
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

	// First get all the types of instance
	input := &ec2.DescribeInstanceTypeOfferingsInput{
		LocationType: types.LocationTypeRegion,
		MaxResults:   aws.Int32(maxLimit),
	}

	var filters []types.Filter
	filters = append(filters, types.Filter{Name: aws.String("location"), Values: []string{region}})
	if d.EqualsQualString("instance_type") != "" {
		filters = append(filters, types.Filter{Name: aws.String("instance-type"), Values: []string{d.EqualsQualString("instance_type")}})
	}

	input.Filters = filters

	paginator := ec2.NewDescribeInstanceTypeOfferingsPaginator(svc, input, func(o *ec2.DescribeInstanceTypeOfferingsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_instance_type.listAwsInstanceTypesOfferings", "api_error", err)
			return nil, err
		}

		for _, items := range output.InstanceTypeOfferings {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func describeInstanceType(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var instanceType types.InstanceType
	if h.Item != nil {
		data := h.Item.(types.InstanceTypeOffering)
		instanceType = data.InstanceType
	} else {
		instanceType = types.InstanceType(d.EqualsQuals["instance_type"].GetStringValue())
	}

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance_type.describeInstanceType", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// First get all the types of
	params := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []types.InstanceType{
			instanceType,
		},
	}

	// execute get call
	op, err := svc.DescribeInstanceTypes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance_type.describeInstanceType", "api_error", err)
		return nil, err
	}
	if len(op.InstanceTypes) > 0 {
		return op.InstanceTypes[0], nil
	}

	return nil, nil
}

func instanceTypeDataToAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var instanceType types.InstanceType
	switch h.Item.(type) {
	case types.InstanceTypeOffering:
		instanceType = h.Item.(types.InstanceTypeOffering).InstanceType
	case types.InstanceTypeInfo:
		instanceType = h.Item.(types.InstanceTypeInfo).InstanceType
	}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance_type.instanceTypeDataToAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{fmt.Sprintf("arn:%s:ec2:::instance-type/%s", commonColumnData.Partition, instanceType)}

	return akas, nil
}
