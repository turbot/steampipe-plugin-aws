package aws

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

//// TABLE DEFINITION

func tableAwsInstanceType(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_instance_type",
		Description: "AWS EC2 Instance Type",
		List: &plugin.ListConfig{
			Hydrate: listAwsInstanceTypesOfferings,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "instance_type", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "instance_type_pattern", Require: plugin.Optional, Operators: []string{"="}, CacheMatch: query_cache.CacheMatchExact},
			},
			Tags: map[string]string{"service": "ec2", "action": "DescribeInstanceTypeOfferings"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EC2ServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "instance_type",
				Description: "The instance type. For more information, see [ Instance Types ](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html) in the Amazon Elastic Compute Cloud User Guide.",
				Type:        proto.ColumnType_STRING,
			},
			// In the query "select * from aws_ec2_instance_type where instance_type = 't2*'", the API fetches the result but returns empty rows due to PostgreSQL-level filtering.
			// The 'instance_type' column contains values like 't2.small', not 't2*', leading to a mismatch between the column value ('t2.small') and the wildcard pattern used in the query ('t2*').
			// To resolve this issue, the 'instance_type_pattern' column has been added, allowing for proper filtering using wildcard patterns.
			{
				Name:        "instance_type_pattern",
				Description: "The instance type pattern includes wildcards, such as 'c5-*', 't2*', and 'm5*'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("instance_type_pattern"),
			},
			{
				Name:        "auto_recovery_supported",
				Description: "Indicates whether auto recovery is supported.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "bare_metal",
				Description: "Indicates whether the instance is a bare metal instance type.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "burstable_performance_supported",
				Description: "Indicates whether the instance type is a burstable performance instance type.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "current_generation",
				Description: "Indicates whether the instance type is current generation.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "dedicated_hosts_supported",
				Description: "Indicates whether Dedicated Hosts are supported on the instance type.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "free_tier_eligible",
				Description: "Indicates whether the instance type is eligible for the free tier.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "nitro_enclaves_support",
				Description: "Indicates whether Nitro Enclaves is supported.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "nitro_tpm_support",
				Description: "Indicates whether NitroTPM is supported.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "hibernation_supported",
				Description: "Indicates whether On-Demand hibernation is supported.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "hypervisor",
				Description: "The hypervisor for the instance type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_storage_supported",
				Description: "Describes the instance storage for the instance type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ebs_info",
				Description: "Describes the Amazon EBS settings for the instance type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "location_type",
				Description: "Type of the location, for example: region, availability-zone, etc.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "memory_info",
				Description: "Describes the memory for the instance type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_info",
				Description: "Describes the network settings for the instance type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "placement_group_info",
				Description: "Describes the placement group settings for the instance type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "processor_info",
				Description: "Describes the processor.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "supported_root_device_types",
				Description: "The supported root device types.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "supported_usage_classes",
				Description: "Indicates whether the instance type is offered for spot or On-Demand.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "supported_virtualization_types",
				Description: "The supported virtualization types.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "v_cpu_info",
				Description: "Describes the vCPU configurations for the instance type.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VCpuInfo"),
			},
			{
				Name:        "gpu_info",
				Description: "Describes the GPU accelerator settings for the instance type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "fpga_info",
				Description: "Describes the FPGA accelerator settings for the instance type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "inference_accelerator_info",
				Description: "Describes the Inference accelerator settings for the instance type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "instance_storage_info",
				Description: "Describes the instance storage for the instance type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "media_accelerator_info",
				Description: "Describes the media accelerator settings for the instance type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "neuron_info",
				Description: "Describes the Neuron accelerator settings for the instance type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "nitro_tpm_info",
				Description: "Describes the supported NitroTPM versions for the instance type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "supported_boot_modes",
				Description: "The supported boot modes.",
				Type:        proto.ColumnType_JSON,
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
				Hydrate:     instanceTypeDataToAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsInstanceTypesOfferings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	// Create EC2 client
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance_type.listAwsInstanceTypesOfferings", "connection_error", err)
		return nil, err
	}

	// If the service is nil, the region is unsupported
	if svc == nil {
		return nil, nil
	}

	// Set the maximum limit for results
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

	// Prepare the input for EC2 DescribeInstanceTypeOfferings API
	input := &ec2.DescribeInstanceTypeOfferingsInput{
		LocationType: types.LocationTypeRegion,
		MaxResults:   aws.Int32(maxLimit),
		Filters:      []types.Filter{{Name: aws.String("location"), Values: []string{region}}},
	}

	// Fetch instance type offering for a particular instance type
	if d.EqualsQualString("instance_type") != "" {
		input.Filters = append(input.Filters, types.Filter{Name: aws.String("instance-type"), Values: []string{d.EqualsQualString("instance_type")}})
	}

	// Fetch instance types offerings
	instanceTypes, err := fetchInstanceTypeOfferings(ctx, d, svc, input, maxLimit)
	if err != nil {
		return nil, err
	}

	// Apply pattern matching on instance types if provided in query
	filteredInstanceTypes := filterInstanceTypesByPattern(ctx, instanceTypes, d.EqualsQualString("instance_type_pattern"))

	// Batch process the instance types in groups of 100
	err = batchDescribeInstanceTypes(ctx, d, filteredInstanceTypes, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance_type.listAwsInstanceTypesOfferings", "batch_process_error", err)
		return nil, err
	}

	return nil, nil
}

