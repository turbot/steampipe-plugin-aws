package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	snsv1 "github.com/aws/aws-sdk-go/service/sns"

	"github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSnsTopic(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sns_topic",
		Description: "AWS SNS Topic",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("topic_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFound", "InvalidParameter"}),
			},
			Hydrate: getTopicAttributes,
			Tags:    map[string]string{"service": "sns", "action": "GetTopicAttributes"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSnsTopics,
			Tags:    map[string]string{"service": "sns", "action": "ListTopics"},
		},

		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listTagsForSnsTopic,
				Tags: map[string]string{"service": "sns", "action": "ListTagsForResource"},
			},
			{
				Func: getTopicAttributes,
				Tags: map[string]string{"service": "sns", "action": "GetTopicAttributes"},
			},
		},

		GetMatrixItemFunc: SupportedRegionMatrix(snsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "topic_arn",
				Description: "Amazon Resource Name (ARN) of the Topic.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.TopicArn"),
			},
			{
				Name:        "display_name",
				Description: "The human-readable name used in the From field for notifications to email and email-json endpoints.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.DisplayName"),
			},
			{
				Name:        "application_failure_feedback_role_arn",
				Description: "IAM role for failed deliveries of notification messages sent to topics with platform application endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.ApplicationFailureFeedbackRoleArn"),
			},
			{
				Name:        "application_success_feedback_role_arn",
				Description: "IAM role for successful deliveries of notification messages sent to topics with platform application endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.ApplicationSuccessFeedbackRoleArn"),
			},
			{
				Name:        "application_success_feedback_sample_rate",
				Description: "Sample rate for successful deliveries of notification messages sent to topics with platform application endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.ApplicationSuccessFeedbackSampleRate"),
			},
			{
				Name:        "firehose_failure_feedback_role_arn",
				Description: "IAM role for failed deliveries of notification messages sent to topics with kinesis data firehose endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.FirehoseFailureFeedbackRoleArn"),
			},
			{
				Name:        "firehose_success_feedback_role_arn",
				Description: "IAM role for successful deliveries of notification messages sent to topics with kinesis data firehose endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.FirehoseSuccessFeedbackRoleArn"),
			},
			{
				Name:        "firehose_success_feedback_sample_rate",
				Description: "Sample rate for successful deliveries of notification messages sent to topics with firehose endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.FirehoseSuccessFeedbackSampleRate"),
			},
			{
				Name:        "http_failure_feedback_role_arn",
				Description: "IAM role for failed deliveries of notification messages sent to topics with http endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.HTTPFailureFeedbackRoleArn"),
			},
			{
				Name:        "http_success_feedback_role_arn",
				Description: "IAM role for successful deliveries of notification messages sent to topics with http endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.HTTPSuccessFeedbackRoleArn"),
			},
			{
				Name:        "http_success_feedback_sample_rate",
				Description: "Sample rate for successful deliveries of notification messages sent to topics with http endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.HTTPSuccessFeedbackSampleRate"),
			},
			{
				Name:        "lambda_failure_feedback_role_arn",
				Description: "IAM role for failed deliveries of notification messages sent to topics with lambda endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.LambdaFailureFeedbackRoleArn"),
			},
			{
				Name:        "lambda_success_feedback_role_arn",
				Description: "IAM role for successful deliveries of notification messages sent to topics with lambda endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.LambdaSuccessFeedbackRoleArn"),
			},
			{
				Name:        "lambda_success_feedback_sample_rate",
				Description: "Sample rate for successful deliveries of notification messages sent to topics with lambda endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.LambdaSuccessFeedbackSampleRate"),
			},
			{
				Name:        "owner",
				Description: "The AWS account ID of the topic's owner.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.Owner"),
			},
			{
				Name:        "sqs_failure_feedback_role_arn",
				Description: "IAM role for failed deliveries of notification messages sent to topics with sqs endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.SQSFailureFeedbackRoleArn"),
			},
			{
				Name:        "sqs_success_feedback_role_arn",
				Description: "IAM role for successful deliveries of notification messages sent to topics with sqs endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.SQSSuccessFeedbackRoleArn"),
			},
			{
				Name:        "sqs_success_feedback_sample_rate",
				Description: "Sample rate for successful deliveries of notification messages sent to topics with sqs endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.SQSSuccessFeedbackSampleRate"),
			},
			{
				Name:        "subscriptions_confirmed",
				Description: "The number of confirmed subscriptions for the topic.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.SubscriptionsConfirmed"),
			},
			{
				Name:        "subscriptions_deleted",
				Description: "The number of deleted subscriptions for the topic.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.SubscriptionsDeleted"),
			},
			{
				Name:        "subscriptions_pending",
				Description: "The number of subscriptions pending confirmation for the topic.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.SubscriptionsPending"),
			},
			{
				Name:        "kms_master_key_id",
				Description: "The ID of an AWS-managed customer master key (CMK) for Amazon SNS or a custom CMK.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.KmsMasterKeyId"),
			},
			{
				Name:        "tags_src",
				Description: "The list of tags associated with the topic.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForSnsTopic,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "policy",
				Description: "The topic's access control policy (i.e. Resource IAM Policy).",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.Policy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.Policy").Transform(unescape).Transform(policyToCanonical),
			},

			{
				Name:        "delivery_policy",
				Description: "The JSON object of the topic's delivery policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.DeliveryPolicy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "effective_delivery_policy",
				Description: "The effective delivery policy, taking system defaults into account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.EffectiveDeliveryPolicy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.TopicArn").Transform(arnToTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForSnsTopic,
				Transform:   transform.From(handleSNSTopicTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Attributes.TopicArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSnsTopics(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := SNSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sns_topic.listAwsSnsTopics", "get_client_error", err)
		return nil, err
	}

	params := &sns.ListTopicsInput{}
	// Does not support limit
	paginator := sns.NewListTopicsPaginator(svc, params, func(o *sns.ListTopicsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)

		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		if err != nil {
			plugin.Logger(ctx).Error("aws_sns_topic.listAwsSnsTopics", "api_error", err)
			return nil, err
		}
		for _, topic := range output.Topics {
			d.StreamListItem(ctx, &sns.GetTopicAttributesOutput{
				Attributes: map[string]string{
					"TopicArn": *topic.TopicArn,
				},
			})
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getTopicAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		data := h.Item.(*sns.GetTopicAttributesOutput)
		arn = types.SafeString(data.Attributes["TopicArn"])
	} else {
		arn = d.EqualsQuals["topic_arn"].GetStringValue()
	}

	if arn == "" {
		return nil, nil
	}

	// Get client
	svc, err := SNSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sns_topic.getTopicAttributes", "get_client_error", err)
		return nil, err
	}

	// Build params
	params := &sns.GetTopicAttributesInput{
		TopicArn: aws.String(arn),
	}

	op, err := svc.GetTopicAttributes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sns_topic.getTopicAttributes", "api_error", err)
		return nil, err
	}
	return op, nil
}

func listTagsForSnsTopic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	topicAttributesOutput := h.Item.(*sns.GetTopicAttributesOutput)

	// Get client
	svc, err := SNSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sns_topic.listTagsForSnsTopic", "get_client_error", err)
		return nil, err
	}

	// Build param
	param := &sns.ListTagsForResourceInput{
		ResourceArn: aws.String(topicAttributesOutput.Attributes["TopicArn"]),
	}

	// Next token is not supported
	// AWS supports upto 50 tags
	topicTags, err := svc.ListTagsForResource(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sns_topic.listTagsForSnsTopic", "api_error", err)
		return nil, err
	}
	return topicTags, nil
}

//// TRANSFORM FUNCTIONS

func handleSNSTopicTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*sns.ListTagsForResourceOutput)
	if len(tags.Tags) == 0 {
		return nil, nil
	}

	turbotTagsMap := map[string]string{}
	for _, i := range tags.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return turbotTagsMap, nil
}
