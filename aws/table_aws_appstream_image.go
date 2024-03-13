package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/aws/aws-sdk-go-v2/service/appstream/types"
	appstreamv1 "github.com/aws/aws-sdk-go/service/appstream"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAppStreamImage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_appstream_image",
		Description: "AWS AppStream Image",
		List: &plugin.ListConfig{
			Hydrate: listAppStreamImages,
			Tags:    map[string]string{"service": "appstream", "action": "DescribeImages"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "arn",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAppStreamTags,
				Tags: map[string]string{"service": "appstream", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(appstreamv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "appstream_agent_version",
				Description: "The version of the AppStream 2.0 agent to use for instances that are launched from this image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "base_image_arn",
				Description: "The ARN of the image from which this image was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The time the image was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description to display.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The image name to display.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_builder_name",
				Description: "The name of the image builder that was used to create the private image. If the image is shared, this value is null.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_builder_supported",
				Description: "Indicates whether an image builder can be launched from this image.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "platform",
				Description: "The operating system platform of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_base_image_released_date",
				Description: "The release date of the public base image. For private images, this date is the release date of the base image from which the image was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "state",
				Description: "The image starts in the PENDING state. If image creation succeeds, the state is AVAILABLE. If image creation fails, the state is FAILED.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "visibility",
				Description: "Indicates whether the image is public or private.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "applications",
				Description: "The applications associated with the image.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "image_errors",
				Description: "Describes the errors that are returned when a new image can't be created.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "image_permissions",
				Description: "The permissions to provide to the destination AWS account for the specified image.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "state_change_reason",
				Description: "The reason why the last state change occurred.",
				Type:        proto.ColumnType_JSON,
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
				Hydrate:     getAppStreamTags,
				Transform:   transform.FromValue(),
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

func listAppStreamImages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := AppStreamClient(ctx, d)
	if err != nil {
		logger.Error("aws_appstream_image.listAppStreamImages", "connection_error", err)
		return nil, err
	}

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(25)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	params := &appstream.DescribeImagesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if d.Quals["name"] != nil {
		for _, q := range d.Quals["name"].Quals {
			value := q.Value.GetStringValue()
			if q.Operator == "=" {
				params.Names = append(params.Names, value)
			}
		}
	}
	if d.Quals["arn"] != nil {
		for _, q := range d.Quals["arn"].Quals {
			value := q.Value.GetStringValue()
			if q.Operator == "=" {
				params.Arns = append(params.Arns, value)
			}
		}
	}

	paginator := appstream.NewDescribeImagesPaginator(svc, params, func(o *appstream.DescribeImagesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_appstream_image.listAppStreamImages", "api_error", err)
			return nil, err
		}
		for _, image := range output.Images {
			d.StreamListItem(ctx, image)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAppStreamTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.Image).Arn
	} else {
		return nil, nil
	}

	// Create Session
	svc, err := AppStreamClient(ctx, d)
	if err != nil {
		logger.Error("aws_appstream_image.getAppStreamTags", "connection_error", err)
		return nil, err
	}

	params := &appstream.ListTagsForResourceInput{
		ResourceArn: &arn,
	}

	tags, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_appstream_image.getAppStreamTags", "api_error", err)
		return nil, err
	}

	return tags.Tags, nil
}
