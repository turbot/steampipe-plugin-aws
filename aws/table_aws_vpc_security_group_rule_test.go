package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestSecurityGroupRuleHashCode(t *testing.T) {
	testCases := []struct {
		name string
		rule types.SecurityGroupRule
		want string
	}{
		{
			name: "ingress tcp port range with IPv4 CIDR",
			rule: types.SecurityGroupRule{
				IpProtocol: aws.String("tcp"),
				IsEgress:   aws.Bool(false),
				FromPort:   aws.Int32(22),
				ToPort:     aws.Int32(22),
				CidrIpv4:   aws.String("10.0.0.0/16"),
			},
			want: "ingress_tcp_22_22_10.0.0.0/16",
		},
		{
			name: "egress all traffic uses -1 ports",
			rule: types.SecurityGroupRule{
				IpProtocol: aws.String("-1"),
				IsEgress:   aws.Bool(true),
				FromPort:   aws.Int32(-1),
				ToPort:     aws.Int32(-1),
				CidrIpv4:   aws.String("0.0.0.0/0"),
			},
			want: "egress_-1_-1_-1_0.0.0.0/0",
		},
		{
			name: "no ports with IPv6 CIDR",
			rule: types.SecurityGroupRule{
				IpProtocol: aws.String("-1"),
				IsEgress:   aws.Bool(true),
				CidrIpv6:   aws.String("::/0"),
			},
			want: "egress_-1_::/0",
		},
		{
			name: "from port without to port",
			rule: types.SecurityGroupRule{
				IpProtocol: aws.String("icmp"),
				IsEgress:   aws.Bool(false),
				FromPort:   aws.Int32(8),
				CidrIpv4:   aws.String("10.0.0.0/8"),
			},
			want: "ingress_icmp_8_0_10.0.0.0/8",
		},
		{
			name: "self-referencing group",
			rule: types.SecurityGroupRule{
				IpProtocol:          aws.String("tcp"),
				IsEgress:            aws.Bool(false),
				FromPort:            aws.Int32(443),
				ToPort:              aws.Int32(443),
				GroupId:             aws.String("sg-111"),
				ReferencedGroupInfo: &types.ReferencedSecurityGroup{GroupId: aws.String("sg-111")},
			},
			want: "ingress_tcp_443_443_sg-111",
		},
		{
			name: "prefix list",
			rule: types.SecurityGroupRule{
				IpProtocol:   aws.String("tcp"),
				IsEgress:     aws.Bool(true),
				FromPort:     aws.Int32(443),
				ToPort:       aws.Int32(443),
				GroupId:      aws.String("sg-111"),
				PrefixListId: aws.String("pl-123"),
			},
			want: "egress_tcp_443_443_pl-123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := securityGroupRuleHashCode(tc.rule); got != tc.want {
				t.Errorf("securityGroupRuleHashCode() = %q, want %q", got, tc.want)
			}
		})
	}
}
