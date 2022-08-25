package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsEc2ASG(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_autoscaling_group",
		Description: "AWS EC2 Autoscaling Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ValidationError"}),
			},
			Hydrate: getAwsEc2AutoscalingGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsEc2AutoscalingGroup,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Auto Scaling group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AutoScalingGroupName"),
			},
			{
				Name:        "autoscaling_group_arn",
				Description: "The Amazon Resource Name (ARN) of the Auto Scaling group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AutoScalingGroupARN"),
			},
			{
				Name:        "status",
				Description: "The current state of the group when the DeleteAutoScalingGroup operation is in progress.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The date and time group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "new_instances_protected_from_scale_in",
				Description: "Indicates whether newly launched instances are protected from termination by Amazon EC2 Auto Scaling when scaling in.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "launch_configuration_name",
				Description: "The name of the associated launch configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_cooldown",
				Description: "The duration of the default cooldown period, in seconds.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "desired_capacity",
				Description: "The desired size of the group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "max_instance_lifetime",
				Description: "The maximum amount of time, in seconds, that an instance can be in service.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "max_size",
				Description: "The maximum size of the group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "min_size",
				Description: "The minimum size of the group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "health_check_grace_period",
				Description: "The amount of time, in seconds, that Amazon EC2 Auto Scaling waits before checking the health status of an EC2 instance that has come into service.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "health_check_type",
				Description: "The service to use for the health checks. The valid values are EC2 and ELB. If you configure an Auto Scaling group to use ELB health checks, it considers the instance unhealthy if it fails either the EC2 status checks or the load balancer health checks.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "placement_group",
				Description: "The name of the placement group into which to launch your instances, if any.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_linked_role_arn",
				Description: "The Amazon Resource Name (ARN) of the service-linked role that the Auto Scaling group uses to call other AWS services on your behalf.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceLinkedRoleARN"),
			},
			{
				Name:        "vpc_zone_identifier",
				Description: "One or more subnet IDs, if applicable, separated by commas.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VPCZoneIdentifier"),
			},
			{
				Name:        "launch_template_name",
				Description: "The launch template name for the group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplate.LaunchTemplateName"),
			},
			{
				Name:        "launch_template_id",
				Description: "The ID of the launch template.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplate.LaunchTemplateId"),
			},
			{
				Name:        "launch_template_version",
				Description: "The version number, $Latest, or $Default.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplate.Version"),
			},
			{
				Name:        "on_demand_allocation_strategy",
				Description: "Indicates how to allocate instance types to fulfill On-Demand capacity. The only valid value is prioritized, which is also the default value. This strategy uses the order of instance types in the overrides to define the launch priority of each instance type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MixedInstancesPolicy.InstancesDistribution.OnDemandAllocationStrategy"),
			},
			{
				Name:        "on_demand_base_capacity",
				Description: "The minimum amount of the Auto Scaling group's capacity that must be fulfilled by On-Demand Instances. This base portion is provisioned first as group scales. Defaults to 0 if not specified.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MixedInstancesPolicy.InstancesDistribution.OnDemandBaseCapacity"),
			},
			{
				Name:        "on_demand_percentage_above_base_capacity",
				Description: "Controls the percentages of On-Demand Instances and Spot Instances for your additional capacity beyond OnDemandBaseCapacity. Expressed as a number (for example, 20 specifies 20% On-Demand Instances, 80% Spot Instances). Defaults to 100 if not specified. If set to 100, only On-Demand Instances are provisioned.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MixedInstancesPolicy.InstancesDistribution.OnDemandPercentageAboveBaseCapacity"),
			},
			{
				Name:        "spot_allocation_strategy",
				Description: "Indicates how to allocate instances across Spot Instance pools. If the allocation strategy is lowest-price, the Auto Scaling group launches instances using the Spot pools with the lowest price, and evenly allocates your instances across the number of Spot pools that you specify. If the allocation strategy is capacity-optimized, the Auto Scaling group launches instances using Spot pools that are optimally chosen based on the available Spot capacity. Defaults to lowest-price if not specified.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MixedInstancesPolicy.InstancesDistribution.SpotAllocationStrategy"),
			},
			{
				Name:        "spot_instance_pools",
				Description: "The number of Spot Instance pools across which to allocate your Spot Instances.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MixedInstancesPolicy.InstancesDistribution.SpotInstancePools"),
			},
			{
				Name:        "spot_max_price",
				Description: "The maximum price per unit hour that user is willing to pay for a Spot Instance. If the value of this parameter is blank (which is the default), the maximum Spot price is set at the On-Demand price.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MixedInstancesPolicy.InstancesDistribution.SpotMaxPrice"),
			},
			{
				Name:        "mixed_instances_policy_launch_template_name",
				Description: "The ID of the launch template for mixed instances policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification.LaunchTemplateName"),
			},
			{
				Name:        "mixed_instances_policy_launch_template_id",
				Description: "The name of the launch template for mixed instances policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification.LaunchTemplateId"),
			},
			{
				Name:        "mixed_instances_policy_launch_template_version",
				Description: "The version of the launch template for mixed instances policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification.Version"),
			},
			{
				Name:        "mixed_instances_policy_launch_template_overrides",
				Description: "Any parameters that is specified in the list override the same parameters in the launch template.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MixedInstancesPolicy.LaunchTemplate.Overrides"),
			},
			{
				Name:        "availability_zones",
				Description: "One or more Availability Zones for the group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "load_balancer_names",
				Description: "One or more load balancers associated with the group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromCamel().Transform(handleLoadBalancersEmptyData),
			},
			{
				Name:        "target_group_arns",
				Description: "The Amazon Resource Names (ARN) of the target groups for your load balancer.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TargetGroupARNs").Transform(handleTargetGroupArnsEmptyData),
			},
			{
				Name:        "instances",
				Description: "The EC2 instances associated with the group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(handleInstancesEmptyData),
			},
			{
				Name:        "enabled_metrics",
				Description: "The metrics enabled for the group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(handleEnableMetricEmptyData),
			},
			{
				Name:        "policies",
				Description: "A set of scaling policies for the specified Auto Scaling group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2AutoscalingGroupPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "termination_policies",
				Description: "The termination policies for the group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "suspended_processes",
				Description: "The suspended processes associated with the group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(handleSuspendedProcessesEmptyData),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Auto Scaling Group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(handleASGTagsEmptyCheck),
			},

			// Steampipe Standard Columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getASGTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AutoScalingGroupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AutoScalingGroupARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsEc2AutoscalingGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := AutoScalingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_autoscaling_group.listAwsEc2AutoscalingGroup", "connection_error", err)
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

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(svc, input, func(o *autoscaling.DescribeAutoScalingGroupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_autoscaling_group.listAwsEc2AutoscalingGroup", "api_error", err)
			return nil, err
		}

		for _, items := range output.AutoScalingGroups {
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

func getAwsEc2AutoscalingGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Session
	svc, err := AutoScalingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_autoscaling_group.getAwsEc2AutoscalingGroup", "connection_error", err)
		return nil, err
	}

	// Build params
	params := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{name},
	}

	rowData, err := svc.DescribeAutoScalingGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_autoscaling_group.getAwsEc2AutoscalingGroup", "api_error", err)
		return nil, err
	}

	if len(rowData.AutoScalingGroups) > 0 {
		return rowData.AutoScalingGroups[0], nil
	}

	return nil, nil
}

