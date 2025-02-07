package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected/types"

	wellarchitectedEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedConsolidatedReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_consolidated_report",
		Description: "AWS Well-Architected Consolidated Report",
		List: &plugin.ListConfig{
			Hydrate: listWellArchitectedConsolidatedReports,
			Tags:    map[string]string{"service": "wellarchitected", "action": "GetConsolidatedReport"},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:       "include_shared_resources",
					Require:    plugin.Optional,
					CacheMatch: query_cache.CacheMatchExact,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listWellArchitectedConsolidatedReportBase64,
				Tags: map[string]string{"service": "wellarchitected", "action": "GetConsolidatedReport"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedEndpoint.AWS_WELLARCHITECTED_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "workload_name",
				Description: "The name of the workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workload_arn",
				Description: "The ARN for the workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workload_id",
				Description: "The ID assigned to the workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "include_shared_resources",
				Description: "Set to true to have shared resources included in the report.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromQual("include_shared_resources"),
			},
			{
				Name:        "base64_string",
				Description: "The Base64-encoded string representation of a lens review report. This data can be used to create a PDF file. Only returned by GetConsolidatedReport when PDF format is requested.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listWellArchitectedConsolidatedReportBase64,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "lenses_applied_count",
				Description: "The total number of lenses applied to the workload.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "metric_type",
				Description: "The metric type of a metric in the consolidated report. Currently only WORKLOAD metric types are supported.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The date and time when the consolidated report was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "lenses",
				Description: "The metrics for the lenses in the workload.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "risk_counts",
				Description: "A map from risk names to the count of how many questions have that rating.",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

//// LIST FUNCTION

func listWellArchitectedConsolidatedReports(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_consolidated_report.listWellArchitectedConsolidatedReports", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(15)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &wellarchitected.GetConsolidatedReportInput{
		Format:     types.ReportFormatJson,
		MaxResults: aws.Int32(maxLimit),
	}

	// The default value for IncludeSharedResources in input param is false.
	if d.EqualsQuals["include_shared_resources"] != nil {
		input.IncludeSharedResources = aws.Bool(d.EqualsQuals["include_shared_resources"].GetBoolValue())
	}

	paginator := wellarchitected.NewGetConsolidatedReportPaginator(svc, input, func(o *wellarchitected.GetConsolidatedReportPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wellarchitected_consolidated_report.listWellArchitectedConsolidatedReports", "api_error", err)
			return nil, err
		}

		for _, items := range output.Metrics {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func listWellArchitectedConsolidatedReportBase64(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_consolidated_report.listWellArchitectedConsolidatedReportBase64", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &wellarchitected.GetConsolidatedReportInput{
		Format:     types.ReportFormatPdf,
		MaxResults: aws.Int32(15),
	}

	// The default value for IncludeSharedResources in input param is false.
	if d.EqualsQuals["include_shared_resources"] != nil {
		input.IncludeSharedResources = aws.Bool(d.EqualsQuals["include_shared_resources"].GetBoolValue())
	}

	paginator := wellarchitected.NewGetConsolidatedReportPaginator(svc, input, func(o *wellarchitected.GetConsolidatedReportPaginatorOptions) {
		o.Limit = int32(15)
		o.StopOnDuplicateToken = true
	})

	var pdfFormatbase64Encoded []*string
	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wellarchitected_consolidated_report.listWellArchitectedConsolidatedReportBase64", "api_error", err)
			return nil, err
		}

		if output.Base64String != nil {
			pdfFormatbase64Encoded = append(pdfFormatbase64Encoded, output.Base64String)
		}
	}

	return pdfFormatbase64Encoded, nil
}
