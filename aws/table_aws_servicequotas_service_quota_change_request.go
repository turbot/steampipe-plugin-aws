package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServiceQuotasServiceQuotaChangeRequest(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicequotas_service_quota_change_request",
		Description: "AWS ServiceQuotas Service Quota Change Request",
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchResourceException"}),
		},
		Get: &plugin.GetConfig{
			Hydrate:    getServiceQuotaChangeRequest,
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchResourceException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceQuotaChangeRequests,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchResourceException"}),
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

	// Create Session
	svc, err := ServiceQuotasClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_service_quota_change_request.listServiceQuotaChangeRequests", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(100)
	input := &servicequotas.ListRequestedServiceQuotaChangeHistoryInput{}
	if d.KeyColumnQuals["service_code"] != nil {
		input.ServiceCode = aws.String(d.KeyColumnQuals["service_code"].GetStringValue())
	}
	if d.KeyColumnQuals["status"] != nil {
		input.Status = types.RequestStatus(d.KeyColumnQuals["status"].GetStringValue())
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = 1
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := servicequotas.NewListRequestedServiceQuotaChangeHistoryPaginator(svc, input, func(o *servicequotas.ListRequestedServiceQuotaChangeHistoryPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_servicequotas_service_quota_change_request.listServiceQuotaChangeRequests", "api_error", err)
			return nil, err
		}

		for _, requestedQuota := range output.RequestedQuotas {
			d.StreamListItem(ctx, requestedQuota)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getServiceQuotaChangeRequest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	id := d.KeyColumnQuals["id"].GetStringValue()

	// check if id is empty
	if id == "" {
		return nil, nil
	}

	// Create service
	svc, err := ServiceQuotasClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_service_quota_change_request.getServiceQuotaChangeRequest", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &servicequotas.GetRequestedServiceQuotaChangeInput{
		RequestId: aws.String(id),
	}

	// Get call
	data, err := svc.GetRequestedServiceQuotaChange(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("getServiceQuotaChangeRequest", "get", err)
		return nil, err
	}

	return *data.RequestedQuota, nil
}

func getServiceQuotaChangeRequestTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	quota := h.Item.(types.RequestedServiceQuotaChange)

	// Create service
	svc, err := ServiceQuotasClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_service_quota_change_request.getServiceQuotaChangeRequestTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &servicequotas.ListTagsForResourceInput{
		ResourceARN: quota.QuotaArn,
	}

	data, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_service_quota_change_request.getServiceQuotaChangeRequestTags", "api_error", err)
		return nil, err
	}

	return data.Tags, nil
}

func getServiceQuotaChangeRequestAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(types.RequestedServiceQuotaChange)

	// Get common columns
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_service_quota_change_request.getServiceQuotaChangeRequestAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := fmt.Sprintf("arn:/%s:servicequotas:%s:%s:changeRequest/%s", commonColumnData.Partition, region, commonColumnData.AccountId, *data.Id)

	return []string{arn}, nil
}

//// TRANSFORM FUNCTIONS

func serviceQuotaChangeRequestTagsToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if len(tags) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
