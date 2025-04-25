---
title: "Steampipe Table: aws_elasticache_redis_metric_list_based_cmds_hourly - Query AWS ElastiCache Redis Metrics using SQL"
description: "Allows users to query ElastiCache Redis Metrics on an hourly basis, providing data on list-based commands executed in the ElastiCache Redis environment."
folder: "ElastiCache"
---

# Table: aws_elasticache_redis_metric_list_based_cmds_hourly - Query AWS ElastiCache Redis Metrics using SQL

The AWS ElastiCache Redis Metrics service allows you to monitor, isolate, and diagnose performance issues in your ElastiCache Redis environments using SQL. It provides important insights into the operational health of your ElastiCache Redis instances by collecting and analyzing key database performance metrics. This service enables efficient troubleshooting and performance optimization of your ElastiCache Redis environments.

## Table Usage Guide

The `aws_elasticache_redis_metric_list_based_cmds_hourly` table in Steampipe provides you with information about list-based command metrics within AWS ElastiCache Redis. This table allows you, as a DevOps engineer, to query command-specific details on an hourly basis, including the number of commands processed, the latency of commands, and associated metadata. You can utilize this table to gather insights on command performance, such as identifying high latency commands, tracking the frequency of command usage, and more. The schema outlines the various attributes of the ElastiCache Redis command metrics, including the cache cluster id, the metric name, and the timestamp.

The `aws_elasticache_redis_metric_list_based_cmds_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Determine the performance trends of your ElastiCache Redis clusters by analyzing hourly metrics. This can help in identifying patterns, optimizing resource usage and planning for capacity upgrades.

```sql+postgres
select
  cache_cluster_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count,
  sum
from
  aws_elasticache_redis_metric_list_based_cmds_hourly
order by
  cache_cluster_id,
  timestamp;
```

```sql+sqlite
select
  cache_cluster_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count,
  sum
from
  aws_elasticache_redis_metric_list_based_cmds_hourly
order by
  cache_cluster_id,
  timestamp;
```

### listbasedcmds sum over 100 
This query is useful for monitoring your AWS ElastiCache Redis clusters by identifying instances where the sum of list-based commands executed per hour exceeds 100. This can help in optimizing your cache usage by pinpointing areas of high command activity.

```sql+postgres
select
  cache_cluster_id,
  timestamp,
  round(minimum::numeric,2) as min_listbasedcmds,
  round(maximum::numeric,2) as max_listbasedcmds,
  round(average::numeric,2) as avg_listbasedcmds,
  round(sum::numeric,2) as sum_listbasedcmds
from
  aws_elasticache_redis_metric_list_based_cmds_hourly
where sum > 100
order by
  cache_cluster_id,
  timestamp;
```

```sql+sqlite
select
  cache_cluster_id,
  timestamp,
  round(minimum,2) as min_listbasedcmds,
  round(maximum,2) as max_listbasedcmds,
  round(average,2) as avg_listbasedcmds,
  round(sum,2) as sum_listbasedcmds
from
  aws_elasticache_redis_metric_list_based_cmds_hourly
where sum > 100
order by
  cache_cluster_id,
  timestamp;
```