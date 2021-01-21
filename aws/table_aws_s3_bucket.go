package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsS3Bucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_bucket",
		Description: "AWS S3 Bucket",
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.SingleColumn("name"),
			ItemFromKey: bucketFromKey,
			Hydrate:     getS3Bucket,
		},
		List: &plugin.ListConfig{
			Hydrate: listS3Buckets,
		},
		HydrateDependencies: []plugin.HydrateDependencies{
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
		},
		Columns: awsS3Columns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The user friendly name of the bucket",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "The date and tiem when bucket was created.",
				Type:        proto.ColumnType_DATETIME,
			},
			{
				Name:        "bucket_policy_is_public",
				Description: "The policy status for an Amazon S3 bucket, indicating whether the bucket is public",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Hydrate:     getBucketIsPublic,
				Transform:   transform.FromField("PolicyStatus.IsPublic").Transform(handleNilString).Transform(transform.ToBool),
			},
			{
				Name:        "versioning_enabled",
				Description: "The versioning state of a bucket",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketVersioning,
				Transform:   transform.FromField("Status").Transform(handleNilString).Transform(transform.ToBool),
			},
			{
				Name:        "versioning_mfa_delete",
				Description: "The MFA Delete status of the versioning state",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketVersioning,
				Transform:   transform.FromField("MFADelete").Transform(handleNilString).Transform(transform.ToBool),
			},
			{
				Name:        "block_public_acls",
				Description: "Specifies whether Amazon S3 should block public access control lists (ACLs) for this bucket and objects in this bucket",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketPublicAccessBlock,
				Transform:   transform.FromField("BlockPublicAcls"),
			},
			{
				Name:        "block_public_policy",
				Description: "Specifies whether Amazon S3 should block public bucket policies for this bucket. If TRUE it causes Amazon S3 to reject calls to PUT Bucket policy if the specified bucket policy allows public access",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketPublicAccessBlock,
				Transform:   transform.FromField("BlockPublicPolicy"),
			},
			{
				Name:        "ignore_public_acls",
				Description: "Specifies whether Amazon S3 should ignore public ACLs for this bucket and objects in this bucket. Setting this element to TRUE causes Amazon S3 to ignore all public ACLs on this bucket and objects in this bucket",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketPublicAccessBlock,
				Transform:   transform.FromField("IgnorePublicAcls"),
			},
			{
				Name:        "restrict_public_buckets",
				Description: "Specifies whether Amazon S3 should restrict public bucket policies for this bucket. Setting this element to TRUE restricts access to this bucket to only AWS service principals and authorized users within this account if the bucket has a public policy",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketPublicAccessBlock,
				Transform:   transform.FromField("RestrictPublicBuckets"),
			},
			{
				Name:        "server_side_encryption_configuration",
				Description: "The default encryption configuration for an Amazon S3 bucket",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketEncryption,
				Transform:   transform.FromField("ServerSideEncryptionConfiguration"),
			},
			{
				Name:        "acl",
				Description: "The access control list (ACL) of a bucket",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketACL,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "lifecycle_rules",
				Description: "The lifecycle configuration information of the bucket",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketLifecycle,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "logging",
				Description: "The logging status of a bucket and the permissions users have to view and modify that status",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketLogging,
				Transform:   transform.FromField("LoggingEnabled"),
			},
			{
				Name:        "policy",
				Description: "The resource IAM access document for the bucket",
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
				Description: "The replication configuration of a bucket",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketReplication,
				Transform:   transform.FromField("ReplicationConfiguration"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to bucket",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketTagging,
				Transform:   transform.FromField("TagSet"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketTagging,
				Transform:   transform.FromField("TagSet").Transform(s3TagsToTurbotTags),
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
				Transform:   transform.From(s3NameToAkas),
			},
			{
				Name:        "region",
				Description: "The AWS Region in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBucketLocation,
				Transform:   transform.FromField("LocationConstraint"),
			},
		}),
	}
}

//// BUILD HYDRATE INPUT

func bucketFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &s3.Bucket{
		Name: &name,
	}
	return item, nil
}

//// LIST FUNCTION

func listS3Buckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listS3Buckets")
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &s3.ListBucketsInput{}
	bucketsResult, err := svc.ListBuckets(input)
	if err != nil {
		return nil, err
	}

	for _, bucket := range bucketsResult.Buckets {
		d.StreamListItem(ctx, bucket)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

// do not have a get call for s3 bucket.
// using list api call to create get function
func getS3Bucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listS3Buckets")
	bucket := h.Item.(*s3.Bucket)
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &s3.ListBucketsInput{}
	bucketsResult, err := svc.ListBuckets(input)
	if err != nil {
		return nil, err
	}

	for _, item := range bucketsResult.Buckets {
		if *item.Name == *bucket.Name {
			return item, nil
		}
	}

	return nil, err
}

func getBucketLocation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketLocation")
	bucket := h.Item.(*s3.Bucket)
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketLocationInput{
		Bucket: bucket.Name,
	}

	// Specifies the Region where the bucket resides. For a list of all the Amazon
	// S3 supported location constraints by Region, see Regions and Endpoints (https://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region).

	location, err := svc.GetBucketLocation(params)
	if err != nil {
		return nil, err
	}

	if location != nil && location.LocationConstraint != nil {
		return location, nil
	}

	// Buckets in Region us-east-1 have a LocationConstraint of null.
	return &s3.GetBucketLocationOutput{
		LocationConstraint: aws.String("us-east-1"),
	}, nil
}

