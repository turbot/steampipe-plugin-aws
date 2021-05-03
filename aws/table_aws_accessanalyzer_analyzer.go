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
		Description: "AWS IAM Access Analyzer",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException", "ValidationException", "InvalidParameter"}),
			Hydrate:           getAwsAccessAnalyzer,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsAccessAnalyzers,
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

func listAwsAccessAnalyzers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsAccessAnalyzers", "AWS_REGION", region)

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

func getAwsAccessAnalyzer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsAccessAnalyzer")

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
		plugin.Logger(ctx).Debug("getAwsAccessAnalyzer", "ERROR", err)
		return nil, err
	}

	return data.Analyzer, nil
}
