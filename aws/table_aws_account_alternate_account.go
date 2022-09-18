package aws

import (
	"context"

	go_kit_packs "github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"

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
			Hydrate: getAwsAccountAlternateContact,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "data",
				Description: "The full name of the primary contact address.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

// type accountAlternateContactData = struct {
// 	AlternateContact types.AlternateContact
// 	ContactsAccountId    *string
// }

//// LIST FUNCTION

func getAwsAccountAlternateContact(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := AccountClient(ctx, d)
	if err != nil {
		logger.Error("aws_account_alternate_account.getAwsAccountAlternateContact", "service_creation_error", err)
		return nil, err
	}

	if svc == nil {
		return nil, nil
	}


	input := &account.GetAlternateContactInput{}
	input.AlternateContactType = "BILLING"
	if d.KeyColumnQuals["contacts_account_id"] != nil {
		input.AccountId = go_kit_packs.String(d.KeyColumnQuals["contacts_account_id"].GetStringValue())
	}

	op, err := svc. GetAlternateContact(ctx, input)
	if err != nil {
		logger.Error("aws_account_alternate_account.getAwsAccountAlternateContact", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, *op.AlternateContact)

	return nil, nil
}
