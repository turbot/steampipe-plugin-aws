package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2RegionalSettings(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_regional_settings",
		Description: "AWS EC2 Regional Settings",
		List: &plugin.ListConfig{
			Hydrate: listEc2RegionalSettings,
		},
		GetMatrixItemFunc: BuildRegionList,
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

	d.StreamListItem(ctx, region)
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDefaultEBSVolumeEncryption(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_regional_settings.getDefaultEBSVolumeEncryption", "connection_error", err)
		return nil, err
	}
	params := &ec2.GetEbsEncryptionByDefaultInput{}
	defaultEncryption, err := svc.GetEbsEncryptionByDefault(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			// Return default ebs key alias for disabled regions
			if ae.ErrorCode() == "AuthFailure" {
				return false, nil
			}
		}
		plugin.Logger(ctx).Error("aws_ec2_regional_settings.getDefaultEBSVolumeEncryption", "api_error", err)
		return nil, err
	}
	return defaultEncryption.EbsEncryptionByDefault, nil
}

func getDefaultEBSVolumeEncryptionKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_regional_settings.getDefaultEBSVolumeEncryptionKey", "connection_error", err)
		return nil, err
	}
	params := &ec2.GetEbsDefaultKmsKeyIdInput{}
	defaultEncryptionKey, err := svc.GetEbsDefaultKmsKeyId(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_regional_settings.getDefaultEBSVolumeEncryptionKey", "api_error", err)
		if strings.Contains(err.Error(), "AuthFailure") {
			return "alias/aws/ebs", nil
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
