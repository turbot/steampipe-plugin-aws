package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsOrganizationsDelegatedServicesForAccount(_ context.Context) *plugin.Table {

	return &plugin.Table{
		Name:        "aws_organizations_delegated_services_for_account",
		Description: "AWS Organizations Delegated Services For Account",
		List: &plugin.ListConfig{
			ParentHydrate: listOrganizationsDelegatedAdmins, // Use Delegated Administrator as parent per recommendation. Referenced table_aws_cloudwatch_log_stream. This function can be found in table_aws_organizations_delegated_administrator
			Hydrate:       listDelegatedServices,
			Tags:          map[string]string{"service": "organizations", "action": "ListDelegatedServicesForAccount"},
			KeyColumns: []*plugin.KeyColumn{ // Make delegated_account_id optional, user can still query `where` using this column.
				{Name: "delegated_account_id", Require: plugin.Optional},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "delegated_account_id",
				Description: "The unique identifier (account ID) of the delegated administrator account for which services are listed.",
				Type:        proto.ColumnType_STRING,
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
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServicePrincipal"),
			},
		}),
	}
}

// Define a struct to hold the DelegatedService and the AccountId used to fetch it.\
type delegatedServiceInfo struct {
	types.DelegatedService
	DelegatedAccountId string
}


func listDelegatedServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	delegatedAccount := h.Item.(types.DelegatedAdministrator) // Get delegated administrator details

	delegatedAccountId := d.EqualsQualString("delegated_account_id") // Get delegated_account_id from where statement

	// Minimize API calls
	if delegatedAccountId != "" {
		if delegatedAccountId != *delegatedAccount.Id {
			return nil, nil
		}
	}

	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_delegated_services_for_account.ListDelegatedServicesForAccount", "client_error", err)
		return nil, err
	}

	// Limiting the result
	maxItems := int32(20)

	// Reduce the page size if a smaller limit is provided
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}
	params := &organizations.ListDelegatedServicesForAccountInput{
		AccountId:  delegatedAccount.Id,
		MaxResults: &maxItems,
	}

	paginator := organizations.NewListDelegatedServicesForAccountPaginator(svc, params, func(o *organizations.ListDelegatedServicesForAccountPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_organizations_delegated_services_for_account.ListDelegatedServicesForAccount", "api_error", err)
			return nil, err
		}

		for _, service := range output.DelegatedServices {
			// Stream a new struct that includes the AccountId used for the API call
			d.StreamListItem(ctx, delegatedServiceInfo{
				DelegatedService:   service,
				DelegatedAccountId: *delegatedAccount.Id,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}