package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlobalAcceleratorListener(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_globalaccelerator_listener",
		Description: "AWS Global Accelerator Listener",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException"}),
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
		Columns: awsGlobalRegionColumns([]*plugin.Column{
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
	AcceleratorArn string
	Listener       types.Listener
}

//// LIST FUNCTION

func listGlobalAcceleratorListeners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accelerator := h.Item.(types.Accelerator)
	acceleratorArn := aws.String(*accelerator.AcceleratorArn)

	// Create session
	svc, err := GlobalAcceleratorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_listener.listGlobalAcceleratorListeners", "connection_error", err)
		return nil, err
	}

	maxItems := int32(100)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input := &globalaccelerator.ListListenersInput{
		MaxResults:     &maxItems,
		AcceleratorArn: acceleratorArn,
	}

	paginator := globalaccelerator.NewListListenersPaginator(svc, input, func(o *globalaccelerator.ListListenersPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_globalaccelerator_listener.listGlobalAcceleratorListeners", "api_error", err)
			return nil, err
		}

		for _, listener := range output.Listeners {
			d.StreamListItem(ctx, &turbotListener{*acceleratorArn, listener})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlobalAcceleratorListener(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	arn := d.KeyColumnQuals["arn"].GetStringValue()

	// check if arn is empty
	if arn == "" {
		return nil, nil
	}

	// Create session
	svc, err := GlobalAcceleratorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_listener.getGlobalAcceleratorListener", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &globalaccelerator.DescribeListenerInput{
		ListenerArn: aws.String(arn),
	}

	// Get call
	data, err := svc.DescribeListener(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_listener.getGlobalAcceleratorListener", "api_error", err)
		return nil, err
	}
	return *data.Listener, nil
}