// Helper function to fetch instance type offerings using pagination
func fetchInstanceTypeOfferings(ctx context.Context, d *plugin.QueryData, svc *ec2.Client, input *ec2.DescribeInstanceTypeOfferingsInput, maxLimit int32) ([]types.InstanceType, error) {
	var instanceTypes []types.InstanceType

	paginator := ec2.NewDescribeInstanceTypeOfferingsPaginator(svc, input, func(o *ec2.DescribeInstanceTypeOfferingsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_instance_type.fetchInstanceTypeOfferings", "api_error", err)
			return nil, err
		}

		for _, item := range output.InstanceTypeOfferings {
			instanceTypes = append(instanceTypes, item.InstanceType)
		}
	}

	return instanceTypes, nil
}

// Helper function to filter instance types by a pattern like t2-*, m5-*, etc.
func filterInstanceTypesByPattern(_ context.Context, instanceTypes []types.InstanceType, pattern string) []types.InstanceType {
	if pattern == "" {
		return instanceTypes
	}

	// The regex pattern "t3*" does not work as expected when matching the string "t3.small". This is because '*' in regex matches zero or more occurrences of the preceding character,
	// allowing it to match "t3.small" due to the presence of the '.' character. To correct this, we replace '*' with '.+' to match any characters after "t3" appropriately.
	// The following patterns were validated:
	// - "c7*" matches (e.g., c7i.2xlarge, c7gn.xlarge, c7i-flex.2xlarge),
	// - "c7i*" matches (e.g., c7i.2xlarge, c7i-flex.2xlarge),
	// - "c7i.*" matches (e.g., c7i.8xlarge, c7i.2xlarge),
	// - "c7i-*" matches (e.g., c7i-flex.2xlarge).
	pattern = strings.ReplaceAll(pattern, ".", "\\.")
	pattern = strings.ReplaceAll(pattern, "*", ".+")
	var matchedInstanceTypes []types.InstanceType
	re := regexp.MustCompile(pattern)

	for _, instanceType := range instanceTypes {
		if re.MatchString(string(instanceType)) {
			matchedInstanceTypes = append(matchedInstanceTypes, instanceType)
		}
	}

	return matchedInstanceTypes
}

// Helper function to batch describe instance types in groups of 100
func batchDescribeInstanceTypes(ctx context.Context, d *plugin.QueryData, instanceTypes []types.InstanceType, region string) error {
	batchSize := 100
	if d.QueryContext.Limit != nil && *d.QueryContext.Limit > 0 && *d.QueryContext.Limit < 100 {
		batchSize = int(*d.QueryContext.Limit)
	}
	for i := 0; i < len(instanceTypes); i += batchSize {
		end := i + batchSize
		if end > len(instanceTypes) {
			end = len(instanceTypes)
		}
		batch := instanceTypes[i:end]

		err := describeInstanceTypes(ctx, d, batch, region)
		if err != nil {
			return err
		}
	}

	return nil
}

// Describe instance types and stream the results
func describeInstanceTypes(ctx context.Context, d *plugin.QueryData, instanceTypes []types.InstanceType, region string) error {
	svc, err := EC2ClientForRegion(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance_type.describeInstanceTypes", "connection_error", err)
		return err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil
	}

	// Create input for DescribeInstanceTypes API
	params := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: instanceTypes,
	}

	// Fetch the instance types
	op, err := svc.DescribeInstanceTypes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance_type.describeInstanceTypes", "api_error", err)
		return err
	}

	// Stream each item from the response
	for _, item := range op.InstanceTypes {
		d.StreamListItem(ctx, item)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil
		}
	}

	return nil
}

//// HYDRATE FUNCTION

func instanceTypeDataToAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceType := h.Item.(types.InstanceTypeInfo).InstanceType

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance_type.instanceTypeDataToAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{fmt.Sprintf("arn:%s:ec2:::instance-type/%s", commonColumnData.Partition, instanceType)}

	return akas, nil
}
