package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	elbv2v1 "github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2ApplicationLoadBalancerListener(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_load_balancer_listener",
		Description: "AWS EC2 Load Balancer Listener",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ListenerNotFound", "LoadBalancerNotFound", "ValidationError"}),
			},
			Hydrate: getEc2LoadBalancerListener,
			Tags:    map[string]string{"service": "elasticloadbalancing", "action": "DescribeLoadBalancers"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEc2LoadBalancers,
			Hydrate:       listEc2LoadBalancerListeners,
			Tags:          map[string]string{"service": "elasticloadbalancing", "action": "DescribeLoadBalancers"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(elbv2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the listener.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ListenerArn"),
			},
			{
				Name:        "load_balancer_arn",
				Description: "The Amazon Resource Name (ARN) of the load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port",
				Description: "The port on which the load balancer is listening.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "protocol",
				Description: "The protocol for connections from clients to the load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssl_policy",
				Description: "The security policy that defines which protocols and ciphers are supported.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "alpn_policy",
				Description: "The name of the Application-Layer Protocol Negotiation (ALPN) policy.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "certificates",
				Description: "The default certificate for the listener.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "default_actions",
				Description: "The default actions for the listener.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "mutual_authentication",
				Description: "The mutual authentication configuration information.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getEc2ApplicationLoadBalancerListenerTurbotTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ListenerArn").Transform(arnToAkas),
			},
		}),
	}
}

//// PARENT LIST FUNCTION

func listEc2LoadBalancers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_load_balancer_listener.listEc2LoadBalancers", "connection_error", err)
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

	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{
		PageSize: aws.Int32(maxLimit),
	}

	paginator := elasticloadbalancingv2.NewDescribeLoadBalancersPaginator(svc, input, func(o *elasticloadbalancingv2.DescribeLoadBalancersPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_load_balancer_listener.listEc2LoadBalancers", "api_error", err)
			return nil, err
		}

		for _, items := range output.LoadBalancers {
			d.StreamListItem(ctx, items)
		}

	}

	return nil, err
}

//// LIST FUNCTION

func listEc2LoadBalancerListeners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of load balancer
	loadBalancerDetails := h.Item.(types.LoadBalancer)

	// Create Session
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_load_balancer_listener.listEc2LoadBalancerListeners", "connection_error", err)
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

	input := &elasticloadbalancingv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(string(*loadBalancerDetails.LoadBalancerArn)),
		PageSize:        aws.Int32(maxLimit),
	}

	paginator := elasticloadbalancingv2.NewDescribeListenersPaginator(svc, input, func(o *elasticloadbalancingv2.DescribeListenersPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_load_balancer_listener.listEc2LoadBalancerListeners", "api_error", err)
			return nil, err
		}

		for _, items := range output.Listeners {
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

func getEc2LoadBalancerListener(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	listenerArn := d.EqualsQuals["arn"].GetStringValue()

	// Create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_load_balancer_listener.getEc2LoadBalancerListener", "connection_error", err)
		return nil, err
	}

	params := &elasticloadbalancingv2.DescribeListenersInput{
		ListenerArns: []string{listenerArn},
	}

	op, err := svc.DescribeListeners(ctx, params)
	if err != nil {
		return nil, err
	}

	if op.Listeners != nil && len(op.Listeners) > 0 {
		return op.Listeners[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS ////

func getEc2ApplicationLoadBalancerListenerTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.Listener)
	splitID := strings.Split(string(*data.ListenerArn), "/")
	title := splitID[2] + "-" + splitID[4]
	return title, nil
}
