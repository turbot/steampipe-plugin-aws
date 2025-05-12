---
title: "Steampipe Table: aws_ecs_cluster_metric_cpu_utilization_daily - Query AWS ECS Cluster Metrics using SQL"
description: "Allows users to query AWS Elastic Container Service (ECS) Cluster Metrics, specifically CPU utilization on a daily basis."
folder: "ECS"
---

# Table: aws_ecs_cluster_metric_cpu_utilization_daily - Query AWS ECS Cluster Metrics using SQL

The AWS ECS Cluster is a logical grouping of tasks or services. It allows you to manage and scale a set of services or tasks together in an AWS environment. The CPU Utilization Metric provides data about the CPU usage of the services or tasks in the cluster, helping you monitor and optimize resource allocation on a daily basis.

## Table Usage Guide

The `aws_ecs_cluster_metric_cpu_utilization_daily` table in Steampipe provides you with information about CPU utilization metrics within AWS Elastic Container Service (ECS) clusters. This table allows you, as a DevOps engineer, to query CPU utilization details on a daily basis, including the average, maximum, and minimum utilization. You can utilize this table to monitor and analyze CPU usage trends, identify potential performance issues, and optimize resource allocation. The schema outlines the various attributes of the CPU utilization metric for you, including the timestamp, period, unit, and statistical values.

## Examples

### Basic info
Explore the daily CPU usage patterns across your AWS ECS clusters. This can help you understand resource utilization trends and plan for capacity adjustments.

```sql+postgres
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

```sql+sqlite
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
Explore instances where the average CPU utilization exceeds 80% in AWS ECS clusters. This can help in identifying potential performance issues and aid in capacity planning.

```sql+postgres
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

```sql+sqlite
select
  cluster_name,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
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
Identify instances where the average daily CPU utilization of AWS ECS clusters is less than 1%. This can help in understanding underutilized resources and potentially save costs by downsizing or eliminating unnecessary clusters.

```sql+postgres
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

```sql+sqlite
select
  cluster_name,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_ecs_cluster_metric_cpu_utilization_daily
where
  average < 1
order by
  cluster_name,
  timestamp;
```