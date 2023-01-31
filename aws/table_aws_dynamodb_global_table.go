package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	dynamodbv1 "github.com/aws/aws-sdk-go/service/dynamodb"

	go_kit_pack "github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsDynamoDBGlobalTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_global_table",
		Description: "AWS DynamoDB Global Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("global_table_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getDynamoDBGlobalTable,
		},
		List: &plugin.ListConfig{
			Hydrate: listDynamoDBGlobalTables,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "global_table_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(dynamodbv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "global_table_name",
				Description: "The global table name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "global_table_arn",
				Description: "The unique identifier of the global table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamoDBGlobalTable,
			},
			{
				Name:        "global_table_status",
				Description: "The current state of the global table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamoDBGlobalTable,
			},
			{
				Name:        "creation_date_time",
				Description: "The creation time of the global table.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getDynamoDBGlobalTable,
			},
			{
				Name:        "replication_group",
				Description: "The Regions where the global table has replicas.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamoDBGlobalTable,
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
				Hydrate:     getDynamoDBGlobalTable,
				Transform:   transform.FromField("GlobalTableArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listDynamoDBGlobalTables(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := DynamoDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_global_table.listDynamoDBGlobalTables", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &dynamodb.ListGlobalTablesInput{
		Limit: aws.Int32(maxLimit),
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["global_table_name"] != nil {
		input.ExclusiveStartGlobalTableName = go_kit_pack.String(equalQuals["global_table_name"].GetStringValue())
	}

	tables, err := svc.ListGlobalTables(ctx, input)

	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_global_table.listDynamoDBGlobalTables", "api_error", err)
		return nil, err
	}

	for _, globalTable := range tables.GlobalTables {
		d.StreamListItem(ctx, types.GlobalTableDescription{
			GlobalTableName: globalTable.GlobalTableName,
		})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			break
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDynamoDBGlobalTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var name string
	if h.Item != nil {
		data := h.Item.(types.GlobalTableDescription)
		name = go_kit_pack.SafeString(data.GlobalTableName)
	} else {
		name = d.KeyColumnQuals["global_table_name"].GetStringValue()
	}

	// Create Session
	svc, err := DynamoDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_global_table.getDynamoDBGlobalTable", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &dynamodb.DescribeGlobalTableInput{
		GlobalTableName: aws.String(name),
	}

	item, err := svc.DescribeGlobalTable(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_global_table.getDynamoDBGlobalTable", "api_error", err)
		return nil, err
	}

	return item.GlobalTableDescription, nil
}
