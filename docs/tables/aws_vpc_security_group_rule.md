# Table: aws_vpc_security_group_rule

Security group rules define the inbound and outbound traffic to the instances.

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
