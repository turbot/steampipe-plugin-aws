package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/synthetics"
	"github.com/aws/aws-sdk-go-v2/service/synthetics/types"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v6/query_cache"
)

func tableAwsSyntheticsCanaryRun(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_synthetics_canary_run",
		Description: "AWS CloudWatch Synthentics Canary Run",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "name",
					Require: plugin.Required,
				},
				{
					Name:       "dry_run_id",
					Require:    plugin.Optional,
					CacheMatch: query_cache.CacheMatchExact,
				},
				{
					Name:       "run_type",
					Require:    plugin.Optional,
					CacheMatch: query_cache.CacheMatchExact,
				},
			},
			Hydrate: listSyntheticsCanaryRuns,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags: map[string]string{"service": "synthetics", "action": "GetCanaryRuns"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SYNTHETICS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "artifact_s3_location",
				Description: "The S3 location where artifacts are stored for the canary run.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ArtifactS3Location"),
			},
			{
				Name:        "browser_type",
				Description: "The browser type for the canary run.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dry_run_config",
				Description: "The dry run configuration for the canary run.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "id",
				Description: "The unique ID for the canary run.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the canary.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "retry_attempt",
				Description: "The number of retries for the canary run.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "scheduled_run_id",
				Description: "The ID of the scheduled canary run.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the canary run.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "timeline",
				Description: "The timeline information for the canary run.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dry_run_id",
				Description: "The ID of the canary dry run.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("dry_run_id"),
			},
			{
				Name:        "run_type",
				Description: "The type of the canary dry run.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("run_type"),
			},
		}),
	}
}

//// LIST FUNCTION

func listSyntheticsCanaryRuns(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := SyntheticsClient(ctx, d, d.EqualsQualString(matrixKeyRegion))
	if err != nil {
		plugin.Logger(ctx).Error("aws_synthetics_canary_run.listSyntheticsCanaryRuns", "client_error", err)
		return nil, err
	}

	// Compile inputs
	input := &synthetics.GetCanaryRunsInput{}
	input.Name = aws.String(d.EqualsQuals["name"].GetStringValue())

	if d.EqualsQuals["dry_run_id"] != nil {
		if d.EqualsQuals["dry_run_id"].GetStringValue() != "" {
			input.DryRunId = aws.String(d.EqualsQuals["dry_run_id"].GetStringValue())
		}
	}

	if d.EqualsQuals["run_type"] != nil {
		if d.EqualsQuals["run_type"].GetStringValue() != "" {
			input.RunType = types.RunType(d.EqualsQuals["run_type"].GetStringValue())
		}
	}

	// Create paginator
	paginator := synthetics.NewGetCanaryRunsPaginator(svc, input, func(o *synthetics.GetCanaryRunsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// Apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_synthetics_canary_run.listSyntheticsCanaryRuns", "api_error", err)
			return nil, err
		}

		for _, runDetail := range output.CanaryRuns {
			d.StreamListItem(ctx, runDetail)

			// Context may be cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
