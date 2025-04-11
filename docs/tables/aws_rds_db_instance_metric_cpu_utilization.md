---
title: "Steampipe Table: aws_rds_db_instance_metric_cpu_utilization - Query Amazon RDS DBInstanceCPUUtilization using SQL"
description: "Allows users to query Amazon RDS DBInstanceCPUUtilization to fetch data about CPU utilization metrics for RDS DB instances."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_cpu_utilization - Query Amazon RDS DBInstanceCPUUtilization using SQL

The Amazon RDS DBInstanceCPUUtilization is a resource that provides metrics about the percentage of CPU utilization for an Amazon RDS instance. This metric allows you to monitor the compute load on your DB instance. By using SQL queries, you can retrieve and analyze these metrics to optimize the performance and resource usage of your Amazon RDS instances.

## Table Usage Guide

The `aws_rds_db_instance_metric_cpu_utilization` table in Steampipe provides you with information about CPU utilization metrics for Amazon RDS DB instances. This table allows you, as a DevOps engineer, database administrator, or other technical professional, to query CPU utilization-specific details, including average, maximum, and minimum CPU utilization, along with associated timestamps. You can utilize this table to monitor and analyze the CPU usage of your RDS DB instances, helping you optimize resource usage, identify potential performance bottlenecks, and make informed scaling decisions. The schema outlines the various attributes of the CPU utilization metrics, including the DB instance identifier, period, unit, and statistical data.

The `aws_rds_db_instance_metric_cpu_utilization` table provides you with metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Explore the CPU utilization metrics of your AWS RDS database instances over time. This can help you understand usage patterns, identify potential performance issues, and plan for capacity management.

```sql+postgres
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization
order by
  db_instance_identifier,
  timestamp;
```

```sql+sqlite
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization
order by
  db_instance_identifier,
  timestamp;
```

### CPU Over 80% average
Identify instances where the average CPU utilization of an AWS RDS database instance exceeds 80%. This can aid in assessing the performance of your databases and pinpointing any potential areas of concern that may require optimization or resource scaling.

```sql+postgres
select
  db_instance_identifier,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization
where average > 80
order by
  db_instance_identifier,
  timestamp;
```

```sql+sqlite
select
  db_instance_identifier,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization
where average > 80
order by
  db_instance_identifier,
  timestamp;
```