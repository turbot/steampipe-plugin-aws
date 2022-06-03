package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodeBuildProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codebuild_project",
		Description: "AWS CodeBuild Project",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidInputException"}),
			},
			Hydrate: getCodeBuildProject,
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeBuildProjects,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the build project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the build project.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description that makes the build project easy to identify.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "concurrent_build_limit",
				Description: "The maximum number of concurrent builds that are allowed for this project.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "created",
				Description: "When the build project was created, expressed in Unix time format.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_modified",
				Description: "When the build project's settings were last modified, expressed in Unix time format.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "encryption_key",
				Description: "The AWS Key Management Service (AWS KMS) customer master key (CMK) to be.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "queued_timeout_in_minutes",
				Description: "The number of minutes a build is allowed to be queued before it times out.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "service_role",
				Description: "The ARN of the AWS Identity and Access Management (IAM) role that enables AWS CodeBuild to interact with dependent AWS services on behalf of the AWS account.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_version",
				Description: "A version of the build input to be built for this project.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timeout_in_minutes",
				Description: "How long, in minutes, from 5 to 480 (8 hours), for AWS CodeBuild to wait before timing out any related build that did not get marked as completed.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "artifacts",
				Description: "Information about the build output artifacts for the build project.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "badge",
				Description: "Information about the build badge for the build project.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "build_batch_config",
				Description: "A ProjectBuildBatchConfig object that defines the batch build options for the project.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cache",
				Description: "Information about the cache for the build project.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "environment",
				Description: "Information about the build environment for this build project.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "file_system_locations",
				Description: "An array of ProjectFileSystemLocation objects for a CodeBuild build project.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "logs_config",
				Description: "Information about logs for the build project. A project can create logs in Amazon CloudWatch Logs, an S3 bucket or both.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "project_visibility",
				Description: "Visibility of the build project.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "secondary_artifacts",
				Description: "An array of ProjectArtifacts objects.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "secondary_source_versions",
				Description: "An array of ProjectSource objects.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "secondary_sources",
				Description: "An array of ProjectSource objects.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source",
				Description: "Information about the build input source code for this build project.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_config",
				Description: "Information about the VPC configuration that AWS CodeBuild accesses.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "webhook",
				Description: " Information about a webhook that connects repository events to a build project in AWS CodeBuild.",
				Hydrate:     getCodeBuildProject,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tag key and value pairs associated with this build project.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeBuildProject,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
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
				Hydrate:     getCodeBuildProject,
				Transform:   transform.From(codeBuildProjectTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeBuildProject,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listCodeBuildProjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CodeBuildService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListProjectsPages(
		&codebuild.ListProjectsInput{},
		func(page *codebuild.ListProjectsOutput, isLast bool) bool {
			for _, result := range page.Projects {
				d.StreamListItem(ctx, &codebuild.Project{
					Name: result,
				})

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

func getCodeBuildProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCodeBuildProject")

	var name string
	if h.Item != nil {
		name = *h.Item.(*codebuild.Project).Name
	} else {
		quals := d.KeyColumnQuals
		name = quals["name"].GetStringValue()
	}

	// get service
	svc, err := CodeBuildService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &codebuild.BatchGetProjectsInput{
		Names: []*string{aws.String(name)},
	}

	// Get call
	op, err := svc.BatchGetProjects(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getCodeBuildProject__", "ERROR", err)
		return nil, err
	}

	if op.Projects != nil && len(op.Projects) > 0 {
		return op.Projects[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func codeBuildProjectTurbotTags(_ context.Context, d *transform.TransformData) (interface{},
	error) {
	data := d.HydrateItem.(*codebuild.Project)

	if data.Tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if data.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range data.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}

	}
	return turbotTagsMap, nil
}
