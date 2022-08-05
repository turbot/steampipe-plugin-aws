package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsConfigRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_rule",
		Description: "AWS Config Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchConfigRuleException", "ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getConfigRule,
		},
		List: &plugin.ListConfig{
			Hydrate: listConfigRules,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name that you assign to the AWS Config rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigRuleName"),
			},
			{
				Name:        "rule_id",
				Description: "The ID of the AWS Config rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigRuleId"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the AWS Config rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigRuleArn"),
			},
			{
				Name:        "rule_state",
				Description: "It indicate the evaluation status for the AWS Config rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigRuleState"),
			},
			{
				Name:        "created_by",
				Description: "Service principal name of the service that created the rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description that you provide for the AWS Config rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "maximum_execution_frequency",
				Description: "The maximum frequency with which AWS Config runs evaluations for a rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "compliance_by_config_rule",
				Description: "The compliance information of the config rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComplianceByConfigRules,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "input_parameters",
				Description: "A string, in JSON format, that is passed to the AWS Config rule Lambda function.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "scope",
				Description: "Defines which resources can trigger an evaluation for the rule. The scope can include one or more resource types, a combination of one resource type and one resource ID, or a combination of a tag key and value. Specify a scope to constrain the resources that can trigger an evaluation for the rule. If you do not specify a scope, evaluations are triggered when any resource in the recording group changes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source",
				Description: "Provides the rule owner (AWS or customer), the rule identifier, and the notifications that cause the function to evaluate your AWS resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getConfigRuleTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigRuleName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getConfigRuleTags,
				Transform:   transform.From(configRuleTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConfigRuleArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listConfigRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ConfigService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &configservice.DescribeConfigRulesInput{}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.ConfigRuleNames = []*string{aws.String(equalQuals["name"].GetStringValue())}
	}

	err = svc.DescribeConfigRulesPages(
		input,
		func(page *configservice.DescribeConfigRulesOutput, lastPage bool) bool {
			for _, rule := range page.ConfigRules {
				d.StreamListItem(ctx, rule)

				// Context can be cancelled due to manual cancellation or the limit has been hit
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

func getConfigRule(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getConfigRule")

	// Create Session
	svc, err := ConfigService(ctx, d)
	if err != nil {
		return nil, err
	}
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Build params
	params := &configservice.DescribeConfigRulesInput{
		ConfigRuleNames: []*string{aws.String(name)},
	}

	op, err := svc.DescribeConfigRules(params)
	if err != nil {
		return nil, err
	}

	if op != nil {
		return op.ConfigRules[0], nil
	}

	return nil, nil
}

func getConfigRuleTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getConfigRuleTags")

	// Create Session
	svc, err := ConfigService(ctx, d)
	if err != nil {
		return nil, err
	}
	ruleArn := h.Item.(*configservice.ConfigRule).ConfigRuleArn

	// Build params
	params := &configservice.ListTagsForResourceInput{
		ResourceArn: ruleArn,
	}

	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getComplianceByConfigRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComplianceByConfigRules")

	// Create Session
	svc, err := ConfigService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getComplianceByConfigRules", "connection", err)
		return nil, err
	}
	ruleName := h.Item.(*configservice.ConfigRule).ConfigRuleName

	// Build params
	params := &configservice.DescribeComplianceByConfigRuleInput{
		ConfigRuleNames: []*string{ruleName},
	}

	op, err := svc.DescribeComplianceByConfigRule(params)
	if err != nil {
		plugin.Logger(ctx).Error("getComplianceByConfigRules", "DescribeComplianceByConfigRule", err)
		return nil, err
	}

	return op.ComplianceByConfigRules, nil
}

//// TRANSFORM FUNCTIONS

func configRuleTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*configservice.ListTagsForResourceOutput)

	if data.Tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	turbotTagsMap := map[string]string{}
	for _, i := range data.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return turbotTagsMap, nil
}
