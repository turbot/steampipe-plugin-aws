package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact/types"

	codeartifactEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodeArtifactDomain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codeartifact_domain",
		Description: "AWS Code Artifact Domain",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Required,
				},
				{
					Name:    "owner",
					Require: plugin.Optional,
				},
			},
			Hydrate: getCodeArtifactDomain,
			Tags:    map[string]string{"service": "codeartifact", "action": "DescribeDomain"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeArtifactDomains,
			Tags:    map[string]string{"service": "codeartifact", "action": "ListDomains"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCodeArtifactDomainTags,
				Tags: map[string]string{"service": "codeartifact", "action": "ListTagsForResource"},
			},
			{
				Func: getCodeArtifactDomainPermissionsPolicy,
				Tags: map[string]string{"service": "codeartifact", "action": "GetDomainPermissionsPolicy"},
			},
			{
				Func: getCodeArtifactDomain,
				Tags: map[string]string{"service": "codeartifact", "action": "DescribeDomain"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(codeartifactEndpoint.AWS_CODEARTIFACT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "asset_size_bytes",
				Description: "The total size of all assets in the domain.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCodeArtifactDomain,
			},
			{
				Name:        "created_time",
				Description: "A timestamp that contains the date and time the domain was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "encryption_key",
				Description: "The key used to encrypt the domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner",
				Description: "The 12-digit account number of the Amazon Web Services account that owns the domain. It does not include dashes or spaces.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_count",
				Description: "The number of repositories in the domain.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCodeArtifactDomain,
			},
			{
				Name:        "s3_bucket_arn",
				Description: "The Amazon Resource Name (ARN) of the Amazon S3 bucket that is used to store package assets in the domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeArtifactDomain,
				Transform:   transform.FromField("S3BucketArn"),
			},
			{
				Name:        "status",
				Description: "A string that contains the status of the domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy",
				Description: "An CodeArtifact resource policy that contains a resource ARN, document details, and a revision.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeArtifactDomainPermissionsPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "policy_std",
				Description: "Contains the contents of the resource-based policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeArtifactDomainPermissionsPolicy,
				Transform:   transform.FromValue().Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeArtifactDomainTags,
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
				Hydrate:     getCodeArtifactDomainTags,
				Transform:   transform.FromField("Tags").Transform(codeArtifactDomainTurbotTags),
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

func listCodeArtifactDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create session
	svc, err := CodeArtifactClient(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_domain.listCodeArtifactDomains", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(100)
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

	input := codeartifact.ListDomainsInput{
		MaxResults: &maxLimit,
	}

	paginator := codeartifact.NewListDomainsPaginator(svc, &input, func(o *codeartifact.ListDomainsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codeartifact_domain.listCodeArtifactDomains", "api_error", err)
			return nil, err
		}

		for _, domain := range output.Domains {
			item := &types.DomainDescription{
				Arn:           domain.Arn,
				CreatedTime:   domain.CreatedTime,
				EncryptionKey: domain.EncryptionKey,
				Name:          domain.Name,
				Owner:         domain.Owner,
				Status:        domain.Status,
			}
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getCodeArtifactDomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var name, owner string
	if h.Item != nil {
		data := h.Item.(*types.DomainDescription)
		name = *data.Name
		owner = *data.Owner
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
		owner = d.EqualsQuals["owner"].GetStringValue()
	}

	if name == "" {
		return nil, nil
	}

	// Build the params
	params := &codeartifact.DescribeDomainInput{
		Domain: &name,
	}
	if owner != "" {
		params.DomainOwner = &owner
	}

	// Create session
	svc, err := CodeArtifactClient(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomain", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Get call
	data, err := svc.DescribeDomain(ctx, params)
	if err != nil {
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomain", "api_error", err)
		return nil, err
	}
	return data.Domain, nil
}

func getCodeArtifactDomainTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	data := h.Item.(*types.DomainDescription)

	// Create Session
	svc, err := CodeArtifactClient(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomainTags", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &codeartifact.ListTagsForResourceInput{
		ResourceArn: aws.String(*data.Arn),
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomainTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getCodeArtifactDomainPermissionsPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	data := h.Item.(*types.DomainDescription)

	// Create Session
	svc, err := CodeArtifactClient(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomainPermissionsPolicy", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &codeartifact.GetDomainPermissionsPolicyInput{
		Domain:      aws.String(*data.Name),
		DomainOwner: aws.String(*data.Owner),
	}

	// Get call
	op, err := svc.GetDomainPermissionsPolicy(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFoundException") {
			return nil, nil
		}
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomainPermissionsPolicy", "api_error", err)
		return nil, err
	}
	if op != nil && op.Policy != nil {
		return op.Policy.Document, nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func codeArtifactDomainTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
