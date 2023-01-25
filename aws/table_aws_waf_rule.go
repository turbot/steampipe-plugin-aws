package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/waf"
	"github.com/aws/aws-sdk-go-v2/service/waf/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsWAFRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_waf_rule",
		Description: "AWS WAF Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("rule_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"WAFNonexistentItemException"}),
			},
			Hydrate: getAwsWAFRule,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWAFRules,
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name for the rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "rule_id",
				Description: "The ID of the Rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metric_name",
				Description: "The name of the metric for the Rule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsWAFRule,
			},
			{
				Name:        "predicates",
				Description: "The Predicates object contains one Predicate element for each ByteMatchSet,IPSet, or SqlInjectionMatchSet object that you want to include in a Rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWAFRule,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWAFRuleTags,
				Transform:   transform.FromField("TagInfoForResource.TagList"),
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
				Hydrate:     getAwsWAFRuleTags,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(wafRuleTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWAFRuleAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsWAFRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := WAFClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule.listAwsWAFRules", "get_client_error", err)
		return nil, err
	}

	// List call
	params := &waf.ListRulesInput{
		Limit: int32(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	// Minimunm limit is 0
	// https://docs.aws.amazon.com/waf/latest/APIReference/API_waf_ListRules.html
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(params.Limit) {
			params.Limit = int32(*limit)
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := svc.ListRules(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_waf_rule.listAwsWAFRules", "api_error", err)
			return nil, err
		}
		for _, rule := range response.Rules {
			d.StreamListItem(ctx, rule)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if response.NextMarker != nil {
			pagesLeft = true
			params.NextMarker = response.NextMarker
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsWAFRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := WAFClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule.getAwsWAFRule", "get_client_error", err)
		return nil, err
	}

	var id string
	if h.Item != nil {
		id = ruleData(h.Item)
	} else {
		id = d.KeyColumnQuals["rule_id"].GetStringValue()
	}

	// Build the params
	param := &waf.GetRuleInput{
		RuleId: &id,
	}

	// Get call
	data, err := svc.GetRule(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule.getAwsWAFRule", "get_client_error", err)
		return nil, err
	}

	return data.Rule, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go/issues/3513
func getAwsWAFRuleTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// id := ruleData(h.Item)

	var id string

	if h.Item != nil {
		id = ruleData(h.Item)
	} else {
		id = d.KeyColumnQuals["rule_id"].GetStringValue()
	}

	commonAwsColumns, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule.getAwsWAFRuleTags", "api_error", err)
		return nil, err
	}
	commonColumnData := commonAwsColumns.(*awsCommonColumnData)

	// Create Session
	svc, err := WAFClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule.getAwsWAFRuleTags", "get_client_error", err)
		return nil, err
	}

	aka := "arn:" + commonColumnData.Partition + ":waf::" + commonColumnData.AccountId + ":rule" + "/" + id

	// Build param with maximum limit set
	params := &waf.ListTagsForResourceInput{
		ResourceARN: &aka,
		Limit:       int32(100),
	}
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.listWafRuleGroups", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsWAFRuleAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := ruleData(h.Item)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule.getAwsWAFRuleAkas", "api_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":waf::" + commonColumnData.AccountId + ":rule/" + id

	return []string{aka}, nil
}

//// TRANSFORM FUNCTION

func wafRuleTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.HydrateItem.(*waf.ListTagsForResourceOutput)

	if tagList.TagInfoForResource.TagList == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList.TagInfoForResource.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func ruleData(item interface{}) string {
	switch item := item.(type) {
	case types.RuleSummary:
		return *item.RuleId
	case *types.Rule:
		return *item.RuleId
	}
	return ""
}
