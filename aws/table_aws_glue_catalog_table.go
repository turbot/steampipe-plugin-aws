package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueCatalogTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_catalog_table",
		Description: "AWS Glue Catalog Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "database_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlueCatalogTable,
		},
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "catalog_id", Require: plugin.Optional},
				{Name: "database_name", Require: plugin.Optional},
			},
			ParentHydrate: listGlueCatalogDatabases,
			Hydrate:       listGlueCatalogTables,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The table name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "catalog_id",
				Description: "The ID of the Data Catalog in which the table resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time when the table definition was created in the data catalog.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description of the table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_by",
				Description: "The person or entity who created the table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "database_name",
				Description: "The name of the database where the table metadata resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_registered_with_lake_formation",
				Description: "Indicates whether the table has been registered with lake formation.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_access_time",
				Description: "The last time that the table was accessed. This is usually taken from HDFS, and might not be reliable.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_analyzed_time",
				Description: "The last time that column statistics were computed for this table.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "owner",
				Description: "The owner of the table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "retention",
				Description: "The retention time for this table.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "table_type",
				Description: "The type of this table (EXTERNAL_TABLE, VIRTUAL_VIEW, etc.).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "The last time that the table was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "view_expanded_text",
				Description: "If the table is a view, the expanded text of the view otherwise null.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "view_original_text",
				Description: "If the table is a view, the original text of the view otherwise null.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parameters",
				Description: "These key-value pairs define properties associated with the table.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "partition_keys",
				Description: "A list of columns by which the table is partitioned.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "storage_descriptor",
				Description: "A storage descriptor containing information about the physical storage of this table.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "target_table",
				Description: "A TableIdentifier structure that describes a target table for resource linking.",
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
				Hydrate:     getGlueCatalogTableAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlueCatalogTables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Info("aws_glue_catalog_table.listGlueCatalogTables")
	database := h.Item.(types.Database)

	// Create session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		logger.Error("aws_glue_catalog_table.listGlueCatalogTables", "connection_error", err)
		return nil, err
	}

	if d.KeyColumnQualString("catalog_id") != "" && *database.CatalogId != d.KeyColumnQualString("catalog_id") {
		return nil, nil
	}
	if d.KeyColumnQualString("database_name") != "" && *database.Name != d.KeyColumnQualString("database_name") {
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(100)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit=1
			} else {
				maxLimit = int32(*limit) 
			}
		}
	}
	input := &glue.GetTablesInput{
		MaxResults: aws.Int32(maxLimit),
		DatabaseName: database.Name,
		CatalogId:    database.CatalogId,
	}
	// List call
	equalQuals := d.KeyColumnQuals
	if equalQuals["catalog_id"] != nil {
		input.CatalogId = aws.String(equalQuals["catalog_id"].GetStringValue())
	}
	paginator:=glue.NewGetTablesPaginator(svc,input,func(o *glue.GetTablesPaginatorOptions){
		o.Limit=maxLimit
		o.StopOnDuplicateToken=true
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_glue_catalog_table.listGlueCatalogDatabases", "api_error", err)
			return nil, err
		}
		for _, table  := range output.TableList{
			d.StreamListItem(ctx, table)
			plugin.Logger(ctx).Error("aws_glue_catalog_table.listGlueCatalogDatabases", "api_error", err)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil,nil
			}
		}
		
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getGlueCatalogTable(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()
	databaseName := d.KeyColumnQuals["database_name"].GetStringValue()

	// Empty check
	if name == "" || databaseName == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &glue.GetTableInput{
		Name:         aws.String(name),
		DatabaseName: aws.String(databaseName),
	}

	// Get call
	data, err := svc.GetTable(ctx,params)
	if err != nil {
		plugin.Logger(ctx).Error("getGlueCatalogTable", "ERROR", err)
		return nil, err
	}
	return data.Table, nil
}

func getGlueCatalogTableAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(types.Table)

	// Get common columns
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":table/" + *data.DatabaseName + "/" + *data.Name

	return []string{aka}, nil
}
