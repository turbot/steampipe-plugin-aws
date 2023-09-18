package aws

import (
	"context"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/go-kit/helpers"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2Instance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_instance",
		Description: "AWS EC2 Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("instance_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInstanceID.NotFound", "InvalidInstanceID.Unavailable", "InvalidInstanceID.Malformed"}),
			},
			Hydrate: getEc2Instance,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeInstances"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2Instance,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeInstances"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "hypervisor", Require: plugin.Optional},
				{Name: "iam_instance_profile_arn", Require: plugin.Optional},
				{Name: "image_id", Require: plugin.Optional},
				{Name: "instance_lifecycle", Require: plugin.Optional},
				{Name: "instance_state", Require: plugin.Optional},
				{Name: "instance_type", Require: plugin.Optional},
				{Name: "monitoring_state", Require: plugin.Optional},
				{Name: "outpost_arn", Require: plugin.Optional},
				{Name: "placement_availability_zone", Require: plugin.Optional},
				{Name: "placement_group_name", Require: plugin.Optional},
				{Name: "public_dns_name", Require: plugin.Optional},
				{Name: "ram_disk_id", Require: plugin.Optional},
				{Name: "root_device_name", Require: plugin.Optional},
				{Name: "root_device_type", Require: plugin.Optional},
				{Name: "subnet_id", Require: plugin.Optional},
				{Name: "placement_tenancy", Require: plugin.Optional},
				{Name: "virtualization_type", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getInstanceDisableAPITerminationData,
				Tags: map[string]string{"service": "ec2", "action": "DescribeInstanceAttribute"},
			},
			{
				Func: getInstanceInitiatedShutdownBehavior,
				Tags: map[string]string{"service": "ec2", "action": "DescribeInstanceAttribute"},
			},
			{
				Func: getInstanceKernelID,
				Tags: map[string]string{"service": "ec2", "action": "DescribeInstanceAttribute"},
			},
			{
				Func: getInstanceRAMDiskID,
				Tags: map[string]string{"service": "ec2", "action": "DescribeInstanceAttribute"},
			},
			{
				Func: getInstanceSriovNetSupport,
				Tags: map[string]string{"service": "ec2", "action": "DescribeInstanceAttribute"},
			},
			{
				Func: getInstanceUserData,
				Tags: map[string]string{"service": "ec2", "action": "DescribeInstanceAttribute"},
			},
			{
				Func: getEc2LaunchTemplateData,
				Tags: map[string]string{"service": "ec2", "action": "GetLaunchTemplateData"},
			},
			{
				Func: getInstanceStatus,
				Tags: map[string]string{"service": "ec2", "action": "DescribeInstanceStatus"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "instance_id",
				Description: "The ID of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEc2InstanceARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "instance_type",
				Description: "The instance type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_state",
				Description: "The state of the instance (pending | running | shutting-down | terminated | stopping | stopped).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State.Name"),
			},
			{
				Name:        "monitoring_state",
				Description: "Indicates whether detailed monitoring is enabled (disabled | enabled).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Monitoring.State"),
			},
			{
				Name:        "disable_api_termination",
				Default:     false,
				Description: "If the value is true, instance can't be terminated through the Amazon EC2 console, CLI, or API.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getInstanceDisableAPITerminationData,
				Transform:   transform.FromField("DisableApiTermination.Value"),
			},
			{
				Name:        "ami_launch_index",
				Description: "The AMI launch index, which can be used to find this instance in the launch group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "architecture",
				Description: "The architecture of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "boot_mode",
				Description: "The boot mode of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capacity_reservation_id",
				Description: "The ID of the Capacity Reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capacity_reservation_specification",
				Description: "Information about the Capacity Reservation targeting option.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "client_token",
				Description: "The idempotency token you provided when you launched the instance, if applicable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cpu_options_core_count",
				Description: "The number of CPU cores for the instance.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CpuOptions.CoreCount"),
			},
			{
				Name:        "cpu_options_threads_per_core",
				Description: "The number of threads per CPU core.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CpuOptions.ThreadsPerCore"),
			},
			{
				Name:        "ebs_optimized",
				Description: "Indicates whether the instance is optimized for Amazon EBS I/O. This optimization provides dedicated throughput to Amazon EBS and an optimized configuration stack to provide optimal I/O performance. This optimization isn't available with all instance types.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "ena_support",
				Description: "Specifies whether enhanced networking with ENA is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "hypervisor",
				Description: "The hypervisor type of the instance. The value xen is used for both Xen and Nitro hypervisors.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_instance_profile_arn",
				Description: "The Amazon Resource Name (ARN) of IAM instance profile associated with the instance, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IamInstanceProfile.Arn"),
			},
			{
				Name:        "iam_instance_profile_id",
				Description: "The ID of the instance profile associated with the instance, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IamInstanceProfile.Id"),
			},
			{
				Name:        "image_id",
				Description: "The ID of the AMI used to launch the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_initiated_shutdown_behavior",
				Description: "Indicates whether an instance stops or terminates when you initiate shutdown from the instance (using the operating system command for system shutdown).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInstanceInitiatedShutdownBehavior,
				Transform:   transform.FromField("InstanceInitiatedShutdownBehavior.Value"),
			},
			{
				Name:        "instance_lifecycle",
				Description: "Indicates whether this is a spot instance or a scheduled instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kernel_id",
				Description: "The kernel ID",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInstanceKernelID,
				Transform:   transform.FromField("KernelId.Value"),
			},
			{
				Name:        "key_name",
				Description: "The name of the key pair, if this instance was launched with an associated key pair.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "launch_time",
				Description: "The time the instance was launched.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "outpost_arn",
				Description: "The Amazon Resource Name (ARN) of the Outpost, if applicable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "placement_affinity",
				Description: "The affinity setting for the instance on the Dedicated Host.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.Affinity"),
			},
			{
				Name:        "placement_availability_zone",
				Description: "The Availability Zone of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.AvailabilityZone"),
			},
			{
				Name:        "placement_group_id",
				Description: "The ID of the placement group that the instance is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.GroupId"),
			},
			{
				Name:        "placement_group_name",
				Description: "The name of the placement group the instance is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.GroupName"),
			},
			{
				Name:        "placement_host_id",
				Description: "The ID of the Dedicated Host on which the instance resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.HostId"),
			},
			{
				Name:        "placement_host_resource_group_arn",
				Description: "The ARN of the host resource group in which to launch the instances.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.HostResourceGroupArn"),
			},
			{
				Name:        "placement_partition_number",
				Description: "The ARN of the host resource group in which to launch the instances.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Placement.PartitionNumber"),
			},
			{
				Name:        "placement_tenancy",
				Description: "The tenancy of the instance (if the instance is running in a VPC). An instance with a tenancy of dedicated runs on single-tenant hardware.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.Tenancy"),
			},
			{
				Name:        "platform",
				Description: "The value is 'Windows' for Windows instances; otherwise blank.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "platform_details",
				Description: "The platform details value for the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_ip_address",
				Description: "The private IPv4 address assigned to the instance.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "private_dns_name",
				Description: "The private DNS hostname name assigned to the instance. This DNS hostname can only be used inside the Amazon EC2 network. This name is not available until the instance enters the running state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_dns_name",
				Description: "The public DNS name assigned to the instance. This name is not available until the instance enters the running state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_ip_address",
				Description: "The public IPv4 address, or the Carrier IP address assigned to the instance, if applicable.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "ram_disk_id",
				Description: "The RAM disk ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInstanceRAMDiskID,
				Transform:   transform.FromField("RamdiskId.Value"),
			},
			{
				Name:        "root_device_name",
				Description: "The device name of the root device volume (for example, /dev/sda1).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "root_device_type",
				Description: "The root device type used by the AMI. The AMI can use an EBS volume or an instance store volume.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_dest_check",
				Description: "Specifies whether to enable an instance launched in a VPC to perform NAT. This controls whether source/destination checking is enabled on the instance.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "spot_instance_request_id",
				Description: "If the request is a Spot Instance request, the ID of the request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sriov_net_support",
				Description: "Indicates whether enhanced networking with the Intel 82599 Virtual Function interface is enabled.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInstanceSriovNetSupport,
				Transform:   transform.FromField("SriovNetSupport.Value"),
			},
			{
				Name:        "state_code",
				Description: "The reason code for the state change.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("State.Code"),
			},
			{
				Name:        "state_transition_reason",
				Description: "The reason for the most recent state transition.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_transition_time",
				Description: "The date and time, the instance state was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.From(ec2InstanceStateChangeTime),
			},
			{
				Name:        "subnet_id",
				Description: "The ID of the subnet in which the instance is running.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tpm_support",
				Description: "If the instance is configured for NitroTPM support, the value is v2.0.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "usage_operation",
				Description: "The usage operation value for the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "usage_operation_update_time",
				Description: "The time that the usage operation was last updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_data",
				Description: "The user data of the instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInstanceUserData,
				Transform:   transform.FromField("UserData.Value").Transform(base64DecodedData),
			},
			{
				Name:        "virtualization_type",
				Description: "The virtualization type of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC in which the instance is running.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "block_device_mappings",
				Description: "Block device mapping entries for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "elastic_gpu_associations",
				Description: "The Elastic GPU associated with the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "elastic_inference_accelerator_associations",
				Description: "The elastic inference accelerator associated with the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "enclave_options",
				Description: "Indicates whether the instance is enabled for Amazon Web Services Nitro Enclaves.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "hibernation_options",
				Description: "Indicates whether the instance is enabled for hibernation.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "launch_template_data",
				Description: "The configuration data of the specified instance.",
				Hydrate:     getEc2LaunchTemplateData,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "licenses",
				Description: "The license configurations for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "maintenance_options",
				Description: "The metadata options for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metadata_options",
				Description: "The metadata options for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_interfaces",
				Description: "The network interfaces for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_dns_name_options",
				Description: "The options for the instance hostname.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "product_codes",
				Description: "The product codes attached to this instance, if applicable.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_groups",
				Description: "The security groups for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "instance_status",
				Description: "The status of an instance. Instance status includes scheduled events, status checks and instance state information.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getInstanceStatus,
				Transform:   transform.FromField("InstanceStatuses[0]"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getEc2InstanceTurbotTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getEc2InstanceTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEc2InstanceARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2Instance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.listEc2Instance", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			// select * from aws_ec2_instance limit 1
			// Error: InvalidParameterValue: Value ( 1 ) for parameter maxResults is invalid. Expecting a value greater than 5.
			// 		status code: 400, request id: a84912d9-f5fd-403f-8e37-7f7b3f6faba6
			if limit < 5 {
				maxLimit = 5
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int32(maxLimit),
	}
	filters := buildEc2InstanceFilter(d.EqualsQuals)

	if len(filters) != 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeInstancesPaginator(svc, input, func(o *ec2.DescribeInstancesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
				// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_instance.listEc2Instance", "api_error", err)
			return nil, err
		}

		for _, items := range output.Reservations {
			for _, instance := range items.Instances {

				d.StreamListItem(ctx, instance)
				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2Instance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	instanceID := d.EqualsQuals["instance_id"].GetStringValue()

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getEc2Instance", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	op, err := svc.DescribeInstances(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getEc2Instance", "api_error", err)
		return nil, err
	}

	if op.Reservations != nil && len(op.Reservations) > 0 {
		if op.Reservations[0].Instances != nil && len(op.Reservations[0].Instances) > 0 {
			return op.Reservations[0].Instances[0], nil
		}
	}
	return nil, nil
}

func getEc2InstanceARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(types.Instance)
	region := d.EqualsQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getEc2InstanceARN", "getCommonColumns_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":instance/" + *instance.InstanceId

	return arn, nil
}

func getInstanceDisableAPITerminationData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(types.Instance)

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceDisableAPITerminationData", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  types.InstanceAttributeNameDisableApiTermination,
	}

	instanceData, err := svc.DescribeInstanceAttribute(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceDisableAPITerminationData", "api_error", err)
		return nil, err
	}

	return instanceData, nil
}

