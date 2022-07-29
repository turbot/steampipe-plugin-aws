package aws

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func tableAwsDynamoDBTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_table",
		Description: "AWS DynamoDB Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getDynamboDbTable,
		},
		List: &plugin.ListConfig{
			Hydrate: listDynamboDbTables,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the table.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TableName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that uniquely identifies the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("TableArn"),
			},
			{
				Name:        "table_id",
				Description: "Unique identifier for the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("TableId"),
			},
			{
				Name:        "creation_date_time",
				Description: "The date and time when the table was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("CreationDateTime"),
			},
			{
				Name:        "table_class",
				Description: "The table class of the specified table. Valid values are STANDARD and STANDARD_INFREQUENT_ACCESS.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("TableClassSummary.TableClass"),
			},
			{
				Name:        "table_status",
				Description: "The current state of the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("TableStatus"),
			},
			{
				Name:        "billing_mode",
				Description: "Controls how AWS charges for read and write throughput and manage capacity.",
				Type:        proto.ColumnType_STRING,
				Default:     "PROVISIONED",
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("BillingModeSummary.BillingMode"),//.Transform(getTableBillingMode),
				// If it is not available then it should default to  "PROVISIONED"
				// Billing mode can only be PAY_PER_REQUEST or PROVISIONED
			},
			{
				Name:        "item_count",
				Description: "Number of items in the table. Note that this is an approximate value.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("ItemCount"),
			},
			{
				Name:        "global_table_version",
				Description: "Represents the version of global tables (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GlobalTables.html) in use, if the table is replicated across AWS Regions.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("GlobalTableVersion"),
			},
			{
				Name:        "read_capacity",
				Description: "The maximum number of strongly consistent reads consumed per second before DynamoDB returns a ThrottlingException.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("ProvisionedThroughput.ReadCapacityUnits"),
			},
			{
				Name:        "write_capacity",
				Description: "The maximum number of writes consumed per second before DynamoDB returns a ThrottlingException.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("ProvisionedThroughput.WriteCapacityUnits"),
			},
			{
				Name:        "latest_stream_arn",
				Description: "The Amazon Resource Name (ARN) that uniquely identifies the latest stream for this table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("LatestStreamArn"),
			},
			{
				Name:        "latest_stream_label",
				Description: "A timestamp, in ISO 8601 format, for this stream.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("LatestStreamLabel"),
			},
			{
				Name:        "table_size_bytes",
				Description: "Size of the table in bytes. Note that this is an approximate value.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("TableSizeBytes"),
			},
			{
				Name:        "archival_summary",
				Description: "Contains information about the table archive.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("ArchivalSummary"),
			},
			{
				Name:        "attribute_definitions",
				Description: "An array of AttributeDefinition objects. Each of these objects describes one attribute in the table and index key schema.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("AttributeDefinitions"),
			},
			{
				Name:        "key_schema",
				Description: "The primary key structure for the table.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("KeySchema"),
			},
			{
				Name:        "sse_description",
				Description: "The description of the server-side encryption status on the specified table.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("SSEDescription"),
			},
			{
				Name:        "continuous_backups_status",
				Description: "The continuous backups status of the table. ContinuousBackupsStatus can be one of the following states: ENABLED, DISABLED.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDescribeContinuousBackups,
				Transform:   transform.FromField("ContinuousBackupsDescription.ContinuousBackupsStatus"),
			},
			{
				Name:        "point_in_time_recovery_description",
				Description: "The description of the point in time recovery settings applied to the table.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDescribeContinuousBackups,
				Transform:   transform.FromField("ContinuousBackupsDescription.PointInTimeRecoveryDescription"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the table.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTableTagging,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTableTagging,
				Transform:   transform.From(getTableTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TableName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamboDbTable,
				Transform:   transform.FromField("TableArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDynamboDbTables(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := DynamoDbClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.listDynamboDbTables", "service_client_error", err)
		return nil, err
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

	input := &dynamodb.ListTablesInput{
		Limit: aws.Int32(maxLimit),
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.ExclusiveStartTableName = aws.String(equalQuals["name"].GetStringValue())
	}

	paginator := dynamodb.NewListTablesPaginator(svc, input, func(o *dynamodb.ListTablesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dynamodb_table.listDynamboDbTables", "api_error", err)
			return nil, err
		}

		for _, items := range output.TableNames {
			d.StreamListItem(ctx, &types.TableDescription{
				TableName: aws.String(items),
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDynamboDbTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var name string
	if h.Item != nil {
		name = *h.Item.(*types.TableDescription).TableName
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := DynamoDbClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.getDynamboDbTable", "service_client_error", err)
		return nil, err
	}

	params := &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	}

	rowData, err := svc.DescribeTable(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.getDynamboDbTable", "api_error", err)
		return nil, err
	}

	if rowData.Table != nil {
		return rowData.Table, nil
	}

	return nil, nil
}

func getDescribeContinuousBackups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	table := *h.Item.(*types.TableDescription).TableName

	// Create Session
	svc, err := DynamoDbClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.getDescribeContinuousBackups", "service_client_error", err)
		return nil, err
	}

	params := &dynamodb.DescribeContinuousBackupsInput{
		TableName: aws.String(table),
	}

	op, err := svc.DescribeContinuousBackups(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "TableNotFoundException") {
			// If a table is archived then continuous backups can't be queried for it
			return dynamodb.DescribeContinuousBackupsOutput{}, nil
		}
		plugin.Logger(ctx).Error("aws_dynamodb_table.getDescribeContinuousBackups", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getTableTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	tableName := *h.Item.(*types.TableDescription).TableName

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Create Session
	svc, err := DynamoDbClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.getTableTagging", "service_client_error", err)
		return nil, err
	}

	tableArn := "arn:" + commonColumnData.Partition + ":dynamodb:" + region + ":" + commonColumnData.AccountId + ":table/" + tableName

	params := &dynamodb.ListTagsOfResourceInput{
		ResourceArn: &tableArn,
	}

	pagesLeft := true
	tags := []types.Tag{}
	for pagesLeft {
		result, err := svc.ListTagsOfResource(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dynamodb_table.getTableTagging", "api_error", err)
			return nil, err
		}
		tags = append(tags, result.Tags...)
		if result.NextToken != nil {
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return tags, nil
}

//// TRANSFORM FUNCTIONS

func getTableTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
