---
title: "Steampipe Table: aws_elasticache_redis_metric_get_type_cmds_hourly - Query AWS ElastiCache Redis Metrics using SQL"
description: "Allows users to query ElastiCache Redis Metrics on an hourly basis. This includes information on GET type commands executed in the selected ElastiCache Redis cluster during the last hour."
folder: "ElastiCache"
---

# Table: aws_elasticache_redis_metric_get_type_cmds_hourly - Query AWS ElastiCache Redis Metrics using SQL

The AWS ElastiCache Redis Metrics service provides valuable insights into the performance of your Redis data stores. It allows you to monitor key performance metrics, including the number of 'get type' commands executed per hour. These metrics can help you optimize the performance and efficiency of your Redis data stores.

## Table Usage Guide

The `aws_elasticache_redis_metric_get_type_cmds_hourly` table in Steampipe provides you with information about the GET type commands executed in your selected AWS ElastiCache Redis cluster during the last hour. This table allows you, whether you're a DevOps engineer, database administrator, or other IT professional, to query and analyze the hourly GET type command metrics. This gives you insights into the performance and usage patterns of your ElastiCache Redis clusters. The schema outlines the various attributes of the ElastiCache Redis Metrics for you, including the average, maximum, minimum, sample count, and sum of GET type commands.

Your `aws_elasticache_redis_metric_get_type_cmds_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Analyze the performance of your AWS ElastiCache Redis clusters over time to ensure optimal resource utilization and response times. This practical application allows you to monitor and manage your clusters effectively, leading to improved performance and cost efficiency.

```sql+postgres
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

```sql+sqlite
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
Explore the performance of your AWS ElastiCache Redis clusters by identifying instances where the sum of 'get type' commands exceeds 100 in an hour. This can help in understanding usage patterns and planning for capacity upgrades or optimizations.

```sql+postgres
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

```sql+sqlite
select
  cache_cluster_id,
  timestamp,
  round(minimum,2) as min_gettypecmds,
  round(maximum,2) as max_gettypecmds,
  round(average,2) as avg_gettypecmds,
  round(sum,2) as sum_gettypecmds
from
  aws_elasticache_redis_metric_get_type_cmds_hourly
where sum > 100
order by
  cache_cluster_id,
  timestamp;
```