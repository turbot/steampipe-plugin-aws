package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/account/types"
	go_kit_packs "github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"

	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAccountAlternateContact(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_account_alternate_contact",
		Description: "AWS Account Alternate Contact",
		List: &plugin.ListConfig{
			Hydrate: listAwsAccountAlternateContacts,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:       "linked_account_id",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:    "contact_type",
					Require: plugin.Optional,
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
				Name:        "contact_type",
				Description: "The type of alternate contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlternateContact.AlternateContactType"),
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
				Name:        "contact_title",
				Description: "The title associated with this alternate contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlternateContact.Title"),
			},
			{
				Name:        "linked_account_id",
				Description: "Account ID to get alternate contact details for.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LinkedAccountID"),
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlternateContact.Name"),
			},
		}),
	}
}

type accountAlternateContactData = struct {
	AlternateContact types.AlternateContact
	LinkedAccountID  string
}

//// LIST FUNCTION

func listAwsAccountAlternateContacts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Account management APIs are not supported in GovCloud as of 2022-09-23
	if commonColumnData.Partition == "aws-us-gov" {
		return nil, nil
	}

	// Create service
	svc, err := AccountClient(ctx, d)
	if err != nil {
		logger.Error("aws_account_alternate_contact.listAwsAccountAlternateContacts", "service_creation_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var linkedAccountID string
	if d.KeyColumnQuals["linked_account_id"] != nil {
		linkedAccountID = d.KeyColumnQuals["linked_account_id"].GetStringValue()
	} else {
		linkedAccountID = commonColumnData.AccountId
	}

	contactTypes := []string{"BILLING", "OPERATIONS", "SECURITY"}
	if d.KeyColumnQuals["contact_type"] != nil {
		contactTypes = []string{d.KeyColumnQuals["contact_type"].GetStringValue()}
	}

	/*
		If calling from the org management account and the management account ID is
		given, the following error is thrown:
		Error: operation error Account: GetAlternateContact, https response error StatusCode: 403, RequestID: 01cb2b09-8b6a-4073-baba-5b9511632d2e, AccessDeniedException: User: arn:aws:iam::123456789012:user/steampipe-test is not authorized to perform: account:GetAlternateContact (The management account can only be managed using the standalone context from the management account.) (SQLSTATE HV000)

		If calling from a linked account and any account ID is given (even its own)
		the following error is thrown:
		Error: operation error Account: GetAlternateContact, https response error StatusCode: 403, RequestID: 875c3f06-611d-43e7-9d87-0f910dddea22, AccessDeniedException: User: arn:aws:iam::123456789012:user/steampipe-test is not authorized to perform: account:GetAlternateContact (SQLSTATE HV000)
	*/
	input := &account.GetAlternateContactInput{}
	if linkedAccountID != commonColumnData.AccountId {
		input.AccountId = go_kit_packs.String(linkedAccountID)
	}

	for _, contactType := range contactTypes {
		input.AlternateContactType = types.AlternateContactType(contactType)
		op, err := svc.GetAlternateContact(ctx, input)
		if err != nil {
			logger.Error("aws_account_alternate_contact.listAwsAccountAlternateContacts", "contact_type", contactType, "api_error", err)
			return nil, err
		}

		d.StreamListItem(ctx, &accountAlternateContactData{*op.AlternateContact, linkedAccountID})
	}

	return nil, nil
}
