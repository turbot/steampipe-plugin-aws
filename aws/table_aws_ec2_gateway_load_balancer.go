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

func tableAwsEc2GatewayLoadBalancer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_gateway_load_balancer",
		Description: "AWS EC2 Gateway Load Balancer",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"LoadBalancerNotFound", "ValidationError"}),
			Hydrate:           getEc2GatewayLoadBalancer,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2GatewayLoadBalancers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "arn", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
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
	svc, err := ELBv2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &elbv2.DescribeLoadBalancersInput{
		PageSize: aws.Int64(400),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["arn"] != nil {
		input.LoadBalancerArns = []*string{aws.String(equalQuals["arn"].GetJsonbValue())}
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.PageSize {
			if *limit < 1 {
				input.PageSize = aws.Int64(1)
			} else {
				input.PageSize = limit
			}
		}
	}

	// List call
	err = svc.DescribeLoadBalancersPages(
		input,
		func(page *elbv2.DescribeLoadBalancersOutput, isLast bool) bool {
			for _, gatewayLoadBalancer := range page.LoadBalancers {
				// Filtering the response to return only gateway load balancers
				if strings.ToLower(*gatewayLoadBalancer.Type) == "gateway" {
					d.StreamListItem(ctx, gatewayLoadBalancer)

					// Context can be cancelled due to manual cancellation or the limit has been hit
					if d.QueryStatus.RowsRemaining(ctx) == 0 {
						return false
					}
				}
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2GatewayLoadBalancer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := ELBv2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	loadBalancerName := quals["name"].GetStringValue()

	params := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{aws.String(loadBalancerName)},
	}

	op, err := svc.DescribeLoadBalancers(params)
	if err != nil {
		return nil, err
	}

	if op.LoadBalancers != nil && len(op.LoadBalancers) > 0 && strings.ToLower(*op.LoadBalancers[0].Type) == "gateway" {
		return op.LoadBalancers[0], nil
	}
	return nil, nil
}

func getAwsEc2GatewayLoadBalancerTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2GatewayLoadBalancerTags")

	gatewayLoadBalancer := h.Item.(*elbv2.LoadBalancer)

	// Create service
	svc, err := ELBv2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &elbv2.DescribeTagsInput{
		ResourceArns: []*string{aws.String(*gatewayLoadBalancer.LoadBalancerArn)},
	}

	loadBalancerData, err := svc.DescribeTags(params)
	if err != nil {
		return nil, err
	}

	if loadBalancerData.TagDescriptions != nil && len(loadBalancerData.TagDescriptions) > 0 {
		return loadBalancerData.TagDescriptions[0].Tags, nil
	}

	return nil, nil
}

func getAwsEc2GatewayLoadBalancerAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2GatewayLoadBalancerAttributes")

	gatewayLoadBalancer := h.Item.(*elbv2.LoadBalancer)

	// Create service
	svc, err := ELBv2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &elbv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: gatewayLoadBalancer.LoadBalancerArn,
	}

	loadBalancerData, err := svc.DescribeLoadBalancerAttributes(params)
	if err != nil {
		return nil, err
	}

	return loadBalancerData, nil
}

//// TRANSFORM FUNCTIONS

func getEc2GatewayLoadBalancerTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	gatewayLoadBalancerTags := d.HydrateItem.([]*elbv2.Tag)

	if gatewayLoadBalancerTags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range gatewayLoadBalancerTags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
