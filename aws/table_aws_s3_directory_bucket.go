package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsS3DirectoryBucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_directory_bucket",
		Description: "AWS S3 Directory Bucket",
		List: &plugin.ListConfig{
			Hydrate: listS3DirectoryBuckets,
			Tags:    map[string]string{"service": "s3", "action": "ListDirectoryBuckets"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_S3_SERVICE_ID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getS3DirectoryBucketPolicy,
				Tags: map[string]string{"service": "s3", "action": "GetBucketPolicy"},
			},
			{
				Func: getS3DirectoryBucketLifecycle,
				Tags: map[string]string{"service": "s3", "action": "GetBucketLifecycleConfiguration"},
			},
			{
				Func: getS3DirectoryBucketEncryption,
				Tags: map[string]string{"service": "s3", "action": "GetBucketEncryption"},
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the directory bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the directory bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getS3DirectoryBucketArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "creation_date",
				Description: "The date when the directory bucket was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "policy",
				Description: "The resource IAM access document for the directory bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getS3DirectoryBucketPolicy,
				Transform:   transform.FromField("Policy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getS3DirectoryBucketPolicy,
				Transform:   transform.FromField("Policy").Transform(policyToCanonical),
			},
			{
				Name:        "lifecycle_rules",
				Description: "The lifecycle configuration information of the directory bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getS3DirectoryBucketLifecycle,
				Transform:   transform.FromField("Rules"),
			},
			{
				Name:        "server_side_encryption_configuration",
				Description: "The default encryption configuration for the directory bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getS3DirectoryBucketEncryption,
				Transform:   transform.FromField("ServerSideEncryptionConfiguration.Rules"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getS3DirectoryBucketArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listS3DirectoryBuckets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the region from the matrix
	region := d.EqualsQualString(matrixKeyRegion)

	// Create service
	svc, err := S3Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_directory_bucket.listS3DirectoryBuckets", "get_client_error", err, "region", region)
		return nil, err
	}

	input := &s3.ListDirectoryBucketsInput{}

	// Paginate through results
	paginator := s3.NewListDirectoryBucketsPaginator(svc, input, func(o *s3.ListDirectoryBucketsPaginatorOptions) {
		o.Limit = 1000
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// For unsupported region we are encountering Error: aws: operation error S3: ListDirectoryBuckets, https response error StatusCode: 0, RequestID: , HostID: , request send failed, Get "https://s3express-control.ap-southeast-3.amazonaws.com/?max-directory-buckets=1000&x-id=ListDirectoryBuckets": lookup s3express-control.ap-southeast-3.amazonaws.com on 192.168.31.1:53: no such host
			if strings.Contains(err.Error(), "no such host") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_s3_directory_bucket.listS3DirectoryBuckets", "api_error", err, "region", region)
			return nil, err
		}

		for _, bucket := range output.Buckets {
			d.StreamListItem(ctx, bucket)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getS3DirectoryBucketPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucket := h.Item.(types.Bucket)
	region := d.EqualsQualString(matrixKeyRegion)

	// Create service
	svc, err := S3Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_directory_bucket.getS3DirectoryBucketPolicy", "get_client_error", err, "region", region)
		return nil, err
	}

	params := &s3.GetBucketPolicyInput{
		Bucket: bucket.Name,
	}

	policy, err := svc.GetBucketPolicy(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NoSuchBucketPolicy" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_directory_bucket.getS3DirectoryBucketPolicy", "api_error", err, "bucket", *bucket.Name)
		return nil, err
	}

	return policy, nil
}

func getS3DirectoryBucketLifecycle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucket := h.Item.(types.Bucket)
	region := d.EqualsQualString(matrixKeyRegion)

	// Create service
	svc, err := S3Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_directory_bucket.getS3DirectoryBucketLifecycle", "get_client_error", err, "region", region)
		return nil, err
	}

	params := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: bucket.Name,
	}

	lifecycleConfiguration, err := svc.GetBucketLifecycleConfiguration(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NoSuchLifecycleConfiguration" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_directory_bucket.getS3DirectoryBucketLifecycle", "api_error", err, "bucket", *bucket.Name)
		return nil, err
	}

	return lifecycleConfiguration, nil
}

func getS3DirectoryBucketEncryption(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucket := h.Item.(types.Bucket)
	region := d.EqualsQualString(matrixKeyRegion)

	// Create service
	svc, err := S3Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_directory_bucket.getS3DirectoryBucketEncryption", "get_client_error", err, "region", region)
		return nil, err
	}

	params := &s3.GetBucketEncryptionInput{
		Bucket: bucket.Name,
	}

	encryptionConfiguration, err := svc.GetBucketEncryption(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ServerSideEncryptionConfigurationNotFoundError" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_directory_bucket.getS3DirectoryBucketEncryption", "api_error", err, "bucket", *bucket.Name)
		return nil, err
	}

	return encryptionConfiguration, nil
}

func getS3DirectoryBucketArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucket := h.Item.(types.Bucket)

	region := d.EqualsQualString(matrixKeyRegion)

	commonInfo, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonData := commonInfo.(*awsCommonColumnData)

	arn := fmt.Sprintf("arn:%s:s3express:%s:%s:bucket/%s", commonData.Partition, region, commonData.AccountId, *bucket.Name)

	return arn, nil
}
