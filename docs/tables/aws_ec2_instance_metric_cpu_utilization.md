---
title: "Table: aws_ec2_instance_metric_cpu_utilization - Query AWS EC2 Instance Metrics using SQL"
description: "Allows users to query EC2 Instance CPU Utilization metrics from AWS CloudWatch."
---

# Table: aws_ec2_instance_metric_cpu_utilization - Query AWS EC2 Instance Metrics using SQL

The `aws_ec2_instance_metric_cpu_utilization` table in Steampipe provides information about CPU utilization metrics of EC2 instances within AWS CloudWatch. This table allows DevOps engineers, system administrators, and other technical professionals to query CPU-specific details, including the instance's average, maximum, and minimum CPU utilization. Users can utilize this table to gather insights on instance performance, such as identifying instances with high CPU utilization, analyzing CPU usage patterns, and more. The schema outlines the various attributes of the EC2 instance CPU utilization metrics, including the instance ID, namespace, metric name, and statistics.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_instance_metric_cpu_utilization` table, you can use the `.inspect aws_ec2_instance_metric_cpu_utilization` command in Steampipe.

**Key columns**:

- `instance_id`: The ID of the instance. This is a key identifier and can be used to join this table with other EC2 instance tables.
- `namespace`: The namespace of the metric. This is useful for identifying the specific AWS service the metric is associated with.
- `metric_name`: The name of the metric. This is useful for identifying the specific type of utilization metric (in this case, CPU utilization).

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
  aws_ec2_instance_metric_cpu_utilization
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
  aws_ec2_instance_metric_cpu_utilization
where average > 80
order by
  instance_id,
  timestamp;
```