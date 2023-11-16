---
title: "Table: aws_elasticache_redis_metric_curr_connections_hourly - Query AWS ElastiCache Redis using SQL"
description: "Allows users to query ElastiCache Redis current connections metrics on an hourly basis."
---

# Table: aws_elasticache_redis_metric_curr_connections_hourly - Query AWS ElastiCache Redis using SQL

The `aws_elasticache_redis_metric_curr_connections_hourly` table in Steampipe provides information about the hourly current connections metrics of ElastiCache Redis within AWS. This table allows DevOps engineers, database administrators, and other technical professionals to query the current number of client connections, excluding connections from read replicas, to a Redis instance. Users can utilize this table to monitor usage patterns, detect possible connection leaks, and optimize resource allocation based on connection demands. The schema outlines the various attributes of the ElastiCache Redis current connections metrics, including the timestamp, average, maximum, minimum, and sample count.

The `aws_elasticache_redis_metric_curr_connections_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_redis_metric_curr_connections_hourly` table, you can use the `.inspect aws_elasticache_redis_metric_curr_connections_hourly` command in Steampipe.

### Key columns:

* `timestamp`: This is the time when the metric data was received. It can be used to track the historical data of connections.
* `average`: This column represents the average number of client connections for a given hour. It can be used to analyze the average demand for connections over time.
* `maximum`: This column shows the maximum number of client connections for a given hour. It can be used to identify peak usage times and potential connection leaks.

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
  aws_elasticache_redis_metric_curr_connections_hourly
order by
  cache_cluster_id,
  timestamp;
```

### currconnections Over 100 average

```sql
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

