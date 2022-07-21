package aws

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsS3Object(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_object",
		Description: "List AWS S3 Objects in S3 buckets by bucket name.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"bucket", "key"}),
			Hydrate:    getAWSS3Object,
		},
		List: &plugin.ListConfig{
			Hydrate:       listAWSS3Objects,
			ParentHydrate: listS3Buckets,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "bucket", Require: plugin.Optional},
				{Name: "key", Require: plugin.Optional},
				{Name: "prefix", Require: plugin.Optional},
			},
		},
		Columns: awsDefaultColumns([]*plugin.Column{
			{
				Name:        "key",
				Description: "The name that you assign to an object. You use the object key to retrieve the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "The entity tag of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag"),
			},
			{
				Name:        "storage_class",
				Description: "The class of storage used to store the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "Size in bytes of the object.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "last_modified",
				Description: "Creation date of the object.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "prefix",
				Description: "The prefix of the key of the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bucket",
				Description: "The name of the container bucket of this object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BucketName"),
			},
			{
				Name:        "owner",
				Description: "The owner of the object.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Owner"),
				Hydrate:     getS3ObjectACL,
			},
			{
				Name:        "acl",
				Description: "ACLs define which AWS accounts or groups are granted access along with the type of access.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
				Hydrate:     getS3ObjectACL,
			},
			{
				Name:        "retention",
				Description: "A retention period protects an object version for a fixed amount of time.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
				Hydrate:     getS3ObjectRetention,
			},
			{
				Name:        "legal_hold",
				Description: "Like a retention period, a legal hold prevents an object version from being overwritten or deleted. A legal hold remains in effect until removed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue(),
				Hydrate:     getS3ObjectLegalHold,
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TagSet").Transform(s3TagsToTurbotTags),
				Hydrate:     getS3ObjectTagSet,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the object.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TagSet"),
				Hydrate:     getS3ObjectTagSet,
			},
			{
				Name:        "torrent",
				Description: "Returns the Bencode of the torrent. You can get torrent only for objects that are less than 5 GB in size, and that are not encrypted using server-side encryption with a customer-provided encryption key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue(),
				Hydrate:     getS3ObjectTorrent,
			},

			{
				Name:        "checksum",
				Description: "The checksum or digest of the object.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Checksum"),
				Hydrate:     getS3ObjectAttributes,
			},
			{
				Name:        "parts",
				Description: "A collection of parts associated with a multipart upload.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ObjectParts"),
				Hydrate:     getS3ObjectAttributes,
			},
			{
				Name:        "delete_marker",
				Description: "Specifies whether the object retrieved was (true) or was not (false) a delete marker.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DeleteMarker"),
				Hydrate:     getS3ObjectAttributes,
			},
			{
				Name:        "content_type",
				Description: "A standard MIME type describing the format of the object data.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContentType"),
				Hydrate:     getS3ObjectContent,
			},
			{
				Name:        "bucket_key_enabled",
				Description: "Indicates whether the object uses an S3 Bucket Key for server-side encryption with Amazon Web Services KMS (SSE-KMS)",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("BucketKeyEnabled"),
				Hydrate:     getS3ObjectContent,
			},
			{
				Name:        "metadata",
				Description: "A map of metadata to store with the object in S3.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Metadata"),
				Hydrate:     getS3ObjectContent,
			},
			{
				Name:        "content_encoding",
				Description: "Specifies what content encodings have been applied to the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContentEncoding"),
				Hydrate:     getS3ObjectContent,
			},
			{
				Name:        "content_length",
				Description: "Size of the body in bytes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContentLength"),
				Hydrate:     getS3ObjectContent,
			},
			{
				Name:        "replication_status",
				Description: "Amazon S3 can return this if your request involves a bucket that is either a source or destination in a replication rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReplicationStatus"),
				Hydrate:     getS3ObjectContent,
			},
			{
				Name:        "restore",
				Description: "Provides information about object restoration action and expiration time of the restored object copy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Restore"),
				Hydrate:     getS3ObjectContent,
			},
			{
				Name:        "server_side_encryption",
				Description: "The server-side encryption algorithm used when storing this object in Amazon S3.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerSideEncryption"),
				Hydrate:     getS3ObjectContent,
			},
			{
				Name:        "website_redirection_location",
				Description: "If the bucket is configured as a website, redirects requests for this object  to another object in the same bucket or to an external URL.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WebsiteRedirectLocation"),
				Hydrate:     getS3ObjectContent,
			},

			// steampipe fields
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Key"),
			},
		}),
	}
}

