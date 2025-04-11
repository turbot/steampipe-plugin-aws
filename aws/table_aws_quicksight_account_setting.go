package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsQuickSightAccountSetting(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_quicksight_account_setting",
		Description: "AWS QuickSight Account Setting",
		List: &plugin.ListConfig{
			Hydrate: listAwsQuickSightAccountSettings,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "region", Require: plugin.Required},
				{Name: "quicksight_account_id", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "quicksight", "action": "DescribeAccountSettings"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_QUICKSIGHT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "account_name",
				Description: "The account name displayed for the account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountSettings.AccountName"),
			},
			// As we have already a column "account_id" as a common column for all the tables, we have renamed the column to "quicksight_account_id"
			{
				Name:        "quicksight_account_id",
				Description: "The ID for the Amazon Web Services account that contains the settings that you want to list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("quicksight_account_id"),
			},
			{
				Name:        "edition",
				Description: "The edition of Amazon QuickSight that you're currently subscribed to: Enterprise edition or Standard edition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountSettings.Edition"),
			},
			{
				Name:        "default_namespace",
				Description: "The default namespace for this AWS account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountSettings.DefaultNamespace"),
			},
			{
				Name:        "notification_email",
				Description: "The email address that QuickSight uses to send notifications.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountSettings.NotificationEmail"),
			},
			{
				Name:        "termination_protection_enabled",
				Description: "A boolean value that indicates whether termination protection is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AccountSettings.TerminationProtectionEnabled"),
			},
			{
				Name:        "public_sharing_enabled",
				Description: "A boolean value that indicates whether public sharing is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AccountSettings.PublicSharingEnabled"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountSettings.AccountName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsQuickSightAccountSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_account_setting.listAwsQuickSightAccountSettings", "connection_error", err)
		return nil, err
	}

	accountId := d.EqualsQuals["quicksight_account_id"].GetStringValue()
	// Get AWS Account ID
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	if accountId == "" {
		accountId = commonColumnData.AccountId
	}

	// Build the params
	params := &quicksight.DescribeAccountSettingsInput{
		AwsAccountId: aws.String(accountId),
	}

	// Get call
	data, err := svc.DescribeAccountSettings(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_account_setting.listAwsQuickSightAccountSettings", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, data)

	return nil, nil
}
