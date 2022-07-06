package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3control"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsS3AccountSettings(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_account_settings",
		Description: "AWS S3 Account Block Public Access Settings",
		List: &plugin.ListConfig{
			Hydrate: listS3Account,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "block_public_acls",
				Description: "Specifies whether Amazon S3 should block public access control lists (ACLs) for this bucket and objects in this bucket",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAccountBucketPublicAccessBlock,
				Transform:   transform.FromField("BlockPublicAcls"),
			},
			{
				Name:        "block_public_policy",
				Description: "Specifies whether Amazon S3 should block public bucket policies for this bucket. If TRUE it causes Amazon S3 to reject calls to PUT Bucket policy if the specified bucket policy allows public access",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAccountBucketPublicAccessBlock,
				Transform:   transform.FromField("BlockPublicPolicy"),
			},
			{
				Name:        "ignore_public_acls",
				Description: "Specifies whether Amazon S3 should ignore public ACLs for this bucket and objects in this bucket. Setting this element to TRUE causes Amazon S3 to ignore all public ACLs on this bucket and objects in this bucket",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAccountBucketPublicAccessBlock,
				Transform:   transform.FromField("IgnorePublicAcls"),
			},
			{
				Name:        "restrict_public_buckets",
				Description: "Specifies whether Amazon S3 should restrict public bucket policies for this bucket. Setting this element to TRUE restricts access to this bucket to only AWS service principals and authorized users within this account if the bucket has a public policy",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAccountBucketPublicAccessBlock,
				Transform:   transform.FromField("RestrictPublicBuckets"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Default:     "S3 Account Level Block Public Access Settings",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(s3AccountDataToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listS3Account(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	// returning the current account details as a list item
	d.StreamListItem(ctx, commonData)
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAccountBucketPublicAccessBlock(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAccountBucketPublicAccessBlock")
	s3Account := h.Item.(*awsCommonColumnData)

	// Create Session
	svc, err := S3ControlService(ctx, d, GetDefaultAwsRegionV1(d))
	if err != nil {
		return nil, err
	}

	params := &s3control.GetPublicAccessBlockInput{
		AccountId: &s3Account.AccountId,
	}

	defaultAccessBlock := &s3control.PublicAccessBlockConfiguration{
		BlockPublicAcls:       aws.Bool(false),
		BlockPublicPolicy:     aws.Bool(false),
		IgnorePublicAcls:      aws.Bool(false),
		RestrictPublicBuckets: aws.Bool(false),
	}

	accessBlock, err := svc.GetPublicAccessBlock(params)
	if err != nil {
		// If the GetPublicAccessBlock is called on an account ( that was created
		// before Public Access Block setting was introduced ), sometime it
		// fails with  NoSuchPublicAccessBlockConfiguration error
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NoSuchPublicAccessBlockConfiguration" {
				return defaultAccessBlock, nil
			}
		}
		return nil, err
	}

	return accessBlock.PublicAccessBlockConfiguration, nil
}

//// Transform Functions

func s3AccountDataToAkas(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("s3AccountDataToAkas")
	accountInfo := d.HydrateItem.(*awsCommonColumnData)

	akas := []string{"arn:" + accountInfo.Partition + ":s3::" + accountInfo.AccountId + ":account"}

	return akas, nil
}
