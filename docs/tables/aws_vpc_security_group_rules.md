# Table: aws_vpc_security_group_rules

Security group rules defines the inbound and outbound traffic to the instances.(using DescribeSecurityGroupRules)

## Examples

## List of security group rules whose inbound access is open to the internet

```sql
select
  group_id,
  security_group_rule_id,
  is_egress
from
  aws_vpc_security_group_rules
where
  is_egress is false
  and cidr_ipv4 = '0.0.0.0/0';
```

## List of security group rules whose SSH and RDP access is not restricted from the internet

```sql
select
  group_id,
  security_group_rule_id,
  ip_protocol,
  from_port,
  to_port,
  cidr_ipv4
from
  aws_vpc_security_group_rules
where
  is_egress is false
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

## List of security group rules join security groups

```sql
select
  sg.group_name,
  sgr.group_id,
  sgr.security_group_rule_id
from
  aws_vpc_security_group_rules as sgr
  join aws_vpc_security_group as sg on sg.group_id = sgr.group_id;
```