func getBucketIsPublic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketIsPublic")
	bucket := h.Item.(*s3.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, *location.LocationConstraint)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketPolicyStatusInput{
		Bucket: bucket.Name,
	}

	policyStatus, err := svc.GetBucketPolicyStatus(params)

	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NoSuchBucketPolicy" {
				return &s3.GetBucketPolicyStatusOutput{}, nil
			}
		}
		return nil, err
	}

	return policyStatus, nil
}

func getBucketVersioning(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketVersioning")
	bucket := h.Item.(*s3.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, *location.LocationConstraint)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketVersioningInput{
		Bucket: bucket.Name,
	}

	versioning, err := svc.GetBucketVersioning(params)
	if err != nil {
		return nil, err
	}

	return versioning, nil
}

func getBucketEncryption(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketEncryption")
	bucket := h.Item.(*s3.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, *location.LocationConstraint)
	if err != nil {
		return nil, err
	}
	params := &s3.GetBucketEncryptionInput{
		Bucket: bucket.Name,
	}

	encryption, err := svc.GetBucketEncryption(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ServerSideEncryptionConfigurationNotFoundError" {
				return nil, nil
			}
		}
		return nil, err
	}

	return encryption, nil
}

func getBucketPublicAccessBlock(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketPublicAccessBlock")
	bucket := h.Item.(*s3.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, *location.LocationConstraint)
	if err != nil {
		return nil, err
	}

	params := &s3.GetPublicAccessBlockInput{
		Bucket: bucket.Name,
	}

	defaultAccessBlock := &s3.PublicAccessBlockConfiguration{
		BlockPublicAcls:       aws.Bool(false),
		BlockPublicPolicy:     aws.Bool(false),
		IgnorePublicAcls:      aws.Bool(false),
		RestrictPublicBuckets: aws.Bool(false),
	}

	accessBlock, err := svc.GetPublicAccessBlock(params)
	if err != nil {
		// If the GetPublicAccessBlock is called on buckets which were created before Public Access Block setting was
		// introduced, sometime it fails with error NoSuchPublicAccessBlockConfiguration
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NoSuchPublicAccessBlockConfiguration" {
				return defaultAccessBlock, nil
			}
		}
		return nil, err
	}

	return accessBlock.PublicAccessBlockConfiguration, nil
}

func getBucketACL(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketACL")
	bucket := h.Item.(*s3.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, *location.LocationConstraint)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketAclInput{
		Bucket: bucket.Name,
	}

	acl, err := svc.GetBucketAcl(params)
	if err != nil {
		return nil, err
	}

	return acl, nil
}

func getBucketLifecycle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketLifecycle")
	bucket := h.Item.(*s3.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, *location.LocationConstraint)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: bucket.Name,
	}

	lifecycleConfiguration, err := svc.GetBucketLifecycleConfiguration(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NoSuchLifecycleConfiguration" {
				return nil, nil
			}
		}
		return nil, err
	}

	return lifecycleConfiguration, nil
}

func getBucketLogging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketLogging")
	bucket := h.Item.(*s3.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, *location.LocationConstraint)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketLoggingInput{
		Bucket: bucket.Name,
	}

	logging, err := svc.GetBucketLogging(params)
	if err != nil {
		return nil, err
	}
	return logging, nil
}

func getBucketPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketPolicy")
	bucket := h.Item.(*s3.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, *location.LocationConstraint)
	if err != nil {
		return nil, err
	}
	params := &s3.GetBucketPolicyInput{
		Bucket: bucket.Name,
	}

	bucketPolicy, err := svc.GetBucketPolicy(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NoSuchBucketPolicy" {
				return &s3.GetBucketPolicyOutput{}, nil
			}
		}
		return nil, err
	}

	return bucketPolicy, nil
}

func getBucketReplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketReplication")
	bucket := h.Item.(*s3.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, *location.LocationConstraint)
	if err != nil {
		return nil, err
	}
	params := &s3.GetBucketReplicationInput{
		Bucket: bucket.Name,
	}

	replication, err := svc.GetBucketReplication(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ReplicationConfigurationNotFoundError" {
				return &s3.GetBucketReplicationOutput{}, nil
			}
		}
		return nil, err
	}

	return replication, nil
}

func getBucketTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketTagging")
	bucket := h.Item.(*s3.Bucket)
	location := h.HydrateResults["getBucketLocation"].(*s3.GetBucketLocationOutput)

	// Create Session
	svc, err := S3Service(ctx, d.ConnectionManager, *location.LocationConstraint)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketTaggingInput{
		Bucket: bucket.Name,
	}

	bucketTags, _ := svc.GetBucketTagging(params)
	if err != nil {
		return nil, err
	}

	return bucketTags, nil
}

//// TRANSFORM FUNCTIONS

func s3NameToAkas(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("s3NameToAkas")
	bucket := d.HydrateItem.(*s3.Bucket)
	return []string{"arn:aws:s3:::" + *bucket.Name}, nil
}

func s3TagsToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("s3TagsToTurbotTags")
	tags := d.Value.([]*s3.Tag)

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
