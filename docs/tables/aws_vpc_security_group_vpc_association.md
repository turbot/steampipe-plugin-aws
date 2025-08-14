---
title: "Steampipe Table: aws_vpc_security_group_vpc_association - Query AWS VPC Security Group VPC Associations using SQL"
description: "Allows users to query AWS VPC Security Group VPC Associations to retrieve details about the associations between security groups and VPCs."
folder: "Config"
---

# Table: aws_vpc_security_group_vpc_association - Query AWS VPC Security Group VPC Associations using SQL

The AWS VPC Security Group VPC Association is a feature that allows security groups to be associated with multiple VPCs. This enables you to use the same security group rules across different VPCs, providing consistent security policies and simplifying security group management.

## Table Usage Guide

The `aws_vpc_security_group_vpc_association` table in Steampipe provides you with information about the associations between security groups and VPCs. This table allows you, as a DevOps engineer, to query association-specific details, including which security groups are associated with which VPCs, the state of these associations, and the VPC owner information. You can utilize this table to gather insights on security group usage across VPCs, identify cross-VPC security group associations, and more.

## Examples

### Basic info
Determine the associations between security groups and VPCs, and gain insights into the security group usage across your AWS infrastructure. This can help enhance your understanding of the security group distribution and identify potential security considerations.

```sql+postgres
select
  group_id,
  vpc_id,
  vpc_owner_id,
  state
from
  aws_vpc_security_group_vpc_association
order by
  group_id,
  vpc_id;
```

```sql+sqlite
select
  group_id,
  vpc_id,
  vpc_owner_id,
  state
from
  aws_vpc_security_group_vpc_association
order by
  group_id,
  vpc_id;
```

### List security group associations for a specific VPC
Identify all security groups associated with a specific VPC, allowing you to understand the security group coverage for that VPC.

```sql+postgres
select
  group_id,
  vpc_id,
  state
from
  aws_vpc_security_group_vpc_association
where
  vpc_id = 'vpc-0841c64fb5d4f3f43';
```

```sql+sqlite
select
  group_id,
  vpc_id,
  state
from
  aws_vpc_security_group_vpc_association
where
  vpc_id = 'vpc-0841c64fb5d4f3f43';
```

### List VPC associations for a specific security group
Discover all VPCs that a specific security group is associated with, helping you understand the scope of that security group's usage.

```sql+postgres
select
  group_id,
  vpc_id,
  vpc_owner_id,
  state
from
  aws_vpc_security_group_vpc_association
where
  group_id = 'sg-0be099bd3551846d1';
```

```sql+sqlite
select
  group_id,
  vpc_id,
  vpc_owner_id,
  state
from
  aws_vpc_security_group_vpc_association
where
  group_id = 'sg-0be099bd3551846d1';
```

### List active security group VPC associations
Explore security groups that are actively associated with VPCs, providing insights into current security group usage.

```sql+postgres
select
  group_id,
  vpc_id,
  vpc_owner_id
from
  aws_vpc_security_group_vpc_association
where
  state = 'associated';
```

```sql+sqlite
select
  group_id,
  vpc_id,
  vpc_owner_id
from
  aws_vpc_security_group_vpc_association
where
  state = 'associated';
```
