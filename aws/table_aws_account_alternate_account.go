package aws

import (
	"context"

	go_kit_packs "github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/aws/aws-sdk-go-v2/service/account/types"

	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAccountAlternateContact(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_account_alternate_account",
		Description: "AWS Account Alternate Contact",
		List: &plugin.ListConfig{
			Hydrate: listAwsAccountAlternateContacts,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:       "contact_account_id",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
			},
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name associated with this alternate contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlternateContact.Name"),
			},
			{
				Name:        "alternate_contact_type",
				Description: "The type of alternate contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlternateContact.AlternateContactType"),
			},
			{
				Name:        "contact_account_id",
				Description: "Account ID to get contact details for.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email_address",
				Description: "The email address associated with this alternate contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlternateContact.EmailAddress"),
			},
			{
				Name:        "phone_number",
				Description: "The phone number associated with this alternate contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlternateContact.PhoneNumber"),
			},
			{
				Name:        "title",
				Description: "The title associated with this alternate contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlternateContact.Title"),
			},
		}),
	}
}

type alternateAccountContactData = struct {
	AlternateContact types.AlternateContact
	ContactAccountId    *string
}

//// LIST FUNCTION

func listAwsAccountAlternateContacts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := AccountClient(ctx, d)
	if err != nil {
		logger.Error("aws_account_alternate_account.listAwsAccountAlternateContacts", "service_creation_error", err)
		return nil, err
	}

	if svc == nil {
		return nil, nil
	}

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	// List call for BILLING type alternate contact

	input_billing := &account.GetAlternateContactInput{}
	input_billing.AlternateContactType = "BILLING"
	if d.KeyColumnQuals["contact_account_id"] != nil {
		input_billing.AccountId = go_kit_packs.String(d.KeyColumnQuals["contact_account_id"].GetStringValue())
	}

	var contactAccountId string
	if d.KeyColumnQuals["contact_account_id"] == nil {
		contactAccountId = commonColumnData.AccountId
	} else {
		contactAccountId = *input_billing.AccountId
	}

	op_billing, err_billing := svc.GetAlternateContact(ctx, input_billing)
	if err_billing != nil {
		logger.Error("aws_account_alternate_account.listAwsAccountAlternateContacts", "billing_api_error", err_billing)
		return nil, err_billing
	}

	d.StreamListItem(ctx, &alternateAccountContactData{*op_billing.AlternateContact, &contactAccountId})

	// List call for SECURITY type alternate contact

	input_security := &account.GetAlternateContactInput{}
	input_security.AlternateContactType = "SECURITY"

	if d.KeyColumnQuals["contact_account_id"] != nil {
		input_security.AccountId = go_kit_packs.String(d.KeyColumnQuals["contact_account_id"].GetStringValue())
	}

	if d.KeyColumnQuals["contact_account_id"] == nil {
		contactAccountId = commonColumnData.AccountId
	} else {
		contactAccountId = *input_security.AccountId
	}

	op_security, err_security := svc.GetAlternateContact(ctx, input_security)
	if err_security != nil {
		logger.Error("aws_account_alternate_account.listAwsAccountAlternateContacts", "security_api_error", err_security)
		return nil, err_security
	}

	d.StreamListItem(ctx, &alternateAccountContactData{*op_security.AlternateContact, &contactAccountId})

	// List call for OPERATIONS type alternate contact

	input_operations := &account.GetAlternateContactInput{}
	input_operations.AlternateContactType = "OPERATIONS"

	if d.KeyColumnQuals["contact_account_id"] != nil {
		input_operations.AccountId = go_kit_packs.String(d.KeyColumnQuals["contact_account_id"].GetStringValue())
	}

	if d.KeyColumnQuals["contact_account_id"] == nil {
		contactAccountId = commonColumnData.AccountId
	} else {
		contactAccountId = *input_operations.AccountId
	}

	op_operations, err_operations := svc.GetAlternateContact(ctx, input_operations)
	if err_operations != nil {
		logger.Error("aws_account_alternate_account.listAwsAccountAlternateContacts", "operations_api_error", err_operations)
		return nil, err_operations
	}

	d.StreamListItem(ctx, &alternateAccountContactData{*op_operations.AlternateContact, &contactAccountId})

	return nil, nil
}