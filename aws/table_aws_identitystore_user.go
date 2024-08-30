package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/identitystore/types"

	identitystorev1 "github.com/aws/aws-sdk-go/service/identitystore"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsIdentityStoreUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_identitystore_user",
		Description: "AWS Identity Store User",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"identity_store_id", "id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getIdentityStoreUser,
			Tags:    map[string]string{"service": "identitystore", "action": "DescribeUser"},
		},
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"identity_store_id"}),
			Hydrate:    listIdentityStoreUsers,
			Tags:       map[string]string{"service": "identitystore", "action": "ListUsers"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(identitystorev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "identity_store_id",
				Description: "The globally unique identifier for the identity store.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "Contains the user’s display name value.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.UserName"),
			},
			{
				Name:        "display_name",
				Description: "A string containing the name of the user that is formatted for display when the user is referenced.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.DisplayName"),
			},
			{
				Name:        "nick_name",
				Description: "A string containing an alternate name for the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.NickName"),
			},
			{
				Name:        "locale",
				Description: "A string containing the geographical region or location of the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.Locale"),
			},
			{
				Name:        "preferred_language",
				Description: "A string containing the preferred language of the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.PreferredLanguage"),
			},
			{
				Name:        "profile_url",
				Description: "A string containing a URL that might be associated with the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.ProfileUrl"),
			},
			{
				Name:        "timezone",
				Description: "A string containing the time zone of the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.Timezone"),
			},
			{
				Name:        "user_type",
				Description: "A string indicating the type of user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.UserType"),
			},
			{
				Name:        "id",
				Description: "The identifier for a user in the identity store.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.UserId"),
			},

			{
				Name:        "addresses",
				Description: "A list of Address objects containing addresses associated with the user.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("User.Addresses"),
			},
			{
				Name:        "emails",
				Description: "A list of Email objects containing email addresses associated with the user.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("User.Emails"),
			},
			{
				Name:        "external_ids",
				Description: "A list of ExternalId objects that contains the identifiers issued to this resource by an external identity provider.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("User.ExternalIds"),
			},
			{
				Name:        "phone_numbers",
				Description: "A list of PhoneNumber objects containing phone numbers associated with the user.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("User.PhoneNumbers"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.UserName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listIdentityStoreUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	identityStoreId := d.EqualsQuals["identity_store_id"].GetStringValue()

	// Create Session
	svc, err := IdentityStoreClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_identitystore_user.listIdentityStoreUsers", "get_client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	params := &identitystore.ListUsersInput{
		IdentityStoreId: aws.String(identityStoreId),
		MaxResults:      aws.Int32(maxLimit),
	}

	paginator := identitystore.NewListUsersPaginator(svc, params, func(o *identitystore.ListUsersPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_identitystore_user.listIdentityStoreUsers", "api_error", err)
			return nil, err
		}
		for _, user := range output.Users {
			item := &IdentityStoreUser{
				IdentityStoreId: &identityStoreId,
				User:            user,
			}
			d.StreamListItem(ctx, item)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

type IdentityStoreUser struct {
	IdentityStoreId *string
	User            types.User
}

//// HYDRATE FUNCTIONS

func getIdentityStoreUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	userId := d.EqualsQuals["id"].GetStringValue()
	identityStoreId := d.EqualsQuals["identity_store_id"].GetStringValue()

	// Create Session
	svc, err := IdentityStoreClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_identitystore_user.getIdentityStoreUser", "get_client_error", err)
		return nil, err
	}

	params := &identitystore.DescribeUserInput{
		UserId:          aws.String(userId),
		IdentityStoreId: aws.String(identityStoreId),
	}

	op, err := svc.DescribeUser(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_identitystore_user.getIdentityStoreUser", "api_error", err)
		return nil, err
	}

	item := &IdentityStoreUser{
		IdentityStoreId: &identityStoreId,
		User: types.User{
			UserName: op.UserName,
			UserId:   op.UserId,
		},
	}

	return item, nil
}
