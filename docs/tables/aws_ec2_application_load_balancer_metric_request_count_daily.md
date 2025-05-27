---
title: "Steampipe Table: aws_ec2_application_load_balancer_metric_request_count_daily - Query AWS EC2 Application Load Balancer using SQL"
description: "Allows users to query daily request count metrics of the AWS EC2 Application Load Balancer."
folder: "EC2"
---

# Table: aws_ec2_application_load_balancer_metric_request_count_daily - Query AWS EC2 Application Load Balancer using SQL

The AWS EC2 Application Load Balancer is a fully managed service that operates at the application layer, the seventh layer of the Open Systems Interconnection (OSI) model. It performs advanced traffic distribution across multiple targets, such as Amazon EC2 instances, containers, and IP addresses. The service also monitors the health of its registered targets and ensures that it routes traffic only to healthy targets.

## Table Usage Guide

The `aws_ec2_application_load_balancer_metric_request_count_daily` table in Steampipe gives you information about the daily request count metrics of the AWS EC2 Application Load Balancer. You can use this table to query and analyze the number of requests processed by the Application Load Balancer on a daily basis. It allows you, as a DevOps engineer, to monitor the load on the balancer, identify potential spikes in traffic, and plan capacity accordingly. The schema outlines the various attributes of the metrics, including the load balancer name, namespace, metric name, and the timestamp of the metric.

The `aws_ec2_application_load_balancer_metric_request_count_daily` table provides you with metric statistics at 24 hour intervals for the most recent 1 year.

## Examples

### Basic info
Explore the performance of your AWS EC2 application load balancers by analyzing daily metrics. This can help you identify patterns, track changes over time, and optimize your load balancing strategy for improved efficiency and performance.

```sql+postgres
select
  name,
  metric_name,
  namespace,
  average,
  maximum,
  minimum,
  sample_count,
  timestamp
from
  aws_ec2_application_load_balancer_metric_request_count_daily
order by
  name,
  timestamp;
```

```sql+sqlite
select
  name,
  metric_name,
  namespace,
  average,
  maximum,
  minimum,
  sample_count,
  timestamp
from
  aws_ec2_application_load_balancer_metric_request_count_daily
order by
  name,
  timestamp;
```

### Intervals averaging less than 100 request count
Identify instances where the average daily request count on your AWS EC2 application load balancer is less than 100. This can help you monitor and manage your load balancer's performance and efficiency.

```sql+postgres
select
  name,
  metric_name,
  namespace,
  maximum,
  minimum,
  average
  sample_count,
  timestamp
from
  aws_ec2_application_load_balancer_metric_request_count_daily
where
  average < 100
order by
  name,
  timestamp;
```

```sql+sqlite
select
  name,
  metric_name,
  namespace,
  maximum,
  minimum,
  average,
  sample_count,
  timestamp
from
  aws_ec2_application_load_balancer_metric_request_count_daily
where
  average < 100
order by
  name,
  timestamp;
```