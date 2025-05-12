---
title: "Steampipe Table: aws_ec2_instance_metric_cpu_utilization - Query AWS EC2 Instance Metrics using SQL"
description: "Allows users to query EC2 Instance CPU Utilization metrics from AWS CloudWatch."
folder: "EC2"
---

# Table: aws_ec2_instance_metric_cpu_utilization - Query AWS EC2 Instance Metrics using SQL

The AWS EC2 Instance Metrics is a feature of Amazon EC2 (Elastic Compute Cloud) that provides detailed reports on the performance of your EC2 instances. These metrics include CPU utilization, which measures the percentage of total CPU time spent on various tasks within the EC2 instance. By querying these metrics using SQL, you can gain insights into your instance's performance and optimize resource usage.

## Table Usage Guide

The `aws_ec2_instance_metric_cpu_utilization` table in Steampipe provides you with information about CPU utilization metrics of EC2 instances within AWS CloudWatch. This table allows you, as a DevOps engineer, system administrator, or other technical professional, to query CPU-specific details, including the instance's average, maximum, and minimum CPU utilization. You can utilize this table to gather insights on instance performance, such as identifying instances with high CPU utilization, analyzing CPU usage patterns, and more. The schema outlines the various attributes of the EC2 instance CPU utilization metrics for you, including the instance ID, namespace, metric name, and statistics.

## Examples

### Basic info
Explore which AWS EC2 instances have varying CPU utilization levels and when these fluctuations occur. This information can help identify instances that may require optimization for improved performance and cost efficiency.

```sql+postgres
select
  instance_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization
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
  aws_ec2_instance_metric_cpu_utilization
order by
  instance_id,
  timestamp;
```

### CPU Over 80% average
Determine the areas in which instances of your AWS EC2 service are experiencing high CPU utilization, specifically where the average CPU usage exceeds 80%. This can help in identifying potential performance issues and optimize resource allocation.

```sql+postgres
select
  instance_id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization
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
  aws_ec2_instance_metric_cpu_utilization
where average > 80
order by
  instance_id,
  timestamp;
```