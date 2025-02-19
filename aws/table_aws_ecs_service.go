package aws

import (
	"context"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	ecsv1 "github.com/aws/aws-sdk-go/service/ecs"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcsService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_service",
		Description: "AWS ECS Service",
		List: &plugin.ListConfig{
			Hydrate:       listEcsServices,
			Tags:          map[string]string{"service": "ecs", "action": "ListServices"},
			ParentHydrate: listEcsClusters,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ClusterNotFoundException"}),
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getEcsServiceTags,
				Tags: map[string]string{"service": "ecs", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ecsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "service_name",
				Description: "The name of the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceArn"),
			},
			{
				Name:        "status",
				Description: "The status of the service. Valid values are: ACTIVE, DRAINING, or INACTIVE.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_arn",
				Description: "The Amazon Resource Name (ARN) of the cluster that hosts the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "task_definition",
				Description: "The task definition to use for tasks in the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The date and time when the service was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by",
				Description: "The principal that created the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deployment_controller_type",
				Description: "The deployment controller type to use. Possible values are: ECS, CODE_DEPLOY, and EXTERNAL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "desired_count",
				Description: "The desired number of instantiations of the task definition to keep running on the service.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "enable_ecs_managed_tags",
				Description: "Specifies whether to enable Amazon ECS managed tags for the tasks in the service.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("EnableECSManagedTags"),
			},
			{
				Name:        "enable_execute_command",
				Description: "Indicates whether or not the execute command functionality is enabled for the service.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "health_check_grace_period_seconds",
				Description: "The period of time, in seconds, that the Amazon ECS service scheduler ignores unhealthy Elastic Load Balancing target health checks after a task has first started.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "launch_type",
				Description: "The launch type on which your service is running. If no value is specified, it will default to EC2.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pending_count",
				Description: "The number of tasks in the cluster that are in the PENDING state.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "platform_family",
				Description: "The operating system that your tasks in the service run on.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "platform_version",
				Description: "The platform version on which to run your service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "propagate_tags",
				Description: "Specifies whether to propagate the tags from the task definition or the service to the task. If no value is specified, the tags are not propagated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role_arn",
				Description: "The ARN of the IAM role associated with the service that allows the Amazon ECS container agent to register container instances with an Elastic Load Balancing load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "running_count",
				Description: "The number of tasks in the cluster that are in the RUNNING state.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "scheduling_strategy",
				Description: "The scheduling strategy to use for the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capacity_provider_strategy",
				Description: "The capacity provider strategy associated with the service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "deployment_configuration",
				Description: "Optional deployment parameters that control how many tasks run during the deployment and the ordering of stopping and starting tasks.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "deployments",
				Description: "The current state of deployments for the service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "events",
				Description: "The event stream for your service. A maximum of 100 of the latest events are displayed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "load_balancers",
				Description: "A list of Elastic Load Balancing load balancer objects, containing the load balancer name, the container name (as it appears in a container definition), and the container port to access from the load balancer.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_configuration",
				Description: "The VPC subnet and security group configuration for tasks that receive their own elastic network interface by using the awsvpc networking mode.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "placement_constraints",
				Description: "The placement constraints for the tasks in the service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "placement_strategy",
				Description: "The placement strategy that determines how tasks for the service are placed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "service_registries",
				Description: "The details of the service discovery registries to assign to this service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "task_sets",
				Description: "Information about a set of Amazon ECS tasks in either an AWS CodeDeploy or an EXTERNAL deployment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "The metadata that you apply to the service to help you categorize and organize them.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsServiceTags,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsServiceTags,
				Transform:   transform.From(getEcsServiceTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type ServiceInfo struct {
	ClusterName *string
	types.Service
}

//// LIST FUNCTION

func listEcsServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ECSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_service.listEcsServices", "connection_error", err)
		return nil, err
	}

	cluster := h.Item.(types.Cluster)

	// Limiting the results
	maxLimit := int32(10)
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
	// DescribeServices API can describe up to 10 services in a single operation. Default MaxResults is 10 for ListServicesInput
	input := &ecs.ListServicesInput{
		Cluster:    cluster.ClusterArn,
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := ecs.NewListServicesPaginator(svc, input, func(o *ecs.ListServicesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	var serviceNames []string
	// List all available ECS services
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ecs_service.listEcsServices", "api_error", err)
			return nil, err
		}

		serviceNames = append(serviceNames, output.ServiceArns...)
	}

	var wg sync.WaitGroup
	serviceCh := make(chan *ecs.DescribeServicesOutput, len(serviceNames))
	errorCh := make(chan error, len(serviceNames))

	for _, serviceData := range serviceNames {
		wg.Add(1)
		go getServiceDataAsync(serviceData, cluster.ClusterArn, svc, &wg, serviceCh, errorCh, ctx)
	}

	// wait for all services to be processed
	wg.Wait()

	// NOTE: close channel before ranging over results
	close(serviceCh)
	close(errorCh)

	for err := range errorCh {
		// return the first error
		return nil, err
	}

	for result := range serviceCh {
		for _, service := range result.Services {
			d.StreamListItem(ctx, ServiceInfo{cluster.ClusterName, service})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

func getServiceDataAsync(serviceData string, clusterARN *string, svc *ecs.Client, wg *sync.WaitGroup, serviceCh chan *ecs.DescribeServicesOutput, errorCh chan error, ctx context.Context) {
	defer wg.Done()
	rowData, err := getEcsService(serviceData, clusterARN, svc, ctx)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		serviceCh <- rowData
	}
}

// Describes the specified services running in your cluster.
// Below API can describe up to 10 services in a single operation.
func getEcsService(serviceData string, clusterARN *string, svc *ecs.Client, ctx context.Context) (*ecs.DescribeServicesOutput, error) {
	params := &ecs.DescribeServicesInput{
		Services: []string{serviceData},
		Cluster:  clusterARN,
	}
	response, err := svc.DescribeServices(ctx, params)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// List api call is not returning the tags for the service, so we need to make a separate api call for getting the tag details
func getEcsServiceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(ServiceInfo)

	if data.ServiceArn == nil {
		return nil, nil
	}

	resourceArn := *data.ServiceArn

	// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-account-settings.html#ecs-resource-ids
	// We can have two arn format for the ESC Service:
	// 1. arn:aws:ecs:<region>:<account_id>:service/<cluster_name>/<service_name> (Newer format)
	// 2. arn:aws:ecs:<region>:<account_id>:service/<service_name> (Older format)
	// While making the API with older service arn we are encountering the error:
	// ERROR: rpc error: code = Unknown desc = my_aws_account: table 'aws_ecs_service' column 'tags' requires hydrate data from getEcsServiceTags,
	// which failed with error operation error ECS: ListTagsForResource, https response error StatusCode: 400,
	// RequestID: 076ed52f-8f0e-43b9-af89-3728995bb52b, InvalidParameterException: Long arn format must be used for tagging operations.
	resourceArnSplitPart := strings.Split(resourceArn, "/")
	if len(resourceArnSplitPart) < 3 {
		// The List Clusters API does not return the cluster name. Therefore, for clusters using the older ARN format, we make a Get API call to retrieve the cluster name.
		// We need the cluster name to make the List tag API call for service.
		h.Item = types.Cluster{ClusterArn: data.ClusterArn}
		res, err := getEcsCluster(ctx, d, h)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ecs_service.getEcsServiceTags.getEcsCluster", "api_error", err)
			return nil, err
		}

		cluster := res.(types.Cluster)
		data.ClusterName = cluster.ClusterName

		resourceArn = resourceArnSplitPart[0] + "/" + *data.ClusterName + "/" + *data.ServiceName
	}

	// Create Session
	svc, err := ECSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_service.getEcsServiceTags", "connection_error", err)
		return nil, err
	}

	params := &ecs.ListTagsForResourceInput{
		ResourceArn: &resourceArn,
	}

	response, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecs_service.getEcsServiceTags", "api_error", err)
		return nil, err
	}

	return response.Tags, nil
}

//// TRANSFORM FUNCTIONS

func getEcsServiceTurbotTags(_ context.Context, d *transform.TransformData) (interface{},
	error) {
	tags := d.HydrateItem.([]types.Tag)

	if len(tags) == 0 {
		return nil, nil
	}

	turbotTagsMap := map[string]string{}
	for _, i := range tags {
		turbotTagsMap[*i.Key] = *i.Value
	}
	return turbotTagsMap, nil
}
