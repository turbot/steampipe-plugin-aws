package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/shield"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsShieldSubscription(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_shield_subscription",
		Description: "AWS Shield Subscription",
		List: &plugin.ListConfig{
			Hydrate: listAwsShieldSubscription,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags:    map[string]string{"service": "shield", "action": "DescribeSubscription"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsShieldSubscriptionState,
				Tags: map[string]string{"service": "shield", "action": "GetSubscriptionState"},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "subscription_state",
				Description: "The current state the subscription.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsShieldSubscriptionState,
			},
			{
				Name:        "start_time",
				Description: "The start time of the subscription.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "end_time",
				Description: "The date and time your subscription will end.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "time_commitment_in_seconds",
				Description: "The length, in seconds, of the Shield Advanced subscription for the account.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "auto_renew",
				Description: "If ENABLED, the subscription will be automatically renewed at the end of the existing subscription period.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "proactive_engagement_status",
				Description: "Status of the proactive engagement of the Shield Response Team (SRT). Indicates if the Shield Response Team (SRT) will use the Shield emergency contact data to notify them about DDoS attacks.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProactiveEngagementStatus").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) of the subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubscriptionArn"),
			},
			{
				Name:        "subscription_limits",
				Description: "The configured limits for your subscription.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "limits",
				Description: "Specifies how many protections of a given type you can create.",
				Type:        proto.ColumnType_JSON,
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubscriptionArn"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SubscriptionArn").Transform(arnToAkas),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listAwsShieldSubscription(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_subscription.listAwsShieldSubscription", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	data, err := svc.DescribeSubscription(ctx, &shield.DescribeSubscriptionInput{})

	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_subscription.getAwsShieldSubscription", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, data.Subscription)

	return nil, nil
}

func getAwsShieldSubscriptionState(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_subscription.getAwsShieldSubscriptionState", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	data, err := svc.GetSubscriptionState(ctx, &shield.GetSubscriptionStateInput{})

	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_subscription.getAwsShieldSubscriptionState", "api_error", err)
		return nil, err
	}

	if data != nil {
		return data, nil
	}

	return nil, nil
}