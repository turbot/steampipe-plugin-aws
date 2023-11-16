---
title: "Table: aws_ecs_cluster_metric_cpu_utilization_hourly - Query AWS ECS Cluster Metrics using SQL"
description: "Allows users to query AWS ECS Cluster CPU Utilization Metrics on an hourly basis."
---

# Table: aws_ecs_cluster_metric_cpu_utilization_hourly - Query AWS ECS Cluster Metrics using SQL

The `aws_ecs_cluster_metric_cpu_utilization_hourly` table in Steampipe provides information about the CPU utilization metrics of AWS ECS (Elastic Container Service) clusters on an hourly basis. This table allows DevOps engineers, data analysts, and other technical professionals to query cluster-specific details, including the average, maximum, and minimum CPU utilization percentages. Users can utilize this table to monitor the performance of their ECS clusters, identify potential resource bottlenecks, and optimize resource allocation. The schema outlines the various attributes of the ECS cluster CPU utilization, including the cluster name, timestamp, average utilization, maximum utilization, and minimum utilization.

The `aws_ecs_cluster_metric_cpu_utilization_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecs_cluster_metric_cpu_utilization_hourly` table, you can use the `.inspect aws_ecs_cluster_metric_cpu_utilization_hourly` command in Steampipe.

### Key columns:

- `cluster_name`: The name of the ECS cluster. This is a key column because it can be used to join this table with other ECS cluster-specific tables.
- `timestamp`: The timestamp for the specific hour when the CPU utilization data was collected. This column is crucial for trend analysis and correlating CPU utilization metrics with other events.
- `average`: The average CPU utilization percentage for the specific hour. This column is important for understanding the typical CPU load on the ECS cluster during that hour.

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
  aws_ecs_cluster_metric_cpu_utilization_hourly
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
  aws_ecs_cluster_metric_cpu_utilization_hourly
where
  average > 80
order by
  cluster_name,
  timestamp;
```

### CPU hourly average < 1%

```sql
select
  cluster_name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_ecs_cluster_metric_cpu_utilization_hourly
where
  average < 1
order by
  cluster_name,
  timestamp;
```
