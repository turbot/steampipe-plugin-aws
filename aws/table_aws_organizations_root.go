package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// The table will return the result if the account is a member of an organization.
// You must use the credentials of an account that belongs to an organization.
// The table will return an empty row if the account isn't a member of an organization instead of AWSOrganizationsNotInUseException.
func tableAwsOrganizationsRoot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_organizations_root",
		Description: "AWS Organizations Root",
		List: &plugin.ListConfig{
			Hydrate: listOrganizationsRoots,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"AWSOrganizationsNotInUseException"}),
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name of the root.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier (ID) for the root.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the root.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy_types",
				Description: "The types of policies that are currently enabled for the root and therefore can be attached to the root or to its OUs or accounts.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
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

func listOrganizationsRoots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_root.listOrganizationsRoots", "client_error", err)
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

	params := &organizations.ListRootsInput{
		MaxResults: &maxItems,
	}

	paginator := organizations.NewListRootsPaginator(svc, params, func(o *organizations.ListRootsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
	// apply rate limiting
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_organizations_root.listOrganizationsRoots", "api_error", err)
			return nil, err
		}

		for _, root := range output.Roots {
			d.StreamListItem(ctx, root)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
