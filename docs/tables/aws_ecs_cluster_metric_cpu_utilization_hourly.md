---
title: "Steampipe Table: aws_ecs_cluster_metric_cpu_utilization_hourly - Query AWS ECS Cluster Metrics using SQL"
description: "Allows users to query AWS ECS Cluster CPU Utilization Metrics on an hourly basis."
folder: "ECS"
---

# Table: aws_ecs_cluster_metric_cpu_utilization_hourly - Query AWS ECS Cluster Metrics using SQL

The AWS ECS Cluster Metrics is a feature of Amazon Elastic Container Service (ECS) that provides CPU utilization data. It allows you to monitor and troubleshoot your applications running on ECS. The CPU Utilization metric represents the percentage of total CPU units that are currently in use on a cluster for an hour.

## Table Usage Guide

The `aws_ecs_cluster_metric_cpu_utilization_hourly` table in Steampipe gives you information about the CPU utilization metrics of AWS ECS (Elastic Container Service) clusters on an hourly basis. This table allows you, as a DevOps engineer, data analyst, or other technical professional, to query cluster-specific details, including the average, maximum, and minimum CPU utilization percentages. You can utilize this table to monitor the performance of your ECS clusters, identify potential resource bottlenecks, and optimize resource allocation. The schema outlines the various attributes of the ECS cluster CPU utilization for you, including the cluster name, timestamp, average utilization, maximum utilization, and minimum utilization.

The `aws_ecs_cluster_metric_cpu_utilization_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Explore the performance of various AWS ECS clusters by tracking their CPU utilization over time. This allows for effective resource management and helps in identifying potential performance issues.

```sql+postgres
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

```sql+sqlite
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
Discover the instances where the average CPU utilization exceeds 80% in your AWS ECS clusters, allowing you to identify potential performance issues and optimize resource allocation.

```sql+postgres
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

```sql+sqlite
select
  cluster_name,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
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
Determine the areas in which AWS ECS clusters are underutilized, by pinpointing instances where the average CPU usage is less than 1% on an hourly basis. This allows for efficient resource management and cost optimization by identifying potential opportunities for downsizing.

```sql+postgres
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

```sql+sqlite
select
  cluster_name,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_ecs_cluster_metric_cpu_utilization_hourly
where
  average < 1
order by
  cluster_name,
  timestamp;
```