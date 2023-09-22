package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlobalAcceleratorAccelerator(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_globalaccelerator_accelerator",
		Description: "AWS Global Accelerator Accelerator",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlobalAcceleratorAccelerator,
			Tags:    map[string]string{"service": "globalaccelerator", "action": "DescribeAccelerator"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGlobalAcceleratorAccelerators,
			Tags:    map[string]string{"service": "globalaccelerator", "action": "ListAccelerators"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getGlobalAcceleratorAcceleratorTags,
				Tags: map[string]string{"service": "globalaccelerator", "action": "ListTagsForResource"},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the accelerator.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the accelerator.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AcceleratorArn"),
			},
			{
				Name:        "created_time",
				Description: "The date and time that the accelerator was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "dns_name",
				Description: "The Domain Name System (DNS) name that Global Accelerator creates that points to your accelerator's static IP addresses.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enabled",
				Description: "Indicates whether the accelerator is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "ip_address_type",
				Description: "The value for the address type must be IPv4.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_sets",
				Description: "The static IP addresses that Global Accelerator associates with the accelerator.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "last_modified_time",
				Description: "The date and time that the accelerator was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status",
				Description: "Describes the deployment status of the accelerator.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the accelerator.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlobalAcceleratorAcceleratorTags,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "accelerator_attributes",
				Description: "Attributes of the accelerator.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlobalAcceleratorAcceleratorAttributes,
				Transform:   transform.FromField("AcceleratorAttributes"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlobalAcceleratorAcceleratorTags,
				Transform:   transform.FromField("Tags").Transform(globalacceleratorAcceleratorTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AcceleratorArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlobalAcceleratorAccelerators(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := GlobalAcceleratorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_accelerator.listGlobalAcceleratorAccelerators", "service_creation_error", err)
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
	input := &globalaccelerator.ListAcceleratorsInput{
		MaxResults: &maxItems,
	}

	paginator := globalaccelerator.NewListAcceleratorsPaginator(svc, input, func(o *globalaccelerator.ListAcceleratorsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_globalaccelerator_accelerator.listGlobalAcceleratorAccelerators", "api_error", err)
			return nil, err
		}

		for _, accelerator := range output.Accelerators {
			d.StreamListItem(ctx, accelerator)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlobalAcceleratorAccelerator(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	arn := d.EqualsQuals["arn"].GetStringValue()

	// check if arn is empty
	if arn == "" {
		return nil, nil
	}

	// Create session
	svc, err := GlobalAcceleratorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_accelerator.getGlobalAcceleratorAccelerator", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &globalaccelerator.DescribeAcceleratorInput{
		AcceleratorArn: aws.String(arn),
	}

	// Get call
	data, err := svc.DescribeAccelerator(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_accelerator.getGlobalAcceleratorAccelerator", "api_error", err)
		return nil, err
	}
	return *data.Accelerator, nil
}

func getGlobalAcceleratorAcceleratorTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accelerator := h.Item.(types.Accelerator)

	// Create Session
	svc, err := GlobalAcceleratorClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &globalaccelerator.ListTagsForResourceInput{
		ResourceArn: accelerator.AcceleratorArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_accelerator.getGlobalAcceleratorAcceleratorTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getGlobalAcceleratorAcceleratorAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accelerator := h.Item.(types.Accelerator)

	// Create Session
	svc, err := GlobalAcceleratorClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &globalaccelerator.DescribeAcceleratorAttributesInput{
		AcceleratorArn: accelerator.AcceleratorArn,
	}

	// Get call
	op, err := svc.DescribeAcceleratorAttributes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("getGlobalAcceleratorAcceleratorAttributes", "api_error", err)
		return nil, err
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func globalacceleratorAcceleratorTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*globalaccelerator.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
