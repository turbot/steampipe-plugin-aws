package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2Settings(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_settings",
		Description: "AWS EC2 Settings",
		List: &plugin.ListConfig{
			Hydrate: listAllAwsRegions,
		},
		Columns: []*plugin.Column{
			{
				Name:        "default_ebs_encryption_enabled",
				Description: "Indicates whether encryption by default is enabled.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getDefaultEBSVolumeEncryption,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "default_ebs_encryption_key",
				Description: "The Amazon Resource Name (ARN) or alias of the default CMK for encryption by default.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDefaultEBSVolumeEncryptionKey,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("EC2 Settings"),
			},
			{
				Name:        "partition",
				Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCommonColumns,
			},
			{
				Name:        "region",
				Description: "The AWS Region in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegionName"),
			},
			{
				Name:        "account_id",
				Description: "The AWS Account ID in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCommonColumns,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

//// LIST FUNCTION

func listAllAwsRegions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultAwsRegion(d)

	// Create Session
	svc, err := Ec2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true),
	}

	// execute list call
	resp, err := svc.DescribeRegions(params)
	if err != nil {
		return nil, err
	}

	for _, region := range resp.Regions {
		d.StreamListItem(ctx, region)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDefaultEBSVolumeEncryption(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDefaultEBSVolumeEncryption")
	data := h.Item.(*ec2.Region)

	// Returning false for disabled regions to avoid permission denied error AuthFailure
	if *data.OptInStatus == "not-opted-in" {
		return false, nil
	}
	// Create session
	svc, err := Ec2Service(ctx, d, *data.RegionName)
	if err != nil {
		return nil, err
	}
	params := &ec2.GetEbsEncryptionByDefaultInput{}
	defaultEncryption, err := svc.GetEbsEncryptionByDefault(params)
	if err != nil {
		return nil, err
	}
	return defaultEncryption.EbsEncryptionByDefault, nil
}

func getDefaultEBSVolumeEncryptionKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDefaultEBSVolumeEncryptionKey")
	data := h.Item.(*ec2.Region)

	// Returning default ebs key alias for disabled regions to avoid permission denied error AuthFailure
	if *data.OptInStatus == "not-opted-in" {
		return "alias/aws/ebs", nil
	}
	// Create session
	svc, err := Ec2Service(ctx, d, *data.RegionName)
	if err != nil {
		return nil, err
	}
	params := &ec2.GetEbsDefaultKmsKeyIdInput{}
	defaultEncryptionKey, err := svc.GetEbsDefaultKmsKeyId(params)
	if err != nil {
		return nil, err
	}
	return defaultEncryptionKey.KmsKeyId, nil
}
