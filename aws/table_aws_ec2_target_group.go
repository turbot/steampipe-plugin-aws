package aws

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2TargetGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_target_group",
		Description: "AWS EC2 Target Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("target_group_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"LoadBalancerNotFound", "TargetGroupNotFound", "ValidationError"}),
			},
			Hydrate: getEc2TargetGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2TargetGroups,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"TargetGroupNotFound", "ValidationError"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "target_group_name", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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

func listEc2TargetGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ELBv2Service(ctx, d, h)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &elbv2.DescribeTargetGroupsInput{
		PageSize: aws.Int64(400),
	}

	// Additional Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["target_group_name"] != nil {
		input.Names = []*string{aws.String(equalQuals["target_group_name"].GetStringValue())}
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.PageSize {
			if *limit < 1 {
				input.PageSize = types.Int64(1)
			} else {
				input.PageSize = limit
			}
		}
	}

	// List call
	err = svc.DescribeTargetGroupsPages(
		input,
		func(page *elbv2.DescribeTargetGroupsOutput, isLast bool) bool {
			for _, targetGroup := range page.TargetGroups {
				d.StreamListItem(ctx, targetGroup)

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

func getEc2TargetGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2TargetGroup")

	targetGroupArn := d.KeyColumnQuals["target_group_arn"].GetStringValue()

	// create service
	svc, err := ELBv2Service(ctx, d, h)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &elbv2.DescribeTargetGroupsInput{
		TargetGroupArns: []*string{aws.String(targetGroupArn)},
	}

	op, err := svc.DescribeTargetGroups(params)
	if err != nil {
		return nil, err
	}

	if op.TargetGroups != nil && len(op.TargetGroups) > 0 {
		return op.TargetGroups[0], nil
	}
	return nil, nil
}

func getAwsEc2TargetGroupTargetHealthDescription(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2TargetGroupTargetHealthDescription")

	targetGroup := h.Item.(*elbv2.TargetGroup)

	// create service
	svc, err := ELBv2Service(ctx, d, h)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: targetGroup.TargetGroupArn,
	}

	op, err := svc.DescribeTargetHealth(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getAwsEc2TargetGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2TargetGroupTags")

	targetGroup := h.Item.(*elbv2.TargetGroup)

	// create service
	svc, err := ELBv2Service(ctx, d, h)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &elbv2.DescribeTagsInput{
		ResourceArns: []*string{aws.String(*targetGroup.TargetGroupArn)},
	}

	op, err := svc.DescribeTags(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func targetGroupTagsToTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*elbv2.DescribeTagsOutput)
	var turbotTagsMap map[string]string
	if data.TagDescriptions != nil && len(data.TagDescriptions) > 0 {
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
	data := d.HydrateItem.(*elbv2.DescribeTagsOutput)
	if len(data.TagDescriptions) < 1 {
		return nil, nil
	}

	var tags []*elbv2.Tag
	if data.TagDescriptions != nil && len(data.TagDescriptions) > 0 {
		for _, tag := range data.TagDescriptions {
			tags = append(tags, tag.Tags...)
		}
	}
	return tags, nil
}
