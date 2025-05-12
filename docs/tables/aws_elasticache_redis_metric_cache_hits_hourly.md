---
title: "Steampipe Table: aws_elasticache_redis_metric_cache_hits_hourly - Query Amazon ElastiCache Redis Cache Hits using SQL"
description: "Allows users to query Amazon ElastiCache Redis Cache Hits on an hourly basis."
folder: "ElastiCache"
---

# Table: aws_elasticache_redis_metric_cache_hits_hourly - Query Amazon ElastiCache Redis Cache Hits using SQL

The Amazon ElastiCache Redis is a web service that makes it easy to set up, manage, and scale a distributed in-memory data store or cache environment in the cloud. It provides a high-performance, scalable, and cost-effective caching solution, while removing the complexity associated with managing a distributed cache environment. The 'Cache Hits' metric specifically provides the number of successful read-only key lookups in the main dictionary on an hourly basis.

## Table Usage Guide

The `aws_elasticache_redis_metric_cache_hits_hourly` table in Steampipe provides you with information about the cache hits metrics of Amazon ElastiCache Redis instances on an hourly basis. This table allows you as a system administrator or a DevOps engineer to monitor and analyze the performance of Redis cache nodes by querying the cache hits metrics. You can utilize this table to gather insights on cache hits, such as the number of successful lookup of keys in the cache, and to understand the efficiency of your cache configurations. The schema outlines the various attributes of the cache hits metrics for you, including the timestamp, cache hits, dimensions, and more.

The `aws_elasticache_redis_metric_cache_hits_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Determine the efficiency of your AWS ElastiCache Redis instances by analyzing cache hit metrics over time. This can help optimize performance and resource utilization.

```sql+postgres
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

```sql+sqlite
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
The query is used to monitor the performance of your AWS ElastiCache Redis clusters by identifying instances where the sum of cache hits falls below 10 in an hour. This can help you pinpoint potential issues and optimize your cache usage for improved application performance.

```sql+postgres
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

```sql+sqlite
select
  cache_cluster_id,
  timestamp,
  round(sum,2) as sum_cachehits,
  round(average,2) as average_cachehits,
  sample_count
from
  aws_elasticache_redis_metric_cache_hits_hourly
where sum < 10
order by
  cache_cluster_id,
  timestamp;
```

### CacheHit hourly average < 100
Explore the performance of your AWS ElastiCache Redis clusters by identifying instances where the hourly average of cache hits is less than 100. This can help pinpoint potential areas of concern and optimize the usage of your cache clusters.

```sql+postgres
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

```sql+sqlite
select
  cache_cluster_id,
  timestamp,
  round(minimum,2) as min_cachehits,
  round(maximum,2) as max_cachehits,
  round(average,2) as avg_cachehits,
  sample_count
from
  aws_elasticache_redis_metric_cache_hits_hourly
where average < 100
order by
  cache_cluster_id,
  timestamp;
```