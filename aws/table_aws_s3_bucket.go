package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsS3Bucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_bucket",
		Description: "AWS S3 Bucket",
		List: &plugin.ListConfig{
			Hydrate: listS3Buckets,
		},
		// Note: No Get for S3 buckets, since it must list all the buckets
		// anyway just to get the creation_date which is only available via the
		// list call.
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func:    getBucketIsPublic,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
			{
				Func:    getBucketVersioning,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
			{
				Func:    getBucketEncryption,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
			{
				Func:    getBucketPublicAccessBlock,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
			{
				Func:    getBucketACL,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
			{
				Func:    getBucketLifecycle,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
			{
				Func:    getBucketLogging,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
			{
				Func:    getBucketPolicy,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
			{
				Func:    getBucketReplication,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
			{
				Func:    getBucketTagging,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
			{
				Func:    getObjectLockConfiguration,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
			{
				Func:    getS3BucketEventNotificationConfigurations,
				Depends: []plugin.HydrateFunc{getBucketLocation},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The user friendly name of the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the AWS S3 Bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBucketARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "creation_date",
				Description: "The date and time when bucket was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "bucket_policy_is_public",
				Description: "The policy status for an Amazon S3 bucket, indicating whether the bucket is public.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Hydrate:     getBucketIsPublic,
				Transform:   transform.FromField("PolicyStatus.IsPublic"),
			},
			{
				Name:        "versioning_enabled",
				Description: "The versioning state of a bucket.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketVersioning,
				Transform:   transform.FromField("Status").Transform(handleNilString).Transform(transform.ToBool),
			},
			{
				Name:        "versioning_mfa_delete",
				Description: "The MFA Delete status of the versioning state.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketVersioning,
				Transform:   transform.FromField("MFADelete").Transform(handleNilString).Transform(transform.ToBool),
			},
			{
				Name:        "block_public_acls",
				Description: "Specifies whether Amazon S3 should block public access control lists (ACLs) for this bucket and objects in this bucket.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketPublicAccessBlock,
				Transform:   transform.FromField("BlockPublicAcls"),
			},
			{
				Name:        "block_public_policy",
				Description: "Specifies whether Amazon S3 should block public bucket policies for this bucket. If TRUE it causes Amazon S3 to reject calls to PUT Bucket policy if the specified bucket policy allows public access.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketPublicAccessBlock,
				Transform:   transform.FromField("BlockPublicPolicy"),
			},
			{
				Name:        "ignore_public_acls",
				Description: "Specifies whether Amazon S3 should ignore public ACLs for this bucket and objects in this bucket. Setting this element to TRUE causes Amazon S3 to ignore all public ACLs on this bucket and objects in this bucket.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketPublicAccessBlock,
				Transform:   transform.FromField("IgnorePublicAcls"),
			},
			{
				Name:        "restrict_public_buckets",
				Description: "Specifies whether Amazon S3 should restrict public bucket policies for this bucket. Setting this element to TRUE restricts access to this bucket to only AWS service principals and authorized users within this account if the bucket has a public policy.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketPublicAccessBlock,
				Transform:   transform.FromField("RestrictPublicBuckets"),
			},
			{
				Name:        "event_notification_configuration",
				Description: "A container for specifying the notification configuration of the bucket. If this element is empty, notifications are turned off for the bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getS3BucketEventNotificationConfigurations,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_side_encryption_configuration",
				Description: "The default encryption configuration for an Amazon S3 bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketEncryption,
				Transform:   transform.FromField("ServerSideEncryptionConfiguration"),
			},
			{
				Name:        "acl",
				Description: "The access control list (ACL) of a bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketACL,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "lifecycle_rules",
				Description: "The lifecycle configuration information of the bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketLifecycle,
				Transform:   transform.FromField("Rules"),
			},
			{
				Name:        "logging",
				Description: "The logging status of a bucket and the permissions users have to view and modify that status.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketLogging,
				Transform:   transform.FromField("LoggingEnabled"),
			},
			{
				Name:        "object_lock_configuration",
				Description: "The specified bucket's object lock configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getObjectLockConfiguration,
			},
			{
				Name:        "policy",
				Description: "The resource IAM access document for the bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketPolicy,
				Transform:   transform.FromField("Policy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketPolicy,
				Transform:   transform.FromField("Policy").Transform(policyToCanonical),
			},
			{
				Name:        "replication",
				Description: "The replication configuration of a bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketReplication,
				Transform:   transform.FromField("ReplicationConfiguration"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketTagging,
				Transform:   transform.FromField("TagSet"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketTagging,
				Transform:   transform.FromField("TagSet").Transform(handleS3TagsToTurbotTags),
			},
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
				Hydrate:     getBucketARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
			{
				Name:        "region",
				Description: "The AWS Region in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBucketLocation,
				Transform:   transform.FromField("LocationConstraint"),
			},
		}),
	}
}

func listS3Buckets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Unlike most services, S3 buckets are a global list. They can be retrieved
	// from any single region.  We must list buckets from the default region to
	// get the actual creation_time of the bucket, in all other regions the list
	// returns the time when the bucket was last modified. See
	// https://www.marksayson.com/blog/s3-bucket-creation-dates-s3-master-regions/
	defaultRegion, err := getLastResortRegion(ctx, d, h)
	if err != nil {
		return nil, err
	}
	svc, err := S3Client(ctx, d, defaultRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.listS3Buckets", "get_client_error", err, "defaultRegion", defaultRegion)
		return nil, err
	}

	// execute list call
	input := &s3.ListBucketsInput{}
	bucketsResult, err := svc.ListBuckets(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.listS3Buckets", "api_error", err, "defaultRegion", defaultRegion)
		return nil, err
	}

	for _, bucket := range bucketsResult.Buckets {
		d.StreamListItem(ctx, bucket)
		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getS3BucketEventNotificationConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	name := h.Item.(types.Bucket).Name
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getS3BucketEventNotificationConfigurations", "client_error", err)
		return nil, err
	}

	// Build param
	input := &s3.GetBucketNotificationConfigurationInput{Bucket: name}

	notificationDetails, err := svc.GetBucketNotificationConfiguration(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getS3BucketEventNotificationConfigurations", "api_error", err)
		return nil, err
	}

	if notificationDetails != nil {
		output := map[string]any{}
		output["EventBridgeConfiguration"] = notificationDetails.EventBridgeConfiguration
		output["LambdaFunctionConfigurations"] = notificationDetails.LambdaFunctionConfigurations
		output["QueueConfigurations"] = notificationDetails.QueueConfigurations
		output["TopicConfigurations"] = notificationDetails.TopicConfigurations
		return output, nil
	}

	return nil, nil
}

func getBucketLocation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucket := h.Item.(types.Bucket)

	// Unlike most services, S3 buckets are a global list. They can be retrieved
	// from any single region. It's best to use the client region of the user
	// (e.g. closest to them).
	clientRegion, err := getDefaultRegion(ctx, d, h)
	if err != nil {
		return nil, err
	}
	svc, err := S3Client(ctx, d, clientRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketLocation", "get_client_error", err, "clientRegion", clientRegion)
		return nil, err
	}

	params := &s3.GetBucketLocationInput{Bucket: bucket.Name}

	// Specifies the Region where the bucket resides. For a list of all the Amazon
	// S3 supported location constraints by Region, see Regions and Endpoints (https://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region).
	location, err := svc.GetBucketLocation(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketLocation", "bucket_name", *bucket.Name, "clientRegion", clientRegion, "api_error", err)
		return nil, err
	}

	if location != nil && location.LocationConstraint != "" {
		// Buckets in eu-west-1 created through the AWS CLI or other API driven methods can return a location of "EU",
		// so we need to convert back
		if location.LocationConstraint == "EU" {
			return &s3.GetBucketLocationOutput{
				LocationConstraint: "eu-west-1",
			}, nil
		}
		return location, nil
	}

	// Buckets in us-east-1 have a LocationConstraint of null
	return &s3.GetBucketLocationOutput{
		LocationConstraint: "us-east-1",
	}, nil
}

func getBucketIsPublic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	bucket := h.Item.(types.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketIsPublic", "client_error", err)
		return nil, err
	}

	params := &s3.GetBucketPolicyStatusInput{Bucket: bucket.Name}
	policyStatus, err := svc.GetBucketPolicyStatus(ctx, params)
	if err != nil {
		var a smithy.APIError
		if errors.As(err, &a) {
			if a.ErrorCode() == "NoSuchBucketPolicy" {
				return &s3.GetBucketPolicyStatusOutput{}, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketIsPublic", "api_error", err)
		return nil, err
	}

	return policyStatus, nil
}

func getBucketVersioning(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	bucket := h.Item.(types.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketVersioning", "client_error", err)
		return nil, err
	}

	params := &s3.GetBucketVersioningInput{Bucket: bucket.Name}
	versioning, err := svc.GetBucketVersioning(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketVersioning", "api_error", err)
		return nil, err
	}

	return versioning, nil
}

func getBucketEncryption(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	bucket := h.Item.(types.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketEncryption", "client_error", err)
		return nil, err
	}
	params := &s3.GetBucketEncryptionInput{
		Bucket: bucket.Name,
	}

	encryption, err := svc.GetBucketEncryption(ctx, params)
	if err != nil {
		var a smithy.APIError
		if errors.As(err, &a) {
			if a.ErrorCode() == "ServerSideEncryptionConfigurationNotFoundError" {
				return nil, nil
			}
		}
		return nil, err
	}

	return encryption, nil
}

func getBucketPublicAccessBlock(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	bucket := h.Item.(types.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketPublicAccessBlock", "client_error", err)
		return nil, err
	}

	params := &s3.GetPublicAccessBlockInput{Bucket: bucket.Name}
	defaultAccessBlock := &types.PublicAccessBlockConfiguration{
		BlockPublicAcls:       false,
		BlockPublicPolicy:     false,
		IgnorePublicAcls:      false,
		RestrictPublicBuckets: false,
	}

	accessBlock, err := svc.GetPublicAccessBlock(ctx, params)
	if err != nil {
		// If the GetPublicAccessBlock is called on buckets which were created before Public Access Block setting was
		// introduced, sometime it fails with error NoSuchPublicAccessBlockConfiguration
		var a smithy.APIError
		if errors.As(err, &a) {
			if a.ErrorCode() == "NoSuchPublicAccessBlockConfiguration" {
				return defaultAccessBlock, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketPublicAccessBlock", "api_error", err)
		return nil, err
	}

	return accessBlock.PublicAccessBlockConfiguration, nil
}

func getBucketACL(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	bucket := h.Item.(types.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketACL", "client_error", err)
		return nil, err
	}

	params := &s3.GetBucketAclInput{Bucket: bucket.Name}

	acl, err := svc.GetBucketAcl(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketACL", "api_error", err)
		return nil, err
	}

	if acl != nil {
		output := map[string]any{}
		output["Grants"] = acl.Grants
		output["Owner"] = acl.Owner
		return &output, nil
	}

	return nil, nil
}

func getBucketLifecycle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	bucket := h.Item.(types.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketLifecycle", "client_error", err)
		return nil, err
	}

	params := &s3.GetBucketLifecycleConfigurationInput{Bucket: bucket.Name}

	lifecycleConfiguration, err := svc.GetBucketLifecycleConfiguration(ctx, params)
	if err != nil {
		var a smithy.APIError
		if errors.As(err, &a) {
			if a.ErrorCode() == "NoSuchLifecycleConfiguration" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketLifecycle", "api_error", err)
		return nil, err
	}

	return lifecycleConfiguration, nil
}

func getBucketLogging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	bucket := h.Item.(types.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketLogging", "client_error", err)
		return nil, err
	}

	params := &s3.GetBucketLoggingInput{Bucket: bucket.Name}
	logging, err := svc.GetBucketLogging(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketLogging", "api_error", err)
		return nil, err
	}
	return logging, nil
}

func getBucketPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	bucket := h.Item.(types.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketPolicy", "client_error", err)
		return nil, err
	}
	params := &s3.GetBucketPolicyInput{
		Bucket: bucket.Name,
	}

	bucketPolicy, err := svc.GetBucketPolicy(ctx, params)
	if err != nil {
		var a smithy.APIError
		if errors.As(err, &a) {
			if a.ErrorCode() == "NoSuchBucketPolicy" {
				return &s3.GetBucketPolicyOutput{}, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketPolicy", "api_error", err)
		return nil, err
	}

	return bucketPolicy, nil
}

func getBucketReplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	bucket := h.Item.(types.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketReplication", "client_error", err)
		return nil, err
	}
	params := &s3.GetBucketReplicationInput{Bucket: bucket.Name}

	replication, err := svc.GetBucketReplication(ctx, params)
	if err != nil {
		var a smithy.APIError
		if errors.As(err, &a) {
			if a.ErrorCode() == "ReplicationConfigurationNotFoundError" {
				return &s3.GetBucketReplicationOutput{}, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketReplication", "api_error", err)
		return nil, err
	}

	return replication, nil
}

func getBucketTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	bucket := h.Item.(types.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketTagging", "client_error", err)
		return nil, err
	}

	params := &s3.GetBucketTaggingInput{Bucket: bucket.Name}

	bucketTags, _ := svc.GetBucketTagging(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketTagging", "api_error", err)
		return nil, err
	}

	return bucketTags, nil
}

func getBucketARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucket := h.Item.(types.Bucket)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getBucketARN", "get_common_columns_error", err)
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":s3:::" + *bucket.Name

	return arn, nil
}

func getObjectLockConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getBucketLocation"] == nil {
		return nil, nil
	}

	bucket := h.Item.(types.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create client
	svc, err := S3Client(ctx, d, string(location.LocationConstraint))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket.getObjectLockConfiguration", "client_error", err)
		return nil, err
	}

	params := &s3.GetObjectLockConfigurationInput{Bucket: bucket.Name}

	data, err := svc.GetObjectLockConfiguration(ctx, params)
	if err != nil {
		var a smithy.APIError
		if errors.As(err, &a) {
			if a.ErrorCode() == "ObjectLockConfigurationNotFoundError" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_bucket.getObjectLockConfiguration", "api_error", err)
		return nil, err
	}

	return data, nil
}

//// TRANSFORM FUNCTIONS

func handleS3TagsToTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.Value.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
