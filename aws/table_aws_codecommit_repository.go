package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodeCommitRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codecommit_repository",
		Description: "AWS CodeCommit Repository",
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"InvalidParameter"}),
			Hydrate:           listCodeCommitRepositories,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "repository_name",
				Description: "The repository's name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_id",
				Description: "The ID of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A comment or description about the repository.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RepositoryDescription"),
			},
			{
				Name:        "creation_date",
				Description: "The date and time the repository was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "clone_url_http",
				Description: "The URL to use for cloning the repository over HTTPS.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "clone_url_ssh",
				Description: "The URL to use for cloning the repository over SSH.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_branch",
				Description: "The repository's default branch name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified_date",
				Description: "The date and time the repository was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RepositoryName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listCodeCommitRepositoryTags,
				Transform:   transform.FromField("Tags"),
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

func listCodeCommitRepositories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := CodeCommitService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List all available repositories
	var repositoryNames []*string
	err = svc.ListRepositoriesPages(
		&codecommit.ListRepositoriesInput{},
		func(page *codecommit.ListRepositoriesOutput, isLast bool) bool {
			for _, data := range page.Repositories {
				repositoryNames = append(repositoryNames, data.RepositoryName)
			}
			return !isLast
		},
	)
	if err != nil {
		return nil, err
	}

	// Build params
	params := &codecommit.BatchGetRepositoriesInput{
		RepositoryNames: repositoryNames,
	}

	// Get details for all available repositories
	result, err := svc.BatchGetRepositories(params)
	if err != nil {
		return nil, err
	}
	for _, repository := range result.Repositories {
		d.StreamListItem(ctx, repository)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func listCodeCommitRepositoryTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listCodeCommitRepositoryTags")

	// Create service
	svc, err := CodeCommitService(ctx, d)
	if err != nil {
		return nil, err
	}
	repositoryARN := h.Item.(*codecommit.RepositoryMetadata).Arn

	// Build the params
	params := &codecommit.ListTagsForResourceInput{
		ResourceArn: repositoryARN,
	}

	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}
