package aws

import (
	"context"

	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	elbv2v1 "github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2GatewayLoadBalancer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_gateway_load_balancer",
		Description: "AWS EC2 Gateway Load Balancer",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"LoadBalancerNotFound", "ValidationError"}),
			},
			Hydrate: getEc2GatewayLoadBalancer,
			Tags:    map[string]string{"service": "elasticloadbalancing", "action": "DescribeClientVpnEndpoints"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2GatewayLoadBalancers,
			Tags:    map[string]string{"service": "elasticloadbalancing", "action": "DescribeClientVpnEndpoints"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationError"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "arn", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsEc2GatewayLoadBalancerAttributes,
				Tags: map[string]string{"service": "elasticloadbalancing", "action": "DescribeLoadBalancerAttributes"},
			},
			{
				Func: getAwsEc2GatewayLoadBalancerTags,
				Tags: map[string]string{"service": "elasticloadbalancing", "action": "DescribeTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(elbv2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the load balancer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancerName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the load balancer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancerArn"),
			},
			{
				Name:        "type",
				Description: "The type of load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enforce_security_group_inbound_rules_on_private_link_traffic",
				Description: "Indicates whether to evaluate inbound security group rules for traffic sent to a Network Load Balancer through Amazon Web Services PrivateLink.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_code",
				Description: "The state of the load balancer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State.Code"),
			},
			{
				Name:        "scheme",
				Description: "The load balancing scheme of gateway load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dns_name",
				Description: "The public DNS name of the gateway load balancer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DNSName"),
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC for the gateway load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The date and time the load balancer was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "ip_address_type",
				Description: "The type of IP addresses used by the subnets for your load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zones",
				Description: "The subnets for the gateway load balancer.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "canonical_hosted_zone_id",
				Description: "The ID of the Amazon Route 53 hosted zone associated with the gateway load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_owned_ipv4_pool",
				Description: "The ID of the customer-owned address pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "security_groups",
				Description: "The IDs of the security groups for the gateway load balancer.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "load_balancer_attributes",
				Description: "Attributes deletion protection and cross_zone of gateway load balancer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2GatewayLoadBalancerAttributes,
				Transform:   transform.FromField("Attributes"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the load balancer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2GatewayLoadBalancerTags,
				Transform:   transform.FromValue(),
			},

			// Standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerArn").Transform(arnToAkas),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancerName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2GatewayLoadBalancerTags,
				Transform:   transform.From(getEc2GatewayLoadBalancerTurbotTags),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2GatewayLoadBalancers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_gateway_load_balancer.listEc2GatewayLoadBalancers", "connection_error", err)
		return nil, err
	}

	input := &elbv2.DescribeLoadBalancersInput{}

	if d.Quals["arn"] != nil {
		arn := getQualsValueByColumn(d.Quals, "arn", "string")
		val, ok := arn.(string)
		if ok {
			input.LoadBalancerArns = []string{val}
		}
	} else {
		// Pagination is not supported when specifying load balancers
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
		input.PageSize = &maxLimit
	}

	paginator := elbv2.NewDescribeLoadBalancersPaginator(svc, input, func(o *elbv2.DescribeLoadBalancersPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_gateway_load_balancer.listEc2GatewayLoadBalancers", "api_error", err)
			return nil, err
		}

		for _, items := range output.LoadBalancers {

			if items.Type == types.LoadBalancerTypeEnumGateway {
				d.StreamListItem(ctx, items)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2GatewayLoadBalancer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_gateway_load_balancer.getEc2GatewayLoadBalancer", "connection_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	loadBalancerName := quals["name"].GetStringValue()

	params := &elbv2.DescribeLoadBalancersInput{
		Names: []string{loadBalancerName},
	}

	op, err := svc.DescribeLoadBalancers(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_gateway_load_balancer.getEc2GatewayLoadBalancer", "api_error", err)
		return nil, err
	}

	if op.LoadBalancers != nil && len(op.LoadBalancers) > 0 && op.LoadBalancers[0].Type == types.LoadBalancerTypeEnumGateway {
		return op.LoadBalancers[0], nil
	}
	return nil, nil
}

func getAwsEc2GatewayLoadBalancerTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	gatewayLoadBalancer := h.Item.(types.LoadBalancer)

	// Create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_gateway_load_balancer.getAwsEc2GatewayLoadBalancerTags", "connection_error", err)
		return nil, err
	}

	params := &elbv2.DescribeTagsInput{
		ResourceArns: []string{*gatewayLoadBalancer.LoadBalancerArn},
	}

	loadBalancerData, err := svc.DescribeTags(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_gateway_load_balancer.getAwsEc2GatewayLoadBalancerTags", "api_error", err)
		return nil, err
	}

	if loadBalancerData.TagDescriptions != nil && len(loadBalancerData.TagDescriptions) > 0 {
		return loadBalancerData.TagDescriptions[0].Tags, nil
	}

	return nil, nil
}

func getAwsEc2GatewayLoadBalancerAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	gatewayLoadBalancer := h.Item.(types.LoadBalancer)

	// Create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_gateway_load_balancer.getAwsEc2GatewayLoadBalancerAttributes", "connection_error", err)
		return nil, err
	}

	params := &elbv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: gatewayLoadBalancer.LoadBalancerArn,
	}

	loadBalancerData, err := svc.DescribeLoadBalancerAttributes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_gateway_load_balancer.getAwsEc2GatewayLoadBalancerAttributes", "api_error", err)
		return nil, err
	}

	return loadBalancerData, nil
}

//// TRANSFORM FUNCTIONS

func getEc2GatewayLoadBalancerTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	gatewayLoadBalancerTags := d.HydrateItem.([]types.Tag)

	if gatewayLoadBalancerTags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range gatewayLoadBalancerTags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
