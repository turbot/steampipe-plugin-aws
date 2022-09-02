package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsEcrImage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecr_image",
		Description: "AWS ECR Image",
		List: &plugin.ListConfig{
			ParentHydrate: listAwsEcrRepositories,
			Hydrate:       listAwsEcrImages,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_name", Require: plugin.Optional},
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

	repositoryName := h.Item.(*ecr.Repository).RepositoryName

	repoName := d.KeyColumnQuals["repository_name"].GetStringValue()

	if repoName != "" {
		if repoName != *repositoryName {
			return nil, nil
		}
	}

	// Create Session
	svc, err := EcrService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_image.listAwsEcrImages", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ecr.DescribeImagesInput{
		RepositoryName: repositoryName,
		MaxResults:     aws.Int64(100),
	}

	if d.KeyColumnQuals["registry_id"].GetStringValue() != "" {
		params.RegistryId = aws.String(d.KeyColumnQuals["registry_id"].GetStringValue())
	}

	err = svc.DescribeImagesPages(
		params,
		func(page *ecr.DescribeImagesOutput, isLast bool) bool {
			for _, image := range page.ImageDetails {
				d.StreamListItem(ctx, image)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_image.listAwsEcrImages", "api_error", err)
		return nil, err
	}

	return nil, nil
}
