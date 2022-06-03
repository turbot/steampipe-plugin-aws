package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsEc2LaunchConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_launch_configuration",
		Description: "AWS EC2 Launch Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException", "ValidationError"}),
			},
			Hydrate: getAwsEc2LaunchConfiguration,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsEc2LaunchConfigurations,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the launch configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchConfigurationName"),
			},
			{
				Name:        "launch_configuration_arn",
				Description: "The Amazon Resource Name (ARN) of the launch configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchConfigurationARN"),
			},
			{
				Name:        "created_time",
				Description: "The creation date and time for the launch configuration.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "image_id",
				Description: "The ID of the Amazon Machine Image (AMI) to use to launch EC2 instances.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_type",
				Description: "The instance type for the instances.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associate_public_ip_address",
				Description: "For Auto Scaling groups that are running in a VPC, specifies whether to assign a public IP address to the group's instances.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "kernel_id",
				Description: "The ID of the kernel associated with the AMI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_name",
				Description: "The name of the key pair to be associated with instances.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ramdisk_id",
				Description: "The ID of the RAM disk associated with the AMI.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "ebs_optimized",
				Description: "Specifies whether the launch configuration is optimized for EBS I/O (true) or not (false).",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "classic_link_vpc_id",
				Description: "The ID of a ClassicLink-enabled VPC to link EC2-Classic instances to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClassicLinkVPCId"),
			},
			{
				Name:        "spot_price",
				Description: "The maximum hourly price to be paid for any Spot Instance launched to fulfill the request. Spot Instances are launched when the price you specified exceeds the current Spot price.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_data",
				Description: "The Base64-encoded user data to make available to the launched EC2 instances.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserData").Transform(base64DecodedData),
			},
			{
				Name:        "placement_tenancy",
				Description: "The tenancy of the instance, either default or dedicated. An instance with dedicated tenancy runs on isolated, single-tenant hardware and can only be launched into a VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_instance_profile",
				Description: "The name or the Amazon Resource Name (ARN) of the instance profile associated with the IAM role for the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_monitoring_enabled",
				Description: "Describes whether detailed monitoring is enabled for the Auto Scaling instances.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("InstanceMonitoring.Enabled"),
				Default:     false,
			},
			{
				Name:        "metadata_options_http_endpoint",
				Description: "This parameter enables or disables the HTTP metadata endpoint on instances. If the parameter is not specified, the default state is enabled.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MetadataOptions.HttpEndpoint"),
			},
			{
				Name:        "metadata_options_put_response_hop_limit",
				Description: "The desired HTTP PUT response hop limit for instance metadata requests. The larger the number, the further instance metadata requests can travel.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MetadataOptions.HttpPutResponseHopLimit"),
			},
			{
				Name:        "metadata_options_http_tokens",
				Description: "The state of token usage for your instance metadata requests. If the parameter is not specified in the request, the default state is optional.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MetadataOptions.HttpTokens"),
			},
			{
				Name:        "block_device_mappings",
				Description: "A block device mapping, which specifies the block devices for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "classic_link_vpc_security_groups",
				Description: "The IDs of one or more security groups for the VPC specified in ClassicLinkVPCId.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClassicLinkVPCSecurityGroups"),
			},
			{
				Name:        "security_groups",
				Description: "A list that contains the security groups to assign to the instances in the Auto Scaling group.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchConfigurationName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LaunchConfigurationARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsEc2LaunchConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := AutoScalingService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &autoscaling.DescribeLaunchConfigurationsInput{
		MaxRecords: aws.Int64(100),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxRecords {
			if *limit < 1 {
				input.MaxRecords = aws.Int64(1)
			} else {
				input.MaxRecords = limit
			}
		}
	}

	// List call
	err = svc.DescribeLaunchConfigurationsPages(
		input,
		func(page *autoscaling.DescribeLaunchConfigurationsOutput, isLast bool) bool {
			for _, launchConfiguration := range page.LaunchConfigurations {
				d.StreamListItem(ctx, launchConfiguration)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsEc2LaunchConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2LaunchConfiguration")

	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Session
	svc, err := AutoScalingService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build params
	params := &autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: []*string{aws.String(name)},
	}

	// panic(params)

	rowData, err := svc.DescribeLaunchConfigurations(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAwsEc2LaunchConfiguration", "ERROR", err)
		return nil, err
	}

	if len(rowData.LaunchConfigurations) > 0 && rowData.LaunchConfigurations[0] != nil {
		return rowData.LaunchConfigurations[0], nil
	}

	return nil, nil
}
