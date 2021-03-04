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

func tableAwsEc2ApplicationLoadBalancer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_application_load_balancer",
		Description: "AWS EC2 Application Load Balancer",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("arn"),
			ShouldIgnoreError: isNotFoundError([]string{"LoadBalancerNotFound"}),
			Hydrate:           getEc2ApplicationLoadBalancer,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2ApplicationLoadBalancers,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name of the Load Balancer that was provided during resource creation",
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
				Name:        "vpc_id",
				Description: "The ID of the VPC for the load balancer",
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
				Hydrate:     getAwsEc2ApplicationLoadBalancerAttributes,
				Transform:   transform.FromField("Attributes"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the load balancer",
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
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listEc2ApplicationLoadBalancers", "AWS_REGION", region)

	// Create Session
	svc, err := ELBv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeLoadBalancersPages(
		&elbv2.DescribeLoadBalancersInput{},
		func(page *elbv2.DescribeLoadBalancersOutput, isLast bool) bool {
			for _, applicationLoadBalancer := range page.LoadBalancers {
				// Filtering the response to return only application load balancers
				if strings.ToLower(*applicationLoadBalancer.Type) == "application" {
					d.StreamListItem(ctx, applicationLoadBalancer)
				}
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2ApplicationLoadBalancer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	loadBalancerArn := d.KeyColumnQuals["arn"].GetStringValue()

	// Create service
	svc, err := ELBv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &elbv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []*string{aws.String(loadBalancerArn)},
	}

	op, err := svc.DescribeLoadBalancers(params)
	if err != nil {
		return nil, err
	}

	if op.LoadBalancers != nil && len(op.LoadBalancers) > 0 {
		return op.LoadBalancers[0], nil
	}
	return nil, nil
}

func getAwsEc2ApplicationLoadBalancerAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2ApplicationLoadBalancerAttributes")
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	applicationLoadBalancer := h.Item.(*elbv2.LoadBalancer)

	// Create service
	svc, err := ELBv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &elbv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: applicationLoadBalancer.LoadBalancerArn,
	}

	loadBalancerData, err := svc.DescribeLoadBalancerAttributes(params)
	if err != nil {
		return nil, err
	}

	return loadBalancerData, nil
}

func getAwsEc2ApplicationLoadBalancerTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2ApplicationLoadBalancerTags")
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	applicationLoadBalancer := h.Item.(*elbv2.LoadBalancer)

	// Create service
	svc, err := ELBv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &elbv2.DescribeTagsInput{
		ResourceArns: []*string{aws.String(*applicationLoadBalancer.LoadBalancerArn)},
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

//// TRANSFORM FUNCTIONS ////

func getEc2ApplicationLoadBalancerTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	applicationLoadBalancerTags := d.HydrateItem.([]*elbv2.Tag)

	if applicationLoadBalancerTags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range applicationLoadBalancerTags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
