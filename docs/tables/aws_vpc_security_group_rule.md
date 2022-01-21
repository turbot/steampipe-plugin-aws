# Table: aws_vpc_security_group_rule

Security group rules defines the inbound and outbound traffic to the instances.

## Examples

## List of security groups whose inbound access is open to the internet

```sql
select
  security_group_rule_id,
  group_id,
  type
from
  aws_vpc_security_group_rule
where
  cidr_ipv4 = '0.0.0.0/0'
  and not is_egress;
```


## List of security groups whose SSH and RDP access is not restricted from the internet

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