func getAWSS3Object(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAWSS3Object")

	bucketName := d.KeyColumnQualString("bucket")
	key := d.KeyColumnQualString("key")
	// Create Session,
	bucketLocation, err := resolveBucketRegion(ctx, d, &bucketName)
	if err != nil {
		return nil, err
	}
	svc, err := S3Service(ctx, d, *bucketLocation.LocationConstraint)
	if err != nil {
		return nil, err
	}

	// we need to retain this, since a few fields in the objects will always be `nil` if this is `nil`
	_, err = svc.GetObjectLockConfigurationWithContext(ctx, &s3.GetObjectLockConfigurationInput{
		Bucket: &bucketName,
	})
	bucketHasLockConfig := true
	if err != nil {
		if strings.Contains(err.Error(), "ObjectLockConfigurationNotFoundError") {
			bucketHasLockConfig = false
		} else {
			return nil, err
		}
	}

	input := &s3.GetObjectAttributesInput{}
	input.
		SetBucket(bucketName).
		SetKey(key).
		SetObjectAttributes(aws.StringSlice(s3.ObjectAttributes_Values()))

	output, err := svc.GetObjectAttributesWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	row := &s3ObjectRow{
		BucketName:          &bucketName,
		BucketRegion:        bucketLocation.LocationConstraint,
		bucketHasLockConfig: bucketHasLockConfig,
		StorageClass:        output.StorageClass,
		Key:                 &key,
		ETag:                output.ETag,
		Size:                output.ObjectSize,
		LastModified:        output.LastModified,
		Prefix:              derivePrefix(ctx, key, d.KeyColumnQualString("prefix")),
	}

	return row, nil
}

func listAWSS3Objects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	plugin.Logger(ctx).Trace("listS3Objects")

	bucket := h.Item.(*s3.Bucket)
	bucketName := d.KeyColumnQualString("bucket")

	// if a bucket name was provided and this is not the same bucket from the parentHydrate, skip
	if len(bucketName) > 0 && !strings.EqualFold(bucketName, *bucket.Name) {
		return nil, nil
	}

	bucketName = *bucket.Name

	// Create Session,
	bucketLocation, err := resolveBucketRegion(ctx, d, &bucketName)
	if err != nil {
		return nil, err
	}
	svc, err := S3Service(ctx, d, *bucketLocation.LocationConstraint)
	if err != nil {
		return nil, err
	}

	// we need to retain this, since a few fields in the objects will always be `nil` if this is `nil`
	_, err = svc.GetObjectLockConfigurationWithContext(ctx, &s3.GetObjectLockConfigurationInput{
		Bucket: &bucketName,
	})
	bucketHasLockConfig := true
	if err != nil {
		if strings.Contains(err.Error(), "ObjectLockConfigurationNotFoundError") {
			bucketHasLockConfig = false
		} else {
			return nil, err
		}
	}

	limit := (int64(1000))
	if d.QueryContext.Limit != nil && limit > *d.QueryContext.Limit {
		limit = *d.QueryContext.Limit
	}

	input := &s3.ListObjectsV2Input{
		Bucket:     &bucketName,
		MaxKeys:    &limit,
		FetchOwner: aws.Bool(false),
	}

	if len(d.KeyColumnQualString("prefix")) > 0 {
		p := d.KeyColumnQualString("prefix")
		input.Prefix = &p
	}

	if len(d.KeyColumnQualString("key")) > 0 {
		p := d.KeyColumnQualString("key")
		// overwrite the prefix with the full key
		input.Prefix = &p
	}

	err = svc.ListObjectsV2PagesWithContext(ctx, input, func(objectList *s3.ListObjectsV2Output, b bool) bool {
		for _, object := range objectList.Contents {

			row := &s3ObjectRow{
				StorageClass:        object.StorageClass,
				Key:                 object.Key,
				ETag:                object.ETag,
				Size:                object.Size,
				LastModified:        object.LastModified,
				BucketName:          &bucketName,
				BucketRegion:        bucketLocation.LocationConstraint,
				bucketHasLockConfig: bucketHasLockConfig,
				Prefix:              derivePrefix(ctx, *object.Key, d.KeyColumnQualString("prefix")),
			}

			d.StreamListItem(ctx, row)

			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				// do not continue
				return false
			}
		}

		// go to the next page
		return true
	})

	return nil, err
}

