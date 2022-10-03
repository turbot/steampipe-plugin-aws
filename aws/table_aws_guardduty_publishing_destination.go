package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/guardduty"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsGuardDutyPublishingDestination(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_publishing_destination",
		Description: "AWS GuardDuty Publishing Destination",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"detector_id", "destination_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidInputException", "NoSuchEntityException", "BadRequestException"}),
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
		GetMatrixItemFunc: BuildRegionList,
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

type DestinationInfo = struct {
	DestinationId                   *string
	DestinationType                 *string
	DestinationArn                  *string
	KmsKeyArn                       *string
	Status                          *string
	PublishingFailureStartTimestamp *int64
	DetectorId                      string
}

//// LIST FUNCTION

func listGuardDutyPublishingDestinations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := h.Item.(detectorInfo).DetectorID

	// Create session
	svc, err := GuardDutyService(ctx, d)
	if err != nil {
		return nil, err
	}

	equalQuals := d.KeyColumnQuals

	// Minimize the API call with the given detector_id
	if equalQuals["detector_id"] != nil {
		if equalQuals["detector_id"].GetStringValue() != "" {
			if equalQuals["detector_id"].GetStringValue() != "" && equalQuals["detector_id"].GetStringValue() != id {
				return nil, nil
			}
		} else if len(getListValues(equalQuals["detector_id"].GetListValue())) > 0 {
			if !strings.Contains(fmt.Sprint(getListValues(equalQuals["detector_id"].GetListValue())), id) {
				return nil, nil
			}
		}
	}

	input := &guardduty.ListPublishingDestinationsInput{
		DetectorId: &id,
		MaxResults: aws.Int64(50),
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
	err = svc.ListPublishingDestinationsPages(
		input,
		func(page *guardduty.ListPublishingDestinationsOutput, isLast bool) bool {
			for _, destination := range page.Destinations {
				d.StreamLeafListItem(ctx, DestinationInfo{
					DestinationId:   destination.DestinationId,
					DestinationType: destination.DestinationType,
					Status:          destination.Status,
					DetectorId:      id,
				})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTION

func getGuardDutyPublishingDestination(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getGuardDutyPublishingDestination")

	// Create Session
	svc, err := GuardDutyService(ctx, d)
	if err != nil {
		return nil, err
	}
	var detectorID string
	var id string
	if h.Item != nil {
		detectorID = h.Item.(DestinationInfo).DetectorId
		id = *h.Item.(DestinationInfo).DestinationId
	} else {
		detectorID = d.KeyColumnQuals["detector_id"].GetStringValue()
		id = d.KeyColumnQuals["destination_id"].GetStringValue()
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
	data, err := svc.DescribePublishingDestination(params)
	if err != nil {
		logger.Error("getGuardDutyPublishingDestination", "err", err)
		return nil, err
	}

	return DestinationInfo{
		DestinationId:                   &id,
		DestinationType:                 data.DestinationType,
		DestinationArn:                  data.DestinationProperties.DestinationArn,
		KmsKeyArn:                       data.DestinationProperties.KmsKeyArn,
		Status:                          data.Status,
		PublishingFailureStartTimestamp: data.PublishingFailureStartTimestamp,
		DetectorId:                      detectorID,
	}, nil
}

func getPublishingDestinationArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(DestinationInfo)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":guardduty:" + region + ":" + commonColumnData.AccountId + ":detector" + "/" + data.DetectorId + "/publishingDestination" + "/" + *data.DestinationId

	return aka, nil
}
