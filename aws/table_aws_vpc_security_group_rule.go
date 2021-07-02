package aws

import (
	"context"

	"github.com/turbot/go-kit/types"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcSecurityGroupRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_security_group_rule",
		Description: "AWS VPC Security Group Rule",
		// TODO -- get call returning a list of items
		List: &plugin.ListConfig{
			ParentHydrate: listVpcSecurityGroups,
			Hydrate:       listSecurityGroupRules,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "group_name",
				Description: "The name of the security group to which rule belongs.",
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
				Transform:   transform.FromField("Group.OwnerId"),
			},
			{
				Name:        "ip_protocol",
				Description: "The IP protocol name (tcp, udp, icmp, icmpv6) or number [see Protocol Numbers ](http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml). Use -1 to specify all protocols. When authorizing security group rules, specifying -1 or a protocol number other than tcp, udp, icmp, or icmpv6 allows traffic on all ports, regardless of any port range specified. For tcp, udp, and icmp, a port range is specified. For icmpv6, the port range is optional. If port range is omitted, traffic for all types and codes is allowed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Permission.IpProtocol"),
			},
			{
				Name:        "from_port",
				Description: "The start of port range for the TCP and UDP protocols, or an ICMP/ICMPv6 type number. A value of -1 indicates all ICMP/ICMPv6 types.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Permission.FromPort"),
			},
			{
				Name:        "to_port",
				Description: "The end of port range for the TCP and UDP protocols, or an ICMP/ICMPv6 code. A value of -1 indicates all ICMP/ICMPv6 codes.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Permission.ToPort"),
			},
			{
				Name:        "cidr_ip",
				Description: "The IPv4 CIDR range. It can be either a CIDR range or a source security group, not both. A single IPv4 address is denoted by /32 prefix length.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("IPRange.CidrIp"),
			},
			{
				Name:        "cidr_ipv6",
				Description: "The IPv6 CIDR range. It can be either a CIDR range or a source security group, not both. A single IPv6 address is denoted by /128 prefix length.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("Ipv6Range.CidrIpv6"),
			},
			{
				Name:        "pair_group_id",
				Description: "The ID of the security group that references this user ID group pair.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserIDGroupPair.GroupId"),
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
				Transform:   transform.FromField("UserIDGroupPair.PeeringStatus"),
			},
			{
				Name:        "pair_user_id",
				Description: "The ID of an AWS account. For a referenced security group in another VPC, the account ID of the referenced security group is returned in the response. If the referenced security group is deleted, this value is not returned.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserIDGroupPair.UserId"),
			},
			{
				Name:        "pair_vpc_id",
				Description: "The ID of the VPC for the referenced security group, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserIDGroupPair.VpcId"),
			},
			{
				Name:        "pair_vpc_peering_connection_id",
				Description: "The ID of the VPC peering connection, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserIDGroupPair.VpcPeeringConnectionId"),
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
	plugin.Logger(ctx).Trace("listSecurityGroupRules")
	securityGroup := h.Item.(*ec2.SecurityGroup)

	if securityGroup.IpPermissions != nil {
		for _, permission := range securityGroup.IpPermissions {
			for _, sgRule := range rowSourceFromIPPermission(securityGroup, permission, "ingress") {
				d.StreamLeafListItem(ctx, sgRule)
			}
		}
	}

	if securityGroup.IpPermissionsEgress != nil {
		for _, permission := range securityGroup.IpPermissionsEgress {
			for _, sgRule := range rowSourceFromIPPermission(securityGroup, permission, "egress") {
				d.StreamLeafListItem(ctx, sgRule)
			}
		}

	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityGroupRuleTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	sgRule := h.Item.(*vpcSecurityGroupRulesRowData)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	// To create uninque aka
	hashCode := sgRule.Type + "_" + *sgRule.Permission.IpProtocol
	if sgRule.Permission.FromPort != nil {
		hashCode = hashCode + "_" + types.IntToString(sgRule.Permission.FromPort) + "_" + types.IntToString(sgRule.Permission.ToPort)
	}

	if sgRule.IPRange != nil && sgRule.IPRange.CidrIp != nil {
		hashCode = hashCode + "_" + *sgRule.IPRange.CidrIp
	} else if sgRule.Ipv6Range != nil && sgRule.Ipv6Range.CidrIpv6 != nil {
		hashCode = hashCode + "_" + *sgRule.Ipv6Range.CidrIpv6
	} else if sgRule.Group != nil && *sgRule.UserIDGroupPair.GroupId == *sgRule.Group.GroupId {
		hashCode = hashCode + "_" + *sgRule.Group.GroupId
	} else if sgRule.UserIDGroupPair != nil && *sgRule.UserIDGroupPair.GroupId == *sgRule.Group.GroupId {
		hashCode = hashCode + "_" + *sgRule.UserIDGroupPair.GroupId
	}

	// generate aka for the rule
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + *sgRule.Group.OwnerId + ":security-group/" + *sgRule.Group.GroupId + ":" + hashCode}

	title := *sgRule.Group.GroupId + "_" + hashCode

	turbotData := map[string]interface{}{
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
	Type            string
}

// TODO - Get call
// func getSecurityGroupRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	// get the required name field value from keyValues
// 	sgGroup := h.Item.(*ec2.SecurityGroup)
// 	defaultRegion := GetDefaultRegion()

// 	// get service
// 	svc, err := Ec2Service(ctx, d, defaultRegion)
// 	if err != nil {
// 		return nil, err
// 	}

// 	params := &ec2.DescribeSecurityGroupsInput{
// 		GroupIds: []*string{sgGroup.GroupId},
// 	}

// 	items, err := svc.DescribeSecurityGroups(params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var rowSource []interface{}
// 	if items.SecurityGroups != nil && len(items.SecurityGroups) > 0 {
// 		if items.SecurityGroups[0].IpPermissions != nil {
// 			for _, permission := range items.SecurityGroups[0].IpPermissions {
// 				rowSource = append(rowSource, rowSourceFromIPPermission(items.SecurityGroups[0], permission, "ingress")...)
// 			}
// 		}

// 		if items.SecurityGroups[0].IpPermissionsEgress != nil {
// 			for _, permission := range items.SecurityGroups[0].IpPermissionsEgress {
// 				rowSource = append(rowSource, rowSourceFromIPPermission(items.SecurityGroups[0], permission, "ingress")...)
// 			}

// 		}
// 		return rowSource, nil
// 	}

// 	return []interface{}{}, nil
// }
