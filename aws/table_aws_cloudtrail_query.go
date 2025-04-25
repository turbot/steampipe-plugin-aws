package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudTrailQuery(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudtrail_query",
		Description: "AWS CloudTrail Query",
		Get: &plugin.GetConfig{
			Hydrate:    getCloudTrailQuery,
			Tags:       map[string]string{"service": "cloudtrail", "action": "DescribeQuery"},
			KeyColumns: plugin.AllColumns([]string{"event_data_store_arn", "query_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EventDataStoreNotFoundException", "QueryIdNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCloudTrailEventDataStores,
			Hydrate:       listCloudTrailLakeQueries,
			Tags:          map[string]string{"service": "cloudtrail", "action": "ListQueries"},
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
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCloudTrailQuery,
				Tags: map[string]string{"service": "cloudtrail", "action": "DescribeQuery"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CLOUDTRAIL_SERVICE_ID),
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
				Name:        "delivery_s3_uri",
				Description: "The URI for the S3 bucket where CloudTrail delivered query results, if applicable.",
				Hydrate:     getCloudTrailQuery,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delivery_status",
				Description: "The delivery status.",
				Hydrate:     getCloudTrailQuery,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "error_message",
				Description: "The error message returned if a query failed.",
				Hydrate:     getCloudTrailQuery,
				Type:        proto.ColumnType_STRING,
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
				Hydrate:     getCloudTrailQuery,
				Transform:   transform.FromField("QueryStatistics.BytesScanned"),
			},
			{
				Name:        "events_matched",
				Description: "The number of events that matched a query.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudTrailQuery,
				Transform:   transform.FromField("QueryStatistics.EventsMatched"),
			},
			{
				Name:        "events_scanned",
				Description: "The number of events that the query scanned in the event data store.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudTrailQuery,
				Transform:   transform.FromField("QueryStatistics.EventsScanned"),
			},
			{
				Name:        "execution_time_in_millis",
				Description: "The query's run time, in milliseconds.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudTrailQuery,
				Transform:   transform.FromField("QueryStatistics.ExecutionTimeInMillis"),
			},
			{
				Name:        "query_string",
				Description: "The SQL code of a query.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudTrailQuery,
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
type ListQueryInfo struct {
	EventDataStoreArn *string
	types.Query
}

type GetQueryInfo struct {
	EventDataStoreArn *string
	*cloudtrail.DescribeQueryOutput
}

//// LIST FUNCTION

func listCloudTrailLakeQueries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	eventDataStore := h.Item.(types.EventDataStore)

	if d.EqualsQualString("event_data_store_arn") != "" {
		if d.EqualsQualString("event_data_store_arn") != *eventDataStore.EventDataStoreArn {
			return nil, nil
		}
	}

	// Get client
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_query.listCloudTrailLakeQueries", "client_error", err)
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

	if d.EqualsQualString("query_status") != "" {
		input.QueryStatus = types.QueryStatus(d.EqualsQualString("query_status"))
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
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

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
			d.StreamListItem(ctx, &ListQueryInfo{eventDataStore.EventDataStoreArn, item})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudTrailQuery(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var eventDataSourceArn, queryId string
	if h.Item != nil {
		data := h.Item.(*ListQueryInfo)
		eventDataSourceArn = *data.EventDataStoreArn
		queryId = *data.QueryId
	} else {
		eventDataSourceArn = d.EqualsQualString("event_data_store_arn")
		queryId = d.EqualsQualString("query_id")
	}

	// Create session
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_query.getCloudTrailQuery", "client_error", err)
		return nil, err
	}

	params := &cloudtrail.DescribeQueryInput{
		EventDataStore: aws.String(eventDataSourceArn),
		QueryId:        aws.String(queryId),
	}

	// execute list call
	op, err := svc.DescribeQuery(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_query.getCloudTrailQuery", "api_error", err)
		return nil, err
	}

	return &GetQueryInfo{aws.String(eventDataSourceArn), op}, nil
}
