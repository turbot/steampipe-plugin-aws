---
title: "Table: aws_ecs_cluster_metric_cpu_utilization_daily - Query AWS ECS Cluster Metrics using SQL"
description: "Allows users to query AWS Elastic Container Service (ECS) Cluster Metrics, specifically CPU utilization on a daily basis."
---

# Table: aws_ecs_cluster_metric_cpu_utilization_daily - Query AWS ECS Cluster Metrics using SQL

The `aws_ecs_cluster_metric_cpu_utilization_daily` table in Steampipe provides information about CPU utilization metrics within AWS Elastic Container Service (ECS) clusters. This table allows DevOps engineers to query CPU utilization details on a daily basis, including the average, maximum, and minimum utilization. Users can utilize this table to monitor and analyze CPU usage trends, identify potential performance issues, and optimize resource allocation. The schema outlines the various attributes of the CPU utilization metric, including the timestamp, period, unit, and statistical values.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecs_cluster_metric_cpu_utilization_daily` table, you can use the `.inspect aws_ecs_cluster_metric_cpu_utilization_daily` command in Steampipe.

### Key columns:

- `timestamp`: This column records the time of each data point. It is crucial for tracking CPU usage over time and identifying usage patterns.
- `average`: This column shows the average CPU utilization for each day. It is useful for understanding the typical CPU usage.
- `maximum`: This column displays the maximum CPU utilization for each day. It aids in identifying peak usage times and potential performance issues.

## Examples

### Basic info

```sql
select
  cluster_name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_ecs_cluster_metric_cpu_utilization_daily
order by
  cluster_name,
  timestamp;
```

### CPU Over 80% average

```sql
select
  cluster_name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_ecs_cluster_metric_cpu_utilization_daily
where
  average > 80
order by
  cluster_name,
  timestamp;
```

### CPU daily average < 1%

```sql
select
  cluster_name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_ecs_cluster_metric_cpu_utilization_daily
where
  average < 1
order by
  cluster_name,
  timestamp;
```
