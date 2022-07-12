package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWafRuleGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_waf_rule_group",
		Description: "AWS WAF Rule Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"rule_group_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NonexistentItemException", "WAFNonexistentItemException"}),
			},
			Hydrate: getAwsWafRuleGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWafRuleGroups,
		},
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
				Hydrate:     getAwsWafRuleGroup,
			},
			{
				Name:        "activated_rules",
				Description: "A list of activated rules associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafRuleGroupActivatedRules,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAwsWafRuleGroup,
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
				Hydrate:     listTagsForAwsWafRuleGroup,
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

func listAwsWafRuleGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := WAFService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.listAwsWafRuleGroups", "service_creation_error", err)
		return nil, err
	}

	// List all rule groups
	pagesLeft := true
	params := &waf.ListRuleGroupsInput{
		Limit: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.Limit {
			if *limit < 1 {
				params.Limit = aws.Int64(1)
			} else {
				params.Limit = limit
			}
		}
	}

	for pagesLeft {
		response, err := svc.ListRuleGroups(params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_waf_rule_group.listAwsWafRuleGroups", "api_error", err)
			return nil, err
		}

		for _, ruleGroups := range response.RuleGroups {
			d.StreamListItem(ctx, ruleGroups)

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

func getAwsWafRuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		data := classicRuleGroupData(h.Item, ctx, d, h)
		id = data["rule_group_id"]
	} else {
		id = d.KeyColumnQuals["rule_group_id"].GetStringValue()
	}

	// Create Session
	svc, err := WAFService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.getAwsWafRuleGroup", "service_creation_error", err)
		return nil, err
	}

	params := &waf.GetRuleGroupInput{
		RuleGroupId: aws.String(id),
	}

	op, err := svc.GetRuleGroup(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.getAwsWafRuleGroup", "api_error", err)
		return nil, err
	}

	return op.RuleGroup, nil
}

func getAwsWafRuleGroupActivatedRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		data := classicRuleGroupData(h.Item, ctx, d, h)
		id = data["rule_group_id"]
	} else {
		id = d.KeyColumnQuals["rule_group_id"].GetStringValue()
	}

	// Create Session
	svc, err := WAFService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.getAwsWafRuleGroupRules", "service_creation_error", err)
		return nil, err
	}

	params := &waf.ListActivatedRulesInRuleGroupInput{
		RuleGroupId: aws.String(id),
	}

	op, err := svc.ListActivatedRulesInRuleGroup(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.getAwsWafRuleGroupRules", "api_error", err)
		return nil, err
	}

	return op, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go/issues/3513
func listTagsForAwsWafRuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := classicRuleGroupData(h.Item, ctx, d, h)

	// Create session
	svc, err := WAFService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.listTagsForAwsWafRuleGroup", "service_creation_error", err)
		return nil, err
	}

	// Build param with maximum limit set
	param := &waf.ListTagsForResourceInput{
		ResourceARN: aws.String(data["Arn"]),
		Limit:       aws.Int64(100),
	}

	ruleGroupTags, err := svc.ListTagsForResource(param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.listTagsForAwsWafRuleGroup", "api_error", err)
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

func classicRuleGroupData(item interface{}, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) map[string]string {
	data := map[string]string{}

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_rule_group.classicRuleGroupData", "cache_error", err)
		return nil
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	switch item := item.(type) {
	case *waf.RuleGroup:
		data["rule_group_id"] = *item.RuleGroupId

		// arn:aws:waf::account:rulegroup/name/ID
		data["Arn"] = "arn:aws:waf::" + commonColumnData.AccountId + ":" + "rulegroup/" + *item.RuleGroupId
		data["Name"] = *item.Name
	case *waf.RuleGroupSummary:
		data["rule_group_id"] = *item.RuleGroupId

		// arn:aws:waf::account:rulegroup/name/ID
		data["Arn"] = "arn:aws:waf::" + commonColumnData.AccountId + ":" + "rulegroup/" + *item.RuleGroupId
		data["Name"] = *item.Name
	}
	return data
}
