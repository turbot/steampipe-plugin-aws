package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/networkfirewall"
	"github.com/aws/aws-sdk-go-v2/service/networkfirewall/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsNetworkFirewallPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_networkfirewall_firewall_policy",
		Description: "AWS Network Firewall Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"arn", "name"}),
			Hydrate:    getNetworkFirewallPolicy,
			Tags:       map[string]string{"service": "network-firewall", "action": "DescribeFirewallPolicy"},
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkFirewallPolicies,
			Tags:    map[string]string{"service": "network-firewall", "action": "ListFirewallPolicies"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getNetworkFirewallPolicy,
				Tags: map[string]string{"service": "network-firewall", "action": "DescribeFirewallPolicy"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_NETWORK_FIREWALL_SERVICE_ID),
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
	svc, err := NetworkFirewallClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_networkfirewall_firewall_policy.listNetworkFirewallPolicies", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
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

	input := &networkfirewall.ListFirewallPoliciesInput{
		MaxResults: &maxLimit,
	}

	paginator := networkfirewall.NewListFirewallPoliciesPaginator(svc, input, func(o *networkfirewall.ListFirewallPoliciesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_networkfirewall_firewall_policy.listNetworkFirewallPolicies", "api_error", err)
			return nil, err
		}

		for _, policy := range output.FirewallPolicies {
			d.StreamListItem(ctx, policy)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getNetworkFirewallPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var name, arn string
	if h.Item != nil {
		name = *h.Item.(types.FirewallPolicyMetadata).Name
		arn = *h.Item.(types.FirewallPolicyMetadata).Arn
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
		arn = d.EqualsQuals["arn"].GetStringValue()
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
	svc, err := NetworkFirewallClient(ctx, d)
	if err != nil {
		logger.Error("aws_networkfirewall_firewall_policy.getNetworkFirewallPolicy", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Get call
	data, err := svc.DescribeFirewallPolicy(ctx, params)
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
