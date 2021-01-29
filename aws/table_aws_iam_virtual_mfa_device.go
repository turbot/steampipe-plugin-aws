package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamVirtualMfaDevice(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_virtual_mfa_device",
		Description: "AWS IAM Virtual MFA device",
		List: &plugin.ListConfig{
			Hydrate: listIamVirtualMFADevices,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "serial_number",
				Description: "The serial number associated with VirtualMFADevice.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_date",
				Description: "The date and time on which the virtual MFA device was enabled.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "user_id",
				Description: "The user id of the user associated with this virtual MFA device.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.UserId"),
			},
			{
				Name:        "user_name",
				Description: "The friendly name of the user associated with this virtual MFA device.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.UserName"),
			},
			{
				Name:        "user",
				Description: "Details of the IAM user associated with this virtual MFA device.",
				Type:        proto.ColumnType_JSON,
			},

			// {
			// 	Name:        "base_32_string",
			// 	Description: "The friendly name identifying the user",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromField("Base32StringSeed"),
			// },
			// {
			// 	Name:        "QRCodePNG",
			// 	Description: "A QR code PNG image that encodes otpauth://totp/$virtualMFADeviceName@$AccountName?secret=$Base32String where $virtualMFADeviceName is one of the create call arguments. AccountName is the user name if set (otherwise, the account ID otherwise), and Base32String is the seed in base32 format. The Base32String value is base64-encoded.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromField("QRCodePNG"),
			// },

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SerialNumber"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SerialNumber").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listIamVirtualMFADevices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IAMService(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	err = svc.ListVirtualMFADevicesPages(
		&iam.ListVirtualMFADevicesInput{},
		func(page *iam.ListVirtualMFADevicesOutput, lastPage bool) bool {
			for _, mfaDevice := range page.VirtualMFADevices {
				d.StreamListItem(ctx, mfaDevice)
			}
			return true
		},
	)

	return nil, err
}
