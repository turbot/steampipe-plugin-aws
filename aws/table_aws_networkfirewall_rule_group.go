package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/networkfirewall"
	"github.com/aws/aws-sdk-go-v2/service/networkfirewall/types"

	networkfirewallEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsNetworkFirewallRuleGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_networkfirewall_rule_group",
		Description: "AWS Network Firewall Rule Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"arn", "rule_group_name"}),
			Hydrate:    getNetworkFirewallRuleGroup,
			Tags:       map[string]string{"service": "network-firewall", "action": "DescribeRuleGroup"},
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkFirewallRuleGroups,
			Tags:    map[string]string{"service": "network-firewall", "action": "ListRuleGroups"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getNetworkFirewallRuleGroup,
				Tags: map[string]string{"service": "network-firewall", "action": "DescribeRuleGroup"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(networkfirewallEndpoint.NETWORK_FIREWALLServiceID),
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
				Name:        "sns_topic",
				Description: "The Amazon resource name (ARN) of the Amazon Simple Notification Service SNS topic that's used to record changes to the managed rule group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroupResponse.SnsTopic"),
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
				Name:        "analysis_results",
				Description: "The list of analysis results for AnalyzeRuleGroup.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroupResponse.AnalysisResults"),
			},
			{
				Name:        "encryption_configuration",
				Description: "A complex type that contains the Amazon Web Services KMS encryption configuration settings for your rule group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallRuleGroup,
				Transform:   transform.FromField("RuleGroupResponse.EncryptionConfiguration"),
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
	// Create session
	svc, err := NetworkFirewallClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_networkfirewall_rule_group.listNetworkFirewallRuleGroups", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &networkfirewall.ListRuleGroupsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := networkfirewall.NewListRuleGroupsPaginator(svc, input, func(o *networkfirewall.ListRuleGroupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_networkfirewall_rule_group.listNetworkFirewallRuleGroups", "api_error", err)
			return nil, err
		}

		for _, items := range output.RuleGroups {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getNetworkFirewallRuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name, arn string
	if h.Item != nil {
		name = *h.Item.(types.RuleGroupMetadata).Name
		arn = *h.Item.(types.RuleGroupMetadata).Arn
	} else {
		name = d.EqualsQuals["rule_group_name"].GetStringValue()
		arn = d.EqualsQuals["arn"].GetStringValue()
	}
	// Create session
	svc, err := NetworkFirewallClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_networkfirewall_rule_group.getNetworkFirewallRuleGroup", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
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
	data, err := svc.DescribeRuleGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_networkfirewall_rule_group.getNetworkFirewallRuleGroup", "api_error", err)
		return nil, err
	}
	return data, nil
}

//// TRANSFORM FUNCTIONS

func networkFirewallRuleGroupTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
