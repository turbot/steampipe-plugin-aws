package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsVpcSecurityGroupRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_security_group_rule",
		Description: "AWS VPC Security Group Rule",
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidGroup.NotFound", "InvalidSecurityGroupRuleId.Malformed", "InvalidSecurityGroupRuleId.NotFound"}),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("security_group_rule_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidGroup.NotFound", "InvalidSecurityGroupRuleId.Malformed", "InvalidSecurityGroupRuleId.NotFound"}),
			},
			Hydrate: getSecurityGroupRule,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeSecurityGroupRules"},
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityGroupRules,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeSecurityGroupRules"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "group_id",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getReferencedSecurityGroupDetails,
				Tags: map[string]string{"service": "ec2", "action": "DescribeSecurityGroups"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "security_group_rule_id",
				Description: "The ID of the security group rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "group_id",
				Description: "The ID of the security group to which rule belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_egress",
				Description: "Indicates whether the security group rule is an outbound rule.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "type",
				Description: "Type of the rule ( ingress | egress).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IsEgress").Transform(setRuleType),
			},
			{
				Name:        "group_owner_id",
				Description: "The ID of the Amazon Web Services account that owns the security group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The security group rule description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_protocol",
				Description: "The IP protocol name (tcp, udp, icmp, icmpv6) or number [see Protocol Numbers ](http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml). Use -1 to specify all protocols. When authorizing security group rules, specifying -1 or a protocol number other than tcp, udp, icmp, or icmpv6 allows traffic on all ports, regardless of any port range specified. For tcp, udp, and icmp, a port range is specified. For icmpv6, the port range is optional. If port range is omitted, traffic for all types and codes is allowed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "from_port",
				Description: "The start of port range for the TCP and UDP protocols, or an ICMP/ICMPv6 type number. A value of -1 indicates all ICMP/ICMPv6 types.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "to_port",
				Description: "The end of port range for the TCP and UDP protocols, or an ICMP/ICMPv6 code. A value of -1 indicates all ICMP/ICMPv6 codes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "cidr_ipv4",
				Description: "The IPv4 CIDR range.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "cidr_ipv6",
				Description: "The IPv6 CIDR range.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "referenced_group_id",
				Description: "The ID of the referenced security group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.GroupId"),
			},
			{
				Name:        "referenced_peering_status",
				Description: "The status of a VPC peering connection, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.PeeringStatus"),
			},
			{
				Name:        "referenced_user_id",
				Description: "The ID of an AWS account. For a referenced security group in another VPC, the account ID of the referenced security group is returned in the response. If the referenced security group is deleted, this value is not returned.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.UserId"),
			},
			{
				Name:        "referenced_vpc_id",
				Description: "The ID of the VPC for the referenced security group, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.VpcId"),
			},
			{
				Name:        "referenced_vpc_peering_connection_id",
				Description: "The ID of the VPC peering connection, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.VpcPeeringConnectionId"),
			},
			{
				Name:        "prefix_list_id",
				Description: "The ID of the prefix list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSecurityGroupRuleTurbotData,
				Transform:   transform.FromField("Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSecurityGroupRuleTurbotData,
				Transform:   transform.FromField("Akas"),
			},
		}),
	}
}

//// LIST FUNCTION

func listSecurityGroupRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group_rule.listSecurityGroupRules", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			// As per API Docs MaxResults value can be between 5 and 1000
			if limit < 5 {
				maxLimit = int32(5)
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeSecurityGroupRulesInput{}

	groupId := d.EqualsQuals["group_id"].GetStringValue()
	if groupId != "" {
		input.Filters = []types.Filter{
			{
				Name:   aws.String("group-id"),
				Values: []string{groupId},
			},
		}
	}

	paginator := ec2.NewDescribeSecurityGroupRulesPaginator(svc, input, func(o *ec2.DescribeSecurityGroupRulesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_security_group_rule.listSecurityGroupRules", "api_error", err)
			return nil, err
		}

		for _, securityGroupRule := range output.SecurityGroupRules {
			d.StreamListItem(ctx, securityGroupRule)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityGroupRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	ruleID := d.EqualsQuals["security_group_rule_id"].GetStringValue()

	// check if rule id is empty
	if ruleID == "" {
		return nil, nil
	}

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group_rule.getSecurityGroupRule", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeSecurityGroupRulesInput{
		SecurityGroupRuleIds: []string{ruleID},
	}

	// Get call
	op, err := svc.DescribeSecurityGroupRules(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group_rule.getSecurityGroupRule", "api_error", err)
		return nil, err
	}

	if len(op.SecurityGroupRules) > 0 {
		return op.SecurityGroupRules[0], nil
	}

	return nil, nil
}

func getSecurityGroupDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	sgRule := h.Item.(types.SecurityGroupRule)

	// Build the params
	params := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{*sgRule.GroupId},
	}

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group_rule.getSecurityGroupDetails", "connection_error", err)
		return nil, err
	}

	op, err := svc.DescribeSecurityGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group_rule.getSecurityGroupDetails", "api_error", err)
		return nil, err
	}

	if len(op.SecurityGroups) > 0 {
		return op.SecurityGroups[0], nil
	}

	return nil, nil
}

func getReferencedSecurityGroupDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	sgRule := h.Item.(types.SecurityGroupRule)
	if sgRule.ReferencedGroupInfo == nil {
		return nil, nil
	}

	// Build the params
	params := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{*sgRule.ReferencedGroupInfo.GroupId},
	}

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group_rule.getReferencedSecurityGroupDetails", "connection_error", err)
		return nil, err
	}

	op, err := svc.DescribeSecurityGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group_rule.getReferencedSecurityGroupDetails", "api_error", err)
		return nil, err
	}

	if len(op.SecurityGroups) > 0 {
		return op.SecurityGroups[0], nil
	}

	return nil, nil
}

func getSecurityGroupRuleTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	sgRule := h.Item.(types.SecurityGroupRule)
	region := d.EqualsQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_security_group_rule.getSecurityGroupRuleTurbotData", "common_data_error", err)
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	// Create a unique AKA

	hashCode := "_" + *sgRule.IpProtocol
	if *sgRule.IsEgress {
		hashCode = "egress" + hashCode
	} else {
		hashCode = "ingress" + hashCode
	}

	if sgRule.FromPort != nil {
		hashCode = fmt.Sprintf("%s_%d_%d", hashCode, sgRule.FromPort, sgRule.ToPort)
	}

	if sgRule.CidrIpv4 != nil {
		hashCode = hashCode + "_" + *sgRule.CidrIpv4
	} else if sgRule.CidrIpv6 != nil {
		hashCode = hashCode + "_" + *sgRule.CidrIpv6
	} else if sgRule.ReferencedGroupInfo != nil && sgRule.ReferencedGroupInfo.GroupId != nil && *sgRule.ReferencedGroupInfo.GroupId == *sgRule.GroupId {
		hashCode = hashCode + "_" + *sgRule.ReferencedGroupInfo.GroupId
	} else if sgRule.PrefixListId != nil {
		hashCode = hashCode + "_" + *sgRule.PrefixListId
	}

	// generate aka for the rule
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + *sgRule.GroupOwnerId + ":security-group-rule/" + *sgRule.SecurityGroupRuleId + ":" + hashCode}

	title := *sgRule.SecurityGroupRuleId + "_" + hashCode

	turbotData := map[string]interface{}{
		"Akas":  akas,
		"Title": title,
	}

	return turbotData, nil
}

func setRuleType(_ context.Context, d *transform.TransformData) (interface{}, error) {
	value := d.Value.(*bool)
	if !*value {
		return "ingress", nil
	}
	return "egress", nil
}
