---
title: "Steampipe Table: aws_vpc_security_group_rule - Query AWS VPC Security Group Rule using SQL"
description: "Allows users to query AWS VPC Security Group Rule, providing detailed information about security group rules within Amazon Virtual Private Cloud (VPC)."
folder: "VPC"
---

# Table: aws_vpc_security_group_rule - Query AWS VPC Security Group Rule using SQL

The AWS VPC Security Group Rule is a feature within Amazon's Virtual Private Cloud (VPC) service that allows you to manage inbound and outbound traffic for your instances and subnets. These rules are used to control the traffic flow in a security group, offering flexibility to permit or deny specified traffic. This ensures a secure environment by providing a robust set of firewall rules at the instance and subnet level.

## Table Usage Guide

The `aws_vpc_security_group_rule` table in Steampipe provides you with information about security group rules within Amazon Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, security analyst, or system administrator, to query rule-specific details, including rule type, IP protocol, port range, and associated metadata. You can utilize this table to gather insights on security group rules, such as rules with open IP ranges, verification of port ranges, and more. The schema outlines the various attributes of the security group rule for you, including the rule ID, security group ID, IP range, and associated tags.

## Examples

## List inbound security group rules open to the Internet

```sql+postgres
select
  security_group_rule_id,
  group_id,
  ip_protocol,
  from_port,
  to_port
from
  aws_vpc_security_group_rule
where
  cidr_ipv4 = '0.0.0.0/0'
  and not is_egress;
```

```sql+sqlite
select
  security_group_rule_id,
  group_id,
  ip_protocol,
  from_port,
  to_port
from
  aws_vpc_security_group_rule
where
  cidr_ipv4 = '0.0.0.0/0'
  and not is_egress;
```

## List ingress security group rules that open SSH and RDP access from the Internet

```sql+postgres
select
  security_group_rule_id,
  group_id,
  ip_protocol,
  from_port,
  to_port,
  cidr_ipv4
from
  aws_vpc_security_group_rule
where
  not is_egress
  and cidr_ipv4 = '0.0.0.0/0'
  and (
    (
      ip_protocol = '-1' -- all traffic
      and from_port is null
    )
    or (
      from_port <= 22
      and to_port >= 22
    )
    or (
      from_port <= 3389
      and to_port >= 3389
    )
  );
```

```sql+sqlite
select
  security_group_rule_id,
  group_id,
  ip_protocol,
  from_port,
  to_port,
  cidr_ipv4
from
  aws_vpc_security_group_rule
where
  not is_egress
  and cidr_ipv4 = '0.0.0.0/0'
  and (
    (
      ip_protocol = '-1' -- all traffic
      and from_port is null
    )
    or (
      from_port <= 22
      and to_port >= 22
    )
    or (
      from_port <= 3389
      and to_port >= 3389
    )
  );
```

### List security group rules with additional security group details
Discover the segments that have security group rules, along with additional details of the security group. This is particularly useful in understanding the security configurations in your environment, helping you to enhance protection measures.

```sql+postgres
select
  r.security_group_rule_id,
  r.ip_protocol,
  r.from_port,
  r.to_port,
  r.cidr_ipv4,
  r.group_id,
  sg.group_name
from
  aws_vpc_security_group_rule as r,
  aws_vpc_security_group as sg 
where
  r.group_id = sg.group_id;
```

```sql+sqlite
select
  r.security_group_rule_id,
  r.ip_protocol,
  r.from_port,
  r.to_port,
  r.cidr_ipv4,
  r.group_id,
  sg.group_name
from
  aws_vpc_security_group_rule as r
join
  aws_vpc_security_group as sg 
on
  r.group_id = sg.group_id;
```