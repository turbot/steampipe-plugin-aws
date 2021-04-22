package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/glue"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsGlueDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_catalog_database",
		Description: "AWS Glue Database",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("Name"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidInputException", "EntityNotFoundException", "InternalServiceException", "OperationTimeoutException", "EncryptionException"}),
			Hydrate:           getAwsGlueDatabase,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsGlueDatabases,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "catalog_id",
				Description: "The ID of the Data Catalog in which the database resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_table_default_permissions",
				Description: "Creates a set of default permissions on the table for principals.",
				Type:        proto.ColumnType_JSON,
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
				Name:        "name",
				Description: "The name of the database. For Hive compatibility, this is folded to lowercase when it is stored.",
				Type:        proto.ColumnType_STRING,
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
			// Standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsGlueDatabaseAkas,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsGlueDatabases(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsGlueDatabases", "AWS_REGION", region)

	// Create session
	svc, err := GlueService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.GetDatabasesPages(
		&glue.GetDatabasesInput{},
		func(page *glue.GetDatabasesOutput, isLast bool) bool {
			for _, database := range page.DatabaseList {
				d.StreamListItem(ctx, database)

			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsGlueDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsGlueDatabase")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	var name string
	if h.Item != nil {
		name = databaseName(h.Item)
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := GlueService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &glue.GetDatabaseInput{
		Name: &name,
	}

	// Get call
	data, err := svc.GetDatabase(params)
	if err != nil {
		logger.Debug("getAwsGlueDatabase", "ERROR", err)
		return nil, err
	}

	return data.Database, nil
}

func getAwsGlueDatabaseAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsGlueDatabaseAkas")
	name := databaseName(h.Item)
	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":glue:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":database"

	if strings.HasPrefix(name, "/") {
		aka = aka + name
	} else {
		aka = aka + "/" + name
	}

	return []string{aka}, nil
}

func databaseName(item interface{}) string {
	switch item.(type) {
	case *glue.Database:
		return *item.(*glue.Database).Name
	}
	return ""
}
