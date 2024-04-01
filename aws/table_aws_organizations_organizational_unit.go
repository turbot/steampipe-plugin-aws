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

func tableAwsOrganizationsOrganizationalUnit(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_organizations_organizational_unit",
		Description: "AWS Organizations Organizational Unit",
		List: &plugin.ListConfig{
			ParentHydrate: listOrganizationsRoots,
			Hydrate:       listOrganizationsOrganizationalUnits,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ParentNotFoundException", "InvalidInputException"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:      "parent_id",
					Require:   plugin.Optional,
					Operators: []string{"="},
				},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name of this OU.",
				Hydrate:     getOrganizationsOrganizationalUnit,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier (ID) associated with this OU.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of this OU.",
				Hydrate:     getOrganizationsOrganizationalUnit,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parent_id",
				Description: "The unique identifier (ID) of the root or OU whose child OUs you want to list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path",
				Description: "The OU path is a string representation that uniquely identifies the hierarchical location of an Organizational Unit within the AWS Organizations structure.",
				Type:        proto.ColumnType_LTREE,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOrganizationsOrganizationalUnit,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOrganizationsOrganizationalUnit,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type OrganizationalUnit struct {
	types.OrganizationalUnit
	Path     string
	ParentId string
}

//// LIST FUNCTION

func listOrganizationsOrganizationalUnits(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	parentId := *h.Item.(types.Root).Id

	// Check if the parentId is provided
	// The unique identifier (ID) of the root or OU whose child OUs you want to list.
	if d.EqualsQualString("parent_id") != "" {
		parentId = d.EqualsQualString("parent_id")
	}

	// empty check
	if parentId == "" {
		return nil, nil
	}

	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_organizational_unit.listOrganizationsOrganizationalUnits", "client_error", err)
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

	// Call the recursive function to list all nested OUs
	rootPath := parentId
	err = listAllNestedOUs(ctx, d, svc, parentId, maxItems, rootPath)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_organizational_unit.listOrganizationsOrganizationalUnits", "recursive_call_error", err)
		return nil, err
	}

	return nil, nil
}

func listAllNestedOUs(ctx context.Context, d *plugin.QueryData, svc *organizations.Client, parentId string, maxItems int32, currentPath string) error {
	params := &organizations.ListOrganizationalUnitsForParentInput{
		ParentId:   aws.String(parentId),
		MaxResults: &maxItems,
	}

	paginator := organizations.NewListOrganizationalUnitsForParentPaginator(svc, params, func(o *organizations.ListOrganizationalUnitsForParentPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, unit := range output.OrganizationalUnits {
			ouPath := strings.Replace(currentPath, "-", "_", -1) + "." + strings.Replace(*unit.Id, "-", "_", -1)
			d.StreamListItem(ctx, OrganizationalUnit{unit, ouPath, parentId})

			// Recursively list units for this child
			err := listAllNestedOUs(ctx, d, svc, *unit.Id, maxItems, ouPath)
			if err != nil {
				return err
			}

			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}
	}

	return nil
}

//// HYDRATE FUNCTIONS

func getOrganizationsOrganizationalUnit(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	orgUnitId := *h.Item.(OrganizationalUnit).Id

	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_organizational_unit.getOrganizationsOrganizationalUnit", "client_error", err)
		return nil, err
	}

	params := &organizations.DescribeOrganizationalUnitInput{
		OrganizationalUnitId: aws.String(orgUnitId),
	}

	op, err := svc.DescribeOrganizationalUnit(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_organizational_unit.getOrganizationsOrganizationalUnit", "api_error", err)
		return nil, err
	}

	return *op.OrganizationalUnit, nil
}
