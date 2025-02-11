package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

//// TABLE DEFINITION

func tableAwsEc2ApplicationLoadBalancerListenerRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_load_balancer_listener_rule",
		Description: "AWS EC2 Load Balancer Listener Rule",
		List: &plugin.ListConfig{
			Hydrate: listEc2LoadBalancerListenerRules,
			// We are encountering a ValidationError with the specified resource ARN if the resource is not available in the queried region.
			//Error: aws: operation error Elastic Load Balancing v2: DescribeRules, https response error StatusCode: 400, RequestID: 7de7321c-afe9-4030-bde5-57f0d56c29f4, api error ValidationError: 'arn:aws:elasticloadbalancing:us-east-1:xxxxxxxxxxxx:listener/app/f7cc8cdc44ff910b/f7cc8cdc44ff910b/c9418b57592205f0' is not a valid listener ARN
			// Therefore, we need to add 'ValidationError' to the ignore config.
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ListenerNotFound", "RuleNotFound", "ValidationError"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "listener_arn", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "arn", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "elasticloadbalancing", "action": "DescribeRules"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ELASTICLOADBALANCING_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RuleArn"),
			},
			{
				Name:        "priority",
				Description: "The priority of the rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "listener_arn",
				Description: "The Amazon Resource Name (ARN) of the listener.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("listener_arn"),
			},
			{
				Name:        "is_default",
				Description: "Indicates whether this is the default rule.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "actions",
				Description: "The actions. Each rule must include exactly one of the following types of actions: forward , redirect , or fixed-response , and it must be the last action to be performed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "conditions",
				Description: "The conditions. Each rule can include zero or one of the following conditions: http-request-method , host-header , path-pattern , and source-ip , and zero or more of the following conditions: http-header and query-string.",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2LoadBalancerListenerRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	listenerArn := d.EqualsQualString("listener_arn")
	ruleArn := d.EqualsQualString("arn")

	// You must specify either a listener or one or more rules.
	// https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_DescribeRules.html
	if listenerArn == "" && ruleArn == "" {
		return nil, fmt.Errorf("you must specify either 'listener_arn' or 'rule_arn' in the query parameter to query this table")
	}
	// Create Session
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_load_balancer_listener_rule.listEc2LoadBalancerListenerRules", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(400)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &elasticloadbalancingv2.DescribeRulesInput{
		PageSize: aws.Int32(maxLimit),
	}

	if listenerArn != "" {
		input.ListenerArn = &listenerArn
	}
	if ruleArn != "" {
		input.RuleArns = []string{ruleArn}
	}

	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		items, err := svc.DescribeRules(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_load_balancer_listener_rule.listEc2LoadBalancerListenerRules", "api_error", err)
			return nil, err
		}

		for _, item := range items.Rules {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if items.NextMarker != nil {
			input.Marker = items.NextMarker
		} else {
			break
		}
	}

	return nil, err
}
