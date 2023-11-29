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

	"github.com/turbot/go-kit/helpers"
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
			ParentHydrate: listAwsEcrRepositories,
			Hydrate:       listAwsEcrImageScanFindings,
			Tags:          map[string]string{"service": "ecr", "action": "DescribeImageScanFindings"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"RepositoryNotFoundException", "ImageNotFoundException", "ScanNotFoundException"}),
			},
			// Ideally image_tag and image_digest could both be used as optional
			// key columns, but the query planner only works with required key
			// columns when there are multiple. We chose image_tag instead of
			// image_digest as it's more common/friendly to use.
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_name", Require: plugin.Optional},
				{Name: "image_tag", Require: plugin.Required},
				// {Name: "account_id", Require: plugin.Required},
				// {Name: "region", Require: plugin.Required},
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

	repository := h.Item.(types.Repository)

	repoName := d.EqualsQuals["repository_name"].GetStringValue()
	plugin.Logger(ctx).Error("Repositoy Name =====>> ", repoName)
	if repoName != "" {
		if repoName != *repository.RepositoryName {
			return nil, nil
		}
	}

	// In the case of an aggregator connection, the API is called multiple times with given repository names, regardless of the account IDs where the repository is available. This behavior is by design, as Steampipe iterates the API call per connection, irrespective of the given repository name. However, this approach consumes time when returning the results. To prevent unnecessary API calls for aggregator connections, we need to implement this check.

	// For example, let's consider an aggregator connection 'aws_all' that aggregates two connections, 'acc_a' and 'acc_b'.
	// Suppose a repository named 'test_repo' is available in 'acc_a' but not in 'acc_b'.
	// In this scenario, the API call should not occur for the connection 'acc_b'.
	// // plugin.Logger(ctx).Error("Account ID =====>> ", d.EqualsQualString("account_id"))
	// if d.EqualsQualString("account_id") != "" {
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_image_scan_finding.listAwsEcrImageScanFindings", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	plugin.Logger(ctx).Info("Okkk  =====>> Account ID: ", commonColumnData.AccountId, "Registry ID: ", *repository.RegistryId)
	if commonColumnData.AccountId != *repository.RegistryId {
		return nil, nil
	}
	// }

	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		return nil, err
	}

	imageTag := d.EqualsQuals["image_tag"]
	plugin.Logger(ctx).Trace("aws_ecr_image_scan_finding.listAwsEcrImageScanFindings", "repositoryName", d.EqualsQualString("repository_name"), "imageTag", imageTag)

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
		RepositoryName: repository.RepositoryName,
		ImageId: &types.ImageIdentifier{
			ImageTag: aws.String(imageTag.GetStringValue()),
		},
	}

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
		if err != nil {
			plugin.Logger(ctx).Error("aws_ecr_image_scan_finding.listAwsEcrImageScanFindings", "api_error", err)
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if helpers.StringSliceContains([]string{"RepositoryNotFoundException", "ImageNotFoundException", "ScanNotFoundException"}, ae.ErrorCode()) {
					return nil, nil
				}
			}
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
