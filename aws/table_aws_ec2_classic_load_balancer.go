package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2ClassicLoadBalancer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_classic_load_balancer",
		Description: "AWS EC2 Classic Load Balancer",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"LoadBalancerNotFound"}),
			ItemFromKey:       loadBalancerNameFromKey,
			Hydrate:           getEc2ClassicLoadBalancer,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2ClassicLoadBalancers,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name of the Load Balancer",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancerName"),
			},
			{
				Name:        "scheme",
				Description: "The load balancing scheme of load balancer",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The date and time the load balancer was created",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC for the load balancer",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VPCId"),
			},
			{
				Name:        "access_log_emit_interval",
				Description: "The interval for publishing the access logs",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsEc2ClassicLoadBalancerAttributes,
				Transform:   transform.FromField("LoadBalancerAttributes.AccessLog.EmitInterval"),
			},
			{
				Name:        "access_log_enabled",
				Description: "Specifies whether access logs are enabled for the load balancer",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsEc2ClassicLoadBalancerAttributes,
				Transform:   transform.FromField("LoadBalancerAttributes.AccessLog.Enabled"),
			},
			{
				Name:        "access_log_s3_bucket_name",
				Description: "The name of the Amazon S3 bucket where the access logs are stored",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsEc2ClassicLoadBalancerAttributes,
				Transform:   transform.FromField("LoadBalancerAttributes.AccessLog.S3BucketName"),
			},
			{
				Name:        "access_log_s3_bucket_prefix",
				Description: "The logical hierarchy you created for your Amazon S3 bucket",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsEc2ClassicLoadBalancerAttributes,
				Transform:   transform.FromField("LoadBalancerAttributes.AccessLog.S3BucketPrefix"),
			},
			{
				Name:        "canonical_hosted_zone_name",
				Description: "The name of the Amazon Route 53 hosted zone for the load balancer",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "canonical_hosted_zone_name_id",
				Description: "The ID of the Amazon Route 53 hosted zone for the load balancer",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CanonicalHostedZoneNameID"),
			},
			{
				Name:        "connection_draining_enabled",
				Description: "Specifies whether connection draining is enabled for the load balancer",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsEc2ClassicLoadBalancerAttributes,
				Transform:   transform.FromField("LoadBalancerAttributes.ConnectionDraining.Enabled"),
			},
			{
				Name:        "connection_draining_timeout",
				Description: "The maximum time, in seconds, to keep the existing connections open before deregistering the instances",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsEc2ClassicLoadBalancerAttributes,
				Transform:   transform.FromField("LoadBalancerAttributes.ConnectionDraining.Timeout"),
			},
			{
				Name:        "connection_settings_idle_timeout",
				Description: "The time, in seconds, that the connection is allowed to be idle (no data has been sent over the connection) before it is closed by the load balancer",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsEc2ClassicLoadBalancerAttributes,
				Transform:   transform.FromField("LoadBalancerAttributes.ConnectionSettings.IdleTimeout"),
			},
			{
				Name:        "cross_zone_load_balancing_enabled",
				Description: "Specifies whether cross-zone load balancing is enabled for the load balancer",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsEc2ClassicLoadBalancerAttributes,
				Transform:   transform.FromField("LoadBalancerAttributes.CrossZoneLoadBalancing.Enabled"),
			},
			{
				Name:        "dns_name",
				Description: "The DNS name of the load balancer",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DNSName"),
			},
			{
				Name:        "health_check_interval",
				Description: "The approximate interval, in seconds, between health checks of an individual instance",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HealthCheck.Interval"),
			},
			{
				Name:        "health_check_timeout",
				Description: "The amount of time, in seconds, during which no response means a failed health check",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HealthCheck.Timeout"),
			},
			{
				Name:        "healthy_threshold",
				Description: "The number of consecutive health checks successes required before moving the instance to the Healthy state",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HealthCheck.HealthyThreshold"),
			},
			{
				Name:        "heath_check_target",
				Description: "The instance being checked",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("HealthCheck.Target"),
			},
			{
				Name:        "source_security_group_name",
				Description: "The name of the security group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceSecurityGroup.GroupName"),
			},
			{
				Name:        "source_security_group_owner_alias",
				Description: "The owner of the security group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceSecurityGroup.OwnerAlias"),
			},
			{
				Name:        "unhealthy_threshold",
				Description: "The number of consecutive health check failures required before moving the instance to the Unhealthy state",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HealthCheck.UnhealthyThreshold"),
			},
			{
				Name:        "additional_attributes",
				Description: "A list of additional attributes",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2ClassicLoadBalancerAttributes,
				Transform:   transform.FromField("LoadBalancerAttributes.AdditionalAttributes"),
			},
			{
				Name:        "app_cookie_stickiness_policies",
				Description: "A list of the stickiness policies created using CreateAppCookieStickinessPolicy",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Policies.AppCookieStickinessPolicies"),
			},
			{
				Name:        "availability_zones",
				Description: "A list of the Availability Zones for the load balancer",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "backend_server_descriptions",
				Description: "A list of information about your EC2 instances",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "instances",
				Description: "A list of the IDs of the instances for the load balancer",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lb_cookie_stickiness_policies",
				Description: "A list of the stickiness policies created using CreateLBCookieStickinessPolicy",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Policies.LBCookieStickinessPolicies"),
			},
			{
				Name:        "listener_descriptions",
				Description: "A list of the listeners for the load balancer",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "other_policies",
				Description: "A list of policies other than the stickiness policies",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Policies.OtherPolicies"),
			},
			{
				Name:        "security_groups",
				Description: "A list of the security groups for the load balancer",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subnets",
				Description: "A list of the IDs of the subnets for the load balancer",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the load balancer",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2ClassicLoadBalancerTags,
				Transform:   transform.FromValue(),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2ClassicLoadBalancerTags,
				Transform:   transform.From(getEc2ClassicLoadBalancerTurbotTags),
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
				Hydrate:     getEc2ClassicLoadBalancerTurbotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// BUILD HYDRATE INPUT

func loadBalancerNameFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	loadBalancerName := quals["name"].GetStringValue()
	item := &elb.LoadBalancerDescription{
		LoadBalancerName: &loadBalancerName,
	}
	return item, nil
}

