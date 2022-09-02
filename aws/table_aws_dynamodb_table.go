package aws

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func tableAwsDynamoDBTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_table",
		Description: "AWS DynamoDB Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException"}),
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
		GetMatrixItemFunc: BuildRegionList,
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
				Transform:   transform.FromField("BillingModeSummary.BillingMode").Transform(getTableBillingMode),
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
				Name:        "streaming_destination",
				Description: "Provides information about the status of Kinesis streaming.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTableStreamingDestination,
				Transform:   transform.FromValue(),
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
	svc, err := DynamoDbService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ListTablesInput{
		Limit: aws.Int64(100),
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.ExclusiveStartTableName = types.String(equalQuals["name"].GetStringValue())
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.Limit {
			if *limit < 1 {
				input.Limit = types.Int64(1)
			} else {
				input.Limit = limit
			}
		}
	}

	err = svc.ListTablesPages(
		input,
		func(page *dynamodb.ListTablesOutput, lastPage bool) bool {
			for _, table := range page.TableNames {
				d.StreamListItem(ctx, &dynamodb.TableDescription{
					TableName: table,
				})

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("ListTablesPages", "list", err)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDynamboDbTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDynamboDbTable")

	var name string
	if h.Item != nil {
		data := h.Item.(*dynamodb.TableDescription)
		name = types.SafeString(data.TableName)
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := DynamoDbService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	}

	rowData, err := svc.DescribeTable(params)
	if err != nil {
		plugin.Logger(ctx).Debug("[DEBUG] getDynamboDbTable__", "ERROR", err)
		return nil, err
	}

	if rowData.Table != nil {
		return rowData.Table, nil
	}

	return nil, nil
}

func getDescribeContinuousBackups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDescribeContinuousBackups")
	table := h.Item.(*dynamodb.TableDescription)

	// Create Session
	svc, err := DynamoDbService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &dynamodb.DescribeContinuousBackupsInput{
		TableName: table.TableName,
	}

	op, err := svc.DescribeContinuousBackups(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			// If a table is archived then continuous backups can't be queried for it
			if a.Code() == "TableNotFoundException" {
				return dynamodb.DescribeContinuousBackupsOutput{}, nil
			}
		}
		return nil, err
	}

	return op, nil
}

func getTableTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getTableTagging")
	region := d.KeyColumnQualString(matrixKeyRegion)
	table := h.Item.(*dynamodb.TableDescription)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Create Session
	svc, err := DynamoDbService(ctx, d)
	if err != nil {
		return nil, err
	}

	tableArn := "arn:" + commonColumnData.Partition + ":dynamodb:" + region + ":" + commonColumnData.AccountId + ":table/" + *table.TableName

	params := &dynamodb.ListTagsOfResourceInput{
		ResourceArn: &tableArn,
	}

	pagesLeft := true
	tags := []*dynamodb.Tag{}
	for pagesLeft {
		result, err := svc.ListTagsOfResource(params)
		if err != nil {
			plugin.Logger(ctx).Error("ListTagsOfResource", "tag", err)
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

func getTableStreamingDestination(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	tableName := h.Item.(*dynamodb.TableDescription).TableName

	// Create Session
	svc, err := DynamoDbService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DescribeKinesisStreamingDestinationInput{
		TableName: tableName,
	}

	op, err := svc.DescribeKinesisStreamingDestination(input)

	if err != nil {
		plugin.Logger(ctx).Error("getTableStreamingDestination", err)
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func getTableBillingMode(_ context.Context, d *transform.TransformData) (interface{}, error) {
	billingMode := "PROVISIONED"
	if d.Value != nil {
		billingMode = *d.Value.(*string)
	}

	return billingMode, nil
}

func getTableTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.([]*dynamodb.Tag)

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
