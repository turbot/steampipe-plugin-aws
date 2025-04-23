---
title: "Steampipe Table: aws_elasticache_redis_metric_engine_cpu_utilization_daily - Query AWS ElastiCache Redis Metrics using SQL"
description: "Allows users to query ElastiCache Redis Metrics and provides daily statistics for Engine CPU Utilization."
folder: "ElastiCache"
---

# Table: aws_elasticache_redis_metric_engine_cpu_utilization_daily - Query AWS ElastiCache Redis Metrics using SQL

The AWS ElastiCache Redis Metrics service is a tool that allows you to collect, track, and analyze performance metrics for your running ElastiCache instances. It provides valuable information about CPU utilization, helping you understand how your applications are using your cache and where bottlenecks are occurring. This data can help you make informed decisions about scaling and optimizing your ElastiCache instances for better application performance.

## Table Usage Guide

The `aws_elasticache_redis_metric_engine_cpu_utilization_daily` table in Steampipe provides you with daily statistical data about the CPU utilization of an Amazon ElastiCache Redis engine. This table allows you, as a DevOps engineer or data analyst, to query and analyze the CPU usage patterns of your ElastiCache Redis instances. This enables you to identify potential performance bottlenecks and optimize resource allocation. The schema outlines the various attributes of the CPU utilization metrics for you, including the timestamp, minimum, maximum, and average CPU usage, as well as the standard deviation.

The `aws_elasticache_redis_metric_engine_cpu_utilization_daily` table provides you with metric statistics at 24-hour intervals for the last year.

## Examples

### Basic info
Analyze the daily CPU utilization of AWS ElastiCache Redis clusters to understand their performance trends and capacity planning. This allows you to identify instances where resource usage may be high and adjust accordingly to ensure optimal functioning.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which your AWS ElastiCache Redis instances are utilizing more than 80% of the CPU on average. This allows you to identify potential performance issues and optimize resource allocation.

```sql+postgres
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

```sql+sqlite
select
  cache_cluster_id,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_elasticache_redis_metric_engine_cpu_utilization_daily
where average > 80
order by
  cache_cluster_id,
  timestamp;
```

### CPU daily average < 2%
Identify instances where the daily average CPU utilization is less than 2% in your AWS ElastiCache Redis clusters. This is useful in understanding underutilized resources, which can help optimize costs and resource allocation.

```sql+postgres
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

```sql+sqlite
select
  cache_cluster_id,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_elasticache_redis_metric_engine_cpu_utilization_daily
where average < 2
order by
  cache_cluster_id,
  timestamp;
```