package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"

	athenav1 "github.com/aws/aws-sdk-go/service/athena"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsAthenaDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_athena_database",
		Description: "AWS Athena Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "catalog_name"}),
			Hydrate:    getAwsAthenaDatabase,
			Tags:       map[string]string{"service": "athena", "action": "GetDatabase"},
		},
		List: &plugin.ListConfig{
			Hydrate:       listAwsAthenaDatabases,
			Tags:          map[string]string{"service": "athena", "action": "ListDatabases"},
			ParentHydrate: listAwsAthenaDataCatalogs,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(athenav1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Name"),
			},
			{
				Name:        "catalog_name",
				Description: "The name of the data catalog that contains the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CatalogName"),
			},
			{
				Name:        "description",
				Description: "A description of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Description"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAthenaDatabaseAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsAthenaDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := AthenaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_athena_database.listAwsAthenaDatabases", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Get the parent data catalog
	if h.Item == nil {
		return nil, nil
	}
	catalog := h.Item.(types.DataCatalogSummary)
	catalogName := *catalog.CatalogName

	maxResults := int32(50)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxResults {
			maxResults = limit
		}
	}

	input := &athena.ListDatabasesInput{
		MaxResults:  &maxResults,
		CatalogName: aws.String(catalogName),
	}

	paginator := athena.NewListDatabasesPaginator(svc, input, func(o *athena.ListDatabasesPaginatorOptions) {
		o.Limit = maxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_athena_database.listAwsAthenaDatabases", "api_error", err)
			return nil, err
		}

		for _, database := range output.DatabaseList {
			// Get the full database details
			params := &athena.GetDatabaseInput{
				DatabaseName: database.Name,
				CatalogName:  aws.String(catalogName),
			}

			op, err := svc.GetDatabase(ctx, params)
			if err != nil {
				plugin.Logger(ctx).Error("aws_athena_database.listAwsAthenaDatabases", "get_database_error", err, "database_name", *database.Name)
				continue
			}

			// Create a map to hold the database data with catalog name
			dbData := map[string]interface{}{
				"Database":    op.Database,
				"CatalogName": catalogName,
			}
			d.StreamListItem(ctx, dbData)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsAthenaDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := AthenaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_athena_database.getAwsAthenaDatabase", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var name, catalogName string
	if h != nil && h.Item != nil {
		// We're being called as part of a parent hydrate
		catalog := h.Item.(types.DataCatalogSummary)
		catalogName = *catalog.CatalogName
		name = d.EqualsQualString("name")
	} else {
		// Direct get operation
		name = d.EqualsQualString("name")
		catalogName = d.EqualsQualString("catalog_name")
	}

	// Empty check
	if name == "" || catalogName == "" {
		return nil, nil
	}

	params := &athena.GetDatabaseInput{
		DatabaseName: aws.String(name),
		CatalogName:  aws.String(catalogName),
	}

	op, err := svc.GetDatabase(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_athena_database.getAwsAthenaDatabase", "api_error", err, "params", params)
		return nil, err
	}

	return op.Database, nil
}

func getAwsAthenaDatabaseAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	// Get the database data from the hydrate data
	data, ok := h.Item.(map[string]interface{})
	if !ok {
		plugin.Logger(ctx).Error("aws_athena_database.getAwsAthenaDatabaseAkas", "invalid_hydrate_data", "Expected map[string]interface{}")
		return nil, nil
	}

	database, ok := data["Database"].(*types.Database)
	if !ok {
		plugin.Logger(ctx).Error("aws_athena_database.getAwsAthenaDatabaseAkas", "invalid_database_data", "Expected *types.Database")
		return nil, nil
	}

	catalogName, ok := data["CatalogName"].(string)
	if !ok {
		plugin.Logger(ctx).Error("aws_athena_database.getAwsAthenaDatabaseAkas", "invalid_catalog_name", "Expected string")
		return nil, nil
	}

	// Get common columns that will be returned for all resources
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Build ARN
	arn := "arn:" + commonColumnData.Partition + ":athena:" + region + ":" + commonColumnData.AccountId + ":datacatalog/" + catalogName + "/database/" + *database.Name

	return []string{arn}, nil
}
