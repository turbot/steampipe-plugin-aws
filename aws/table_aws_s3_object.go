package aws

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsS3Object(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_object",
		Description: "List AWS S3 Objects in S3 buckets by bucket name.",
		List: &plugin.ListConfig{
			Hydrate: listS3Objects,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "bucket", Require: plugin.Required},
				{Name: "key", Require: plugin.Optional},
				{Name: "prefix", Require: plugin.Optional},
			},
		},
		Columns: awsDefaultColumns([]*plugin.Column{
			{Name: "key", Description: "The name that you assign to an object. You use the object key to retrieve the object.", Type: proto.ColumnType_STRING},
			{Name: "etag", Description: "The entity tag of the object.", Type: proto.ColumnType_STRING, Transform: transform.FromField("ETag")},
			{Name: "storage_class", Description: "The class of storage used to store the object.", Type: proto.ColumnType_STRING},
			{Name: "size", Description: "Size in bytes of the object.", Type: proto.ColumnType_INT},
			{Name: "last_modified", Description: "Creation date of the object.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "prefix", Description: "The prefix of the key of the object.", Type: proto.ColumnType_STRING},
			{Name: "bucket", Description: "The name of the container bucket of this object.", Type: proto.ColumnType_STRING, Transform: transform.FromQual("bucket")},
			{Name: "acl", Description: "ACLs define which AWS accounts or groups are granted access along with the type of access.", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Hydrate: getS3ObjectACL},
			{Name: "retention", Description: "A retention period protects an object version for a fixed amount of time.", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Hydrate: getS3ObjectRetention},
			{Name: "legal_hold", Description: "Like a retention period, a legal hold prevents an object version from being overwritten or deleted. A legal hold remains in effect until removed.", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: getS3ObjectLegalHold},
			{Name: "tags", Description: "The tag set associated with an object.", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: getS3ObjectTagSet},
			{Name: "torrent", Description: "Returns the Bencode of the torrent. You can get torrent only for objects that are less than 5 GB in size, and that are not encrypted using server-side encryption with a customer-provided encryption key.", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: getS3ObjectTorrent},
		}),
	}
}

func listS3Objects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listS3Objects")

	bucketName := d.KeyColumnQualString("bucket")

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
		Bucket:  &bucketName,
		MaxKeys: &limit,
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
			derivedPrefix := ""
			if strings.Contains(*object.Key, "/") {
				derivedPrefix = (*object.Key)[:strings.LastIndex(*object.Key, "/")]
				if len(d.KeyColumnQualString("prefix")) > 0 {
					p := d.KeyColumnQualString("prefix")
					derivedPrefix = p
				}
			}

			row := &s3ObjectRow{
				Object:              *object,
				Prefix:              &derivedPrefix,
				BucketName:          &bucketName,
				BucketRegion:        bucketLocation.LocationConstraint,
				bucketHasLockConfig: bucketHasLockConfig,
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

func getS3ObjectACL(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	s3Object := h.Item.(*s3ObjectRow)
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
	taggingOutput, err := svc.GetObjectTaggingWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return taggingOutput.TagSet, nil
}

func getS3ObjectLegalHold(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	s3Object := h.Item.(*s3ObjectRow)

	if !s3Object.bucketHasLockConfig {
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
