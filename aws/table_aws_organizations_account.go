package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsOrganizationsAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_organizations_account",
		Description: "AWS Organizations Account",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"AccountNotFoundException", "InvalidInputException"}),
			},
			Hydrate: getOrganizationsAccount,
		},
		List: &plugin.ListConfig{
			Hydrate: listOrganizationsAccounts,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The description of the permission set.",
				Type:        proto.ColumnType_STRING,
			},
			// This description has added text for better clarification on ID type
			{
				Name:        "id",
				Description: "The unique identifier (account ID) of the member account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the account in the organization.",
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
				Name:      "tags_src",
				Type:      proto.ColumnType_JSON,
				Hydrate:   getOrganizationsAccountTags,
				Transform: transform.FromValue(),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOrganizationsAccountTags,
				Transform:   transform.From(getOrganizationsResourceTurbotTags),
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
	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_account.listOrganizationsAccounts", "client_error", err)
		return nil, err
	}

	// The maximum number for MaxResults parameter is not defined by the API
	// We have set the MaxResults to 1000 based on our test
	maxItems := int32(20)
	params := &organizations.ListAccountsInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	params.MaxResults = &maxItems
	paginator := organizations.NewListAccountsPaginator(svc, params, func(o *organizations.ListAccountsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_organizations_account.listOrganizationsAccounts", "api_error", err)
			return nil, err
		}

		for _, account := range output.Accounts {
			d.StreamListItem(ctx, account)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOrganizationsAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	accountId := d.KeyColumnQuals["id"].GetStringValue()

	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_account.getOrganizationsAccount", "client_error", err)
		return nil, err
	}

	params := &organizations.DescribeAccountInput{
		AccountId: &accountId,
	}

	op, err := svc.DescribeAccount(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_account.getOrganizationsAccount", "api_error", err)
		return nil, err
	}

	return *op.Account, nil
}

func getOrganizationsAccountTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	resourceId := *h.Item.(types.Account).Id

	tags, err := getOrganizationsResourceTags(ctx, d, resourceId)
	return tags, err
}

func getOrganizationsResourceTags(ctx context.Context, d *plugin.QueryData, resourceId string) (interface{}, error) {

	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_account.getOrganizationsResourceTags", "client_error", err)
		return nil, err
	}

	params := &organizations.ListTagsForResourceInput{
		ResourceId: &resourceId,
	}

	tags := []types.Tag{}

	paginator := organizations.NewListTagsForResourcePaginator(svc, params, func(o *organizations.ListTagsForResourcePaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_organizations_account.getOrganizationsResourceTags", "api_error", err)
			return nil, err
		}

		tags = append(tags, output.Tags...)
	}

	return tags, err
}

//// TRANSFORM FUNCTIONS

func getOrganizationsResourceTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.([]types.Tag)
	tagsMap := map[string]string{}

	for _, tag := range tags {
		tagsMap[*tag.Key] = *tag.Value
	}

	return tagsMap, nil
}
