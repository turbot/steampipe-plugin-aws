---
title: "Steampipe Table: aws_elasticache_replication_group - Query AWS ElastiCache Replication Groups using SQL"
description: "Allows users to query AWS ElastiCache Replication Groups to retrieve information related to their configuration, status, and associated resources."
folder: "ElastiCache"
---

# Table: aws_elasticache_replication_group - Query AWS ElastiCache Replication Groups using SQL

The AWS ElastiCache Replication Group is a feature of AWS ElastiCache that allows you to create a group of one or more cache clusters that are managed as a single entity. This enables the automatic partitioning of your data across multiple shards, providing enhanced performance, reliability, and scalability. Replication groups also support automatic failover, providing a high level of data availability.

## Table Usage Guide

The `aws_elasticache_replication_group` table in Steampipe provides you with information about replication groups within AWS ElastiCache. This table allows you, as a DevOps engineer, to query group-specific details, including configuration, status, and associated resources. You can utilize this table to gather insights on replication groups, such as their current status, associated cache clusters, node types, and more. The schema outlines the various attributes of the replication group for you, including the replication group ID, status, description, and associated tags.

## Examples

### Basic info
Determine the areas in which automatic failover is enabled in AWS ElastiCache, as well as whether authentication tokens are being used, to enhance security and ensure data redundancy. This query helps in identifying potential vulnerabilities and improving disaster recovery strategies.

```sql+postgres
select
  replication_group_id,
  description,
  cache_node_type,
  cluster_enabled,
  auth_token_enabled,
  automatic_failover
from
  aws_elasticache_replication_group;
```

```sql+sqlite
select
  replication_group_id,
  description,
  cache_node_type,
  cluster_enabled,
  auth_token_enabled,
  automatic_failover
from
  aws_elasticache_replication_group;
```

### List replication groups that are not encrypted at rest
Identify instances where replication groups in AWS ElastiCache are not encrypted at rest. This is useful to ensure data security by pinpointing potential vulnerabilities.

```sql+postgres
select
  replication_group_id,
  cache_node_type,
  at_rest_encryption_enabled
from
  aws_elasticache_replication_group
where
  not at_rest_encryption_enabled;
```

```sql+sqlite
select
  replication_group_id,
  cache_node_type,
  at_rest_encryption_enabled
from
  aws_elasticache_replication_group
where
  at_rest_encryption_enabled = 0;
```

### List replication groups with multi-AZ disabled
Determine the areas in which replication groups have multi-AZ disabled to assess potential vulnerabilities in your AWS ElastiCache setup.

```sql+postgres
select
  replication_group_id,
  cache_node_type,
  multi_az
from
  aws_elasticache_replication_group
where
  multi_az = 'disabled';
```

```sql+sqlite
select
  replication_group_id,
  cache_node_type,
  multi_az
from
  aws_elasticache_replication_group
where
  multi_az = 'disabled';
```

### List replication groups whose backup retention period is less than 30 days
Determine the areas in which backup retention periods for replication groups fall short of a 30-day standard, allowing for timely adjustments to ensure data safety.

```sql+postgres
select
  replication_group_id,
  snapshot_retention_limit,
  snapshot_window,
  snapshotting_cluster_id
from
  aws_elasticache_replication_group
where
  snapshot_retention_limit < 30;
```

```sql+sqlite
select
  replication_group_id,
  snapshot_retention_limit,
  snapshot_window,
  snapshotting_cluster_id
from
  aws_elasticache_replication_group
where
  snapshot_retention_limit < 30;
```

### List replication groups by node type
Explore which node types are used in your replication groups and determine their frequency. This can help optimize resource allocation and improve system performance.

```sql+postgres
select
  cache_node_type,
  count (*)
from
  aws_elasticache_replication_group
group by
  cache_node_type;
```

```sql+sqlite
select
  cache_node_type,
  count (*)
from
  aws_elasticache_replication_group
group by
  cache_node_type;
```

### List member clusters for each replication group
Explore the relationships within your replication groups by identifying which member clusters belong to each group. This helps in understanding the distribution and organization of your data across different clusters.

```sql+postgres
select
  replication_group_id,
  jsonb_array_elements_text(member_clusters) as member_clusters
from
  aws_elasticache_replication_group;
```

```sql+sqlite
select
  replication_group_id,
  json_each.value as member_clusters
from
  aws_elasticache_replication_group,
  json_each(aws_elasticache_replication_group.member_clusters);
```

### Find all ElastiCache update actions for replication groups
Retrieve a list of all service updates that are associated with ElastiCache replication groups.

```sql+postgres
select
  c.replication_group_id,
  c.engine,
  c.engine_version,
  c.cache_node_type,
  a.service_update_name,
  a.service_update_severity,
  a.service_update_status,
  a.service_update_type,
  a.service_update_recommended_apply_by_date
from
  aws_elasticache_cluster as c
  join aws_elasticache_update_action as a on c.replication_group_id = a.replication_group_id
where
  a.service_update_status = 'available'
order by
  a.service_update_severity,
  a.service_update_recommended_apply_by_date;
```

```sql+sqlite
select
  c.replication_group_id,
  c.engine,
  c.engine_version,
  c.cache_node_type,
  a.service_update_name,
  a.service_update_severity,
  a.service_update_status,
  a.service_update_type,
  a.service_update_recommended_apply_by_date
from
  aws_elasticache_cluster as c
  join aws_elasticache_update_action as a on c.replication_group_id = a.replication_group_id
where
  a.service_update_status = 'available'
order by
  a.service_update_severity,
  a.service_update_recommended_apply_by_date;
```