func getInstanceInitiatedShutdownBehavior(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(types.Instance)

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceInitiatedShutdownBehavior", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  types.InstanceAttributeNameInstanceInitiatedShutdownBehavior,
	}

	instanceData, err := svc.DescribeInstanceAttribute(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceInitiatedShutdownBehavior", "api_error", err)
		return nil, err
	}

	return instanceData, nil
}

func getInstanceKernelID(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(types.Instance)

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceKernelID", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  types.InstanceAttributeNameKernel,
	}

	instanceData, err := svc.DescribeInstanceAttribute(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceKernelID", "api_error", err)
		return nil, err
	}

	return instanceData, nil
}

func getInstanceRAMDiskID(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(types.Instance)

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceRAMDiskID", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  types.InstanceAttributeNameRamdisk,
	}

	instanceData, err := svc.DescribeInstanceAttribute(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceRAMDiskID", "api_error", err)
		return nil, err
	}

	return instanceData, nil
}

func getInstanceSriovNetSupport(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(types.Instance)

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceSriovNetSupport", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  types.InstanceAttributeNameSriovNetSupport,
	}

	instanceData, err := svc.DescribeInstanceAttribute(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceSriovNetSupport", "api_error", err)
		return nil, err
	}

	return instanceData, nil
}

