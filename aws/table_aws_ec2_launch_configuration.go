package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsEc2LaunchConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_launch_configuration",
		Description: "AWS EC2 Launch Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationError"}),
			},
			Hydrate: getAwsEc2LaunchConfiguration,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsEc2LaunchConfigurations,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
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
	svc, err := AutoScalingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_launch_configuration.listAwsEc2LaunchConfigurations", "connection_error")
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &autoscaling.DescribeLaunchConfigurationsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	paginator := autoscaling.NewDescribeLaunchConfigurationsPaginator(svc, input, func(o *autoscaling.DescribeLaunchConfigurationsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_launch_configuration.listAwsEc2LaunchConfigurations", "api_error", err)
			return nil, err
		}

		for _, items := range output.LaunchConfigurations {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsEc2LaunchConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Session
	svc, err := AutoScalingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_launch_configuration.getAwsEc2LaunchConfiguration", "connection_error")
		return nil, err
	}

	// Build params
	params := &autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: []string{name},
	}

	// panic(params)

	rowData, err := svc.DescribeLaunchConfigurations(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_launch_configuration.getAwsEc2LaunchConfiguration", "api_error", err)
		return nil, err
	}

	if len(rowData.LaunchConfigurations) > 0 {
		return rowData.LaunchConfigurations[0], nil
	}

	return nil, nil
}
