package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLightsailInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lightsail_instance",
		Description: "AWS Lightsail Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getLightsailInstance,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidResourceName", "DoesNotExist"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listLightsailInstances,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "blueprint_id",
				Description: "The blueprint ID (e.g., os_amlinux_2016_03).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "blueprint_name",
				Description: "The friendly name of the blueprint (e.g., Amazon Linux).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bundle_id",
				Description: "The bundle for the instance (e.g., micro_1_0).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The timestamp when the instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "hardware",
				Description: "The size of the vCPU and the amount of RAM for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ip_address_type",
				Description: "The IP address type of the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ip_v6_addresses",
				Description: "The IPv6 addresses of the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "is_static_ip",
				Description: "A Boolean value indicating whether this instance has a static IP assigned to it.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "availability_zone",
				Description: "The Availability Zone where the instance is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location.AvailabilityZone"),
			},
			{
				Name:        "metadata_options",
				Description: "The metadata options for the Amazon Lightsail instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "networking",
				Description: "Information about the public ports and monthly data transfer rates for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_ip_address",
				Description: "The private IP address of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_ip_address",
				Description: "The public IP address of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The type of resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssh_key_name",
				Description: "The name of the SSH key being used to connect to the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_code",
				Description: "The status code for the instance.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("State.Code"),
			},
			{
				Name:        "state_name",
				Description: "The status of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State.Name"),
			},
			{
				Name:        "support_code",
				Description: "The support code.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "username",
				Description: "The user name for connecting to the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
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
				Hydrate:     getLightsailInstance,
				Transform:   transform.FromField("Tags").Transform(getLightsailInstanceTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getLightsailInstance,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listLightsailInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := LightsailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_instance.listLightsailInstances", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &lightsail.GetInstancesInput{}

	// List call
	for {
		resp, err := svc.GetInstances(ctx, input)

		if err != nil {
			plugin.Logger(ctx).Error("aws_lightsail_instance.listLightsailInstances", "query_error", err)
			return nil, nil
		}

		for _, item := range resp.Instances {
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextPageToken != nil {
			input.PageToken = resp.NextPageToken
		} else {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getLightsailInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := LightsailClient(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var name string
	if h.Item != nil {
		name = *h.Item.(types.Instance).Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	params := &lightsail.GetInstanceInput{
		InstanceName: aws.String(name),
	}

	detail, err := svc.GetInstance(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_instance.getLightsailInstance", "api_error", err)
		return nil, err
	}
	return detail.Instance, nil
}

//// TRANSFORM FUNCTIONS

func getLightsailInstanceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.Value.([]types.Tag)
	var turbotTagsMap map[string]string
	if tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}
