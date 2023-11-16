---
title: "Table: aws_vpc_security_group_rule - Query AWS VPC Security Group Rule using SQL"
description: "Allows users to query AWS VPC Security Group Rule, providing detailed information about security group rules within Amazon Virtual Private Cloud (VPC)."
---

# Table: aws_vpc_security_group_rule - Query AWS VPC Security Group Rule using SQL

The `aws_vpc_security_group_rule` table in Steampipe provides information about security group rules within Amazon Virtual Private Cloud (VPC). This table allows DevOps engineers, security analysts, and system administrators to query rule-specific details, including rule type, IP protocol, port range, and associated metadata. Users can utilize this table to gather insights on security group rules, such as rules with open IP ranges, verification of port ranges, and more. The schema outlines the various attributes of the security group rule, including the rule ID, security group ID, IP range, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_security_group_rule` table, you can use the `.inspect aws_vpc_security_group_rule` command in Steampipe.

### Key columns:

- `rule_id`: This is the unique identifier for the security group rule. This column is important as it allows users to uniquely identify and query specific security group rules.
- `group_id`: This is the ID of the security group that the rule belongs to. This column is useful for joining this table with the `aws_vpc_security_group` table to get more information about the security group.
- `ip_protocol`: This column indicates the IP protocol used in the rule. This is useful for analyzing the protocols used in security group rules.

## Examples

## List inbound security group rules open to the Internet

```sql
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

```sql
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

```sql
select
  r.security_group_rule_id,
  r.ip_protocol,
  r.from_port,
  r.to_port,
  r.cidr_ipv4,
  r.group_id,
  sg.group_name,
  sg.vpc_id 
from
  aws_vpc_security_group_rule as r,
  aws_vpc_security_group as sg 
where
  r.group_id = sg.group_id;
```
