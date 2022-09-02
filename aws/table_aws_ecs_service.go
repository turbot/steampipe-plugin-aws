package aws

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcsService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecs_service",
		Description: "AWS ECS Service",
		List: &plugin.ListConfig{
			Hydrate:       listEcsServices,
			ParentHydrate: listEcsClusters,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ClusterNotFoundException"}),
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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

//// LIST FUNCTION

func listEcsServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EcsService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get cluster details
	cluster := h.Item.(*ecs.Cluster)

	// DescribeServices API can describe up to 10 services in a single operation. Default MaxResults is 10 for ListServicesInput
	input := &ecs.ListServicesInput{
		Cluster:    cluster.ClusterArn,
		MaxResults: aws.Int64(10),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List all available ECS services
	var serviceNames [][]*string
	err = svc.ListServicesPages(
		input,
		func(page *ecs.ListServicesOutput, isLast bool) bool {
			if len(page.ServiceArns) != 0 {
				// Create a chunk of array of size 10
				serviceNames = append(serviceNames, page.ServiceArns)
			}
			return !isLast
		},
	)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	serviceCh := make(chan *ecs.DescribeServicesOutput, len(serviceNames))
	errorCh := make(chan error, len(serviceNames))

	for _, serviceData := range serviceNames {
		wg.Add(1)
		go getServiceDataAsync(serviceData, cluster.ClusterArn, svc, &wg, serviceCh, errorCh)
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
			d.StreamListItem(ctx, service)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

func getServiceDataAsync(serviceData []*string, clusterARN *string, svc *ecs.ECS, wg *sync.WaitGroup, serviceCh chan *ecs.DescribeServicesOutput, errorCh chan error) {
	defer wg.Done()
	rowData, err := getEcsService(serviceData, clusterARN, svc)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		serviceCh <- rowData
	}
}

// Describes the specified services running in your cluster.
// Below API can describe up to 10 services in a single operation.
func getEcsService(serviceData []*string, clusterARN *string, svc *ecs.ECS) (*ecs.DescribeServicesOutput, error) {
	params := &ecs.DescribeServicesInput{
		Services: serviceData,
		Cluster:  clusterARN,
	}
	response, err := svc.DescribeServices(params)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// List api call is not returning the tags for the service, so we need to make a separate api call for getting the tag details
func getEcsServiceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*ecs.Service)

	if data.ServiceArn == nil {
		return nil, nil
	}

	// Create Session
	svc, err := EcsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getEcsServiceTags", "connection_error", err)
		return nil, err
	}

	params := &ecs.ListTagsForResourceInput{
		ResourceArn: data.ServiceArn,
	}

	response, err := svc.ListTagsForResource(params)
	if err != nil {
		plugin.Logger(ctx).Error("getEcsServiceTags", err)
		return nil, err
	}

	return response.Tags, nil
}

//// TRANSFORM FUNCTIONS

func getEcsServiceTurbotTags(_ context.Context, d *transform.TransformData) (interface{},
	error) {
	tags := d.HydrateItem.([]*ecs.Tag)

	if len(tags) == 0 {
		return nil, nil
	}

	turbotTagsMap := map[string]string{}
	for _, i := range tags {
		turbotTagsMap[*i.Key] = *i.Value
	}
	return turbotTagsMap, nil
}
