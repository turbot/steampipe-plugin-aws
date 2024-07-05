package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table behavior:
// 1. Uses `aws_organizations_root` as the parent hydrate function.
// 2. Avoids using the "ListAccounts" API call to retrieve all accounts
//    within the organization.
// 3. If `parent_id` is specified in the query parameter,
//    lists the accounts for that specific parent.
// 4. If `parent_id` is not specified, lists the accounts under all
//    Organizational Units (OUs), including the Root ID.

// Reason for this table design:
// 1. To address the issue described here:
//    https://github.com/turbot/steampipe-plugin-aws/issues/2235
// 2. As per the query plan for the query:
//    select id, parent_id, title from aws_organizations_account WHERE parent_id IN (select id from aws_organizations_organizational_unit WHERE parent_id='ou-wxnb-wofu2g1q') limit 2
//    Postgres does not provide the `parent_id` value with our previous design where the `parent_id` column value was populated from the query parameter.
//    This resulted in an empty result set due to Steampipe's level filtration on the `parent_id` value ("FromQual()" returns a null value) and the value passed in the query parameter mismatch.

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
			Tags:    map[string]string{"service": "organizations", "action": "DescribeAccount"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listOrganizationsRoots,
			Hydrate:       listOrganizationsAccounts,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "parent_id", Require: plugin.Optional, CacheMatch: "exact"},
			},
			Tags: map[string]string{"service": "organizations", "action": "ListAccounts"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getOrganizationsAccountTags,
				Tags: map[string]string{"service": "organizations", "action": "ListTagsForResource"},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name of the account.",
				Type:        proto.ColumnType_STRING,
			},
			// This description has added text for better clarification on ID type
			{
				Name:        "id",
				Description: "The unique identifier (account ID) of the member account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parent_id",
				Description: "The unique identifier (ID) for the parent root or organization unit (OU) whose accounts you want to list.",
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

type OrgAccount struct {
	types.Account
	ParentId *string
}

//// LIST FUNCTION

func listOrganizationsAccounts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	parentId := *h.Item.(types.Root).Id

	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_account.listOrganizationsAccounts", "client_error", err)
		return nil, err
	}

	// Limit the result
	maxItems := int32(20)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	var parents []types.OrganizationalUnit

	// Lists the accounts in an organization that are contained by the specified target root or organizational unit (OU).
	// If you specify the root, you get a list of all the accounts that aren't in any OU.
	// If you specify an OU, you get a list of all the accounts in only that OU and not in any child OUs.
	if d.EqualsQualString("parent_id") != "" {
		parents = []types.OrganizationalUnit{{
			Id: aws.String(d.EqualsQualString("parent_id")),
		}}
	}

	// Restrict listing of parents when "parent_id" is provided.
	if !(len(parents) > 0) {
		// Call the recursive function to list all nested OUs
		rootPath := parentId
		res, err := listAllOusByParent(ctx, d, svc, parentId, maxItems, rootPath)
		if err != nil {
			return nil, err
		}

		parents = append(res, types.OrganizationalUnit{Id: aws.String(rootPath)})
	}

	// Iterate the API call based on the parent entities
	for _, p := range parents {
		params := &organizations.ListAccountsForParentInput{
			ParentId:   p.Id,
			MaxResults: &maxItems,
		}
		paginator := organizations.NewListAccountsForParentPaginator(svc, params, func(o *organizations.ListAccountsForParentPaginatorOptions) {
			o.Limit = maxItems
			o.StopOnDuplicateToken = true
		})

		for paginator.HasMorePages() {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			output, err := paginator.NextPage(ctx)
			if err != nil {
				plugin.Logger(ctx).Error("aws_organizations_account.listOrganizationsAccounts", "api_error", err)
				return nil, err
			}

			for _, account := range output.Accounts {
				d.StreamListItem(ctx, OrgAccount{account, p.Id})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOrganizationsAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	accountId := d.EqualsQuals["id"].GetStringValue()

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

	// The "parent_id" column value will not be populated by the GET API call because its response does not include parent ID information.
	// We can populate this value by making a separate hydrated function.
	// However, to avoid unnecessary API calls for all rows.
	// If our LIST API call is executed, it will correctly populate the parent_id column value.
	// In the case of a GET API call, the Parent ID will not be available. Therefore, we make an additional API call to populate the parent_id column value when only the GET configuration is used.
	parent, err := getParentForAccount(ctx, svc, accountId)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_account.getParentForAccount", "api_error", err)
		return nil, err
	}

	return OrgAccount{*op.Account, parent}, nil
}

func getParentForAccount(ctx context.Context, client *organizations.Client, accountId string) (*string, error) {
	// Pagination is not needed here because an account will always have a single parent.
	parent, err := client.ListParents(ctx, &organizations.ListParentsInput{ChildId: &accountId})
	if err != nil {
		return nil, err
	}
	if len(parent.Parents) > 0 {
		return parent.Parents[0].Id, nil
	}
	return nil, nil
}

// List all nested Organizational Units (OUs)
// We cannot use the existing table "aws_organizations_organizational_unit" as the parent hydrate because it already has a parent hydrate, and Steampipe does not support parent hydrate chaining.
func listAllOusByParent(ctx context.Context, d *plugin.QueryData, svc *organizations.Client, parentId string, maxItems int32, currentPath string) ([]types.OrganizationalUnit, error) {
	params := &organizations.ListOrganizationalUnitsForParentInput{
		ParentId:   aws.String(parentId),
		MaxResults: &maxItems,
	}

	var units []types.OrganizationalUnit
	paginator := organizations.NewListOrganizationalUnitsForParentPaginator(svc, params, func(o *organizations.ListOrganizationalUnitsForParentPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_organizations_account.listAllOusByParent", "api_error", err)
			return nil, err
		}

		for _, unit := range output.OrganizationalUnits {
			ouPath := strings.Replace(currentPath, "-", "_", -1) + "." + strings.Replace(*unit.Id, "-", "_", -1)
			units = append(units, unit)

			// Recursively list units for this child
			childUnits, err := listAllOusByParent(ctx, d, svc, *unit.Id, maxItems, ouPath)
			if err != nil {
				return nil, err
			}
			// Append child units to the main units slice
			units = append(units, childUnits...)
		}
	}

	return units, nil
}

func getOrganizationsAccountTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	resourceId := *h.Item.(OrgAccount).Id

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
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

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
