package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsConfigRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_rule",
		Description: "AWS Config Rule",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"NoSuchConfigRuleException", "ResourceNotFoundException", "ValidationException"}),
			Hydrate:           getConfigRule,
		},
		List: &plugin.ListConfig{
			Hydrate: listConfigRules,
		},
		GetMatrixItem: BuildRegionList,
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
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listConfigRules", "AWS_REGION", region)

	// Create Session
	svc, err := ConfigService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	op, err := svc.DescribeConfigRules(&configservice.DescribeConfigRulesInput{})
	if err != nil {
		return nil, err
	}

	for _, rule := range op.ConfigRules {
		d.StreamListItem(ctx, rule)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getConfigRule(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getConfigRule")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := ConfigService(ctx, d, region)
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

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := ConfigService(ctx, d, region)
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
