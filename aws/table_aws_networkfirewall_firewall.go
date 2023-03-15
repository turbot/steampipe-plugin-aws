package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/networkfirewall"
	"github.com/aws/aws-sdk-go-v2/service/networkfirewall/types"

	networkfirewallv1 "github.com/aws/aws-sdk-go/service/networkfirewall"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsNetworkFirewallFirewall(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_networkfirewall_firewall",
		Description: "AWS Network Firewall Firewall",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"arn", "name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidRequestException", "ValidationException"}),
			},
			Hydrate: getNetworkFirewallFirewall,
		},
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"vpc_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidRequestException", "ValidationException"}),
			},
			Hydrate: listNetworkFirewallFirewalls,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(networkfirewallv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The last time that the firewall policy was changed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.FirewallArn"),
			},
			{
				Name:        "name",
				Description: "The descriptive name of the firewall. You can't change the name of a firewall after you create it.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.FirewallName"),
			},
			{
				Name:        "vpc_id",
				Description: "The unique identifier of the VPC where the firewall is in use.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.VpcId"),
			},

			{
				Name:        "delete_protection",
				Description: "A flag indicating whether it is possible to delete the firewall.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.DeleteProtection"),
			},
			{
				Name:        "description",
				Description: "A description of the firewall.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.Description"),
			},
			{
				Name:        "id",
				Description: "The unique identifier for the firewall.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.FirewallId"),
			},
			{
				Name:        "policy_arn",
				Description: "The public subnets that Network Firewall is using for the firewall. Each subnet must belong to a different Availability Zone.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.FirewallPolicyArn"),
			},
			{
				Name:        "policy_change_protection",
				Description: "A setting indicating whether the firewall is protected against a change to the firewall policy association.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.FirewallPolicyChangeProtection"),
			},
			{
				Name:        "subnet_change_protection",
				Description: "A setting indicating whether the firewall is protected against changes to the subnet associations.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.SubnetChangeProtection"),
			},

			{
				Name:        "encryption_configuration",
				Description: "A complex type that contains the Amazon Web Services KMS encryption configuration settings for your firewall.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.EncryptionConfiguration"),
			},
			{
				Name:        "firewall_status",
				Description: "Detailed information about the current status of a Firewall.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallFirewall,
			},
			{
				Name:        "subnet_mappings",
				Description: "The public subnets that Network Firewall is using for the firewall.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.SubnetMappings"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the firewall",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.FirewallName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.Tags").Transform(networkFirewallTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallFirewall,
				Transform:   transform.FromField("Firewall.FirewallArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listNetworkFirewallFirewalls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// vpcId := d.EqualsQuals["vpc_id"].GetStringValue()

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

	input := &networkfirewall.ListFirewallsInput{
		MaxResults: &maxLimit,
	}

	if d.Quals["vpc_id"] != nil {
		for _, q := range d.Quals["vpc_id"].Quals {
			value := q.Value.GetStringValue()
			if q.Operator == "=" {
				input.VpcIds = append(input.VpcIds, value)
			}
		}
	}

	// Create session
	svc, err := NetworkFirewallClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_networkfirewall_firewall.listNetworkFirewallFirewalls", "connection_error", err)
		return nil, err
	}

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

	paginator := networkfirewall.NewListFirewallsPaginator(svc, input, func(o *networkfirewall.ListFirewallsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_networkfirewall_firewall.listNetworkFirewallFirewalls", "api_error", err)
			return nil, err
		}

		for _, firewall := range output.Firewalls {
			d.StreamListItem(ctx, firewall)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getNetworkFirewallFirewall(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var name, arn string
	if h.Item != nil {
		name = *h.Item.(types.FirewallMetadata).FirewallName
		arn = *h.Item.(types.FirewallMetadata).FirewallArn
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
		arn = d.EqualsQuals["arn"].GetStringValue()
	}

	// Build the params
	// Can pass in ARN, name, or both
	params := &networkfirewall.DescribeFirewallInput{}
	if arn != "" {
		params.FirewallArn = aws.String(arn)
	} else if name != "" {
		params.FirewallName = aws.String(name)
	} else {
		return nil, nil
	}

	// Create session
	svc, err := NetworkFirewallClient(ctx, d)
	if err != nil {
		logger.Error("aws_networkfirewall_firewall.getNetworkFirewallFirewall", "connection_error", err)
		return nil, err
	}

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

	// Get call
	data, err := svc.DescribeFirewall(ctx, params)
	if err != nil {
		logger.Error("aws_networkfirewall_firewall.getNetworkFirewallFirewall", "api_error", err)
		return nil, err
	}

	return data, nil
}

//// TRANSFORM FUNCTIONS

func networkFirewallTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	firewallTags := d.Value.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if firewallTags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range firewallTags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
