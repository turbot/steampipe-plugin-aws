package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"

	ecrv1 "github.com/aws/aws-sdk-go/service/ecr"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcrRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecr_repository",
		Description: "AWS ECR Repository",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("repository_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"RepositoryNotFoundException", "RepositoryPolicyNotFoundException", "LifecyclePolicyNotFoundException"}),
			},
			Hydrate: getAwsEcrRepositories,
			Tags:    map[string]string{"service": "ecr", "action": "DescribeRepositories"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsEcrRepositories,
			Tags:    map[string]string{"service": "ecr", "action": "DescribeRepositories"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "registry_id", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsEcrRepositories,
				Tags: map[string]string{"service": "ecr", "action": "DescribeRepositories"},
			},
			{
				Func: listAwsEcrRepositoryTags,
				Tags: map[string]string{"service": "ecr", "action": "ListTagsForResource"},
			},
			{
				Func: getAwsEcrRepositoryPolicy,
				Tags: map[string]string{"service": "ecr", "action": "GetRepositoryPolicy"},
			},
			{
				Func: getAwsEcrDescribeImages,
				Tags: map[string]string{"service": "ecr", "action": "DescribeImages"},
			},
			{
				Func: getAwsEcrDescribeImageScanningFindings,
				Tags: map[string]string{"service": "ecr", "action": "DescribeImageScanFindings"},
			},
			{
				Func: getAwsEcrRepositoryLifecyclePolicy,
				Tags: map[string]string{"service": "ecr", "action": "GetLifecyclePolicy"},
			},
			{
				Func: getAwsEcrRepositoryScanningConfiguration,
				Tags: map[string]string{"service": "ecr", "action": "BatchGetRepositoryScanningConfiguration"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ecrv1.EndpointsID),
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
				Name:        "image_tag_mutability",
				Description: "The tag mutability setting for the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_evaluated_at",
				Description: "The time stamp of the last time that the lifecycle policy was run.",
				Hydrate:     getAwsEcrRepositoryLifecyclePolicy,
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "max_results",
				Description: "The maximum number of repository results returned by DescribeRepositories.",
				Hydrate:     getAwsEcrRepositories,
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "encryption_configuration",
				Description: "The encryption configuration for the repository.",
				Hydrate:     getAwsEcrRepositories,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "image_details",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release, use the aws_ecr_image table instead. A list of ImageDetail objects that contain data about the image.",
				Hydrate:     getAwsEcrDescribeImages,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "repository_scanning_configuration",
				Description: "Gets the scanning configuration for one or more repositories.",
				Hydrate:     getAwsEcrRepositoryScanningConfiguration,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "image_scanning_configuration",
				Description: "The image scanning configuration for a repository.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "image_scanning_findings",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release, use the aws_ecr_image_scan_finding table instead. Scan findings for an image.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEcrDescribeImageScanningFindings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "lifecycle_policy",
				Description: "The JSON lifecycle policy text.",
				Hydrate:     getAwsEcrRepositoryLifecyclePolicy,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LifecyclePolicyText"),
			},
			{
				Name:        "policy",
				Description: "The JSON repository policy text associated with the repository.",
				Hydrate:     getAwsEcrRepositoryPolicy,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyText"),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Hydrate:     getAwsEcrRepositoryPolicy,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyText").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Repository.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsEcrRepositoryTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsEcrRepositoryTags,
				Transform:   transform.FromField("Tags").Transform(ecrTagListToTurbotTags),
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

func listAwsEcrRepositories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.listAwsEcrRepositories", "connection_error", err)
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

	input := &ecr.DescribeRepositoriesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	equalQuals := d.EqualsQuals
	if equalQuals["registry_id"] != nil {
		input.RegistryId = aws.String(equalQuals["registry_id"].GetStringValue())
	}

	paginator := ecr.NewDescribeRepositoriesPaginator(svc, input, func(o *ecr.DescribeRepositoriesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ecr_repository.listAwsEcrRepositories", "api_error", err)
			return nil, err
		}

		for _, items := range output.Repositories {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

////  HYDRATE FUNCTIONS

func getAwsEcrRepositories(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var name string
	if h.Item != nil {
		name = *h.Item.(types.Repository).RepositoryName
	} else {
		name = d.EqualsQuals["repository_name"].GetStringValue()
	}

	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrRepositories", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{name},
	}

	// Get call
	data, err := svc.DescribeRepositories(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrRepositories", "api_error", err)
		return nil, err
	}
	return data.Repositories[0], nil
}

func listAwsEcrRepositoryTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	resourceArn := h.Item.(types.Repository).RepositoryArn

	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.listAwsEcrRepositoryTags", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ecr.ListTagsForResourceInput{
		ResourceArn: resourceArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.listAwsEcrRepositoryTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getAwsEcrRepositoryPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsEcrRepositoryPolicy")

	repositoryName := h.Item.(types.Repository).RepositoryName

	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrRepositoryPolicy", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ecr.GetRepositoryPolicyInput{
		RepositoryName: repositoryName,
	}

	// Get call
	op, err := svc.GetRepositoryPolicy(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "RepositoryPolicyNotFoundException") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrRepositoryPolicy", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsEcrDescribeImages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	repositoryName := h.Item.(types.Repository).RepositoryName

	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrDescribeImages", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ecr.DescribeImagesInput{
		RepositoryName: repositoryName,
		MaxResults:     aws.Int32(100),
	}

	// Get call
	op, err := svc.DescribeImages(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrDescribeImages", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsEcrDescribeImageScanningFindings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	getAwsEcrDescribeImageDetails := plugin.HydrateFunc(getAwsEcrDescribeImages)
	imageDetails, err := getAwsEcrDescribeImageDetails(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrDescribeImageDetails", "api_error", err)
		return nil, err
	}
	images := imageDetails.(*ecr.DescribeImagesOutput)

	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrDescribeImageScanningFindings", "connection_error", err)
		return nil, err
	}

	// Build the params
	// As per doc the max result value can be between 1-1000 but as per testing it returns only 100 result per page
	params := &ecr.DescribeImageScanFindingsInput{
		MaxResults: aws.Int32(100),
	}

	var result []ecr.DescribeImageScanFindingsOutput

	for _, image := range images.ImageDetails {
		var scanningDetails *ecr.DescribeImageScanFindingsOutput

		params.RepositoryName = image.RepositoryName
		params.ImageId = &types.ImageIdentifier{
			ImageDigest: image.ImageDigest,
		}

		paginator := ecr.NewDescribeImageScanFindingsPaginator(svc, params, func(o *ecr.DescribeImageScanFindingsPaginatorOptions) {
			o.Limit = 100
			o.StopOnDuplicateToken = true
		})

		// List call
		for paginator.HasMorePages() {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			scan, err := paginator.NextPage(ctx)
			if err != nil {
				if strings.Contains(err.Error(), "ScanNotFoundException") {
					return result, nil
				}
				plugin.Logger(ctx).Error("aws_ecr_repository.DescribeImageScanFindingsPages", "api_error", err)
				return nil, err
			}
			if scanningDetails != nil {
				if *scanningDetails.ImageId.ImageDigest == *image.ImageDigest {
					if scanningDetails.ImageScanFindings.EnhancedFindings != nil {
						scanningDetails.ImageScanFindings.EnhancedFindings = append(scan.ImageScanFindings.EnhancedFindings, scan.ImageScanFindings.EnhancedFindings...)
					} else if scanningDetails.ImageScanFindings.Findings != nil {
						scanningDetails.ImageScanFindings.Findings = append(scanningDetails.ImageScanFindings.Findings, scan.ImageScanFindings.Findings...)
					}
				}
				for k, v := range scan.ImageScanFindings.FindingSeverityCounts {
					scanningDetails.ImageScanFindings.FindingSeverityCounts[k] = v + scanningDetails.ImageScanFindings.FindingSeverityCounts[k]
				}
			} else {
				scanningDetails = scan
			}
		}
		if scanningDetails != nil {
			result = append(result, *scanningDetails)
		}
	}

	return result, nil
}

func getAwsEcrRepositoryLifecyclePolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	repositoryName := h.Item.(types.Repository).RepositoryName

	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrRepositoryLifecyclePolicy", "connection_error", err)
		return nil, err
	}
	// Build the params
	params := &ecr.GetLifecyclePolicyInput{
		RepositoryName: repositoryName,
	}
	// Get call
	op, err := svc.GetLifecyclePolicy(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "LifecyclePolicyNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrRepositoryLifecyclePolicy", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getAwsEcrRepositoryScanningConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	repositoryName := h.Item.(types.Repository).RepositoryName

	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrRepositoryScanningConfiguration", "connection_error", err)
		return nil, err
	}
	// Build the params
	params := &ecr.BatchGetRepositoryScanningConfigurationInput{
		RepositoryNames: []string{*repositoryName},
	}
	// Get call
	op, err := svc.BatchGetRepositoryScanningConfiguration(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_repository.getAwsEcrRepositoryScanningConfiguration", "api_error", err)
		return nil, err
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func ecrTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*ecr.ListTagsForResourceOutput)

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
