package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/mediastore"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsMediaStoreContainer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_media_store_container",
		Description: "AWS Media Store Container",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidParameter", "ContainerNotFoundException", "ContainerInUseException"}),
			},
			Hydrate: getMediaStoreContainer,
		},
		List: &plugin.ListConfig{
			Hydrate: listMediaStoreContainers,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ContainerInUseException"}),
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := MediaStoreService(ctx, d)
	if err != nil {
		logger.Error("listMediaStoreContainers", "error_MediaStoreService", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Set MaxResults to the maximum number allowed
	input := mediastore.ListContainersInput{
		MaxResults: aws.Int64(100),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
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

	err = svc.ListContainersPages(
		&input,
		func(page *mediastore.ListContainersOutput, lastPage bool) bool {
			for _, container := range page.Containers {
				d.StreamListItem(ctx, container)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	if err != nil {
		logger.Error("listMediaStoreContainers", "error_ListContainersPages", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getMediaStoreContainer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	containerName := d.KeyColumnQuals["name"].GetStringValue()
	// check if name is empty
	if containerName == "" {
		return nil, nil
	}

	// Create service
	svc, err := MediaStoreService(ctx, d)
	if err != nil {
		logger.Error("getMediaStoreContainer", "error_MediaStoreService", err)
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
	data, err := svc.DescribeContainer(params)
	if err != nil {
		logger.Error("getMediaStoreContainer", "error_DescribeContainer", err)
		return nil, err
	}

	return data.Container, nil
}

func getMediaStoreContainerPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getMediaStoreContainerPolicy")

	var containerName string
	if h.Item != nil {
		containerName = *h.Item.(*mediastore.Container).Name
	}

	// Create Session
	svc, err := MediaStoreService(ctx, d)
	if err != nil {
		logger.Error("getMediaStoreContainerPolicy", "error_MediaStoreService", err)
		return nil, err
	}

	// Build the params
	params := &mediastore.GetContainerPolicyInput{
		ContainerName: &containerName,
	}

	// Get call
	data, err := svc.GetContainerPolicy(params)
	if err != nil {
		logger.Error("getMediaStoreContainerPolicy", "error_GetContainerPolicy", err)
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "PolicyNotFoundException" || a.Code() == "ContainerInUseException" {
				return nil, nil
			}
		}
		return nil, err
	}

	return data, nil
}

func listMediaStoreContainerTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listMediaStoreContainerTags")

	var arn string
	if h.Item != nil {
		arn = *h.Item.(*mediastore.Container).ARN
	}

	// Create Session
	svc, err := MediaStoreService(ctx, d)
	if err != nil {
		logger.Error("listMediaStoreContainerTags", "error_MediaStoreService", err)
		return nil, err
	}

	// Build the params
	params := &mediastore.ListTagsForResourceInput{
		Resource: &arn,
	}

	// Get call
	data, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Error("listMediaStoreContainerTags", "error_ListTagsForResource", err)
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ContainerInUseException" {
				return nil, nil
			}
		}
		return nil, err
	}

	return data, nil
}

//// TRANSFORM FUNCTION

func containerTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("containerTurbotTags")

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
