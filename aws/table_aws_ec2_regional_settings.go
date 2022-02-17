package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2RegionalSettings(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_regional_settings",
		Description: "AWS EC2 Regional Settings",
		List: &plugin.ListConfig{
			Hydrate: listEc2RegionalSettings,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
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

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getEc2SettingTitle),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2RegionalSettings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listEc2RegionalSettings", "AWS_REGION", region)

	d.StreamListItem(ctx, region)
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDefaultEBSVolumeEncryption(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDefaultEBSVolumeEncryption")

	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listEc2RegionalSettings", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}
	params := &ec2.GetEbsEncryptionByDefaultInput{}
	defaultEncryption, err := svc.GetEbsEncryptionByDefault(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			// Returning false for disabled regions
			if a.Code() == "AuthFailure" {
				return false, nil
			}
		}
		return nil, err
	}
	return defaultEncryption.EbsEncryptionByDefault, nil
}

func getDefaultEBSVolumeEncryptionKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDefaultEBSVolumeEncryptionKey")

	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listEc2RegionalSettings", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}
	params := &ec2.GetEbsDefaultKmsKeyIdInput{}
	defaultEncryptionKey, err := svc.GetEbsDefaultKmsKeyId(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			// Returning default ebs key alias for disabled regions
			if a.Code() == "AuthFailure" {
				return "alias/aws/ebs", nil
			}
		}
		return nil, err
	}
	return defaultEncryptionKey.KmsKeyId, nil
}

//// TRANSFORM FUNCTIONS

func getEc2SettingTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	region := d.MatrixItem[matrixKeyRegion]

	title := region.(string) + " EC2 Settings"
	return title, nil
}
