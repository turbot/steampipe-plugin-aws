package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"

	athenav1 "github.com/aws/aws-sdk-go/service/athena"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsAthenaQueryExecution(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_athena_query_execution",
		Description: "AWS Athena Query Execution",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAwsAthenaQueryExecution,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsAthenaWorkGroups,
			Hydrate:       listAwsAthenaQueryExecutions,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(athenav1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for each query execution.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("QueryExecution.QueryExecutionId"),
			},
			{
				Name:        "workgroup",
				Description: "The name of the workgroup in which the query ran.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.WorkGroup"),
			},
			{
				Name:        "catalog",
				Description: "The name of the data catalog used in the query execution.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.QueryExecutionContext.Catalog"),
			},
			{
				Name:        "database",
				Description: "The name of the data database used in the query execution.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.QueryExecutionContext.Database"),
			},
			{
				Name:        "query",
				Description: "The SQL query statements which the query execution ran.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Query"),
			},
			{
				Name:        "effective_engine_version",
				Description: "The engine version on which the query runs.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.EngineVersion.EffectiveEngineVersion"),
			},
			{
				Name:        "selected_engine_version",
				Description: "The engine version requested by the users.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.EngineVersion.SelectedEngineVersion"),
			},
			{
				Name:        "execution_parameters",
				Description: "A list of values for the parameters in a query.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.ExecutionParameters"),
			},
			{
				Name:        "statement_type",
				Description: "The type of query statement that was run.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.StatementType"),
			},
			{
				Name:        "substatement_type",
				Description: "The kind of query statement that was run.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.SubstatementType"),
			},
			{
				Name:        "state",
				Description: "The state of query execution.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Status.State"),
			},
			{
				Name:        "state_change_reason",
				Description: "Further detail about the status of the query.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Status.StateChangeReason"),
			},
			{
				Name:        "submission_date_time",
				Description: "The date and time that the query was submitted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Status.SubmissionDateTime"),
			},
			{
				Name:        "completion_date_time",
				Description: "The date and time that the query completed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Status.CompletionDateTime"),
			},
			{
				Name:        "error_message",
				Description: "Contains a short description of the error that occurred.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Status.AthenaError.ErrorMessage"),
			},
			{
				Name:        "error_type",
				Description: "An integer value that provides specific information about an Athena query error.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Status.AthenaError.ErrorType"),
			},
			{
				Name:        "error_category",
				Description: "An integer value that specifies the category of a query failure error.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Status.AthenaError.ErrorCategory"),
			},
			{
				Name:        "retryable",
				Description: "True if the query might succeed if resubmitted.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Status.AthenaError.Retryable"),
			},
			{
				Name:        "data_manifest_location",
				Description: "The location and file name of a data manifest file.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Statistics.DataManifestLocation"),
			},
			{
				Name:        "data_scanned_in_bytes",
				Description: "The number of bytes in the data that was queried.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Statistics.DataScannedInBytes"),
			},
			{
				Name:        "engine_execution_time_in_millis",
				Description: "The number of milliseconds that the query took to execute.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Statistics.EngineExecutionTimeInMillis"),
			},
			{
				Name:        "query_planning_time_in_millis",
				Description: "The number of milliseconds that Athena took to plan the query processing flow.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Statistics.QueryPlanningTimeInMillis"),
			},
			{
				Name:        "query_queue_time_in_millis",
				Description: "The number of milliseconds that the query was in your query queue waiting for resources.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Statistics.QueryQueueTimeInMillis"),
			},
			{
				Name:        "service_processing_time_in_millis",
				Description: "The number of milliseconds that Athena took to finalize and publish the query results after the query engine finished running the query.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Statistics.ServiceProcessingTimeInMillis"),
			},
			{
				Name:        "total_execution_time_in_millis",
				Description: "The number of milliseconds that Athena took to run the query.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Statistics.TotalExecutionTimeInMillis"),
			},
			{
				Name:        "reused_previous_result",
				Description: "True if a previous query result was reused; false if the result was generated.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Statistics.ResultReuseInformation.ReusedPreviousResult"),
			},
			{
				Name:        "s3_acl_option",
				Description: "The Amazon S3 canned ACL that Athena should specify when storing query results.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.ResultConfiguration.AclConfiguration.S3AclOption"),
			},
			{
				Name:        "encryption_option",
				Description: "Indicates whether Amazon S3 server-side encryption with Amazon S3-managed keys (SSE_S3), server-side encryption with KMS-managed keys (SSE_KMS), or client-side encryption with KMS-managed keys (CSE_KMS) is used.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.ResultConfiguration.EncryptionConfiguration.EncryptionOption"),
			},
			{
				Name:        "kms_key",
				Description: "For SSE_KMS and CSE_KMS, this is the KMS key ARN or ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.ResultConfiguration.EncryptionConfiguration.KmsKey"),
			},
			{
				Name:        "expected_bucket_owner",
				Description: "The Amazon Web Services account ID that you expect to be the owner of the Amazon S3 bucket specified by ResultConfiguration$OutputLocation.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.ResultConfiguration.ExpectedBucketOwner"),
			},
			{
				Name:        "output_location",
				Description: "The location in Amazon S3 where your query results are stored.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.ResultConfiguration.OutputLocation"),
			},
			{
				Name:        "result_reuse_by_age_enabled",
				Description: "True if previous query results can be reused when the query is run.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.ResultReuseConfiguration.ResultReuseByAgeConfiguration.Enabled"),
			},
			{
				Name:        "result_reuse_by_age_mag_age_in_minutes",
				Description: "Specifies, in minutes, the maximum age of a previous query result that Athena should consider for reuse. The default is 60.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.ResultReuseConfiguration.ResultReuseByAgeConfiguration.MaxAgeInMinutes"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsAthenaQueryExecutions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := AthenaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_athena_query_execution.listAwsAthenaQueryExecutions", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxResults := int32(50)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxResults {
			maxResults = int32(limit)
		}
	}

	workgroup := h.Item.(types.WorkGroup)

	input := athena.ListQueryExecutionsInput{
		MaxResults: aws.Int32(maxResults),
		WorkGroup:  workgroup.Name,
	}

	paginator := athena.NewListQueryExecutionsPaginator(svc, &input, func(o *athena.ListQueryExecutionsPaginatorOptions) {
		o.Limit = maxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_athena_query_execution.listAwsAthenaQueryExecutions", "api_error", err)
			return nil, err
		}

		for _, queryExecutionId := range output.QueryExecutionIds {
			id := strings.Clone(queryExecutionId)
			var queryExecution = types.QueryExecution{QueryExecutionId: &id}

			d.StreamListItem(ctx, athena.GetQueryExecutionOutput{QueryExecution: &queryExecution})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsAthenaQueryExecution(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var id string
	if h.Item != nil {
		id = *h.Item.(athena.GetQueryExecutionOutput).QueryExecution.QueryExecutionId
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
	}

	// Empty input check
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	svc, err := AthenaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_athena_query_execution.getAwsAthenaQueryExecution", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build params
	params := &athena.GetQueryExecutionInput{
		QueryExecutionId: aws.String(id),
	}

	rowData, err := svc.GetQueryExecution(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_athena_query_execution.getAwsAthenaQueryExecution", "api_error", err)
		return nil, err
	}

	return rowData, nil
}
