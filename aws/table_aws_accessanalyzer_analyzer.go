package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/accessanalyzer"
	"github.com/aws/aws-sdk-go-v2/service/accessanalyzer/types"

	accessanalyzerv1 "github.com/aws/aws-sdk-go/service/accessanalyzer"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAccessAnalyzer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_accessanalyzer_analyzer",
		Description: "AWS Access Analyzer",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException", "InvalidParameter"}),
			},
			Hydrate: getAccessAnalyzer,
		},
		List: &plugin.ListConfig{
			Hydrate: listAccessAnalyzers,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "type",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(accessanalyzerv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Analyzer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the analyzer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the analyzer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of analyzer, which corresponds to the zone of trust chosen for the analyzer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "A timestamp for the time at which the analyzer was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_resource_analyzed",
				Description: "The resource that was most recently analyzed by the analyzer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_resource_analyzed_at",
				Description: "The time at which the most recently analyzed resource was analyzed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status_reason",
				Description: "The statusReason provides more details about the current status of the analyzer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "findings",
				Description: "A list of findings retrieved from the analyzer that match the filter criteria specified, if any.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAccessAnalyzerFindings,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAccessAnalyzers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := AccessAnalyzerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_accessanalyzer_analyzer.listAccessAnalyzers", "client_error", err)
		return nil, err
	}

	// The maximum number for MaxResults parameter is not defined by the API
	// We have set the MaxResults to 1000 based on our test
	maxItems := int32(1000)
	input := &accessanalyzer.ListAnalyzersInput{}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["type"] != nil {
		input.Type = types.Type(equalQuals["type"].GetStringValue())
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = &maxItems

	paginator := accessanalyzer.NewListAnalyzersPaginator(svc, input, func(o *accessanalyzer.ListAnalyzersPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_accessanalyzer_analyzer.listAccessAnalyzers", "api_error", err)
			return nil, err
		}

		for _, analyzer := range output.Analyzers {
			d.StreamListItem(ctx, analyzer)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAccessAnalyzer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	if strings.TrimSpace(name) == "" {
		return nil, nil
	}

	// Create Session
	svc, err := AccessAnalyzerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_accessanalyzer_analyzer.getAccessAnalyzer", "client_error", err)
		return nil, err
	}

	// Build the params
	params := &accessanalyzer.GetAnalyzerInput{
		AnalyzerName: &name,
	}

	// Get call
	data, err := svc.GetAnalyzer(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Debug("aws_accessanalyzer_analyzer.getAccessAnalyzer", "api_error", err)
		return nil, err
	}

	return *data.Analyzer, nil
}

func listAccessAnalyzerFindings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(types.AnalyzerSummary)

	// Create Session
	svc, err := AccessAnalyzerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_accessanalyzer_analyzer.listAccessAnalyzerFindings", "client_error", err)
		return nil, err
	}

	var findings []types.FindingSummary
	input := &accessanalyzer.ListFindingsInput{AnalyzerArn: data.Arn}

	paginator := accessanalyzer.NewListFindingsPaginator(svc, input, func(o *accessanalyzer.ListFindingsPaginatorOptions) {
		o.Limit = 1000
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_accessanalyzer_analyzer.listAccessAnalyzerFindings", "api_error", err)
			return nil, err
		}
		findings = append(findings, output.Findings...)
	}

	return findings, nil
}