//// LIST FUNCTION

func listEc2ClassicLoadBalancers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listEc2ClassicLoadBalancers", "AWS_REGION", region)

	// Create Session
	svc, err := ELBService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeLoadBalancersPages(
		&elb.DescribeLoadBalancersInput{},
		func(page *elb.DescribeLoadBalancersOutput, isLast bool) bool {
			for _, classicLoadBalancer := range page.LoadBalancerDescriptions {
				d.StreamListItem(ctx, classicLoadBalancer)
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2ClassicLoadBalancer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	classicLoadBalancer := h.Item.(*elb.LoadBalancerDescription)

	// Create service
	svc, err := ELBService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{aws.String(*classicLoadBalancer.LoadBalancerName)},
	}

	op, err := svc.DescribeLoadBalancers(params)
	if err != nil {
		return nil, err
	}

	if op.LoadBalancerDescriptions != nil && len(op.LoadBalancerDescriptions) > 0 {
		return op.LoadBalancerDescriptions[0], nil
	}
	return nil, nil
}

func getAwsEc2ClassicLoadBalancerAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2ClassicLoadBalancerAttributes")
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	classicLoadBalancer := h.Item.(*elb.LoadBalancerDescription)

	// Create service
	svc, err := ELBService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &elb.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: classicLoadBalancer.LoadBalancerName,
	}

	loadBalancerData, err := svc.DescribeLoadBalancerAttributes(params)
	if err != nil {
		return nil, err
	}

	return loadBalancerData, nil
}

func getAwsEc2ClassicLoadBalancerTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2ClassicLoadBalancerTags")
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	classicLoadBalancer := h.Item.(*elb.LoadBalancerDescription)

	// Create service
	svc, err := ELBService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &elb.DescribeTagsInput{
		LoadBalancerNames: []*string{aws.String(*classicLoadBalancer.LoadBalancerName)},
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

func getEc2ClassicLoadBalancerTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2ClassicLoadBalancerTurbotAkas")
	classicLoadBalancer := h.Item.(*elb.LoadBalancerDescription)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":elasticloadbalancing:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":loadbalancer/" + *classicLoadBalancer.LoadBalancerName}

	return akas, nil
}

//// TRANSFORM FUNCTIONS ////

func getEc2ClassicLoadBalancerTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	classicLoadBalancerTags := d.HydrateItem.([]*elb.Tag)

	if classicLoadBalancerTags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range classicLoadBalancerTags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
