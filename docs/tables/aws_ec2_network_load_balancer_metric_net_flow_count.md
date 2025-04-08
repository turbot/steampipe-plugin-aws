---
title: "Steampipe Table: aws_ec2_network_load_balancer_metric_net_flow_count - Query AWS EC2 Network Load Balancer Metrics using SQL"
description: "Allows users to query AWS EC2 Network Load Balancer Metrics for net flow count data. This includes information such as the number of new or terminated flows per minute from a network load balancer."
folder: "ELB"
---

# Table: aws_ec2_network_load_balancer_metric_net_flow_count - Query AWS EC2 Network Load Balancer Metrics using SQL

The AWS EC2 Network Load Balancer is a high-performance load balancer that operates at the transport layer (Layer 4). It is designed to handle volatile traffic patterns and millions of requests per second for your applications. It can automatically scale to meet the needs of your applications, and you can enable cross-zone load balancing to distribute traffic evenly across all registered instances in all enabled Availability Zones.

## Table Usage Guide

The `aws_ec2_network_load_balancer_metric_net_flow_count` table in Steampipe provides you with information about the net flow count metrics of AWS EC2 Network Load Balancers. This table allows you, as a DevOps engineer, to query net flow count-specific details, including the number of new or terminated flows per minute. You can utilize this table to gather insights on network load balancing, such as monitoring the amount of traffic processed by your load balancer, identifying trends in network traffic, and more. The schema outlines the various attributes of the net flow count metric, including the load balancer name, namespace, metric name, and dimensions.

The `aws_ec2_network_load_balancer_metric_net_flow_count` table provides you with metric statistics at 5 min intervals for the most recent 5 days.

## Examples

### Basic info
Analyze the metrics of your AWS EC2 network load balancer to understand its performance over time. This will help you identify instances where the load balance may be skewed or inefficient, allowing for timely adjustments and improved resource management.

```sql+postgres
select
  name,
  metric_name,
  namespace,
  maximum,
  minimum,
  sample_count,
  timestamp
from
  aws_ec2_network_load_balancer_metric_net_flow_count
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
  sample_count,
  timestamp
from
  aws_ec2_network_load_balancer_metric_net_flow_count
order by
  name,
  timestamp;
```

### Intervals where net flow count < 100
Explore instances where the average network load balance metric net flow count is less than 100, which can be useful in identifying periods of low network traffic for AWS EC2 instances. This can be beneficial in optimizing resource allocation and understanding usage patterns.

```sql+postgres
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
  aws_ec2_network_load_balancer_metric_net_flow_count
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
  aws_ec2_network_load_balancer_metric_net_flow_count
where
  average < 100
order by
  name,
  timestamp;
```