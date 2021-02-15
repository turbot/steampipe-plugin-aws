package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
)

func tableAwsCloudwatchEventRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_event_rule",
		Description: "AWS CloudWatch Event Rule",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"NoSuchEntity", "InvalidParameter"}),
			ItemFromKey:       eventRuleFromKey,
			Hydrate:           getCloudwatchEventRule,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudwatchEventRules,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the event rule",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the event rule",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the role that is used for target invocation",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleArn"),
			},
			{
				Name:        "description",
				Description: "The description of the rule",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_bus_name",
				Description: "The name or ARN of the event bus associated with the rule",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "managed_by",
				Description: "If the rule was created on behalf of your account by an AWS service, this field displays the principal name of the service that created the rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedBy"),
			},
			{
				Name:        "event_pattern",
				Description: "The event pattern of the rule",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schedule_expression",
				Description: "The scheduling expression",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ScheduleExpression"),
			},
			{
				Name:        "state",
				Description: "The state of the rule",
				Type:        proto.ColumnType_STRING,
			},

			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEventRuleTagging,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// BUILD HYDRATE INPUT

func eventRuleFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &cloudwatchevents.Rule{
		Name: &name,
	}
	return item, nil
}

//// LIST FUNCTION

func listCloudwatchEventRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listCloudwatchEventRules", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := CloudWatchEventsService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	ruleOutput, err := svc.ListRules(&cloudwatchevents.ListRulesInput{})
	for _, rule := range ruleOutput.Rules {
		d.StreamListItem(ctx, rule)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudwatchEventRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudwatchEventRule")

	defaultRegion := GetDefaultRegion()
	rule := h.Item.(*cloudwatchevents.Rule)

	// Create session
	svc, err := CloudWatchEventsService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &cloudwatchevents.ListRulesInput{
		NamePrefix: rule.Name,
	}

	// execute list call
	item, err := svc.ListRules(params)
	if err != nil {
		return nil, err
	}
	for _, rule := range item.Rules {
		return rule, nil
	}
	return nil, nil
}

func getEventRuleTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEventRuleTagging")
	rule := h.Item.(*cloudwatchevents.Rule)
	defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := CloudWatchEventsService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &cloudwatchevents.ListTagsForResourceInput{
		ResourceARN: rule.Arn,
	}

	// List resource tags
	ruleTags, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}
	return ruleTags, nil
}
