package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"

	codecommitEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodeCommitRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codecommit_repository",
		Description: "AWS CodeCommit Repository",
		List: &plugin.ListConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameter"}),
			},
			Hydrate: listCodeCommitRepositories,
			Tags:    map[string]string{"service": "codecommit", "action": "ListRepositories"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listCodeCommitRepositoryTags,
				Tags: map[string]string{"service": "codecommit", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(codecommitEndpoint.CODECOMMITServiceID),
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
				Name:        "kms_key_id",
				Description: "The ID of the Key Management Service encryption key used to encrypt and decrypt the repository.",
				Type:        proto.ColumnType_STRING,
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
	svc, err := CodeCommitClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codecommit_repository.listCodeCommitRepositories", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &codecommit.ListRepositoriesInput{}
	// List all available repositories
	var repositoryNames []*string

	paginator := codecommit.NewListRepositoriesPaginator(svc, input, func(o *codecommit.ListRepositoriesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codecommit_repository.listCodeCommitRepositories", "api_error", err)
			return nil, err
		}

		for _, items := range output.Repositories {
			repositoryNames = append(repositoryNames, items.RepositoryName)
		}
	}

	if len(repositoryNames) <= 0 {
		return nil, nil
	}

	passedRepositoryNames := 0
	nameLeft := true
	for nameLeft {
		// BatchGetRepositories api can take maximum 25 number of repository name at a time.
		var names []*string
		if len(repositoryNames) > passedRepositoryNames {
			if (len(repositoryNames) - passedRepositoryNames) >= 25 {
				names = repositoryNames[passedRepositoryNames : passedRepositoryNames+25]
				passedRepositoryNames += 25
			} else {
				names = repositoryNames[passedRepositoryNames:]
				nameLeft = false
			}
		}

		// Build params
		params := &codecommit.BatchGetRepositoriesInput{
			RepositoryNames: aws.ToStringSlice(names),
		}

		// Get details for all available repositories
		result, err := svc.BatchGetRepositories(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codecommit_repository.BatchGetRepositories", "api_error", err)
			return nil, err
		}
		for _, repository := range result.Repositories {
			d.StreamListItem(ctx, repository)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func listCodeCommitRepositoryTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := CodeCommitClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codecommit_repository.listCodeCommitRepositoryTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}
	repositoryARN := h.Item.(types.RepositoryMetadata).Arn

	// Build the params
	params := &codecommit.ListTagsForResourceInput{
		ResourceArn: repositoryARN,
	}

	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codecommit_repository.listCodeCommitRepositoryTags", "api_error", err)
		return nil, err
	}

	return op, nil
}
