package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAwsS3ObjectData(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_object_data",
		Description: "List AWS S3 Objects in S3 buckets by bucket name.",
		List: &plugin.ListConfig{
			Hydrate:    getAWSS3ObjectData,
			KeyColumns: plugin.KeyColumnSlice{
				// {Name: "bucket", Require: plugin.Required},
				// {Name: "key", Require: plugin.Required},

				// {Name: "sse_customer_key", Require: plugin.Optional},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			// {
			// 	Name:        "key",
			// 	Description: "The name assigned to an object.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromQual("key"),
			// },
			// {
			// 	Name:        "bucket",
			// 	Description: "The name of the container bucket of the object.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromQual("bucket"),
			// },
			// {
			// 	Name:        "content_encoding",
			// 	Description: "Specifies what content encodings have been applied to the object.",
			// 	Type:        proto.ColumnType_STRING,
			// },
			// {
			// 	Name:        "content_type",
			// 	Description: "A standard MIME type describing the format of the object data.",
			// 	Type:        proto.ColumnType_STRING,
			// },
			// {
			// 	Name:        "sse_customer_key",
			// 	Description: "If server-side encryption is set on the object, use this to provide the 256-bit, base64-encoded encryption key which Amazon S3 will use to decrypt your data.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromQual("sse_customer_key"),
			// },
			// {
			// 	Name:        "sse_customer_algorithm",
			// 	Description: "If server-side encryption with a customer-provided encryption key was requested, the response will include this header confirming the encryption algorithm used.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromField("SSECustomerAlgorithm"),
			// },
			// {
			// 	Name:        "sse_kms_key_id",
			// 	Description: "If present, specifies the ID of the Amazon Web Services Key Management Service (Amazon Web Services KMS) symmetric customer managed key that was used for the object.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromField("SSEKMSKeyId"),
			// },
			// {
			// 	Name:        "data",
			// 	Description: "The raw bytes of the object as a string. An UTF8 encoded string is sent, if the bytes entirely consists of valid UTF8 runes, an UTF8 is sent otherwise the bas64 encoding of the bytes is sent.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromMethod("ReadBody"),
			// },

			// // steampipe fields
			// {
			// 	Name:        "title",
			// 	Description: resourceInterfaceDescription("title"),
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromQual("key"),
			// },
		}),
	}
}

func getAWSS3ObjectData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// bucketName := d.KeyColumnQualString("bucket")
	// key := d.KeyColumnQualString("key")

	// region, err := resolveBucketRegion(ctx, d, &bucketName)
	// if err != nil {
	// 	return nil, err
	// }

	// svc, err := S3Service(ctx, d, *region.LocationConstraint)
	// if err != nil {
	// 	return nil, err
	// }

	// input := &s3.GetObjectInput{
	// 	Bucket: &bucketName,
	// 	Key:    &key,
	// }

	// if len(d.KeyColumnQualString("sse_customer_key")) > 0 {
	// 	// Refer:
	// 	// https://docs.aws.amazon.com/AmazonS3/latest/userguide/ServerSideEncryptionCustomerKeys.html#specifying-s3-c-encryption
	// 	input.
	// 		// the 256-bit, base64-encoded encryption key
	// 		SetSSECustomerKey(d.KeyColumnQualString("sse_customer_key")).
	// 		// value must be "AES256"
	// 		SetSSECustomerAlgorithm("AES256")
	// }

	// output, err := svc.GetObjectWithContext(ctx, input)
	// if err != nil {
	// 	return nil, err
	// }

	// return &s3ObjectContent{
	// 	GetObjectOutput: *output,
	// }, nil
	return nil, nil
}
