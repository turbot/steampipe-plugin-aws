package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsECRRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecr_repository",
		Description: "AWS ECR Repository",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("repository_name"),
			ShouldIgnoreError: isNotFoundError([]string{"RepositoryNotFoundException"}),
			Hydrate:           getAwsECRRepositories,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsECRRepositories,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "repository_name",
				Description: "The name of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "registry_id",
				Description: "The AWS account ID associated with the registry that contains the repositories to be described.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_results",
				Description: "The maximum number of repository results returned by DescribeRepositories.",
				Hydrate:     getAwsECRRepositories,
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "repository_arn",
				Description: "The Amazon Resource Name (ARN) that identifies the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_uri",
				Description: "The URI for the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encryption_configuration",
				Description: "The encryption configuration for the repository.",
				Hydrate:     getAwsECRRepositories,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "created_at",
				Description: "The date and time, in JavaScript date format, when the repository was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "image_tag_mutability",
				Description: "The tag mutability setting for the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_scanning_configuration",
				Description: "The image scanning configuration for a repository.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "image_details",
				Description: "A list of ImageDetail objects that contain data about the image.",
				Hydrate:     getAwsECRDescribeImages,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Repository",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsECRRepositoryTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsECRRepositoryTags,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RepositoryName"),
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

func listAwsECRRepositories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsECRRepositories", "AWS_REGION", region)

	// Create Session
	svc, err := EcrService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeRepositoriesPages(
		&ecr.DescribeRepositoriesInput{},
		func(page *ecr.DescribeRepositoriesOutput, isLast bool) bool {
			for _, repository := range page.Repositories {
				d.StreamListItem(ctx, repository)

			}
			return !isLast
		},
	)

	return nil, err
}

////  HYDRATE FUNCTIONS

func getAwsECRRepositories(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsECRRepositories")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	var name string
	if h.Item != nil {
		name = *h.Item.(*ecr.Repository).RepositoryName
	} else {
		name = d.KeyColumnQuals["repository_name"].GetStringValue()
	}

	// Create Session
	svc, err := EcrService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []*string{aws.String(name)},
	}

	// Get call
	data, err := svc.DescribeRepositories(params)
	if err != nil {
		logger.Debug("getAwsECRRepositories", "ERROR", err)
		return nil, err
	}
	return data.Repositories[0], nil
}

func getAwsECRRepositoryTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsECRRepositoryTags")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	resourceArn := h.Item.(*ecr.Repository).RepositoryArn

	// Create Session
	svc, err := EcrService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ecr.ListTagsForResourceInput{
		ResourceArn: resourceArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Debug("getAwsECRRepositoryTags", "ERROR", err)
		return nil, err
	}
	return op, nil
}

func getRepositoryPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getRepositoryPolicy")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	repositoryName := h.Item.(*ecr.Repository).RepositoryName

	// Create Session
	svc, err := EcrService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ecr.GetRepositoryPolicyInput{
		RepositoryName: repositoryName,
	}

	// Get call
	op, err := svc.GetRepositoryPolicy(params)
	if err != nil {
		logger.Debug("getRepositoryPolicy", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getAwsECRDescribeImages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsECRDescribeImages")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	repositoryName := h.Item.(*ecr.Repository).RepositoryName

	// Create Session
	svc, err := EcrService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ecr.DescribeImagesInput{
		RepositoryName: repositoryName,
	}

	// Get call
	op, err := svc.DescribeImages(params)
	if err != nil {
		logger.Debug("getAwsECRDescribeImages", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getLifecyclePolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getLifecyclePolicy")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	repositoryName := h.Item.(*ecr.Repository).RepositoryName

	// Create Session
	svc, err := EcrService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ecr.GetLifecyclePolicyInput{
		RepositoryName: repositoryName,
	}

	// Get call
	op, err := svc.GetLifecyclePolicy(params)
	if err != nil {
		logger.Debug("getLifecyclePolicy", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func ecrTagListToTurbotTags(ctx context.Context, d *transform.TransformData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ecrTagListToTurbotTags")
	tags := h.Item.(*ecr.ListTagsForResourceOutput)

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
