package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsEventBridgeRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eventbridge_rule",
		Description: "AWS EventBridge Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getAwsEventBridgeRule,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsEventBridgeBuses,
			Hydrate:       listAwsEventBridgeRules,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "event_bus_name", Require: plugin.Optional},
				{Name: "name_prefix", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Hydrate:     getAwsEventBridgeRule,
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
				Hydrate:     getAwsEventBridgeRule,
			},
			{
				Name:        "name_prefix",
				Description: "Specifying this limits the results to only those event rules with names that start with the specified prefix.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNamePrefixValue,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "targets",
				Description: "The targets assigned to the rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEventBridgeTargetByRule,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEventBridgeRuleTags,
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
				Hydrate:     getAwsEventBridgeRuleTags,
				Transform:   transform.FromField("Tags").Transform(eventBridgeTagListToTurbotTags),
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

//// LIST FUNCTION

func listAwsEventBridgeRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var eventBusName string
	if h.Item != nil {
		data := h.Item.(*eventbridge.DescribeEventBusOutput)
		eventBusName = types.SafeString(data.Name)
	} else {
		eventBusName = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Get client
	svc, err := EventBridgeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_rule.listAwsEventBridgeRules", "get_client_error", err)
		return nil, err
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
	if eventBusName != "" {
		params.EventBusName = &eventBusName
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["name_prefix"] != nil {
		params.NamePrefix = aws.String(equalQuals["name_prefix"].GetStringValue())
	}

	for pagesLeft {
		output, err := svc.ListRules(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eventbridge_rule.listAwsEventBridgeRules", "api_error", err)
			return nil, err
		}
		for _, rule := range output.Rules {
			d.StreamListItem(ctx, &eventbridge.DescribeRuleOutput{
				Name:         rule.Name,
				Arn:          rule.Arn,
				Description:  rule.Description,
				State:        rule.State,
				EventBusName: rule.EventBusName,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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

func getAwsEventBridgeRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var name string
	if h.Item != nil {
		name = *h.Item.(*eventbridge.DescribeRuleOutput).Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := EventBridgeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_bus.getAwsEventBridgeRule", "get_client_error", err)
		return nil, err
	}
	// Build the params
	params := &eventbridge.DescribeRuleInput{
		Name: &name,
	}

	// Get call
	data, err := svc.DescribeRule(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_rule.getAwsEventBridgeRule", "api_error", err)
		return nil, err
	}

	return data, nil
}

func getAwsEventBridgeTargetByRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	eventbusname := h.Item.(*eventbridge.DescribeRuleOutput).EventBusName
	name := h.Item.(*eventbridge.DescribeRuleOutput).Name

	// Create Session
	svc, err := EventBridgeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_bus.getAwsEventBridgeTargetByRule", "get_client_error", err)
		return nil, err
	}
	// Build the params
	params := &eventbridge.ListTargetsByRuleInput{
		EventBusName: eventbusname,
		Rule:         name,
	}

	data, err := svc.ListTargetsByRule(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_rule.getAwsEventBridgeTargetByRule", "api_error", err)
		return nil, err
	}

	return data, nil
}

func getAwsEventBridgeRuleTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	arn := h.Item.(*eventbridge.DescribeRuleOutput).Arn

	// Create Session
	svc, err := EventBridgeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_bus.getAwsEventBridgeRuleTags", "get_client_error", err)
		return nil, err
	}

	// Build the params
	params := &eventbridge.ListTagsForResourceInput{
		ResourceARN: arn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_rule.getAwsEventBridgeRuleTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getNamePrefixValue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if d.KeyColumnQuals["name_prefix"].GetStringValue() != "" {
		return d.KeyColumnQuals["name_prefix"].GetStringValue(), nil
	} else {
		return h.Item.(*eventbridge.DescribeRuleOutput).Name, nil
	}
}

//// TRANSFORM FUNCTIONS

func eventBridgeTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("eventBridgeTagListToTurbotTags")
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
