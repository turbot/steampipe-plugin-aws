package aws

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/smithy-go"

	ecrv1 "github.com/aws/aws-sdk-go/service/ecr"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcrImageScanFinding(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecr_image_scan_finding",
		Description: "AWS ECR Image Scan Finding",
		List: &plugin.ListConfig{
			Hydrate: listAwsEcrImageScanFindings,
			Tags:    map[string]string{"service": "ecr", "action": "DescribeImageScanFindings"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"RepositoryNotFoundException", "ImageNotFoundException", "ScanNotFoundException"}),
			},
			// Ideally image_tag and image_digest could both be used as optional
			// key columns, but the query planner only works with required key
			// columns when there are multiple. We chose image_tag instead of
			// image_digest as it's more common/friendly to use.
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_name", Require: plugin.Required},
				{Name: "image_tag", Require: plugin.AnyOf},
				{Name: "image_digest", Require: plugin.AnyOf},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ecrv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "repository_name",
				Description: "The name of the repository.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("repository_name"),
			},
			{
				Name:        "image_tag",
				Description: "The image tag",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_digest",
				Description: "The image digest",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name associated with the finding, usually a CVE number.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageScanFinding.Name"),
			},
			{
				Name:        "severity",
				Description: "The finding severity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageScanFinding.Severity"),
			},
			{
				Name:        "attributes",
				Description: "A collection of attributes of the host from which the finding is generated.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ImageScanFinding.Attributes"),
			},
			{
				Name:        "description",
				Description: "The description of the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageScanFinding.Description"),
			},
			{
				Name:        "uri",
				Description: "A link containing additional details about the security vulnerability.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageScanFinding.Uri"),
			},
			{
				Name:        "image_scan_status",
				Description: "The current state of the scan",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageScanStatus.Status"),
			},
			{
				Name:        "image_scan_status_description",
				Description: "The description of the image scan status.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageScanStatus.Description"),
			},
			{
				Name:        "image_scan_completed_at",
				Description: "The date and time, in JavaScript date format, when the repository was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "vulnerability_source_updated_at",
				Description: "The date and time, in JavaScript date format, when the repository was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageScanFinding.Name"),
			},
		}),
	}
}

// // LIST FUNCTION
func listAwsEcrImageScanFindings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		return nil, err
	}

	imageTag := d.EqualsQuals["image_tag"]
	imageDigest := d.EqualsQuals["image_digest"]
	repositoryName := d.EqualsQuals["repository_name"]

	if imageTag == nil && imageDigest == nil {
		return nil, errors.New("image_tag or image_digest must be provided")
	}
	
	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &ecr.DescribeImageScanFindingsInput{
		MaxResults:     aws.Int32(maxLimit),
		RepositoryName: aws.String(repositoryName.GetStringValue()),
	}

	imageInfo := &types.ImageIdentifier{
		ImageTag: aws.String(imageTag.GetStringValue()),
	}

	// Ideally, both image_tag and image_digest could be used.
	// However, they cannot be passed together simultaneously.
	// 1. If ImageTag is provided, it takes precedence and is used as the input parameter.
	// 2. If both ImageTag and ImageDigest are provided, ImageTag will be prioritized to keep the existing table behavior unchanged.
	// 3. If only ImageDigest is provided, the ImageDigest value will be used as the input parameter.
	if imageTag != nil {
		imageInfo.ImageTag = aws.String(imageTag.GetStringValue())
	}
	if imageTag != nil && imageDigest != nil {
		imageInfo.ImageTag = aws.String(imageTag.GetStringValue())
	}
	if imageTag == nil && imageDigest != nil {
		imageInfo.ImageDigest = aws.String(imageDigest.GetStringValue())
		imageInfo.ImageTag = nil
	}

	input.ImageId = imageInfo

	paginator := ecr.NewDescribeImageScanFindingsPaginator(svc, input, func(o *ecr.DescribeImageScanFindingsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	type ImageScanFindingsOutput struct {
		ImageDigest                  *string
		ImageScanCompletedAt         *time.Time
		ImageScanFinding             types.ImageScanFinding
		ImageScanStatus              types.ImageScanStatus
		ImageTag                     *string
		VulnerabilitySourceUpdatedAt *time.Time
	}

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		// In the case of parent hydrate use, the ignore error config in the list config is not functioning, so we need to handle the error here.
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if ae.ErrorCode() == "ScanNotFoundException" || ae.ErrorCode() == "RepositoryNotFoundException" || ae.ErrorCode() == "ImageNotFoundException" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_ecr_image_scan_finding.listAwsEcrImageScanFindings", "api_error", err)
			return nil, err
		}

		// If the scan is in progress and no findings are available yet, ImageScanFindings is nil
		if output.ImageScanFindings == nil {
			return nil, nil
		}

		for _, finding := range output.ImageScanFindings.Findings {
			result := &ImageScanFindingsOutput{
				ImageDigest:      output.ImageId.ImageDigest,
				ImageScanFinding: finding,
				ImageScanStatus:  *output.ImageScanStatus,
				ImageTag:         output.ImageId.ImageTag,
			}
			if output.ImageScanFindings.ImageScanCompletedAt != nil {
				result.ImageScanCompletedAt = output.ImageScanFindings.ImageScanCompletedAt
			}
			if output.ImageScanFindings.VulnerabilitySourceUpdatedAt != nil {
				result.VulnerabilitySourceUpdatedAt = output.ImageScanFindings.VulnerabilitySourceUpdatedAt
			}
			d.StreamListItem(ctx, result)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}
