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

func tableAwsNetworkFirewallPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_networkfirewall_firewall_policy",
		Description: "AWS Network Firewall Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"arn", "name"}),
			Hydrate:    getNetworkFirewallPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkFirewallPolicies,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The descriptive name of the rule group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "FirewallPolicyResponse.FirewallPolicyName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the rule group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Arn", "FirewallPolicyResponse.FirewallPolicyArn"),
			},
			{
				Name:        "firewall_policy_id",
				Description: "The unique identifier for the firewall policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallPolicy,
				Transform:   transform.FromField("FirewallPolicyResponse.FirewallPolicyId"),
			},
			{
				Name:        "consumed_stateful_rule_capacity",
				Description: "The number of capacity units currently consumed by the policy's stateful rules.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getNetworkFirewallPolicy,
				Transform:   transform.FromField("FirewallPolicyResponse.ConsumedStatefulRuleCapacity"),
			},
			{
				Name:        "consumed_stateless_rule_capacity",
				Description: "The number of capacity units currently consumed by the policy's stateless rules.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getNetworkFirewallPolicy,
				Transform:   transform.FromField("FirewallPolicyResponse.ConsumedStatelessRuleCapacity"),
			},
			{
				Name:        "description",
				Description: "A description of the firewall policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallPolicy,
				Transform:   transform.FromField("FirewallPolicyResponse.Description"),
			},
			{
				Name:        "firewall_policy_status",
				Description: "The current status of the firewall policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallPolicy,
				Transform:   transform.FromField("FirewallPolicyResponse.FirewallPolicyStatus"),
			},
			{
				Name:        "last_modified_time",
				Description: "The last time that the firewall policy was changed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getNetworkFirewallPolicy,
				Transform:   transform.FromField("FirewallPolicyResponse.LastModifiedTime"),
			},
			{
				Name:        "number_of_associations",
				Description: "The number of firewall policies that use this rule group.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getNetworkFirewallPolicy,
				Transform:   transform.FromField("FirewallPolicyResponse.NumberOfAssociations"),
			},

			{
				Name:        "encryption_configuration",
				Description: "A complex type that contains the Amazon Web Services KMS encryption configuration settings for your firewall policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallPolicy,
				Transform:   transform.FromField("FirewallPolicyResponse.EncryptionConfiguration"),
			},
			{
				Name:        "firewall_policy",
				Description: "The policy for the specified firewall policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallPolicy,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallPolicy,
				Transform:   transform.FromField("FirewallPolicyResponse.Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "FirewallPolicyResponse.FirewallPolicyName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallPolicy,
				Transform:   transform.FromField("Tags").Transform(networkFirewallPolicyTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn", "FirewallPolicyResponse.FirewallPolicyArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listNetworkFirewallPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := NetworkFirewallService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_networkfirewall_firewall_policy.listNetworkFirewallPolicies", "service_creation_error", err)
		return nil, err
	}

	input := &networkfirewall.ListFirewallPoliciesInput{
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
	err = svc.ListFirewallPoliciesPages(
		input,
		func(page *networkfirewall.ListFirewallPoliciesOutput, isLast bool) bool {
			for _, policy := range page.FirewallPolicies {
				d.StreamListItem(ctx, policy)

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

func getNetworkFirewallPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var name, arn string
	if h.Item != nil {
		name = *h.Item.(*networkfirewall.FirewallPolicyMetadata).Name
		arn = *h.Item.(*networkfirewall.FirewallPolicyMetadata).Arn
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
		arn = d.KeyColumnQuals["arn"].GetStringValue()
	}

	// Build the params
	// Can pass in ARN, name, or both
	params := &networkfirewall.DescribeFirewallPolicyInput{}
	if name != "" {
		params.FirewallPolicyName = aws.String(name)
	}
	if arn != "" {
		params.FirewallPolicyArn = aws.String(arn)
	}

	// Create session
	svc, err := NetworkFirewallService(ctx, d)
	if err != nil {
		logger.Error("aws_networkfirewall_firewall_policy.getNetworkFirewallPolicy", "service_creation_error", err)
		return nil, err
	}

	// Get call
	data, err := svc.DescribeFirewallPolicy(params)
	if err != nil {
		logger.Error("aws_networkfirewall_firewall_policy.getNetworkFirewallPolicy", "api_error", err)
		return nil, err
	}
	return data, nil
}

//// TRANSFORM FUNCTIONS

func networkFirewallPolicyTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	policy := d.HydrateItem.(*networkfirewall.DescribeFirewallPolicyOutput)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if policy.FirewallPolicyResponse.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range policy.FirewallPolicyResponse.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
