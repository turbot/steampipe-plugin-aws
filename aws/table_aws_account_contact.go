package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"

	// "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go-v2/service/account"
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
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "data",
				Description: "The name of the workgroup.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAccountContact,
				Transform:   transform.FromValue(),
			},
		}),
	}
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

	// getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	// commonData, err := getCommonColumnsCached(ctx, d, h)
	// if err != nil {
	// 	return nil, err
	// }
	// commonColumnData := commonData.(*awsCommonColumnData)

	// logger.Debug("aws_account_contact.input", "input", *commonColumnData)

	// // Reduce the basic request limit down if the user has only requested a small number of rows

	// input := &account.GetContactInformationInput{
	// 	AccountId: &commonColumnData.AccountId,
	// }

	// logger.Debug("aws_account_contact.input", "input", *input.AccountId)
	// execute list call
	op, err := svc.GetContactInformation(ctx, &account.GetContactInformationInput{})
	if err != nil {
		return nil, err
	}
	logger.Debug("op", "op", *op.ContactInformation)
	return op.ContactInformation, nil
}
