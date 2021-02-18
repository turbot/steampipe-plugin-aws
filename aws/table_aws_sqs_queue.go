package aws

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//// TABLE DEFINITION

func tableAwsSqsQueue(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sqs_queue",
		Description: "AWS SQS Queue",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("queue_url"),
			ItemFromKey:       sqsQueueFromKey,
			ShouldIgnoreError: isNotFoundError([]string{"AWS.SimpleQueueService.NonExistentQueue"}),
			Hydrate:           getQueueAttributes,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSqsQueues,
		},
		GetMatrixItem: BuildRegionList,
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
				Description: "The limit of how many bytes a message can contain before Amazon SQS rejects it",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.MaximumMessageSize"),
			},
			{
				Name:        "message_retention_seconds",
				Description: "The length of time, in seconds, for which Amazon SQS retains a message",
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
				Name:        "visibility_timeout_seconds",
				Description: "The visibility timeout for the queue in seconds.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.VisibilityTimeout"),
			},
			{
				Name:        "policy",
				Description: "The resource IAM policy of the queue",
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
				Description: "Mentions whether content-based deduplication is enabled for the queue",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getQueueAttributes,
				Transform:   transform.FromField("Attributes.ContentBasedDeduplication"),
			},
			{
				Name:        "kms_master_key_id",
				Description: "the ID of an AWS-managed customer master key (CMK) for Amazon SQS or a custom CMK.",
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

//// BUILD HYDRATE INPUT

func sqsQueueFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	queueURL := quals["queue_url"].GetStringValue()
	item := &sqs.GetQueueAttributesOutput{
		Attributes: map[string]*string{
			"QueueUrl": &queueURL,
		},
	}
	return item, nil
}

//// LIST FUNCTION

func listAwsSqsQueues(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsSqsQueues", "AWS_REGION", region)

	// Create session
	svc, err := SQSService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	err = svc.ListQueuesPages(
		&sqs.ListQueuesInput{},
		func(page *sqs.ListQueuesOutput, lastPage bool) bool {
			for _, queueURL := range page.QueueUrls {
				d.StreamListItem(ctx, &sqs.GetQueueAttributesOutput{
					Attributes: map[string]*string{
						"QueueUrl": queueURL,
					},
				})
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getQueueAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getQueueAttributes")
	queueAttributesOutput := h.Item.(*sqs.GetQueueAttributesOutput)
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create session
	svc, err := SQSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &sqs.GetQueueAttributesInput{
		QueueUrl:       queueAttributesOutput.Attributes["QueueUrl"],
		AttributeNames: []*string{aws.String("All")},
	}

	op, err := svc.GetQueueAttributes(input)
	if err != nil {
		return nil, err
	}

	// Add QueueUrl info to the output as it is missing from GetQueueAttributesOutput
	op.Attributes["QueueUrl"] = queueAttributesOutput.Attributes["QueueUrl"]

	return op, nil
}

func listQueueTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listQueueTags")
	queueAttributesOutput := h.Item.(*sqs.GetQueueAttributesOutput)
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create session
	svc, err := SQSService(ctx, d, region)
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
