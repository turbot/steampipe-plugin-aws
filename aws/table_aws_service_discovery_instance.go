package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"

	servicediscoveryv1 "github.com/aws/aws-sdk-go/service/servicediscovery"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServiceDiscoveryInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_service_discovery_instance",
		Description: "AWS Service Discovery Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"service_id", "id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InstanceNotFound"}),
			},
			Hydrate: getServiceDiscoveryInstance,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listServiceDiscoveryServices,
			Hydrate:       listServiceDiscoveryInstances,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "service_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(servicediscoveryv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_id",
				Description: "The ID of the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ec2_instance_id",
				Description: "The Amazon EC2 instance ID for the instance. When the AWS_EC2_INSTANCE_ID attribute is specified, then the AWS_INSTANCE_IPV4 attribute contains the primary private IPv4 address.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.AWS_EC2_INSTANCE_ID"),
			},
			{
				Name:        "alias_dns_name",
				Description: "For an alias record that routes traffic to an Elastic Load Balancing load balancer, the DNS name that's associated with the load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_cname",
				Description: "A CNAME record, the domain name that Route 53 returns in response to DNS queries (for example, example.com ).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.AWS_INSTANCE_CNAME"),
			},
			{
				Name:        "init_health_status",
				Description: "If the service configuration includes HealthCheckCustomConfig, you can optionally use AWS_INIT_HEALTH_STATUS to specify the initial status of the custom health check, HEALTHY or UNHEALTHY. If you don't specify a value for AWS_INIT_HEALTH_STATUS, the initial status is HEALTHY.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.AWS_INIT_HEALTH_STATUS"),
			},
			{
				Name:        "instance_ipv4",
				Description: "For an A record, the IPv4 address that Route 53 returns in response to DNS queries.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("Attributes.AWS_INSTANCE_IPV4"),
			},
			{
				Name:        "instance_ipv6",
				Description: "For an AAAA record, the IPv6 address that Route 53 returns in response to DNS queries.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("Attributes.AWS_INSTANCE_IPV6"),
			},
			{
				Name:        "instance_port",
				Description: "For an SRV record, the value that Route 53 returns for the port. In addition, if the service includes HealthCheckConfig, the port on the endpoint that Route 53 sends requests to.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Attributes.AWS_INSTANCE_PORT"),
			},
			{
				Name:        "attributes",
				Description: "Attributes of the instance.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
		}),
	}
}

type ServiceInstanceInfo struct {
	Id         *string
	ServiceId  *string
	Attributes map[string]string
}

//// LIST FUNCTION

func listServiceDiscoveryInstances(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service := h.Item.(types.ServiceSummary)

	// Restrict API call for other service IDs
	if d.EqualsQualString("service_id") != "" {
		if d.EqualsQualString("service_id") != *service.Id {
			return nil, nil
		}
	}

	// Create Client
	svc, err := ServiceDiscoveryClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_service_discovery_instance.listServiceDiscoveryInstances", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &servicediscovery.ListInstancesInput{
		ServiceId:  service.Id,
		MaxResults: &maxLimit,
	}

	paginator := servicediscovery.NewListInstancesPaginator(svc, input, func(o *servicediscovery.ListInstancesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_service_discovery_instance.listServiceDiscoveryInstances", "api_error", err)
			return nil, err
		}

		for _, item := range output.Instances {
			d.StreamListItem(ctx, &ServiceInstanceInfo{
				Id:         item.Id,
				ServiceId:  service.Id,
				Attributes: item.Attributes,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getServiceDiscoveryInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	serviceId := d.EqualsQualString("service_id")
	instanceId := d.EqualsQualString("id")

	// Empty check
	if serviceId == "" || instanceId == "" {
		return nil, nil
	}

	// Create client
	svc, err := ServiceDiscoveryClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_service_discovery_instance.getServiceDiscoveryInstance", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &servicediscovery.GetInstanceInput{
		InstanceId: &instanceId,
		ServiceId:  &serviceId,
	}

	// Get call
	op, err := svc.GetInstance(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_service_discovery_instance.getServiceDiscoveryInstance", "api_error", err)
		return nil, err
	}

	if op.Instance != nil {
		return &ServiceInstanceInfo{Id: op.Instance.Id, ServiceId: &serviceId, Attributes: op.Instance.Attributes}, nil
	}

	return nil, nil
}
