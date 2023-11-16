---
title: "Table: aws_elasticache_redis_metric_get_type_cmds_hourly - Query AWS ElastiCache Redis Metrics using SQL"
description: "Allows users to query ElastiCache Redis Metrics on an hourly basis. This includes information on GET type commands executed in the selected ElastiCache Redis cluster during the last hour."
---

# Table: aws_elasticache_redis_metric_get_type_cmds_hourly - Query AWS ElastiCache Redis Metrics using SQL

The `aws_elasticache_redis_metric_get_type_cmds_hourly` table in Steampipe provides information about the GET type commands executed in the selected AWS ElastiCache Redis cluster during the last hour. This table allows DevOps engineers, database administrators, and other IT professionals to query and analyze the hourly GET type command metrics, providing insights into the performance and usage patterns of the ElastiCache Redis clusters. The schema outlines the various attributes of the ElastiCache Redis Metrics, including the average, maximum, minimum, sample count, and sum of GET type commands.

The `aws_elasticache_redis_metric_get_type_cmds_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_redis_metric_get_type_cmds_hourly` table, you can use the `.inspect aws_elasticache_redis_metric_get_type_cmds_hourly` command in Steampipe.

Key columns:

- `title`: The title of the metric. This column can be used to join this table with other tables that contain metric details.
- `average`: The average number of GET type commands executed per hour. This column is useful for identifying usage patterns and performance analysis.
- `timestamp`: The timestamp of the metric. This column can be used to join this table with other tables that contain time-series data.

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
  aws_elasticache_redis_metric_get_type_cmds_hourly
order by
  cache_cluster_id,
  timestamp;
```

### gettypecmds sum 0ver 100 

```sql
select
  cache_cluster_id,
  timestamp,
  round(minimum::numeric,2) as min_gettypecmds,
  round(maximum::numeric,2) as max_gettypecmds,
  round(average::numeric,2) as avg_gettypecmds,
  round(sum::numeric,2) as sum_gettypecmds
from
  aws_elasticache_redis_metric_get_type_cmds_hourly
where sum > 100
order by
  cache_cluster_id,
  timestamp;
```
