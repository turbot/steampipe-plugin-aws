package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/guardduty"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type detectorInfo = struct {
	guardduty.GetDetectorOutput
	DetectorID string
}

//// TABLE DEFINITION

func tableAwsGuardDutyDetector(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_detector",
		Description: "AWS GuardDuty Detector",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("detector_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidInputException", "BadRequestException"}),
			},
			Hydrate: getGuardDutyDetector,
		},
		List: &plugin.ListConfig{
			Hydrate: listGuardDutyDetectors,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "detector_id",
				Description: "The ID of the detector.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DetectorID"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the detector.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGuardDutyDetectorARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "status",
				Description: "The detector status.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGuardDutyDetector,
			},
			{
				Name:        "created_at",
				Description: "The timestamp of when the detector was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getGuardDutyDetector,
			},
			{
				Name:        "finding_publishing_frequency",
				Description: "The publishing frequency of the finding.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGuardDutyDetector,
			},
			{
				Name:        "service_role",
				Description: "The GuardDuty service role.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGuardDutyDetector,
			},
			{
				Name:        "updated_at",
				Description: "The last-updated timestamp for the detector.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getGuardDutyDetector,
			},
			{
				Name:        "data_sources",
				Description: "Describes which data sources are enabled for the detector.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGuardDutyDetector,
			},
			{
				Name:        "master_account",
				Description: "Contains information about the administrator account and invitation.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGuardDutyDetectorMasterAccount,
				Transform:   transform.FromValue(),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DetectorID"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGuardDutyDetector,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGuardDutyDetectorARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listGuardDutyDetectors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_detector.listGuardDutyDetectors", "get_client_error", err)
		return nil, err
	}

	params := &guardduty.ListDetectorsInput{
		MaxResults: int32(50),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(params.MaxResults) {
			params.MaxResults = int32(*limit)
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := svc.ListDetectors(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_guardduty_detector.listGuardDutyDetectors", "api_error", err)
			return nil, err
		}
		for _, item := range response.DetectorIds {
			d.StreamListItem(ctx, detectorInfo{DetectorID: item})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if response.NextToken != nil {
			pagesLeft = true
			params.NextToken = response.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGuardDutyDetector(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		id = h.Item.(detectorInfo).DetectorID
	} else {
		quals := d.KeyColumnQuals
		id = quals["detector_id"].GetStringValue()
	}

	// Create Session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_detector.getGuardDutyDetector", "client_error", err)
		return nil, err
	}

	params := &guardduty.GetDetectorInput{
		DetectorId: &id,
	}

	op, err := svc.GetDetector(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_detector.getGuardDutyDetector", "api_error", err)
		return nil, err
	}

	return detectorInfo{*op, id}, nil
}

func getGuardDutyDetectorMasterAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := h.Item.(detectorInfo).DetectorID

	// Create Session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_detector.listGuardDutyDetectors", "get_client_error", err)
		return nil, err
	}

	params := &guardduty.GetAdministratorAccountInput{
		DetectorId: &id,
	}

	op, err := svc.GetAdministratorAccount(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_detector.getGuardDutyDetectorMasterAccount", "api_error", err)
		return nil, err
	}

	return op.Administrator, nil
}

func getGuardDutyDetectorARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(detectorInfo)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_detector.getGuardDutyDetectorARN", "error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":guardduty:" + region + ":" + commonColumnData.AccountId + ":detector/" + data.DetectorID

	return arn, nil
}
