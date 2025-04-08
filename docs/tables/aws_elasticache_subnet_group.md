---
title: "Steampipe Table: aws_elasticache_subnet_group - Query AWS ElastiCache Subnet Groups using SQL"
description: "Allows users to query AWS ElastiCache Subnet Groups, providing details about each subnet group within their ElastiCache service, including the associated VPC, subnets, and status."
folder: "ElastiCache"
---

# Table: aws_elasticache_subnet_group - Query AWS ElastiCache Subnet Groups using SQL

The AWS ElastiCache Subnet Group is a feature of Amazon ElastiCache that allows you to specify a set of subnets in your Virtual Private Cloud (VPC) network where you can launch your ElastiCache clusters. A subnet group is a collection of subnets (typically private) that you can designate for your clusters running in a VPC, providing the ability to partition ElastiCache services across your VPC subnets and Availability Zones. This facilitates high availability, fault tolerance, and scalability.

## Table Usage Guide

The `aws_elasticache_subnet_group` table in Steampipe provides you with information about each subnet group within the AWS ElastiCache service. This table enables you, as a DevOps engineer, to query subnet group-specific details, such as the associated VPC, subnets, and status. You can utilize this table to gather insights on subnet groups, such as their availability status, the number of subnets within each group, and the VPC they're associated with. The schema outlines the various attributes of the subnet group for you, including the subnet group name, description, VPC ID, and subnet IDs.

## Examples

### Basic info
Discover the segments that make up your AWS Elasticache subnet groups, including their descriptions and associated regions, to better manage your cloud resources. This could be particularly useful for understanding your resource allocation and planning for future infrastructure needs.

```sql+postgres
select
  cache_subnet_group_name,
  cache_subnet_group_description,
  region,
  account_id
from
  aws_elasticache_subnet_group;
```

```sql+sqlite
select
  cache_subnet_group_name,
  cache_subnet_group_description,
  region,
  account_id
from
  aws_elasticache_subnet_group;
```


### Get network info for each subnet group
Explore the specific network configuration for each subnet group in your AWS ElastiCache service. This query is useful in identifying the availability zone, identifier, and outpost details for each subnet, which can assist in network management and troubleshooting.

```sql+postgres
select
  vpc_id,
  sub -> 'SubnetAvailabilityZone' ->> 'Name' as subnet_availability_zone,
  sub ->> 'SubnetIdentifier' as subnet_identifier,
  sub ->> 'SubnetOutpost' as subnet_outpost
from
  aws_elasticache_subnet_group,
  jsonb_array_elements(subnets) as sub;
```

```sql+sqlite
select
  vpc_id,
  json_extract(sub.value, '$.SubnetAvailabilityZone.Name') as subnet_availability_zone,
  json_extract(sub.value, '$.SubnetIdentifier') as subnet_identifier,
  json_extract(sub.value, '$.SubnetOutpost') as subnet_outpost
from
  aws_elasticache_subnet_group,
  json_each(subnets) as sub;
```


### List ElastiCache clusters in each subnet group
Determine the distribution of ElastiCache clusters within your subnet groups to better understand your AWS resource allocation. This can help in optimizing resource usage and managing costs effectively.

```sql+postgres
select
  c.cache_cluster_id,
  sg.cache_subnet_group_name,
  sg.vpc_id
from
  aws_elasticache_subnet_group as sg
  join aws_elasticache_cluster as c on sg.cache_subnet_group_name = c.cache_subnet_group_name;
```

```sql+sqlite
select
  c.cache_cluster_id,
  sg.cache_subnet_group_name,
  sg.vpc_id
from
  aws_elasticache_subnet_group as sg
  join aws_elasticache_cluster as c on sg.cache_subnet_group_name = c.cache_subnet_group_name;
```