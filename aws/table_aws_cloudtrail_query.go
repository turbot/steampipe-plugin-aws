package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudTrailQuery(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudtrail_query",
		Description: "AWS CloudTrail Query",
		Get: &plugin.GetConfig{
			Hydrate:    getCloudTrailLakeQuery,
			KeyColumns: plugin.AllColumns([]string{"event_data_store_arn", "query_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EventDataStoreNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCloudTrailEventDataStores,
			Hydrate:       listCloudTrailLakeQueries,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "event_data_store_arn",
					Require: plugin.Optional,
				},
				{
					Name:    "query_status",
					Require: plugin.Optional,
				},
				{
					Name:      "creation_time",
					Require:   plugin.Optional,
					Operators: []string{"=", "<=", "<", ">", ">="},
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "query_id",
				Description: "The ID of the query.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_data_store_arn",
				Description: "The ID of the event data store.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time of the query.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationTime", "QueryStatistics.CreationTime"),
			},
			{
				Name:        "query_status",
				Description: "The status of a query. Values for QueryStatus include QUEUED, RUNNING, FINISHED, FAILED, TIMED_OUT, or CANCELLED.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bytes_scanned",
				Description: "Gets metadata about a query, including the number of events that were matched, the total number of events scanned, the query run time in milliseconds, and the query's creation time.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudTrailLakeQuery,
				Transform:   transform.FromField("QueryStatistics.BytesScanned"),
			},
			{
				Name:        "events_matched",
				Description: "The number of events that matched a query.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudTrailLakeQuery,
				Transform:   transform.FromField("QueryStatistics.EventsMatched"),
			},
			{
				Name:        "events_scanned",
				Description: "The number of events that the query scanned in the event data store.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudTrailLakeQuery,
				Transform:   transform.FromField("QueryStatistics.EventsScanned"),
			},
			{
				Name:        "execution_time_in_millis",
				Description: "The query's run time, in milliseconds.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudTrailLakeQuery,
				Transform:   transform.FromField("QueryStatistics.ExecutionTimeInMillis"),
			},
			{
				Name:        "query_string",
				Description: "The SQL code of a query.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudTrailLakeQuery,
			},

			// steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("QueryId"),
			},
		}),
	}
}

// We need custom structure
// 1. EventDataStoreArn is required for joing the queries
// 2. list/get call do not provide the EventDataStoreArn
// 3. Both list/get API call return different stracture with dfferent data
type QueryInfo struct {
	EventDataStoreArn *string
	CreationTime      *time.Time
	QueryId           *string
	QueryStatus       types.QueryStatus
	QueryString       *string
	QueryStatistics   *types.QueryStatisticsForDescribeQuery
}

//// LIST FUNCTION

func listCloudTrailLakeQueries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	eventDataStore := h.Item.(types.EventDataStore)

	if d.KeyColumnQualString("event_data_store_arn") != "" {
		if d.KeyColumnQualString("event_data_store_arn") != *eventDataStore.EventDataStoreArn {
			return nil, nil
		}
	}

	// Get client
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Info("aws_cloudtrail_query.listCloudTrailLakeQueries", "client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit

		}
	}

	input := &cloudtrail.ListQueriesInput{
		MaxResults:     aws.Int32(maxLimit),
		EventDataStore: eventDataStore.EventDataStoreArn,
	}

	if d.KeyColumnQualString("query_status") != "" {
		input.QueryStatus = types.QueryStatus(d.KeyColumnQualString("query_status"))
	}

	if d.Quals["creation_time"] != nil {
		for _, q := range d.Quals["creation_time"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case "<", "<=", "=":
				input.EndTime = aws.Time(timestamp)
			case ">", ">=":
				input.StartTime = aws.Time(timestamp)
			}
		}
	}

	paginator := cloudtrail.NewListQueriesPaginator(svc, input, func(o *cloudtrail.ListQueriesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	if paginator.HasMorePages() {
		op, err := paginator.NextPage(ctx)
		if err != nil {
			// You cannot act on an event data store that is inactive. This error could not be caught by configuring it in ignore config
			if strings.Contains(err.Error(), "InactiveEventDataStoreException") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_cloudtrail_query.listCloudTrailLakeQueries", "api_error", err)
			return nil, err
		}

		for _, item := range op.Queries {
			d.StreamListItem(ctx, &QueryInfo{eventDataStore.EventDataStoreArn, item.CreationTime, item.QueryId, item.QueryStatus, aws.String(""), &types.QueryStatisticsForDescribeQuery{}})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudTrailLakeQuery(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var eventDataSourceArn, queryId string
	if h.Item != nil {
		data := h.Item.(*QueryInfo)
		eventDataSourceArn = *data.EventDataStoreArn
		queryId = *data.QueryId
	} else {
		eventDataSourceArn = d.KeyColumnQualString("event_data_store_arn")
		queryId = d.KeyColumnQualString("query_id")
	}

	// Create session
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Info("aws_cloudtrail_query.getCloudTrailLakeQuery", "client_error", err)
		return nil, err
	}

	params := &cloudtrail.DescribeQueryInput{
		EventDataStore: aws.String(eventDataSourceArn),
		QueryId:        aws.String(queryId),
	}

	// execute list call
	op, err := svc.DescribeQuery(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Info("aws_cloudtrail_query.getCloudTrailLakeQuery", "api_error", err)
		return nil, err
	}

	return &QueryInfo{aws.String(eventDataSourceArn), nil, op.QueryId, op.QueryStatus, op.QueryString, op.QueryStatistics}, nil
}
