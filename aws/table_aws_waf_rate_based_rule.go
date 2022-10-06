package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/waf"
	"github.com/aws/aws-sdk-go-v2/service/waf/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION
func tableAwsWafRateBasedRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_waf_rate_based_rule",
		Description: "AWS WAF Rate Based Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("rule_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"WAFNonexistentItemException", "ValidationException"}),
			},
			Hydrate: getAwsWafRateBasedRule,
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
	// Create session
	svc, err := WAFClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rate_based_rule.listAwsWafRateBasedRules", "get_client_error", err)
		return nil, err
	}

	// List call
	params := &waf.ListRateBasedRulesInput{
		Limit: int32(20),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	// Minimunm limit is 0
	// https://docs.aws.amazon.com/waf/latest/APIReference/API_waf_ListRateBasedRules.html
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(params.Limit) {
			params.Limit = int32(*limit)
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := svc.ListRateBasedRules(ctx, params)
		if err != nil {
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

////  HYDRATE FUNCTIONS

func getAwsWafRateBasedRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string

	if h.Item != nil {
		id = rateBasedRuleData(h.Item)
	} else {
		id = d.KeyColumnQuals["rule_id"].GetStringValue()
	}

	// Create Session
	svc, err := WAFClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rate_based_rule.getAwsWafRateBasedRule", "get_client_error", err)
		return nil, err
	}

	// Build the params
	params := &waf.GetRateBasedRuleInput{
		RuleId: &id,
	}

	// Get call
	data, err := svc.GetRateBasedRule(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rate_based_rule.getAwsWafRateBasedRule", "api_error", err)
		return nil, err
	}

	return data.Rule, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go/issues/3513
func listAwsWafRateBasedRuleTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
	svc, err := WAFClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rate_based_rule.listAwsWafRateBasedRuleTags", "get_client_error", err)
		return nil, err
	}

	aka := "arn:" + commonColumnData.Partition + ":waf::" + commonColumnData.AccountId + ":ratebasedrule" + "/" + id

	// Build param with maximum limit set
	params := &waf.ListTagsForResourceInput{
		ResourceARN: &aka,
		Limit:       int32(100),
	}
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getAwsWafRateBasedRuleAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := rateBasedRuleData(h.Item)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)

	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rate_based_rule.getAwsWafRateBasedRuleAkas", "get_client_error", err)
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":waf::" + commonColumnData.AccountId + ":ratebasedrule/" + id

	return []string{aka}, nil
}

//// TRANSFORM FUNCTIONS

func wafRateBasedRuletagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
	case types.RuleSummary:
		return *item.RuleId
	case *types.RateBasedRule:
		return *item.RuleId
	}

	return ""
}
