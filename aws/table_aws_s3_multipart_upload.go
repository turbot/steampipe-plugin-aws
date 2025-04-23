package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAwsS3MultipartUpload(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_multipart_upload",
		Description: "AWS S3 Multipart Upload",
		List: &plugin.ListConfig{
			Hydrate: listS3MultipartUploads,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bucket_name", Require: plugin.Required, CacheMatch: query_cache.CacheMatchExact},
			},
			Tags: map[string]string{"service": "s3", "action": "ListMultipartUploads"},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "bucket_name",
				Description: "The name of the bucket to which the multipart upload was initiated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("bucket_name"),
			},
			{
				Name:        "key",
				Description: "The object key for which the multipart upload was initiated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "upload_id",
				Description: "Upload ID that identifies the multipart upload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "initiated",
				Description: "Date and time at which the multipart upload was initiated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "storage_class",
				Description: "The class of storage used to store the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "initiator_id",
				Description: "The ID of the initiator.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Initiator.ID"),
			},
			{
				Name:        "initiator_display_name",
				Description: "Display name of the initiator.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Initiator.DisplayName"),
			},
			{
				Name:        "owner_id",
				Description: "The ID of the owner.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.ID"),
			},
			{
				Name:        "owner_display_name",
				Description: "Display name of the owner.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.DisplayName"),
			},

			// Steampipe standard columns
			{
				Name:        "region",
				Description: "The AWS Region in which the object is located.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBucketRegionForObjects,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UploadId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listS3MultipartUploads(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucketName := d.EqualsQuals["bucket_name"].GetStringValue()

	// Bucket location will be nil if getBucketLocationForObjects returned an error but
	// was ignored through ignore_error_codes config arg
	bucketRegion, err := doGetBucketRegion(ctx, d, h, bucketName)
	if err != nil {
		return nil, err
	} else if bucketRegion == "" {
		return nil, nil
	}

	// Create Session
	svc, err := S3Client(ctx, d, bucketRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_multipart_upload.listS3MultipartUploads", "connection_error", err)
		return nil, err
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &s3.ListMultipartUploadsInput{
		Bucket:     aws.String(bucketName),
		MaxUploads: aws.Int32(maxLimit),
	}

	paginator := s3.NewListMultipartUploadsPaginator(svc, input, func(o *s3.ListMultipartUploadsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3_multipart_upload.listS3MultipartUploads", "api_error", err)
			return nil, err
		}

		for _, upload := range output.Uploads {
			d.StreamListItem(ctx, upload)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
