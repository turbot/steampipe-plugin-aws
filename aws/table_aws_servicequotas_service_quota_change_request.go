package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServiceQuotasServiceQuotaChangeRequest(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicequotas_service_quota_change_request",
		Description: "AWS ServiceQuotas Service Quota Change Request",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchResourceException"}),
			},
			Hydrate: getServiceQuotaChangeRequest,
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceQuotaChangeRequests,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchResourceException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "service_code", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "case_id",
				Description: "The case ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The state of the quota increase request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "quota_name",
				Description: "The quota name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "quota_code",
				Description: "The quota code.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "quota_arn",
				Description: "The arn of the service quota.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "desired_value",
				Description: "The increased value for the quota.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "created",
				Description: "The date and time when the quota increase request was received and the case ID was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "global_quota",
				Description: "Indicates whether the quota is global.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_updated",
				Description: "The date and time of the most recent change.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "requester",
				Description: "The IAM identity of the requester.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_name",
				Description: "The service name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_code",
				Description: "The service identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "unit",
				Description: "The unit of measurement.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags associated with the change request.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceQuotaChangeRequestTags,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceQuotaChangeRequestTags,
				Transform:   transform.From(serviceQuotaChangeRequestTagsToTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("QuotaName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceQuotaChangeRequestAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listServiceQuotaChangeRequests(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listServiceQuotaChangeRequests")

	// Create Session
	svc, err := ServiceQuotasRegionalService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &servicequotas.ListRequestedServiceQuotaChangeHistoryInput{
		MaxResults: aws.Int64(100),
	}

	if d.KeyColumnQuals["service_code"] != nil {
		input.ServiceCode = aws.String(d.KeyColumnQuals["service_code"].GetStringValue())
	}
	if d.KeyColumnQuals["status"] != nil {
		input.Status = aws.String(d.KeyColumnQuals["status"].GetStringValue())
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
	err = svc.ListRequestedServiceQuotaChangeHistoryPages(
		input,
		func(page *servicequotas.ListRequestedServiceQuotaChangeHistoryOutput, isLast bool) bool {
			for _, requestedQuota := range page.RequestedQuotas {
				d.StreamListItem(ctx, requestedQuota)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listServiceQuotaChangeRequests", "list", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getServiceQuotaChangeRequest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServiceQuotaChangeRequest")

	id := d.KeyColumnQuals["id"].GetStringValue()

	// check if id is empty
	if id == "" {
		return nil, nil
	}

	// Create service
	svc, err := ServiceQuotasRegionalService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &servicequotas.GetRequestedServiceQuotaChangeInput{
		RequestId: aws.String(id),
	}

	// Get call
	data, err := svc.GetRequestedServiceQuotaChange(params)
	if err != nil {
		plugin.Logger(ctx).Error("getServiceQuotaChangeRequest", "get", err)
		return nil, err
	}

	return data.RequestedQuota, nil
}

func getServiceQuotaChangeRequestTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServiceQuotaChangeRequestTags")

	quota := h.Item.(*servicequotas.RequestedServiceQuotaChange)

	// Create service
	svc, err := ServiceQuotasRegionalService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &servicequotas.ListTagsForResourceInput{
		ResourceARN: quota.QuotaArn,
	}

	data, err := svc.ListTagsForResource(params)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchResourceException") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("getServiceQuotaChangeRequestTags", "error", err)
		return nil, err
	}

	return data.Tags, nil
}

func getServiceQuotaChangeRequestAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServiceQuotaChangeRequestAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(*servicequotas.RequestedServiceQuotaChange)

	// Get common columns
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":servicequotas:" + region + ":" + commonColumnData.AccountId + ":changeRequest/" + *data.Id}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func serviceQuotaChangeRequestTagsToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("serviceQuotaChangeRequestTagsToTurbotTags")
	tags := d.HydrateItem.([]*servicequotas.Tag)

	if tags == nil {
		return nil, nil
	}
	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
