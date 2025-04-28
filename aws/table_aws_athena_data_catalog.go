package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"

	athenav1 "github.com/aws/aws-sdk-go/service/athena"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsAthenaDataCatalog(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_athena_data_catalog",
		Description: "AWS Athena Data Catalog",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getAwsAthenaDataCatalog,
			Tags:       map[string]string{"service": "athena", "action": "GetDataCatalog"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsAthenaDataCatalogs,
			Tags:    map[string]string{"service": "athena", "action": "ListDataCatalogs"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(athenav1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the data catalog.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CatalogName"),
			},
			{
				Name:        "type",
				Description: "The type of data catalog (e.g., GLUE, LAMBDA, HIVE).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description of the data catalog.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parameters",
				Description: "Specifies the Lambda function or functions to use for the data catalog.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CatalogName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAthenaDataCatalogAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsAthenaDataCatalogs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := AthenaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_athena_data_catalog.listAwsAthenaDataCatalogs", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxResults := int32(50)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxResults {
			maxResults = limit
		}
	}

	input := &athena.ListDataCatalogsInput{
		MaxResults: &maxResults,
	}

	paginator := athena.NewListDataCatalogsPaginator(svc, input, func(o *athena.ListDataCatalogsPaginatorOptions) {
		o.Limit = maxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_athena_data_catalog.listAwsAthenaDataCatalogs", "api_error", err)
			return nil, err
		}

		for _, catalog := range output.DataCatalogsSummary {
			d.StreamListItem(ctx, catalog)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsAthenaDataCatalog(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := AthenaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_athena_data_catalog.getAwsAthenaDataCatalog", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var name string
	if h != nil {
		catalog := h.Item.(types.DataCatalogSummary)
		name = *catalog.CatalogName
	} else {
		name = d.EqualsQualString("name")
	}

	// Empty check
	if name == "" {
		return nil, nil
	}

	params := &athena.GetDataCatalogInput{
		Name: aws.String(name),
	}

	op, err := svc.GetDataCatalog(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_athena_data_catalog.getAwsAthenaDataCatalog", "api_error", err)
		return nil, err
	}

	return op.DataCatalog, nil
}

func getAwsAthenaDataCatalogAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	data := h.Item.(types.DataCatalogSummary)

	// Get common columns that will be returned for all resources
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Build ARN
	arn := "arn:" + commonColumnData.Partition + ":athena:" + region + ":" + commonColumnData.AccountId + ":datacatalog/" + *data.CatalogName

	return []string{arn}, nil
}
