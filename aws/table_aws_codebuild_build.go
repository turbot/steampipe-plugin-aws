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
				Description: "The ARN of the build.",
				Type:        proto.ColumnType_STRING,
			},			
			{
				Name:        "id",
				Description: "The unique identifier of the  build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "build_batch_arn",
				Description: "The ARN of the batch build that this build is a member of, if applicable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "build_complete",
				Description: "Indicates if the build is complete.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "build_number",
				Description: "The number of the build.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "build_status",
				Description: "The status of the build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "encryption_key",
				Description: "The Key Management Service customer master key (CMK) to be used for encrypting the build output artifacts.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "end_time",
				Description: "The date and time that the build process ended, expressed in Unix time forma.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "project_name",
				Description: "The name of the build project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "queued_timeout_in_minutes",
				Description: "Specifies the amount of time, in minutes, that a build is allowed to be queued before it times out.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "source_version",
				Description: "The identifier of the version of the source code to be built.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_time",
				Description: "The date and time that the build started.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "timeout_in_minutes",
				Description: "How long, in minutes, for CodeBuild to wait before timing out this build if it does not get marked as completed.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "artifacts",
				Description: "A BuildArtifacts object the defines the build artifacts for this build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cache",
				Description: "Information about the cache for the build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "current_phase",
				Description: "The current build phase.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "debug_session",
				Description: "Contains information about the debug session for this build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "environment",
				Description: "Information about the build environment for this build project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "exported_environment_variables",
				Description: "A list of exported environment variables for this build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "file_system_locations",
				Description: "An array of ProjectFileSystemLocation objects for a CodeBuild build project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "initiator",
				Description: "The entity that started the build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "logs",
				Description: "Information about the build's logs in CloudWatch Logs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_interfaces",
				Description: " Describes a network interface.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "phases",
				Description: "Information about all previous build phases that are complete and information about any current build phase that is not yet complete.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "report_arns",
				Description: "An array of the ARNs associated with this build's reports.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resolved_source_version",
				Description: "The identifier of the resolved version of this build's source code.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "secondary_artifacts",
				Description: "An array of BuildArtifacts objects the define the build artifacts for this build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "secondary_source_versions",
				Description: "An array of ProjectSourceVersion objects.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "secondary_sources",
				Description: "An array of ProjectSource objects that define the sources for the build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source",
				Description: "Information about the build input source code for the build project.",
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.FromField("Arn"),
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
// The reason for this is the BatchGetBuilds call can accept up to 100 IDs. If we moved it out to another
// hydrate functions we may save a request or two if we only wanted to retrieve the IDs but the tradeoff is we need
// to get any other info an API call per codebuild build would need to be made. So in the case where we need to get
// all info for less then 100 instances including the BatchGetBuild request here, and batching requests means only making
// two API calls as opposed to 101.
func listCodeBuildBuilds(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listCodeBuildBuilds")

	// Create Session
	svc, err := CodeBuildClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codebuild_build.listCodeBuildBuild", "get_client_error", err)
		return nil, err
	}

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
