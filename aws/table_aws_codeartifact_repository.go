package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/codeartifact"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
)

//// TABLE DEFINITION

func tableAwsCodeArtifactRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codeartifact_repository",
		Description: "AWS Code Artifact Repository",
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
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the repository.",
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
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the repository.",
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

	region := d.KeyColumnQualString(matrixKeyRegion)
	serviceId := "codeartifact"
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	// Create session
	svc, err := CodeArtifactService(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_repository.listCodeArtifactRepositories", "service_creation_error", err)
		return nil, err
	}

	input := &codeartifact.ListRepositoriesInput{
		MaxResults: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.ListRepositoriesPages(
		input,
		func(page *codeartifact.ListRepositoriesOutput, isLast bool) bool {
			for _, repository := range page.Repositories {
				d.StreamListItem(ctx, repository)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
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

func getCodeArtifactRepository(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.KeyColumnQualString(matrixKeyRegion)
	serviceId := "codeartifact"
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	var domainName, name, owner string
	if h.Item != nil {
		data := repositoryData(h.Item, ctx, d, h)
		name = data["Name"]
		owner = data["DomainOwner"]
		domainName = data["DomainName"]
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
		params.SetDomainOwner(owner)
	}

	// Create session
	svc, err := CodeArtifactService(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepository", "service_creation_error", err)
		return nil, err
	}

	// Get call
	data, err := svc.DescribeRepository(params)
	if err != nil {
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepository", "api_error", err)
		return nil, err
	}
	return data.Repository, nil
}

func getCodeArtifactRepositoryEndpoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.KeyColumnQualString(matrixKeyRegion)
	serviceId := "codeartifact"
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	resultData := []string{}
	repository := h.HydrateResults["getCodeArtifactRepository"].(*codeartifact.RepositoryDescription)

	if len(repository.ExternalConnections) == 0 {
		return nil, nil
	}

	for _, item := range repository.ExternalConnections {

		// Build the params
		params := &codeartifact.GetRepositoryEndpointInput{
			Repository:  aws.String(*repository.Name),
			Domain:      aws.String(*repository.DomainName),
			DomainOwner: aws.String(*repository.DomainOwner),
			Format:      aws.String(*item.PackageFormat),
		}

		// Create session
		svc, err := CodeArtifactService(ctx, d)
		if err != nil {
			logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryEndpoint", "service_creation_error", err)
			return nil, err
		}

		// Get call
		data, err := svc.GetRepositoryEndpoint(params)

		if err != nil {
			logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryEndpoint", "api_error", err)
			return nil, err
		}
		resultData = append(resultData, *data.RepositoryEndpoint)
	}

	return resultData, nil
}

func getCodeArtifactRepositoryTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	data := repositoryData(h.Item, ctx, d, h)
	region := d.KeyColumnQualString(matrixKeyRegion)
	serviceId := "codeartifact"
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	// Create Session
	svc, err := CodeArtifactService(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryTags", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &codeartifact.ListTagsForResourceInput{
		ResourceArn: aws.String(data["Arn"]),
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getCodeArtifactRepositoryPermissionsPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	data := repositoryData(h.Item, ctx, d, h)
	region := d.KeyColumnQualString(matrixKeyRegion)
	serviceId := "codeartifact"
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	// Create Session
	svc, err := CodeArtifactService(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryPermissionsPolicy", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &codeartifact.GetRepositoryPermissionsPolicyInput{
		Repository:  aws.String(data["Name"]),
		DomainOwner: aws.String(data["DomainOwner"]),
		Domain:      aws.String(data["DomainName"]),
	}

	// Get call
	op, err := svc.GetRepositoryPermissionsPolicy(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ResourceNotFoundException" {
				return nil, nil
			}
			return nil, err
		}
		logger.Error("aws_codeartifact_repository.getCodeArtifactRepositoryPermissionsPolicy", "api_error", err)
	}
	return op.Policy.Document, nil
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

func repositoryData(item interface{}, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case *codeartifact.RepositorySummary:
		data["Arn"] = *item.Arn
		data["Name"] = *item.Name
		data["DomainOwner"] = *item.DomainOwner
		data["DomainName"] = *item.DomainName
	case *codeartifact.RepositoryDescription:
		data["Arn"] = *item.Arn
		data["Name"] = *item.Name
		data["DomainOwner"] = *item.DomainOwner
		data["DomainName"] = *item.DomainName
	}
	return data
}
