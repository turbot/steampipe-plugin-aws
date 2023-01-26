package aws

import (
	"context"

	go_kit_types "github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"

	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/aws/aws-sdk-go-v2/service/account/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAccountContact(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_account_contact",
		Description: "AWS Account Contact",
		List: &plugin.ListConfig{
			Hydrate: listAwsAccountContacts,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:       "linked_account_id",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "full_name",
				Description: "The full name of the primary contact address.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.FullName"),
			},
			{
				Name:        "address_line_1",
				Description: "The first line of the primary contact address",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.AddressLine1"),
			},
			{
				Name:        "address_line_2",
				Description: "The second line of the primary contact address, if any.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.AddressLine2"),
			},
			{
				Name:        "address_line_3",
				Description: "The third line of the primary contact address, if any.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.AddressLine3"),
			},
			{
				Name:        "company_name",
				Description: "The name of the company associated with the primary contact information, if any.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.CompanyName"),
			},
			{
				Name:        "city",
				Description: "The city of the primary contact address.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.City"),
			},
			{
				Name:        "country_code",
				Description: "The ISO-3166 two-letter country code for the primary contact address.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.CountryCode"),
			},
			{
				Name:        "district_or_county",
				Description: "The district or county of the primary contact address, if any.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.DistrictOrCounty"),
			},
			{
				Name:        "phone_number",
				Description: "The phone number of the primary contact information. The number will be validated and, in some countries, checked for activation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.PhoneNumber"),
			},
			{
				Name:        "postal_code",
				Description: "The postal code of the primary contact address.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.PostalCode"),
			},
			{
				Name:        "state_or_region",
				Description: "The state or region of the primary contact address. This field is required in selected countries.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.StateOrRegion"),
			},
			{
				Name:        "website_url",
				Description: "The URL of the website associated with the primary contact information, if any.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.WebsiteUrl"),
			},
			{
				Name:        "linked_account_id",
				Description: "Account ID to get contact details for.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LinkedAccountID"),
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactInformation.FullName"),
			},
		}),
	}
}

type accountContactData = struct {
	ContactInformation types.ContactInformation
	LinkedAccountID    string
}

//// LIST FUNCTION

func listAwsAccountContacts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	commonData, err := getCommonColumns(ctx, d, h)
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
		logger.Error("aws_account_contact.listAwsAccountContacts", "service_creation_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var linkedAccountID string
	if d.EqualsQuals["linked_account_id"] != nil {
		linkedAccountID = d.EqualsQuals["linked_account_id"].GetStringValue()
	} else {
		linkedAccountID = commonColumnData.AccountId
	}

	/*
		If calling from the org management account and the management account ID is
		given, the following error is thrown:
		Error: operation error Account: GetContactInformation, https response error StatusCode: 403, RequestID: 01cb2b09-8b6a-4073-baba-5b9511632d2e, AccessDeniedException: User: arn:aws:iam::123456789012:user/steampipe-test is not authorized to perform: account:GetContactInformation (The management account can only be managed using the standalone context from the management account.) (SQLSTATE HV000)

		If calling from a linked account and any account ID is given (even its own)
		the following error is thrown:
		Error: operation error Account: GetContactInformation, https response error StatusCode: 403, RequestID: 875c3f06-611d-43e7-9d87-0f910dddea22, AccessDeniedException: User: arn:aws:iam::123456789012:user/steampipe-test is not authorized to perform: account:GetContactInformation (SQLSTATE HV000)
	*/
	input := &account.GetContactInformationInput{}
	if linkedAccountID != commonColumnData.AccountId {
		input.AccountId = go_kit_types.String(linkedAccountID)
	}

	op, err := svc.GetContactInformation(ctx, input)
	if err != nil {
		logger.Error("aws_account_contact.listAwsAccountContacts", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, &accountContactData{*op.ContactInformation, linkedAccountID})

	return nil, nil
}
