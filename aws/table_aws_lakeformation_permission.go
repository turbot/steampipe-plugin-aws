package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lakeformation"
	lakeformationv1 "github.com/aws/aws-sdk-go/service/lakeformation"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLakeformationPermission(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lakeformation_permission",
		Description: "AWS Lake Formation Permissions.",
		List: &plugin.ListConfig{
			Hydrate: listLakeformationPermissions,
			Tags:    map[string]string{"service": "lakeformation", "action": "ListPermissions"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(lakeformationv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			// Principal Information
			{
				Name:        "principal_identifier",
				Description: "The identifier of the principal to whom permissions are granted.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Principal.DataLakePrincipalIdentifier"),
			},

			// Permissions
			{
				Name:        "permissions",
				Description: "The permissions granted to the principal.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "permissions_with_grant_option",
				Description: "The permissions granted with the grant option.",
				Type:        proto.ColumnType_JSON,
			},

			// Condition
			{
				Name:        "condition_expression",
				Description: "The condition expression associated with the permission.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Condition.Expression"),
			},

			// Additional Details
			{
				Name:        "additional_details_resource_share",
				Description: "Resource share details associated with the permission.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AdditionalDetails.ResourceShare"),
			},

			// Last Updated
			{
				Name:        "last_updated",
				Description: "The timestamp of when the permission was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastUpdated"),
			},
			{
				Name:        "last_updated_by",
				Description: "The principal that last updated the permission.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LastUpdatedBy"),
			},

			// Resource Information
			{
				Name:        "resource_catalog_id",
				Description: "The catalog ID associated with the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.Catalog.Id"),
			},

			// Database Resource
			{
				Name:        "database_catalog_id",
				Description: "The catalog ID associated with the database resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.Database.CatalogId"),
			},
			{
				Name:        "database_name",
				Description: "The name of the database resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.Database.Name"),
			},

			// Table Resource
			{
				Name:        "table_catalog_id",
				Description: "The catalog ID associated with the table resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.Table.CatalogId"),
			},
			{
				Name:        "table_database_name",
				Description: "The database name associated with the table resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.Table.DatabaseName"),
			},
			{
				Name:        "table_name",
				Description: "The name of the table resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.Table.Name"),
			},

			// Table with Columns Resource
			{
				Name:        "table_with_columns_catalog_id",
				Description: "The catalog ID associated with the table with columns resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.TableWithColumns.CatalogId"),
			},
			{
				Name:        "table_with_columns_database_name",
				Description: "The database name associated with the table with columns resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.TableWithColumns.DatabaseName"),
			},
			{
				Name:        "table_with_columns_name",
				Description: "The name of the table with columns resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.TableWithColumns.Name"),
			},
			{
				Name:        "table_with_columns_column_names",
				Description: "The list of column names associated with the table with columns resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Resource.TableWithColumns.ColumnNames"),
			},

			// Data Location Resource
			{
				Name:        "data_location_catalog_id",
				Description: "The catalog ID associated with the data location resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.DataLocation.CatalogId"),
			},
			{
				Name:        "data_location_resource_arn",
				Description: "The ARN of the data location resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.DataLocation.ResourceArn"),
			},

			// LF Tag Resource
			{
				Name:        "lf_tag_catalog_id",
				Description: "The catalog ID associated with the LF tag resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.LFTag.CatalogId"),
			},
			{
				Name:        "lf_tag_key",
				Description: "The key of the LF tag resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.LFTag.TagKey"),
			},
			{
				Name:        "lf_tag_values",
				Description: "The values of the LF tag resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Resource.LFTag.TagValues"),
			},

			// LF Tag Expression
			{
				Name:        "lf_tag_expression_catalog_id",
				Description: "The catalog ID associated with the LF tag expression resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.LFTagExpression.CatalogId"),
			},
			{
				Name:        "lf_tag_expression_name",
				Description: "The name of the LF tag expression resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.LFTagExpression.Name"),
			},

			/// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Principal.DataLakePrincipalIdentifier"),
			},
		}),
	}
}

//// LIST FUNCTION

func listLakeformationPermissions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := LakeFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lakeformation_permission.listLakeformationPermissions", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(100)
	input := &lakeformation.ListPermissionsInput{}

	// Reduce the request limit based on user input
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
	input.MaxResults = aws.Int32(maxItems)
	paginator := lakeformation.NewListPermissionsPaginator(svc, input, func(o *lakeformation.ListPermissionsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// Apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lakeformation_permission.listLakeformationPermissions", "api_error", err)
			return nil, err
		}

		for _, permission := range output.PrincipalResourcePermissions {
			d.StreamListItem(ctx, permission)

			// Stop if the context is cancelled
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
