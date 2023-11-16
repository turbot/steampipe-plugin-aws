---
title: "Table: aws_elasticache_replication_group - Query AWS ElastiCache Replication Groups using SQL"
description: "Allows users to query AWS ElastiCache Replication Groups to retrieve information related to their configuration, status, and associated resources."
---

# Table: aws_elasticache_replication_group - Query AWS ElastiCache Replication Groups using SQL

The `aws_elasticache_replication_group` table in Steampipe provides information about replication groups within AWS ElastiCache. This table allows DevOps engineers to query group-specific details, including configuration, status, and associated resources. Users can utilize this table to gather insights on replication groups, such as their current status, associated cache clusters, node types, and more. The schema outlines the various attributes of the replication group, including the replication group ID, status, description, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_replication_group` table, you can use the `.inspect aws_elasticache_replication_group` command in Steampipe.

**Key columns**:

- `replication_group_id`: The identifier for the replication group. This column can be used to join this table with other tables that contain information about specific replication groups.
- `status`: The current state of this replication group - creating, available, modifying, deleting, etc. This column is useful for tracking the status of replication groups and identifying any potential issues.
- `node_type`: The current node type of the replication group, such as cache.t2.micro. This column is important for understanding the resources allocated to each replication group.

## Examples

### Basic info

```sql
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

```sql
select
  replication_group_id,
  cache_node_type,
  at_rest_encryption_enabled
from
  aws_elasticache_replication_group
where
  not at_rest_encryption_enabled;
```


### List replication groups with multi-AZ disabled

```sql
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

```sql
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

```sql
select
  cache_node_type,
  count (*)
from
  aws_elasticache_replication_group
group by
  cache_node_type;
```


### List member clusters for each replication group

```sql
select
  replication_group_id,
  jsonb_array_elements_text(member_clusters) as member_clusters
from
  aws_elasticache_replication_group;
```
