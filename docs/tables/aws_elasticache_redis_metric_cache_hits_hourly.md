---
title: "Table: aws_elasticache_redis_metric_cache_hits_hourly - Query Amazon ElastiCache Redis Cache Hits using SQL"
description: "Allows users to query Amazon ElastiCache Redis Cache Hits on an hourly basis."
---

# Table: aws_elasticache_redis_metric_cache_hits_hourly - Query Amazon ElastiCache Redis Cache Hits using SQL

The `aws_elasticache_redis_metric_cache_hits_hourly` table in Steampipe provides information about the cache hits metrics of Amazon ElastiCache Redis instances on an hourly basis. This table allows system administrators and DevOps engineers to monitor and analyze the performance of Redis cache nodes by querying the cache hits metrics. Users can utilize this table to gather insights on cache hits, such as the number of successful lookup of keys in the cache, and to understand the efficiency of their cache configurations. The schema outlines the various attributes of the cache hits metrics, including the timestamp, cache hits, dimensions, and more.

The `aws_elasticache_redis_metric_cache_hits_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_redis_metric_cache_hits_hourly` table, you can use the `.inspect aws_elasticache_redis_metric_cache_hits_hourly` command in Steampipe.

### Key columns:

- `timestamp`: This is the timestamp for the cache hits data point. It is useful for tracking the cache hits over time and can be used to join with tables that contain time-based information.
- `cache_hits`: This is the number of successful lookup of keys in the cache. It is crucial for understanding the efficiency of the cache and can be used to join with tables that contain cache performance data.
- `dimensions`: This contains information about the name and region of the cache. It can be used to join with other tables that contain cache details based on name and region.

## Examples

### Basic info

```sql
select
  cache_cluster_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_elasticache_redis_metric_cache_hits_hourly
order by
  cache_cluster_id,
  timestamp;
```

### CacheHit sum below 10 

```sql
select
  cache_cluster_id,
  timestamp,
  round(sum::numeric,2) as sum_cachehits,
  round(average::numeric,2) as average_cachehits,
  sample_count
from
  aws_elasticache_redis_metric_cache_hits_hourly
where sum < 10
order by
  cache_cluster_id,
  timestamp;
```

### CacheHit hourly average < 100

```sql
select
  cache_cluster_id,
  timestamp,
  round(minimum::numeric,2) as min_cachehits,
  round(maximum::numeric,2) as max_cachehits,
  round(average::numeric,2) as avg_cachehits,
  sample_count
from
  aws_elasticache_redis_metric_cache_hits_hourly
where average < 100
order by
  cache_cluster_id,
  timestamp;
```