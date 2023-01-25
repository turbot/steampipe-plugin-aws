package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"

	gluev1 "github.com/aws/aws-sdk-go/service/glue"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueCatalogDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_catalog_database",
		Description: "AWS Glue Catalog Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlueCatalogDatabase,
		},
		List: &plugin.ListConfig{
			Hydrate: listGlueCatalogDatabases,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(gluev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the database. For Hive compatibility, this is folded to lowercase when it is stored.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "catalog_id",
				Description: "The ID of the Data Catalog in which the database resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time at which the metadata database was created in the catalog.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description of the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location_uri",
				Description: "The location of the database (for example, an HDFS path).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_table_default_permissions",
				Description: "Creates a set of default permissions on the table for principals.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "parameters",
				Description: "These key-value pairs define parameters and properties of the database.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "target_database",
				Description: "A DatabaseIdentifier structure that describes a target database for resource linking.",
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
				Hydrate:     getGlueCatalogDatabaseAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlueCatalogDatabases(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_catalog_database.listGlueCatalogDatabases", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(100)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	input := &glue.GetDatabasesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	// List call
	equalQuals := d.KeyColumnQuals
	if equalQuals["catalog_id"] != nil {
		input.CatalogId = aws.String(equalQuals["catalog_id"].GetStringValue())
	}
	paginator := glue.NewGetDatabasesPaginator(svc, input, func(o *glue.GetDatabasesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_glue_catalog_database.listGlueCatalogDatabases", "api_error", err)
			return nil, err
		}
		for _, database := range output.DatabaseList {
			d.StreamListItem(ctx, database)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueCatalogDatabase(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_catalog_database.getGlueCatalogDatabase", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &glue.GetDatabaseInput{
		Name: aws.String(name),
	}

	// Get call
	data, err := svc.GetDatabase(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_catalog_database.getGlueCatalogDatabase", "api_error", err)
		return nil, err
	}

	return *data.Database, nil
}

func getGlueCatalogDatabaseAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(types.Database)

	// Get common columns

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_catalog_database.getGlueCatalogDatabaseAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":database/" + *data.Name

	return []string{aka}, nil
}
