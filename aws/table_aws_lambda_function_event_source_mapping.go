package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	lambdav1 "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"strings"
)

const (
	EV_SOURCE_MAPPING_TABLE_NAME = "aws_lambda_function_event_source_mapping"
)

func tableAwsLambdaFunctionEventSourceMapping(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        EV_SOURCE_MAPPING_TABLE_NAME,
		Description: "AWS Lambda Function Event Source Mappings",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("uuid"),
			Hydrate:    getAwsLambdaFunctionTrigger,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsLambdaFunctionTriggers,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(lambdav1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "uuid",
				Description: "The identifier of the function trigger",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UUID"),
			},
			{
				Name:        "event_source_arn",
				Description: "The Amazon Resource Name (ARN) of the event source",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EventSourceArn"),
			},
			{
				Name:        "function_arn",
				Description: "The ARN of the Lambda function",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FunctionArn"),
			},
			{
				Name:        "function_name",
				Description: "The name of the Lambda function",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getFunctionNameFromArn),
			},
			//{
			//	Name:        "service_name",
			//	Description: "The service name of the event source",
			//	Type:        proto.ColumnType_STRING,
			//	Transform:   transform.From(getServiceNameFromArn),
			//},
			//{
			//	Name:        "service_component_name",
			//	Description: "The service component name of the event source",
			//	Type:        proto.ColumnType_STRING,
			//	Transform:   transform.From(getServiceComponentNameFromArn),
			//},
			//{
			//	Name:        "service_component_type",
			//	Description: "The service component type of the event source",
			//	Type:        proto.ColumnType_STRING,
			//	Transform:   transform.From(getServiceComponentTypeFromArn),
			//},
			{
				Name:        "batch_size",
				Description: "The maximum number of records in each batch that Lambda pulls from your stream or queue and sends to your function",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("BatchSize"),
			},
			{
				Name:        "bisect_batch_on_function_error",
				Description: "If the function returns an error, split the batch in two and retry",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("BisectBatchOnFunctionError"),
			},
			{
				Name:        "last_modified",
				Description: "The date that the event source mapping was last updated or that its state changed",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastModified"),
			},
			{
				Name:        "last_processing_result",
				Description: "The result of the last Lambda invocation of your function",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LastProcessingResult"),
			},
			{
				Name:        "maximum_batching_window_in_seconds",
				Description: "The maximum amount of time, in seconds, that Lambda spends gathering records before invoking the function",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MaximumBatchingWindowInSeconds"),
			},
			{
				Name:        "maximum_record_age_in_seconds",
				Description: "Discard records older than the specified age",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MaximumRecordAgeInSeconds"),
			},
			{
				Name:        "maximum_retry_attempts",
				Description: "Discard records after the specified number of retries",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MaximumRetryAttempts"),
			},
			{
				Name:        "parallelization_factor",
				Description: "The number of batches to process concurrently from each shard",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ParallelizationFactor"),
			},
			{
				Name:        "state",
				Description: "The state of the event source mapping",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State"),
			},
			{
				Name:        "state_transition_reason",
				Description: "Indicates whether a user or Lambda made the last change to the event source mapping",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StateTransitionReason"),
			},
			{
				Name:        "tumbling_window_in_seconds",
				Description: "The duration in seconds of a processing window",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("TumblingWindowInSeconds"),
			},
			{
				Name:        "function_response_types",
				Description: "A list of current response type enums applied to the event source mapping",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FunctionResponseTypes"),
			},
			{
				Name:        "source_access_configurations",
				Description: "An array of the authentication protocol, VPC components, or virtual host to secure and define your event source",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SourceAccessConfigurations"),
			},
			{
				Name:        "destination_config",
				Description: "An Amazon SQS queue or Amazon SNS topic destination for discarded records",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DestinationConfig"),
			},
			{
				Name:        "filter_criteria",
				Description: "An object that defines the filter criteria that determine whether Lambda should process an event",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FilterCriteria"),
			},
			// ===== Very Event Source Specific Mapping Configurations
			{
				Name:        "amazon_managed_kafka_event_source_config",
				Description: "Specific configuration settings for an Amazon Managed Streaming for Apache Kafka",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AmazonManagedKafkaEventSourceConfig"),
			},
			{
				Name:        "queues",
				Description: "The name of the Amazon MQ broker destination queue to consume",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Queues"),
			},
			{
				Name:        "scaling_config",
				Description: "The scaling configuration for the event source",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ScalingConfig"),
			},
			{
				Name:        "self_managed_event_source",
				Description: "The self-managed Apache Kafka cluster for your event source",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SelfManagedEventSource"),
			},
			{
				Name:        "self_managed_kafka_event_source_config",
				Description: "Specific configuration settings for a self-managed Apache Kafka event source",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SelfManagedKafkaEventSourceConfig"),
			},
			{
				Name:        "starting_position",
				Description: "The position in a stream from which to start reading",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StartingPosition"),
			},
			{
				Name:        "starting_position_timestamp",
				Description: "The position in a stream from which to start reading",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("StartingPositionTimestamp"),
			},
			{
				Name:        "topics",
				Description: "The name of the Kafka topic[s]",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Topics"),
			},
		}),
	}
}

