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

func tableAwsGlobalAcceleratorAccelerator(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_globalaccelerator_accelerator",
		Description: "AWS Global Accelerator Accelerator",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlobalAcceleratorAccelerator,
		},
		List: &plugin.ListConfig{
			Hydrate: listGlobalAcceleratorAccelerators,
		},
		Columns: awsColumns([]*plugin.Column{
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
	plugin.Logger(ctx).Trace("listGlobalAcceleratorAccelerators")

	// Create session
	svc, err := GlobalAcceleratorService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_accelerator.listGlobalAcceleratorAccelerators", "service_creation_error", err)
		return nil, err
	}

	input := &globalaccelerator.ListAcceleratorsInput{
		MaxResults: aws.Int64(100),
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
	err = svc.ListAcceleratorsPages(
		input,
		func(page *globalaccelerator.ListAcceleratorsOutput, isLast bool) bool {
			for _, accelerator := range page.Accelerators {
				d.StreamListItem(ctx, accelerator)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_accelerator.listGlobalAcceleratorAccelerators", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlobalAcceleratorAccelerator(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGlobalAcceleratorAccelerator")

	arn := d.KeyColumnQuals["arn"].GetStringValue()

	// check if arn is empty
	if arn == "" {
		return nil, nil
	}

	// Create session
	svc, err := GlobalAcceleratorService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_accelerator.getGlobalAcceleratorAccelerator", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &globalaccelerator.DescribeAcceleratorInput{
		AcceleratorArn: aws.String(arn),
	}

	// Get call
	data, err := svc.DescribeAccelerator(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_globalaccelerator_accelerator.getGlobalAcceleratorAccelerator", "api_error", err)
		return nil, err
	}
	return data.Accelerator, nil
}

func getGlobalAcceleratorAcceleratorTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getGlobalAcceleratorAcceleratorTags")

	accelerator := h.Item.(*globalaccelerator.Accelerator)

	// Create Session
	svc, err := GlobalAcceleratorService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &globalaccelerator.ListTagsForResourceInput{
		ResourceArn: accelerator.AcceleratorArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Debug("getGlobalAcceleratorAcceleratorTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getGlobalAcceleratorAcceleratorAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getGlobalAcceleratorAcceleratorAttributes")

	accelerator := h.Item.(*globalaccelerator.Accelerator)

	// Create Session
	svc, err := GlobalAcceleratorService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &globalaccelerator.DescribeAcceleratorAttributesInput{
		AcceleratorArn: accelerator.AcceleratorArn,
	}

	// Get call
	op, err := svc.DescribeAcceleratorAttributes(params)
	if err != nil {
		logger.Debug("getGlobalAcceleratorAcceleratorAttributes", "api_error", err)
		return nil, err
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func globalacceleratorAcceleratorTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("globalacceleratorAcceleratorTurbotTags")

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
