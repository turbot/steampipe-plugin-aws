package aws

import (
	"context"

	go_kit_packs "github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"

	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/aws/aws-sdk-go-v2/service/account/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAccountContact(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_account_contact",
		Description: "AWS Account Contact",
		List: &plugin.ListConfig{
			Hydrate: getAwsAccountContact,
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
				Name:        "contact_account_id",
				Description: "Account ID to get contact details for.",
				Type:        proto.ColumnType_STRING,
				// Transform:   transform.FromQual("contact_account_id"),
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
		}),
	}
}

type accountContactData = struct {
	ContactInformation types.ContactInformation
	ContactAccountId    *string
}

//// LIST FUNCTION

func getAwsAccountContact(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := AccountClient(ctx, d)
	if err != nil {
		logger.Error("aws_account_contact.getAwsAccountContact", "service_creation_error", err)
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

	input := &account.GetContactInformationInput{}
	if d.KeyColumnQuals["contact_account_id"] != nil {
		input.AccountId = go_kit_packs.String(d.KeyColumnQuals["contact_account_id"].GetStringValue())
	}

	var contactAccountId string
	if d.KeyColumnQuals["contact_account_id"] == nil {
		contactAccountId = commonColumnData.AccountId
	} else {
		contactAccountId = *input.AccountId
	}

	op, err := svc.GetContactInformation(ctx, input)
	if err != nil {
		logger.Error("aws_account_contact.getAwsAccountContact", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, &accountContactData{*op.ContactInformation, &contactAccountId})

	return nil, nil
}
