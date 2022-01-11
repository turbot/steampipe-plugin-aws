package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/accessanalyzer"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAccessAnalyzer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_accessanalyzer_analyzer",
		Description: "AWS Access Analyzer",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException", "ValidationException", "InvalidParameter"}),
			Hydrate:           getAccessAnalyzer,
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
		GetMatrixItem: BuildRegionList,
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
	logger := plugin.Logger(ctx)

	// Create session
	svc, err := AccessAnalyzerService(ctx, d)
	if err != nil {
		logger.Trace("listAccessAnalyzers", "connection error", err)
		return nil, err
	}

	// The maximum number for MaxResults parameter is not defined by the API
	// We have set the MaxResults to 1000 based on our test
	input := &accessanalyzer.ListAnalyzersInput{
		MaxResults: aws.Int64(1000),
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["type"] != nil {
		input.Type = types.String(equalQuals["type"].GetStringValue())
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = types.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.ListAnalyzersPages(
		input,
		func(page *accessanalyzer.ListAnalyzersOutput, isLast bool) bool {
			for _, analyzer := range page.Analyzers {
				d.StreamListItem(ctx, analyzer)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAccessAnalyzer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAccessAnalyzer")

	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Session
	svc, err := AccessAnalyzerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &accessanalyzer.GetAnalyzerInput{
		AnalyzerName: &name,
	}

	// Get call
	data, err := svc.GetAnalyzer(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAccessAnalyzer", "ERROR", err)
		return nil, err
	}

	return data.Analyzer, nil
}

func listAccessAnalyzerFindings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAccessAnalyzerFindings")

	data := h.Item.(*accessanalyzer.AnalyzerSummary)

	// Create Session
	svc, err := AccessAnalyzerService(ctx, d)
	if err != nil {
		return nil, err
	}

	var findings []*accessanalyzer.FindingSummary
	err = svc.ListFindingsPages(
		&accessanalyzer.ListFindingsInput{
			AnalyzerArn: data.Arn,
		},
		func(page *accessanalyzer.ListFindingsOutput, isLast bool) bool {
			findings = append(findings, page.Findings...)
			return !isLast
		},
	)
	if err != nil {
		return nil, err
	}

	return findings, nil
}
