package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecrpublic"
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
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"RepositoryNotFoundException", "RepositoryPolicyNotFoundException", "LifecyclePolicyNotFoundException"}),
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
	svc, err := EcrPublicService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &ecrpublic.DescribeRepositoriesInput{
		MaxResults: aws.Int64(1000),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["registry_id"] != nil {
		input.RegistryId = aws.String(equalQuals["registry_id"].GetStringValue())
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 5 {
				input.MaxResults = aws.Int64(5)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.DescribeRepositoriesPages(
		input,
		func(page *ecrpublic.DescribeRepositoriesOutput, isLast bool) bool {
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

////  HYDRATE FUNCTIONS

func getAwsEcrpublicRepository(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsEcrpublicRepository")

	region := d.KeyColumnQualString(matrixKeyRegion)

	// https://docs.aws.amazon.com/AmazonECR/latest/public/getting-started-cli.html
	// DescribeRepositories command is only supported in us-east-1
	if region != "us-east-1" {
		return nil, nil
	}

	name := d.KeyColumnQuals["repository_name"].GetStringValue()

	// Create Session
	svc, err := EcrPublicService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ecrpublic.DescribeRepositoriesInput{
		RepositoryNames: []*string{aws.String(name)},
	}

	// Get call
	data, err := svc.DescribeRepositories(params)
	if err != nil {
		logger.Debug("getAwsEcrpublicRepository", "ERROR", err)
		return nil, err
	}
	if data.Repositories != nil && len(data.Repositories) > 0 {
		return data.Repositories[0], nil
	}

	return nil, nil
}

func listAwsEcrpublicRepositoryTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listAwsEcrpublicRepositoryTags")

	resourceArn := h.Item.(*ecrpublic.Repository).RepositoryArn

	// Create Session
	svc, err := EcrPublicService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ecrpublic.ListTagsForResourceInput{
		ResourceArn: resourceArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Debug("listAwsEcrpublicRepositoryTags", "ERROR", err)
		return nil, err
	}
	return op, nil
}

func getAwsEcrpublicRepositoryPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsEcrpublicRepositoryPolicy")

	repositoryName := h.Item.(*ecrpublic.Repository).RepositoryName

	// Create Session
	svc, err := EcrPublicService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ecrpublic.GetRepositoryPolicyInput{
		RepositoryName: repositoryName,
	}

	// Get call
	op, err := svc.GetRepositoryPolicy(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "RepositoryPolicyNotFoundException" {
				return nil, nil
			}
			return nil, err
		}
	}
	return op, nil
}

func getAwsEcrpublicDescribeImages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsEcrpublicDescribeImages")

	repositoryName := h.Item.(*ecrpublic.Repository).RepositoryName

	// Create Session
	svc, err := EcrPublicService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ecrpublic.DescribeImagesInput{
		RepositoryName: repositoryName,
	}

	// Get call
	op, err := svc.DescribeImages(params)
	if err != nil {
		logger.Debug("getAwsEcrpublicDescribeImages", "ERROR", err)
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
