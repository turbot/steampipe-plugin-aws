package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func tableAwsDynamoDBGlobalTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_global_table",
		Description: "AWS DynamoDB Global Table",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("global_table_name"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			ItemFromKey:       globalTableFromKey,
			Hydrate:           getDynamboDbGlobalTable,
		},
		List: &plugin.ListConfig{
			Hydrate: listDynamboDbGlobalTables,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "global_table_name",
				Description: "The global table name",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "global_table_arn",
				Description: "The unique identifier of the global table",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamboDbGlobalTable,
			},
			{
				Name:        "global_table_status",
				Description: "The current state of the global table",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamboDbGlobalTable,
			},
			{
				Name:        "creation_date_time",
				Description: "The creation time of the global table",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getDynamboDbGlobalTable,
			},
			{
				Name:        "replication_group",
				Description: "The Regions where the global table has replicas",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamboDbGlobalTable,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GlobalTableName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamboDbGlobalTable,
				Transform:   transform.FromField("GlobalTableArn").Transform(arnToAkas),
			},
		}),
	}
}

//// ITEM FROM KEY

func globalTableFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["global_table_name"].GetStringValue()
	item := &dynamodb.GlobalTableDescription{
		GlobalTableName: &name,
	}

	return item, nil
}

//// LIST FUNCTION

func listDynamboDbGlobalTables(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listDynamboDbGlobalTables", "AWS_REGION", region)

	// Create Session
	svc, err := DynamoDbService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	tables, err := svc.ListGlobalTables(&dynamodb.ListGlobalTablesInput{})

	for _, globalTable := range tables.GlobalTables {
		d.StreamListItem(ctx, &dynamodb.GlobalTableDescription{
			GlobalTableName: globalTable.GlobalTableName,
		})
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDynamboDbGlobalTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	logger.Trace("getDynamboDbGlobalTable")
	globalTable := h.Item.(*dynamodb.GlobalTableDescription)
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listDynamboDbGlobalTables", "AWS_REGION", region)

	// Create Session
	svc, err := DynamoDbService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &dynamodb.DescribeGlobalTableInput{
		GlobalTableName: globalTable.GlobalTableName,
	}

	item, err := svc.DescribeGlobalTable(params)
	if err != nil {
		logger.Debug("[DEBUG] getDynamboDbGlobalTable__", "ERROR", err)
		return nil, err
	}

	return item.GlobalTableDescription, nil
}
