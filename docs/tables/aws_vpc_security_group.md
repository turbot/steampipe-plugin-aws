---
title: "Steampipe Table: aws_vpc_security_group - Query AWS VPC Security Groups using SQL"
description: "Allows users to query AWS VPC Security Groups and retrieve data such as group ID, name, description, owner ID, and associated VPC ID. This table can be used to gain insights on security group configurations, policies, and related metadata."
folder: "VPC"
---

# Table: aws_vpc_security_group - Query AWS VPC Security Groups using SQL

An AWS VPC Security Group acts as a virtual firewall for your instance to control inbound and outbound traffic. When you launch an instance in a VPC, you can assign up to five security groups to the instance. Security groups act at the instance level, not the subnet level, therefore each instance in a subnet in your VPC can be assigned to a different set of security groups.

## Table Usage Guide

The `aws_vpc_security_group` table in Steampipe provides you with information about Security Groups within AWS Virtual Private Cloud (VPC). This table enables you, as a DevOps engineer, to query security group-specific details, including configurations, associated policies, and related metadata. You can utilize this table to gather insights on security groups, such as understanding the security rules applied, verifying the security policies, and more. The schema outlines for you the various attributes of the security group, including the group ID, name, description, owner ID, and associated VPC ID.

## Examples

### Basic ingress rule info
Review the configuration of your network's security settings to understand which ports are open and the protocols being used. This can help in identifying potential security vulnerabilities and ensuring that the network is adequately protected.

```sql+postgres
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

```sql+sqlite
select
  group_name,
  vpc_id,
  json_extract(perm.value, '$.FromPort') as from_port,
  json_extract(perm.value, '$.ToPort') as to_port,
  json_extract(perm.value, '$.IpProtocol') as ip_protocol,
  json_extract(perm.value, '$.IpRanges') as ip_ranges,
  json_extract(perm.value, '$.Ipv6Ranges') as ipv6_ranges,
  json_extract(perm.value, '$.UserIdGroupPairs') as user_id_group_pairs,
  json_extract(perm.value, '$.PrefixListIds') as prefix_list_ids
from
  aws_vpc_security_group as sg,
  json_each(ip_permissions) as perm;
```

### List of security groups whose SSH and RDP access is not restricted from the internet
Explore which security groups have unrestricted SSH and RDP access from the internet, allowing you to identify potential security risks and tighten access controls. This is crucial for maintaining security standards and preventing unauthorized access.

```sql+postgres
select
  sg.group_name,
  sg.group_id,
  sgr.type,
  sgr.ip_protocol,
  sgr.from_port,
  sgr.to_port,
  cidr_ipv4
from
  aws_vpc_security_group as sg
  join aws_vpc_security_group_rule as sgr on sg.group_id = sgr.group_id
where
  sgr.type = 'ingress'
  and sgr.cidr_ipv4 = '0.0.0.0/0'
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

```sql+sqlite
select
  sg.group_name,
  sg.group_id,
  sgr.type,
  sgr.ip_protocol,
  sgr.from_port,
  sgr.to_port,
  cidr_ipv4
from
  aws_vpc_security_group as sg
  join aws_vpc_security_group_rule as sgr on sg.group_id = sgr.group_id
where
  sgr.type = 'ingress'
  and sgr.cidr_ipv4 = '0.0.0.0/0'
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
Gain insights into the distribution of security groups across your Virtual Private Clouds (VPCs) to manage resources and improve security measures effectively.

```sql+postgres
select
  vpc_id,
  count(vpc_id) as count
from
  aws_vpc_security_group
group by
  vpc_id;
```

```sql+sqlite
select
  vpc_id,
  count(vpc_id) as count
from
  aws_vpc_security_group
group by
  vpc_id;
```


### List of security groups whose name is prefixed with 'launch wizard'
Identify instances where security groups have names prefixed with 'launch wizard'. This can help in managing and organizing your security groups effectively, particularly in large-scale AWS environments.

```sql+postgres
select
  group_name,
  group_id
from
  aws_vpc_security_group
where
  group_name like '%launch-wizard%';
```

```sql+sqlite
select
  group_name,
  group_id
from
  aws_vpc_security_group
where
  group_name like '%launch-wizard%';
```