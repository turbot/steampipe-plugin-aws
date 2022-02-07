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

func tableAwsServiceQuotasService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicequotas_service",
		Description: "AWS ServiceQuota Service",
		List: &plugin.ListConfig{
			Hydrate: listServiceQuotasServices,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "service_code",
				Description: "The service identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_name",
				Description: "The name of the service.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listServiceQuotasServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listServiceQuotasServices")

	// Create Session
	svc, err := ServiceQuotasService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &servicequotas.ListServicesInput{
		MaxResults: aws.Int64(100),
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
	err = svc.ListServicesPages(
		input,
		func(page *servicequotas.ListServicesOutput, isLast bool) bool {
			for _, service := range page.Services {
				d.StreamListItem(ctx, service)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listServiceQuotasServices", "list", err)
		return nil, err
	}

	return nil, nil
}