func getInstanceUserData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(types.Instance)

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceUserData", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  types.InstanceAttributeNameUserData,
	}

	instanceData, err := svc.DescribeInstanceAttribute(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceUserData", "api_error", err)
		return nil, err
	}

	return instanceData, nil
}

func getEc2LaunchTemplateData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of load balancer
	instance := h.Item.(types.Instance)

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getEc2LaunchTemplateData", "connection_error", err)
		return nil, err
	}

	params := &ec2.GetLaunchTemplateDataInput{
		InstanceId: instance.InstanceId,
	}

	op, err := svc.GetLaunchTemplateData(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getEc2LaunchTemplateData", "api_error", err)
		return nil, err
	}

	return op.LaunchTemplateData, err
}

func getInstanceStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(types.Instance)

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceStatus", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeInstanceStatusInput{
		InstanceIds:         []string{*instance.InstanceId},
		IncludeAllInstances: aws.Bool(true),
	}

	instanceData, err := svc.DescribeInstanceStatus(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getInstanceStatus", "api_error", err)
		return nil, err
	}

	return instanceData, nil
}

//// TRANSFORM FUNCTIONS

func getEc2InstanceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(types.Instance)
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

func getEc2InstanceTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.Instance)
	title := data.InstanceId
	if data.Tags != nil {
		for _, i := range data.Tags {
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}
	return title, nil
}

