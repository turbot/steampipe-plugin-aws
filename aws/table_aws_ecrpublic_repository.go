package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcrpublicRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecrpublic_repository",
		Description: "AWS ECR Public Repository",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("repository_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"RepositoryNotFoundException", "RepositoryPolicyNotFoundException", "LifecyclePolicyNotFoundException"}),
			},
			Hydrate: getAwsEcrpublicRepository,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsEcrpublicRepositories,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "registry_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "repository_name",
				Description: "The name of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "registry_id",
				Description: "The AWS account ID associated with the public registry that contains the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that identifies the repository.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RepositoryArn"),
			},
			{
				Name:        "repository_uri",
				Description: "The URI for the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The date and time, in JavaScript date format, when the repository was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "image_details",
				Description: "A list of ImageDetail objects that contain data about the image.",
				Hydrate:     getAwsEcrpublicDescribeImages,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy",
				Description: "The JSON repository policy text associated with the repository.",
				Hydrate:     getAwsEcrpublicRepositoryPolicy,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyText"),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Hydrate:     getAwsEcrpublicRepositoryPolicy,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyText").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the repository.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsEcrpublicRepositoryTags,
				Transform:   transform.FromField("Tags"),
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
				Hydrate:     listAwsEcrpublicRepositoryTags,
				Transform:   transform.FromField("Tags").Transform(ecrpublicTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RepositoryArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsEcrpublicRepositories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// https://docs.aws.amazon.com/AmazonECR/latest/public/getting-started-cli.html
	// DescribeRepositories command is only supported in us-east-1
	region := d.KeyColumnQualString(matrixKeyRegion)

	if region != "us-east-1" {
		return nil, nil
	}

	// Create Session
	svc, err := ECRPublicClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecrpublic_repository.listAwsEcrpublicRepositories", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				maxLimit = 5
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ecrpublic.DescribeRepositoriesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["registry_id"] != nil {
		input.RegistryId = aws.String(equalQuals["registry_id"].GetStringValue())
	}

	paginator := ecrpublic.NewDescribeRepositoriesPaginator(svc, input, func(o *ecrpublic.DescribeRepositoriesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ecrpublic_repository.listAwsEcrpublicRepositories", "api_error", err)
			return nil, err
		}

		for _, items := range output.Repositories {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

////  HYDRATE FUNCTIONS

func getAwsEcrpublicRepository(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	region := d.KeyColumnQualString(matrixKeyRegion)

	// https://docs.aws.amazon.com/AmazonECR/latest/public/getting-started-cli.html
	// DescribeRepositories command is only supported in us-east-1
	if region != "us-east-1" {
		return nil, nil
	}

	name := d.KeyColumnQuals["repository_name"].GetStringValue()

	// Create Session
	svc, err := ECRPublicClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecrpublic_repository.getAwsEcrpublicRepository", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ecrpublic.DescribeRepositoriesInput{
		RepositoryNames: []string{name},
	}

	// Get call
	data, err := svc.DescribeRepositories(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecrpublic_repository.getAwsEcrpublicRepository", "api_error", err)
		return nil, err
	}
	if data.Repositories != nil && len(data.Repositories) > 0 {
		return data.Repositories[0], nil
	}

	return nil, nil
}

func listAwsEcrpublicRepositoryTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	resourceArn := h.Item.(types.Repository).RepositoryArn

	// Create Session
	svc, err := ECRPublicClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecrpublic_repository.listAwsEcrpublicRepositoryTags", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ecrpublic.ListTagsForResourceInput{
		ResourceArn: resourceArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecrpublic_repository.listAwsEcrpublicRepositoryTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getAwsEcrpublicRepositoryPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	repositoryName := h.Item.(types.Repository).RepositoryName

	// Create Session
	svc, err := ECRPublicClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecrpublic_repository.getAwsEcrpublicRepositoryPolicy", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ecrpublic.GetRepositoryPolicyInput{
		RepositoryName: repositoryName,
	}

	// Get call
	op, err := svc.GetRepositoryPolicy(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "RepositoryPolicyNotFoundException") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_ecrpublic_repository.getAwsEcrpublicRepositoryPolicy", "api_error", err)
		return nil, err

	}
	return op, nil
}

func getAwsEcrpublicDescribeImages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	repositoryName := h.Item.(types.Repository).RepositoryName

	// Create Session
	svc, err := ECRPublicClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecrpublic_repository.getAwsEcrpublicDescribeImages", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ecrpublic.DescribeImagesInput{
		RepositoryName: repositoryName,
	}

	// Get call
	op, err := svc.DescribeImages(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecrpublic_repository.getAwsEcrpublicDescribeImages", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func ecrpublicTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ecrpublicTagListToTurbotTags")
	tags := d.HydrateItem.(*ecrpublic.ListTagsForResourceOutput)

	if tags == nil {
		return nil, nil
	}

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
