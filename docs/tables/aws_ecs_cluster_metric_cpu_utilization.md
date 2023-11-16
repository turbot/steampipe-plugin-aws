---
title: "Table: aws_ecs_cluster_metric_cpu_utilization - Query AWS ECS Cluster Metrics using SQL"
description: "Allows users to query ECS Cluster CPU Utilization Metrics for a specified period."
---

# Table: aws_ecs_cluster_metric_cpu_utilization - Query AWS ECS Cluster Metrics using SQL

The `aws_ecs_cluster_metric_cpu_utilization` table in Steampipe provides information about CPU utilization metrics of AWS Elastic Container Service (ECS) clusters. This table allows DevOps engineers, system administrators, and other technical professionals to query CPU utilization-specific details, including the average, maximum, and minimum CPU utilization, along with the corresponding timestamps. Users can utilize this table to monitor CPU usage trends, identify potential performance issues, and optimize resource allocation. The schema outlines various attributes of the CPU utilization metric, including the cluster name, period, timestamp, and average, minimum, and maximum CPU utilization.

The `aws_ecs_cluster_metric_cpu_utilization` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecs_cluster_metric_cpu_utilization` table, you can use the `.inspect aws_ecs_cluster_metric_cpu_utilization` command in Steampipe.

### Key columns:

- `cluster_name`: This is the name of the ECS cluster. It can be used to join this table with other ECS related tables.
- `timestamp`: This is the timestamp for the corresponding CPU utilization metric. It is useful for tracking CPU usage over time.
- `average`: This represents the average CPU utilization for the specified period. It is essential for understanding the typical CPU usage of the ECS cluster.

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
  aws_ecs_cluster_metric_cpu_utilization
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
  aws_ecs_cluster_metric_cpu_utilization
where
  average > 80
order by
  cluster_name,
  timestamp;
```
