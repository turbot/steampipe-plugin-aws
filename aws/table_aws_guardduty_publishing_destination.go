package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"

	guarddutyv1 "github.com/aws/aws-sdk-go/service/guardduty"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type DestinationInfo = struct {
	DestinationId                   *string
	DestinationType                 types.DestinationType
	DestinationArn                  *string
	KmsKeyArn                       *string
	Status                          types.PublishingStatus
	PublishingFailureStartTimestamp *int64
	DetectorId                      string
}

func tableAwsGuardDutyPublishingDestination(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_publishing_destination",
		Description: "AWS GuardDuty Publishing Destination",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"detector_id", "destination_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException", "NoSuchEntityException", "BadRequestException"}),
			},
			Hydrate: getGuardDutyPublishingDestination,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listGuardDutyDetectors,
			Hydrate:       listGuardDutyPublishingDestinations,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "detector_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(guarddutyv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "destination_id",
				Description: "The ID of the publishing destination.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the publishing destination.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getPublishingDestinationArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "detector_id",
				Description: "The ID of the detector.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "destination_arn",
				Description: "The ARN of the resource to publish to.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGuardDutyPublishingDestination,
			},
			{
				Name:        "destination_type",
				Description: "The type of publishing destination. Currently, only Amazon S3 buckets are supported.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_arn",
				Description: "The ARN of the KMS key to use for encryption.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGuardDutyPublishingDestination,
			},
			{
				Name:        "publishing_failure_start_timestamp",
				Description: "The time, in epoch millisecond format, at which GuardDuty was first unable to publish findings to the destination.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getGuardDutyPublishingDestination,
				Transform:   transform.FromField("PublishingFailureStartTimestamp").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "status",
				Description: "The status of the publishing destination.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DestinationId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPublishingDestinationArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listGuardDutyPublishingDestinations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := h.Item.(detectorInfo).DetectorID

	// Create session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_publishing_destination.listGuardDutyPublishingDestinations", "get_client_error", err)
		return nil, err
	}

	equalQuals := d.EqualsQuals
	// Minimize the API call with the given detector_id
	if equalQuals["detector_id"] != nil {
		if equalQuals["detector_id"].GetStringValue() != id {
			return nil, nil
		}
	}

	maxItems := int32(50)
	input := &guardduty.ListPublishingDestinationsInput{
		DetectorId: &id,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			input.MaxResults = limit
		}
	}

	paginator := guardduty.NewListPublishingDestinationsPaginator(svc, input, func(o *guardduty.ListPublishingDestinationsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_guardduty_publishing_destination.listGuardDutyPublishingDestinations", "api_error", err)
			return nil, err
		}

		for _, item := range output.Destinations {
			d.StreamListItem(ctx, DestinationInfo{
				DestinationId:   item.DestinationId,
				DestinationType: item.DestinationType,
				Status:          item.Status,
				DetectorId:      id,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getGuardDutyPublishingDestination(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_publishing_destination.getGuardDutyPublishingDestination", "get_client_error", err)
		return nil, err
	}

	var detectorID string
	var id string

	if h.Item != nil {
		detectorID = h.Item.(DestinationInfo).DetectorId
		id = *h.Item.(DestinationInfo).DestinationId
	} else {
		detectorID = d.EqualsQuals["detector_id"].GetStringValue()
		id = d.EqualsQuals["destination_id"].GetStringValue()
	}

	// Empty check
	if detectorID == "" || id == "" {
		return nil, nil
	}

	// Build the params
	params := &guardduty.DescribePublishingDestinationInput{
		DetectorId:    &detectorID,
		DestinationId: &id,
	}

	// Get call
	data, err := svc.DescribePublishingDestination(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_publishing_destination.getPublishingDestinationArn", "api_error", err)
		return nil, err
	}

	return DestinationInfo{
		DestinationId:                   &id,
		DestinationType:                 data.DestinationType,
		DestinationArn:                  data.DestinationProperties.DestinationArn,
		KmsKeyArn:                       data.DestinationProperties.KmsKeyArn,
		Status:                          data.Status,
		PublishingFailureStartTimestamp: aws.Int64(data.PublishingFailureStartTimestamp),
		DetectorId:                      detectorID,
	}, nil
}

func getPublishingDestinationArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(DestinationInfo)
	region := d.EqualsQualString(matrixKeyRegion)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_publishing_destination.getPublishingDestinationArn", "api_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := fmt.Sprintf("arn:%s:guardduty:%s:%s:detector/%s/publishingDestination/%s", commonColumnData.Partition, region, commonColumnData.AccountId, data.DetectorId, *data.DestinationId)

	return aka, nil
}
