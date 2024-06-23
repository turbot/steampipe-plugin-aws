package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codestarnotifications"
	"github.com/aws/aws-sdk-go-v2/service/codestarnotifications/types"

	turbot_types "github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodestarNotificationRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codestar_notification_rule",
		Description: "AWS CodeStar Notification Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getCodeStarNotificationRule,
			Tags:    map[string]string{"service": "codestar-notifications", "action": "DescribeNotificationRule"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeStarNotificationRules,
			Tags:    map[string]string{"service": "codestar-notifications", "action": "ListNotificationRules"},
		},
		HydrateConfig:     []plugin.HydrateConfig{},
		GetMatrixItemFunc: SupportedRegionMatrix("codestar-notifications"), // unsure how to get EndpointsID from the package
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the notification rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique ID of the notification rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the notification rule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeStarNotificationRule,
			},
			{
				Name:        "resource",
				Description: "The Amazon Resource Name (ARN) of the resource associated with the notification rule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeStarNotificationRule,
			},
			{
				Name:        "detail_type",
				Description: "The level of detail included in the notifications for this resource. BASIC will include only the contents of the event as it would appear in Amazon CloudWatch. FULL will include any supplemental information provided by AWS CodeStar Notifications and/or the service for the resource for which the notification is created.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeStarNotificationRule,
			},
			{
				Name:        "status",
				Description: "The status of the notification rule. Valid statuses are on (sending notifications) or off (not sending notifications).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeStarNotificationRule,
			},
			{
				Name:        "created_by",
				Description: "The name or email alias of the person who created the notification rule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeStarNotificationRule,
			},
			{
				Name:        "event_type_id",
				Description: "Specifies that only notification rules with the given event type enabled are returned.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("event_type_id"),
			},
			{
				Name:        "target_address",
				Description: "Specifies that only notification rules with a target with the given address are returned.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("target_address"),
			},
			{
				Name:        "created_timestamp",
				Description: "The date and time the notification rule was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCodeStarNotificationRule,
			},
			{
				Name:        "last_modified_timestamp",
				Description: "The date and time the notification rule was most recently updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCodeStarNotificationRule,
			},
			{
				Name:        "event_types",
				Description: "A list of the event types associated with the notification rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeStarNotificationRule,
			},
			{
				Name:        "targets",
				Description: "A list of targets associated with the notification rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeStarNotificationRule,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeStarNotificationRule,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeStarNotificationRule,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeStarNotificationRule,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listCodeStarNotificationRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CodeStarNotificationsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codestar_notification_rule.listCodeStarNotificationRules", "connection_error", err)
		return nil, err
	}

	params := &codestarnotifications.ListNotificationRulesInput{}

	filters := buildCodeStarNotificationRulesFilter(d.EqualsQuals)
	if len(filters) != 0 {
		params.Filters = filters
	}

	paginator := codestarnotifications.NewListNotificationRulesPaginator(svc, params, func(o *codestarnotifications.ListNotificationRulesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codestar_notification_rule.listCodeStarNotificationRules", "api_error", err)
			return nil, err
		}
		for _, rule := range output.NotificationRules {
			d.StreamListItem(ctx, &codestarnotifications.DescribeNotificationRuleOutput{
				Arn: rule.Arn,
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

func getCodeStarNotificationRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		data := h.Item.(*codestarnotifications.DescribeNotificationRuleOutput)
		arn = turbot_types.SafeString(data.Arn)
	} else {
		arn = d.EqualsQuals["arn"].GetStringValue()
	}

	if arn == "" {
		return nil, nil
	}

	// Get client
	svc, err := CodeStarNotificationsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codestar_notification_rule.getCodeStarNotificationRule", "connection_error", err)
		return nil, err
	}

	// Build params
	params := &codestarnotifications.DescribeNotificationRuleInput{
		Arn: aws.String(arn),
	}

	op, err := svc.DescribeNotificationRule(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codestar_notification_rule.getCodeStarNotificationRule", "api_error", err)
		return nil, err
	}
	return op, nil
}

//// UTILITY FUNCTIONS

func buildCodeStarNotificationRulesFilter(equalQuals plugin.KeyColumnEqualsQualMap) []types.ListNotificationRulesFilter {
	filters := make([]types.ListNotificationRulesFilter, 0)

	filterQuals := map[string]types.ListNotificationRulesFilterName{
		"event_type_id":  "EVENT_TYPE_ID",
		"created_by":     "CREATED_BY",
		"resource":       "RESOURCE",
		"target_address": "TARGET_ADDRESS",
	}

	for columnName, filterName := range filterQuals {
		if equalQuals[columnName] != nil {
			filter := types.ListNotificationRulesFilter{
				Name:  filterName,
				Value: aws.String(equalQuals[columnName].GetStringValue()),
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
