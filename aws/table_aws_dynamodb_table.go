package aws

import (
	"context"
	"errors"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsDynamoDBTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dynamodb_table",
		Description: "AWS DynamoDB Table",
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"TableNotFoundException"}),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getDynamoDBTable,
		},
		List: &plugin.ListConfig{
			Hydrate: listDynamoDBTables,
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
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("TableArn"),
			},
			{
				Name:        "table_id",
				Description: "Unique identifier for the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("TableId"),
			},
			{
				Name:        "creation_date_time",
				Description: "The date and time when the table was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("CreationDateTime"),
			},
			{
				Name:        "table_class",
				Description: "The table class of the specified table. Valid values are STANDARD and STANDARD_INFREQUENT_ACCESS.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("TableClassSummary.TableClass"),
			},
			{
				Name:        "table_status",
				Description: "The current state of the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("TableStatus"),
			},
			// If it is not available then it should default to "PROVISIONED"
			// Possible values are "PAY_PER_REQUEST" or "PROVISIONED"
			{
				Name:        "billing_mode",
				Description: "Controls how AWS charges for read and write throughput and manage capacity.",
				Type:        proto.ColumnType_STRING,
				Default:     "PROVISIONED",
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("BillingModeSummary.BillingMode").Transform(getTableBillingMode),
			},
			{
				Name:        "item_count",
				Description: "Number of items in the table. Note that this is an approximate value.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("ItemCount"),
			},
			{
				Name:        "global_table_version",
				Description: "Represents the version of global tables (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GlobalTables.html) in use, if the table is replicated across AWS Regions.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("GlobalTableVersion"),
			},
			{
				Name:        "read_capacity",
				Description: "The maximum number of strongly consistent reads consumed per second before DynamoDB returns a ThrottlingException.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("ProvisionedThroughput.ReadCapacityUnits"),
			},
			{
				Name:        "write_capacity",
				Description: "The maximum number of writes consumed per second before DynamoDB returns a ThrottlingException.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("ProvisionedThroughput.WriteCapacityUnits"),
			},
			{
				Name:        "latest_stream_arn",
				Description: "The Amazon Resource Name (ARN) that uniquely identifies the latest stream for this table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("LatestStreamArn"),
			},
			{
				Name:        "latest_stream_label",
				Description: "A timestamp, in ISO 8601 format, for this stream.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("LatestStreamLabel"),
			},
			{
				Name:        "table_size_bytes",
				Description: "Size of the table in bytes. Note that this is an approximate value.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("TableSizeBytes"),
			},
			{
				Name:        "archival_summary",
				Description: "Contains information about the table archive.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("ArchivalSummary"),
			},
			{
				Name:        "attribute_definitions",
				Description: "An array of AttributeDefinition objects. Each of these objects describes one attribute in the table and index key schema.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("AttributeDefinitions"),
			},
			{
				Name:        "key_schema",
				Description: "The primary key structure for the table.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("KeySchema"),
			},
			{
				Name:        "sse_description",
				Description: "The description of the server-side encryption status on the specified table.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDynamoDBTable,
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
				Hydrate:     getDynamoDBTable,
				Transform:   transform.FromField("TableArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDynamoDBTables(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := DynamoDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.listDynamoDBTables", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &dynamodb.ListTablesInput{
		Limit: aws.Int32(100),
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.ExclusiveStartTableName = aws.String(equalQuals["name"].GetStringValue())
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.Limit {
			if limit < 1 {
				input.Limit = aws.Int32(1)
			} else {
				input.Limit = aws.Int32(limit)
			}
		}
	}

	paginator := dynamodb.NewListTablesPaginator(svc, input, func(o *dynamodb.ListTablesPaginatorOptions) {
		o.Limit = *input.Limit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dynamodb_table.listDynamoDBTables", "api_error", err)
			return nil, err
		}

		for _, table := range output.TableNames {
			d.StreamListItem(ctx, types.TableDescription{
				TableName: aws.String(table),
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDynamoDBTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var name string
	if h.Item != nil {
		data := h.Item.(types.TableDescription)
		if data.TableName != nil {
			name = *data.TableName
		}
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := DynamoDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.getDynamoDBTable", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	}

	rowData, err := svc.DescribeTable(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.getDynamoDBTable", "api_error", err)
		return nil, err
	}

	if rowData.Table != nil {
		return *rowData.Table, nil
	}

	return nil, nil
}

func getDescribeContinuousBackups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	table := h.Item.(types.TableDescription)

	// Create Session
	svc, err := DynamoDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.getDescribeContinuousBackups", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &dynamodb.DescribeContinuousBackupsInput{
		TableName: table.TableName,
	}

	op, err := svc.DescribeContinuousBackups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.getDescribeContinuousBackups", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getTableTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	region := d.KeyColumnQualString(matrixKeyRegion)
	table := h.Item.(types.TableDescription)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	tableArn := "arn:" + commonColumnData.Partition + ":dynamodb:" + region + ":" + commonColumnData.AccountId + ":table/" + *table.TableName

	// Create Session
	svc, err := DynamoDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.getTableTagging", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &dynamodb.ListTagsOfResourceInput{
		ResourceArn: &tableArn,
	}

	pagesLeft := true
	tags := []types.Tag{}
	for pagesLeft {
		result, err := svc.ListTagsOfResource(ctx, params)
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) {
				// Added to support regex in not found errors
				if ok, _ := path.Match("ResourceNotFoundException", ae.ErrorCode()); ok {
					return nil, nil
				}
			}
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

func getTableStreamingDestination(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	tableName := h.Item.(types.TableDescription).TableName

	// Create Session
	svc, err := DynamoDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.getTableStreamingDestination", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &dynamodb.DescribeKinesisStreamingDestinationInput{
		TableName: tableName,
	}

	op, err := svc.DescribeKinesisStreamingDestination(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dynamodb_table.getTableStreamingDestination", "api_error", err)
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func getTableBillingMode(_ context.Context, d *transform.TransformData) (interface{}, error) {
	billingMode := types.BillingModeProvisioned
	if d.Value != nil {
		billingMode = d.Value.(types.BillingMode)
	}

	return billingMode, nil
}

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
