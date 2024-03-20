package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"

	guarddutyv1 "github.com/aws/aws-sdk-go/service/guardduty"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type filterInfo = struct {
	guardduty.GetFilterOutput
	Name       string
	DetectorId string
}

func tableAwsGuardDutyFilter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_filter",
		Description: "AWS GuardDuty Filter",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"detector_id", "name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException", "NoSuchEntityException", "BadRequestException"}),
			},
			Hydrate: getAwsGuardDutyFilter,
			Tags:    map[string]string{"service": "guardduty", "action": "GetFilter"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listGuardDutyDetectors,
			Hydrate:       listAwsGuardDutyFilters,
			Tags:          map[string]string{"service": "guardduty", "action": "ListFilters"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "detector_id", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsGuardDutyFilter,
				Tags: map[string]string{"service": "guardduty", "action": "GetFilter"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(guarddutyv1.EndpointsID),
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

//// LIST FUNCTION

func listAwsGuardDutyFilters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := h.Item.(detectorInfo).DetectorID

	// Create session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_filter.listAwsGuardDutyFilters", "get_client_error", err)
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
	params := &guardduty.ListFiltersInput{
		DetectorId: &id,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			params.MaxResults = aws.Int32(limit)
		}
	}

	paginator := guardduty.NewListFiltersPaginator(svc, params, func(o *guardduty.ListFiltersPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_guardduty_filter.listAwsGuardDutyFilters", "api_error", err)
			return nil, err
		}

		for _, item := range output.FilterNames {
			d.StreamListItem(ctx, filterInfo{Name: item, DetectorId: id})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getAwsGuardDutyFilter(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_filter.getAwsGuardDutyFilter", "get_client_error", err)
		return nil, err
	}

	var detectorID string
	var name string

	if h.Item != nil {
		detectorID = h.Item.(filterInfo).DetectorId
		name = h.Item.(filterInfo).Name
	} else {
		detectorID = d.EqualsQuals["detector_id"].GetStringValue()
		name = d.EqualsQuals["name"].GetStringValue()
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
	data, err := svc.GetFilter(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_filter.getAwsGuardDutyFilter", "api_error", err)
		return nil, err
	}

	return filterInfo{*data, name, detectorID}, nil
}

//// TRANSFORM FUNCTION

func getAwsGuardDutyFilterAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(filterInfo)
	region := d.EqualsQualString(matrixKeyRegion)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_filter.getAwsGuardDutyFilterAkas", "api_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := fmt.Sprintf("arn:%s:guardduty:%s:%s:detector/%s/filter/%s", commonColumnData.Partition, region, commonColumnData.AccountId, data.DetectorId, data.Name)

	return []string{aka}, nil
}
