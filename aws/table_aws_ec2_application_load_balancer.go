package aws

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2ApplicationLoadBalancer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_application_load_balancer",
		Description: "AWS EC2 Application Load Balancer",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"LoadBalancerNotFound", "ValidationError"}),
			},
			Hydrate: getEc2ApplicationLoadBalancer,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2ApplicationLoadBalancers,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"LoadBalancerNotFound"}),
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
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := ELBv2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &elbv2.DescribeLoadBalancersInput{}

	// Additional Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.Names = []*string{aws.String(equalQuals["name"].GetStringValue())}
	} else {
		// If the names will be provided in param then page limit cannot be set, API throws an error
		// ValidationError: Pagination is not supported when specifying load balancers
		input.PageSize = aws.Int64(400)
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
	}

	if equalQuals["arn"] != nil {
		input.LoadBalancerArns = []*string{aws.String(equalQuals["arn"].GetStringValue())}
	}

	// List call
	err = svc.DescribeLoadBalancersPages(
		input,
		func(page *elbv2.DescribeLoadBalancersOutput, isLast bool) bool {
			for _, applicationLoadBalancer := range page.LoadBalancers {
				// Filtering the response to return only application load balancers
				if strings.ToLower(*applicationLoadBalancer.Type) == "application" {
					d.StreamListItem(ctx, applicationLoadBalancer)

					// Context may get cancelled due to manual cancellation or if the limit has been reached
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

func getEc2ApplicationLoadBalancer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	loadBalancerArn := d.KeyColumnQuals["arn"].GetStringValue()

	// check if arn is empty
	if loadBalancerArn == "" {
		return nil, nil
	}

	// Create service
	svc, err := ELBv2Service(ctx, d)
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

	applicationLoadBalancer := h.Item.(*elbv2.LoadBalancer)

	// Create service
	svc, err := ELBv2Service(ctx, d)
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

	applicationLoadBalancer := h.Item.(*elbv2.LoadBalancer)

	// Create service
	svc, err := ELBv2Service(ctx, d)
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
