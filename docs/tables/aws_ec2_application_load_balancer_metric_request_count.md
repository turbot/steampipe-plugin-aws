---
title: "Steampipe Table: aws_ec2_application_load_balancer_metric_request_count - Query AWS EC2 Application Load Balancer Metrics using SQL"
description: "Allows users to query AWS EC2 Application Load Balancer Metrics, specifically the request count."
folder: "EC2"
---

# Table: aws_ec2_application_load_balancer_metric_request_count - Query AWS EC2 Application Load Balancer Metrics using SQL

The AWS EC2 Application Load Balancer is a component of the Elastic Load Balancing service that automatically distributes incoming application traffic across multiple targets, such as Amazon EC2 instances, containers, IP addresses, and Lambda functions. It can handle the varying load of your application traffic in a single Availability Zone or across multiple Availability Zones. The Application Load Balancer offers two types of load balancer that cater to specific needs - Classic Load Balancer (CLB) for simple load balancing across multiple Amazon EC2 instances and Application Load Balancer (ALB) for applications needing advanced routing capabilities, microservices, and container-based architectures.

## Table Usage Guide

The `aws_ec2_application_load_balancer_metric_request_count` table in Steampipe provides you with information about the request count metrics of Application Load Balancers within Amazon Elastic Compute Cloud (EC2). This table allows you as a DevOps engineer, system administrator, or other technical professional to query specific details about the number of requests processed by your Application Load Balancers. You can utilize this table to gather insights on load balancing performance and to monitor the traffic your applications are receiving. The schema outlines the various attributes of the request count metric, including the load balancer name, namespace, metric name, and dimensions.

The `aws_ec2_application_load_balancer_metric_request_count` table provides you with metric statistics at 5 min intervals for the most recent 5 days.

## Examples

### Basic info
Explore the performance of your application load balancers on AWS EC2 by analyzing metrics such as average, maximum, and minimum request counts. This allows you to assess the load on your balancers and make informed decisions about scaling and resource allocation.

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
  aws_ec2_application_load_balancer_metric_request_count
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
  aws_ec2_application_load_balancer_metric_request_count
order by
  name,
  timestamp;
```

### Intervals averaging less than 100 net flow count
Gain insights into application load balancer metrics where the average request count is less than 100, to understand the performance and traffic patterns. This can be beneficial for optimizing resource allocation and managing load effectively.

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
  aws_ec2_application_load_balancer_metric_request_count
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
  aws_ec2_application_load_balancer_metric_request_count
where
  average < 100
order by
  name,
  timestamp;
```