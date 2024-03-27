package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas/types"

	servicequotasv1 "github.com/aws/aws-sdk-go/service/servicequotas"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServiceQuotasServiceQuota(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicequotas_service_quota",
		Description: "AWS Service Quotas Service Quota",
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchResourceException"}),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"service_code", "quota_code"}),
			Hydrate:    getServiceQuota,
			Tags:       map[string]string{"service": "servicequotas", "action": "GetServiceQuota"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchResourceException"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listServiceQuotasServices,
			Hydrate:       listServiceQuotas,
			Tags:          map[string]string{"service": "servicequotas", "action": "ListServiceQuotas"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchResourceException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "service_code", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getServiceQuotaTags,
				Tags: map[string]string{"service": "servicequotas", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(servicequotasv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
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
				Name:        "global_quota",
				Description: "Indicates whether the quota is global.",
				Type:        proto.ColumnType_BOOL,
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
				Name:        "adjustable",
				Description: "Indicates whether the quota value can be increased.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "unit",
				Description: "The unit of measurement.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "The quota value.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "error_reason",
				Description: "The error code and error reason.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "period",
				Description: "The period of time for the quota.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags associated with the service quota.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceQuotaTags,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "usage_metric",
				Description: "Information about the measurement.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceQuotaTags,
				Transform:   transform.From(serviceQuotaTagsToTurbotTags),
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
				Transform:   transform.FromField("QuotaArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listServiceQuotas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service := h.Item.(types.ServiceInfo)

	// Create Session
	svc, err := ServiceQuotasClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_service_quota.listServiceQuotas", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	serviceCode := d.EqualsQuals["service_code"].GetStringValue()
	// Filter the serviceCode if user provided value for it
	if serviceCode != "" && serviceCode != *service.ServiceCode {
		return nil, nil
	}

	maxItems := int32(100)
	input := &servicequotas.ListServiceQuotasInput{
		ServiceCode: service.ServiceCode,
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
	paginator := servicequotas.NewListServiceQuotasPaginator(svc, input, func(o *servicequotas.ListServiceQuotasPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_servicequotas_service_quota.listServiceQuotas", "api_error", err)
			return nil, err
		}

		for _, quota := range output.Quotas {
			d.StreamListItem(ctx, quota)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getServiceQuota(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	quotaCode := d.EqualsQuals["quota_code"].GetStringValue()
	serviceCode := d.EqualsQuals["service_code"].GetStringValue()

	// check if quotaCode or serviceCode is empty
	if quotaCode == "" || serviceCode == "" {
		return nil, nil
	}

	// Create service
	svc, err := ServiceQuotasClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_service_quota.getServiceQuota", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String(quotaCode),
		ServiceCode: aws.String(serviceCode),
	}

	// Get call
	data, err := svc.GetServiceQuota(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_service_quota.getServiceQuota", "api_error", err)
		return nil, err
	}

	return *data.Quota, nil
}

func getServiceQuotaTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quota := h.Item.(types.ServiceQuota)

	// Create service
	svc, err := ServiceQuotasClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_service_quota.getServiceQuotaTags", "connection_error", err)
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
		plugin.Logger(ctx).Error("aws_servicequotas_service_quota.getServiceQuotaTags", "connection_error", err)
		return nil, err
	}

	return data.Tags, nil
}

//// TRANSFORM FUNCTIONS

func serviceQuotaTagsToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
