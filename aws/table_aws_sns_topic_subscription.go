package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSnsTopicSubscription(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sns_topic_subscription",
		Description: "AWS SNS Topic Subscription",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("subscription_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NotFound", "InvalidParameter"}),
			},
			Hydrate: getSubscriptionAttributes,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSnsTopicSubscriptions,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "subscription_arn",
				Description: "Amazon Resource Name of the subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.SubscriptionArn").NullIfEqual("PendingConfirmation"),
			},
			{
				Name:        "topic_arn",
				Description: "The topic ARN that the subscription is associated with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.TopicArn"),
			},
			{
				Name:        "owner",
				Description: "The AWS account ID of the subscription's owner.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.Owner"),
			},
			{
				Name:        "protocol",
				Description: "The subscription's protocol.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.Protocol"),
			},
			{
				Name:        "endpoint",
				Description: "The subscription's endpoint (format depends on the protocol).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.Endpoint"),
			},
			{
				Name:        "confirmation_was_authenticated",
				Description: "Reflects authentication status of the subscription.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getSubscriptionAttributes,
				Transform:   transform.FromField("Attributes.ConfirmationWasAuthenticated"),
			},
			{
				Name:        "pending_confirmation",
				Description: "Reflects the confirmation status of the subscription. True if the subscription hasn't been confirmed.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getSubscriptionAttributes,
				Transform:   transform.FromField("Attributes.PendingConfirmation"),
			},
			{
				Name:        "raw_message_delivery",
				Description: "true if raw message delivery is enabled for the subscription.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getSubscriptionAttributes,
				Transform:   transform.FromField("Attributes.RawMessageDelivery"),
			},
			{
				Name:        "delivery_policy",
				Description: "The JSON of the subscription's delivery policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSubscriptionAttributes,
				Transform:   transform.FromField("Attributes.DeliveryPolicy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "effective_delivery_policy",
				Description: "The JSON of the effective delivery policy that takes into account the topic delivery policy and account system defaults.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSubscriptionAttributes,
				Transform:   transform.FromField("Attributes.EffectiveDeliveryPolicy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "redrive_policy",
				Description: "When specified, sends undeliverable messages to the specified Amazon SQS dead-letter queue. Messages that can't be delivered due to client errors (for example, when the subscribed endpoint is unreachable) or server errors (for example, when the service that powers the subscribed endpoint becomes unavailable) are held in the dead-letter queue for further analysis or reprocessing.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSubscriptionAttributes,
				Transform:   transform.FromField("Attributes.RedrivePolicy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "filter_policy",
				Description: "The filter policy JSON that is assigned to the subscription.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSubscriptionAttributes,
				Transform:   transform.FromField("Attributes.FilterPolicy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.SubscriptionArn").Transform(arnToTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Attributes.SubscriptionArn").Transform(subscriptionArnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSnsTopicSubscriptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsSnsTopicSubscriptions")

	// Create Session
	svc, err := SNSService(ctx, d)
	if err != nil {
		return nil, err
	}

	err = svc.ListSubscriptionsPages(
		&sns.ListSubscriptionsInput{},
		func(page *sns.ListSubscriptionsOutput, lastPage bool) bool {
			for _, subscription := range page.Subscriptions {
				d.StreamListItem(ctx, &sns.GetSubscriptionAttributesOutput{
					Attributes: map[string]*string{
						"Endpoint":        subscription.Endpoint,
						"Owner":           subscription.Owner,
						"Protocol":        subscription.Protocol,
						"SubscriptionArn": subscription.SubscriptionArn,
						"TopicArn":        subscription.TopicArn,
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

func getSubscriptionAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSubscriptionAttributes")

	var arn string
	if h.Item != nil {
		data := h.Item.(*sns.GetSubscriptionAttributesOutput)
		arn = types.SafeString(data.Attributes["SubscriptionArn"])
	} else {
		arn = d.KeyColumnQuals["subscription_arn"].GetStringValue()
	}

	// Create session
	svc, err := SNSService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &sns.GetSubscriptionAttributesInput{
		SubscriptionArn: aws.String(arn),
	}

	// As of 7th september 2020, Next token is not supported in go
	op, err := svc.GetSubscriptionAttributes(input)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			plugin.Logger(ctx).Trace("SubErrorCode", a.Code())
			if a.Code() == "NotFound" || a.Code() == "InvalidParameter" {
				return nil, nil
			}
			return nil, err
		}
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func subscriptionArnToAkas(_ context.Context, d *transform.TransformData) (interface{}, error) {
	arn := types.SafeString(d.Value)

	if arn == "PendingConfirmation" {
		return []string{}, nil
	}

	return []string{arn}, nil
}
