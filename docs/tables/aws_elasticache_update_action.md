---
title: "Steampipe Table: aws_elasticache_update_action - Query Amazon ElastiCache Update Actions using SQL"
description: "Allows users to query Amazon ElastiCache Update Actions, providing information about service updates for ElastiCache clusters and nodes within the AWS account."
folder: "ElastiCache"
---

# Table: aws_elasticache_update_action - Query Amazon ElastiCache Update Actions using SQL

The Amazon ElastiCache Update Action is a feature of AWS's ElastiCache service that allows users to manage and schedule updates for their ElastiCache clusters and nodes. This resource is designed to ensure that your caching infrastructure remains secure, performant, and compliant with the latest updates and patches.

## Table Usage Guide

The `aws_elasticache_update_action` table in Steampipe provides you with information about service updates for ElastiCache clusters and nodes within your AWS account. This table enables you, as a DevOps engineer, database administrator, or IT professional, to query update-specific details, including the status of updates, recommended application dates, and more. You can utilize this table to gather insights on updates, such as their severity, type, and compliance requirements. The schema outlines the various attributes of the ElastiCache Update Action for you, including the update action ID, associated cache cluster and replication group IDs, and relevant dates.

## Examples

### List all pending ElastiCache update actions
Retrieve a list of all service updates that are pending for your ElastiCache clusters and nodes. This is essential for ensuring that your caching infrastructure is up to date.

```sql+postgres
select
  cache_cluster_id,
  replication_group_id,
  engine,
  estimated_update_time,
  nodes_updated,
  service_update_name,
  service_update_recommended_apply_by_date,
  service_update_release_date,
  service_update_severity,
  service_update_status,
  service_update_type,
  sla_met,
  update_action_available_date,
  update_action_status,
  update_action_status_modified_date
from
  aws_elasticache_update_action;
```

```sql+sqlite
select
  cache_cluster_id,
  replication_group_id,
  engine,
  estimated_update_time,
  nodes_updated,
  service_update_name,
  service_update_recommended_apply_by_date,
  service_update_release_date,
  service_update_severity,
  service_update_status,
  service_update_type,
  sla_met,
  update_action_available_date,
  update_action_status,
  update_action_status_modified_date
from
  aws_elasticache_update_action;
```

### Find all ElastiCache update actions with important severity
Retrieve a list of all service updates that have an important severity level. This can help you prioritize and address critical updates first.

```sql+postgres
select 
  cache_cluster_id,
  replication_group_id,
  engine,
  service_update_severity,
  nodes_updated,
  service_update_name,
  service_update_recommended_apply_by_date,
  service_update_release_date
from 
  aws_elasticache_update_action 
where 
  service_update_severity='important'
```

```sql+sqlite
select 
  cache_cluster_id,
  replication_group_id,
  engine,
  service_update_severity,
  nodes_updated,
  service_update_name,
  service_update_recommended_apply_by_date,
  service_update_release_date 
from 
  aws_elasticache_update_action 
where 
  service_update_severity='important'
```

### Find all ElastiCache update actions for a specific cache cluster
Retrieve a list of all service updates that are associated with a specific ElastiCache cluster. This can help you track updates for a specific cluster.

```sql+postgres
select 
  cache_cluster_id,
  replication_group_id,
  estimated_update_time,
  service_update_recommended_apply_by_date,
  service_update_release_date,
  service_update_type
from 
  aws_elasticache_update_action 
where 
  cache_cluster_id='minutes-auth-qa-ec'
```

```sql+sqlite
select 
  cache_cluster_id,
  replication_group_id,
  estimated_update_time,
  service_update_recommended_apply_by_date,
  service_update_release_date,
  service_update_type
from 
  aws_elasticache_update_action 
where 
  cache_cluster_id='minutes-auth-qa-ec'
```

### Find all ElastiCache update actions for a specific time range
The range of time specified to search for service updates that are in available status

```sql+postgres
select
  cache_cluster_id,
  replication_group_id,
  service_update_name,
  service_update_status,
  update_action_status,
  service_update_release_date
from
  aws_elasticache_update_action
where
  service_update_release_date >= '2024-02-05T18:59:59+08:00' and service_update_release_date <= '2025-02-05T18:59:59+08:00';
```

```sql+sqlite
select
  cache_cluster_id,
  replication_group_id,
  service_update_name,
  service_update_status,
  update_action_status,
  service_update_release_date
from
  aws_elasticache_update_action
where
  service_update_release_date > '2024-02-05T18:59:59+08:00' and service_update_release_date < '2025-02-05T18:59:59+08:00';
```

### Find all ElastiCache update actions for replication groups
Retrieve a list of all service updates that are associated with ElastiCache replication groups.

```sql+postgres
select
  c.cache_cluster_id,
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

### Find all ElastiCache update actions for cache clusters
Retrieve a list of all service updates that are associated with ElastiCache clusters.

```sql+postgres
select
  c.cache_cluster_id,
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
  join aws_elasticache_update_action as a on c.cache_cluster_id = a.cache_cluster_id
where
  a.service_update_status = 'available'
order by
  a.service_update_severity,
  a.service_update_recommended_apply_by_date;
```

```sql+sqlite
select
  c.cache_cluster_id,
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
  join aws_elasticache_update_action as a on c.cache_cluster_id = a.cache_cluster_id
where
  a.service_update_status = 'available'
order by
  a.service_update_severity,
  a.service_update_recommended_apply_by_date;
```