func derivePrefix(ctx context.Context, key string, knownPrefix string) string {
	derivedPrefix := knownPrefix
	if strings.Contains(key, "/") && len(knownPrefix) == 0 {
		derivedPrefix = (key)[:strings.LastIndex(key, "/")]
	}
	return derivedPrefix
}

// Column hydrates

func getS3ObjectContent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	s3Object := h.Item.(*s3ObjectRow)

	svc, err := S3Service(ctx, d, *s3Object.BucketRegion)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Trace("fetching content for ", *s3Object.Key)

	input := &s3.GetObjectInput{
		Bucket: s3Object.BucketName,
		Key:    s3Object.Key,
	}

	output, err := svc.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return &s3ObjectContent{
		GetObjectOutput: *output,
	}, nil
}

func getS3ObjectAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	row := h.Item.(*s3ObjectRow)
	svc, err := S3Service(ctx, d, *row.BucketRegion)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Trace("fetching attributes for", *row.Key)

	input := &s3.GetObjectAttributesInput{
		Bucket:           row.BucketName,
		Key:              row.Key,
		ObjectAttributes: aws.StringSlice(s3.ObjectAttributes_Values()),
	}
	output, err := svc.GetObjectAttributesWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return &s3ObjectAttributes{
		GetObjectAttributesOutput: *output,
		parentRow:                 row,
	}, nil
}

func getS3ObjectACL(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	s3Object := h.Item.(*s3ObjectRow)

	if s3Object.isOutpostObject() {
		return nil, nil
	}

	svc, err := S3Service(ctx, d, *s3Object.BucketRegion)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Trace("fetching ACL for ", s3Object.Key)

	input := &s3.GetObjectAclInput{
		Bucket: s3Object.BucketName,
		Key:    s3Object.Key,
	}
	return svc.GetObjectAclWithContext(ctx, input)
}

func getS3ObjectTorrent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	s3Object := h.Item.(*s3ObjectRow)

	if s3Object.isOutpostObject() {
		return nil, nil
	}

	if !s3Object.bucketHasLockConfig {
		return nil, nil
	}

	svc, err := S3Service(ctx, d, *s3Object.BucketRegion)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Trace("fetching torrent for ", s3Object.Key)

	input := &s3.GetObjectTorrentInput{
		Bucket: s3Object.BucketName,
		Key:    s3Object.Key,
	}
	torrentOutput, err := svc.GetObjectTorrentWithContext(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("torrent bytes error", err)
		return nil, err
	}

	bodyBytes, err := io.ReadAll(torrentOutput.Body)
	if err != nil {
		plugin.Logger(ctx).Error("torrent bytes error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("torrent bytes", bodyBytes)

	return string(bodyBytes), nil
}

func getS3ObjectTagSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	s3Object := h.Item.(*s3ObjectRow)

	if !s3Object.bucketHasLockConfig {
		return nil, nil
	}

	svc, err := S3Service(ctx, d, *s3Object.BucketRegion)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Trace("fetching tag set for ", s3Object.Key)

	input := &s3.GetObjectTaggingInput{
		Bucket: s3Object.BucketName,
		Key:    s3Object.Key,
	}

	return svc.GetObjectTaggingWithContext(ctx, input)
}

func getS3ObjectLegalHold(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	s3Object := h.Item.(*s3ObjectRow)

	if !s3Object.bucketHasLockConfig {
		return nil, nil
	}

	if s3Object.isOutpostObject() {
		return nil, nil
	}

	svc, err := S3Service(ctx, d, *s3Object.BucketRegion)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Trace("fetching legal hold for ", s3Object.Key)

	input := &s3.GetObjectLegalHoldInput{
		Bucket: s3Object.BucketName,
		Key:    s3Object.Key,
	}
	legalHoldOutput, err := svc.GetObjectLegalHoldWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return legalHoldOutput.LegalHold.Status, nil
}

func getS3ObjectRetention(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	s3Object := h.Item.(*s3ObjectRow)

	if !s3Object.bucketHasLockConfig {
		return nil, nil
	}

	if s3Object.isOutpostObject() {
		return nil, nil
	}

	svc, err := S3Service(ctx, d, *s3Object.BucketRegion)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Trace("fetching Object Retention for ", s3Object.Key)

	input := &s3.GetObjectRetentionInput{
		Bucket: s3Object.BucketName,
		Key:    s3Object.Key,
	}
	retentionOutput, err := svc.GetObjectRetentionWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return retentionOutput.Retention, nil
}
