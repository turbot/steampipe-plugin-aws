---
title: "Table: aws_rds_db_instance_metric_cpu_utilization - Query Amazon RDS DBInstanceCPUUtilization using SQL"
description: "Allows users to query Amazon RDS DBInstanceCPUUtilization to fetch data about CPU utilization metrics for RDS DB instances."
---

# Table: aws_rds_db_instance_metric_cpu_utilization - Query Amazon RDS DBInstanceCPUUtilization using SQL

The `aws_rds_db_instance_metric_cpu_utilization` table in Steampipe provides information about CPU utilization metrics for Amazon RDS DB instances. This table allows DevOps engineers, database administrators, and other technical professionals to query CPU utilization-specific details, including average, maximum, and minimum CPU utilization, along with associated timestamps. Users can utilize this table to monitor and analyze the CPU usage of their RDS DB instances, helping them optimize resource usage, identify potential performance bottlenecks, and make informed scaling decisions. The schema outlines the various attributes of the CPU utilization metrics, including the DB instance identifier, period, unit, and statistical data.

The `aws_rds_db_instance_metric_cpu_utilization` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_metric_cpu_utilization` table, you can use the `.inspect aws_rds_db_instance_metric_cpu_utilization` command in Steampipe.

### Key columns:

- `db_instance_identifier`: The identifier for the DB instance. This column is important as it uniquely identifies the DB instance and can be used to join this table with other tables related to RDS DB instances.
- `timestamp`: The timestamp for the data point. This column is useful for tracking CPU utilization over time and identifying patterns or anomalies.
- `average`: The average of the metric values that correspond to the specified time period. This column is crucial for understanding the typical CPU utilization of the DB instance.

## Examples

### Basic info

```sql
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

```sql
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