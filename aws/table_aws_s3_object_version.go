package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsS3ObjectVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_object_version",
		Description: "AWS S3 Object Version",
		List: &plugin.ListConfig{
			Hydrate: listS3ObjectVersions,
			Tags:    map[string]string{"service": "s3", "action": "ListObjectVersions"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bucket_name", Require: plugin.Required, CacheMatch: "exact"},
				{Name: "key", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getBucketRegionForObjects,
				Tags: map[string]string{"service": "s3", "action": "HTTPHeadBucket"},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "bucket_name",
				Description: "The name of the container bucket of this object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("bucket_name"),
			},
			{
				Name:        "key",
				Description: "The object key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_class",
				Description: "The class of storage used to store the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_id",
				Description: "The entity tag is an MD5 hash of that version of the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_display_name",
				Description: "Container for the display name of the owner.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.DisplayName"),
			},
			{
				Name:        "owner_id",
				Description: "The entity tag is an MD5 hash of that version of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.ID"),
			},
			{
				Name:        "etag",
				Description: "Version ID of an object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag"),
			},
			{
				Name:        "size",
				Description: "Size in bytes of the object.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "is_latest",
				Description: "Specifies whether the object is (true) or is not (false) the latest version of an object.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_modified",
				Description: "Date and time the object was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "checksum_algorithm",
				Description: "The algorithm that was used to create a checksum of the object.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VersionId"),
			},
			{
				Name:        "region",
				Description: "The AWS Region in which the object is located.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBucketRegionForObjects,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

func listS3ObjectVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucketName := d.EqualsQuals["bucket_name"].GetStringValue()
	bucketRegion := h.HydrateResults["getBucketRegionForObjects"].(string)

	svc, err := S3Client(ctx, d, bucketRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object_version.listS3ObjectVersions", "get_client_error", err)
		return nil, err
	}

	// default supported max value is 1000 by ListObjectVersions
	maxItems := int32(1000)

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	input := &s3.ListObjectVersionsInput{
		Bucket:  aws.String(bucketName),
		MaxKeys: maxItems,
	}

	if d.EqualsQualString("key") != "" {
		input.Prefix = aws.String(d.EqualsQualString("key"))
	}

	// execute list call
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		objects, err := svc.ListObjectVersions(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3_object_version.ListObjectVersions", "api_error", err)
			return nil, err
		}

		for _, version := range objects.Versions {
			d.StreamListItem(ctx, version)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		input.KeyMarker = objects.NextVersionIdMarker
		if objects.NextVersionIdMarker == nil {
			break
		}
	}

	return nil, err
}
