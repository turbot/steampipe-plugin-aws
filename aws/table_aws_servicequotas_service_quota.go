package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServiceQuotasServiceQuota(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicequotas_service_quota",
		Description: "AWS ServiceQuotas Service Quota",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"service_code", "quota_code"}),
			ShouldIgnoreError: isNotFoundError([]string{"ValidationException"}),
			Hydrate:           getServiceQuota,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listServiceQuotasServices,
			Hydrate:       listServiceQuotas,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "service_code", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
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
				Name:        "usage_metric",
				Description: "Information about the measurement.",
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
				Hydrate:     getServiceQuotasServiceQuotaArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listServiceQuotas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listServiceQuotas")

	// Create Session
	svc, err := ServiceQuotasService(ctx, d)
	if err != nil {
		return nil, err
	}
	service := h.Item.(*servicequotas.ServiceInfo)
	serviceCode := d.KeyColumnQuals["service_code"].GetStringValue()

	if serviceCode != "" && serviceCode != *service.ServiceCode {
		return nil, nil
	}

	input := &servicequotas.ListServiceQuotasInput{
		MaxResults:  aws.Int64(100),
		ServiceCode: service.ServiceCode,
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

	// check if quotaCode or serviceCode is empty
	if quotaCode == "" || serviceCode == "" {
		return nil, nil
	}

	// Create service
	svc, err := ServiceQuotasService(ctx, d)
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
		plugin.Logger(ctx).Error("DescribeWorkspaces", "ERROR", err)
		return nil, err
	}

	return data.Quota, nil
}

// https://docs.aws.amazon.com/servicequotas/latest/userguide/identity-access-management.html
func getServiceQuotasServiceQuotaArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServiceQuotasServiceQuotaArn")
	region := d.KeyColumnQualString(matrixKeyRegion)
	quota := h.Item.(*servicequotas.ServiceQuota)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":servicequotas:" + region + ":" + commonColumnData.AccountId + ":" + *quota.ServiceCode + "/" + *quota.QuotaCode

	return arn, nil
}
