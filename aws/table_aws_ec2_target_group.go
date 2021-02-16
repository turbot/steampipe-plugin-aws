package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2TargetGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_target_group",
		Description: "AWS EC2 Target Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("target_group_arn"),
			ShouldIgnoreError: isNotFoundError([]string{"LoadBalancerNotFound", "TargetGroupNotFound"}),
			ItemFromKey:       targetGroupFromKey,
			Hydrate:           getEc2TargetGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2TargetGroups,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "target_group_name",
				Description: "The name of the target group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_group_arn",
				Description: "The Amazon Resource Name (ARN) of the target group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_type",
				Description: "The type of target that is specified when registering targets with this target group. The possible values are instance (register targets by instance ID), ip (register targets by IP address), or lambda (register a single Lambda function as a target)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "load_balancer_arns",
				Description: "The Amazon Resource Names (ARN) of the load balancers that route traffic to this target group",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "port",
				Description: "The port on which the targets are listening. Not used if the target is a Lambda function",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC for the targets",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "protocol",
				Description: "The protocol to use for routing traffic to the targets",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "matcher_http_code",
				Description: "The HTTP codes to use when checking for a successful response from a target",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Matcher.HttpCode"),
			},
			{
				Name:        "matcher_grpc_code",
				Description: "The gRPC codes to use when checking for a successful response from a target",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Matcher.GrpcCode"),
			},
			{
				Name:        "healthy_threshold_count",
				Description: "The number of consecutive health checks successes required before considering an unhealthy target healthy",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "unhealthy_threshold_count",
				Description: "The number of consecutive health checks successes required before considering an unhealthy target healthy",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "health_check_enabled",
				Description: "Indicates whether health checks are enabled",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "health_check_interval_seconds",
				Description: "The approximate amount of time, in seconds, between health checks of an individual target",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "health_check_path",
				Description: "The destination for health checks on the targets",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_check_port",
				Description: "The port to use to connect with the target",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_check_protocol",
				Description: "The protocol to use to connect with the target. The GENEVE, TLS, UDP, and TCP_UDP protocols are not supported for health checks",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_check_timeout_seconds",
				Description: "The amount of time, in seconds, during which no response means a failed health check",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "target_health_descriptions",
				Description: "Contains information about the health of the targets",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2TargetGroupTargetHealthDescription,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with target group",
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

//// BUILD HYDRATE INPUT

func targetGroupFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	targetGroupArn := quals["target_group_arn"].GetStringValue()
	item := &elbv2.TargetGroup{
		TargetGroupArn: &targetGroupArn,
	}
	return item, nil
}

//// LIST FUNCTION

func listEc2TargetGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listEc2TargetGroups", "AWS_REGION", region)

	// Create Session
	svc, err := ELBv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeTargetGroupsPages(
		&elbv2.DescribeTargetGroupsInput{},
		func(page *elbv2.DescribeTargetGroupsOutput, isLast bool) bool {
			for _, targetGroup := range page.TargetGroups {
				d.StreamListItem(ctx, targetGroup)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2TargetGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2TargetGroup")
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	targetGroup := h.Item.(*elbv2.TargetGroup)

	// create service
	svc, err := ELBv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &elbv2.DescribeTargetGroupsInput{
		TargetGroupArns: []*string{aws.String(*targetGroup.TargetGroupArn)},
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
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	targetGroup := h.Item.(*elbv2.TargetGroup)

	// create service
	svc, err := ELBv2Service(ctx, d, region)
	if err != nil {
		return nil, err
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
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	targetGroup := h.Item.(*elbv2.TargetGroup)

	// create service
	svc, err := ELBv2Service(ctx, d, region)
	if err != nil {
		return nil, err
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
	if data.TagDescriptions != nil && len(data.TagDescriptions) > 0 {
		if data.TagDescriptions[0].Tags != nil {
			return data.TagDescriptions[0].Tags, nil
		}
	}
	return nil, nil
}
