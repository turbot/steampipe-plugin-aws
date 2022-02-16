package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamVirtualMfaDevice(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_virtual_mfa_device",
		Description: "AWS IAM Virtual MFA device",
		List: &plugin.ListConfig{
			Hydrate: listIamVirtualMFADevices,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "assignment_status",
					Require: plugin.Optional,
				},
			},
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
				Name:        "assignment_status",
				Description: "The status (Unassigned or Assigned) of the device.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getAssignmentStatus),
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
			{
				Name:        "tags_src",
				Description: "A list of tags attached with the MFA device.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamMfaDeviceTags,
				Transform:   transform.FromField("Tags"),
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
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamMfaDeviceTags,
				Transform:   transform.From(virtualMfaDeviceTurbotTags),
			},
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
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &iam.ListVirtualMFADevicesInput{
		MaxItems: aws.Int64(1000),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["assignment_status"] != nil {
		if equalQuals["assignment_status"].GetStringValue() != "" {
			input.AssignmentStatus = aws.String(equalQuals["assignment_status"].GetStringValue())
		}
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxItems {
			if *limit < 1 {
				input.MaxItems = aws.Int64(1)
			} else {
				input.MaxItems = limit
			}
		}
	}

	err = svc.ListVirtualMFADevicesPages(
		input,
		func(page *iam.ListVirtualMFADevicesOutput, _ bool) bool {
			for _, mfaDevice := range page.VirtualMFADevices {
				d.StreamListItem(ctx, mfaDevice)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamMfaDeviceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIamMfaDeviceTags")

	data := h.Item.(*iam.VirtualMFADevice)

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &iam.ListMFADeviceTagsInput{
		SerialNumber: data.SerialNumber,
	}

	op, err := svc.ListMFADeviceTags(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getIamMfaDeviceTags__", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func virtualMfaDeviceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*iam.ListMFADeviceTagsOutput)
	var turbotTagsMap map[string]string
	if data.Tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range data.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

func getAssignmentStatus(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*iam.VirtualMFADevice)
	if data.User != nil {
		return "Assigned", nil
	}
	return "Unassigned", nil
}
