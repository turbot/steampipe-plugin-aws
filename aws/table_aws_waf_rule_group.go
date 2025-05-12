package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/waf"
	"github.com/aws/aws-sdk-go-v2/service/waf/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWafRuleGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_waf_rule_group",
		Description: "AWS WAF Rule Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"rule_group_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NonexistentItemException", "WAFNonexistentItemException"}),
			},
			Hydrate: getWafRuleGroup,
			Tags:    map[string]string{"service": "waf", "action": "GetRuleGroup"},
		},
		List: &plugin.ListConfig{
			Hydrate: listWafRuleGroups,
			Tags:    map[string]string{"service": "waf", "action": "ListRuleGroups"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getWafRuleGroup,
				Tags: map[string]string{"service": "waf", "action": "GetRuleGroup"},
			},
			{
				Func: getWafRuleGroupActivatedRules,
				Tags: map[string]string{"service": "sso", "action": "ListActivatedRulesInRuleGroup"},
			},
			{
				Func: listTagsForWafRuleGroup,
				Tags: map[string]string{"service": "waf", "action": "ListTagsForResource"},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the rule group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the entity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafRuleGroupArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "rule_group_id",
				Description: "A unique identifier for the rule group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metric_name",
				Description: "A friendly name or description for the metrics for this RuleGroup.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafRuleGroup,
			},
			{
				Name:        "activated_rules",
				Description: "A list of activated rules associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafRuleGroupActivatedRules,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForWafRuleGroup,
				Transform:   transform.FromField("TagInfoForResource.TagList"),
			},

			// steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "RuleGroup.Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForWafRuleGroup,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(classicRuleGroupTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafRuleGroupArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listWafRuleGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := WAFClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.listAwsWAFRules", "get_client_error", err)
		return nil, err
	}

	// List all rule groups
	pagesLeft := true
	params := &waf.ListRuleGroupsInput{
		Limit: int32(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(params.Limit) {
			params.Limit = int32(*limit)
		}
	}

	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		response, err := svc.ListRuleGroups(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_waf_rule_group.listWafRuleGroups", "api_error", err)
			return nil, err
		}

		for _, ruleGroups := range response.RuleGroups {
			d.StreamListItem(ctx, ruleGroups)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
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

func getWafRuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		data := classicRuleGroupData(h.Item, ctx, d, h)
		id = data["rule_group_id"]
	} else {
		id = d.EqualsQuals["rule_group_id"].GetStringValue()
	}

	// Create session
	svc, err := WAFClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.getWafRuleGroup", "get_client_error", err)
		return nil, err
	}

	params := &waf.GetRuleGroupInput{
		RuleGroupId: &(id),
	}

	op, err := svc.GetRuleGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.getWafRuleGroup", "api_error", err)
		return nil, err
	}

	return op.RuleGroup, nil
}

func getWafRuleGroupActivatedRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := classicRuleGroupData(h.Item, ctx, d, h)
	id := data["rule_group_id"]

	// Create session
	svc, err := WAFClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.getWafRuleGroupActivatedRules", "get_client_error", err)
		return nil, err
	}

	params := &waf.ListActivatedRulesInRuleGroupInput{
		RuleGroupId: &(id),
	}

	op, err := svc.ListActivatedRulesInRuleGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.getWafRuleGroupActivatedRules", "api_error", err)
		return nil, err
	}

	return op, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go-v2/issues/3513
func listTagsForWafRuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := classicRuleGroupData(h.Item, ctx, d, h)

	// Create session
	svc, err := WAFClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.listTagsForWafRuleGroup", "get_client_error", err)
		return nil, err
	}

	// Build param with maximum limit set
	param := &waf.ListTagsForResourceInput{
		ResourceARN: aws.String(data["Arn"]),
		Limit:       int32(100),
	}

	ruleGroupTags, err := svc.ListTagsForResource(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.listTagsForWafRuleGroup", "api_error", err)
		return nil, err
	}
	return ruleGroupTags, nil
}

func getWafRuleGroupArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := classicRuleGroupData(h.Item, ctx, d, h)
	return data["Arn"], nil
}

//// TRANSFORM FUNCTIONS

func classicRuleGroupTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*waf.ListTagsForResourceOutput)

	if len(data.TagInfoForResource.TagList) < 1 {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if data.TagInfoForResource.TagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range data.TagInfoForResource.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func classicRuleGroupData(item interface{}, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) map[string]string {
	data := map[string]string{}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.classicRuleGroupData", "cache_error", err)
		return nil
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	switch item := item.(type) {
	case *types.RuleGroup:
		data["rule_group_id"] = *item.RuleGroupId
		// arn:aws:waf::account:rulegroup/name/ID
		data["Arn"] = "arn:" + commonColumnData.Partition + ":waf::" + commonColumnData.AccountId + ":rulegroup/" + *item.RuleGroupId
		data["Name"] = *item.Name

	case types.RuleGroupSummary:
		data["rule_group_id"] = *item.RuleGroupId
		// arn:aws:waf::account:rulegroup/name/ID
		data["Arn"] = "arn:" + commonColumnData.Partition + ":waf::" + commonColumnData.AccountId + ":rulegroup/" + *item.RuleGroupId
		data["Name"] = *item.Name
	}

	return data
}
