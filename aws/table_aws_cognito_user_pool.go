package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	cognitoidentityproviderv1 "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCognitoUserPool(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cognito_user_pool",
		Description: "AWS Cognito User Pool",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCognitoUserPool,
			Tags:       map[string]string{"service": "cognito-idp", "action": "DescribeUserPool"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCognitoUserPools,
			Tags:    map[string]string{"service": "cognito-idp", "action": "ListUserPools"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCognitoUserPool,
				Tags: map[string]string{"service": "cognito-idp", "action": "DescribeUserPool"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cognitoidentityproviderv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the user pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the user pool.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_recovery_setting",
				Description: "The available verified method a user can use to recover their password when they call ForgotPassword.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "admin_create_user_config",
				Description: "The configuration for AdminCreateUser requests.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "alias_attributes",
				Description: "The attributes that are aliased in a user pool.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_verified_attributes",
				Description: "The attributes that are auto-verified in a user pool.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "The date the user pool was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "custom_domain",
				Description: "A custom domain name that you provide to Amazon Cognito.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deletion_protection",
				Description: "When active, DeletionProtection prevents accidental deletion of your user pool.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "device_configuration",
				Description: "The device-remembering configuration for a user pool.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "domain",
				Description: "The domain prefix, if the user pool has a domain associated with it.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email_configuration",
				Description: "The email configuration of your user pool.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "estimated_number_of_users",
				Description: "A number estimating the size of the user pool.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "lambda_config",
				Description: "The Lambda triggers associated with the user pool.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "last_modified_date",
				Description: "The date the user pool was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "mfa_configuration",
				Description: "Multi-Factor Authentication (MFA) configuration for the User Pool.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the Cognito User Pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policies",
				Description: "The policies associated with the user pool.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "schema_attributes",
				Description: "A container with the schema attributes of a user pool.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "sms_authentication_message",
				Description: "The contents of the SMS authentication message.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sms_configuration",
				Description: "The SMS configuration with the settings that your Amazon Cognito user pool must use to send an SMS message from your Amazon Web Services account through Amazon Simple Notification Service.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "sms_configuration_failure",
				Description: "The reason why the SMS configuration can't send the messages to your users.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of a user pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_attribute_update_settings",
				Description: "The settings for updates to user attributes.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "user_pool_add_ons",
				Description: "The user pool add-ons.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "username_attributes",
				Description: "Specifies whether a user can use an email address or phone number as a username when they sign up.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "username_configuration",
				Description: "Case sensitivity of the username input for the selected sign-in option.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "verification_message_template",
				Description: "The template for verification messages.",
				Hydrate:     getCognitoUserPool,
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCognitoUserPool,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCognitoUserPool,
				Transform:   transform.FromField("UserPoolTags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCognitoUserPools(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CognitoIdentityProviderClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_user_pool.listCognitoUserPools", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		plugin.Logger(ctx).Debug("aws_cognito_user_pool.listCognitoUserPools", "unsupported_region")
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(60)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	input := &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: maxLimit,
	}
	// List call
	paginator := cognitoidentityprovider.NewListUserPoolsPaginator(svc, input, func(o *cognitoidentityprovider.ListUserPoolsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cognito_user_pool.listCognitoUserPools", "api_error", err)
			return nil, err
		}
		plugin.Logger(ctx).Debug("aws_cognito_user_pool.listCognitoUserPools", "user_pools", fmt.Sprintf("%#v", output.UserPools))

		for _, userPool := range output.UserPools {
			plugin.Logger(ctx).Debug("aws_cognito_user_pool.listCognitoUserPools", "user_pool", fmt.Sprintf("%#v", userPool))
			d.StreamListItem(ctx, userPool)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCognitoUserPool(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		data := h.Item.(types.UserPoolDescriptionType)
		if data.Id != nil {
			id = *data.Id
		}
	} else {
		id = d.EqualsQualString("id")
	}
	plugin.Logger(ctx).Debug("aws_cognito_user_pool.getCognitoUserPool", "id", id)

	// check if id is empty
	if id == "" {
		return nil, nil
	}

	// Create Session
	svc, err := CognitoIdentityProviderClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_user_pool.getCognitoUserPool", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: aws.String(id),
	}

	// Get call
	data, err := svc.DescribeUserPool(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_user_pool.getCognitoUserPool", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("aws_cognito_user_pool.getCognitoUserPool", "user_pool", fmt.Sprintf("%#v", *data.UserPool))
	return *data.UserPool, nil
}
