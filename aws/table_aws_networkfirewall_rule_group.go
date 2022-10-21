package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/networkfirewall"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
)

//// TABLE DEFINITION

func tableAwsNetworkFirewallRuleGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_networkfirewall_rule_group",
		Description: "AWS Network Firewall Rule Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"arn", "rule_group_name"}),
			Hydrate:    getNetworkFirewallRuleGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkFirewallRuleGroups,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "rule_group_name",
				Description: "The descriptive name of the rule group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the rule group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Arn"),
			},
			{
				Name:        "capacity",
				Description: "The maximum operating resources that this rule group can use. Rule group capacity is fixed at creation. When you update a rule group, you are limited to this capacity. When you reference a rule group from a firewall policy, Network Firewall reserves this capacity for the rule group.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroupResponse.Capacity"),
			},
			{
				Name:        "consumed_capacity",
				Description: "The number of capacity units currently consumed by the rule group rules.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroupResponse.ConsumedCapacity"),
			},
			{
				Name:        "description",
				Description: "A description of the rule group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroupResponse.Description"),
			},
			{
				Name:        "number_of_associations",
				Description: "The number of firewall policies that use this rule group.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroupResponse.NumberOfAssociations"),
			},
			{
				Name:        "rule_group_id",
				Description: "The unique identifier for the rule group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroupResponse.RuleGroupId"),
			},
			{
				Name:        "rule_group_status",
				Description: "Detailed information about the current status of a rule group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroupResponse.RuleGroupStatus"),
			},
			{
				Name:        "rule_variables",
				Description: "Settings that are available for use in the rules in the rule group. You can only use these for stateful rule groups.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroup.RuleVariables"),
			},
			{
				Name:        "rules_source",
				Description: "The stateful rules or stateless rules for the rule group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroup.RulesSource"),
			},
			{
				Name:        "stateful_rule_options",
				Description: "Additional options governing how Network Firewall handles the rule group. You can only use these for stateful rule groups.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroup.StatefulRuleOptions"),
			},
			{
				Name:        "type",
				Description: "Indicates whether the rule group is stateless or stateful. If the rule group is stateless, it contains stateless rules. If it is stateful, it contains stateful rules.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroupResponse.Type"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroupResponse.Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "RuleGroupName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("Tags").Transform(networkFirewallRuleGroupTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listNetworkFirewallRuleGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listNetworkFirewallRuleGroups")

	// Create session
	svc, err := NetworkFirewallService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &networkfirewall.ListRuleGroupsInput{
		MaxResults: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.ListRuleGroupsPages(
		input,
		func(page *networkfirewall.ListRuleGroupsOutput, isLast bool) bool {
			for _, rule_group := range page.RuleGroups {
				d.StreamListItem(ctx, rule_group)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getNetworkFirewallRuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getNetworkFirewallRuleGroup")

	var name, arn string
	if h.Item != nil {
		name = *h.Item.(*networkfirewall.RuleGroupMetadata).Name
		arn = *h.Item.(*networkfirewall.RuleGroupMetadata).Arn
	} else {
		name = d.KeyColumnQuals["rule_group_name"].GetStringValue()
		arn = d.KeyColumnQuals["arn"].GetStringValue()
	}
	// Create session
	svc, err := NetworkFirewallService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	// Can pass in ARN, name, or both
	params := &networkfirewall.DescribeRuleGroupInput{}
	if name != "" {
		params.RuleGroupName = aws.String(name)
	}
	if arn != "" {
		params.RuleGroupArn = aws.String(arn)
	}

	// Get call
	data, err := svc.DescribeRuleGroup(params)
	if err != nil {
		logger.Debug("getNetworkFirewallRuleGroup", "ERROR", err)
		return nil, err
	}
	return data, nil
}

//// TRANSFORM FUNCTIONS

func networkFirewallRuleGroupTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("networkFirewallRuleGroupTurbotTags")
	ruleGroup := d.HydrateItem.(*networkfirewall.DescribeRuleGroupOutput)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if ruleGroup.RuleGroupResponse.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range ruleGroup.RuleGroupResponse.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
