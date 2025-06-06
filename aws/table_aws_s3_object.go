package aws

import (
	"context"
	"encoding/base64"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAwsS3Object(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_object",
		Description: "List AWS S3 Objects in S3 buckets by bucket name.",
		List: &plugin.ListConfig{
			Hydrate: listS3Objects,
			Tags:    map[string]string{"service": "s3", "action": "ListObjectsV2"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bucket_name", Require: plugin.Required, CacheMatch: query_cache.CacheMatchExact},
				{Name: "prefix", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},

				// If you encrypt an object by using server-side encryption with customer-provided
				// encryption keys (SSE-C) when you store the object in Amazon S3, then when you
				// GET the object, you must use the following query parameter:
				{Name: "sse_customer_algorithm", Require: plugin.Optional},
				{Name: "sse_customer_key", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "sse_customer_key_md5", Require: plugin.Optional},
			},
		},

		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getBucketRegionForObjects,
				Tags: map[string]string{"service": "s3", "action": "HTTPHeadBucket"},
			},
			{
				Func:    getS3Object,
				Depends: []plugin.HydrateFunc{getBucketRegionForObjects},
				Tags:    map[string]string{"service": "s3", "action": "GetObject"},
			},
			{
				Func:    headS3Object,
				Depends: []plugin.HydrateFunc{getBucketRegionForObjects},
				Tags:    map[string]string{"service": "s3", "action": "HeadObject"},
			},
			{
				Func:    getS3ObjectAttributes,
				Depends: []plugin.HydrateFunc{getBucketRegionForObjects},
				Tags:    map[string]string{"service": "s3", "action": "GetObjectAttributes"},
			},
			{
				Func:    getS3ObjectACL,
				Depends: []plugin.HydrateFunc{getBucketRegionForObjects},
				Tags:    map[string]string{"service": "s3", "action": "GetObjectAcl"},
			},
			{
				Func:    getS3ObjectTagging,
				Depends: []plugin.HydrateFunc{getBucketRegionForObjects},
				Tags:    map[string]string{"service": "s3", "action": "GetObjectTagging"},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "key",
				Description: "The name that you assign to an object. You use the object key to retrieve the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the AWS S3 Object.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getObjectARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "bucket_name",
				Description: "The name of the container bucket of this object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("bucket_name"),
			},
			{
				Name:        "last_modified",
				Description: "Last modified time of the object.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "storage_class",
				Description: "The class of storage used to store the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_id",
				Description: "The version ID of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VersionId"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "accept_ranges",
				Description: "Indicates that a range of bytes was specified.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AcceptRanges"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "body",
				Description: "The raw bytes of the object data as a string. If the bytes entirely consists of valid UTF8 runes, an UTF8 is sent otherwise the bas64 encoding of the bytes is sent.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue().Transform(readObjectBody),
				Hydrate:     getS3Object,
			},
			{
				Name:        "bucket_key_enabled",
				Description: "Indicates whether the object uses an S3 Bucket Key for server-side encryption with Amazon Web Services KMS (SSE-KMS)",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("BucketKeyEnabled"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "cache_control",
				Description: "Specifies caching behavior along the request/reply chain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CacheControl"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "checksum_crc32",
				Description: "The base64-encoded, 32-bit CRC32 checksum of the object. This will only be present if it was uploaded with the object. With multipart uploads, this may not be a checksum value of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ChecksumCRC32"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "checksum_crc32c",
				Description: "The base64-encoded, 32-bit CRC32C checksum of the object. This will only be present if it was uploaded with the object. With multipart uploads, this may not be a checksum value of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ChecksumCRC32C"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "checksum_sha1",
				Description: "The base64-encoded, 160-bit SHA-1 digest of the object. This will only be present if it was uploaded with the object. With multipart uploads, this may not be a checksum value of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ChecksumSHA1"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "checksum_sha256",
				Description: "The base64-encoded, 256-bit SHA-256 digest of the object. This will only be present if it was uploaded with the object. With multipart uploads, this may not be a checksum value of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ChecksumSHA256"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "content_disposition",
				Description: "Specifies presentational information for the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContentDisposition"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "content_encoding",
				Description: "Specifies what content encodings have been applied to the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContentEncoding"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "content_language",
				Description: "The language the content is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContentLanguage"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "content_length",
				Description: "Size of the body in bytes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContentLength"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "content_range",
				Description: "The portion of the object returned in the response.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContentRange"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "content_type",
				Description: "A standard MIME type describing the format of the object data.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContentType"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "delete_marker",
				Description: "Specifies whether the object retrieved was (true) or was not (false) a delete marker.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DeleteMarker"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "expiration",
				Description: "If the object expiration is configured (see PUT Bucket lifecycle), the response includes this header. It includes the expiry-date and rule-id key-value pairs providing object expiration information. The value of the rule-id is URL-encoded.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Expiration"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "expires",
				Description: "The date and time at which the object is no longer cacheable.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Expires"),
			},
			{
				Name:        "etag",
				Description: "The entity tag of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "object_lock_legal_hold_status",
				Description: "Like a retention period, a legal hold prevents an object version from being overwritten or deleted. A legal hold remains in effect until removed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ObjectLockLegalHoldStatus"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "object_lock_mode",
				Description: "The Object Lock mode currently in place for this object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ObjectLockMode"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "object_lock_retain_until_date",
				Description: "The date and time when this object's Object Lock will expire.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ObjectLockRetainUntilDate"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "prefix",
				Description: "The prefix of the key of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("prefix"),
			},
			{
				Name:        "replication_status",
				Description: "Amazon S3 can return this if your request involves a bucket that is either a source or destination in a replication rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReplicationStatus"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "request_charged",
				Description: "If present, indicates that the requester was successfully charged for the request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RequestCharged"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "restore",
				Description: "Provides information about object restoration action and expiration time of the restored object copy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Restore"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "server_side_encryption",
				Description: "The server-side encryption algorithm used when storing this object in Amazon S3.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerSideEncryption"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "size",
				Description: "Size in bytes of the object.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "sse_customer_algorithm",
				Description: "If server-side encryption with a customer-provided encryption key was requested, the response will include this header confirming the encryption algorithm used.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSECustomerAlgorithm"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "sse_customer_key",
				Description: "Specifies the customer-provided encryption key for Amazon S3 to use in encrypting data. This value is used to store the object and then it is discarded; Amazon S3 does not store the encryption key. The key must be appropriate for use with the algorithm specified in the x-amz-server-side-encryption-customer-algorithm header.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("sse_customer_key"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "sse_customer_key_md5",
				Description: "If server-side encryption with a customer-provided encryption key was requested, the response will include this header to provide round-trip message integrity verification of the customer-provided encryption key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSECustomerKeyMD5"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "sse_kms_key_id",
				Description: "If present, specifies the ID of the Amazon Web Services Key Management Service(Amazon Web Services KMS) symmetric customer managed key that was used for the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSEKMSKeyId"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "tag_count",
				Description: "The number of tags, if any, on the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TagCount"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "website_redirection_location",
				Description: "If the bucket is configured as a website, redirects requests for this object  to another object in the same bucket or to an external URL.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WebsiteRedirectLocation"),
				Hydrate:     headS3Object,
			},
			{
				Name:        "acl",
				Description: "ACLs define which AWS accounts or groups are granted access along with the type of access.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
				Hydrate:     getS3ObjectACL,
			},
			{
				Name:        "checksum",
				Description: "The checksum or digest of the object.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Checksum"),
				Hydrate:     getS3ObjectAttributes,
			},
			{
				Name:        "metadata",
				Description: "A map of metadata to store with the object in S3.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Metadata"),
				Hydrate:     getS3Object,
			},
			{
				Name:        "object_parts",
				Description: "A collection of parts associated with a multipart upload.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ObjectParts"),
				Hydrate:     getS3ObjectAttributes,
			},
			{
				Name:        "owner",
				Description: "The owner of the object.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the object.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TagSet"),
				Hydrate:     getS3ObjectTagging,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TagSet").Transform(handleS3TagsToTurbotTags),
				Hydrate:     getS3ObjectTagging,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Key"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getObjectARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
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

func getBucketRegionForObjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucketName := d.EqualsQuals["bucket_name"].GetStringValue()

	return doGetBucketRegion(ctx, d, h, bucketName)
}

// We cannot define the Hydrate Config for the List Hydrate call with hydrate dependencies.
// Therefore, we need to make an explicit hydrate call to get the bucket location.
func listS3Objects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucketName := d.EqualsQuals["bucket_name"].GetStringValue()

	// Bucket location will be nil if getBucketLocationForObjects returned an error but
	// was ignored through ignore_error_codes config arg
	bucketRegion, err := doGetBucketRegion(ctx, d, h, bucketName)
	if err != nil {
		return nil, err
	} else if bucketRegion == "" {
		return nil, nil
	}

	svc, err := S3Client(ctx, d, bucketRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object.listS3Objects", "get_client_error", err)
		return nil, err
	}

	// default supported max value is 1000 by ListObjectsV2
	maxItems := int32(1000)

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	input := &s3.ListObjectsV2Input{
		Bucket:     aws.String(bucketName),
		MaxKeys:    aws.Int32(maxItems),
		FetchOwner: aws.Bool(true),
	}

	equalQuals := d.EqualsQuals
	if equalQuals["prefix"] != nil {
		if equalQuals["prefix"].GetStringValue() != "" {
			input.Prefix = aws.String(equalQuals["prefix"].GetStringValue())
		}
	}

	// execute list call
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		objects, err := svc.ListObjectsV2(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3_object.ListObjectsV2", "api_error", err)
			return nil, err
		}

		for _, object := range objects.Contents {
			d.StreamListItem(ctx, object)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		input.ContinuationToken = objects.NextContinuationToken
		if objects.NextContinuationToken == nil {
			break
		}
	}

	return nil, err
}

