package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas/types"

	servicequotasEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsServiceQuotasService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicequotas_service",
		Description: "AWS Service Quotas Service",
		List: &plugin.ListConfig{
			Hydrate: listServiceQuotasServices,
			Tags:    map[string]string{"service": "servicequotas", "action": "ListServices"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(servicequotasEndpoint.AWS_SERVICEQUOTAS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "service_name",
				Description: "Specifies the service name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_code",
				Description: "Specifies the service identifier.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceQuotaServiceAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

func listServiceQuotasServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ServiceQuotasClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_service.listServiceQuotasServices", "connection_error", err)
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &servicequotas.ListServicesInput{
		MaxResults: aws.Int32(100),
	}

	paginator := servicequotas.NewListServicesPaginator(svc, input, func(o *servicequotas.ListServicesPaginatorOptions) {
		o.Limit = 100
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_servicequotas_service.listServiceQuotasServices", "api_error", err)
			return nil, err
		}

		for _, service := range output.Services {
			d.StreamListItem(ctx, service)
		}
	}

	return nil, nil
}

func getServiceQuotaServiceAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	data := h.Item.(types.ServiceInfo)

	// Get common columns

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_service.getServiceQuotaServiceAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := fmt.Sprintf("arn:%s:servicequotas:%s:%s:%s", commonColumnData.Partition, region, commonColumnData.AccountId, *data.ServiceCode)

	return []string{arn}, nil
}
