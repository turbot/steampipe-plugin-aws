package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION
func tableAwsWafRateBasedRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_waf_rate_based_rule",
		Description: "AWS WAF Rate Based Rule",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("rule_id"),
			ShouldIgnoreError: isNotFoundError([]string{"WAFNonexistentItemException", "ValidationException"}),
			Hydrate:           getAwsWafRateBasedRule,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWafRateBasedRules,
		},
		Columns: awsColumns([]*plugin.Column{
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
				Description: "The name or description for the metrics for a RateBasedRule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsWafRateBasedRule,
			},
			{
				Name:        "rate_key",
				Description: "The field that AWS WAF uses to determine if requests are likely arriving from single source and thus subject to rate monitoring.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsWafRateBasedRule,
			},
			{
				Name:        "rate_limit",
				Description: "The maximum number of requests, which have an identical value in the field specified by the RateKey, allowed in a five-minute period.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsWafRateBasedRule,
			},
			{
				Name:        "predicates",
				Description: "The Predicates object contains one Predicate element for each ByteMatchSet, IPSet or SqlInjectionMatchSet object that you want to include in a RateBasedRule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafRateBasedRule,
				Transform:   transform.FromField("MatchPredicates"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsWafRateBasedRuleTags,
				Transform:   transform.FromField("TagInfoForResource.TagList"),
			},

			// Steampipe standard columns
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
				Hydrate:     listAwsWafRateBasedRuleTags,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(wafRateBasedRuletagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafRateBasedRuleAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsWafRateBasedRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsWafRateBasedRules")
	// Create session
	svc, err := WAFService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	params := &waf.ListRateBasedRulesInput{}
	pagesLeft := true
	for pagesLeft {
		response, err := svc.ListRateBasedRules(params)
		if err != nil {
			return nil, err
		}
		for _, rule := range response.Rules {
			d.StreamListItem(ctx, rule)
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

////  HYDRATE FUNCTIONS

func getAwsWafRateBasedRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsWafRateBasedRule")
	var id string
	if h.Item != nil {
		id = rateBasedRuleData(h.Item)
	} else {
		id = d.KeyColumnQuals["rule_id"].GetStringValue()
	}
	// Create Session
	svc, err := WAFService(ctx, d)
	if err != nil {
		return nil, err
	}
	// Build the params
	params := &waf.GetRateBasedRuleInput{
		RuleId: &id,
	}
	// Get call
	data, err := svc.GetRateBasedRule(params)
	if err != nil {
		logger.Debug("getAwsWafRateBasedRule", "ERROR", err)
		return nil, err
	}
	return data.Rule, nil
}

func listAwsWafRateBasedRuleTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsWafRateBasedRuleTags")

	var id string
	if h.Item != nil {
		id = rateBasedRuleData(h.Item)
	} else {
		id = d.KeyColumnQuals["rule_id"].GetStringValue()
	}

	commonAwsColumns, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonAwsColumns.(*awsCommonColumnData)
	// Create Session
	svc, err := WAFService(ctx, d)
	if err != nil {
		return nil, err
	}
	aka := "arn:" + commonColumnData.Partition + ":waf::" + commonColumnData.AccountId + ":ratebasedrule" + "/" + id

	params := &waf.ListTagsForResourceInput{
		ResourceARN: &aka,
	}
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}
	return op, nil
}

func getAwsWafRateBasedRuleAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsWafRateBasedRuleAkas")

	id := rateBasedRuleData(h.Item)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":waf::" + commonColumnData.AccountId + ":ratebasedrule/" + id
	return []string{aka}, nil
}

//// TRANSFORM FUNCTIONS

func wafRateBasedRuletagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("wafRateBasedRuletagListToTurbotTags")
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

func rateBasedRuleData(item interface{}) string {
	switch item := item.(type) {
	case *waf.RuleSummary:
		return *item.RuleId
	case *waf.RateBasedRule:
		return *item.RuleId
	}
	return ""
}
