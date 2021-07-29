package aws

import (
	"context"
	"regexp"
	"strings"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2Instance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_instance",
		Description: "AWS EC2 Instance",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("instance_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidInstanceID.NotFound", "InvalidInstanceID.Unavailable", "InvalidInstanceID.Malformed"}),
			Hydrate:           getEc2Instance,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2Instance,
		},
		GetMatrixItem: BuildRegionList,
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
				Name:        "metadata_options",
				Description: "The metadata options for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "outpost_arn",
				Description: "The Amazon Resource Name (ARN) of the Outpost, if applicable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "placement_availability_zone",
				Description: "The Availability Zone of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.AvailabilityZone"),
			},
			{
				Name:        "placement_group_name",
				Description: "The name of the placement group the instance is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.GroupName"),
			},
			{
				Name:        "placement_tenancy",
				Description: "The tenancy of the instance (if the instance is running in a VPC). An instance with a tenancy of dedicated runs on single-tenant hardware.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.Tenancy"),
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
				Name:        "block_device_mappings",
				Description: "Block device mapping entries for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_interfaces",
				Description: "The network interfaces for the instance.",
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listEc2Instance", "AWS_REGION", region)

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeInstancesPages(
		&ec2.DescribeInstancesInput{},
		func(page *ec2.DescribeInstancesOutput, isLast bool) bool {
			if page.Reservations != nil && len(page.Reservations) > 0 {
				for _, reservation := range page.Reservations {
					for _, instance := range reservation.Instances {
						d.StreamListItem(ctx, instance)
					}
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2Instance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2Instance")

	region := d.KeyColumnQualString(matrixKeyRegion)
	instanceID := d.KeyColumnQuals["instance_id"].GetStringValue()

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}

	op, err := svc.DescribeInstances(params)
	if err != nil {
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
	plugin.Logger(ctx).Trace("getEc2InstanceARN")
	instance := h.Item.(*ec2.Instance)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":instance/" + *instance.InstanceId

	return arn, nil
}

func getInstanceDisableAPITerminationData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getInstanceDisableAPITerminationData")
	region := d.KeyColumnQualString(matrixKeyRegion)
	instance := h.Item.(*ec2.Instance)

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  aws.String(ec2.InstanceAttributeNameDisableApiTermination),
	}

	instanceData, err := svc.DescribeInstanceAttribute(params)
	if err != nil {
		return nil, err
	}

	return instanceData, nil
}

func getInstanceInitiatedShutdownBehavior(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getInstanceInitiatedShutdownBehavior")
	region := d.KeyColumnQualString(matrixKeyRegion)
	instance := h.Item.(*ec2.Instance)

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  aws.String(ec2.InstanceAttributeNameInstanceInitiatedShutdownBehavior),
	}

	instanceData, err := svc.DescribeInstanceAttribute(params)
	if err != nil {
		return nil, err
	}

	return instanceData, nil
}

func getInstanceKernelID(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getInstanceKernelID")
	region := d.KeyColumnQualString(matrixKeyRegion)
	instance := h.Item.(*ec2.Instance)

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  aws.String(ec2.InstanceAttributeNameKernel),
	}

	instanceData, err := svc.DescribeInstanceAttribute(params)
	if err != nil {
		return nil, err
	}

	return instanceData, nil
}

func getInstanceRAMDiskID(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getInstanceRAMDiskID")
	region := d.KeyColumnQualString(matrixKeyRegion)
	instance := h.Item.(*ec2.Instance)

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  aws.String(ec2.InstanceAttributeNameRamdisk),
	}

	instanceData, err := svc.DescribeInstanceAttribute(params)
	if err != nil {
		return nil, err
	}

	return instanceData, nil
}

func getInstanceSriovNetSupport(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getInstanceSriovNetSupport")
	region := d.KeyColumnQualString(matrixKeyRegion)
	instance := h.Item.(*ec2.Instance)

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  aws.String(ec2.InstanceAttributeNameSriovNetSupport),
	}

	instanceData, err := svc.DescribeInstanceAttribute(params)
	if err != nil {
		return nil, err
	}

	return instanceData, nil
}

func getInstanceUserData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getInstanceUserData")
	region := d.KeyColumnQualString(matrixKeyRegion)
	instance := h.Item.(*ec2.Instance)

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeInstanceAttributeInput{
		InstanceId: instance.InstanceId,
		Attribute:  aws.String(ec2.InstanceAttributeNameUserData),
	}

	instanceData, err := svc.DescribeInstanceAttribute(params)
	if err != nil {
		return nil, err
	}

	return instanceData, nil
}

func getInstanceStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getInstanceStatus")
	region := d.KeyColumnQualString(matrixKeyRegion)
	instance := h.Item.(*ec2.Instance)

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeInstanceStatusInput{
		InstanceIds:         []*string{instance.InstanceId},
		IncludeAllInstances: types.Bool(true),
	}

	instanceData, err := svc.DescribeInstanceStatus(params)
	if err != nil {
		return nil, err
	}

	return instanceData, nil
}

//// TRANSFORM FUNCTIONS

func getEc2InstanceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(*ec2.Instance)
	return ec2TagsToMap(instance.Tags)
}

func getEc2InstanceTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ec2.Instance)
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
	data := d.HydrateItem.(*ec2.Instance)

	if *data.StateTransitionReason != "" {
		if helpers.StringSliceContains([]string{"shutting-down", "stopped", "stopping", "terminated"}, *data.State.Name) {
			// User initiated (2019-09-12 16:38:34 GMT)
			regexExp := regexp.MustCompile(`\((.*?) *\)`)
			stateTransitionTime := regexExp.FindStringSubmatch(*data.StateTransitionReason)[1]
			stateTransitionTimeInUTC := strings.Replace(strings.Replace(stateTransitionTime, " ", "T", 1), " GMT", "Z", 1)
			return stateTransitionTimeInUTC, nil
		}
	}
	return data.LaunchTime, nil
}
