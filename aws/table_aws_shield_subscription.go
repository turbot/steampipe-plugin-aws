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
			Tags:    map[string]string{"service": "shield", "action": "DescribeSubscription"},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "start_time",
				Description: "The start time of the subscription.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Subscription.StartTime"),
			},
			{
				Name:        "end_time",
				Description: "The date and time your subscription will end.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Subscription.EndTime"),
			},
			{
				Name:        "time_commitment_in_seconds",
				Description: "The length, in seconds, of the Shield Advanced subscription for the account.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Subscription.TimeCommitmentInSeconds"),
			},
			{
				Name:        "auto_renew",
				Description: "If ENABLED, the subscription will be automatically renewed at the end of the existing subscription period.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subscription.AutoRenew"),
			},
			{
				Name:        "limits",
				Description: "Specifies how many protections of a given type you can create.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Subscription.Limits"),
			},
			{
				Name:        "proactive_engagement_status",
				Description: "If the Shield Response Team (SRT) will use email and phone to notify contacts about escalations to the SRT and to initiate proactive customer support.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subscription.ProactiveEngagementStatus"),
			},
			{
				Name:        "subscription_limits",
				Description: "The configured limits for your subscription.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Subscription.SubscriptionLimits"),
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) of the subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subscription.SubscriptionArn"),
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subscription.SubscriptionArn"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Subscription.SubscriptionArn").Transform(arnToAkas),
			},
		}),
	}
}

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

	d.StreamListItem(ctx, data)

	return nil, nil
}