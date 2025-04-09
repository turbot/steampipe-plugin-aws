package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsQuickSightAccountSettings(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_quicksight_account_settings",
		Description: "AWS QuickSight Account Settings",
		List: &plugin.ListConfig{
			Hydrate: listAwsQuickSightAccountSettings,
			Tags:    map[string]string{"service": "quicksight", "action": "DescribeAccountSettings"},
			// TODO do we need to add the account id as a qualifier?
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_QUICKSIGHT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "account_name",
				Description: "The account name displayed for the account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountSettings.AccountName"),
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
				Name:        "authentication_methods",
				Description: "The authentication methods that are enabled for this QuickSight account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccountSettings.AuthenticationMethods"),
			},
			{
				Name:        "data_lake_cross_service_role_arn",
				Description: "ARN of the cross service role that QuickSight data lake can assume for management of resources.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountSettings.DataLakeCrossServiceRoleArn"),
			},
			{
				Name:        "data_lake_metadata_catalog_cross_service_role_arn",
				Description: "ARN of the cross service role that QuickSight can assume for metadata catalog integration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountSettings.DataLakeMetadataCatalogCrossServiceRoleArn"),
			},
			{
				Name:        "custom_identity_configuration",
				Description: "The custom identity provider that's currently configured for your QuickSight account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccountSettings.CustomIdentityConfiguration"),
			},
			{
				Name:        "domain_configuration",
				Description: "The domain configuration for the QuickSight account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccountSettings.DomainConfiguration"),
			},
			{
				Name:        "vpc_connection_config",
				Description: "The VPC connection configuration for the QuickSight account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccountSettings.VpcConnectionConfiguration"),
			},
			{
				Name:        "private_spice_capacity_configuration",
				Description: "The private SPICE capacity configuration for the QuickSight account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccountSettings.PrivateSpiceCapacityConfiguration"),
			},
			{
				Name:        "public_sharing_configuration",
				Description: "The public sharing configuration for the QuickSight account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccountSettings.PublicSharingConfiguration"),
			},
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
		plugin.Logger(ctx).Error("aws_quicksight_account_settings.listAwsQuickSightAccountSettings", "connection_error", err)
		return nil, err
	}

	// Get AWS Account ID
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Build the params
	params := &quicksight.DescribeAccountSettingsInput{
		AwsAccountId: aws.String(commonColumnData.AccountId),
	}

	// Get call
	data, err := svc.DescribeAccountSettings(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_quicksight_account_settings.listAwsQuickSightAccountSettings", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, data)

	return nil, nil
} 