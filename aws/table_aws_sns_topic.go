package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSnsTopic(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sns_topic",
		Description: "AWS SNS Topic",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("topic_arn"),
			ShouldIgnoreError: isNotFoundError([]string{"NotFound", "InvalidParameter"}),
			ItemFromKey:       snsTopicFromKey,
			Hydrate:           getTopicAttributes,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSnsTopics,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "topic_arn",
				Description: "Amazon Resource Name (ARN) of the Topic",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.TopicArn"),
			},
			{
				Name:        "display_name",
				Description: "The human-readable name used in the From field for notifications to email and email-json endpoints",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.DisplayName"),
			},
			{
				Name:        "owner",
				Description: "The AWS account ID of the topic's owner",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.Owner"),
			},
			{
				Name:        "subscriptions_confirmed",
				Description: "The number of confirmed subscriptions for the topic",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.SubscriptionsConfirmed"),
			},
			{
				Name:        "subscriptions_deleted",
				Description: "The number of deleted subscriptions for the topic",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.SubscriptionsDeleted"),
			},
			{
				Name:        "subscriptions_pending",
				Description: "The number of subscriptions pending confirmation for the topic",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.SubscriptionsPending"),
			},
			{
				Name:        "kms_master_key_id",
				Description: "The ID of an AWS-managed customer master key (CMK) for Amazon SNS or a custom CMK",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.KmsMasterKeyId"),
			},
			{
				Name:        "tags_src",
				Description: "The list of tags associated with the topic",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForSnsTopic,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "policy",
				Description: "The topic's access control policy (i.e. Resource IAM Policy)",
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
				Description: "The JSON object of the topic's delivery policy",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTopicAttributes,
				Transform:   transform.FromField("Attributes.DeliveryPolicy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "effective_delivery_policy",
				Description: "The effective delivery policy, taking system defaults into account",
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
				Transform:   transform.From(snsTopicTurbotTags),
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

//// BUILD HYDRATE INPUT

func snsTopicFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	arn := quals["topic_arn"].GetStringValue()
	item := &sns.GetTopicAttributesOutput{
		Attributes: map[string]*string{
			"TopicArn": &arn,
		},
	}
	return item, nil
}

//// LIST FUNCTION

func listAwsSnsTopics(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listAwsSnsTopics", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := SNSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListTopicsPages(
		&sns.ListTopicsInput{},
		func(page *sns.ListTopicsOutput, lastPage bool) bool {
			for _, topic := range page.Topics {
				d.StreamListItem(ctx, &sns.GetTopicAttributesOutput{
					Attributes: map[string]*string{
						"TopicArn": topic.TopicArn,
					},
				})
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getTopicAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getTopicAttributes")
	topicAttributesOutput := h.Item.(*sns.GetTopicAttributesOutput)
	defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := SNSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// Build params
	param := &sns.GetTopicAttributesInput{
		TopicArn: topicAttributesOutput.Attributes["TopicArn"],
	}

	op, err := svc.GetTopicAttributes(param)
	if err != nil {
		logger.Trace("getTopicAttributes__", "Error", err)
		return nil, err
	}
	return op, nil
}

func listTagsForSnsTopic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listTagsForSnsTopic")
	topicAttributesOutput := h.Item.(*sns.GetTopicAttributesOutput)
	defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := SNSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// Build param
	param := &sns.ListTagsForResourceInput{
		ResourceArn: topicAttributesOutput.Attributes["TopicArn"],
	}

	// Next token is not supported
	// AWS supports upto 50 tags
	topicTags, err := svc.ListTagsForResource(param)

	if err != nil {
		return nil, err
	}
	return topicTags, nil
}

//// TRANSFORM FUNCTIONS

func snsTopicTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*sns.ListTagsForResourceOutput)
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}
