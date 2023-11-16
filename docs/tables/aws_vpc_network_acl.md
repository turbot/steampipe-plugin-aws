---
title: "Table: aws_vpc_network_acl - Query AWS VPC Network ACL using SQL"
description: "Allows users to query AWS VPC Network ACLs to retrieve detailed information about network access control lists in a specific AWS VPC."
---

# Table: aws_vpc_network_acl - Query AWS VPC Network ACL using SQL

The `aws_vpc_network_acl` table in Steampipe provides information about Network Access Control Lists (ACLs) within Amazon Virtual Private Cloud (VPC). This table allows DevOps engineers, security analysts, and system administrators to query ACL-specific details, including rules, associations, and related metadata. Users can utilize this table to gather insights on ACLs, such as rule configurations, associated subnets, and more. The schema outlines the various attributes of the Network ACL, including the ACL ID, VPC ID, default status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_network_acl` table, you can use the `.inspect aws_vpc_network_acl` command in Steampipe.

**Key columns**:

- `network_acl_id`: The ID of the network ACL. This is the primary key and can be used to join this table with other tables that contain network ACL information.
- `vpc_id`: The ID of the VPC for the network ACL. This can be used to join with other tables that contain VPC information, allowing for comprehensive queries across your VPC infrastructure.
- `default`: Indicates whether this is the default network ACL for the VPC. This can be useful in queries to check for custom ACLs or to ensure default ACLs have not been modified.

## Examples

### List the attached VPC IDs for each network ACL

```sql
select
  network_acl_id,
  arn,
  vpc_id
from
  aws_vpc_network_acl;
```


### List the default NACL associated with the VPCs

```sql
select
  network_acl_id,
  vpc_id,
  is_default
from
  aws_vpc_network_acl
where
  is_default = true;
```


### Subnet associated with each network ACL

```sql
select
  network_acl_id,
  vpc_id,
  association ->> 'SubnetId' as subnet_id,
  association ->> 'NetworkAclAssociationId' as network_acl_association_id
from
  aws_vpc_network_acl
  cross join jsonb_array_elements(associations) as association;
```