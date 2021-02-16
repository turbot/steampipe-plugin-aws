package aws

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2ApplicationLoadBalancerListener(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_load_balancer_listener",
		Description: "AWS EC2 Load Balancer Listener",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("arn"),
			ShouldIgnoreError: isNotFoundError([]string{"ListenerNotFound", "LoadBalancerNotFound"}),
			ItemFromKey:       listenerArnFromKey,
			Hydrate:           getEc2LoadBalancerListener,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEc2LoadBalancers,
			Hydrate:       listEc2LoadBalancerListeners,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the listener",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ListenerArn"),
			},
			{
				Name:        "load_balancer_arn",
				Description: "The Amazon Resource Name (ARN) of the load balancer",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port",
				Description: "The port on which the load balancer is listening",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "protocol",
				Description: "The protocol for connections from clients to the load balancer",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssl_policy",
				Description: "The security policy that defines which protocols and ciphers are supported",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "alpn_policy",
				Description: "The name of the Application-Layer Protocol Negotiation (ALPN) policy",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "certificates",
				Description: "The default certificate for the listener",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "default_actions",
				Description: "The default actions for the listener",
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

//// BUILD HYDRATE INPUT

func listenerArnFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	listenerArn := quals["arn"].GetStringValue()
	item := &elbv2.Listener{
		ListenerArn: &listenerArn,
	}
	return item, nil
}

//// PARENT LIST FUNCTION

func listEc2LoadBalancers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listEc2ApplicationLoadBalancers", "AWS_REGION", defaultRegion)

	// Create Session
	svc, err := ELBv2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeLoadBalancersPages(
		&elbv2.DescribeLoadBalancersInput{},
		func(page *elbv2.DescribeLoadBalancersOutput, isLast bool) bool {
			for _, loadBalancer := range page.LoadBalancers {
				d.StreamListItem(ctx, loadBalancer)
			}
			return !isLast
		},
	)
	return nil, err
}

//// LIST FUNCTION

func listEc2LoadBalancerListeners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listEc2LoadBalancerListeners", "AWS_REGION", defaultRegion)

	// Get the details of load balancer
	loadBalancerDetails := h.Item.(*elbv2.LoadBalancer)

	// Create Session
	svc, err := ELBv2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeListenersPages(
		&elbv2.DescribeListenersInput{
			LoadBalancerArn: aws.String(string(*loadBalancerDetails.LoadBalancerArn)),
		},
		func(page *elbv2.DescribeListenersOutput, isLast bool) bool {
			for _, listener := range page.Listeners {
				d.StreamLeafListItem(ctx, listener)
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2LoadBalancerListener(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	loadBalancerListener := h.Item.(*elbv2.Listener)

	// Create service
	svc, err := ELBv2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &elbv2.DescribeListenersInput{
		ListenerArns: []*string{aws.String(*loadBalancerListener.ListenerArn)},
	}

	op, err := svc.DescribeListeners(params)
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
	data := d.HydrateItem.(*elbv2.Listener)
	splitID := strings.Split(string(*data.ListenerArn), "/")
	title := splitID[2] + "-" + splitID[4]
	return title, nil
}
