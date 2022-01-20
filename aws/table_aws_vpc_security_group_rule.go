package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcSecurityGroupRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_security_group_rule",
		Description: "AWS VPC Security Group Rule",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("security_group_rule_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidGroupId.Malformed", "InvalidGroupId.NotFound", "InvalidGroup.NotFound"}),
			Hydrate:           getSecurityGroupRule,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityGroupRules,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "security_group_rule_id",
				Description: "The ID of the security group rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "group_name",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Group.GroupName"),
			},
			{
				Name:        "group_id",
				Description: "The ID of the security group to which rule belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Group.GroupId"),
			},
			{
				Name:        "type",
				Description: "Type of the rule ( ingress | egress).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC for the security group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Group.VpcId"),
			},
			{
				Name:        "owner_id",
				Description: "The AWS account ID of the owner of the security group to which rule belongs.",
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
				Description: "The IPv4 CIDR range. It can be either a CIDR range or a source security group, not both. A single IPv4 address is denoted by /32 prefix length.",
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
				Description: "The ID of the security group that references this user ID group pair.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.GroupId"),
			},
			{
				Name:        "referenced_group_id",
				Description: "The ID of the security group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.GroupId"),
			},
			{
				Name:        "pair_group_name",
				Description: "The name of the security group that references this user ID group pair.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserIDGroupPair.GroupName"),
			},
			{
				Name:        "pair_peering_status",
				Description: "The status of a VPC peering connection, if applicable.",
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
				Description: "The ID of an AWS account. For a referenced security group in another VPC, the account ID of the referenced security group is returned in the response. If the referenced security group is deleted, this value is not returned.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.UserId"),
			},
			{
				Name:        "referenced_user_id",
				Description: "The Amazon Web Services account ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.UserId"),
			},
			{
				Name:        "pair_vpc_id",
				Description: "The ID of the VPC for the referenced security group, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.VpcId"),
			},
			{
				Name:        "referenced_vpc_id",
				Description: "The ID of the VPC for the referenced security group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.VpcId"),
			},
			{
				Name:        "pair_vpc_peering_connection_id",
				Description: "The ID of the VPC peering connection, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.VpcPeeringConnectionId"),
			},
			{
				Name:        "referenced_vpc_peering_connection_id",
				Description: "The ID of the VPC peering connection.",
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

	// List call
	err = svc.DescribeSecurityGroupRulesPages(
		&ec2.DescribeSecurityGroupRulesInput{},
		func(page *ec2.DescribeSecurityGroupRulesOutput, isLast bool) bool {
			for _, securityGroupRule := range page.SecurityGroupRules {
				d.StreamListItem(ctx, securityGroupRule)
			}
			return !isLast
		},
	)

	return nil, err
}

func getSecurityGroupRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSecurityGroupRule")

	region := d.KeyColumnQualString(matrixKeyRegion)
	ruleID := d.KeyColumnQuals["security_group_rule_id"].GetStringValue()

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
		plugin.Logger(ctx).Debug("DescribeSecurityGroupRules__", "ERROR", err)
		return nil, err
	}

	if op.SecurityGroupRules != nil && len(op.SecurityGroupRules) > 0 {
		return op.SecurityGroupRules[0], nil
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityGroupRuleTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	sgRule := h.Item.(*ec2.SecurityGroupRule)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	// To create uninque aka
	hashCode := "_" + *sgRule.IpProtocol
	if *sgRule.IsEgress {
		hashCode = "ingress" + hashCode
	} else {
		hashCode = "egress" + hashCode
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

	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + *sgRule.GroupOwnerId + ":security-group-rule/" + *sgRule.SecurityGroupRuleId
	// generate aka for the rule
	akas := []string{arn + ":" + hashCode}

	title := *sgRule.SecurityGroupRuleId + "_" + hashCode

	turbotData := map[string]interface{}{
		"Arn":   arn,
		"Akas":  akas,
		"Title": title,
	}

	return turbotData, nil
}

// custom struct for security group rule
type vpcSecurityGroupRulesRowData struct {
	Group           *ec2.SecurityGroup
	Permission      *ec2.IpPermission
	IPRange         *ec2.IpRange
	Ipv6Range       *ec2.Ipv6Range
	UserIDGroupPair *ec2.UserIdGroupPair
	PrefixListId    *ec2.PrefixListId
	Type            string
}
