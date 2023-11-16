---
title: "Table: aws_elasticache_redis_metric_list_based_cmds_hourly - Query AWS ElastiCache Redis Metrics using SQL"
description: "Allows users to query ElastiCache Redis Metrics on an hourly basis, providing data on list-based commands executed in the ElastiCache Redis environment."
---

# Table: aws_elasticache_redis_metric_list_based_cmds_hourly - Query AWS ElastiCache Redis Metrics using SQL

The `aws_elasticache_redis_metric_list_based_cmds_hourly` table in Steampipe provides information about list-based command metrics within AWS ElastiCache Redis. This table allows DevOps engineers to query command-specific details on an hourly basis, including the number of commands processed, the latency of commands, and associated metadata. Users can utilize this table to gather insights on command performance, such as identifying high latency commands, tracking the frequency of command usage, and more. The schema outlines the various attributes of the ElastiCache Redis command metrics, including the cache cluster id, the metric name, and the timestamp.

The `aws_elasticache_redis_metric_list_based_cmds_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_redis_metric_list_based_cmds_hourly` table, you can use the `.inspect aws_elasticache_redis_metric_list_based_cmds_hourly` command in Steampipe.

### Key columns:

- `cache_cluster_id`: This is the identifier for the cache cluster. This column is important as it can be used to join this table with other ElastiCache tables to get more detailed information about the specific cache cluster.

- `metric_name`: This column contains the name of the specific list-based command metric. This column is useful for filtering queries to target specific metrics.

- `timestamp`: This column holds the timestamp for the metric data point. It is useful for tracking changes over time and identifying trends or anomalies.

## Examples

### Basic info

```sql
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

```sql
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