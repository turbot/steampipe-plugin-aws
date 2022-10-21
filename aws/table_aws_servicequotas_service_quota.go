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

func tableAwsServiceQuotasServiceQuota(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicequotas_service_quota",
		Description: "AWS ServiceQuotas Service Quota",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"service_code", "quota_code", "region"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchResourceException"}),
			},
			Hydrate: getServiceQuota,
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceQuotas,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchResourceException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "service_code", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildServiceQuotasServicesRegionList,
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
	plugin.Logger(ctx).Trace("listServiceQuotas")

	// Create Session
	svc, err := ServiceQuotasRegionalService(ctx, d)
	if err != nil {
		return nil, err
	}

	matrixServiceCode := d.KeyColumnQualString(matrixKeyServiceCode)
	serviceCode := d.KeyColumnQuals["service_code"].GetStringValue()

	// Filter the serviceCode if user provided value for it
	if serviceCode != "" && serviceCode != matrixServiceCode {
		return nil, nil
	}

	input := &servicequotas.ListServiceQuotasInput{
		MaxResults:  aws.Int64(100),
		ServiceCode: aws.String(matrixServiceCode),
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
	err = svc.ListServiceQuotasPages(
		input,
		func(page *servicequotas.ListServiceQuotasOutput, isLast bool) bool {
			for _, quota := range page.Quotas {
				d.StreamListItem(ctx, quota)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listServiceQuotas", "list", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getServiceQuota(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServiceQuota")

	quotaCode := d.KeyColumnQuals["quota_code"].GetStringValue()
	serviceCode := d.KeyColumnQuals["service_code"].GetStringValue()
	region := d.KeyColumnQuals["region"].GetStringValue()

	// check if quotaCode or serviceCode or region is empty
	if quotaCode == "" || serviceCode == "" || region == "" {
		return nil, nil
	}

	// Filter the serviceCode and region with the provided value
	matrixServiceCode := d.KeyColumnQualString(matrixKeyServiceCode)
	matrixRegion := d.KeyColumnQualString(matrixKeyRegion)
	if serviceCode != matrixServiceCode || region != matrixRegion {
		return nil, nil
	}

	// Create service
	svc, err := ServiceQuotasRegionalService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String(quotaCode),
		ServiceCode: aws.String(serviceCode),
	}

	// Get call
	data, err := svc.GetServiceQuota(params)
	if err != nil {
		plugin.Logger(ctx).Error("getServiceQuota", "get", err)
		return nil, err
	}

	return data.Quota, nil
}

func getServiceQuotaTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServiceQuotaTags")

	quota := h.Item.(*servicequotas.ServiceQuota)

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
		plugin.Logger(ctx).Error("getServiceQuotaTags", "error", err)
		return nil, err
	}

	return data.Tags, nil
}

//// TRANSFORM FUNCTIONS

func serviceQuotaTagsToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("serviceQuotaTagsToTurbotTags")
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
