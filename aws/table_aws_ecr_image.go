package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcrImage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecr_image",
		Description: "AWS ECR Image",
		List: &plugin.ListConfig{
			ParentHydrate: listAwsEcrRepositories,
			Hydrate:       listAwsEcrImages,
			Tags:          map[string]string{"service": "ecr", "action": "DescribeImages"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_name", Require: plugin.Optional},
				{Name: "registry_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_API_ECR_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "repository_name",
				Description: "The name of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "artifact_media_type",
				Description: "The artifact media type of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_digest",
				Description: "The sha256 digest of the image manifest.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_uri",
				Description: "The URI for the image.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getImageURI,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "image_manifest_media_type",
				Description: "The media type of the image manifest.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_pushed_at",
				Description: "The date and time, expressed in standard JavaScript date format, at which the current image was pushed to the repository.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "image_size_in_bytes",
				Description: "The size, in bytes, of the image in the repository.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "last_recorded_pull_time",
				Description: "The date and time, expressed in standard JavaScript date format, when Amazon ECR recorded the last image pull.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "registry_id",
				Description: "The Amazon Web Services account ID associated with the registry to which this image belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_scan_findings_summary",
				Description: "A summary of the last completed image scan.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "image_scan_status",
				Description: "The current state of the scan.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "image_tags",
				Description: "The list of tags associated with this image.",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsEcrImages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	repositoryName := h.Item.(types.Repository).RepositoryName

	repoName := d.EqualsQuals["repository_name"].GetStringValue()

	if repoName != "" {
		if repoName != *repositoryName {
			return nil, nil
		}
	}

	// Limiting the results
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

	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_image.listAwsEcrImages", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ecr.DescribeImagesInput{
		RepositoryName: repositoryName,
		MaxResults:     aws.Int32(maxLimit),
	}

	if d.EqualsQuals["registry_id"].GetStringValue() != "" {
		params.RegistryId = aws.String(d.EqualsQuals["registry_id"].GetStringValue())
	}

	paginator := ecr.NewDescribeImagesPaginator(svc, params, func(o *ecr.DescribeImagesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ecr_image.listAwsEcrImages", "api_error", err)
			return nil, err
		}

		for _, items := range output.ImageDetails {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getImageURI(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	image := h.Item.(types.ImageDetail)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// AWS follows the below image URI format -
	// with tag - {aws_account_id}.dkr.ecr.{region}.amazonaws.com/{repository name}:{first tag}
	// without tag - {aws_account_id}.dkr.ecr.{region}.amazonaws.com/{repository name}@image_digest
	if len(image.ImageTags) == 0 {
		uri := commonColumnData.AccountId + ".dkr.ecr." + commonColumnData.Region + ".amazonaws.com/" + *image.RepositoryName + "@" + *image.ImageDigest

		return uri, nil
	}

	uri := commonColumnData.AccountId + ".dkr.ecr." + commonColumnData.Region + ".amazonaws.com/" + *image.RepositoryName + ":" + image.ImageTags[0]

	return uri, nil
}
