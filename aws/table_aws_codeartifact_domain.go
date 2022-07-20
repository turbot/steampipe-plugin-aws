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
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeArtifactDomains,
		},
		GetMatrixItem: BuildRegionList,
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

	region := d.KeyColumnQualString(matrixKeyRegion)
	serviceId := "codeartifact"
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	// Create session
	svc, err := CodeArtifactService(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_domain.listCodeArtifactDomains", "service_creation_error", err)
		return nil, err
	}

	input := &codeartifact.ListDomainsInput{
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
	err = svc.ListDomainsPages(
		input,
		func(page *codeartifact.ListDomainsOutput, isLast bool) bool {
			for _, domain := range page.Domains {
				d.StreamListItem(ctx, domain)

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

func getCodeArtifactDomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.KeyColumnQualString(matrixKeyRegion)
	serviceId := "codeartifact"
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	var name, owner string
	if h.Item != nil {
		data := domainData(h.Item, ctx, d, h)
		name = data["Name"]
		owner = data["Owner"]
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
		owner = d.KeyColumnQuals["owner"].GetStringValue()
	}

	if name == "" {
		return nil, nil
	}

	// Build the params
	params := &codeartifact.DescribeDomainInput{
		Domain: &name,
	}
	if owner != "" {
		params.SetDomainOwner(owner)
	}

	// Create session
	svc, err := CodeArtifactService(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomain", "service_creation_error", err)
		return nil, err
	}

	// Get call
	data, err := svc.DescribeDomain(params)
	if err != nil {
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomain", "api_error", err)
		return nil, err
	}
	return data.Domain, nil
}

func getCodeArtifactDomainTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	data := domainData(h.Item, ctx, d, h)
	region := d.KeyColumnQualString(matrixKeyRegion)
	serviceId := "codeartifact"
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	// Create Session
	svc, err := CodeArtifactService(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomainTags", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &codeartifact.ListTagsForResourceInput{
		ResourceArn: aws.String(data["Arn"]),
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomainTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getCodeArtifactDomainPermissionsPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	data := domainData(h.Item, ctx, d, h)
	region := d.KeyColumnQualString(matrixKeyRegion)
	serviceId := "codeartifact"
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	// Create Session
	svc, err := CodeArtifactService(ctx, d)
	if err != nil {
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomainPermissionsPolicy", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &codeartifact.GetDomainPermissionsPolicyInput{
		Domain:      aws.String(data["Name"]),
		DomainOwner: aws.String(data["Owner"]),
	}

	// Get call
	op, err := svc.GetDomainPermissionsPolicy(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ResourceNotFoundException" {
				return nil, nil
			}
			return nil, err
		}
		logger.Error("aws_codeartifact_domain.getCodeArtifactDomainPermissionsPolicy", "api_error", err)
	}
	return op.Policy.Document, nil
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

func domainData(item interface{}, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case *codeartifact.DomainSummary:
		data["Arn"] = *item.Arn
		data["Name"] = *item.Name
		data["Owner"] = *item.Owner
	case *codeartifact.DomainDescription:
		data["Arn"] = *item.Arn
		data["Name"] = *item.Name
		data["Owner"] = *item.Owner
	}
	return data
}