func getAwsEc2AutoscalingGroupPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	asg := h.Item.(types.AutoScalingGroup)

	// Create Session
	svc, err := AutoScalingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_autoscaling_group.getAwsEc2AutoscalingGroupPolicy", "connection_error", err)
		return nil, err
	}

	var policies []types.ScalingPolicy

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

	input := &autoscaling.DescribePoliciesInput{
		AutoScalingGroupName: asg.AutoScalingGroupName,
		MaxRecords:           aws.Int32(maxLimit),
	}

	paginator := autoscaling.NewDescribePoliciesPaginator(svc, input, func(o *autoscaling.DescribePoliciesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_autoscaling_group.getAwsEc2AutoscalingGroupPolicy", "api_error", err)
			return nil, err
		}
		for _, policy := range output.ScalingPolicies {
			if len(policy.Alarms) < 1 {
				policy.Alarms = nil
			}

			if len(policy.StepAdjustments) < 1 {
				policy.StepAdjustments = nil
			}
			policies = append(policies, policy)
		}

	}

	return policies, nil
}

//// TRANSFORM FUNCTIONS

func getASGTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	asg := d.HydrateItem.(types.AutoScalingGroup)
	if len(asg.Tags) < 1 {
		return nil, nil
	}

	var turbotTagsMap map[string]string
	if asg.Tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range asg.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

func handleASGTagsEmptyCheck(_ context.Context, d *transform.TransformData) (interface{}, error) {
	asg := d.HydrateItem.(types.AutoScalingGroup)
	if len(asg.Tags) < 1 {
		return nil, nil
	}

	return asg.Tags, nil
}

func handleEnableMetricEmptyData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	asg := d.HydrateItem.(types.AutoScalingGroup)

	if len(asg.EnabledMetrics) > 0 {
		return asg.EnabledMetrics, nil
	}

	return nil, nil
}

func handleInstancesEmptyData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	asg := d.HydrateItem.(types.AutoScalingGroup)

	if len(asg.Instances) > 0 {
		return asg.Instances, nil
	}

	return nil, nil
}

func handleLoadBalancersEmptyData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	asg := d.HydrateItem.(types.AutoScalingGroup)

	if len(asg.LoadBalancerNames) > 0 {
		return asg.LoadBalancerNames, nil
	}

	return nil, nil
}

func handleSuspendedProcessesEmptyData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	asg := d.HydrateItem.(types.AutoScalingGroup)

	if len(asg.SuspendedProcesses) > 0 {
		return asg.SuspendedProcesses, nil
	}

	return nil, nil
}

func handleTargetGroupArnsEmptyData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	asg := d.HydrateItem.(types.AutoScalingGroup)

	if len(asg.TargetGroupARNs) > 0 {
		return asg.TargetGroupARNs, nil
	}

	return nil, nil
}
