---
title: "Table: aws_ec2_instance_metric_cpu_utilization_hourly - Query AWS EC2 Instance Metrics using SQL"
description: "Allows users to query AWS EC2 Instance CPU Utilization metrics on an hourly basis."
---

# Table: aws_ec2_instance_metric_cpu_utilization_hourly - Query AWS EC2 Instance Metrics using SQL

The `aws_ec2_instance_metric_cpu_utilization_hourly` table in Steampipe provides information about the CPU Utilization metrics of EC2 instances in AWS. This table allows DevOps engineers, system administrators, and other technical professionals to query CPU utilization metrics on an hourly basis, which can be useful for monitoring system performance, identifying potential bottlenecks, and planning for capacity. The schema outlines the various attributes of the EC2 instance CPU utilization metrics, including the instance ID, timestamp, maximum, minimum, and average CPU utilization.

The `aws_ec2_instance_metric_cpu_utilization_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_instance_metric_cpu_utilization_hourly` table, you can use the `.inspect aws_ec2_instance_metric_cpu_utilization_hourly` command in Steampipe.

### Key columns:

- `instance_id`: This is the unique identifier of the EC2 instance. It is crucial for joining this table with other EC2-related tables.
- `timestamp`: This column represents the time when the CPU utilization metrics were recorded. It is useful for tracking CPU utilization over time.
- `average`: This column provides the average CPU utilization during the specified hour. It is helpful for identifying trends in CPU usage.

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
  aws_ec2_instance_metric_cpu_utilization_hourly
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
  aws_ec2_instance_metric_cpu_utilization_hourly
where average > 80
order by
  instance_id,
  timestamp;
```

### CPU hourly average < 1%

```sql
select
  instance_id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization_hourly
where average < 1
order by
  instance_id,
  timestamp;
```