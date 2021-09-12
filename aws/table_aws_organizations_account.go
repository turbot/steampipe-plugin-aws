package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsOrganizationsAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_organizations_account",
		Description: "AWS Organizations Account",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"AccountNotFoundException"}),
			Hydrate:           getOrganizationsAccount,
		},
		List: &plugin.ListConfig{
			Hydrate: listOrganizationsAccounts,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier (ID) of the account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email",
				Description: "The email address associated with the AWS account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "joined_method",
				Description: "The method by which the account joined the organization.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "joined_timestamp",
				Description: "The date the account became a part of the organization.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "name",
				Description: "The description of the permission set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the account in the organization.",
				Type:        proto.ColumnType_STRING,
			},
			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOrganizationsAccountTags,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listOrganizationsAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listOrganizationsAccounts")

	// Create session
	svc, err := OrganizationService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &organizations.ListAccountsInput{}

	err = svc.ListAccountsPages(
		params,
		func(page *organizations.ListAccountsOutput, isLast bool) bool {
			for _, account := range page.Accounts {
				d.StreamListItem(ctx, account)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOrganizationsAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getOrganizationsAccount")

	accountId := d.KeyColumnQuals["id"].GetStringValue()

	// Create session
	svc, err := OrganizationService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &organizations.DescribeAccountInput{
		AccountId: aws.String(accountId),
	}

	op, err := svc.DescribeAccount(params)
	if err != nil {
		return nil, err
	}

	return op.Account, nil
}

func getOrganizationsAccountTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getOrganizationsAccountTags")

	resourceId := *h.Item.(*organizations.Account).Id

	tags, err := getOrganizationsResourceTags(ctx, d, resourceId)
	return tags, err
}

func getOrganizationsResourceTags(ctx context.Context, d *plugin.QueryData, resourceId string) (interface{}, error) {
	plugin.Logger(ctx).Trace("getOrganizationsResourceTags")

	// Create Session
	svc, err := OrganizationService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &organizations.ListTagsForResourceInput{
		ResourceId: aws.String(resourceId),
	}

	tags := map[string]string{}
	err = svc.ListTagsForResourcePages(
		params,
		func(page *organizations.ListTagsForResourceOutput, isLast bool) bool {
			for _, i := range page.Tags {
				tags[*i.Key] = *i.Value
			}
			return !isLast
		},
	)

	return &tags, err
}

//// TRANSFORM FUNCTIONS
