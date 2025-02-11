package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2TargetGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_target_group",
		Description: "AWS EC2 Target Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("target_group_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"LoadBalancerNotFound", "TargetGroupNotFound", "ValidationError"}),
			},
			Hydrate: getEc2TargetGroup,
			Tags:    map[string]string{"service": "elasticloadbalancing", "action": "DescribeTargetGroups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2TargetGroups,
			Tags:    map[string]string{"service": "elasticloadbalancing", "action": "DescribeTargetGroups"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"TargetGroupNotFound", "ValidationError"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "target_group_name", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsEc2TargetGroupTargetHealthDescription,
				Tags: map[string]string{"service": "elasticloadbalancing", "action": "DescribeTargetHealth"},
			},
			{
				Func: getAwsEc2TargetGroupTags,
				Tags: map[string]string{"service": "elasticloadbalancing", "action": "DescribeTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ELASTICLOADBALANCING_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "target_group_name",
				Description: "The name of the target group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_group_arn",
				Description: "The Amazon Resource Name (ARN) of the target group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_type",
				Description: "The type of target that is specified when registering targets with this target group. The possible values are instance (register targets by instance ID), ip (register targets by IP address), or lambda (register a single Lambda function as a target).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "load_balancer_arns",
				Description: "The Amazon Resource Names (ARN) of the load balancers that route traffic to this target group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "port",
				Description: "The port on which the targets are listening. Not used if the target is a Lambda function.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC for the target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "protocol",
				Description: "The protocol to use for routing traffic to the target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "protocol_version",
				Description: "The protocol version. The possible values are GRPC , HTTP1 , and HTTP2 .",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address_type",
				Description: "The type of IP address used for this target group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "matcher_http_code",
				Description: "The HTTP codes to use when checking for a successful response from a target.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Matcher.HttpCode"),
			},
			{
				Name:        "matcher_grpc_code",
				Description: "The gRPC codes to use when checking for a successful response from a target.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Matcher.GrpcCode"),
			},
			{
				Name:        "healthy_threshold_count",
				Description: "The number of consecutive health checks successes required before considering an unhealthy target healthy.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "unhealthy_threshold_count",
				Description: "The number of consecutive health checks successes required before considering an unhealthy target healthy.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "health_check_enabled",
				Description: "Indicates whether health checks are enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "health_check_interval_seconds",
				Description: "The approximate amount of time, in seconds, between health checks of an individual target.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "health_check_path",
				Description: "The destination for health checks on the target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_check_port",
				Description: "The port to use to connect with the target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_check_protocol",
				Description: "The protocol to use to connect with the target. The GENEVE, TLS, UDP, and TCP_UDP protocols are not supported for health checks.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_check_timeout_seconds",
				Description: "The amount of time, in seconds, during which no response means a failed health check.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "target_health_descriptions",
				Description: "Contains information about the health of the target.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2TargetGroupTargetHealthDescription,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with target group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2TargetGroupTags,
				Transform:   transform.From(targetGroupRawTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TargetGroupName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2TargetGroupTags,
				Transform:   transform.From(targetGroupTagsToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TargetGroupArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2TargetGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_target_group.listEc2TargetGroups", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(400)
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

	input := &elasticloadbalancingv2.DescribeTargetGroupsInput{
		PageSize: aws.Int32(maxLimit),
	}

	// Additional Filter
	equalQuals := d.EqualsQuals
	if equalQuals["target_group_name"] != nil {
		input.Names = []string{equalQuals["target_group_name"].GetStringValue()}
	}

	paginator := elasticloadbalancingv2.NewDescribeTargetGroupsPaginator(svc, input, func(o *elasticloadbalancingv2.DescribeTargetGroupsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_target_group.listEc2TargetGroups", "api_error", err)
			return nil, err
		}

		for _, items := range output.TargetGroups {
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

func getEc2TargetGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	targetGroupArn := d.EqualsQuals["target_group_arn"].GetStringValue()

	// create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_target_group.getEc2TargetGroup", "connection_error", err)
		return nil, err
	}

	params := &elasticloadbalancingv2.DescribeTargetGroupsInput{
		TargetGroupArns: []string{targetGroupArn},
	}

	op, err := svc.DescribeTargetGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_target_group.getEc2TargetGroup", "api_error", err)
		return nil, err
	}

	if len(op.TargetGroups) > 0 {
		return op.TargetGroups[0], nil
	}
	return nil, nil
}

func getAwsEc2TargetGroupTargetHealthDescription(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	targetGroup := h.Item.(types.TargetGroup)

	// create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_target_group.getAwsEc2TargetGroupTargetHealthDescription", "connection_error", err)
		return nil, err
	}

	params := &elasticloadbalancingv2.DescribeTargetHealthInput{
		TargetGroupArn: targetGroup.TargetGroupArn,
	}

	op, err := svc.DescribeTargetHealth(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_target_group.getAwsEc2TargetGroupTargetHealthDescription", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsEc2TargetGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	targetGroup := h.Item.(types.TargetGroup)

	// create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_target_group.getAwsEc2TargetGroupTags", "connection_error", err)
		return nil, err
	}

	params := &elasticloadbalancingv2.DescribeTagsInput{
		ResourceArns: []string{*targetGroup.TargetGroupArn},
	}

	op, err := svc.DescribeTags(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_target_group.getAwsEc2TargetGroupTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func targetGroupTagsToTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*elasticloadbalancingv2.DescribeTagsOutput)
	var turbotTagsMap map[string]string
	if len(data.TagDescriptions) > 0 {
		if data.TagDescriptions[0].Tags != nil {
			turbotTagsMap = map[string]string{}
			for _, i := range data.TagDescriptions[0].Tags {
				turbotTagsMap[*i.Key] = *i.Value
			}
		}
	}

	return turbotTagsMap, nil
}

func targetGroupRawTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*elasticloadbalancingv2.DescribeTagsOutput)
	if len(data.TagDescriptions) > 0 {
		if data.TagDescriptions[0].Tags != nil {
			return data.TagDescriptions[0].Tags, nil
		}
	}
	return nil, nil
}
