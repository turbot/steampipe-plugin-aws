# Table: aws_vpc_security_group_rule

Security group rules defines the inbound and outbound traffic to the instances.

## Examples

## List of security groups whose inbound access is open to the internet

```sql
select
  group_name,
  group_id,
  type
from
  aws_vpc_security_group_rule
where
  type = 'ingress'
  and cidr_ip = '0.0.0.0/0';
```


## List of security groups whose SSH and RDP access is not restricted from the internet

```sql
select
  group_name,
  group_id,
  ip_protocol,
  from_port,
  to_port,
  cidr_ip
from
  aws_vpc_security_group_rule
where
  type = 'ingress'
  and cidr_ip = '0.0.0.0/0'
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