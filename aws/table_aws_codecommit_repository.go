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
				Name:        "name",
				Description: "The repository's name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RepositoryName"),
			},
			{
				Name:        "id",
				Description: "The ID of the repository.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RepositoryId"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description that makes the build project easy to identify.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RepositoryDescription"),
			},
			{
				Name:        "creation_date",
				Description: "The date and time the repository was created, in timestamp format.",
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
				Description: "The date and time the repository was last modified, in timestamp format.",
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
				Hydrate:     getCodeCommitRepositoryTag,
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

func listCodeCommitRepositories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listCodeCommitRepositories", "AWS_REGION", region)

	// Create session
	svc, err := CodeCommitService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	var repoNames []*string

	// List call
	err = svc.ListRepositoriesPages(
		&codecommit.ListRepositoriesInput{},
		func(page *codecommit.ListRepositoriesOutput, isLast bool) bool {
			if len(page.Repositories) != 0 {
				for _, repo := range page.Repositories {
					repoNames = append(repoNames, repo.RepositoryName)
				}
			}
			return !isLast
		},
	)
	if err != nil {
		return nil, err
	}

	// Build the params
	input := &codecommit.BatchGetRepositoriesInput{
		RepositoryNames: repoNames,
	}

	// Get all repository details
	result, err := svc.BatchGetRepositories(input)
	if err != nil {
		return nil, err
	}

	for _, repository := range result.Repositories {
		d.StreamListItem(ctx, repository)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCodeCommitRepositoryTag(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCodeCommitRepositoryTag")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	arn := h.Item.(*codecommit.RepositoryMetadata).Arn

	// get service
	svc, err := CodeCommitService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &codecommit.ListTagsForResourceInput{
		ResourceArn: arn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getCodeCommitRepositoryTag_", "ERROR", err)
		return nil, err
	}

	return op, nil
}
