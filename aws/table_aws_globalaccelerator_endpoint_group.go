package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/globalaccelerator"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlobalAcceleratorEndpointGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_globalaccelerator_endpoint_group",
		Description: "AWS Global Accelerator Endpoint Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlobalAcceleratorEndpointGroup,
		},
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "listener_arn", Require: plugin.Optional},
			},
			// TODO: Directly getting listeners would be better, but nested parent
			// hydrates are currently not working due to https://github.com/turbot/steampipe-plugin-sdk/issues/394
			//ParentHydrate: listGlobalAcceleratorListeners,
			ParentHydrate: listGlobalAcceleratorAccelerators,
			Hydrate:       listGlobalAcceleratorEndpointGroups,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the endpoint group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointGroup.EndpointGroupArn"),
			},
			{
				Name:        "listener_arn",
				Description: "The Amazon Resource Name (ARN) of parent listener.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "endpoint_descriptions",
				Description: "The list of endpoint objects.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EndpointGroup.EndpointDescriptions"),
			},
			{
				Name:        "endpoint_group_region",
				Description: "The AWS Region where the endpoint group is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointGroup.EndpointGroupRegion"),
			},
			{
				Name:        "health_check_interval_seconds",
				Description: "The time—10 seconds or 30 seconds—between health checks for each endpoint.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("EndpointGroup.HealthCheckIntervalSeconds"),
			},
			{
				Name:        "health_check_path",
				Description: "If the protocol is HTTP/S, then this value provides the ping path that Global Accelerator uses for the destination on the endpoints for health checks.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointGroup.HealthCheckPath"),
			},
			{
				Name:        "health_check_port",
				Description: "The port that Global Accelerator uses to perform health checks on endpoints that are part of this endpoint group.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("EndpointGroup.HealthCheckPort"),
			},
			{
				Name:        "health_check_protocol",
				Description: "The protocol that Global Accelerator uses to perform health checks on endpoints that are part of this endpoint group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointGroup.HealthCheckProtocol"),
			},
			{
				Name:        "port_overrides",
				Description: "Overrides for destination ports used to route traffic to an endpoint.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EndpointGroup.PortOverrides"),
			},
			{
				Name:        "threshold_count",
				Description: "The number of consecutive health checks required to set the state of a healthy endpoint to unhealthy, or to set an unhealthy endpoint to healthy.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("EndpointGroup.ThresholdCount"),
			},
			{
				Name:        "traffic_dial_percentage",
				Description: "The percentage of traffic to send to an AWS Region.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("EndpointGroup.TrafficDialPercentage"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointGroup.EndpointGroupArn").Transform(arnToTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EndpointGroup.EndpointGroupArn").Transform(arnToAkas),
			},
		}),
	}
}

type turbotEndpointGroup struct {
	ListenerArn   *string
	EndpointGroup *globalaccelerator.EndpointGroup
}

//// LIST FUNCTION

func listGlobalAcceleratorEndpointGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listGlobalAcceleratorEndpointGroups")

	accelerator := h.Item.(*globalaccelerator.Accelerator)
	acceleratorArn := aws.String(*accelerator.AcceleratorArn)

	// Create session
	svc, err := GlobalAcceleratorService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_endpoint_group.listGlobalAcceleratorEndpointGroups", "service_creation_error", err)
		return nil, err
	}

	// First get accelerator listener ARNs
	listenerArns := []*string{}

	listenersInput := &globalaccelerator.ListListenersInput{
		MaxResults:     aws.Int64(100),
		AcceleratorArn: acceleratorArn,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *listenersInput.MaxResults {
			if *limit < 1 {
				listenersInput.MaxResults = aws.Int64(1)
			} else {
				listenersInput.MaxResults = limit
			}
		}
	}

	// List listeners call
	err = svc.ListListenersPages(
		listenersInput,
		func(page *globalaccelerator.ListListenersOutput, isLast bool) bool {
			for _, listener := range page.Listeners {

				listenerArns = append(listenerArns, listener.ListenerArn)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_endpoint_group.listGlobalAcceleratorEndpointGroups", "api_error", err)
		return nil, err
	}

	// Now get endpoint groups for each listener
	for _, listenerArn := range listenerArns {

		endpointGroupsInput := &globalaccelerator.ListEndpointGroupsInput{
			MaxResults:  aws.Int64(100),
			ListenerArn: listenerArn,
		}

		// Reduce the basic request limit down if the user has only requested a small number of rows
		limit := d.QueryContext.Limit
		if d.QueryContext.Limit != nil {
			if *limit < *endpointGroupsInput.MaxResults {
				if *limit < 1 {
					endpointGroupsInput.MaxResults = aws.Int64(1)
				} else {
					endpointGroupsInput.MaxResults = limit
				}
			}
		}

		// List endpoint groups call
		err = svc.ListEndpointGroupsPages(
			endpointGroupsInput,
			func(page *globalaccelerator.ListEndpointGroupsOutput, isLast bool) bool {
				for _, endpointGroup := range page.EndpointGroups {
					d.StreamListItem(ctx, &turbotEndpointGroup{listenerArn, endpointGroup})

					// Context may get cancelled due to manual cancellation or if the limit has been reached
					if d.QueryStatus.RowsRemaining(ctx) == 0 {
						return false
					}
				}
				return !isLast
			},
		)

		if err != nil {
			plugin.Logger(ctx).Error("aws_globalaccelerator_endpoint_group.listGlobalAcceleratorEndpointGroups", "api_error", err)
			return nil, err
		}

	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlobalAcceleratorEndpointGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGlobalAcceleratorEndpointGroup")

	arn := d.KeyColumnQuals["arn"].GetStringValue()

	// check if arn is empty
	if arn == "" {
		return nil, nil
	}

	// Create session
	svc, err := GlobalAcceleratorService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_endpoint_group.getGlobalAcceleratorEndpointGroup", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &globalaccelerator.DescribeEndpointGroupInput{
		EndpointGroupArn: aws.String(arn),
	}

	// Get call
	data, err := svc.DescribeEndpointGroup(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_endpoint_group.getGlobalAcceleratorEndpointGroup", "api_error", err)
		return nil, err
	}
	return data.EndpointGroup, nil
}
