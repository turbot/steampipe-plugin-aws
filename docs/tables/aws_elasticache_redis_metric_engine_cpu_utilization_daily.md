---
title: "Table: aws_elasticache_redis_metric_engine_cpu_utilization_daily - Query AWS ElastiCache Redis Metrics using SQL"
description: "Allows users to query ElastiCache Redis Metrics and provides daily statistics for Engine CPU Utilization."
---

# Table: aws_elasticache_redis_metric_engine_cpu_utilization_daily - Query AWS ElastiCache Redis Metrics using SQL

The `aws_elasticache_redis_metric_engine_cpu_utilization_daily` table in Steampipe provides daily statistical data about the CPU utilization of an Amazon ElastiCache Redis engine. This table allows DevOps engineers and data analysts to query and analyze the CPU usage patterns of their ElastiCache Redis instances, enabling them to identify potential performance bottlenecks and optimize resource allocation. The schema outlines the various attributes of the CPU utilization metrics, including the timestamp, minimum, maximum, and average CPU usage, as well as the standard deviation.

The `aws_elasticache_redis_metric_engine_cpu_utilization_daily` table provides metric statistics at 24 hour intervals for the last year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_redis_metric_engine_cpu_utilization_daily` table, you can use the `.inspect aws_elasticache_redis_metric_engine_cpu_utilization_daily` command in Steampipe.

Key columns:

- `title`: The title of the metric, which is 'EngineCPUUtilization'. This column is important as it confirms the type of metric being queried.
- `timestamp`: The timestamp for the metric data point. This column is useful for tracking CPU usage over time and identifying usage patterns.
- `average`: The average CPU utilization for the given day. This column provides a general idea of the CPU load on the ElastiCache Redis engine for that day.

## Examples

### Basic info

```sql
select
  cache_cluster_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_elasticache_redis_metric_engine_cpu_utilization_daily
order by
  cache_cluster_id,
  timestamp;
```

### CPU Over 80% average

```sql
select
  cache_cluster_id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_elasticache_redis_metric_engine_cpu_utilization_daily
where average > 80
order by
  cache_cluster_id,
  timestamp;
```

### CPU daily average < 2%

```sql
select
  cache_cluster_id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_elasticache_redis_metric_engine_cpu_utilization_daily
where average < 2
order by
  cache_cluster_id,
  timestamp;
```
