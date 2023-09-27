package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	elbv2v1 "github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2NetworkLoadBalancer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_network_load_balancer",
		Description: "AWS EC2 Network Load Balancer",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"LoadBalancerNotFound", "ValidationError"}),
			},
			Hydrate: getEc2NetworkLoadBalancer,
			Tags:    map[string]string{"service": "elasticloadbalancing", "action": "DescribeLoadBalancers"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2NetworkLoadBalancers,
			Tags:    map[string]string{"service": "elasticloadbalancing", "action": "DescribeLoadBalancers"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"LoadBalancerNotFound", "ValidationError"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsEc2NetworkLoadBalancerAttributes,
				Tags: map[string]string{"service": "elasticloadbalancing", "action": "DescribeLoadBalancerAttributes"},
			},
			{
				Func: getAwsEc2NetworkLoadBalancerTags,
				Tags: map[string]string{"service": "elasticloadbalancing", "action": "DescribeTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(elbv2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name of the Load Balancer",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancerName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the load balancer",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancerArn"),
			},
			{
				Name:        "type",
				Description: "The type of load balancer",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scheme",
				Description: "The load balancing scheme of load balancer",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "canonical_hosted_zone_id",
				Description: "The ID of the Amazon Route 53 hosted zone associated with the load balancer",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The date and time the load balancer was created",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "customer_owned_ipv4_pool",
				Description: "The ID of the customer-owned address pool",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dns_name",
				Description: "The public DNS name of the load balancer",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DNSName"),
			},
			{
				Name:        "ip_address_type",
				Description: "The type of IP addresses used by the subnets for your load balancer",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_code",
				Description: "Current state of the load balancer",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State.Code"),
			},
			{
				Name:        "state_reason",
				Description: "A description of the state",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State.Reason"),
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC for the load balancer",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zones",
				Description: "The subnets for the load balancer",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_groups",
				Description: "The IDs of the security groups for the load balancer",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "load_balancer_attributes",
				Description: "The AWS account ID of the image owner",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2NetworkLoadBalancerAttributes,
				Transform:   transform.FromField("Attributes"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the load balancer",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2NetworkLoadBalancerTags,
				Transform:   transform.FromValue(),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2NetworkLoadBalancerTags,
				Transform:   transform.From(getEc2NetworkLoadBalancerTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancerName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2NetworkLoadBalancers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		logger.Error("aws_ec2_network_load_balancer.listEc2NetworkLoadBalancers", "connection error", err)
		return nil, err
	}

	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}

	equalQuals := d.EqualsQuals
	if equalQuals["name"] != nil {
		input.Names = []string{equalQuals["name"].GetStringValue()}
	} else {
		// If the names will be provided in param then page limit can not be set, api throws error
		// ValidationError: Pagination is not supported when specifying load balancers
		// Limiting the results
		maxLimit := int32(400)
		if d.QueryContext.Limit != nil {
			limit := int32(*d.QueryContext.Limit)
			if limit < maxLimit {
				if limit < 1 {
					input.PageSize = aws.Int32(1)
				} else {
					input.PageSize = aws.Int32(limit)
				}
			}
		}
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
			plugin.Logger(ctx).Error("aws_api_gateway_rest_api.listRestAPI", "api_error", err)
			return nil, err
		}

		for _, items := range output.LoadBalancers {
			if items.Type == types.LoadBalancerTypeEnumNetwork {

				d.StreamListItem(ctx, items)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}

	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2NetworkLoadBalancer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	loadBalancerArn := d.EqualsQuals["arn"].GetStringValue()

	// Create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_rest_api.getEc2NetworkLoadBalancer", "connection", err)
		return nil, err
	}

	params := &elasticloadbalancingv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []string{loadBalancerArn},
	}

	op, err := svc.DescribeLoadBalancers(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_rest_api.getEc2NetworkLoadBalancer", "api_error", err)
		return nil, err
	}

	if op.LoadBalancers != nil && len(op.LoadBalancers) > 0 {
		return op.LoadBalancers[0], nil
	}
	return nil, nil
}

func getAwsEc2NetworkLoadBalancerAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2NetworkLoadBalancerAttributes")

	networkLoadBalancer := h.Item.(types.LoadBalancer)

	// Create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_rest_api.getAwsEc2NetworkLoadBalancerAttributes", "connection_error", err)
		return nil, err
	}

	params := &elasticloadbalancingv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: networkLoadBalancer.LoadBalancerArn,
	}

	loadBalancerData, err := svc.DescribeLoadBalancerAttributes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_rest_api.getAwsEc2NetworkLoadBalancerAttributes", "api_error", err)
		return nil, err
	}

	return loadBalancerData, nil
}

func getAwsEc2NetworkLoadBalancerTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	networkLoadBalancer := h.Item.(types.LoadBalancer)

	// Create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_rest_api.getAwsEc2NetworkLoadBalancerTags", "connection_error", err)
		return nil, err
	}

	params := &elasticloadbalancingv2.DescribeTagsInput{
		ResourceArns: []string{*networkLoadBalancer.LoadBalancerArn},
	}

	loadBalancerData, err := svc.DescribeTags(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_rest_api.getAwsEc2NetworkLoadBalancerTags", "api_error", err)
		return nil, err
	}

	if loadBalancerData.TagDescriptions != nil && len(loadBalancerData.TagDescriptions) > 0 {
		return loadBalancerData.TagDescriptions[0].Tags, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS ////

func getEc2NetworkLoadBalancerTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	networkLoadBalancerTags := d.HydrateItem.([]types.Tag)

	if networkLoadBalancerTags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range networkLoadBalancerTags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
