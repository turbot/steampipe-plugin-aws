package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsS3AccountSettings(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_account_settings",
		Description: "AWS S3 Account Block Public Access Settings",
		List: &plugin.ListConfig{
			Hydrate: listS3Account,
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
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

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_account_settings.listS3Account", "common_data_error", err)
		return nil, err
	}

	// returning the current account details as a list item
	d.StreamListItem(ctx, commonData)
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAccountBucketPublicAccessBlock(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	s3Account := h.Item.(*awsCommonColumnData)

	// Unlike most services, S3 buckets are a global list. They can be retrieved
	// from any single region. It's best to use the client region of the user
	// (e.g. closest to them).
	clientRegion, err := getDefaultRegion(ctx, d, h)
	if err != nil {
		return nil, err
	}
	svc, err := S3ControlClient(ctx, d, clientRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_account_settings.getAccountBucketPublicAccessBlock", "get_client_error", err, "clientRegion", clientRegion)
		return nil, err
	}

	params := &s3control.GetPublicAccessBlockInput{
		AccountId: &s3Account.AccountId,
	}

	defaultAccessBlock := &types.PublicAccessBlockConfiguration{
		BlockPublicAcls:       false,
		BlockPublicPolicy:     false,
		IgnorePublicAcls:      false,
		RestrictPublicBuckets: false,
	}

	accessBlock, err := svc.GetPublicAccessBlock(ctx, params)
	if err != nil {
		// If the GetPublicAccessBlock is called on an account ( that was created
		// before Public Access Block setting was introduced ), sometime it
		// fails with  NoSuchPublicAccessBlockConfiguration error
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NoSuchPublicAccessBlockConfiguration" {
				return defaultAccessBlock, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_account_settings.getAccountBucketPublicAccessBlock", "api_error", err, "clientRegion", clientRegion)
		return nil, err
	}

	return accessBlock.PublicAccessBlockConfiguration, nil
}

//// Transform Functions

func s3AccountDataToAkas(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	accountInfo := d.HydrateItem.(*awsCommonColumnData)

	akas := []string{"arn:" + accountInfo.Partition + ":s3::" + accountInfo.AccountId + ":account"}

	return akas, nil
}
