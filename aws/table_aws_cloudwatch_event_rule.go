package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudwatchEventRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_event_rule",
		Description: "AWS CloudWatch Event Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "event_bus_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getAwsCloudWatchEventRule,
			Tags:    map[string]string{"service": "events", "action": "DescribeRule"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsCloudWatchEventRules,
			Tags:    map[string]string{"service": "events", "action": "ListRules"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "event_bus_name", Require: plugin.Optional},
				{Name: "name_prefix", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsCloudWatchEventTargetsByRule,
				Tags: map[string]string{"service": "events", "action": "ListTargetsByRule"},
			},
			{
				Func: getAwsCloudWatchEventRuleTags,
				Tags: map[string]string{"service": "events", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_EVENTS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_bus_name",
				Description: "The name or ARN of the event bus associated with the rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_by",
				Description: "The account ID of the user that created the rule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsCloudWatchEventRule,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role associated with the rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schedule_expression",
				Description: "The scheduling expression. For example, 'cron(0 20 * * ? *)', 'rate(5 minutes)'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "managed_by",
				Description: "If this is a managed rule, created by an AWS service on your behalf, this field displays the principal name of the AWS service that created the rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_pattern",
				Description: "The event pattern of the rule.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "name_prefix",
				Description: "Specifying this limits the results to only those event rules with names that start with the specified prefix.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudWatchNamePrefixValue,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "targets",
				Description: "The targets assigned to the rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsCloudWatchEventTargetsByRule,
				Transform:   transform.FromField("Targets"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsCloudWatchEventRuleTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
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
				Hydrate:     getAwsCloudWatchEventRuleTags,
				Transform:   transform.FromField("Tags").Transform(cloudWatchEventTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsCloudWatchEventRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := EventBridgeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_event_rule.listAwsCloudWatchEventRules", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	pagesLeft := true
	params := &eventbridge.ListRulesInput{
		// Default to the maximum allowed
		Limit: aws.Int32(maxLimit),
	}

	equalQuals := d.EqualsQuals
	if equalQuals["event_bus_name"] != nil {
		params.EventBusName = aws.String(equalQuals["event_bus_name"].GetStringValue())
	}
	if equalQuals["name_prefix"] != nil {
		params.NamePrefix = aws.String(equalQuals["name_prefix"].GetStringValue())
	}
	// API doesn't support aws-go-sdk-v2 paginator as of date
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := svc.ListRules(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_event_rule.listAwsCloudWatchEventRules", "api_error", err)
			return nil, err
		}
		for _, rule := range output.Rules {
			d.StreamListItem(ctx, &eventbridge.DescribeRuleOutput{
				Name:               rule.Name,
				Arn:                rule.Arn,
				Description:        rule.Description,
				State:              rule.State,
				EventBusName:       rule.EventBusName,
				ManagedBy:          rule.ManagedBy,
				ScheduleExpression: rule.ScheduleExpression,
				RoleArn:            rule.RoleArn,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if output.NextToken != nil {
			pagesLeft = true
			params.NextToken = output.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsCloudWatchEventRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name, eventBusName string

	if h.Item != nil {
		name = *h.Item.(*eventbridge.DescribeRuleOutput).Name
		if h.Item.(*eventbridge.DescribeRuleOutput).EventBusName != nil {
			eventBusName = *h.Item.(*eventbridge.DescribeRuleOutput).EventBusName
		}
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
		eventBusName = d.EqualsQuals["event_bus_name"].GetStringValue()
	}

	// Create Session
	svc, err := EventBridgeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_event_rule.getAwsCloudWatchEventRule", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &eventbridge.DescribeRuleInput{
		Name: &name,
	}

	// Use the event bus name if specified
	if eventBusName != "" {
		params.EventBusName = &eventBusName
	}

	// Get call
	data, err := svc.DescribeRule(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_event_rule.getAwsCloudWatchEventRule", "api_error", err)
		return nil, err
	}

	return data, nil
}

func getAwsCloudWatchEventTargetsByRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	rule := h.Item.(*eventbridge.DescribeRuleOutput)
	name := rule.Name
	var eventBusName *string

	if rule.EventBusName != nil {
		eventBusName = rule.EventBusName
	}

	// Create Session
	svc, err := EventBridgeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_event_rule.getAwsCloudWatchEventTargetsByRule", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &eventbridge.ListTargetsByRuleInput{
		Rule: name,
	}

	// Use event bus name if available
	if eventBusName != nil {
		params.EventBusName = eventBusName
	}

	data, err := svc.ListTargetsByRule(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_event_rule.getAwsCloudWatchEventTargetsByRule", "api_error", err)
		return nil, err
	}

	return data, nil
}

func getAwsCloudWatchEventRuleTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := h.Item.(*eventbridge.DescribeRuleOutput).Arn

	// Create Session
	svc, err := EventBridgeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_event_rule.getAwsCloudWatchEventRuleTags", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &eventbridge.ListTagsForResourceInput{
		ResourceARN: arn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_event_rule.getAwsCloudWatchEventRuleTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func cloudWatchEventTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.HydrateItem.(*eventbridge.ListTagsForResourceOutput)

	if tagList.Tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

//// UTILITY FUNCTIONS

func getCloudWatchNamePrefixValue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if d.EqualsQuals["name_prefix"].GetStringValue() != "" {
		return d.EqualsQuals["name_prefix"].GetStringValue(), nil
	} else {
		return h.Item.(*eventbridge.DescribeRuleOutput).Name, nil
	}
}
