package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsVpcSecurityGroupRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_security_group_rule",
		Description: "AWS VPC Security Group Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("security_group_rule_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidSecurityGroupRuleId.Malformed", "InvalidSecurityGroupRuleId.NotFound"}),
			},
			Hydrate: getSecurityGroupRule,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityGroupRules,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "group_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "security_group_rule_id",
				Description: "The ID of the security group rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "group_name",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release. The name of the security group to which rule belongs.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSecurityGroupDetails,
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
				Name:        "vpc_id",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release. The ID of the VPC for the security group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSecurityGroupDetails,
			},
			{
				Name:        "owner_id",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release, use group_owner_id instead. The AWS account ID of the owner of the security group to which rule belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GroupOwnerId"),
			},
			{
				Name:        "group_owner_id",
				Description: "The AWS account ID of the owner of the security group to which rule belongs.",
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
				Name:        "cidr_ip",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release, use cidr_ipv4 instead. The IPv4 CIDR range. It can be either a CIDR range or a source security group, not both. A single IPv4 address is denoted by /32 prefix length.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("CidrIpv4"),
			},
			{
				Name:        "cidr_ipv4",
				Description: "The IPv4 CIDR range. It can be either a CIDR range or a source security group, not both. A single IPv4 address is denoted by /32 prefix length.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "cidr_ipv6",
				Description: "The IPv6 CIDR range. It can be either a CIDR range or a source security group, not both. A single IPv6 address is denoted by /128 prefix length.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "pair_group_id",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release, use referenced_group_id instead. The ID of the referenced security group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.GroupId"),
			},
			{
				Name:        "referenced_group_id",
				Description: "The ID of the referenced security group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.GroupId"),
			},
			{
				Name:        "pair_group_name",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release. The name of the referenced security group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getReferencedSecurityGroupDetails,
				Transform:   transform.FromField("GroupName"),
			},
			{
				Name:        "pair_peering_status",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release, use referenced_peering_status instead. Please use the referenced_peering_status column instead. The status of a VPC peering connection, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.PeeringStatus"),
			},
			{
				Name:        "referenced_peering_status",
				Description: "The status of a VPC peering connection, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.PeeringStatus"),
			},
			{
				Name:        "pair_user_id",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release, use referenced_user_id instead. The ID of an AWS account. For a referenced security group in another VPC, the account ID of the referenced security group is returned in the response. If the referenced security group is deleted, this value is not returned.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.UserId"),
			},
			{
				Name:        "referenced_user_id",
				Description: "The ID of an AWS account. For a referenced security group in another VPC, the account ID of the referenced security group is returned in the response. If the referenced security group is deleted, this value is not returned.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.UserId"),
			},
			{
				Name:        "pair_vpc_id",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release, use referenced_vpc_id instead. The ID of the VPC for the referenced security group, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.VpcId"),
			},
			{
				Name:        "referenced_vpc_id",
				Description: "The ID of the VPC for the referenced security group, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.VpcId"),
			},
			{
				Name:        "pair_vpc_peering_connection_id",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release, use referenced_vpc_peering_connection_id instead. The ID of the VPC peering connection, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.VpcPeeringConnectionId"),
			},
			{
				Name:        "referenced_vpc_peering_connection_id",
				Description: "The ID of the VPC peering connection, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.VpcPeeringConnectionId"),
			},
			{
				Name:        "prefix_list_id",
				Description: "The ID of the referenced prefix list.",
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listSecurityGroupRules", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Additonal Filter
	// As per API Docs MaxResults value can be between 5 and 1000
	input := &ec2.DescribeSecurityGroupRulesInput{
		MaxResults: aws.Int64(1000),
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 5 {
				input.MaxResults = types.Int64(5)
			} else {
				input.MaxResults = limit
			}
		}
	}

	groupId := d.KeyColumnQuals["group_id"].GetStringValue()
	if groupId != "" {
		paramFilter := &ec2.Filter{
			Name:   aws.String("group-id"),
			Values: []*string{aws.String(groupId)},
		}
		input.Filters = []*ec2.Filter{paramFilter}
	}

	// List call
	err = svc.DescribeSecurityGroupRulesPages(
		input,
		func(page *ec2.DescribeSecurityGroupRulesOutput, isLast bool) bool {
			for _, securityGroupRule := range page.SecurityGroupRules {
				d.StreamListItem(ctx, securityGroupRule)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listSecurityGroupRules", "list", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityGroupRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSecurityGroupRule")

	region := d.KeyColumnQualString(matrixKeyRegion)
	ruleID := d.KeyColumnQuals["security_group_rule_id"].GetStringValue()

	// check if rule id is empty
	if ruleID == "" {
		return nil, nil
	}

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeSecurityGroupRulesInput{
		SecurityGroupRuleIds: []*string{aws.String(ruleID)},
	}

	// Get call
	op, err := svc.DescribeSecurityGroupRules(params)
	if err != nil {
		plugin.Logger(ctx).Error("getSecurityGroupRule", "get", err)
		return nil, err
	}

	if len(op.SecurityGroupRules) > 0 {
		return op.SecurityGroupRules[0], nil
	}

	return nil, nil
}

func getSecurityGroupDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSecurityGroupDetails")

	region := d.KeyColumnQualString(matrixKeyRegion)
	sgRule := h.Item.(*ec2.SecurityGroupRule)

	// Build the params
	params := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{aws.String(*sgRule.GroupId)},
	}

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	op, err := svc.DescribeSecurityGroups(params)
	if err != nil {
		// Unlikely, but handle any NotFound errors
		if strings.Contains(err.Error(), "InvalidGroup.NotFound") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("getSecurityGroupDetails", "ERROR", err)
		return nil, err
	}

	if len(op.SecurityGroups) > 0 {
		return op.SecurityGroups[0], nil
	}

	return nil, nil
}

func getReferencedSecurityGroupDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getReferencedSecurityGroupDetails")

	region := d.KeyColumnQualString(matrixKeyRegion)
	sgRule := h.Item.(*ec2.SecurityGroupRule)

	if sgRule.ReferencedGroupInfo == nil {
		return nil, nil
	}

	// Build the params
	params := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{aws.String(*sgRule.ReferencedGroupInfo.GroupId)},
	}

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	op, err := svc.DescribeSecurityGroups(params)
	if err != nil {
		// If the referenced security group is in another account, a NotFound error
		// will be returned
		if strings.Contains(err.Error(), "InvalidGroup.NotFound") {
			plugin.Logger(ctx).Error("getReferencedSecurityGroupDetails", "ERROR", err)
			return nil, nil
		}
		return nil, err
	}

	if len(op.SecurityGroups) > 0 {
		return op.SecurityGroups[0], nil
	}

	return nil, nil
}

func getSecurityGroupRuleTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	sgRule := h.Item.(*ec2.SecurityGroupRule)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
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
		hashCode = hashCode + "_" + types.IntToString(sgRule.FromPort) + "_" + types.IntToString(sgRule.ToPort)
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
