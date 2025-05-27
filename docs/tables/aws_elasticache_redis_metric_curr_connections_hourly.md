---
title: "Steampipe Table: aws_elasticache_redis_metric_curr_connections_hourly - Query AWS ElastiCache Redis using SQL"
description: "Allows users to query ElastiCache Redis current connections metrics on an hourly basis."
folder: "ElastiCache"
---

# Table: aws_elasticache_redis_metric_curr_connections_hourly - Query AWS ElastiCache Redis using SQL

The AWS ElastiCache Redis service provides a fully managed in-memory data store, compatible with Redis or Memcached. It improves the performance of web applications by retrieving data from fast, managed, in-memory data stores, instead of relying on slower disk-based databases. ElastiCache Redis supports data structures such as strings, hashes, lists, sets, sorted sets with range queries, bitmaps, hyperloglogs, geospatial indexes with radius queries and streams.

## Table Usage Guide

The `aws_elasticache_redis_metric_curr_connections_hourly` table in Steampipe provides you with information about the hourly current connections metrics of ElastiCache Redis within AWS. This table allows you, as a DevOps engineer, database administrator, or other technical professional, to query the current number of client connections, excluding connections from read replicas, to a Redis instance. You can utilize this table to monitor usage patterns, detect possible connection leaks, and optimize resource allocation based on connection demands. The schema outlines the various attributes of the ElastiCache Redis current connections metrics for you, including the timestamp, average, maximum, minimum, and sample count.

The `aws_elasticache_redis_metric_curr_connections_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Explore which AWS ElastiCache Redis clusters have the most connections over time. This information can help you understand the load on your clusters and identify any unusual spikes in connections that could indicate a problem.

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
  aws_elasticache_redis_metric_curr_connections_hourly
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
  aws_elasticache_redis_metric_curr_connections_hourly
order by
  cache_cluster_id,
  timestamp;
```

### currconnections Over 100 average
Explore the performance of your AWS ElastiCache Redis clusters by identifying instances where the average number of connections exceeds 100 in an hour. This can help in understanding the load on your clusters and take necessary actions if they are consistently over-utilized.

```sql+postgres
select
  cache_cluster_id,
  timestamp,
  round(minimum::numeric,2) as min_currconnections,
  round(maximum::numeric,2) as max_currconnections,
  round(average::numeric,2) as avg_currconnections,
  sample_count
from
  aws_elasticache_redis_metric_curr_connections_hourly
where average > 100
order by
  cache_cluster_id,
  timestamp;
```

```sql+sqlite
select
  cache_cluster_id,
  timestamp,
  round(minimum,2) as min_currconnections,
  round(maximum,2) as max_currconnections,
  round(average,2) as avg_currconnections,
  sample_count
from
  aws_elasticache_redis_metric_curr_connections_hourly
where average > 100
order by
  cache_cluster_id,
  timestamp;
```