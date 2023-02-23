package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codebuild"

	codebuildv1 "github.com/aws/aws-sdk-go/service/codebuild"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCodeBuildBuild(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codebuild_build",
		Description: "AWS CodeBuild Build",
		List: &plugin.ListConfig{
			Hydrate: listCodeBuildBuilds,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(codebuildv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The ARN of the batch build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "artifacts",
				Description: "A BuildArtifacts object the defines the build artifacts for this batch build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "build_batch_configuration",
				Description: "Contains configuration information about a batch build project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "build_batch_number",
				Description: "The number of the batch build.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "batch_build_status",
				Description: "The status of the batch build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "build_groups",
				Description: "An array of BuildGroup objects that define the build groups for the batch build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "build_timeout_in_minutes",
				Description: "Specifies the maximum amount of time, in minutes, that the build in a batch must be completed in.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "cache",
				Description: "Information about the cache for the build project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "complete",
				Description: "Indicates if the batch build is complete.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "current_phase",
				Description: "The current phase of the batch build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "debug_session_enabled",
				Description: "Specifies if session debugging is enabled for this batch build.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "encryption_key",
				Description: "The Key Management Service customer master key (CMK) to be used for encrypting the batch build output artifacts.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "end_time",
				Description: "The date and time that the batch build ended.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "environment",
				Description: "Information about the build environment of the build project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "file_system_locations",
				Description: "An array of ProjectFileSystemLocation objects for the batch build project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "id",
				Description: "The identifier of the batch build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "initiator",
				Description: "The entity that started the batch build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_config",
				Description: "Information about logs for a build project. These can be logs in CloudWatch Logs, built in a specified S3 bucket, or both.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "phases",
				Description: "An array of BuildBatchPhase objects the specify the phases of the batch build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "project_name",
				Description: "The name of the batch build project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "queued_timeout_in_minutes",
				Description: "Specifies the amount of time, in minutes, that the batch build is allowed to be queued before it times out.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "resolved_source_version",
				Description: "The identifier of the resolved version of this batch build's source code.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "secondary_artifacts",
				Description: "An array of BuildArtifacts objects the define the build artifacts for this batch build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "secondary_source_versions",
				Description: "An array of ProjectSourceVersion objects.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "secondary_sources",
				Description: "An array of ProjectSource objects that define the sources for the batch build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "service_role",
				Description: "The name of a service role used for builds in the batch.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source",
				Description: "Information about the build input source code for the build project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_version",
				Description: "The identifier of the version of the source code to be built.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_time",
				Description: "The date and time that the batch build started.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "vpc_config",
				Description: "Information about the VPC configuration that CodeBuild accesses.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BuildBatch.Arn"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

// listCodeBuildBuilds handles both listing and describing of the codebuild builds.
//
// The reason for this is the BatchGetBuildBatches call can accept up to 100 IDs. If we moved it out to another
// hydrate functions we may save a request or two if we only wanted to retrieve the IDs but the tradeoff is we need
// to get any other info an API call per codebuild build would need to be made. So in the case where we need to get
// all info for less then 100 instances including the BatchGetBuild request here, and batching requests means only making
// two API calls as opposed to 101.
func listCodeBuildBuilds(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listCodeBuildBuilds")

	// Create Session
	svc, err := CodeBuildClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// Limiting the results
	// maxLimit := int32(100)
	// if d.QueryContext.Limit != nil {
	// 	limit := int32(*d.QueryContext.Limit)
	// 	if limit < maxLimit {
	// 		if limit < 1 {
	// 			maxLimit = 1
	// 		} else {
	// 			maxLimit = limit
	// 		}
	// 	}
	// }

	input := &codebuild.ListBuildsInput{}

	paginator := codebuild.NewListBuildsPaginator(svc, input, func(o *codebuild.ListBuildsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codebuild_build.listCodeBuildBuild", "api_error", err)
			return nil, err
		}
		plugin.Logger(ctx).Error("listCodeBuildBuild", "output.Ids", output.Ids)
		if len(output.Ids) > 0 {
			params := &codebuild.BatchGetBuildsInput{
				Ids: output.Ids,
			}
			op, err := svc.BatchGetBuilds(ctx, params)
			if err != nil {
				plugin.Logger(ctx).Error("aws_codebuild_build.BatchGetBuilds", "api_error", err)
				return nil, err
			}

			for _, build := range op.Builds {
				d.StreamListItem(ctx, build)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}
