package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsEventBridgeRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eventbridge_rule",
		Description: "AWS EventBridge Rule",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException", "ValidationException"}),
			Hydrate:           getAwsEventBridgeRule,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsEventBridgeRules,
		},
		GetMatrixItem: BuildRegionList,
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

func listAwsEventBridgeRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsEventBridgeRules", "AWS_REGION", region)

	// Create session
	svc, err := EventBridgeService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	param := &eventbridge.ListRulesInput{}
	for {
		response, err := svc.ListRules(&eventbridge.ListRulesInput{})
		if err != nil {
			return nil, err
		}
		for _, rule := range response.Rules {
			d.StreamListItem(ctx, &eventbridge.DescribeRuleOutput{
				Name:         rule.Name,
				Arn:          rule.Arn,
				Description:  rule.Description,
				State:        rule.State,
				EventBusName: rule.EventBusName,
			})
		}
		if response.NextToken == nil {
			break
		}
		param.NextToken = response.NextToken
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsEventBridgeRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsEventBridgeRule")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var name string
	if h.Item != nil {
		name = *h.Item.(*eventbridge.DescribeRuleOutput).Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := EventBridgeService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	// Build the params
	params := &eventbridge.DescribeRuleInput{
		Name: &name,
	}

	// Get call
	data, err := svc.DescribeRule(params)
	if err != nil {
		logger.Debug("getAwsEventBridgeRule", "ERROR", err)
		return nil, err
	}

	return data, nil
}

func getAwsEventBridgeTargetByRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsEventBridgeTargetByRule")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	eventbusname := h.Item.(*eventbridge.DescribeRuleOutput).EventBusName
	name := h.Item.(*eventbridge.DescribeRuleOutput).Name

	// Create Session
	svc, err := EventBridgeService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	// Build the params
	params := &eventbridge.ListTargetsByRuleInput{
		EventBusName: eventbusname,
		Rule:         name,
	}

	data, err := svc.ListTargetsByRule(params)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getAwsEventBridgeRuleTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsEventBridgeRuleTags")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	arn := h.Item.(*eventbridge.DescribeRuleOutput).Arn

	// Create Session
	svc, err := EventBridgeService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &eventbridge.ListTagsForResourceInput{
		ResourceARN: arn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Debug("getAwsEventBridgeRuleTags", "ERROR", err)
		return nil, err
	}

	return op, nil
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
