package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/identitystore"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsIdentityStoreUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_identitystore_user",
		Description: "AWS Identity Store User",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"identity_store_id", "id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getIdentityStoreUser,
		},
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"identity_store_id", "name"}),
			Hydrate:    listIdentityStoreUsers,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Name:        "id",
				Description: "The identifier for a user in the identity store.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.UserId"),
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
	plugin.Logger(ctx).Trace("listIdentityStoreUsers")

	name := d.KeyColumnQuals["name"].GetStringValue()
	identityStoreId := d.KeyColumnQuals["identity_store_id"].GetStringValue()

	// Create session
	svc, err := IdentityStoreService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &identitystore.ListUsersInput{
		IdentityStoreId: aws.String(identityStoreId),
		Filters: []*identitystore.Filter{
			{
				AttributePath:  aws.String("UserName"),
				AttributeValue: aws.String(name),
			},
		},
		MaxResults: aws.Int64(50),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxResults {
			if *limit < 1 {
				params.MaxResults = aws.Int64(1)
			} else {
				params.MaxResults = limit
			}
		}
	}

	err = svc.ListUsersPages(
		params,
		func(page *identitystore.ListUsersOutput, isLast bool) bool {
			for _, user := range page.Users {
				item := &IdentityStoreUser{
					IdentityStoreId: &identityStoreId,
					User:            user,
				}
				d.StreamListItem(ctx, item)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listIdentityStoreUsers", "ListUsersPages_error", err)

	}
	return nil, err
}

type IdentityStoreUser struct {
	IdentityStoreId *string
	User            *identitystore.User
}

//// HYDRATE FUNCTIONS

func getIdentityStoreUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIdentityStoreUser")

	userId := d.KeyColumnQuals["id"].GetStringValue()
	identityStoreId := d.KeyColumnQuals["identity_store_id"].GetStringValue()

	// Create session
	svc, err := IdentityStoreService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &identitystore.DescribeUserInput{
		UserId:          aws.String(userId),
		IdentityStoreId: aws.String(identityStoreId),
	}

	op, err := svc.DescribeUser(params)
	if err != nil {
		plugin.Logger(ctx).Error("getIdentityStoreUser", "DescribeUser_error", err)
		return nil, err
	}

	item := &IdentityStoreUser{
		IdentityStoreId: &identityStoreId,
		User: &identitystore.User{
			UserName: op.UserName,
			UserId:   op.UserId,
		},
	}

	return item, nil
}
