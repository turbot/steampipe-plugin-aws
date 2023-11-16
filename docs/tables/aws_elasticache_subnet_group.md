---
title: "Table: aws_elasticache_subnet_group - Query AWS ElastiCache Subnet Groups using SQL"
description: "Allows users to query AWS ElastiCache Subnet Groups, providing details about each subnet group within their ElastiCache service, including the associated VPC, subnets, and status."
---

# Table: aws_elasticache_subnet_group - Query AWS ElastiCache Subnet Groups using SQL

The `aws_elasticache_subnet_group` table in Steampipe provides information about each subnet group within the AWS ElastiCache service. This table enables DevOps engineers to query subnet group-specific details, such as the associated VPC, subnets, and status. Users can utilize this table to gather insights on subnet groups, such as their availability status, the number of subnets within each group, and the VPC they're associated with. The schema outlines the various attributes of the subnet group, including the subnet group name, description, VPC ID, and subnet IDs.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_subnet_group` table, you can use the `.inspect aws_elasticache_subnet_group` command in Steampipe.

### Key columns:

- `name`: The name of the subnet group. This can be used to join with other tables that reference subnet groups by name.
- `vpc_id`: The ID of the VPC that the subnet group is associated with. This can be used to join with other tables that reference VPCs by ID.
- `subnet_ids`: The IDs of the subnets within the subnet group. This can be used to join with other tables that reference subnets by ID.

## Examples

### Basic info

```sql
select
  cache_subnet_group_name,
  cache_subnet_group_description,
  region,
  account_id
from
  aws_elasticache_subnet_group;
```


### Get network info for each subnet group

```sql
select
  vpc_id,
  sub -> 'SubnetAvailabilityZone' ->> 'Name' as subnet_availability_zone,
  sub ->> 'SubnetIdentifier' as subnet_identifier,
  sub ->> 'SubnetOutpost' as subnet_outpost
from
  aws_elasticache_subnet_group,
  jsonb_array_elements(subnets) as sub;
```


### List ElastiCache clusters in each subnet group

```sql
select
  c.cache_cluster_id,
  sg.cache_subnet_group_name,
  sg.vpc_id
from
  aws_elasticache_subnet_group as sg
  join aws_elasticache_cluster as c on sg.cache_subnet_group_name = c.cache_subnet_group_name;
```
