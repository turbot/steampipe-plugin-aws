---
title: "Steampipe Table: aws_ecs_cluster_metric_cpu_utilization - Query AWS ECS Cluster Metrics using SQL"
description: "Allows users to query ECS Cluster CPU Utilization Metrics for a specified period."
folder: "ECS"
---

# Table: aws_ecs_cluster_metric_cpu_utilization - Query AWS ECS Cluster Metrics using SQL

The AWS ECS Cluster Metrics service allows you to monitor and collect data on CPU utilization in your Amazon Elastic Container Service (ECS) clusters. This feature provides insights into the efficiency of your clusters and can be used to optimize resource usage. You can query these metrics using SQL, allowing for easy integration and analysis of the data.

## Table Usage Guide

The `aws_ecs_cluster_metric_cpu_utilization` table in Steampipe provides you with information about CPU utilization metrics of AWS Elastic Container Service (ECS) clusters. This table allows you, as a DevOps engineer, system administrator, or other technical professional, to query CPU utilization-specific details, including the average, maximum, and minimum CPU utilization, along with the corresponding timestamps. You can utilize this table to monitor CPU usage trends, identify potential performance issues, and optimize resource allocation. The schema outlines various attributes of the CPU utilization metric, including the cluster name, period, timestamp, and average, minimum, and maximum CPU utilization.

The `aws_ecs_cluster_metric_cpu_utilization` table provides you with metric statistics at 5-minute intervals for the most recent 5 days.

## Examples

### Basic info
Analyze the CPU utilization metrics of AWS ECS clusters over time to understand resource usage trends and optimize cluster performance. This information could be useful in identifying patterns, planning capacity, and managing costs effectively.

```sql+postgres
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

```sql+sqlite
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
Identify instances where the average CPU utilization of your AWS ECS clusters exceeds 80%. This can help in managing resources effectively, ensuring optimal performance and avoiding potential bottlenecks.

```sql+postgres
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

```sql+sqlite
select
  cluster_name,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_ecs_cluster_metric_cpu_utilization
where
  average > 80
order by
  cluster_name,
  timestamp;
```