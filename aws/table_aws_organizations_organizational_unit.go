package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// The table will return the Organizational Units for the root account if parent_id is not specified in the query parameter.
// If parent_id is specified in the query parameter then it will return the Organizational Units for the given parent.
func tableAwsOrganizationsOrganizationalUnit(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_organizations_organizational_unit",
		Description: "AWS Organizations Organizational Unit",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOrganizationsOrganizationalUnit,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException"}),
			},
		},
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
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parent_id",
				Description: "The unique identifier (ID) of the root or OU whose child OUs you want to list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getParentId),
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

func listOrganizationsOrganizationalUnits(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	parentId := d.EqualsQualString("parent_id")

	if parentId == "" && h.Item != nil {
		parentId = *h.Item.(types.Root).Id
	}

	// Empty Check
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

	params := &organizations.ListOrganizationalUnitsForParentInput{
		ParentId:   aws.String(parentId),
		MaxResults: &maxItems,
	}

	paginator := organizations.NewListOrganizationalUnitsForParentPaginator(svc, params, func(o *organizations.ListOrganizationalUnitsForParentPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if ae.ErrorCode() == "ParentNotFoundException" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_organizations_organizational_unit.listOrganizationsOrganizationalUnits", "api_error", err)
			return nil, err
		}

		for _, unit := range output.OrganizationalUnits {
			d.StreamListItem(ctx, unit)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOrganizationsOrganizationalUnit(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var orgUnitId string

	if h.Item != nil {
		orgUnitId = *h.Item.(types.OrganizationalUnit).Id
	} else {
		orgUnitId = d.EqualsQuals["id"].GetStringValue()
	}

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

//// TRANSFORM FUNCTION

// This function will be useful if user query the table with 'NOT' operator for the optionsl qual 'parent_id'. Like: "select * from aws_organizations_organizational_unit where parent_id <> 'ou-skjaa-siiewfhgw'"
func getParentId(_ context.Context, d *transform.TransformData) (interface{}, error) {
	quals := d.KeyColumnQuals["parent_id"]
	for _, data := range quals {
		parentId := data.Value.GetStringValue()
		if parentId != "" && data.Operator == "=" {
			return parentId, nil
		}
	}

	if d.HydrateItem != nil {
		data := d.HydrateItem.(types.OrganizationalUnit)
		return strings.Split(*data.Arn, "/")[2], nil
	}

	return nil, nil
}
