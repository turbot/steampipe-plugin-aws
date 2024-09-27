package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/shield"
	"github.com/aws/aws-sdk-go-v2/service/shield/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsShieldEmergencyContact(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_shield_emergency_contact",
		Description: "AWS Shield Emergency Contact",
		List: &plugin.ListConfig{
			Hydrate: listAwsShieldEmergencyContact,
			Tags:    map[string]string{"service": "shield", "action": "DescribeEmergencyContactSettings"},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "email_address",
				Description: "The email address for the contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EmailAddress"),
			},
			{
				Name:        "phone_number",
				Description: "The phone number for the contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PhoneNumber"),
			},
			{
				Name:        "ContactNotes",
				Description: "Additional notes regarding the contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactNotes"),
			},
		}),
	}
}

func listAwsShieldEmergencyContact(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_emergency_contact.listAwsShieldEmergencyContact", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	data, err := svc.DescribeEmergencyContactSettings(ctx, &shield.DescribeEmergencyContactSettingsInput{})

	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_emergency_contact.listAwsShieldEmergencyContact", "api_error", err)
		return nil, err
	}

	for _, contact := range data.EmergencyContactList {
		d.StreamListItem(ctx, &types.EmergencyContact{
			EmailAddress:	contact.EmailAddress,
			PhoneNumber:	contact.PhoneNumber,
			ContactNotes:	contact.ContactNotes,
		})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
