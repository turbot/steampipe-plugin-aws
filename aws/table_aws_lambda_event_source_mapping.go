package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	lambdav1 "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsLambdaEventSourceMapping(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_event_source_mapping",
		Description: "AWS Lambda Event Source Mapping",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("uuid"),
			Hydrate:    getAwsLambdaEventSourceMapping,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsLambdaEventSourceMappings,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "arn",
					Require: plugin.Optional,
				},
				{
					Name:    "function_arn",
					Require: plugin.Optional,
				},
				{
					Name:    "function_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(lambdav1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "uuid",
				Description: "The identifier of the event source mapping.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UUID"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the event source.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EventSourceArn"),
			},
			{
				Name:        "function_arn",
				Description: "The ARN of the Lambda function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "function_name",
				Description: "The name of the Lambda function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getFunctionNameFromArn),
			},
			{
				Name:        "batch_size",
				Description: "The maximum number of records in each batch that Lambda pulls from your stream or queue and sends to your function.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "bisect_batch_on_function_error",
				Description: "If the function returns an error, split the batch in two and retry.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_modified",
				Description: "The date that the event source mapping was last updated or that its state changed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_processing_result",
				Description: "The result of the last Lambda invocation of your function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "maximum_batching_window_in_seconds",
				Description: "The maximum amount of time, in seconds, that Lambda spends gathering records before invoking the function.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "maximum_record_age_in_seconds",
				Description: "Discard records older than the specified age.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "maximum_retry_attempts",
				Description: "Discard records after the specified number of retries.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "parallelization_factor",
				Description: "The number of batches to process concurrently from each shard.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "state",
				Description: "The state of the event source mapping.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_transition_reason",
				Description: "Indicates whether a user or Lambda made the last change to the event source mapping.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "starting_position",
				Description: "The position in a stream from which to start reading. Required for Amazon Kinesis and Amazon DynamoDB Stream event sources. AT_TIMESTAMP is supported only for Amazon Kinesis streams, Amazon DocumentDB, Amazon MSK, and self-managed Apache Kafka.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "starting_position_timestamp",
				Description: "The position in a stream from which to start reading. With StartingPosition set to AT_TIMESTAMP, the time from which to start reading, in Unix time seconds. StartingPositionTimestamp cannot be in the future.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tumbling_window_in_seconds",
				Description: "The duration in seconds of a processing window for DynamoDB and Kinesis Streams event sources. A value of 0 seconds indicates no tumbling window.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "function_response_types",
				Description: "A list of current response type enums applied to the event source mapping.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_access_configurations",
				Description: "An array of the authentication protocol, VPC components, or virtual host to secure and define your event source.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "destination_config",
				Description: "An Amazon SQS queue or Amazon SNS topic destination for discarded records.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "filter_criteria",
				Description: "An object that defines the filter criteria that determine whether Lambda should process an event.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "amazon_managed_kafka_event_source_config",
				Description: "Specific configuration settings for an Amazon Managed Streaming for Apache Kafka.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "queues",
				Description: "The name of the Amazon MQ broker destination queue to consume.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "scaling_config",
				Description: "The scaling configuration for the event source.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "self_managed_event_source",
				Description: "The self-managed Apache Kafka cluster for your event source.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "self_managed_kafka_event_source_config",
				Description: "Specific configuration settings for a self-managed Apache Kafka event source.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "topics",
				Description: "The name of the Kafka topic.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getEventSourceMappingTitle),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsLambdaEventSourceMappings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_event_source_mapping.listAwsLambdaEventSourceMappings", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(10000)
	input := lambda.ListEventSourceMappingsInput{}
	if d.EqualsQualString("arn") != "" {
		input.EventSourceArn = aws.String(d.EqualsQualString("arn"))
	}
	if d.EqualsQualString("function_arn") != "" || d.EqualsQualString("function_name") != "" {
		if d.EqualsQualString("function_arn") != "" {
			input.FunctionName = aws.String(d.EqualsQualString("function_arn"))
		} else if d.EqualsQualString("function_name") != "" {
			input.FunctionName = aws.String(d.EqualsQualString("function_name"))
		}
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	input.MaxItems = aws.Int32(maxItems)
	paginator := lambda.NewListEventSourceMappingsPaginator(svc, &input, func(o *lambda.ListEventSourceMappingsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lambda_event_source_mapping.listAwsLambdaEventSourceMappings", "api_error", err)
			return nil, err
		}

		for _, mapping := range output.EventSourceMappings {
			d.StreamListItem(ctx, mapping)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsLambdaEventSourceMapping(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var uuid string

	if h.Item != nil {
		uuid = *h.Item.(types.EventSourceMappingConfiguration).UUID
	} else {
		uuid = d.EqualsQuals["uuid"].GetStringValue()
	}

	// Empty input check
	if strings.TrimSpace(uuid) == "" {
		return nil, nil
	}

	// Create Session
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_event_source_mapping.getAwsLambdaEventSourceMapping", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &lambda.GetEventSourceMappingInput{
		UUID: &uuid,
	}

	rowData, err := svc.GetEventSourceMapping(ctx, params) // ---> *GetEventSourceMappingOutput
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_event_source_mapping.getAwsLambdaEventSourceMapping", "api_error", err)
		return nil, err
	}

	return rowData, nil

}

//// TRANSFORM FUNCTIONS

func getFunctionNameFromArn(ctx context.Context, td *transform.TransformData) (interface{}, error) {
	arn := *td.HydrateItem.(types.EventSourceMappingConfiguration).FunctionArn
	parts := strings.Split(arn, ":")
	return parts[len(parts)-1], nil
}

func getEventSourceMappingTitle(ctx context.Context, td *transform.TransformData) (interface{}, error) {
	arn := *td.HydrateItem.(types.EventSourceMappingConfiguration).EventSourceArn
	parts := strings.Split(arn, ":")
	return parts[len(parts)-1], nil
}
