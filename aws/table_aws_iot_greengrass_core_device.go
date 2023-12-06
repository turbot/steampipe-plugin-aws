package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/greengrassv2"
	"github.com/aws/aws-sdk-go-v2/service/greengrassv2/types"

	greengrassv1 "github.com/aws/aws-sdk-go/service/greengrassv2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIotGreengrassCoreDevice(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iot_greengrass_core_device",
		Description: "AWS IoT Greengrass Core Device",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("core_device_thing_name"),
			Hydrate:    getIotGreengrassCoreDevice,
			Tags:       map[string]string{"service": "greengrassv2", "action": "GetCoreDevice"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listIotGreengrassCoreDevices,
			Tags:    map[string]string{"service": "greengrassv2", "action": "ListCoreDevices"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "status", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(greengrassv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "core_device_thing_name",
				Description: "The name of the core device. This is also the name of the IoT thing.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "architecture",
				Description: "The computer architecture of the core device.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIotGreengrassCoreDevice,
			},
			{
				Name:        "core_version",
				Description: "The version of the IoT Greengrass Core software that the core device runs.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIotGreengrassCoreDevice,
			},
			{
				Name:        "platform",
				Description: "The operating system platform that the core device runs.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIotGreengrassCoreDevice,
			},
			{
				Name:        "last_status_update_timestamp",
				Description: "The time at which the core device's status last updated, expressed in ISO 8601 format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status",
				Description: "The status of the core device. Core devices can have the following statuses HEALTHY, UNHEALTHY.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CoreDeviceThingName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIotGreengrassCoreDevice,
			},
		}),
	}
}

//// LIST FUNCTION

func listIotGreengrassCoreDevices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IoTGreengrassClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_greengrass_core_device.listIotGreengrassCoreDevices", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
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

	input := &greengrassv2.ListCoreDevicesInput{
		MaxResults: aws.Int32(maxLimit),
	}
	if d.EqualsQualString("status") != "" {
		input.Status = types.CoreDeviceStatus(d.EqualsQualString("status"))
	}

	paginator := greengrassv2.NewListCoreDevicesPaginator(svc, input, func(o *greengrassv2.ListCoreDevicesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iot_greengrass_core_device.listIotGreengrassCoreDevices", "api_error", err)
			return nil, err
		}

		for _, item := range output.CoreDevices {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getIotGreengrassCoreDevice(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	coreDeviceThingName := ""
	if h.Item != nil {
		t := h.Item.(types.CoreDevice)
		coreDeviceThingName = *t.CoreDeviceThingName
	} else {
		coreDeviceThingName = d.EqualsQualString("core_device_thing_name")
	}

	if coreDeviceThingName == "" {
		return nil, nil
	}

	// Create service
	svc, err := IoTGreengrassClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_greengrass_core_device.getIotGreengrassCoreDevice", "connection_error", err)
		return nil, err
	}

	params := &greengrassv2.GetCoreDeviceInput{
		CoreDeviceThingName: aws.String(coreDeviceThingName),
	}

	resp, err := svc.GetCoreDevice(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_greengrass_core_device.getIotGreengrassCoreDevice", "api_error", err)
		return nil, err
	}

	return resp, nil
}
