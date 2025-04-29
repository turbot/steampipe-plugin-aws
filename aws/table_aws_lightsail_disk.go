package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"

	lightsailv1 "github.com/aws/aws-sdk-go/service/lightsail"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLightsailDisk(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lightsail_disk",
		Description: "AWS Lightsail Disk",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getLightsailDisk,
			Tags:       map[string]string{"service": "lightsail", "action": "GetDisk"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidResourceName", "DoesNotExist"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listLightsailDisks,
			Tags:    map[string]string{"service": "lightsail", "action": "GetDisks"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(lightsailv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The unique name of the disk.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the disk.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The date when the disk was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "location",
				Description: "The AWS Region and Availability Zone where the disk is located.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "is_attached",
				Description: "A Boolean value indicating whether the disk is attached.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "iops",
				Description: "The input/output operations per second (IOPS) of the disk.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "path",
				Description: "The disk path.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size_in_gb",
				Description: "The size of the disk in GB.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "state",
				Description: "The status of the disk.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "support_code",
				Description: "The support code. Include this code in your email to support when you have questions about your Lightsail disk.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the disk.",
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
				Transform:   transform.FromField("Tags").Transform(getLightsailDiskTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listLightsailDisks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := LightsailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_disk.listLightsailDisks", "connection_error", err)
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

	input := &lightsail.GetDisksInput{}

	// List call
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := svc.GetDisks(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lightsail_disk.listLightsailDisks", "api_error", err)
			return nil, err
		}

		for _, disk := range output.Disks {
			d.StreamListItem(ctx, disk)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if output.NextPageToken == nil {
			break
		}
		input.PageToken = output.NextPageToken
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getLightsailDisk(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()

	// Create Session
	svc, err := LightsailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_disk.getLightsailDisk", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &lightsail.GetDiskInput{
		DiskName: aws.String(name),
	}

	op, err := svc.GetDisk(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_disk.getLightsailDisk", "api_error", err)
		return nil, err
	}

	return op.Disk, nil
}

func getLightsailDiskTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	var disk types.Disk
	switch v := d.HydrateItem.(type) {
	case types.Disk:
		disk = v
	case *types.Disk:
		disk = *v
	default:
		plugin.Logger(ctx).Error("aws_lightsail_disk.getLightsailDiskTurbotTags", "unexpected_type", fmt.Sprintf("%T", v))
		return nil, fmt.Errorf("unexpected type %T", v)
	}

	turbotTagsMap := map[string]string{}
	if disk.Tags != nil {
		for _, i := range disk.Tags {
			if i.Key != nil && i.Value != nil {
				turbotTagsMap[*i.Key] = *i.Value
			}
		}
	}

	plugin.Logger(ctx).Debug("aws_lightsail_disk.getLightsailDiskTurbotTags", "disk_name", disk.Name, "tags", turbotTagsMap)
	return turbotTagsMap, nil
}
