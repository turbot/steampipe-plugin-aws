package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcSecurityGroupRules(_ context.Context) *plugin.Table {
    return &plugin.Table{
        Name:        "aws_vpc_security_group_rules",
        Description: "AWS VPC Security Group Rules",
        List: &plugin.ListConfig{
            Hydrate:       listSecurityGroupRules1,
        },
        GetMatrixItem: BuildRegionList,
        Columns: awsRegionalColumns([]*plugin.Column{
			{
                Name:        "security_group_rule_id",
                Description: "The ID of the security group rule.",
                Type:        proto.ColumnType_STRING,
            },
            {
                Name:        "group_id",
                Description: "The ID of the security group.",
                Type:        proto.ColumnType_STRING,
            },
			{
                Name:        "group_owner_id",
                Description: "The ID of the Amazon Web Services account that owns the security group.",
                Type:        proto.ColumnType_STRING,
            },
			{
                Name:        "is_egress",
                Description: "Indicates whether the security group rule is an outbound rule.",
                Type:        proto.ColumnType_BOOL,
            },
			{
				Name:        "ip_protocol",
				Description: "The IP protocol name (tcp, udp, icmp, icmpv6) or number [see Protocol Numbers ](http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml). Use -1 to specify all protocols.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "from_port",
				Description: "The start of port range for the TCP and UDP protocols, or an ICMP/ICMPv6 type. A value of -1 indicates all ICMP/ICMPv6 types. If you specify all ICMP/ICMPv6 types, you must specify all codes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "to_port",
				Description: "The end of port range for the TCP and UDP protocols, or an ICMP/ICMPv6 code. A value of -1 indicates all ICMP/ICMPv6 codes. If you specify all ICMP/ICMPv6 types, you must specify all codes.",
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
				Name:        "prefix_list_id",
				Description: "The ID of the prefix list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "referenced_group_id",
				Description: "The ID of the security group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReferencedGroupInfo.GroupId"),
			},
			{
				Name:        "referenced_peering_status",
				Description: "The status of a VPC peering connection, if applicable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "referenced_user_id",
				Description: "The Amazon Web Services account ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "referenced_vpc_id",
				Description: "The ID of the VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "referenced_vpc_peering_connection_id",
				Description: "The ID of the VPC peering connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The security group rule description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the security group rule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSecurityGroupRulesTurbotData,
				Transform:   transform.FromField("Arn"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getVpcSecurityGroupRuleTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSecurityGroupRulesTurbotData,
				Transform:   transform.FromField("Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSecurityGroupRulesTurbotData,
				Transform:   transform.FromField("Akas"),
			},
        }),
    }
}

//// LIST FUNCTION

func listSecurityGroupRules1(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
    region := d.KeyColumnQualString(matrixKeyRegion)
    plugin.Logger(ctx).Trace("listVpcSecurityGroupRules", "AWS_REGION", region)

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

//// HYDRATE FUNCTIONS

func getSecurityGroupRulesTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		"Arn": arn,
		"Akas":  akas,
		"Title": title,
	}

	return turbotData, nil
}

//// TRANSFORM FUNCTIONS

func getVpcSecurityGroupRuleTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	securityGroupRule := d.HydrateItem.(*ec2.SecurityGroupRule)

	// Get the resource tags
	if securityGroupRule.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range securityGroupRule.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}