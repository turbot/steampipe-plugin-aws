package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsEBSEncryptionByDefault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ebs_encryption_by_default",
		Description: "AWS EBS Encryption By Default",
		List: &plugin.ListConfig{
			Hydrate: getAwsEBSEncryptionByDefault,
			Tags:    map[string]string{"service": "ec2", "action": "GetEbsEncryptionByDefault"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "ebs_encryption_by_default",
				Description: "Indicates whether encryption by default is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "ebs_default_kms_key_id",
				Description: "The Amazon Resource Name (ARN) of the default KMS key for encryption by default.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsEBSDefaultKmsKeyId,
				Transform:   transform.FromField("KmsKeyId"),
			},
		}),
	}
}

//// LIST FUNCTION

func getAwsEBSEncryptionByDefault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_encryption_by_default.getAwsEBSEncryptionByDefault", "connection_error", err)
		return nil, err
	}

	// Build params
	params := &ec2.GetEbsEncryptionByDefaultInput{}

	resp, err := svc.GetEbsEncryptionByDefault(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_encryption_by_default.getAwsEBSEncryptionByDefault", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, resp)

	return nil, nil
}

func getAwsEBSDefaultKmsKeyId(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_encryption_by_default.getAwsEBSDefaultKmsKeyId", "connection_error", err)
		return nil, err
	}

	// Build params
	params := &ec2.GetEbsDefaultKmsKeyIdInput{}

	resp, err := svc.GetEbsDefaultKmsKeyId(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_encryption_by_default.getAwsEBSDefaultKmsKeyId", "api_error", err)
		return nil, err
	}

	return resp, nil
}
