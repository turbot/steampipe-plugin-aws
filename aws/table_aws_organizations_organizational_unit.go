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
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOrganizationsOrganizationalUnit,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate:    listOrganizationsOrganizationalUnits,
			KeyColumns: plugin.SingleColumn("parent_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ParentNotFoundException", "InvalidInputException"}),
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
				Hydrate:     getOrganizationsOrganizationalUnit,
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

func listOrganizationsOrganizationalUnits(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get Client
	svc, err := OrganizationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_organizations_organizational_unit.listOrganizationsOrganizationalUnits", "client_error", err)
		return nil, err
	}

	parentId := d.EqualsQualString("parent_id")

	// Empty Check
	if parentId == "" {
		return nil, nil
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
		orgUnitId = *h.Item.(*types.Policy).PolicySummary.Id
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

func getParentId(_ context.Context, d *transform.TransformData) (interface{}, error) {
	quals := d.KeyColumnQuals["parent_id"]
	for _, data := range quals {
		parentId := data.Value.GetStringValue()
		if parentId != "" {
			return parentId, nil
		}
	}

	if d.HydrateItem != nil {
		data := d.HydrateItem.(*types.OrganizationalUnit)
		return strings.Split(*data.Arn, "/")[2], nil
	}

	return nil, nil
}
