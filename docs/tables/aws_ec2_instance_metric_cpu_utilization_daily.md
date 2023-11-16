---
title: "Table: aws_ec2_instance_metric_cpu_utilization_daily - Query AWS EC2 Instances using SQL"
description: "Allows users to query daily CPU utilization metrics of AWS EC2 instances."
---

# Table: aws_ec2_instance_metric_cpu_utilization_daily - Query AWS EC2 Instances using SQL

The `aws_ec2_instance_metric_cpu_utilization_daily` table in Steampipe provides information about the daily CPU utilization metrics of AWS EC2 instances. This table allows DevOps engineers to query instance-specific details, including average, maximum, and minimum CPU utilization, and associated timestamps. Users can utilize this table to gather insights on CPU usage patterns over time, such as instances with high or low CPU utilization, instances with abnormal CPU usage patterns, and more. The schema outlines the various attributes of the CPU utilization metrics, including the instance ID, timestamp, average CPU utilization, maximum CPU utilization, and minimum CPU utilization.

The `aws_ec2_instance_metric_cpu_utilization_daily` table provides metric statistics at 24 hour intervals for the last year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_instance_metric_cpu_utilization_daily` table, you can use the `.inspect aws_ec2_instance_metric_cpu_utilization_daily` command in Steampipe.

Key columns:

- `instance_id`: This is the ID of the EC2 instance. It is a primary key for the table and can be used to join this table with other EC2 related tables.
- `timestamp`: This is the timestamp of the CPU utilization metric. It can be used to join this table with other time-based tables or to perform time-based queries.
- `average`: This is the average CPU utilization for the day. It can be used to identify instances with high or low average CPU utilization over time.

## Examples

### Basic info

```sql
select
  instance_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization_daily
order by
  instance_id,
  timestamp;
```



### CPU Over 80% average

```sql
select
  instance_id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization_daily
where average > 80
order by
  instance_id,
  timestamp;
```

### CPU daily average < 1%

```sql
select
  instance_id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization_daily
where average < 1
order by
  instance_id,
  timestamp;
```