package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamVirtualMfaDevice(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_virtual_mfa_device",
		Description: "AWS IAM Virtual MFA device",
		List: &plugin.ListConfig{
			Hydrate: listIamVirtualMFADevices,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "assignment_status", Require: plugin.Optional},
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
	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_virtual_mfa_device.listIamVirtualMFADevices", "client_error", err)
		return nil, err
	}

	maxItems := int32(1000)
	input := iam.ListVirtualMFADevicesInput{}

	equalQuals := d.KeyColumnQuals
	if equalQuals["assignment_status"] != nil {
		if equalQuals["assignment_status"].GetStringValue() != "" {
			input.AssignmentStatus = types.AssignmentStatusType(equalQuals["assignment_status"].GetStringValue())
		}
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxItems = aws.Int32(maxItems)
	paginator := iam.NewListVirtualMFADevicesPaginator(svc, &input, func(o *iam.ListVirtualMFADevicesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iam_virtual_mfa_device.listIamVirtualMFADevices", "api_error", err)
			return nil, err
		}

		for _, mfaDevice := range output.VirtualMFADevices {
			d.StreamListItem(ctx, mfaDevice)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamMfaDeviceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(types.VirtualMFADevice)

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_virtual_mfa_device.getIamMfaDeviceTags", "client_error", err)
		return nil, err
	}

	params := &iam.ListMFADeviceTagsInput{SerialNumber: data.SerialNumber}

	op, err := svc.ListMFADeviceTags(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Debug("aws_iam_virtual_mfa_device.getIamMfaDeviceTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func virtualMfaDeviceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*iam.ListMFADeviceTagsOutput)
	var turbotTagsMap map[string]string
	if len(data.Tags) == 0 {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range data.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

func getAssignmentStatus(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.VirtualMFADevice)
	if data.User != nil {
		return "Assigned", nil
	}
	return "Unassigned", nil
}
