package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	elbv2v1 "github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2ApplicationLoadBalancer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_application_load_balancer",
		Description: "AWS EC2 Application Load Balancer",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"LoadBalancerNotFound", "ValidationError"}),
			},
			Hydrate: getEc2ApplicationLoadBalancer,
			Tags:    map[string]string{"service": "elasticloadbalancing", "action": "DescribeLoadBalancers"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2ApplicationLoadBalancers,
			Tags:    map[string]string{"service": "elasticloadbalancing", "action": "DescribeLoadBalancers"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"LoadBalancerNotFound"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "arn",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsEc2ApplicationLoadBalancerAttributes,
				Tags: map[string]string{"service": "elasticloadbalancing", "action": "DescribeLoadBalancerAttributes"},
			},
			{
				Func: getAwsEc2ApplicationLoadBalancerTags,
				Tags: map[string]string{"service": "elasticloadbalancing", "action": "DescribeTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(elbv2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name of the Load Balancer that was provided during resource creation.",
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
				Name:        "scheme",
				Description: "The load balancing scheme of load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "canonical_hosted_zone_id",
				Description: "The ID of the Amazon Route 53 hosted zone associated with the load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC for the load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The date and time the load balancer was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "customer_owned_ipv4_pool",
				Description: "The ID of the customer-owned address pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dns_name",
				Description: "The public DNS name of the load balancer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DNSName"),
			},
			{
				Name:        "ip_address_type",
				Description: "The type of IP addresses used by the subnets for your load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_code",
				Description: "Current state of the load balancer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State.Code"),
			},
			{
				Name:        "state_reason",
				Description: "A description of the state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State.Reason"),
			},
			{
				Name:        "availability_zones",
				Description: "The subnets for the load balancer.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_groups",
				Description: "The IDs of the security groups for the load balancer.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "load_balancer_attributes",
				Description: "The AWS account ID of the image owner.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2ApplicationLoadBalancerAttributes,
				Transform:   transform.FromField("Attributes"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the load balancer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2ApplicationLoadBalancerTags,
				Transform:   transform.FromValue(),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2ApplicationLoadBalancerTags,
				Transform:   transform.From(getEc2ApplicationLoadBalancerTurbotTags),
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

func listEc2ApplicationLoadBalancers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_application_load_balancer.listEc2ApplicationLoadBalancers", "connection error", err)
		return nil, err
	}

	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}
	maxLimit := int32(400)

	// Additional Filter
	equalQuals := d.EqualsQuals
	if equalQuals["name"] != nil {
		input.Names = []string{equalQuals["name"].GetStringValue()}
	} else {
		// If the names will be provided in param then page limit cannot be set, API throws an error
		// ValidationError: Pagination is not supported when specifying load balancers
		// Limiting the results
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

	if equalQuals["arn"] != nil {
		input.LoadBalancerArns = []string{equalQuals["arn"].GetStringValue()}
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
			plugin.Logger(ctx).Error("aws_ec2_application_load_balancer.listEc2ApplicationLoadBalancers", "api_error", err)
			return nil, err
		}

		for _, items := range output.LoadBalancers {
			if items.Type == types.LoadBalancerTypeEnumApplication {
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

func getEc2ApplicationLoadBalancer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	loadBalancerArn := d.EqualsQuals["arn"].GetStringValue()

	// check if arn is empty
	if loadBalancerArn == "" {
		return nil, nil
	}

	// Create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_application_load_balancer.getEc2ApplicationLoadBalancer", "connection_error", err)
		return nil, err
	}

	params := &elasticloadbalancingv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []string{loadBalancerArn},
	}

	op, err := svc.DescribeLoadBalancers(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_application_load_balancer.getEc2ApplicationLoadBalancer", "api_error", err)
		return nil, err
	}

	if len(op.LoadBalancers) > 0 {
		return op.LoadBalancers[0], nil
	}
	return nil, nil
}

func getAwsEc2ApplicationLoadBalancerAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	applicationLoadBalancer := h.Item.(types.LoadBalancer)

	// Create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_application_load_balancer.getAwsEc2ApplicationLoadBalancerAttributes", "connection_error", err)
		return nil, err
	}

	params := &elasticloadbalancingv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: applicationLoadBalancer.LoadBalancerArn,
	}

	loadBalancerData, err := svc.DescribeLoadBalancerAttributes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_application_load_balancer.getAwsEc2ApplicationLoadBalancerAttributes", "api_error", err)
		return nil, err
	}

	return loadBalancerData, nil
}

func getAwsEc2ApplicationLoadBalancerTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	applicationLoadBalancer := h.Item.(types.LoadBalancer)

	// Create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_application_load_balancer.getAwsEc2ApplicationLoadBalancerTags", "connection_error", err)
		return nil, err
	}

	params := &elasticloadbalancingv2.DescribeTagsInput{
		ResourceArns: []string{*applicationLoadBalancer.LoadBalancerArn},
	}

	loadBalancerData, err := svc.DescribeTags(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_application_load_balancer.getAwsEc2ApplicationLoadBalancerTags", "api_error", err)
		return nil, err
	}

	if len(loadBalancerData.TagDescriptions) > 0 {
		return loadBalancerData.TagDescriptions[0].Tags, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS ////

func getEc2ApplicationLoadBalancerTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	applicationLoadBalancerTags := d.HydrateItem.([]types.Tag)

	if applicationLoadBalancerTags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range applicationLoadBalancerTags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
