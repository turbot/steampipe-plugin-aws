---
title: "Table: aws_elasticache_redis_metric_engine_cpu_utilization_hourly - Query AWS ElastiCache Redis using SQL"
description: "Allows users to query hourly CPU utilization metrics for AWS ElastiCache Redis."
---

# Table: aws_elasticache_redis_metric_engine_cpu_utilization_hourly - Query AWS ElastiCache Redis using SQL

The `aws_elasticache_redis_metric_engine_cpu_utilization_hourly` table in Steampipe provides information about the hourly CPU utilization metrics for AWS ElastiCache Redis. This table allows DevOps engineers, database administrators, and other technical professionals to query time-series data related to CPU usage, thereby enabling them to monitor performance, identify potential bottlenecks, and optimize resource allocation. The schema outlines various attributes of the CPU utilization metrics, including the timestamp, average, maximum, and minimum CPU utilization, among others.

The `aws_elasticache_redis_metric_engine_cpu_utilization_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_redis_metric_engine_cpu_utilization_hourly` table, you can use the `.inspect aws_elasticache_redis_metric_engine_cpu_utilization_hourly` command in Steampipe.

**Key columns**:

- `title_id`: This is the unique identifier for the metric, which can be used to join this table with other tables that contain metric-specific information.
- `average`: This column holds the average CPU utilization for the hour. It is crucial for understanding the typical load on the CPU over time.
- `timestamp`: This column records the time at which the metric was collected. It is essential for tracking CPU usage trends and identifying peak usage periods.

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
  aws_elasticache_redis_metric_engine_cpu_utilization_hourly
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
  aws_elasticache_redis_metric_engine_cpu_utilization_hourly
where average > 80
order by
  cache_cluster_id,
  timestamp;
```

### CPU hourly average < 2%

```sql
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