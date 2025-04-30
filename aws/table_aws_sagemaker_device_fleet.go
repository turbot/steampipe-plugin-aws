package aws

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"

	sagemakerv1 "github.com/aws/aws-sdk-go/service/sagemaker"

	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSageMakerDeviceFleet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_device_fleet",
		Description: "AWS SageMaker Device Fleet",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("device_fleet_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "NotFoundException", "ResourceNotFound"}),
			},
			Hydrate: getAwsSageMakerDeviceFleet,
			Tags:    map[string]string{"service": "sagemaker", "action": "DescribeDeviceFleet"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSageMakerDeviceFleets,
			Tags:    map[string]string{"service": "sagemaker", "action": "ListDeviceFleets"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsSageMakerDeviceFleet,
				Tags: map[string]string{"service": "sagemaker", "action": "DescribeDeviceFleet"},
			},
			{
				Func: listAwsSageMakerDeviceFleetTags,
				Tags: map[string]string{"service": "sagemaker", "action": "ListTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(sagemakerv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "device_fleet_name",
				Description: "The name of the fleet the device belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the fleet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DeviceFleetArn"),
			},
			{
				Name:        "creation_time",
				Description: "Timestamp of when the device fleet was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_modified_time",
				Description: "Timestamp of when the device fleet was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description of the fleet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerDeviceFleet,
			},
			{
				Name:        "iot_role_alias",
				Description: "The Amazon Resource Name (ARN) that has access to AWS Internet of Things (IoT).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerDeviceFleet,
			},
			{
				Name:        "output_config",
				Description: "The output configuration for storing sample data collected by the fleet.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerDeviceFleet,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role that has access to AWS IoT.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerDeviceFleet,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the device fleet.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSageMakerDeviceFleetTags,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DeviceFleetName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSageMakerDeviceFleetTags,
				Transform:   transform.FromValue().Transform(sageMakerTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DeviceFleetArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSageMakerDeviceFleets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_device_fleet.listAwsSageMakerDeviceFleets", "connection_error", err)
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

	input := &sagemaker.ListDeviceFleetsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := sagemaker.NewListDeviceFleetsPaginator(svc, input, func(o *sagemaker.ListDeviceFleetsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if ae.ErrorCode() == "ThrottlingException" {
					// Implement exponential backoff
					waitTime := time.Second * 10
					plugin.Logger(ctx).Warn("aws_sagemaker_device_fleet.listAwsSageMakerDeviceFleets", "throttling_error", err, "wait_time", waitTime)
					time.Sleep(waitTime)
					continue
				}
			}
			plugin.Logger(ctx).Error("aws_sagemaker_device_fleet.listAwsSageMakerDeviceFleets", "api_error", err)
			return nil, err
		}

		for _, items := range output.DeviceFleetSummaries {
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

func getAwsSageMakerDeviceFleet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		name = *h.Item.(types.DeviceFleetSummary).DeviceFleetName
	} else {
		name = d.EqualsQuals["device_fleet_name"].GetStringValue()
	}

	// Create service
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_device_fleet.getAwsSageMakerDeviceFleet", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &sagemaker.DescribeDeviceFleetInput{
		DeviceFleetName: aws.String(name),
	}

	// Get call
	data, err := svc.DescribeDeviceFleet(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_device_fleet.getAwsSageMakerDeviceFleet", "api_error", err)
		return nil, err
	}
	return data, nil
}

func listAwsSageMakerDeviceFleetTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.DeviceFleetSummary).DeviceFleetArn
	}

	// Create Session
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_device_fleet.listAwsSageMakerDeviceFleetTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(arn),
	}

	pagesLeft := true
	tags := []types.Tag{}
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		keyTags, err := svc.ListTags(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sagemaker_device_fleet.listAwsSageMakerDeviceFleetTags", "api_error", err)
			return nil, err
		}
		tags = append(tags, keyTags.Tags...)

		if keyTags.NextToken != nil {
			params.NextToken = keyTags.NextToken
		} else {
			pagesLeft = false
		}
	}

	return tags, nil
}
