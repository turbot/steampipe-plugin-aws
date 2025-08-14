package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2SpotFleetRequest(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_spot_fleet_request",
		Description: "AWS EC2 Spot Fleet Request",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("spot_fleet_request_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidSpotFleetRequestId.NotFound", "InvalidParameterValue"}),
			},
			Hydrate: getEc2SpotFleetRequest,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeSpotFleetRequests"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2SpotFleetRequests,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeSpotFleetRequests"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_EC2_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "spot_fleet_request_id",
				Description: "The ID of the Spot Fleet request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "spot_fleet_request_state",
				Description: "The state of the Spot Fleet request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "activity_status",
				Description: "The progress of the Spot Fleet request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The creation date and time of the request.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "type",
				Description: "The type of request. Indicates whether the Spot Fleet only requests the target capacity or also attempts to maintain it.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.Type"),
			},
			{
				Name:        "target_capacity",
				Description: "The number of units to request for the Spot Fleet.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SpotFleetRequestConfig.TargetCapacity"),
			},
			{
				Name:        "target_capacity_unit_type",
				Description: "The unit for the target capacity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.TargetCapacityUnitType"),
			},
			{
				Name:        "allocation_strategy",
				Description: "The strategy that determines how to allocate the target Spot Instance capacity across the Spot Instance pools.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.AllocationStrategy"),
			},
			{
				Name:        "iam_fleet_role",
				Description: "The Amazon Resource Name (ARN) of an IAM role that grants the Spot Fleet the permission to request, launch, terminate, and tag instances on your behalf.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.IamFleetRole"),
			},
			{
				Name:        "spot_price",
				Description: "The maximum price per unit hour that you are willing to pay for a Spot Instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.SpotPrice"),
			},
			{
				Name:        "spot_max_total_price",
				Description: "The maximum amount per hour for Spot Instances that you're willing to pay.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.SpotMaxTotalPrice"),
			},
			{
				Name:        "on_demand_target_capacity",
				Description: "The number of On-Demand units to request.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SpotFleetRequestConfig.OnDemandTargetCapacity"),
			},
			{
				Name:        "on_demand_allocation_strategy",
				Description: "The order of the launch template overrides to use in fulfilling On-Demand capacity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.OnDemandAllocationStrategy"),
			},
			{
				Name:        "on_demand_max_total_price",
				Description: "The maximum amount per hour for On-Demand Instances that you're willing to pay.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.OnDemandMaxTotalPrice"),
			},
			{
				Name:        "fulfilled_capacity",
				Description: "The number of units fulfilled by this request compared to the set target capacity.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("SpotFleetRequestConfig.FulfilledCapacity"),
			},
			{
				Name:        "on_demand_fulfilled_capacity",
				Description: "The number of On-Demand units fulfilled by this request compared to the set target On-Demand capacity.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("SpotFleetRequestConfig.OnDemandFulfilledCapacity"),
			},
			{
				Name:        "excess_capacity_termination_policy",
				Description: "Indicates whether running instances should be terminated if you decrease the target capacity of the Spot Fleet request below the current size of the Spot Fleet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.ExcessCapacityTerminationPolicy"),
			},
			{
				Name:        "instance_interruption_behavior",
				Description: "The behavior when a Spot Instance is interrupted.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.InstanceInterruptionBehavior"),
			},
			{
				Name:        "instance_pools_to_use_count",
				Description: "The number of Spot pools across which to allocate your target Spot capacity.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SpotFleetRequestConfig.InstancePoolsToUseCount"),
			},
			{
				Name:        "replace_unhealthy_instances",
				Description: "Indicates whether Spot Fleet should replace unhealthy instances.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SpotFleetRequestConfig.ReplaceUnhealthyInstances"),
			},
			{
				Name:        "terminate_instances_with_expiration",
				Description: "Indicates whether running Spot Instances are terminated when the Spot Fleet request expires.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SpotFleetRequestConfig.TerminateInstancesWithExpiration"),
			},
			{
				Name:        "valid_from",
				Description: "The start date and time of the request, in UTC format.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SpotFleetRequestConfig.ValidFrom"),
			},
			{
				Name:        "valid_until",
				Description: "The end date and time of the request, in UTC format.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SpotFleetRequestConfig.ValidUntil"),
			},
			{
				Name:        "client_token",
				Description: "A unique, case-sensitive identifier that you provide to ensure the idempotency of your listings.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.ClientToken"),
			},
			{
				Name:        "context",
				Description: "Reserved.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestConfig.Context"),
			},
			{
				Name:        "launch_specifications",
				Description: "The launch specifications for the Spot Fleet request.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SpotFleetRequestConfig.LaunchSpecifications"),
			},
			{
				Name:        "launch_template_configs",
				Description: "The launch template and overrides.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SpotFleetRequestConfig.LaunchTemplateConfigs"),
			},
			{
				Name:        "load_balancers_config",
				Description: "One or more Classic Load Balancers and target groups to attach to the Spot Fleet request.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SpotFleetRequestConfig.LoadBalancersConfig"),
			},
			{
				Name:        "spot_maintenance_strategies",
				Description: "The strategies for managing your Spot Instances that are at an elevated risk of being interrupted.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SpotFleetRequestConfig.SpotMaintenanceStrategies"),
			},
			{
				Name:        "tag_specifications",
				Description: "The key-value pair for tagging the Spot Fleet request on creation.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SpotFleetRequestConfig.TagSpecifications"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the spot fleet request.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SpotFleetRequestId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(getEc2SpotFleetRequestTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEc2SpotFleetRequestARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2SpotFleetRequests(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_spot_fleet_request.listEc2SpotFleetRequests", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &ec2.DescribeSpotFleetRequestsInput{
		MaxResults: &maxLimit,
	}

	paginator := ec2.NewDescribeSpotFleetRequestsPaginator(svc, input, func(o *ec2.DescribeSpotFleetRequestsPaginatorOptions) {
		o.Limit = maxLimit
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_spot_fleet_request.listEc2SpotFleetRequests", "api_error", err)
			return nil, err
		}

		for _, spotFleetRequest := range output.SpotFleetRequestConfigs {
			d.StreamListItem(ctx, spotFleetRequest)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEc2SpotFleetRequest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	spotFleetRequestId := d.EqualsQualString("spot_fleet_request_id")

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_spot_fleet_request.getEc2SpotFleetRequest", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeSpotFleetRequestsInput{
		SpotFleetRequestIds: []string{spotFleetRequestId},
	}

	op, err := svc.DescribeSpotFleetRequests(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_spot_fleet_request.getEc2SpotFleetRequest", "api_error", err)
		return nil, err
	}

	if op.SpotFleetRequestConfigs != nil && len(op.SpotFleetRequestConfigs) > 0 {
		return op.SpotFleetRequestConfigs[0], nil
	}

	return nil, nil
}

func getEc2SpotFleetRequestARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	spotFleetRequest := h.Item.(types.SpotFleetRequestConfig)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_spot_fleet_request.getEc2SpotFleetRequestARN", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Build ARN
	// arn:${Partition}:ec2:${Region}:${Account}:spot-instances-request/${SpotInstanceRequestId}
	// https://docs.aws.amazon.com/service-authorization/latest/reference/list_amazonec2.html#amazonec2-spot-fleet-request
	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":spot-fleet-request/" + *spotFleetRequest.SpotFleetRequestId

	return arn, nil
}

func getEc2SpotFleetRequestTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem != nil {
		spotFleetRequest := d.HydrateItem.(types.SpotFleetRequestConfig)
		turbotTags := make(map[string]string)
		for _, tag := range spotFleetRequest.Tags {
			turbotTags[*tag.Key] = *tag.Value
		}
		return turbotTags, nil
	}

	return nil, nil
}
