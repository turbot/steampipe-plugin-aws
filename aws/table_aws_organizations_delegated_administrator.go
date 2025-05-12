package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsOrganizationsDelegatedAdministrator(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_organizations_delegated_administrator",
		Description: "AWS Organizations Delegated Administrator",
		List: &plugin.ListConfig{
			Hydrate: listOrganizationsDelegatedAdmins,
			Tags:    map[string]string{"service": "organizations", "action": "ListDelegatedAdministrators"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier (account ID) of the delegated administrator.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the delegated administrator.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email",
				Description: "The email address associated with the delegated administrator account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The friendly name of the delegated administrator account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the delegated administrator.",
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
				Name:        "delegation_enabled_date",
				Description: "The date when the delegation was enabled.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Standard columns for all tables
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

func listOrganizationsDelegatedAdmins(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_delegated_administrator.ListDelegatedAdministrators", "client_error", err)
		return nil, err
	}

	// Limiting the result
	maxItems := int32(20)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	params := &organizations.ListDelegatedAdministratorsInput{
		MaxResults: &maxItems,
	}

	paginator := organizations.NewListDelegatedAdministratorsPaginator(svc, params, func(o *organizations.ListDelegatedAdministratorsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_organizations_delegated_administrator.ListDelegatedAdministrators", "api_error", err)
			return nil, err
		}

		for _, admin := range output.DelegatedAdministrators {
			d.StreamListItem(ctx, admin)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil

}
