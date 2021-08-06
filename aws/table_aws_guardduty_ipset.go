package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/guardduty"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsGuardDutyIPSet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_ipset",
		Description: "AWS GuardDuty IPSet",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"detector_id", "ipset_id"}),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidInputException", "NoSuchEntityException", "BadRequestException"}),
			Hydrate:           getAwsGuardDutyIPSet,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listGuardDutyDetectors,
			Hydrate:       listAwsGuardDutyIPSets,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name for the IPSet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsGuardDutyIPSet,
			},
			{
				Name:        "detector_id",
				Description: "The ID of the detector.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsGuardDutyIPSet,
				Transform:   transform.FromField("DetectorID"),
			},
			{
				Name:        "ipset_id",
				Description: "The ID of the IPSet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IPSetID"),
			},
			{
				Name:        "format",
				Description: "The format of the file that contains the IPSet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsGuardDutyIPSet,
			},
			{
				Name:        "status",
				Description: "The status of IPSet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsGuardDutyIPSet,
			},
			{
				Name:        "location",
				Description: "The URI of the file that contains the IPSet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsGuardDutyIPSet,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsGuardDutyIPSet,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsGuardDutyIPSet,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsGuardDutyIPSetAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type ipsetInfo = struct {
	guardduty.GetIPSetOutput
	IPSetID    string
	DetectorID string
}

//// LIST FUNCTION

func listAwsGuardDutyIPSets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := h.Item.(detectorInfo).DetectorID

	// Create session
	svc, err := GuardDutyService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListIPSetsPages(
		&guardduty.ListIPSetsInput{DetectorId: &id},
		func(page *guardduty.ListIPSetsOutput, isLast bool) bool {
			for _, parameter := range page.IpSetIds {
				d.StreamLeafListItem(ctx, ipsetInfo{
					IPSetID:    *parameter,
					DetectorID: id,
				})

			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTION

func getAwsGuardDutyIPSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsGuardDutyIPSet")

	// Create Session
	svc, err := GuardDutyService(ctx, d)
	if err != nil {
		return nil, err
	}
	var detectorID string
	var id string
	if h.Item != nil {
		detectorID = h.Item.(ipsetInfo).DetectorID
		id = h.Item.(ipsetInfo).IPSetID
	} else {
		detectorID = d.KeyColumnQuals["detector_id"].GetStringValue()
		id = d.KeyColumnQuals["ipset_id"].GetStringValue()
	}

	// Build the params
	params := &guardduty.GetIPSetInput{
		DetectorId: &detectorID,
		IpSetId:    &id,
	}

	// Get call
	data, err := svc.GetIPSet(params)
	if err != nil {
		return nil, err
	}

	return ipsetInfo{*data, id, detectorID}, nil
}

//// TRANSFORM FUNCTION

func getAwsGuardDutyIPSetAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsGuardDutyIPSetAkas")

	data := h.Item.(ipsetInfo)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":guardduty:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":detector" + "/" + data.DetectorID + "/ipset" + "/" + data.IPSetID

	return []string{aka}, nil
}
