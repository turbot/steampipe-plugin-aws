package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact/types"

	codeartifactv1 "github.com/aws/aws-sdk-go/service/codeartifact"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodeArtifactPackage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codeartifact_package",
		Description: "AWS CodeArtifact Package",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "domain",
					Require: plugin.Required,
				},
				{
					Name:    "repository",
					Require: plugin.Required,
				},
				{
					Name:    "name",
					Require: plugin.Required,
				},
				{
					Name:    "format",
					Require: plugin.Required,
				},
			},
			Hydrate: getCodeArtifactPackage,
			Tags:    map[string]string{"service": "codeartifact", "action": "DescribePackage"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCodeArtifactRepositories,
			Hydrate:       listCodeArtifactPackages,
			Tags:    map[string]string{"service": "codeartifact", "action": "ListPackages"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "domain",
					Require: plugin.Required,
				},
				{
					Name:    "repository",
					Require: plugin.Required,
				},
				{
					Name:       "domain_owner",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "next_token",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "package_prefix",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "format",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "namespace",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCodeArtifactPackage,
				Tags: map[string]string{"service": "codeartifact", "action": "DescribePackage"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(codeartifactv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The description of the repository.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeArtifactPackage,
			},
			{
				Name:        "format",
				Description: "The format of the package.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace",
				Description: "The namespace of the package. The package component that specifies its namespace depends on its type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "origin_configuration",
				Description: "A PackageOriginConfiguration object that contains a PackageOriginRestrictions object that contains information about the upstream and publish package origin restrictions.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "domain",
				Description: "The domain that the package is associated to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository",
				Description: "The repository that the package is associated to.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCodeArtifactPackages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create session
	svc, err := CodeArtifactClient(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_package.listCodeArtifactPackages", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	var domain, repository string
	if h.Item != nil {
		domain = *h.Item.(*types.RepositoryDescription).DomainName
		repository = *h.Item.(*types.RepositoryDescription).Name
	} else {
		domain = d.EqualsQualString("domain")
		repository = d.EqualsQualString("repository")
	}

	input := codeartifact.ListPackagesInput{
		Domain:     &domain,
		Repository: &repository,
		MaxResults: &maxLimit,
	}

	paginator := codeartifact.NewListPackagesPaginator(svc, &input, func(o *codeartifact.ListPackagesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codeartifact_package.listCodeArtifactPackages", "api_error", err)
			return nil, err
		}

		for _, items := range output.Packages {
			d.StreamListItem(ctx, items)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCodeArtifactPackage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	data := h.Item.(*types.PackageDescription)

	PackageName := d.EqualsQualString("package")
	format := d.EqualsQualString("format")

	var domain, repository string
	if h.Item != nil {
		domain = *h.Item.(*types.RepositoryDescription).DomainName
		repository = *h.Item.(*types.RepositoryDescription).Name
	} else {
		domain = d.EqualsQualString("domain")
		repository = d.EqualsQualString("repository")
	}

	if PackageName == "" || format == "" || domain == "" || repository == "" {
		return nil, nil
	}

	// Build the params
	params := &codeartifact.DescribePackageInput{
		Domain:     &domain,
		Repository: &repository,
		Format:     data.Format,
		Package:    data.Name,
	}

	// Create session
	svc, err := CodeArtifactClient(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_package.getCodeArtifactPackage", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Get call
	op, err := svc.DescribePackage(ctx, params)
	if err != nil {
		logger.Error("aws_codeartifact_package.getCodeArtifactPackage", "api_error", err)
		return nil, err
	}
	return op.Package, nil
}
