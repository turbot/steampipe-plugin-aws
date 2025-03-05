package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lakeformation"
	"github.com/aws/aws-sdk-go-v2/service/lakeformation/types"
	lakeformationv1 "github.com/aws/aws-sdk-go/service/lakeformation"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLakeformationPermission(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lakeformation_permission",
		Description: "AWS Lake Formation Permission",
		List: &plugin.ListConfig{
			Hydrate: listLakeformationPermissions,
			// Future Reference:
			// The key column qualifiers can be included for "table_*", "table_with_*", "lf_tag_*", etc.
			// However, handling the values in the input parameters is complex due to their interdependencies.
			// The dependencies work as follows:
			//
			// - If `table_catalog_id` is provided, `table_database_name`, `table_name`, and the table wildcard
			//   must also be included in the query parameters.
			//
			// - If `table_database_name` is provided, `table_name` and the table wildcard
			//   must also be included in the query parameters.
			//
			// - If `table_name` is provided, the table wildcard must also be included in the query parameters.
			//
			// Omitting any required parameter may result in an incomplete or invalid query.
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "principal_identifier", Require: plugin.Optional},
				{Name: "database_catalog_id", Require: plugin.Optional},
				{Name: "database_name", Require: plugin.Optional},
				{Name: "data_location_catalog_id", Require: plugin.Optional},
				{Name: "data_location_resource_arn", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "lakeformation", "action": "ListPermissions"},
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

			// Condition
			{
				Name:        "condition_expression",
				Description: "The condition expression associated with the permission.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Condition.Expression"),
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

			// LF-Tag values
			{
				Name:        "lf_tag_values",
				Description: "The values of the LF tag resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Resource.LFTag.TagValues"),
			},

			// Additional Details
			{
				Name:        "additional_details_resource_share",
				Description: "Resource share details associated with the permission.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AdditionalDetails.ResourceShare"),
			},

			// Steampipe standard columns
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

	resourceInput := buildLakeformationResourceInputFilter(ctx, d.Quals)
	if resourceInput != nil {
		input.Resource = resourceInput
	}

	// Error: aws: operation error LakeFormation: ListPermissions, https response error StatusCode: 400, RequestID: 4dd01db3-77a7-4fcc-86d2-398e7f62fbcf, InvalidInputException: Resource is mandatory if Principal is set in the input.
	if d.EqualsQualString("principal_identifier") != "" && resourceInput != nil {
		input.Principal = &types.DataLakePrincipal{
			DataLakePrincipalIdentifier: aws.String(d.EqualsQualString("principal_identifier")),
		}
	}

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

func buildLakeformationResourceInputFilter(ctx context.Context, quals plugin.KeyColumnQualMap) (resource *types.Resource) {
	resourceInput := &types.Resource{}
	hasValues := false // Track if any value is set

	for columnName := range quals {
		if quals[columnName] != nil {
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if !ok {
				continue
			}

			switch columnName {
			case "database_catalog_id":
				// If we are providing database_catalog_id then we must have to provide database_name
				// Otherwise we will get the error:
				// Error: aws: operation error LakeFormation: ListPermissions, 1 validation error(s) found.
				// - missing required field, ListPermissionsInput.Resource.Database.Name.
				if quals["database_name"] != nil {
					if resourceInput.Database == nil {
						resourceInput.Database = &types.DatabaseResource{}
					}
					resourceInput.Database.CatalogId = aws.String(val)
					hasValues = true
				}
			case "database_name":
				if resourceInput.Database == nil {
					resourceInput.Database = &types.DatabaseResource{}
				}
				resourceInput.Database.Name = aws.String(val)
				hasValues = true

			// If we are providing data_location_catalog_id then we must have to provide data_location_resource_arn
			// Otherwise we will get the error:
			// Error: aws: operation error LakeFormation: ListPermissions, 1 validation error(s) found.
			// - missing required field, ListPermissionsInput.Resource.DataLocation.ResourceArn.
			//  (SQLSTATE HV000
			case "data_location_catalog_id":
				if quals["data_location_resource_arn"] != nil {
					if resourceInput.DataLocation == nil {
						resourceInput.DataLocation = &types.DataLocationResource{}
					}
					resourceInput.DataLocation.CatalogId = aws.String(val)
					hasValues = true
				}
			case "data_location_resource_arn":
				if resourceInput.DataLocation == nil {
					resourceInput.DataLocation = &types.DataLocationResource{}
				}
				resourceInput.DataLocation.ResourceArn = aws.String(val)
				hasValues = true
			}
		}
	}

	// Return the resource only if at least one field is set
	if hasValues {
		return resourceInput
	}

	return nil
}