func getS3Object(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucketName := d.EqualsQuals["bucket_name"].GetStringValue()
	sseCusAlgorithm := d.EqualsQuals["sse_customer_algorithm"].GetStringValue()
	sseCusKeyMd5 := d.EqualsQuals["sse_customer_key_md5"].GetStringValue()
	sseCusKey := d.EqualsQuals["sse_customer_key"].GetStringValue()
	bucketRegion := ""

	// Bucket location will be nil if getBucketLocationForObjects returned an error but
	// was ignored through ignore_error_codes config arg
	res := h.HydrateResults["getBucketRegionForObjects"]
	if res != nil {
		bucketRegion = res.(string)
	}

	// Bucket region empty check
	if bucketRegion == "" {
		return nil, nil
	}

	// Create client
	svc, err := S3Client(ctx, d, bucketRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object.getS3Object", "client_error", err)
		return nil, err
	}

	key := h.Item.(types.Object).Key

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    key,
	}

	// If you encrypt an object by using server-side encryption with customer-provided
	// encryption keys (SSE-C) when you store the object in Amazon S3.
	if sseCusAlgorithm != "" {
		params.SSECustomerAlgorithm = &sseCusAlgorithm
	}
	if sseCusKey != "" {
		params.SSECustomerKey = &sseCusKey
	}
	if sseCusKeyMd5 != "" {
		params.SSECustomerKeyMD5 = &sseCusKeyMd5
	}

	object, err := svc.GetObject(ctx, params)
	if err != nil {
		// if the key is unavailable in the provided bucket
		if strings.Contains(err.Error(), "NoSuchKey") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_s3_object.getS3Object", "api_error", err)
		return nil, err
	}

	return object, nil
}

func headS3Object(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucketName := d.EqualsQuals["bucket_name"].GetStringValue()
	sseCusAlgorithm := d.EqualsQuals["sse_customer_algorithm"].GetStringValue()
	sseCusKeyMd5 := d.EqualsQuals["sse_customer_key_md5"].GetStringValue()
	sseCusKey := d.EqualsQuals["sse_customer_key"].GetStringValue()
	bucketRegion := ""

	// Bucket location will be nil if getBucketLocationForObjects returned an error but
	// was ignored through ignore_error_codes config arg
	res := h.HydrateResults["getBucketRegionForObjects"]
	if res != nil {
		bucketRegion = res.(string)
	}

	// Bucket region empty check
	if bucketRegion == "" {
		return nil, nil
	}

	// Create client
	svc, err := S3Client(ctx, d, bucketRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object.headS3Object", "client_error", err)
		return nil, err
	}

	key := h.Item.(types.Object).Key

	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    key,
	}

	// If you encrypt an object by using server-side encryption with customer-provided
	// encryption keys (SSE-C) when you store the object in Amazon S3.
	if sseCusAlgorithm != "" {
		params.SSECustomerAlgorithm = &sseCusAlgorithm
	}
	if sseCusKey != "" {
		params.SSECustomerKey = &sseCusKey
	}
	if sseCusKeyMd5 != "" {
		params.SSECustomerKeyMD5 = &sseCusKeyMd5
	}

	object, err := svc.HeadObject(ctx, params)
	if err != nil {
		// if the key is unavailable in the provided bucket
		if strings.Contains(err.Error(), "NoSuchKey") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_s3_object.headS3Object", "api_error", err)
		return nil, err
	}

	return object, nil
}

