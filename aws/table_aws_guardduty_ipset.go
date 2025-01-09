package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"

	guarddutyEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException", "NoSuchEntityException", "BadRequestException"}),
			},
			Hydrate: getAwsGuardDutyIPSet,
			Tags:    map[string]string{"service": "guardduty", "action": "GetIPSet"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listGuardDutyDetectors,
			Hydrate:       listAwsGuardDutyIPSets,
			Tags:          map[string]string{"service": "guardduty", "action": "ListIPSets"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "detector_id", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsGuardDutyIPSet,
				Tags: map[string]string{"service": "guardduty", "action": "GetIPSet"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(guarddutyEndpoint.GUARDDUTYServiceID),
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
	equalQuals := d.EqualsQuals

	// Minimize the API call with the given detector id
	if equalQuals["detector_id"] != nil {
		if equalQuals["detector_id"].GetStringValue() != id {
			return nil, nil
		}
	}

	// Create session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_ipset.getAwsGuardDutyIPSet", "get_client_error", err)
		return nil, err
	}

	maxItems := int32(50)
	params := &guardduty.ListIPSetsInput{
		DetectorId: &id,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			params.MaxResults = aws.Int32(limit)
		}
	}

	paginator := guardduty.NewListIPSetsPaginator(svc, params, func(o *guardduty.ListIPSetsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_guardduty_ipset.getAwsGuardDutyIPSet", "api_error", err)
			return nil, err
		}

		for _, item := range output.IpSetIds {
			d.StreamListItem(ctx, ipsetInfo{
				IPSetID:    item,
				DetectorID: id,
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
		detectorID = d.EqualsQuals["detector_id"].GetStringValue()
		id = d.EqualsQuals["ipset_id"].GetStringValue()
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
	region := d.EqualsQualString(matrixKeyRegion)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_ipset.getAwsGuardDutyIPSetAkas", "api_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := fmt.Sprintf("arn:%s:guardduty:%s:%s:detector/%s/ipset/%s", commonColumnData.Partition, region, commonColumnData.AccountId, data.DetectorID, data.IPSetID)

	return []string{aka}, nil
}
