---
title: "Steampipe Table: aws_vpc_network_acl - Query AWS VPC Network ACL using SQL"
description: "Allows users to query AWS VPC Network ACLs to retrieve detailed information about network access control lists in a specific AWS VPC."
folder: "VPC"
---

# Table: aws_vpc_network_acl - Query AWS VPC Network ACL using SQL

The AWS VPC Network ACL is a security layer that controls traffic in and out of a Virtual Private Cloud (VPC). It operates at the subnet level and evaluates traffic based on defined rules in a numbered list. This Network Access Control List (ACL) provides an additional line of defense for your VPC and can be customized to fit your security needs.

## Table Usage Guide

The `aws_vpc_network_acl` table in Steampipe provides you with information about Network Access Control Lists (ACLs) within Amazon Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, security analyst, or system administrator, to query ACL-specific details, including rules, associations, and related metadata. You can utilize this table to gather insights on ACLs, such as rule configurations, associated subnets, and more. The schema outlines the various attributes of the Network ACL for you, including the ACL ID, VPC ID, default status, and associated tags.

## Examples

### List the attached VPC IDs for each network ACL
Explore which network access control lists (ACLs) are associated with each virtual private cloud (VPC) in your AWS environment. This can help you manage and monitor your network security by identifying the VPCs that are linked to each ACL.

```sql+postgres
select
  network_acl_id,
  arn,
  vpc_id
from
  aws_vpc_network_acl;
```

```sql+sqlite
select
  network_acl_id,
  arn,
  vpc_id
from
  aws_vpc_network_acl;
```


### List the default NACL associated with the VPCs
Determine the areas in which the default Network Access Control List (NACL) is associated with the Virtual Private Clouds (VPCs). This is useful to understand the default security settings of your network resources in AWS.

```sql+postgres
select
  network_acl_id,
  vpc_id,
  is_default
from
  aws_vpc_network_acl
where
  is_default = true;
```

```sql+sqlite
select
  network_acl_id,
  vpc_id,
  is_default
from
  aws_vpc_network_acl
where
  is_default = 1;
```


### Subnet associated with each network ACL
Determine the areas in which each network access control list (ACL) is associated with a specific subnet. This is useful for understanding your network's security configuration and identifying any potential vulnerabilities or misconfigurations.

```sql+postgres
select
  network_acl_id,
  vpc_id,
  association ->> 'SubnetId' as subnet_id,
  association ->> 'NetworkAclAssociationId' as network_acl_association_id
from
  aws_vpc_network_acl
  cross join jsonb_array_elements(associations) as association;
```

```sql+sqlite
select
  network_acl_id,
  vpc_id,
  json_extract(association.value, '$.SubnetId') as subnet_id,
  json_extract(association.value, '$.NetworkAclAssociationId') as network_acl_association_id
from
  aws_vpc_network_acl,
  json_each(associations) as association;
```