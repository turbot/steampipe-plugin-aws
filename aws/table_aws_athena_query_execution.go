package aws

import (
	"context"
	// "errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"

	// "github.com/aws/smithy-go"
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
			Hydrate: listAwsAthenaQueryExeuctions,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(athenav1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The id of the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("QueryExecutionId"),
			},
			{
				Name:        "query",
				Description: "The SQL query",
				Type:        proto.ColumnType_STRING,
                                Hydrate:     getAwsAthenaQueryExecution,
				Transform:   transform.FromField("QueryExecution.Query"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsAthenaQueryExeuctions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	input := athena.ListQueryExecutionsInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxResults {
			if limit < 1 {
				maxResults = int32(1)
			} else {
				maxResults = int32(limit)
			}
		}
	}

	input.MaxResults = aws.Int32(maxResults)
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

			d.StreamListItem(ctx, queryExecution)

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
		id = *h.Item.(types.QueryExecution).QueryExecutionId
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
