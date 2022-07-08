package aws

import (
	"context"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsS3Object(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_object",
		Description: "AWS S3 Object",
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

			// {Name: "content_body", Description: "TODO", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getS3ObjectContent).WithCache(), Transform: transform.FromMethod("ReadBody")},
			// {Name: "content_body_parsed", Description: "TODO", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getS3ObjectContent).WithCache(), Transform: transform.FromMethod("ReadBodyParsed")},
			{Name: "content_type", Description: "TODO", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getS3ObjectContent).WithCache()},
			{Name: "content_encoding", Description: "TODO", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getS3ObjectContent).WithCache()},

			{Name: "acl_grants", Description: "A list of ACL grants.", Type: proto.ColumnType_JSON, Transform: transform.FromField("Grants"), Hydrate: plugin.HydrateFunc(getS3ObjectACL).WithCache()},
			{Name: "acl_owner", Description: "The bucket owner's display name and ID.", Type: proto.ColumnType_JSON, Transform: transform.FromField("Owner"), Hydrate: plugin.HydrateFunc(getS3ObjectACL).WithCache()},

			{Name: "retention", Description: "TODO", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Hydrate: getS3ObjectRetention},
			{Name: "legal_hold", Description: "TODO", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: getS3ObjectLegalHold},
		}),
	}
}

func getS3ObjectContent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	s3Object := h.Item.(*s3ObjectRow)
	svc, err := S3Service(ctx, d, *s3Object.BucketRegion)
	if err != nil {
		return nil, err
	}
	input := &s3.GetObjectInput{
		Bucket: s3Object.BucketName,
		Key:    s3Object.Key,
	}

	content, err := svc.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return &s3ObjectContent{
		GetObjectOutput: *content,
		parentRow:       s3Object,
		contentReadLock: &sync.Mutex{},
	}, nil
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

func listS3Objects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listS3Objects")

	bucketName := d.KeyColumnQualString("bucket")

	defer func() {
		r := recover()
		if r != nil {
			plugin.Logger(ctx).Error("panic recover", r)
			plugin.Logger(ctx).Error("panic recover", string(debug.Stack()))
		}
	}()

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
