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

func tableAwsGuardDutyFilter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_filter",
		Description: "AWS GuardDuty Filter",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"detector_id", "name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidInputException", "NoSuchEntityException", "BadRequestException"}),
			},
			Hydrate: getAwsGuardDutyFilter,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listGuardDutyDetectors,
			Hydrate:       listAwsGuardDutyFilters,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "detector_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name for the filter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "detector_id",
				Description: "The ID of the detector.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsGuardDutyFilter,
			},
			{
				Name:        "action",
				Description: "Specifies the action that is to be applied to the findings that match the filter.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsGuardDutyFilter,
			},
			{
				Name:        "description",
				Description: "The description of the filter.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsGuardDutyFilter,
			},
			{
				Name:        "rank",
				Description: "Specifies the position of the filter in the list of current filters. Also specifies the order in which this filter is applied to the findings.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsGuardDutyFilter,
			},
			{
				Name:        "finding_criteria",
				Description: "Represents the criteria to be used in the filter for querying findings.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsGuardDutyFilter,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsGuardDutyFilter,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsGuardDutyFilter,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsGuardDutyFilterAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type filterInfo = struct {
	guardduty.GetFilterOutput
	Name       string
	DetectorId string
}

//// LIST FUNCTION

func listAwsGuardDutyFilters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	input := &guardduty.ListFiltersInput{
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
	err = svc.ListFiltersPages(
		input,
		func(page *guardduty.ListFiltersOutput, isLast bool) bool {
			for _, parameter := range page.FilterNames {
				d.StreamLeafListItem(ctx, filterInfo{
					Name:       *parameter,
					DetectorId: id,
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

func getAwsGuardDutyFilter(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsGuardDutyFilter")

	// Create Session
	svc, err := GuardDutyService(ctx, d)
	if err != nil {
		return nil, err
	}
	var detectorID string
	var name string
	if h.Item != nil {
		detectorID = h.Item.(filterInfo).DetectorId
		name = h.Item.(filterInfo).Name
	} else {
		detectorID = d.KeyColumnQuals["detector_id"].GetStringValue()
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// check if name or detectorID is empty
	if detectorID == "" || name == "" {
		return nil, nil
	}

	// Build the params
	params := &guardduty.GetFilterInput{
		DetectorId: &detectorID,
		FilterName: &name,
	}

	// Get call
	data, err := svc.GetFilter(params)
	if err != nil {
		return nil, err
	}

	return filterInfo{*data, name, detectorID}, nil
}

//// TRANSFORM FUNCTION

func getAwsGuardDutyFilterAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsGuardDutyFilterAkas")

	data := h.Item.(filterInfo)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":guardduty:" + region + ":" + commonColumnData.AccountId + ":detector" + "/" + data.DetectorId + "/filter" + "/" + data.Name

	return []string{aka}, nil
}
