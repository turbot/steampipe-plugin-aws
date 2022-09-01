package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServiceQuotasDefaultServiceQuota(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicequotas_default_service_quota",
		Description: "AWS ServiceQuotas Default Service Quota",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"service_code", "quota_code", "region"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchResourceException"}),
			},
			Hydrate: getDefaultServiceQuota,
		},
		List: &plugin.ListConfig{
			Hydrate: listDefaultServiceQuotas,
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
				Transform:   transform.FromField("QuotaArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDefaultServiceQuotas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDefaultServiceQuotas")

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

	input := &servicequotas.ListAWSDefaultServiceQuotasInput{
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
	err = svc.ListAWSDefaultServiceQuotasPages(
		input,
		func(page *servicequotas.ListAWSDefaultServiceQuotasOutput, isLast bool) bool {
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
		plugin.Logger(ctx).Error("listDefaultServiceQuotas", "list", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDefaultServiceQuota(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDefaultServiceQuota")

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
	params := &servicequotas.GetAWSDefaultServiceQuotaInput{
		QuotaCode:   aws.String(quotaCode),
		ServiceCode: aws.String(serviceCode),
	}

	// Get call
	data, err := svc.GetAWSDefaultServiceQuota(params)
	if err != nil {
		plugin.Logger(ctx).Error("getDefaultServiceQuota", "get", err)
		return nil, err
	}

	return data.Quota, nil
}
