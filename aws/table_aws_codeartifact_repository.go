package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
)

//// TABLE DEFINITION

func tableAwsCodeArtifactRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codeartifact_repository",
		Description: "AWS CodeArtifact Repository",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Required,
				},
				{
					Name:    "domain_name",
					Require: plugin.Required,
				},
				{
					Name:    "domain_owner",
					Require: plugin.Optional,
				},
			},
			Hydrate: getCodeArtifactRepository,
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeArtifactRepositories,
		},
		HydrateDependencies: []plugin.HydrateDependencies{
			{
				Func:    getCodeArtifactRepositoryEndpoints,
				Depends: []plugin.HydrateFunc{getCodeArtifactRepository},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain_name",
				Description: "The name of the domain that contains the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "administrator_account",
				Description: "The Amazon Web Services account ID that manages the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the repository.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeArtifactRepository,
			},
			{
				Name:        "domain_owner",
				Description: "The 12-digit account number of the Amazon Web Services account that owns the repository. It does not include dashes or spaces.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "external_connections",
				Description: "An array of external connections associated with the repository.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeArtifactRepository,
			},
			{
				Name:        "policy",
				Description: "An CodeArtifact resource policy that contains a resource ARN, document details, and a revision.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeArtifactRepositoryPermissionsPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "policy_std",
				Description: "Contains the contents of the resource-based policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeArtifactRepositoryPermissionsPolicy,
				Transform:   transform.FromValue().Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "repository_endpoint",
				Description: "A string that specifies the URL of the returned endpoint.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeArtifactRepositoryEndpoints,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "upstreams",
				Description: "A list of upstream repositories to associate with the repository.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeArtifactRepository,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeArtifactRepositoryTags,
				Transform:   transform.FromField("Tags"),
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
				Hydrate:     getCodeArtifactRepositoryTags,
				Transform:   transform.FromField("Tags").Transform(codeArtifactRepositoryTurbotTags),
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

func listCodeArtifactRepositories(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create session
	svc, err := CodeArtifactClient(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_repository.listCodeArtifactRepositories", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// un-supported region check
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := codeartifact.ListRepositoriesInput{
		MaxResults: &maxLimit,
	}

	paginator := codeartifact.NewListRepositoriesPaginator(svc, &input, func(o *codeartifact.ListRepositoriesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codeartifact_repository.listCodeArtifactRepositories", "api_error", err)
			return nil, err
		}

		for _, repository := range output.Repositories {
			item := &types.RepositoryDescription{
				AdministratorAccount: repository.AdministratorAccount,
				Arn:                  repository.Arn,
				Description:          repository.Description,
				DomainName:           repository.DomainName,
				DomainOwner:          repository.DomainOwner,
				Name:                 repository.Name,
			}
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCodeArtifactRepository(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var domainName, name, owner string
	if h.Item != nil {
		data := h.Item.(*types.RepositoryDescription)
		name = *data.Name
		owner = *data.DomainOwner
		domainName = *data.DomainName
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
		owner = d.KeyColumnQuals["owner"].GetStringValue()
		domainName = d.KeyColumnQuals["domain_name"].GetStringValue()
	}

	if name == "" || domainName == "" {
		return nil, nil
	}

	// Build the params
	params := &codeartifact.DescribeRepositoryInput{
		Repository: &name,
		Domain:     &domainName,
	}
	if owner != "" {
		params.DomainOwner = &owner
	}

	// Create session
	svc, err := CodeArtifactClient(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepository", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// un-supported region check
		return nil, nil
	}

	// Get call
	data, err := svc.DescribeRepository(ctx, params)
	if err != nil {
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepository", "api_error", err)
		return nil, err
	}
	return data.Repository, nil
}

func getCodeArtifactRepositoryEndpoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	resultData := []string{}
	repository := h.HydrateResults["getCodeArtifactRepository"].(*types.RepositoryDescription)

	if len(repository.ExternalConnections) == 0 {
		return nil, nil
	}

	for _, item := range repository.ExternalConnections {

		// Build the params
		params := &codeartifact.GetRepositoryEndpointInput{
			Repository:  repository.Name,
			Domain:      repository.DomainName,
			DomainOwner: repository.DomainOwner,
			Format:      item.PackageFormat,
		}

		// Create session
		svc, err := CodeArtifactClient(ctx, d)
		if err != nil {
			logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryEndpoints", "service_creation_error", err)
			return nil, err
		}
		if svc == nil {
			// un-supported region check
			return nil, nil
		}

		// Get call
		data, err := svc.GetRepositoryEndpoint(ctx, params)

		if err != nil {
			logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryEndpoints", "api_error", err)
			return nil, err
		}
		resultData = append(resultData, *data.RepositoryEndpoint)
	}

	return resultData, nil
}

func getCodeArtifactRepositoryTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	data := h.Item.(*types.RepositoryDescription)

	// Create Session
	svc, err := CodeArtifactClient(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryTags", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// un-supported region check
		return nil, nil
	}

	// Build the params
	params := &codeartifact.ListTagsForResourceInput{
		ResourceArn: aws.String(*data.Arn),
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getCodeArtifactRepositoryPermissionsPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	data := h.Item.(*types.RepositoryDescription)

	// Create Session
	svc, err := CodeArtifactClient(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryPermissionsPolicy", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// un-supported region check
		return nil, nil
	}

	// Build the params
	params := &codeartifact.GetRepositoryPermissionsPolicyInput{
		Repository:  aws.String(*data.Name),
		DomainOwner: aws.String(*data.DomainOwner),
		Domain:      aws.String(*data.DomainName),
	}

	// Get call
	op, err := svc.GetRepositoryPermissionsPolicy(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFoundException") {
			return nil, nil
		}
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryPermissionsPolicy", "api_error", err)
		return nil, err
	}
	if op != nil && op.Policy != nil {
		return op.Policy.Document, nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func codeArtifactRepositoryTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*codeartifact.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}
