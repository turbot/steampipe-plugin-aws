package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafregional"
	"github.com/aws/aws-sdk-go-v2/service/wafregional/types"
	wafregionalv1 "github.com/aws/aws-sdk-go/service/wafregional"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWafRegionalRuleGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafregional_rule_group",
		Description: "AWS WAF Regional Rule Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"rule_group_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NonexistentItemException", "WAFNonexistentItemException"}),
			},
			Hydrate: getWafRegionalRuleGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listWafRegionalRuleGroups,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wafregionalv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the rule group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the entity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafRegionalRuleGroupArn,
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
				Hydrate:     getWafRegionalRuleGroup,
			},
			{
				Name:        "activated_rules",
				Description: "A list of activated rules associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafRegionalRuleGroupActivatedRules,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForWafRegionalRuleGroup,
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
				Hydrate:     listTagsForWafRegionalRuleGroup,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(wafRegionalRuleGroupTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafRegionalRuleGroupArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listWafRegionalRuleGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := WAFRegionalClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_rule_group.listWafRegionalRuleGroups", "client_error", err)
		return nil, err
	}

	// List all rule groups
	pagesLeft := true
	params := &wafregional.ListRuleGroupsInput{
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
		response, err := svc.ListRuleGroups(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wafregional_rule_group.listWafRegionalRuleGroups", "api_error", err)
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

func getWafRegionalRuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		data := wafRegionalRuleGroupData(h.Item, ctx, d, h)
		id = data["rule_group_id"]
	} else {
		id = d.EqualsQuals["rule_group_id"].GetStringValue()
		if id == "" {
			return nil, nil
		}
	}

	// Create session
	svc, err := WAFRegionalClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_rule_group.getWafRegionalRuleGroup", "client_error", err)
		return nil, err
	}

	params := &wafregional.GetRuleGroupInput{
		RuleGroupId: &(id),
	}

	op, err := svc.GetRuleGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_rule_group.getWafRegionalRuleGroup", "api_error", err)
		return nil, err
	}

	return op.RuleGroup, nil
}

func getWafRegionalRuleGroupActivatedRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := wafRegionalRuleGroupData(h.Item, ctx, d, h)
	id := data["rule_group_id"]

	// Create session
	svc, err := WAFRegionalClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_rule_group.getWafRegionalRuleGroupActivatedRules", "client_error", err)
		return nil, err
	}

	params := &wafregional.ListActivatedRulesInRuleGroupInput{
		RuleGroupId: &(id),
	}

	op, err := svc.ListActivatedRulesInRuleGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_rule_group.getWafRegionalRuleGroupActivatedRules", "api_error", err)
		return nil, err
	}

	return op, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go/issues/3513
func listTagsForWafRegionalRuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := wafRegionalRuleGroupData(h.Item, ctx, d, h)

	// Create session
	svc, err := WAFRegionalClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_rule_group.listTagsForWafRegionalRuleGroup", "client_error", err)
		return nil, err
	}

	// Build param with maximum limit set
	param := &wafregional.ListTagsForResourceInput{
		ResourceARN: aws.String(data["Arn"]),
		Limit:       int32(100),
	}

	ruleGroupTags, err := svc.ListTagsForResource(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_rule_group.listTagsForWafRegionalRuleGroup", "api_error", err)
		return nil, err
	}
	return ruleGroupTags, nil
}

func getWafRegionalRuleGroupArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := wafRegionalRuleGroupData(h.Item, ctx, d, h)
	return data["Arn"], nil
}

//// TRANSFORM FUNCTIONS

func wafRegionalRuleGroupTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*wafregional.ListTagsForResourceOutput)

	if data.TagInfoForResource.TagList == nil || len(data.TagInfoForResource.TagList) < 1 {
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

func wafRegionalRuleGroupData(item interface{}, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) map[string]string {
	data := map[string]string{}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_rule_group.wafRegionalRuleGroupData", "cache_error", err)
		return nil
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	region := d.EqualsQualString(matrixKeyRegion)
	switch item := item.(type) {
	case *types.RuleGroup:
		data["rule_group_id"] = *item.RuleGroupId
		data["Arn"] = "arn:" + commonColumnData.Partition + ":waf-regional:" + region + ":" + commonColumnData.AccountId + ":rulegroup/" + *item.RuleGroupId
		data["Name"] = *item.Name

	case types.RuleGroupSummary:
		data["rule_group_id"] = *item.RuleGroupId
		data["Arn"] = "arn:" + commonColumnData.Partition + ":waf-regional:" + region + ":" + commonColumnData.AccountId + ":rulegroup/" + *item.RuleGroupId
		data["Name"] = *item.Name
	}

	return data
}