func ec2InstanceStateChangeTime(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.Instance)

	if *data.StateTransitionReason != "" {
		if helpers.StringSliceContains([]string{"shutting-down", "stopped", "stopping", "terminated"}, string(data.State.Name)) {
			// User initiated (2019-09-12 16:38:34 GMT)
			regexExp := regexp.MustCompile(`\((.*?) *\)`)
			stateTransitionTime := regexExp.FindStringSubmatch(*data.StateTransitionReason)
			if len(stateTransitionTime) >= 1 {
				stateTransitionTimeInUTC := strings.Replace(strings.Replace(stateTransitionTime[1], " ", "T", 1), " GMT", "Z", 1)
				return stateTransitionTimeInUTC, nil
			}
		}
	}
	return data.LaunchTime, nil
}

//// UTILITY FUNCTIONS

// Build ec2 instance list call input filter
func buildEc2InstanceFilter(equalQuals plugin.KeyColumnEqualsQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"hypervisor":                  "hypervisor",
		"iam_instance_profile_arn":    "iam-instance-profile.arn",
		"image_id":                    "image-id",
		"instance_lifecycle":          "instance-lifecycle",
		"instance_state":              "instance-state-name",
		"instance_type":               "instance-type",
		"monitoring_state":            "monitoring-state",
		"outpost_arn":                 "outpost-arn",
		"placement_availability_zone": "availability-zone",
		"placement_group_name":        "placement-group-name",
		"public_dns_name":             "dns-name",
		"ram_disk_id":                 "ramdisk-id",
		"root_device_name":            "root-device-name",
		"root_device_type":            "root-device-type",
		"subnet_id":                   "subnet-id",
		"placement_tenancy":           "tenancy",
		"virtualization_type":         "virtualization-type",
		"vpc_id":                      "vpc-id",
	}

	for columnName, filterName := range filterQuals {
		if equalQuals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			value := equalQuals[columnName]
			if value.GetStringValue() != "" {
				filter.Values = []string{equalQuals[columnName].GetStringValue()}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}

func getListValues(listValue *proto.QualValueList) []*string {
	values := make([]*string, 0)
	if listValue != nil {
		for _, value := range listValue.Values {
			if value.GetStringValue() != "" {
				values = append(values, aws.String(value.GetStringValue()))
			}
		}
	}
	return values
}
