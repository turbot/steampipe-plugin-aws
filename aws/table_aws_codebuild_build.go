package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codebuild"

	codebuildEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

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
			Tags:    map[string]string{"service": "codebuild", "action": "ListBuilds"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(codebuildEndpoint.CODEBUILDServiceID),
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
				Name:        "current_phase",
				Description: "The current build phase.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encryption_key",
				Description: "The Key Management Service customer master key (CMK) to be used for encrypting the build output artifacts.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "end_time",
				Description: "The date and time that the build process ended, expressed in Unix time format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "initiator",
				Description: "The entity that started the build.",
				Type:        proto.ColumnType_STRING,
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
				Name:        "resolved_source_version",
				Description: "The identifier of the resolved version of this build's source code.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_role",
				Description: "The name of a service role used for this build.",
				Type:        proto.ColumnType_STRING,
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
				Name:        "logs",
				Description: "Information about the build's logs in CloudWatch Logs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_interfaces",
				Description: "Describes a network interface.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("NetworkInterface"),
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

			// Steampipe standard columns
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
// The reason for this is the BatchGetBuilds API call can accept up to 100 IDs. If we moved it out to another
// hydrate function we may save a request or two if we only wanted to retrieve the IDs but the tradeoff is that we would need
// to make an API call per codebuild build. So in cases where we need to get all information
// for less then 100 instances, including the BatchGetBuild request here, and batching requests means only making
// two API calls as opposed to 101.

func listCodeBuildBuilds(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := CodeBuildClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codebuild_build.listCodeBuildBuilds", "connection_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	buildId := quals["id"].GetStringValue()

	// If the user specifies a build id in optional quals, restrict BatchGetBuilds for other build ids.
	if buildId != "" {
		// Build param for a single build id
		params := &codebuild.BatchGetBuildsInput{
			Ids: []string{buildId},
		}

		err = getCodeBuildBuild(ctx, d, params)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	input := &codebuild.ListBuildsInput{}

	paginator := codebuild.NewListBuildsPaginator(svc, input, func(o *codebuild.ListBuildsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	var buildIds []string

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codebuild_build.listCodeBuildBuilds", "api_error", err)
			return nil, err
		}
		if len(output.Ids) > 0 {
			// Adding ids to a slice in order to do batch operations
			buildIds = append(buildIds, output.Ids...)
		} else {
			return nil, nil
		}
	}

	if len(buildIds) == 0 {
		return nil, nil
	}

	passedIds := 0
	idLeft := true

	for idLeft {
		// BatchGetBuilds api can take maximum 100 number of build id at a time.
		var ids []string
		if len(buildIds) > passedIds {
			if (len(buildIds) - passedIds) >= 100 {
				ids = buildIds[passedIds : passedIds+100]
				passedIds += 100
			} else {
				ids = buildIds[passedIds:]
				idLeft = false
			}
		}

		if len(ids) <= 0 {
			return nil, nil
		}

		// Build param
		params := &codebuild.BatchGetBuildsInput{
			Ids: ids,
		}

		err = getCodeBuildBuild(ctx, d, params)
		if err != nil {
			return nil, err
		}

	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

// Performing batch operation for the builds
func getCodeBuildBuild(ctx context.Context, d *plugin.QueryData, params *codebuild.BatchGetBuildsInput) error {

	// get service
	svc, err := CodeBuildClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codebuild_build.getCodeBuildBuild", "connection_error", err)
		return err
	}

	// Get call
	op, err := svc.BatchGetBuilds(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codebuild_build.getCodeBuildBuild", "api_error", err)
		return err
	}

	for _, build := range op.Builds {
		d.StreamListItem(ctx, build)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil
		}
	}

	return nil
}
