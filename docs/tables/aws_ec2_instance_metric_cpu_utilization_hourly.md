---
title: "Steampipe Table: aws_ec2_instance_metric_cpu_utilization_hourly - Query AWS EC2 Instance Metrics using SQL"
description: "Allows users to query AWS EC2 Instance CPU Utilization metrics on an hourly basis."
folder: "EC2"
---

# Table: aws_ec2_instance_metric_cpu_utilization_hourly - Query AWS EC2 Instance Metrics using SQL

The AWS EC2 Instance Metrics service provides insights into the performance of your EC2 instances. It allows you to monitor CPU utilization in an hourly manner using SQL queries. This can assist in identifying performance bottlenecks and optimizing resource usage for your EC2 instances.

## Table Usage Guide

The `aws_ec2_instance_metric_cpu_utilization_hourly` table in Steampipe provides you with information about the CPU Utilization metrics of EC2 instances in AWS. This table enables you as a DevOps engineer, system administrator, or other technical professional to query CPU utilization metrics on an hourly basis. This can be useful for you in monitoring system performance, identifying potential bottlenecks, and planning for capacity. The schema outlines the various attributes of the EC2 instance CPU utilization metrics for you, including the instance ID, timestamp, maximum, minimum, and average CPU utilization.

The `aws_ec2_instance_metric_cpu_utilization_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Explore which instances in your AWS EC2 service are experiencing fluctuating CPU utilization over time. This allows you to pinpoint specific locations where performance optimization may be needed to improve overall efficiency.

```sql+postgres
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

```sql+sqlite
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
Identify instances where EC2 instances have an average CPU utilization exceeding 80%. This can help in monitoring and optimizing resource usage, ensuring efficient performance of your AWS infrastructure.

```sql+postgres
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

```sql+sqlite
select
  instance_id,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization_hourly
where average > 80
order by
  instance_id,
  timestamp;
```

### CPU hourly average < 1%
Determine the areas in which your AWS EC2 instances' CPU utilization is less than 1% on average per hour. This can help you identify underutilized resources and optimize your AWS usage for cost-effectiveness.

```sql+postgres
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

```sql+sqlite
select
  instance_id,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization_hourly
where average < 1
order by
  instance_id,
  timestamp;
```