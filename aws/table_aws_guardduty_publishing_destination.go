package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
		MaxResults: int32(50),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(input.MaxResults) {
			input.MaxResults = int32(*limit)
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := svc.ListPublishingDestinations(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_guardduty_publishing_destination.listGuardDutyPublishingDestinations", "api_error", err)
			return nil, err
		}
		for _, item := range response.Destinations {
			d.StreamListItem(ctx, DestinationInfo{
				DestinationId:   item.DestinationId,
				DestinationType: item.DestinationType,
				Status:          item.Status,
				DetectorId:      id,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if response.NextToken != nil {
			pagesLeft = true
			input.NextToken = response.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, err
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
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_publishing_destination.getPublishingDestinationArn", "api_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":guardduty:" + region + ":" + commonColumnData.AccountId + ":detector" + "/" + data.DetectorId + "/publishingDestination" + "/" + *data.DestinationId

	return aka, nil
}