//// HYDRATE FUNCTION

func getAwsLambdaFunctionTrigger(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

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
		plugin.Logger(ctx).Error(EV_SOURCE_MAPPING_TABLE_NAME+".getAwsLambdaFunctionTrigger", "connection_error", err)
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
		plugin.Logger(ctx).Error(EV_SOURCE_MAPPING_TABLE_NAME+".getAwsLambdaFunctionTrigger", "api_error", err)
		return nil, err
	}

	return rowData, nil

}

//// LIST FUNCTION

func listAwsLambdaFunctionTriggers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error(EV_SOURCE_MAPPING_TABLE_NAME+".listAwsLambdaFunctionTriggers", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(10000)
	input := lambda.ListEventSourceMappingsInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
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
			plugin.Logger(ctx).Error(EV_SOURCE_MAPPING_TABLE_NAME+".listAwsLambdaFunctionTriggers", "api_error", err)
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

func getFunctionNameFromArn(ctx context.Context, td *transform.TransformData) (interface{}, error) {
	arn := *td.HydrateItem.(types.EventSourceMappingConfiguration).FunctionArn
	parts := strings.Split(arn, ":")
	return parts[len(parts)-1], nil
}

// Commented these, and the using columns because the implementation is brittle or non-useful in some cases.

//func getServiceNameFromArn(ctx context.Context, td *transform.TransformData) (interface{}, error) {
//	arn := td.HydrateItem.(types.EventSourceMappingConfiguration).EventSourceArn
//	if arn == nil {
//		return "nil", nil
//	}
//	parts := strings.Split(*arn, ":")
//	if len(parts) >= 3 {
//		return parts[2], nil
//	}
//	return *arn, nil
//}

//func getServiceComponentNameFromArn(ctx context.Context, td *transform.TransformData) (interface{}, error) {
//	arn := td.HydrateItem.(types.EventSourceMappingConfiguration).EventSourceArn
//	if arn == nil {
//		return "nil", nil
//	}
//	parts := strings.SplitN(*arn, "/", 2)
//	if len(parts) == 1 {
//		return parts[0], nil
//	}
//	return parts[1], nil
//}
//
//func getServiceComponentTypeFromArn(ctx context.Context, td *transform.TransformData) (interface{}, error) {
//	arn := td.HydrateItem.(types.EventSourceMappingConfiguration).EventSourceArn
//	if arn == nil {
//		return "nil", nil
//	}
//	if strings.HasPrefix(*arn, "arn:aws:sqs:") {
//		return "queue", nil
//	}
//	sub := strings.Split(*arn, "/")[0]
//	parts := strings.Split(sub, ":")
//	return parts[len(parts)-1], nil
//}
