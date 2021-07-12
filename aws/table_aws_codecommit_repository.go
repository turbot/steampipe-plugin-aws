package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
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
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("repository_name"),
			ShouldIgnoreError: isNotFoundError([]string{"RepositoryDoesNotExistException", "InvalidRepositoryNameException"}),
			Hydrate:           getCodeCommitRepository,
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeCommitRepositories,
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
				Hydrate:     getCodeCommitRepository,
			},
			{
				Name:        "description",
				Description: "A comment or description about the repository.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeCommitRepository,
				Transform:   transform.FromField("RepositoryDescription"),
			},
			{
				Name:        "creation_date",
				Description: "The date and time the repository was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCodeCommitRepository,
			},
			{
				Name:        "clone_url_http",
				Description: "The URL to use for cloning the repository over HTTPS.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeCommitRepository,
			},
			{
				Name:        "clone_url_ssh",
				Description: "The URL to use for cloning the repository over SSH.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeCommitRepository,
			},
			{
				Name:        "default_branch",
				Description: "The repository's default branch name.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeCommitRepository,
			},
			{
				Name:        "last_modified_date",
				Description: "The date and time the repository was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCodeCommitRepository,
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
				Hydrate:     getCodeCommitRepository,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listCodeCommitRepositories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listCodeCommitRepositories", "AWS_REGION", region)

	// Create service
	svc, err := CodeCommitService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListRepositoriesPages(
		&codecommit.ListRepositoriesInput{},
		func(page *codecommit.ListRepositoriesOutput, isLast bool) bool {
			for _, repository := range page.Repositories {
				d.StreamListItem(ctx, repository)
			}
			return !isLast
		},
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCodeCommitRepository(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCodeCommitRepository")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Get repository name
	var repositoryName string
	if h.Item != nil {
		repositoryName = *h.Item.(*codecommit.RepositoryNameIdPair).RepositoryName
	} else {
		repositoryName = d.KeyColumnQuals["repository_name"].GetStringValue()
	}

	// Return nil, if no input provided
	if repositoryName == "" {
		return nil, nil
	}

	// Create service
	svc, err := CodeCommitService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build params
	params := &codecommit.GetRepositoryInput{
		RepositoryName: aws.String(repositoryName),
	}

	op, err := svc.GetRepository(params)
	if err != nil {
		return nil, err
	}

	return op.RepositoryMetadata, nil
}

func listCodeCommitRepositoryTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listCodeCommitRepositoryTags")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create service
	svc, err := CodeCommitService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Get common columns
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get repository ARN
	var repositoryARN string
	switch item := h.Item.(type) {
	case *codecommit.RepositoryNameIdPair:
		repositoryARN = "arn:" + commonColumnData.Partition + ":codecommit:" + region + ":" + commonColumnData.AccountId + ":" + *item.RepositoryName
	case *codecommit.RepositoryMetadata:
		repositoryARN = *item.Arn
	}

	// Build the params
	params := &codecommit.ListTagsForResourceInput{
		ResourceArn: aws.String(repositoryARN),
	}

	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}
