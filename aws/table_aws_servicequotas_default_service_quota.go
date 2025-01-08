package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas/types"

	servicequotasEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServiceQuotasDefaultServiceQuota(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicequotas_default_service_quota",
		Description: "AWS Service Quotas Default Service Quota",
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchResourceException"}),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"service_code", "quota_code"}),
			Hydrate:    getDefaultServiceQuota,
			Tags:       map[string]string{"service": "servicequotas", "action": "GetAWSDefaultServiceQuota"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchResourceException"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listServiceQuotasServices,
			Hydrate:       listDefaultServiceQuotas,
			Tags:          map[string]string{"service": "servicequotas", "action": "ListAWSDefaultServiceQuotas"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchResourceException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "service_code", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(servicequotasEndpoint.SERVICEQUOTASServiceID),
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
				Name:        "quota_applied_at_level",
				Description: "Specifies at which level of granularity that the quota value is applied.",
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
				Name:        "usage_metric",
				Description: "Information about the measurement.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "quota_context",
				Description: "The context for this service quota.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
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

func listDefaultServiceQuotas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service := h.Item.(types.ServiceInfo)

	// Create Session
	svc, err := ServiceQuotasClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_default_service_quota.listDefaultServiceQuotas", "connection_error", err)
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
	input := &servicequotas.ListAWSDefaultServiceQuotasInput{
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
	paginator := servicequotas.NewListAWSDefaultServiceQuotasPaginator(svc, input, func(o *servicequotas.ListAWSDefaultServiceQuotasPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_servicequotas_default_service_quota.listDefaultServiceQuotas", "api_error", err)
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

func getDefaultServiceQuota(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	quotaCode := d.EqualsQuals["quota_code"].GetStringValue()
	serviceCode := d.EqualsQuals["service_code"].GetStringValue()

	// check if quotaCode or serviceCode or region is empty
	if quotaCode == "" || serviceCode == "" {
		return nil, nil
	}

	// Create service
	svc, err := ServiceQuotasClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_default_service_quota.getDefaultServiceQuota", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &servicequotas.GetAWSDefaultServiceQuotaInput{
		QuotaCode:   aws.String(quotaCode),
		ServiceCode: aws.String(serviceCode),
	}

	// Get call
	data, err := svc.GetAWSDefaultServiceQuota(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_default_service_quota.getDefaultServiceQuota", "api_error", err)
		return nil, err
	}

	return *data.Quota, nil
}
