---
title: "Table: aws_elasticache_redis_metric_new_connections_hourly - Query AWS ElastiCache Redis Metrics using SQL"
description: "Allows users to query AWS ElastiCache Redis Metrics to get hourly data on new connections."
---

# Table: aws_elasticache_redis_metric_new_connections_hourly - Query AWS ElastiCache Redis Metrics using SQL

The `aws_elasticache_redis_metric_new_connections_hourly` table in Steampipe provides information about AWS ElastiCache Redis Metrics. This table allows DevOps engineers and system administrators to query hourly data about new connections to their AWS ElastiCache Redis instances. Users can utilize this table to monitor connection trends, analyze system performance, and identify potential issues. The schema outlines the various attributes of the metrics, including the cache node ID, timestamp, maximum number of connections, and more.

The `aws_elasticache_redis_metric_new_connections_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_redis_metric_new_connections_hourly` table, you can use the `.inspect aws_elasticache_redis_metric_new_connections_hourly` command in Steampipe.

**Key columns**:

- `cache_node_id`: This is the identifier of the cache node. It can be used to join this table with other tables that contain information about the cache node.
- `timestamp`: This is the timestamp for the data point. It allows users to query data based on specific timeframes.
- `maximum`: This is the maximum number of new connections within the given hour. It can be used to identify peak usage times.

## Examples

### Basic info

```sql
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

```sql
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