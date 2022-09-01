package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/globalaccelerator"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlobalacceleratorListener(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_globalaccelerator_listener",
		Description: "AWS Global Accelerator Listener",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlobalAcceleratorListener,
		},
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "accelerator_arn", Require: plugin.Optional},
			},
			ParentHydrate: listGlobalAcceleratorAccelerators,
			Hydrate:       listGlobalAcceleratorListeners,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the listener.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Listener.ListenerArn"),
			},
			{
				Name:        "accelerator_arn",
				Description: "The Amazon Resource Name (ARN) of parent accelerator.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "client_affinity",
				Description: "Client affinity setting for the listener.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Listener.ClientAffinity"),
			},
			{
				Name:        "port_ranges",
				Description: "The list of port ranges for the connections from clients to the accelerator.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Listener.PortRanges"),
			},
			{
				Name:        "protocol",
				Description: "The protocol for the connections from clients to the accelerator.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Listener.Protocol"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Listener.ListenerArn").Transform(arnToTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Listener.ListenerArn").Transform(arnToAkas),
			},
		}),
	}
}

type turbotListener struct {
	AcceleratorArn *string
	Listener       *globalaccelerator.Listener
}

//// LIST FUNCTION

func listGlobalAcceleratorListeners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listGlobalAcceleratorListeners")

	accelerator := h.Item.(*globalaccelerator.Accelerator)
	acceleratorArn := aws.String(*accelerator.AcceleratorArn)

	// Create session
	svc, err := GlobalAcceleratorService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_listener.listGlobalAcceleratorListeners", "service_creation_error", err)
		return nil, err
	}

	input := &globalaccelerator.ListListenersInput{
		MaxResults:     aws.Int64(100),
		AcceleratorArn: acceleratorArn,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
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

	// List call
	err = svc.ListListenersPages(
		input,
		func(page *globalaccelerator.ListListenersOutput, isLast bool) bool {
			for _, listener := range page.Listeners {
				d.StreamListItem(ctx, &turbotListener{acceleratorArn, listener})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_listener.listGlobalAcceleratorListeners", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlobalAcceleratorListener(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGlobalAcceleratorListener")

	arn := d.KeyColumnQuals["arn"].GetStringValue()

	// check if arn is empty
	if arn == "" {
		return nil, nil
	}

	// Create session
	svc, err := GlobalAcceleratorService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_listener.getGlobalAcceleratorListener", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &globalaccelerator.DescribeListenerInput{
		ListenerArn: aws.String(arn),
	}

	// Get call
	data, err := svc.DescribeListener(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_listener.getGlobalAcceleratorListener", "api_error", err)
		return nil, err
	}
	return data.Listener, nil
}
