package aws

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//// TABLE DEFINITION

func tableAwsSqsQueue(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sqs_queue",
		Description: "AWS SQS Queue",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("queue_url"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"AWS.SimpleQueueService.NonExistentQueue"}),
			},
			Hydrate: getQueueAttributes,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSqsQueues,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "queue_url",
				Description: "The URL of the Amazon SQS queue.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.QueueUrl"),
			},
			{
				Name:        "queue_arn",
				Description: "The Amazon resource name (ARN) of the queue.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.QueueArn"),
			},
			{
				Name:        "fifo_queue",
				Description: "Returns true if the queue is FIFO.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.FifoQueue"),
				Default:     false,
			},
			{
				Name:        "delay_seconds",
				Description: "The default delay on the queue in seconds.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.DelaySeconds"),
			},
			{
				Name:        "max_message_size",
				Description: "The limit of how many bytes a message can contain before Amazon SQS rejects it.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.MaximumMessageSize"),
			},
			{
				Name:        "message_retention_seconds",
				Description: "The length of time, in seconds, for which Amazon SQS retains a message.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.MessageRetentionPeriod"),
			},
			{
				Name:        "receive_wait_time_seconds",
				Description: "The length of time, in seconds, for which the ReceiveMessage action waits for a message to arrive.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.ReceiveMessageWaitTimeSeconds"),
			},
			{
				Name:        "sqs_managed_sse_enabled",
				Description: "Returns true if the queue is using SSE-SQS encryption with SQS-owned encryption keys.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.SqsManagedSseEnabled"),
			},
			{
				Name:        "visibility_timeout_seconds",
				Description: "The visibility timeout for the queue in seconds.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.VisibilityTimeout"),
			},
			{
				Name:        "policy",
				Description: "The resource IAM policy of the queue.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.Policy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.Policy").Transform(unescape).Transform(policyToCanonical),
			},

			{
				Name:        "redrive_policy",
				Description: "The string that includes the parameters for the dead-letter queue functionality of the source queue as a JSON object.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.RedrivePolicy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "content_based_deduplication",
				Description: "Mentions whether content-based deduplication is enabled for the queue.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.ContentBasedDeduplication"),
			},
			{
				Name:        "kms_master_key_id",
				Description: "The ID of an AWS-managed customer master key (CMK) for Amazon SQS or a custom CMK.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.KmsMasterKeyId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listQueueTags,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.QueueUrl").Transform(getAwsSqsQueueTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.QueueArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSqsQueues(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsSqsQueues")

	// Create session
	svc, err := SQSService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &sqs.ListQueuesInput{
		MaxResults: aws.Int64(1000),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	err = svc.ListQueuesPages(
		input,
		func(page *sqs.ListQueuesOutput, lastPage bool) bool {
			for _, queueURL := range page.QueueUrls {
				d.StreamListItem(ctx, &sqs.GetQueueAttributesOutput{
					Attributes: map[string]*string{
						"QueueUrl": queueURL,
					},
				})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getQueueAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getQueueAttributes")

	var queueURL string
	if h.Item != nil {
		data := h.Item.(*sqs.GetQueueAttributesOutput)
		queueURL = types.SafeString(data.Attributes["QueueUrl"])
	} else {
		queueURL = d.KeyColumnQuals["queue_url"].GetStringValue()
	}

	// Create session
	svc, err := SQSService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueURL),
		AttributeNames: []*string{aws.String("All")},
	}

	op, err := svc.GetQueueAttributes(input)
	if err != nil {
		return nil, err
	}

	// Add QueueUrl info to the output as it is missing from GetQueueAttributesOutput
	op.Attributes["QueueUrl"] = aws.String(queueURL)

	return op, nil
}

func listQueueTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listQueueTags")
	queueAttributesOutput := h.Item.(*sqs.GetQueueAttributesOutput)

	// Create session
	svc, err := SQSService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	param := &sqs.ListQueueTagsInput{
		QueueUrl: queueAttributesOutput.Attributes["QueueUrl"],
	}

	queueTags, err := svc.ListQueueTags(param)
	if err != nil {
		return nil, err
	}
	return queueTags, nil
}

//// TRANSFORM FUNCTION

func getAwsSqsQueueTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	queueURL := types.SafeString(d.Value)

	queueName, err := extractNameFromSqsQueueURL(queueURL)
	if err != nil {
		return nil, err
	}

	return queueName, nil
}
