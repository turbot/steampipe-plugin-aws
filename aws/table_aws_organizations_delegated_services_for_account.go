package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsOrganizationsDelegatedServicesForAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_organizations_delegated_services_for_account",
		Description: "AWS Organizations Delegated Services For Account",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("delegated_account_id"),
			Hydrate:    listDelegatedServices,
			Tags:       map[string]string{"service": "organizations", "action": "ListDelegatedServicesForAccount"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "delegated_account_id",
				Description: "The unique identifier (account ID) of the delegated administrator.",
				Type:        proto.ColumnType_STRING,
				Transform: 	 transform.FromQual("delegated_account_id"),
			},
			{
				Name:        "service_principal",
				Description: "The service principal delegated to the administrator account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delegation_enabled_date",
				Description: "The date when the delegation was enabled.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			// Standard columns for all tables
			// TODO
			// I am unsure whether the title and akas below should correspond to 'ServicePrincipal'.
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServicePrincipal"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServicePrincipal").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

func listDelegatedServices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Retrieve the `delegated_account_id` from the user's `WHERE` statement
	delegatedAccountId := d.EqualsQuals["delegated_account_id"].GetStringValue()

	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_delegated_services_for_account.ListDelegatedServicesForAccount", "client_error", err)
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

	params := &organizations.ListDelegatedServicesForAccountInput{
		AccountId:  aws.String(delegatedAccountId),
		MaxResults: &maxItems,
	}

	paginator := organizations.NewListDelegatedServicesForAccountPaginator(svc, params, func(o *organizations.ListDelegatedServicesForAccountPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_organizations_delegated_services_for_account.ListDelegatedServicesForAccount", "api_error", err)
			return nil, err
		}

		for _, service := range output.DelegatedServices {
			d.StreamListItem(ctx, service)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