func getS3ObjectAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucketName := d.EqualsQuals["bucket_name"].GetStringValue()
	sseCusAlgorithm := d.EqualsQuals["sse_customer_algorithm"].GetStringValue()
	sseCusKeyMd5 := d.EqualsQuals["sse_customer_key_md5"].GetStringValue()
	sseCusKey := d.EqualsQuals["sse_customer_key"].GetStringValue()
	bucketRegion := ""

	// Bucket location will be nil if getBucketLocationForObjects returned an error but
	// was ignored through ignore_error_codes config arg
	res := h.HydrateResults["getBucketRegionForObjects"]
	if res != nil {
		bucketRegion = res.(string)
	}

	// Bucket region empty check
	if bucketRegion == "" {
		return nil, nil
	}

	// Create client
	svc, err := S3Client(ctx, d, bucketRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object.getS3ObjectAttributes", "client_error", err)
		return nil, err
	}

	key := h.Item.(types.Object).Key

	params := &s3.GetObjectAttributesInput{
		Bucket:           aws.String(bucketName),
		Key:              key,
		ObjectAttributes: []types.ObjectAttributes{types.ObjectAttributesChecksum, types.ObjectAttributesObjectParts},
	}

	// If you encrypt an object by using server-side encryption with customer-provided
	// encryption keys (SSE-C) when you store the object in Amazon S3.
	if sseCusAlgorithm != "" {
		params.SSECustomerAlgorithm = &sseCusAlgorithm
	}
	if sseCusKey != "" {
		params.SSECustomerKey = &sseCusKey
	}
	if sseCusKeyMd5 != "" {
		params.SSECustomerKeyMD5 = &sseCusKeyMd5
	}

	objectAttributes, err := svc.GetObjectAttributes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object.GetObjectAttributes", "api_error", err)
		return nil, err
	}

	return objectAttributes, nil
}

func getS3ObjectACL(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucketName := d.EqualsQuals["bucket_name"].GetStringValue()
	bucketRegion := ""

	// Bucket location will be nil if getBucketLocationForObjects returned an error but
	// was ignored through ignore_error_codes config arg
	res := h.HydrateResults["getBucketRegionForObjects"]
	if res != nil {
		bucketRegion = res.(string)
	}

	// Bucket region empty check
	if bucketRegion == "" {
		return nil, nil
	}

	object := h.Item.(types.Object)

	// GetObjectAcl is not supported by Amazon S3 on Outposts.
	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectAcl.html
	if isOutpostObject(string(object.StorageClass)) {
		return nil, nil
	}

	// Create client
	svc, err := S3Client(ctx, d, bucketRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object.getS3ObjectACL", "client_error", err)
		return nil, err
	}

	input := &s3.GetObjectAclInput{
		Bucket: aws.String(bucketName),
		Key:    object.Key,
	}

	objectAcl, err := svc.GetObjectAcl(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object.GetObjectAcl", "api_error", err)
		return nil, err
	}

	return objectAcl, nil
}

func getS3ObjectTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucketName := d.EqualsQuals["bucket_name"].GetStringValue()
	bucketRegion := ""

	// Bucket location will be nil if getBucketLocationForObjects returned an error but
	// was ignored through ignore_error_codes config arg
	res := h.HydrateResults["getBucketRegionForObjects"]
	if res != nil {
		bucketRegion = res.(string)
	}

	// Bucket region empty check
	if bucketRegion == "" {
		return nil, nil
	}

	object := h.Item.(types.Object)

	// Create client
	svc, err := S3Client(ctx, d, bucketRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object.getS3ObjectTagging", "client_error", err)
		return nil, err
	}

	input := &s3.GetObjectTaggingInput{
		Bucket: aws.String(bucketName),
		Key:    object.Key,
	}

	tags, err := svc.GetObjectTagging(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object.GetObjectTagging", "api_error", err)
		return nil, err
	}

	return tags, nil
}

func getObjectARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	object := h.Item.(types.Object)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object.getObjectARN", "get_common_columns_error", err)
		return nil, err
	}
	bucketName := d.EqualsQuals["bucket_name"].GetStringValue()
	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":s3:::" + bucketName + "/" + *object.Key

	return arn, nil
}

func isOutpostObject(storageClass string) bool {
	// S3 on Outposts provides a new storage class, OUTPOSTS
	// as in https://docs.aws.amazon.com/AmazonS3/latest/userguide/S3onOutposts.html
	return strings.EqualFold(storageClass, "OUTPOSTS")
}

//// TRANSFORM FUNCTIONS

func readObjectBody(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.Value.(*s3.GetObjectOutput)

	body, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	if utf8.Valid(body) {
		return string(body), nil
	}

	return base64.StdEncoding.EncodeToString(body), nil
}
