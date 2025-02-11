package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediastore"
	"github.com/aws/aws-sdk-go-v2/service/mediastore/types"

	"github.com/aws/smithy-go"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsMediaStoreContainer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_media_store_container",
		Description: "AWS Media Store Container",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameter", "ContainerNotFoundException", "ContainerInUseException"}),
			},
			Hydrate: getMediaStoreContainer,
			Tags:    map[string]string{"service": "mediastore", "action": "DescribeContainer"},
		},
		List: &plugin.ListConfig{
			Hydrate: listMediaStoreContainers,
			Tags:    map[string]string{"service": "mediastore", "action": "ListContainers"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ContainerInUseException"}),
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getMediaStoreContainerPolicy,
				Tags: map[string]string{"service": "mediastore", "action": "GetContainerPolicy"},
			},
			{
				Func: listMediaStoreContainerTags,
				Tags: map[string]string{"service": "mediastore", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_MEDIASTORE_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the container.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the container.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "status",
				Description: "The status of container creation or deletion. The status is one of the following: 'CREATING', 'ACTIVE', or 'DELETING'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "access_logging_enabled",
				Description: "The state of access logging on the container. This value is false by default, indicating that AWS Elemental MediaStore does not send access logs to Amazon CloudWatch Logs. When you enable access logging on the container, MediaStore changes this value to true, indicating that the service delivers access logs for objects stored in that container to CloudWatch Logs.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "creation_time",
				Description: "The Unix timestamp that the container was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "endpoint",
				Description: "The DNS endpoint of the container.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy",
				Description: "The contents of the access policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMediaStoreContainerPolicy,
			},
			{
				Name:        "policy_std",
				Description: "Contains the contents of the access policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMediaStoreContainerPolicy,
				Transform:   transform.FromField("Policy").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the container",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMediaStoreContainerTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
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
				Hydrate:     listMediaStoreContainerTags,
				Transform:   transform.From(containerTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listMediaStoreContainers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create service
	svc, err := MediaStoreClient(ctx, d)
	if err != nil {
		logger.Error("aws_media_store_container.listMediaStoreContainers", "connection_error", err)
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
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	// Set MaxResults to the maximum number allowed
	input := &mediastore.ListContainersInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := mediastore.NewListContainersPaginator(svc, input, func(o *mediastore.ListContainersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_media_store_container.listMediaStoreContainers", "api_error", err)
			return nil, err
		}

		for _, items := range output.Containers {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getMediaStoreContainer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	containerName := d.EqualsQuals["name"].GetStringValue()
	// check if name is empty
	if containerName == "" {
		return nil, nil
	}

	// Create service
	svc, err := MediaStoreClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_media_store_container.getMediaStoreContainer", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &mediastore.DescribeContainerInput{
		ContainerName: &containerName,
	}

	// Get call
	data, err := svc.DescribeContainer(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_media_store_container.getMediaStoreContainer", "api_error", err)
		return nil, err
	}

	return *data.Container, nil
}

func getMediaStoreContainerPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var containerName string
	if h.Item != nil {
		containerName = *h.Item.(types.Container).Name
	}

	// Create Session
	svc, err := MediaStoreClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_media_store_container.getMediaStoreContainerPolicy", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &mediastore.GetContainerPolicyInput{
		ContainerName: &containerName,
	}

	// Get call
	data, err := svc.GetContainerPolicy(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if helpers.StringSliceContains([]string{"PolicyNotFoundException", "ContainerInUseException"}, ae.ErrorCode()) {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_media_store_container.getMediaStoreContainerPolicy", "api_error", err)
		return nil, err
	}

	return data, nil
}

func listMediaStoreContainerTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.Container).ARN
	}

	// Create Session
	svc, err := MediaStoreClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_media_store_container.listMediaStoreContainerTags", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &mediastore.ListTagsForResourceInput{
		Resource: &arn,
	}

	// Get call
	data, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ContainerInUseException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_media_store_container.listMediaStoreContainerTags", "api_error", err)
		return nil, err
	}

	return data, nil
}

//// TRANSFORM FUNCTION

func containerTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {

	// When container is being created or deleted, 'ContainerInUseException' is thrown.
	// When 'ContainerInUseException' is thrown, the hydrated data comes as null.
	// As a result, panic interface conversion error is thrown.
	switch item := d.HydrateItem.(type) {
	case *mediastore.ListTagsForResourceOutput:
		tags := item.Tags
		var turbotTagsMap map[string]string
		if tags != nil {
			turbotTagsMap = map[string]string{}
			for _, i := range tags {
				turbotTagsMap[*i.Key] = *i.Value
			}
		}

		return turbotTagsMap, nil
	}

	return nil, nil
}
