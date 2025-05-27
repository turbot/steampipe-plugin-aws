---
title: "Steampipe Table: aws_ec2_instance_metric_cpu_utilization_daily - Query AWS EC2 Instances using SQL"
description: "Allows users to query daily CPU utilization metrics of AWS EC2 instances."
folder: "EC2"
---

# Table: aws_ec2_instance_metric_cpu_utilization_daily - Query AWS EC2 Instances using SQL

The AWS EC2 Instance is a virtual server in Amazon's Elastic Compute Cloud (EC2) for running applications on the Amazon Web Services (AWS) infrastructure. It provides scalable computing capacity in the AWS Cloud, allowing developers to launch as many or as few virtual servers as needed. The CPU Utilization metric provides the percentage of CPU utilization for an EC2 instance, averaged over a daily period.

## Table Usage Guide

The `aws_ec2_instance_metric_cpu_utilization_daily` table in Steampipe provides you with information about the daily CPU utilization metrics of AWS EC2 instances. This table allows you, as a DevOps engineer, to query instance-specific details, including average, maximum, and minimum CPU utilization, and associated timestamps. You can utilize this table to gather insights on CPU usage patterns over time, such as instances with high or low CPU utilization, instances with abnormal CPU usage patterns, and more. The schema outlines the various attributes of the CPU utilization metrics for you, including the instance ID, timestamp, average CPU utilization, maximum CPU utilization, and minimum CPU utilization.

The `aws_ec2_instance_metric_cpu_utilization_daily` table provides you with metric statistics at 24 hour intervals for the last year.

## Examples

### Basic info
Determine the areas in which daily CPU utilization of AWS EC2 instances fluctuates, allowing for more effective resource management and cost optimization.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which your AWS EC2 instances are utilizing more than 80% of their CPU capacity on average. This can help in identifying potential performance bottlenecks and planning for capacity upgrades.

```sql+postgres
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

```sql+sqlite
select
  instance_id,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization_daily
where average > 80
order by
  instance_id,
  timestamp;
```

### CPU daily average < 1%
Determine the areas in which your AWS EC2 instances are underutilized, specifically where daily average CPU usage is less than 1%. This can help identify potential cost savings by downsizing or eliminating these underused resources.

```sql+postgres
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

```sql+sqlite
select
  instance_id,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization_daily
where average < 1
order by
  instance_id,
  timestamp;
```