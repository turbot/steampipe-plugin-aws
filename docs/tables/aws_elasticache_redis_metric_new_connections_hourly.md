---
title: "Steampipe Table: aws_elasticache_redis_metric_new_connections_hourly - Query AWS ElastiCache Redis Metrics using SQL"
description: "Allows users to query AWS ElastiCache Redis Metrics to get hourly data on new connections."
folder: "ElastiCache"
---

# Table: aws_elasticache_redis_metric_new_connections_hourly - Query AWS ElastiCache Redis Metrics using SQL

The AWS ElastiCache Redis Metrics provides a robust monitoring solution for your applications. It allows you to collect, view, and analyze metrics for your ElastiCache Redis instances through SQL queries. The 'new_connections_hourly' metric specifically measures the number of new connections made to the Redis server per hour, aiding in capacity planning and performance tuning.

## Table Usage Guide

The `aws_elasticache_redis_metric_new_connections_hourly` table in Steampipe provides you with information about AWS ElastiCache Redis Metrics. This table allows you, as a DevOps engineer or system administrator, to query hourly data about new connections to your AWS ElastiCache Redis instances. You can utilize this table to monitor connection trends, analyze system performance, and identify potential issues. The schema outlines the various attributes of the metrics, including the cache node ID, timestamp, maximum number of connections, and more.

The `aws_elasticache_redis_metric_new_connections_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Determine the areas in which AWS ElastiCache Redis clusters have experienced new connections over time. This can help in understanding usage patterns and identifying potential periods of high demand or unusual activity.

```sql+postgres
select
  cache_cluster_id,
  timestamp,
  minimum,
  maximum,
  average
from
  aws_elasticache_redis_metric_new_connections_hourly
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
  average
from
  aws_elasticache_redis_metric_new_connections_hourly
order by
  cache_cluster_id,
  timestamp;
```

### newconnections sum over 10
This query is useful for identifying instances where the total number of new connections to your AWS ElastiCache Redis clusters exceeds 10 within an hour. It allows you to monitor and manage your connection usage, helping to ensure optimal performance and avoid potential overloads.

```sql+postgres
select
  cache_cluster_id,
  timestamp,
  round(minimum::numeric,2) as min_newconnections,
  round(maximum::numeric,2) as max_newconnections,
  round(average::numeric,2) as avg_newconnections,
  round(sum::numeric,2) as sum_newconnections
from
  aws_elasticache_redis_metric_new_connections_hourly
where sum > 10
order by
  cache_cluster_id,
  timestamp;
```

```sql+sqlite
select
  cache_cluster_id,
  timestamp,
  round(minimum,2) as min_newconnections,
  round(maximum,2) as max_newconnections,
  round(average,2) as avg_newconnections,
  round(sum,2) as sum_newconnections
from
  aws_elasticache_redis_metric_new_connections_hourly
where sum > 10
order by
  cache_cluster_id,
  timestamp;
```