package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"
	"github.com/aws/aws-sdk-go-v2/service/quicksight/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsQuickSightUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_quicksight_user",
		Description: "AWS QuickSight User",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "user_name", Require: plugin.Required},
				{Name: "namespace", Require: plugin.Required},
				{Name: "region", Require: plugin.Required},
				{Name: "quicksight_account_id", Require: plugin.Optional},
			},
			Hydrate: getAwsQuickSightUser,
			Tags:    map[string]string{"service": "quicksight", "action": "DescribeUser"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsQuickSightNamespaces,
			Hydrate:       listAwsQuickSightUsers,
			Tags:          map[string]string{"service": "quicksight", "action": "ListUsers"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "quicksight_account_id", Require: plugin.Optional},
				{Name: "namespace", Require: plugin.Optional},
				{Name: "region", Require: plugin.Required},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_QUICKSIGHT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "user_name",
				Description: "The user's user name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the user.",
				Type:        proto.ColumnType_STRING,
			},
			// As we have already a column "account_id" as a common column for all the tables, we have renamed the column to "quicksight_account_id"
			{
				Name:        "quicksight_account_id",
				Description: "The ID for the Amazon Web Services account that the user is in. Currently, you use the ID for the Amazon Web Services account that contains your Amazon QuickSight account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("quicksight_account_id"),
			},
			{
				Name:        "email",
				Description: "The user's email address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role",
				Description: "The user's role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity_type",
				Description: "The type of identity authentication used by the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "active",
				Description: "The active status of the user. When you create an Amazon QuickSight user that's not an IAM user or an Active Directory user, that user is inactive until they sign in and provide a password.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "principal_id",
				Description: "The principal ID of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "external_login_id",
				Description: "The identity ID for the user in the external login provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace",
				Description: "The namespace. Currently, you should set this to default.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "custom_permissions_name",
				Description: "The custom permissions profile associated with this user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "external_login_federation_provider_type",
				Description: "The type of supported external login provider that provides identity to let a user federate into Amazon QuickSight with an associated AWS Identity and Access Management (IAM) role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "external_login_federation_provider_url",
				Description: "The URL of the external login provider.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserName"),
			},
		}),
	}
}

type QuickSightUser struct {
	types.User
	Namespace string
}

//// LIST FUNCTION

func listAwsQuickSightUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_user.listAwsQuickSightUsers", "connection_error", err)
		return nil, err
	}

	// Get namespace from parent or quals
	namespaceInfo := h.Item.(types.NamespaceInfoV2)
	if d.EqualsQuals["namespace"] != nil && d.EqualsQuals["namespace"].GetStringValue() != *namespaceInfo.Name {
		return nil, nil
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

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &quicksight.ListUsersInput{
		AwsAccountId: aws.String(accountId),
		Namespace:    namespaceInfo.Name,
		MaxResults:   aws.Int32(maxLimit),
	}

	paginator := quicksight.NewListUsersPaginator(svc, input, func(o *quicksight.ListUsersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_quicksight_user.listAwsQuickSightUsers", "api_error", err)
			return nil, err
		}

		for _, item := range output.UserList {
			d.StreamListItem(ctx, QuickSightUser{User: item, Namespace: *namespaceInfo.Name})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsQuickSightUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_user.getAwsQuickSightUser", "connection_error", err)
		return nil, err
	}

	userName := d.EqualsQuals["user_name"].GetStringValue()
	namespace := d.EqualsQuals["namespace"].GetStringValue()

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

	params := &quicksight.DescribeUserInput{
		AwsAccountId: aws.String(accountId),
		Namespace:    aws.String(namespace),
		UserName:     aws.String(userName),
	}

	// Get call
	data, err := svc.DescribeUser(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_user.getAwsQuickSightUser", "api_error", err)
		return nil, err
	}

	return QuickSightUser{User: *data.User, Namespace: namespace}, nil
}
