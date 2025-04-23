---
title: "Steampipe Table: aws_elasticache_redis_metric_engine_cpu_utilization_hourly - Query AWS ElastiCache Redis using SQL"
description: "Allows users to query hourly CPU utilization metrics for AWS ElastiCache Redis."
folder: "ElastiCache"
---

# Table: aws_elasticache_redis_metric_engine_cpu_utilization_hourly - Query AWS ElastiCache Redis using SQL

The AWS ElastiCache Redis is a web service that makes it easy to deploy, operate, and scale an in-memory data store or cache in the cloud. It provides a high-performance, scalable, and cost-effective caching solution, while removing the complexity associated with managing a distributed cache environment. This service is primarily used to improve the performance of web applications by retrieving information from fast, managed, in-memory caches, instead of relying entirely on slower disk-based databases.

## Table Usage Guide

The `aws_elasticache_redis_metric_engine_cpu_utilization_hourly` table in Steampipe gives you information about the hourly CPU utilization metrics for AWS ElastiCache Redis. This table enables you, as a DevOps engineer, database administrator, or other technical professional, to query time-series data related to CPU usage. As a result, you can monitor performance, identify potential bottlenecks, and optimize resource allocation. The schema outlines various attributes of the CPU utilization metrics for you, including the timestamp, average, maximum, and minimum CPU utilization, among others.

The `aws_elasticache_redis_metric_engine_cpu_utilization_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Explore the performance of your AWS ElastiCache Redis instances by analyzing CPU utilization over time. This can help optimize resource allocation and identify instances where performance tuning may be required.

```sql+postgres
select
  cache_cluster_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_elasticache_redis_metric_engine_cpu_utilization_hourly
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
  aws_elasticache_redis_metric_engine_cpu_utilization_hourly
order by
  cache_cluster_id,
  timestamp;
```

### CPU Over 80% average
Discover instances where your AWS ElastiCache Redis clusters are experiencing high CPU usage, specifically over 80% on average. This can help identify potential performance issues and allow for proactive troubleshooting.

```sql+postgres
select
  cache_cluster_id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_elasticache_redis_metric_engine_cpu_utilization_hourly
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
  aws_elasticache_redis_metric_engine_cpu_utilization_hourly
where average > 80
order by
  cache_cluster_id,
  timestamp;
```

### CPU hourly average < 2%
Analyze the performance of your AWS ElastiCache Redis clusters by identifying instances where the average CPU usage is less than 2% on an hourly basis. This can help pinpoint potential inefficiencies or underutilized resources, optimizing your cloud infrastructure management and cost efficiency.

```sql+postgres
select
  cache_cluster_id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_elasticache_redis_metric_engine_cpu_utilization_hourly
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
  aws_elasticache_redis_metric_engine_cpu_utilization_hourly
where average < 2
order by
  cache_cluster_id,
  timestamp;
```