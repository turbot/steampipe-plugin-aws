package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
			KeyColumns:        plugin.SingleColumn("detector_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidInputException", "BadRequestException"}),
			Hydrate:           getGuardDutyDetector,
		},
		List: &plugin.ListConfig{
			Hydrate: listGuardDutyDetectors,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "detector_id",
				Description: "The ID of the detector.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DetectorID"),
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
				Hydrate:     getAwsGuardDutyDetectorAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listGuardDutyDetectors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listGuardDutyDetectors", "AWS_REGION", region)

	// Create session
	svc, err := GuardDutyService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	err = svc.ListDetectorsPages(
		&guardduty.ListDetectorsInput{},
		func(page *guardduty.ListDetectorsOutput, isLast bool) bool {
			for _, result := range page.DetectorIds {
				d.StreamListItem(ctx, detectorInfo{
					DetectorID: *result,
				})
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGuardDutyDetector(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getGuardDutyDetector")

	var id string
	if h.Item != nil {
		id = h.Item.(detectorInfo).DetectorID
	} else {
		quals := d.KeyColumnQuals
		id = quals["detector_id"].GetStringValue()
	}

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("matrixRegionmatrixRegion", "matrixRegion", matrixRegion)

	// Create Session
	svc, err := GuardDutyService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &guardduty.GetDetectorInput{
		DetectorId: &id,
	}

	op, err := svc.GetDetector(params)
	if err != nil {
		logger.Debug("getGuardDutyDetector", "ERROR", err)
		return nil, err
	}

	return detectorInfo{*op, id}, nil
}

func getAwsGuardDutyDetectorAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsGuardDutyDetectorAkas")
	data := h.Item.(detectorInfo)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := []string{"arn:" + commonColumnData.Partition + ":guardduty:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":detector/" + data.DetectorID}

	return aka, nil
}
