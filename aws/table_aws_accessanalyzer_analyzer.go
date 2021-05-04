package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/accessanalyzer"

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
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAccessAnalyzers", "AWS_REGION", region)

	// Create session
	svc, err := AccessAnalyzerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListAnalyzersPages(
		&accessanalyzer.ListAnalyzersInput{},
		func(page *accessanalyzer.ListAnalyzersOutput, isLast bool) bool {
			for _, analyzer := range page.Analyzers {
				d.StreamListItem(ctx, analyzer)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAccessAnalyzer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAccessAnalyzer")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Session
	svc, err := AccessAnalyzerService(ctx, d, region)
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

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	data := h.Item.(*accessanalyzer.AnalyzerSummary)

	// Create Session
	svc, err := AccessAnalyzerService(ctx, d, region)
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
