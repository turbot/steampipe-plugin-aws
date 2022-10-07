package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/guardduty"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type ipsetInfo = struct {
	guardduty.GetIPSetOutput
	IPSetID    string
	DetectorID string
}

func tableAwsGuardDutyIPSet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_ipset",
		Description: "AWS GuardDuty IPSet",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"detector_id", "ipset_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidInputException", "NoSuchEntityException", "BadRequestException"}),
			},
			Hydrate: getAwsGuardDutyIPSet,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listGuardDutyDetectors,
			Hydrate:       listAwsGuardDutyIPSets,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "detector_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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

//// LIST FUNCTION

func listAwsGuardDutyIPSets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := h.Item.(detectorInfo).DetectorID

	// Create session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_ipset.getAwsGuardDutyIPSet", "get_client_error", err)
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

	input := &guardduty.ListIPSetsInput{
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
		response, err := svc.ListIPSets(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_guardduty_detector.listGuardDutyDetectors", "api_error", err)
			return nil, err
		}
		for _, item := range response.IpSetIds {
			d.StreamListItem(ctx, ipsetInfo{
				IPSetID:    item,
				DetectorID: id,
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

func getAwsGuardDutyIPSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_ipset.getAwsGuardDutyIPSet", "get_client_error", err)
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
	data, err := svc.GetIPSet(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_ipset.getAwsGuardDutyIPSet", "api_error", err)
		return nil, err
	}

	return ipsetInfo{*data, id, detectorID}, nil
}

//// TRANSFORM FUNCTION

func getAwsGuardDutyIPSetAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(ipsetInfo)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_ipset.getAwsGuardDutyIPSetAkas", "api_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":guardduty:" + region + ":" + commonColumnData.AccountId + ":detector" + "/" + data.DetectorID + "/ipset" + "/" + data.IPSetID

	return []string{aka}, nil
}
