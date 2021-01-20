# Table: aws_vpc_security_group

A security group acts as a virtual firewall for EC2 instances to control incoming and outgoing traffic.

## Examples

### Basic ingress rule info

```sql
select
  group_name,
  vpc_id,
  perm ->> 'FromPort' as from_port,
  perm ->> 'ToPort' as to_port,
  perm ->> 'IpProtocol' as ip_protocol,
  perm ->> 'IpRanges' as ip_ranges,
  perm ->> 'Ipv6Ranges' as ipv6_ranges,
  perm ->> 'UserIdGroupPairs' as user_id_group_pairs,
  perm ->> 'PrefixListIds' as prefix_list_ids
from
  aws_vpc_security_group as sg
  cross join jsonb_array_elements(ip_permissions) as perm;
```


### List of security groups whose SSH and RDP access is not restricted from the internet

```sql
select
  sg.group_name,
  sg.group_id,
  sgr.type,
  sgr.ip_protocol,
  sgr.from_port,
  sgr.to_port,
  cidr_ip
from
  aws_vpc_security_group as sg
  join aws_vpc_security_group_rule as sgr on sg.group_name = sgr.group_name
where
  sgr.type = 'ingress'
  and sgr.cidr_ip = '0.0.0.0/0'
  and (
    (
      sgr.ip_protocol = '-1' -- all traffic
      and sgr.from_port is null
    )
    or (
      sgr.from_port <= 22
      and sgr.to_port >= 22
    )
    or (
      sgr.from_port <= 3389
      and sgr.to_port >= 3389
    )
  );
```


### Count of security groups by VPC ID

```sql
select
  vpc_id,
  count(vpc_id) as count
from
  aws_vpc_security_group
group by
  vpc_id;
```


### List of security groups whose name is prefixed with 'launch wizard'

```sql
select
  group_name,
  group_id
from
  aws_vpc_security_group
where
  group_name like '%launch-wizard%';
```